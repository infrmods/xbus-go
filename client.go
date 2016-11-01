package xbus

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type XBusConfig struct {
	Endpoint   string
	CACertFile string
	CertFile   string
	KeyFile    string
}

type Client struct {
	cli    *http.Client
	config XBusConfig
}

func readFile(path string) ([]byte, error) {
	if f, err := os.Open(path); err == nil {
		defer f.Close()
		return ioutil.ReadAll(f)
	} else {
		return nil, err
	}
}

func NewClient(config XBusConfig) (*Client, error) {
	tlsConfig := &tls.Config{}
	if config.CACertFile != "" {
		if data, err := readFile(config.CACertFile); err == nil {
			pool := x509.NewCertPool()
			if !pool.AppendCertsFromPEM(data) {
				return nil, fmt.Errorf("add cacert fail")
			}
			tlsConfig.RootCAs = pool
		}
	}
	if config.CertFile != "" {
		if config.KeyFile == "" {
			return nil, fmt.Errorf("missing key file")
		}
		if cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile); err == nil {
			tlsConfig.Certificates = []tls.Certificate{cert}
		} else {
			return nil, err
		}
	} else if config.KeyFile != "" {
		return nil, fmt.Errorf("missing cert file")
	}
	httpCli := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	cli := Client{config: config, cli: &httpCli}
	return &cli, nil
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Ok    bool   `json:"ok"`
	Error *Error `json:"error,omitempty"`
}

func (cli *Client) urlFor(path string) string {
	return cli.config.Endpoint + path
}

func (cli *Client) request(method, path string, body io.Reader, v interface{}) error {
	req, err := http.NewRequest(method, cli.urlFor(path), body)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := cli.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		return fmt.Errorf("resquest fail: %s", string(data))
	} else {
		return err
	}
}
