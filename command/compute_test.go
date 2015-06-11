package command

import (
	"testing"
)

var identity *Identity

func TestPrepare(t *testing.T) {
	if !doConnectToApi() {
		t.Skip("Skip test that needs connect to API.")
	}

	identity = NewIdentity()
	identity.ApiUsername = API_USERNAME
	identity.ApiPassword = API_PASSWORD
	identity.ApiTenantId = API_TENANT_ID
	identity.Region = REGION

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
	if !doConnectToApi() {
		t.Skip("Skip test that needs connect to API.")
	}

	compute := NewCompute(identity)
	_, err := compute.List()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestDownload(t *testing.T) {
	if !doConnectToApi() {
		t.Skip("Skip test that needs connect to API.")
	}

	compute := NewCompute(identity)
	err := compute.Download(DOWNLOAD_URL)
	if err != nil {
		t.Errorf(err.Error())
	}
}
