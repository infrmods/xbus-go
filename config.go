package xbus

import (
	"net/url"
	"strconv"
	"strings"
)

type ConfigItem struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Version int64  `json:"version"`
}

type ConfigGetResp struct {
	Response
	Result *struct {
		Config   *ConfigItem `json:"config"`
		Revision int64       `json:"revision"`
	} `json:"result,omitempty"`
}

func (cli *Client) GetConfig(key string) (*ConfigItem, error) {
	var rep ConfigGetResp
	if err := cli.request("GET", "/api/configs/"+key, nil, &rep); err != nil {
		return nil, err
	}
	return rep.Result.Config, nil
}

type ConfigPutResp struct {
	Response
	Result *struct {
		Revision int64 `json:"revision"`
	} `json:"result,omitempty"`
}

func (cli *Client) PutConfig(key, value string, revision int64) error {
	vals := make(url.Values)
	vals.Set("value", value)
	if revision != 0 {
		vals.Set("revision", strconv.FormatInt(revision, 10))
	}

	var rep ConfigPutResp
	if err := cli.request("PUT", "/api/configs/"+key, strings.NewReader(vals.Encode()), &rep); err != nil {
		return err
	}
	return nil
}

func (cli *Client) WatchConfig(key string, revision int64, timeout int64) (*ConfigItem, error) {
	vals := make(url.Values)
	vals.Set("watch", "true")
	if revision > 0 {
		vals.Set("revision", strconv.FormatInt(revision, 10))
	}
	if timeout > 0 {
		vals.Set("timeout", strconv.FormatInt(timeout, 10))
	}

	var rep ConfigGetResp
	if err := cli.request("GET", "/api/configs/"+key+"?"+vals.Encode(), nil, &rep); err == nil {
		return rep.Result.Config, nil
	} else {
		return nil, err
	}
}
