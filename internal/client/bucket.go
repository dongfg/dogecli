package client

import (
    "errors"
    "github.com/dustin/go-humanize"
    "time"
)

type Bucket struct {
    Name       string `json:"name"`
    CTime      int    `json:"ctime"`
    CTimeStr   string
    Region     int `json:"region"`
    RegionName string
    Space      uint64 `json:"space"`
    SpaceHuman string
}

type bucketListResponse struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data struct {
        Buckets []Bucket `json:"buckets"`
    } `json:"data"`
}

var bucketRegionMap map[int]string

func init() {
    bucketRegionMap = make(map[int]string)
    bucketRegionMap[0] = "上海（华东）"
    bucketRegionMap[1] = "北京（华北）"
    bucketRegionMap[2] = "广州（华南）"
    bucketRegionMap[3] = "成都（西南）"

}

func (c *Client) BucketList() ([]Bucket, error) {
    var response bucketListResponse
    err := c.httpGET("/oss/bucket/list.json", &response)
    if err != nil {
        return nil, err
    }
    if response.Code != 200 {
        return nil, errors.New(response.Msg)
    }
    for i := range response.Data.Buckets {
        bucket := &response.Data.Buckets[i]
        bucket.SpaceHuman = humanize.Bytes(bucket.Space)
        bucket.CTimeStr = humanize.Time(time.Unix(int64(bucket.CTime), 0))
        name, ok := bucketRegionMap[bucket.Region]
        if ok {
            bucket.RegionName = name
        } else {
            bucket.RegionName = "未知"
        }
    }
    return response.Data.Buckets, nil
}
