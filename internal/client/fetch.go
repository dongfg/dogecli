package client

import (
	"errors"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

type fileFetchResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ID   string `json:"id"`
		Wait int    `json:"wait"`
	} `json:"data"`
}

type fileFetchState struct {
	ID     string `json:"id"     yaml:"id"`
	Bucket string `json:"bucket" yaml:"bucket"`
	Name   string `json:"name"   yaml:"name"`
	CTime  string `json:"ctime"  yaml:"ctime"`
	Wait   int    `json:"wait"   yaml:"wait"`
	State  string `json:"state"  yaml:"state"`
}

func (c *Client) FileFetch(remoteUrl string, bucket, name string) error {
	var response fileFetchResponse
	body := make(map[string]interface{})
	body["url"] = remoteUrl
	body["bucket"] = bucket
	body["key"] = name

	err := c.httpPostJson("/oss/fetch.json", body, &response)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return errors.New(response.Msg)
	}
	return addFetch(fileFetchState{
		ID:     response.Data.ID,
		Bucket: bucket,
		Name:   name,
		State:  waitToState(response.Data.Wait),
	})
}

func (c *Client) FileFetchList() ([]fileFetchState, error) {
	fetches, err := readAllFetches()
	if err != nil {
		return nil, err
	}
	for _, f := range fetches {
		if f.Wait != -1 {
			_ = c.FileFetchState(f.ID)
		}
	}
	return readAllFetches()
}

func (c *Client) FileFetchState(id string) error {
	var response fileFetchResponse
	params := url.Values{}
	params.Set("id", id)

	err := c.httpGetJson("/oss/fetch/query.json?"+params.Encode(), &response)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return errors.New(response.Msg)
	}
	return updateFetch(id, waitToState(response.Data.Wait))
}

func readAllFetches() ([]fileFetchState, error) {
	var fetches []fileFetchState
	err := viper.UnmarshalKey("fetches", &fetches)
	return fetches, err
}

func writeFetches(fetches []fileFetchState) error {
	viper.Set("fetches", fetches)
	return viper.WriteConfigAs(viper.ConfigFileUsed())
}

func addFetch(fetch fileFetchState) error {
	fetches, err := readAllFetches()
	if err != nil {
		return err
	}
	fetch.CTime = time.Now().Format("2006-01-02 15:04:05")
	fetches = append(fetches, fetch)
	return writeFetches(fetches)
}

func updateFetch(id string, state string) error {
	fetches, err := readAllFetches()
	if err != nil {
		return err
	}
	for i, fetch := range fetches {
		if fetch.ID == id {
			fetches[i].State = state
			break
		}
	}
	return writeFetches(fetches)
}

func waitToState(wait int) string {
	if wait == 0 {
		return "Fetching"
	}
	if wait == -1 {
		return "Finish"
	}
	return "Unknown"
}
