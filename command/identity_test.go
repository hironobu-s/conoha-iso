package command

import (
	"testing"
)

func doConnectToApi() bool {
	if API_USERNAME != "" && API_PASSWORD != "" && API_TENANT_ID != "" {
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
	if !doConnectToApi() {
		t.Skip("Skip test that needs connect to API.")
	}

	identity := NewIdentity()
	identity.ApiUsername = API_USERNAME
	identity.ApiPassword = API_PASSWORD
	identity.ApiTenantId = API_TENANT_ID
	identity.Region = REGION

	if err := identity.Auth(); err != nil {
		t.Errorf(err.Error())
	}

	if identity.Token == "" {
		t.Errorf("Token is not set")
	}
}

func TestAuthNoRegion(t *testing.T) {
	if !doConnectToApi() {
		t.Skip("Skip test that needs connect to API.")
	}

	identity := NewIdentity()
	identity.ApiUsername = API_USERNAME
	identity.ApiPassword = API_PASSWORD
	identity.ApiTenantId = API_TENANT_ID

	if err := identity.Auth(); err == nil {
		t.Errorf("No region specified. Test should be error.")
	}
}
