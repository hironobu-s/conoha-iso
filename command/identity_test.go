package command

import (
	"os"
	"testing"
)

func setTestAuthentication(identity *Identity) bool {
	identity.ApiUsername = os.Getenv("OS_USERNAME")
	identity.ApiPassword = os.Getenv("OS_PASSWORD")
	identity.ApiTenantId = os.Getenv("OS_TENANT_ID")
	if identity.ApiTenantId == "" {
		identity.ApiTenantName = os.Getenv("OS_TENANT_NAME")
	}
	identity.Region = os.Getenv("OS_REGION_NAME")

	if identity.ApiUsername != "" &&
		identity.ApiPassword != "" &&
		(identity.ApiTenantId != "" || identity.ApiTenantName != "") {
		return true
	} else {
		return false
	}
}

func TestNewIdentity(t *testing.T) {
	identity := NewIdentity()
	if identity == nil {
		t.Errorf("Invalid type")
	}
}

func TestAuthOk(t *testing.T) {
	identity := NewIdentity()
	if !setTestAuthentication(identity) {
		t.Skip("Skip test that needs connect to API.")
	}

	if err := identity.Auth(); err != nil {
		t.Errorf(err.Error())
	}

	if identity.Token == "" {
		t.Errorf("Token is not set")
	}
}

func TestAuthNoRegion(t *testing.T) {
	identity := NewIdentity()

	if !setTestAuthentication(identity) {
		t.Skip("Skip test that needs connect to API.")
	}

	// Region is null
	identity.Region = ""

	if err := identity.Auth(); err == nil {
		t.Errorf("No region specified. Test should be error.")
	}
}
