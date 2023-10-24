package alias

import (
	"encoding/json"
	"errors"
	"github.com/raisultan/url-shortener/services/main/internal/config"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
	Alias  string `json:"alias"`
	Error  string `json:"error,omitempty"`
}

type Client struct {
	baseURL string
	client  *http.Client
}

func NewAliasGeneratorClient(cfg config.AliasGenerator) *Client {
	return &Client{
		baseURL: cfg.Address,
		client:  &http.Client{Timeout: cfg.Timeout},
	}
}

func (agc *Client) GenerateAlias() (string, error) {
	resp, err := agc.client.Get(agc.baseURL + "/alias")
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	var aliasResp Response
	err = json.NewDecoder(resp.Body).Decode(&aliasResp)
	if err != nil {
		return "", err
	}

	if aliasResp.Status == "Error" {
		return "", errors.New(aliasResp.Error)
	}

	return aliasResp.Alias, nil
}
