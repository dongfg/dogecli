package client

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dongfg/dogecli/internal/constants"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Client wraps HTTP and credentials
type Client struct {
	httpClient *http.Client
	accessKey  string
	secretKey  string
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		accessKey:  viper.GetString(constants.EnvAccessKey),
		secretKey:  viper.GetString(constants.EnvSecretKey),
	}
}

// accessToken of request
func (c *Client) accessToken(apiPath string, data []byte) string {
	body := ""
	if data != nil {
		body = string(data)
	}
	signStr := apiPath + "\n" + body
	hmacObj := hmac.New(sha1.New, []byte(c.secretKey))
	hmacObj.Write([]byte(signStr))
	sign := hex.EncodeToString(hmacObj.Sum(nil))
	return "TOKEN " + c.accessKey + ":" + sign
}

func (c *Client) httpGetJson(url string, out interface{}) error {
	logrus.Debug("request url: ", url)
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
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	logrus.Debug("response body: ", string(raw))
	// 解码 JSON
	if err := json.Unmarshal(raw, &out); err != nil && err != io.EOF {
		return fmt.Errorf("decoding json: %w", err)
	}
	return nil
}

func (c *Client) httpPostJson(url string, data map[string]interface{}, out interface{}) error {
	logrus.Debug("request url: ", url)
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, constants.ApiHost+url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.accessToken(url, body))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("doing http GET: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	logrus.Debug("response body: ", string(raw))
	// 解码 JSON
	if err := json.Unmarshal(raw, &out); err != nil && err != io.EOF {
		return fmt.Errorf("decoding json: %w", err)
	}
	return nil
}

func (c *Client) fileUpload(url string, filePath string, out interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// 获取文件大小（可选，用于进度条总量显示）
	fi, _ := file.Stat()
	total := fi.Size()
	if total > 10*1024*1024 { // 10MB
		return fmt.Errorf("文件不能超过10M")
	}
	bar := progressbar.NewOptions(
		int(total),
		progressbar.OptionSetDescription("Uploading..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
	)
	// 使用带缓冲的Reader包装文件，减少IO次数
	bufFile := bufio.NewReaderSize(file, 1024*1024) // 1MB缓冲区

	// 创建一个带缓冲的TeeReader，避免进度条阻塞
	// 使用管道异步处理进度更新，防止阻塞主读取流程
	pr, pw := io.Pipe()

	// 启动goroutine处理进度更新
	go func() {
		defer func() {
			_ = pw.Close()
		}()
		// 将数据从缓冲文件复制到管道，同时更新进度条
		_, err := io.Copy(pw, io.TeeReader(bufFile, bar))
		if err != nil && !errors.Is(err, io.ErrClosedPipe) {
			logrus.Warnf("copy data with progress: %v", err)
		}
	}()
	logrus.Debug("request url: ", constants.ApiHost+url)
	req, err := http.NewRequest(http.MethodPut, constants.ApiHost+url, pr)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", c.accessToken(url, nil))
	req.ContentLength = total

	resp, err := c.httpClient.Do(req)
	fmt.Println()
	if err != nil {
		logrus.Debug("request header: ", req.Header)
		return fmt.Errorf("doing http PUT: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	logrus.Debug("response body: ", string(raw))
	// 解码 JSON
	if err := json.Unmarshal(raw, &out); err != nil && err != io.EOF {
		return fmt.Errorf("decoding json: %w", err)
	}
	return nil
}
