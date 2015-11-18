package command

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Compute struct {
	identity *Identity
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

type Servers struct {
	Servers []*Server
}

type Server struct {
	Id   string
	Name string
}

func NewCompute(ident *Identity) *Compute {
	compute := &Compute{
		identity: ident,
	}

	return compute
}

func (cmd *Compute) newApi() (api *Api, err error) {
	api, err = NewApi("compute", cmd.identity.Region)
	if err != nil {
		return nil, err
	}
	api.Token = cmd.identity.Token
	api.TenantId = cmd.identity.ApiTenantId

	return api, err
}

func (cmd *Compute) Insert() error {
	server, err := cmd.selectVps()
	if err != nil {
		return err
	} else if server == nil {
		return fmt.Errorf("Can't detect the server")
	}

	println()

	iso, err := cmd.selectIso()
	if err != nil {
		return err
	} else if iso == nil {
		return fmt.Errorf("No ISO Images.")
	}

	reqjson := map[string]interface{}{
		"mountImage": iso.Path,
	}

	b, err := json.Marshal(reqjson)
	if err != nil {
		return err
	}

	api, err := cmd.newApi()
	if err != nil {
		return err
	}

	if err = api.Prepare("POST", []string{"servers", server.Id, "action"}, b); err != nil {
		return err
	}

	ch := api.Do()
	_ = <-ch

	if err = api.LastError(); err != nil {
		return err
	}

	return nil
}

func (cmd *Compute) Eject() (err error) {
	server, err := cmd.selectVps()
	if err != nil {
		return err
	} else if server == nil {
		return fmt.Errorf("Can't detect the server")
	}

	reqjson := map[string]interface{}{
		"unmountImage": "",
	}

	b, err := json.Marshal(reqjson)
	if err != nil {
		return err
	}

	api, err := cmd.newApi()
	if err != nil {
		return err
	}

	if err = api.Prepare("POST", []string{"servers", server.Id, "action"}, b); err != nil {
		return err
	}

	ch := api.Do()
	_ = <-ch

	if err = api.LastError(); err != nil {
		return err
	}

	return nil
}

func (cmd *Compute) List() (isos *ISOImages, err error) {
	api, err := cmd.newApi()
	if err != nil {
		return nil, err
	}

	if err = api.Prepare("GET", []string{"iso-images"}, nil); err != nil {
		return nil, err
	}

	ch := api.Do()

	res := <-ch
	if err = api.LastError(); err != nil {
		return nil, err
	}

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

	api, err := cmd.newApi()
	if err != nil {
		return err
	}

	if err = api.Prepare("POST", []string{"iso-images"}, b); err != nil {
		return err
	}

	ch := api.Do()

	_ = <-ch
	if err = api.LastError(); err != nil {
		return err
	}

	return nil
}

// -----------------------

func (cmd *Compute) serverList() (servers *Servers, err error) {
	api, err := cmd.newApi()
	if err != nil {
		return nil, err
	}

	if err = api.Prepare("GET", []string{"servers"}, nil); err != nil {
		return nil, err
	}

	ch := api.Do()
	resp := <-ch

	if err = api.LastError(); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func (cmd *Compute) selectVps() (*Server, error) {
	servers, err := cmd.serverList()
	if err != nil {
		return nil, err
	} else if len(servers.Servers) == 0 {
		return nil, fmt.Errorf("No servers found.")
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
		return nil, err
	}

	i, err = strconv.Atoi(no)
	if err != nil {
		return nil, err

	} else if 1 <= i && i <= len(servers.Servers) {
		return servers.Servers[i-1], nil

	} else {
		return nil, fmt.Errorf("Wrong VPS no.")
	}
}

func (cmd *Compute) selectIso() (*ISOImage, error) {
	isos, err := cmd.List()
	if err != nil {
		return nil, err
	} else if len(isos.IsoImages) == 0 {
		return nil, fmt.Errorf("No iso images found.")
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
		return nil, err
	}

	i, err = strconv.Atoi(no)
	if err != nil {
		return nil, err

	} else if 1 <= i && i <= len(isos.IsoImages) {
		return isos.IsoImages[i-1], nil

	} else {
		return nil, fmt.Errorf("Wrong ISO no.")
	}
}
