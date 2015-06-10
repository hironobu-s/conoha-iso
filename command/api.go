package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Regions
const (
	TYO1 = "tyo1"
	SIN1 = "sin1"
	SJC1 = "sjc1"
)

type Region struct {
	Name string
}

func (r *Region) Set(region string) error {
	switch region {
	case TYO1:
		fallthrough
	case SIN1:
		fallthrough
	case SJC1:
		r.Name = region
	default:
		return fmt.Errorf("Undefined region[%s].", region)
	}
	return nil
}

// -----------

const (
	IDENTITY = "identity"
	COMPUTE  = "compute"
)

type Service struct {
	Name    string
	Version string
}

func (s *Service) Set(service string) error {

	switch service {
	case IDENTITY:
		s.Name = IDENTITY
		s.Version = "v2.0"
	case COMPUTE:
		s.Name = COMPUTE
		s.Version = "v2"
	default:
		return fmt.Errorf("Undefined service[%s].", service)
	}
	return nil
}

// -----------

type Api struct {
	Region   Region
	Service  Service
	Token    string
	TenantId string

	request          *http.Request
	lastRequestError error
}

func NewApi(service string, region string) (api *Api, err error) {
	api = &Api{}
	if err = api.Service.Set(service); err != nil {
		return nil, err
	}

	if err = api.Region.Set(region); err != nil {
		return nil, err
	}

	return api, nil
}

func (api *Api) Endpoint(path []string) (*url.URL, error) {
	var strurl string
	if api.TenantId == "" {
		strurl = fmt.Sprintf("https://%s.%s.conoha.io/%s/%s", api.Service.Name, api.Region.Name, api.Service.Version, strings.Join(path, "/"))
	} else {
		strurl = fmt.Sprintf("https://%s.%s.conoha.io/%s/%s/%s", api.Service.Name, api.Region.Name, api.Service.Version, api.TenantId, strings.Join(path, "/"))
	}

	return url.Parse(strurl)
}

func (api *Api) Prepare(method string, path []string, body []byte) error {
	endpoint, err := api.Endpoint(path)
	if err != nil {
		return err
	}

	api.request, err = http.NewRequest(method, endpoint.String(), strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	if api.Token != "" {
		api.request.Header.Set("X-Auth-Token", api.Token)
	}

	if body != nil {
		api.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return nil
}

func (api *Api) Do() chan []byte {
	client := &http.Client{}

	ch := make(chan []byte)

	go func() {
		api.lastRequestError = nil

		resp, err := client.Do(api.request)
		if err != nil {
			api.lastRequestError = err
		}
		defer resp.Body.Close()

		switch {
		case resp.StatusCode >= 400:
			msg := api.extractApiErrorMessage(resp.Body)
			api.lastRequestError = fmt.Errorf("Return %d status code from the server. (%s)", resp.StatusCode, msg)
		}

		res, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			ch <- res
		} else {
			api.lastRequestError = err
			ch <- nil
		}
	}()

	return ch
}

func (api *Api) LastError() error {
	return api.lastRequestError
}

func (api *Api) extractApiErrorMessage(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err.Error()
	}

	errjson := string(b)

	p := strings.Index(errjson, `"message":"`)

	if p >= 0 {
		p += 11
		pe := strings.Index(errjson[p:], `"`)
		if pe < 0 {
			return errjson
		}
		return errjson[p : p+pe]
	} else {
		return errjson
	}
}
