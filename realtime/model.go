package realtime

type StreamResponse struct {
	Code    int      `json:"code"`
	Server  string   `json:"server"`
	Service string   `json:"service"`
	Pid     string   `json:"pid"`
	Streams []Stream `json:"streams"`
}

type Stream struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	VHost     string    `json:"vhost"`
	App       string    `json:"app"`
	TcURL     string    `json:"tcUrl"`
	URL       string    `json:"url"`
	LiveMS    int64     `json:"live_ms"`
	Clients   int       `json:"clients"`
	Frames    int       `json:"frames"`
	SendBytes int64     `json:"send_bytes"`
	RecvBytes int64     `json:"recv_bytes"`
	Kbps      Kbps      `json:"kbps"`
	Publish   Publish   `json:"publish"`
	Video     VideoInfo `json:"video"`
	Audio     AudioInfo `json:"audio"`
}

type Kbps struct {
	Recv30s int `json:"recv_30s"`
	Send30s int `json:"send_30s"`
}

type Publish struct {
	Active bool   `json:"active"`
	CID    string `json:"cid"`
}

type VideoInfo struct {
	Codec   string `json:"codec"`
	Profile string `json:"profile"`
	Level   string `json:"level"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

type AudioInfo struct {
	Codec      string `json:"codec"`
	SampleRate int    `json:"sample_rate"`
	Channel    int    `json:"channel"`
	Profile    string `json:"profile"`
}
