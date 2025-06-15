package utils

import "net/url"

func ExtractPass(rtmpUrl string) (string, error) {
	u, err := url.Parse(rtmpUrl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	return q.Get("pass"), nil
}
