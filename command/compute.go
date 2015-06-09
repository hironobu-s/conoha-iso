package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

type ServerId string

type Servers struct {
	Servers []struct {
		Name string
		Id   ServerId
	}
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

func (cmd *Compute) Insert() error {
	serverId := cmd.selectVps()
	if serverId == "" {
		return fmt.Errorf("Can't detect the server")
	}

	println()

	iso := cmd.selectIso()
	if iso == nil {
		return fmt.Errorf("Can't detect ISO Image ")
	}

	reqjson := map[string]interface{}{
		"mountImage": iso.Path,
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
		endpoint+"/"+cmd.Identity.ApiTenantId+"/servers/"+string(serverId)+"/action",
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

func (cmd *Compute) Eject() error {
	serverId := cmd.selectVps()
	if serverId == "" {
		return fmt.Errorf("Can't detect the server")
	}

	reqjson := map[string]interface{}{
		"unmountImage": "",
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
		endpoint+"/"+cmd.Identity.ApiTenantId+"/servers/"+string(serverId)+"/action",
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

func (cmd *Compute) serverList() (*Servers, error) {
	endpoint, ok := cmd.computeEndpoints[cmd.Identity.Region]
	if !ok {
		return nil, fmt.Errorf("Undefined region \"%s\"", cmd.Identity.Region)
	}

	req, err := http.NewRequest(
		"GET",
		endpoint+"/"+cmd.Identity.ApiTenantId+"/servers",
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
	r, _ := ioutil.ReadAll(resp.Body)

	// -------------

	println(string(r))

	var servers *Servers
	if err = json.Unmarshal(r, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func (cmd *Compute) selectVps() ServerId {
	servers, err := cmd.serverList()
	if err != nil {
		return ""
	} else if len(servers.Servers) == 0 {
		return ""
	}

	var i int
	for i, vps := range servers.Servers {
		fmt.Printf("[%d] %s\n", i+1, vps.Name)
		i++
	}

	if len(servers.Servers) == 1 {
		fmt.Printf("Please select VPS no. [1]: ")
	} else {
		fmt.Printf("Please select VPS no. [1-%d]: ", len(servers.Servers))
	}

	var no string
	if _, err = fmt.Scanf("%s", &no); err != nil {
		return ""
	}

	i, err = strconv.Atoi(no)
	if err != nil {
		return ""

	} else if 1 <= i && i <= len(servers.Servers) {
		return servers.Servers[i-1].Id

	} else {
		return ""
	}
}

func (cmd *Compute) selectIso() *ISOImage {
	isos, err := cmd.List()
	if err != nil {
		return nil
	} else if len(isos.IsoImages) == 0 {
		return nil
	}

	var i int
	for i, iso := range isos.IsoImages {
		fmt.Printf("[%d] %s\n", i+1, iso.Name)
		i++
	}

	if len(isos.IsoImages) == 1 {
		fmt.Printf("Please select ISO no. [1]: ")
	} else {
		fmt.Printf("Please select ISO no. [1-%d]: ", len(isos.IsoImages))
	}

	var no string
	if _, err = fmt.Scanf("%s", &no); err != nil {
		return nil
	}

	i, err = strconv.Atoi(no)
	if err != nil {
		return nil

	} else if 1 <= i && i <= len(isos.IsoImages) {
		return isos.IsoImages[i-1]

	} else {
		return nil
	}
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
