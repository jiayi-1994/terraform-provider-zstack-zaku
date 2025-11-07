package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ZeAuthProviderTransport struct {
	Ak          string
	Sk          string
	ContextPath string

	http.RoundTripper
}

func (p *ZeAuthProviderTransport) WrapTransport(rt http.RoundTripper) http.RoundTripper {
	p.RoundTripper = rt
	return p
}

// RoundTrip implements the http.RoundTripper interface
func (t *ZeAuthProviderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	date := time.Now().Format(time.RFC1123Z)
	hmac := hmac.New(sha1.New, []byte(t.Sk))
	uri := strings.Replace(req.URL.Path, t.ContextPath, "", 1)
	hmac.Write([]byte(fmt.Sprintf("%s\n%s\n%s", req.Method, date, uri)))
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	req.Header.Set("Authorization", fmt.Sprintf("Zstack %s:%s", t.Ak, signature))
	req.Header.Set("Date", date)
	fmt.Printf("%v\n", req.Header)
	return t.RoundTripper.RoundTrip(req)
}
