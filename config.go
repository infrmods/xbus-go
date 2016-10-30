package xbus

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
	if err := cli.simpleReq("GET", "/api/configs/"+key, &rep); err != nil {
		return nil, err
	}
	return rep.Result.Config, nil
}
