package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Regions
const (
	TYO1 = iota
	SIN1
	SJC1
)

type Identity struct {
	ApiUsername string
	ApiPassword string
	ApiTenantId string
	Region      int

	Token        string
	TokenExpires time.Time

	identityEndpoints map[int]string

	*Command
}

func NewIdentity() *Identity {
	identity := &Identity{}

	identity.identityEndpoints = map[int]string{
		TYO1: "https://identity.tyo1.conoha.io/v2.0",
		SIN1: "https://identity.sin1.conoha.io/v2.0",
		SJC1: "https://identity.sjc1.conoha.io/v2.0",
	}

	// identity.computeEndpoints = map[string]string{
	// 	"tyo1": "https://compute.tyo1.conoha.io/v2/",
	// 	"sin1": "https://compute.sin1.conoha.io/v2/",
	// 	"sjc1": "https://compute.sjc1.conoha.io/v2/",
	// }
	return identity
}

// 認証を実行
func (cmd *Identity) Auth() error {

	// アカウント情報
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

	// 認証リクエスト
	endpoint, ok := cmd.identityEndpoints[cmd.Region]
	if !ok {
		return fmt.Errorf("Undefined region \"%s\"", cmd.Region)
	}

	req, err := http.NewRequest(
		"POST",
		endpoint+"/tokens",
		strings.NewReader(string(b)),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	client := &http.Client{}

	// httpリクエスト実行
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

	strjson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = cmd.parseResponse(strjson)
	if err != nil {
		return err
	}

	return nil
}

// レスポンスのJSONをパースする
func (cmd *Identity) parseResponse(strjson []byte) error {
	// jsonパース
	var auth map[string]interface{}
	var ok bool
	var err error

	err = json.Unmarshal(strjson, &auth)
	if err != nil {
		return err
	}

	// 認証失敗など
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

	// トークンの有効期限を取得
	tokenExpires, err := time.Parse(time.RFC3339, t["expires"].(string))
	if err != nil {
		return err
	}

	cmd.Token = token
	cmd.TokenExpires = tokenExpires
	return nil
}
