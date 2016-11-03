package xbus

import (
	"context"
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

func (cli *Client) GetConfig(ctx context.Context, key string) (*ConfigItem, error) {
	var rep ConfigGetResp
	if err := cli.request(ctx, "GET", "/api/configs/"+key, nil, &rep); err != nil {
		return nil, err
	}
	if !rep.Ok {
		return nil, rep.Error
	}
	return rep.Result.Config, nil
}

type ConfigPutResp struct {
	Response
	Result *struct {
		Revision int64 `json:"revision"`
	} `json:"result,omitempty"`
}

func (cli *Client) PutConfig(ctx context.Context, key, value string, revision int64) error {
	vals := make(url.Values)
	vals.Set("value", value)
	if revision != 0 {
		vals.Set("revision", strconv.FormatInt(revision, 10))
	}

	var rep ConfigPutResp
	if err := cli.request(ctx, "PUT", "/api/configs/"+key, strings.NewReader(vals.Encode()), &rep); err != nil {
		return err
	}
	return rep.Error
}

func (cli *Client) WatchConfig(ctx context.Context, key string, revision int64, timeout int64) (*ConfigItem, error) {
	vals := make(url.Values)
	vals.Set("watch", "true")
	if revision > 0 {
		vals.Set("revision", strconv.FormatInt(revision, 10))
	}
	if timeout > 0 {
		vals.Set("timeout", strconv.FormatInt(timeout, 10))
	}

	var rep ConfigGetResp
	if err := cli.request(ctx, "GET", "/api/configs/"+key+"?"+vals.Encode(), nil, &rep); err != nil {
		return nil, err
	}
	if !rep.Ok {
		return nil, rep.Error
	}
	return rep.Result.Config, nil
}
