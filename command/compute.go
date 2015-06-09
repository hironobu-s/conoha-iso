package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Compute struct {
	Identity         *Identity
	computeEndpoints map[int]string

	*Command
}

type ISOImage struct {
	Name  string
	Url   string
	Path  string
	Ctime string
	Size  int64
}

type ISOImages struct {
	IsoImages []*ISOImage `json:"iso-images"`
}

func NewCompute() *Compute {
	compute := &Compute{}

	compute.computeEndpoints = map[int]string{
		TYO1: "https://compute.tyo1.conoha.io/v2/",
		SIN1: "https://compute.sin1.conoha.io/v2/",
		SJC1: "https://compute.sjc1.conoha.io/v2/",
	}
	return compute
}

func (cmd *Compute) List() (*ISOImages, error) {

	endpoint, ok := cmd.computeEndpoints[cmd.Identity.Region]

	if !ok {
		return nil, fmt.Errorf("Undefined region \"%s\"", cmd.Identity.Region)
	}

	req, err := http.NewRequest(
		"GET",
		endpoint+"/"+cmd.Identity.ApiTenantId+"/iso-images",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", cmd.Identity.Token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 400:
		msg := cmd.extractApiErrorMessage(resp.Body)
		return nil, fmt.Errorf("Return %d status code from the server. [%s]", resp.StatusCode, msg)
	}

	res, _ := ioutil.ReadAll(resp.Body)

	var isos *ISOImages
	if err = json.Unmarshal(res, &isos); err != nil {
		return nil, err
	}

	return isos, nil
}

func (cmd *Compute) Download(url string) error {
	var err error

	reqjson := map[string]interface{}{
		"iso-image": map[string]interface{}{
			"url": url,
		},
	}

	b, err := json.Marshal(reqjson)
	if err != nil {
		return err
	}

	endpoint, ok := cmd.computeEndpoints[cmd.Identity.Region]
	if !ok {
		return fmt.Errorf("Undefined region \"%s\"", cmd.Identity.Region)
	}

	req, err := http.NewRequest(
		"POST",
		endpoint+"/"+cmd.Identity.ApiTenantId+"/iso-images",
		strings.NewReader(string(b)),
	)

	req.Header.Set("X-Auth-Token", cmd.Identity.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 400:
		msg := cmd.extractApiErrorMessage(resp.Body)
		return fmt.Errorf("Return %d status code from the server. [%s]", resp.StatusCode, msg)
	}

	return nil
}
