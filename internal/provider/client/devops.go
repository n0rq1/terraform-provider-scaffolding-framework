package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetDevOps() ([]DevOps, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devops", c.endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var devops []DevOps
	if err := json.Unmarshal(body, &devops); err != nil {
		return nil, err
	}
	return devops, nil
}

func (c *Client) CreateDevops(devops DevOps) (*DevOps, error) {
	b, err := json.Marshal(devops)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devops", c.endpoint), strings.NewReader(string(b)))
	
	req.Header.Set("Content-Type", "application/json")
	
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var created DevOps
	if err := json.Unmarshal(body, &created); err != nil {
		return nil, err
	}
	return &created, nil
}