/**

Lib: discord-connect
Author: Z3NTL3
License: GNU
Description:

**/

package discordconnect

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/monaco-io/request/request"
)

type (
	/*
		Context for your Discord application	
	*/
	AppContext struct {
		Client_id string `client_id`
		Client_secret string `client_secret`
		Grant_type string `grant_type`
		Code string `code`
		Redirect_uri string `redirect_uri`
	}

	Client struct {
		rawReq *request.Request
		app_ctx *AppContext
	}

	/*
		Discord API response after ``authorization_code`` exchange
	*/
	Response struct {
		Token Token `json:"access_token"`
		Type string `json:"token_type"`
		ExpiresIn int `json:"expires_in"`
		RefreshToken Token `json:"refresh_token"`
		Scope string `json:"scope"`
	}
)


// Initializes the client
func Initialize(timeout time.Duration, appCtx AppContext) *Client {
	client := &Client{
		app_ctx: &appCtx,
		rawReq: request.New().
		AddHeader(map[string]string{
			"cache-control":"must-revalidate",
			"user-agent": DEFAULT_UA,
			"content-type": "application/x-www-form-urlencoded",
		}).
		AddTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		
	}

	client.rawReq.Ctx().Client.Timeout = timeout
	return client
}

// Sets the obtained 'authorization_code' for the client as defined by Discord.
// Ref: https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-authorization-url-example
func (c *Client) SetCode(code string) {
	c.app_ctx.Code = code
}

// Sets proxy, this proxy will be used while perfoming OAUTH2 tasks.
func (c *Client) SetProxy(proxy string) (error) {
	uri, err := url.Parse(proxy); if err != nil {
		return err
	}
	c.rawReq.Ctx().Client.Transport = &http.Transport{
		Proxy: http.ProxyURL(uri),
	}
	return nil
}

// Delete proxy, so no proxy will be used while perfoming requests towards Discord endpoints
func(c *Client) DelProxy() {
	c.rawReq.Ctx().Client.Transport = http.DefaultTransport
}

func(c *Client) Connect() (error, *Response) {
	payload := make(map[string]string)
	apiResult := new(Response)

	values := reflect.ValueOf(*c.app_ctx)
	types := values.Type()

	for i := 0; i < values.NumField(); i++ {
		payload[string(types.Field(i).Tag)] = values.Field(i).String()
	}

	resp := c.rawReq.POST(DISCORD_OAUTH2_TOKEN_URL).
	AddURLEncodedForm(payload).Send().Scan(apiResult)

	if resp.Error() != nil{
		return resp.Error(), nil
	}
	
	if resp.Code() != 200 {
		return errors.New(
			fmt.Sprintf("Discord API Error - BODY: '%s' CODE: '%d'",
				resp.String(), 
				resp.Code(),
			),
		), nil
	}

	return nil, apiResult
}

func(app *AppContext) Change_ClientID(id string) {
	app.Client_id = id
}

func(app *AppContext) Change_ClientSecret(secret string){
	app.Client_secret = secret
}


func(app *AppContext) Change_GrantType(grant_type string){
	app.Grant_type = grant_type
}

func(app *AppContext) Change_Code(code string){
	app.Code = code
}

func(app *AppContext) Change_REDIR_URI(redir_uri string){
	app.Redirect_uri = redir_uri
}