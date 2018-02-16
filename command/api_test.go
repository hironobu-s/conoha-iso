package command

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// For test only
func (a *Api) Request() *http.Request {
	return a.request
}

func TestRegionSet(t *testing.T) {
	var err error

	for _, s := range []string{TYO1, TYO2, SIN1, SJC1} {
		region := &Region{}
		if err = region.Set(s); err != nil {
			t.Errorf("Test shoud pass in this case. region=%s", s)
		}
	}

	for _, s := range []string{"invalid", "undefined"} {
		region := &Region{}
		if err = region.Set(s); err == nil {
			t.Errorf("Test shoud fail in this case. region=%s", s)
		}
	}
}

func TestServiceSet(t *testing.T) {
	var err error

	for _, s := range []string{IDENTITY, COMPUTE} {
		service := &Service{}
		if err = service.Set(s); err != nil {
			t.Errorf("Test shoud pass in this case. service=%s", s)
		}
	}

	for _, s := range []string{"invalid", "undefined"} {
		service := &Service{}
		if err = service.Set(s); err == nil {
			t.Errorf("Test shoud fail in this case. service=%s", s)
		}
	}
}

func TestApiNew(t *testing.T) {
	for _, s := range []string{IDENTITY, COMPUTE} {
		for _, r := range []string{TYO1, TYO2, SIN1, SJC1} {
			api, err := NewApi(s, r)
			if api == nil || err != nil {
				t.Errorf("Test shoud be pass in this case. service=%s region=%s", r, s)
			}
		}
	}

	for _, s := range []string{"invalid"} {
		for _, r := range []string{"undefined"} {
			api, err := NewApi(s, r)
			if api != nil || err == nil {
				t.Errorf("Test shoud be error in this case. service=%s region=%s", r, s)
			}
		}
	}
}

func TestApiEndpoint(t *testing.T) {
	var url *url.URL

	for _, r := range []string{TYO1, TYO2, SIN1, SJC1} {
		api, _ := NewApi("identity", r)
		url, _ = api.Endpoint([]string{"tokens"})
		expect := fmt.Sprintf("https://identity.%s.conoha.io/v2.0/tokens", r)
		if url.String() != expect {
			t.Errorf("Endpoint url is not correct. [%s]", url.String())
		}
	}

	for _, r := range []string{TYO1, TYO2, SIN1, SJC1} {
		api, _ := NewApi("compute", r)
		url, _ = api.Endpoint([]string{"iso-images"})
		expect := fmt.Sprintf("https://compute.%s.conoha.io/v2/iso-images", r)
		if url.String() != expect {
			t.Errorf("Endpoint url is not correct. [%s]", url.String())
		}
	}
}

func TestApiPrepare(t *testing.T) {
	var err error
	api, _ := NewApi("identity", TYO1)
	if err = api.Prepare("GET", []string{}, nil); err != nil {
		t.Error(err.Error())
	}
}

func TestApiPrepareToken(t *testing.T) {
	var err error
	api, _ := NewApi("identity", TYO1)

	var testtoken = "testtoken"
	api.Token = testtoken
	if err = api.Prepare("GET", []string{}, nil); err != nil {
		t.Error(err.Error())
	}

	req := api.Request()
	if req.Header.Get("X-Auth-Token") != testtoken {
		t.Error("Token is not match.")
	}
}

func TestApiPrepareContentType(t *testing.T) {
	var err error
	api, _ := NewApi("identity", TYO1)

	var body = []byte("testbody")
	if err = api.Prepare("POST", []string{}, body); err != nil {
		t.Error(err.Error())
	}

	req := api.Request()
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Error("Content-type is not set.")
	}
}

func TestApiDo(t *testing.T) {
	var err error
	api, _ := NewApi("identity", TYO1)

	if err = api.Prepare("GET", []string{}, nil); err != nil {
		t.Error("Can't prepared.")
	}

	ch := api.Do()
	resp := <-ch

	if err = api.LastError(); err != nil {
		t.Error(err.Error())
	}

	if !strings.HasPrefix(string(resp), `{"version":`) {
		t.Errorf("Invalid response.")
	}
}

func TestApiLastError(t *testing.T) {
	var err error
	api, _ := NewApi("identity", TYO1)

	// Error 404 Url
	if err = api.Prepare("GET", []string{"hogehoge"}, nil); err != nil {
		t.Error("Can't prepared.")
	}

	ch := api.Do()
	_ = <-ch

	if err = api.LastError(); err == nil {
		t.Error("Test shoud be error in this case. (404 not found)")
	}
}
