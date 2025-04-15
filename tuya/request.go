package tuya

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	headerClientId    = "client_id"
	headerSignMethod  = "sign_method"
	headerTimestamp   = "t"
	headerAccessToken = "access_token"
	headerSign        = "sign"
	signMethod        = "HMAC-SHA256"
)

func (c *TuyaClient) performRequest(method string, url string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest(method, c.host+url, bytes.NewReader(body))
	c.buildHeader(req, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// @see https://developer.tuya.com/en/docs/iot/new-singnature?id=Kbw0q34cs2e5g
func (c *TuyaClient) buildHeader(req *http.Request, body []byte) {
	req.Header.Set(headerClientId, c.clientId)
	req.Header.Set(headerSignMethod, signMethod)
	ts := fmt.Sprint(time.Now().UnixNano() / 1e6)
	req.Header.Set(headerTimestamp, ts)
	if c.token != "" {
		req.Header.Set(headerAccessToken, c.token)
	}
	sign := c.buildSign(req, body, ts)
	req.Header.Set(headerSign, sign)
}

func (c *TuyaClient) buildSign(req *http.Request, body []byte, timestamp string) string {
	headers := c.getHeaderStr(req)
	urlStr := c.getUrlStr(req)
	contentSha256 := c.sha256(body)
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s", req.Method, contentSha256, headers, urlStr)
	signStr := fmt.Sprintf("%s%s%s%s", c.clientId, c.token, timestamp, stringToSign)
	sign := strings.ToUpper(c.hmacSha256(signStr, c.secret))
	return sign
}

func (c *TuyaClient) sha256(data []byte) string {
	sha256Contain := sha256.New()
	sha256Contain.Write(data)
	return hex.EncodeToString(sha256Contain.Sum(nil))
}

func (c *TuyaClient) getUrlStr(req *http.Request) string {
	url := req.URL.Path
	keys := make([]string, 0, 10)
	for key, _ := range req.URL.Query() {
		keys = append(keys, fmt.Sprintf("%s=%s", key, req.URL.Query().Get(key)))
	}
	if len(keys) > 0 {
		sort.Strings(keys)
		keysPart := strings.Join(keys, "&")
		return fmt.Sprintf("%s?%s", url, keysPart)
	} else {
		return url
	}
}

func (c *TuyaClient) getHeaderStr(req *http.Request) string {
	signHeaderKeys := req.Header.Get("Signature-Headers")
	if signHeaderKeys == "" {
		return ""
	}
	keys := strings.Split(signHeaderKeys, ":")
	headers := make([]string, 0, len(keys))
	for _, key := range keys {
		headers = append(headers, fmt.Sprintf("%s:%s", key, req.Header.Get(key)))
	}
	return strings.Join(headers, "\n")
}

func (c *TuyaClient) hmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
