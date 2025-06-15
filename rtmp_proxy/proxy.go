// internal/rtmp/proxy.go
package rtmp_proxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go-postgres-gorm-gin-api/models"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type Proxy struct {
	listenAddr string
	targetAddr string
	ln         net.Listener
	wg         sync.WaitGroup
	quit       chan struct{}
	DB         *gorm.DB
}

func NewProxy(listenAddr, targetAddr string, db *gorm.DB) *Proxy {
	return &Proxy{
		listenAddr: listenAddr,
		targetAddr: targetAddr,
		quit:       make(chan struct{}),
		DB:         db,
	}
}

func (p *Proxy) Start() error {
	ln, err := net.Listen("tcp", p.listenAddr)
	if err != nil {
		return err
	}
	p.ln = ln

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-p.quit:
					return
				default:
					log.Println("accept error:", err)
				}
				continue
			}

			p.wg.Add(1)
			go p.handleConn(conn)
		}
	}()

	log.Println("rtmp proxy started on", p.listenAddr)
	return nil
}

func (p *Proxy) handleConn(client net.Conn) {
	defer p.wg.Done()
	defer client.Close()

	server, err := net.Dial("tcp", p.targetAddr)
	if err != nil {
		log.Println("dial target failed:", err)
		return
	}
	defer server.Close()

	if err := proxyHandshake(client, server); err != nil {
		log.Println("handshake failed:", err)
		return
	}

	// Создаём буфер для анализа первых данных
	var sniffBuf bytes.Buffer
	tee := io.TeeReader(client, &sniffBuf)

	// Считываем часть трафика, не нарушая поток
	firstChunk := make([]byte, 4096)
	_, err = tee.Read(firstChunk)
	if err != nil {
		log.Println("tee read failed:", err)
		return
	}

	// анализируем буфер
	buf := sniffBuf.Bytes()
	tcUrl := extractTcUrl(buf)

	room, pass, err := extractInfo(tcUrl)
	if err != nil {
		return
	}

	ok := p.checkPass(room, pass)

	if !ok {
		return
	}

	// передаём всё что накопили на сервер
	if _, err := server.Write(buf); err != nil {
		log.Println("initial server write failed:", err)
		return
	}

	// теперь запускаем полноценный io.Copy
	done := make(chan struct{})

	go func() {
		io.Copy(server, client)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(client, server)
		done <- struct{}{}
	}()

	<-done
}

func (p *Proxy) checkPass(room, pass string) bool {
	var creds models.StreamCredentials

	err := p.DB.Where("room = ?", room).First(&creds).Error
	if err != nil {
		return false
	}

	return pass == creds.Password
}

func (p *Proxy) Stop() {
	close(p.quit)
	if p.ln != nil {
		p.ln.Close()
	}
	p.wg.Wait()
	log.Println("rtmp proxy stopped")
}

// === helpers ===

func proxyHandshake(client, server net.Conn) error {
	// C0+C1
	c0c1 := make([]byte, 1537)
	if _, err := io.ReadFull(client, c0c1); err != nil {
		return err
	}
	if _, err := server.Write(c0c1); err != nil {
		return err
	}

	// S0+S1+S2
	s0s1s2 := make([]byte, 3073)
	if _, err := io.ReadFull(server, s0s1s2); err != nil {
		return err
	}
	if _, err := client.Write(s0s1s2); err != nil {
		return err
	}

	// C2
	c2 := make([]byte, 1536)
	if _, err := io.ReadFull(client, c2); err != nil {
		return err
	}
	if _, err := server.Write(c2); err != nil {
		return err
	}

	return nil
}

// Грубый парсинг tcUrl из AMF0 пакета (без полной спецификации)
func extractTcUrl(data []byte) string {
	// Ищем строку "tcUrl"
	idx := bytes.Index(data, []byte("tcUrl"))
	if idx == -1 || idx+7 >= len(data) {
		return ""
	}

	// Пропускаем "tcUrl" + 1 байт (0x02 — тип строки)
	buf := data[idx+5:]
	if len(buf) < 3 || buf[0] != 0x02 {
		return ""
	}

	strLen := binary.BigEndian.Uint16(buf[1:3])
	if int(3+strLen) > len(buf) {
		return ""
	}

	tcUrl := string(buf[3 : 3+strLen])
	tcUrl = strings.TrimSpace(tcUrl)
	return tcUrl
}

func extractInfo(rtmpUrl string) (room, pass string, err error) {
	u, err := url.Parse(rtmpUrl)
	if err != nil {
		return "", "", err
	}
	pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", "", fmt.Errorf("invalid path: %s", u.Path)
	}
	room = pathParts[1]
	pass = u.Query().Get("pass")
	return room, pass, nil
}
