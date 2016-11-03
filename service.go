package xbus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ServiceDesc struct {
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Type        string `json:"type"`
	Proto       string `json:"proto,omitempty"`
	Description string `json:"description,omitempty"`
}

type ServiceEndpoint struct {
	Address string `json:"address"`
	Config  string `json:"config,omitempty"`
}

type Service struct {
	Endpoints []ServiceEndpoint `json:"endpoints"`

	ServiceDesc
}

type serviceQueryResponse struct {
	Response
	Result *struct {
		Service  *Service `json:"service"`
		Revision int64    `json:"revision"`
	} `json:"result"`
}

func (cli *Client) GetService(ctx context.Context, name, version string) (*Service, error) {
	var rep serviceQueryResponse
	if err := cli.request(ctx, "GET", fmt.Sprintf("/api/services/%s/%s", name, version), nil, &rep); err != nil {
		return nil, err
	}
	if !rep.Ok {
		return nil, rep.Error
	}
	return rep.Result.Service, nil
}

type allServiceResponse struct {
	Response
	Result *struct {
		Services map[string]Service `json:"services"`
		Revision int64
	} `json:"result"`
}

func (cli *Client) GetAllService(ctx context.Context, name string) (map[string]Service, error) {
	var rep allServiceResponse
	if err := cli.request(ctx, "GET", fmt.Sprintf("/api/services/%s", name), nil, &rep); err != nil {
		return nil, err
	}
	if !rep.Ok {
		return nil, rep.Error
	}
	return rep.Result.Services, nil
}

func (cli *Client) WatchService(ctx context.Context, name, version string, timeout int64) (*Service, error) {
	vals := make(url.Values)
	vals.Set("watch", "true")
	//TODO: revision
	//vals.Set("revision", 0)
	if timeout > 0 {
		vals.Set("timeout", strconv.FormatInt(timeout, 10))
	}

	var rep serviceQueryResponse
	if err := cli.request(ctx, "GET", fmt.Sprintf("/api/services/%s/%s?%s", name, version, vals.Encode()), nil, &rep); err != nil {
		return nil, err
	}
	if !rep.Ok {
		return nil, rep.Error
	}

	return rep.Result.Service, nil
}

type servicePlugResponse struct {
	Response
	Result *struct {
		LeaseID int64 `json:"lease_id"`
		TTL     int64 `json:"ttl"`
	} `json:"result"`
}

func (cli *Client) PlugService(ctx context.Context,
	desc *ServiceDesc, endpoint *ServiceEndpoint, ttl, leaseId int64) (int64, error) {
	vals := make(url.Values)
	if ttl > 0 {
		vals.Set("ttl", strconv.FormatInt(ttl, 10))
	}
	if leaseId > 0 {
		vals.Set("lease_id", strconv.FormatInt(leaseId, 10))
	}
	if data, err := json.Marshal(desc); err == nil {
		vals.Set("desc", string(data))
	} else {
		return 0, err
	}
	if data, err := json.Marshal(endpoint); err == nil {
		vals.Set("endpoint", string(data))
	} else {
		return 0, err
	}
	var rep servicePlugResponse
	if err := cli.request(ctx, "POST",
		fmt.Sprintf("/api/services/%s/%s", desc.Name, desc.Version),
		strings.NewReader(vals.Encode()), &rep); err != nil {
		return 0, err
	}
	if !rep.Ok {
		return 0, rep.Error
	}
	return rep.Result.LeaseID, nil
}

func (cli *Client) PlugAllService(ctx context.Context, desces []ServiceDesc,
	endpoint *ServiceEndpoint, ttl, leaseId int64) (int64, error) {
	vals := make(url.Values)
	if ttl > 0 {
		vals.Set("ttl", strconv.FormatInt(ttl, 10))
	}
	if leaseId > 0 {
		vals.Set("lease_id", strconv.FormatInt(leaseId, 10))
	}
	if data, err := json.Marshal(desces); err == nil {
		vals.Set("desces", string(data))
	} else {
		return 0, err
	}
	if data, err := json.Marshal(endpoint); err == nil {
		vals.Set("endpoint", string(data))
	} else {
		return 0, err
	}
	var rep servicePlugResponse
	if err := cli.request(ctx, "POST",
		"/api/services", strings.NewReader(vals.Encode()), &rep); err != nil {
		return 0, err
	}
	if !rep.Ok {
		return 0, rep.Error
	}
	return rep.Result.LeaseID, nil
}

func (cli *Client) UnplugService(ctx context.Context, name, version, addr string) error {
	var rep Response
	if err := cli.request(ctx, "DELETE", fmt.Sprintf("/api/services/%s/%s/%s", name, version, addr), nil, &rep); err != nil {
		return err
	}
	return rep.Error
}

func (cli *Client) UpdateService(ctx context.Context, name, version, addr string, endpoint *ServiceEndpoint) error {
	vals := make(url.Values)
	if data, err := json.Marshal(endpoint); err == nil {
		vals.Set("endpoint", string(data))
	} else {
		return err
	}

	var rep Response
	if err := cli.request(ctx, "POST",
		fmt.Sprintf("/api/services/%s/%s/%s", name, version, addr),
		strings.NewReader(vals.Encode()), &rep); err != nil {
		return err
	}
	return rep.Error
}
