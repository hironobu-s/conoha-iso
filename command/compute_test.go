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
	_, err := compute.List()
	if err != nil {
		t.Errorf(err.Error())
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
