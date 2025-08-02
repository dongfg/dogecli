package client

import (
	"errors"
	"net/url"

	"github.com/dustin/go-humanize"
)

type File struct {
	Key        string `json:"key"`
	Hash       string `json:"hash"`
	FSize      uint64 `json:"fsize"`
	FSizeHuman string
	Time       string `json:"time"`
	// Type can be file/folder
	Type string `json:"type"`
}

type fileListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Files    []File `json:"files"`
		Continue string `json:"continue"` // 还有剩余文件则返回本次获取的最后一个文件名，可以作为下一次获取的参数传入
	} `json:"data"`
}

type fileUploadResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c *Client) FileList(bucket, prefix, cursor string) ([]File, string, error) {
	var response fileListResponse
	params := url.Values{}
	params.Set("bucket", bucket)
	params.Set("prefix", prefix)
	params.Set("continue", cursor)
	params.Set("limit", "10")
	err := c.httpGetJson("/oss/file/list.json?"+params.Encode(), &response)
	if err != nil {
		return nil, "", err
	}
	if response.Code != 200 {
		return nil, "", errors.New(response.Msg)
	}
	var files []File
	for i := range response.Data.Files {
		f := &response.Data.Files[i]
		if f.Type == "file" {
			f.FSizeHuman = humanize.Bytes(f.FSize)
		} else {
			f.FSizeHuman = "(folder)"
		}
		files = append(files, *f)
	}
	return files, response.Data.Continue, nil
}

func (c *Client) FileUpload(filePath string, bucket, name string) error {
	var response fileUploadResponse
	params := url.Values{}
	params.Set("bucket", bucket)
	params.Set("key", name)

	err := c.fileUpload("/oss/upload/put.json?"+params.Encode(), filePath, &response)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return errors.New(response.Msg)
	}
	return nil
}
