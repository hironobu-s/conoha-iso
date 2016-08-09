package command

import (
	"testing"
)

var identity *Identity

func TestPrepare(t *testing.T) {
	identity = NewIdentity()
	if !setTestAuthentication(identity) {
		identity = nil
		t.Skip("Skip test that needs connect to API.")
	}

	if err := identity.Auth(); err != nil {
		t.Fatalf(err.Error())
	}

	if identity.Token == "" {
		t.Fatalf("No token")
	}
}

func TestNewCompute(t *testing.T) {
	i := NewIdentity()
	compute := NewCompute(i)
	if compute == nil {
		t.Errorf("Invalid type")
	}
}

func TestList(t *testing.T) {
	if identity == nil {
		t.Skip("Skip test that needs connect to API.")
	}

	compute := NewCompute(identity)
	_, err := compute.Servers()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestServer(t *testing.T) {
	if identity == nil {
		t.Skip("Skip test that needs connect to API.")
	}

	compute := NewCompute(identity)
	ss, err := compute.Servers()
	if err != nil {
		t.Errorf(err.Error())
	} else if len(ss.Servers) == 0 {
		t.Skip("Skip test that no servers found.")
		return
	}

	s, err := compute.Server(ss.Servers[0].Id)
	if err != nil {
		t.Errorf(err.Error())
	} else if s.Id != ss.Servers[0].Id {
		t.Errorf("UUID does not match[%s != %s].", s.Id, ss.Servers[0].Id)
	}
}

func TestDownload(t *testing.T) {
	if identity == nil {
		t.Skip("Skip test that needs connect to API.")
	}

	compute := NewCompute(identity)
	err := compute.Download(DOWNLOAD_URL)
	if err != nil {
		t.Errorf(err.Error())
	}
}
