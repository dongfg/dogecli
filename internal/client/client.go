package client

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dongfg/dogecli/internal/constants"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

// Client wraps HTTP and credentials
type Client struct {
	httpClient *http.Client
	accessKey  string
	secretKey  string
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		accessKey:  viper.GetString(constants.EnvAccessKey),
		secretKey:  viper.GetString(constants.EnvSecretKey),
	}
}

// accessToken of request
func (c *Client) accessToken(apiPath string, data map[string]interface{}) string {
	body := ""
	if data != nil {
		bytes, _ := json.Marshal(data)
		body = string(bytes)
	}
	signStr := apiPath + "\n" + body
	hmacObj := hmac.New(sha1.New, []byte(c.secretKey))
	hmacObj.Write([]byte(signStr))
	sign := hex.EncodeToString(hmacObj.Sum(nil))
	return "TOKEN " + c.accessKey + ":" + sign
}

func (c *Client) httpGET(url string, out interface{}) error {
	req, err := http.NewRequest(http.MethodGet, constants.ApiHost+url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.accessToken(url, nil))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("doing http GET: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	// 解码 JSON
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(out); err != nil && err != io.EOF {
		return fmt.Errorf("decoding json: %w", err)
	}
	return nil
}
