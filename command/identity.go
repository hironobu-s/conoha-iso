package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Identity struct {
	ApiUsername string
	ApiPassword string
	ApiTenantId string

	Token        string
	TokenExpires time.Time
	Region       string

	*Command
}

func NewIdentity() *Identity {
	identity := &Identity{}

	return identity
}

func (cmd *Identity) Auth() (err error) {
	authinfo := map[string]interface{}{
		"auth": map[string]interface{}{
			"tenantId": cmd.ApiTenantId,
			"passwordCredentials": map[string]interface{}{
				"username": cmd.ApiUsername,
				"password": cmd.ApiPassword,
			},
		},
	}

	b, err := json.Marshal(authinfo)
	if err != nil {
		return err
	}

	api, err := NewApi("identity", cmd.Region)
	if err != nil {
		return err
	}

	if err = api.Prepare("POST", []string{"tokens"}, b); err != nil {
		return err
	}

	ch := api.Do()
	strjson := <-ch

	if err = api.LastError(); err != nil {
		return err
	}

	err = cmd.parseResponse(strjson)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Identity) parseResponse(strjson []byte) error {
	var auth map[string]interface{}
	var ok bool
	var err error

	err = json.Unmarshal(strjson, &auth)
	if err != nil {
		return err
	}

	if _, ok = auth["error"]; ok {
		obj := auth["error"].(map[string]interface{})
		msg := fmt.Sprintf("%s(%0.0f): %s",
			obj["title"].(string),
			obj["code"].(float64),
			obj["message"].(string),
		)

		err = errors.New(msg)
		return err
	}

	// アクセストークンを取得
	if _, ok = auth["access"]; !ok {
		err = errors.New("Undefined index: access")
		return err
	}
	access := auth["access"].(map[string]interface{})

	if _, ok = access["token"]; !ok {
		err = errors.New("Undefined index: token")
		return err
	}
	t := access["token"].(map[string]interface{})
	token := t["id"].(string)

	tokenExpires, err := time.Parse(time.RFC3339, t["expires"].(string))
	if err != nil {
		return err
	}

	cmd.Token = token
	cmd.TokenExpires = tokenExpires
	return nil
}
