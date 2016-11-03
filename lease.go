package xbus

import (
	"context"
	"fmt"
)

func (cli *Client) KeepAliveLease(ctx context.Context, leaseId int64) error {
	var rep Response
	if err := cli.request(ctx, "POST", fmt.Sprintf("/api/leases/%d", leaseId), nil, &rep); err != nil {
		return err
	}
	return rep.Error
}

func (cli *Client) RevokeLease(ctx context.Context, leaseId int64) error {
	var rep Response
	if err := cli.request(ctx, "DELETE", fmt.Sprintf("/api/leases/%d", leaseId), nil, &rep); err != nil {
		return err
	}
	return rep.Error
}
