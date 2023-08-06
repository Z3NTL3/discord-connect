# Discord Connect
Implements Discord OAUTH2 within ease. **OOP-like API!**

### Explore API
<a href="https://pkg.go.dev/github.com/Z3NTL3/discord-connect">Documentation</a>

#### Download
``go get github.com/Z3NTL3/discord-connect``

### Example
```go
package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	discordconnect "github.com/Z3NTL3/discord-connect"
)

func main() {
	application := discordconnect.AppContext{
		// Dummy app context
		Client_id: "1131879305987760208",
		Client_secret: "9spr9szfgLvEsV0TjsqH0OaIJ97lVjMR",
		Code: "7ey3za64UPN4DA0xgIYTAGZ4HKSWOW",
		Grant_type: "authorization_code",
		Redirect_uri: "http://localhost:2000/api/v1/discord/callback",
	}

	client := discordconnect.Initialize(time.Duration(time.Second * 5), application)
	
	// client.SetProxy("http://root:root@localhost:9050")
	err, resp := client.Connect(); if err != nil {
		log.Fatal(err)
	}

	ctx := reflect.ValueOf(*resp)
	for i := 0; i < ctx.NumField(); i++ {
		if ctx.Field(i).CanInt() {
			fmt.Println("["+ctx.Type().Field(i).Name+"]", ctx.Field(i).Int())
			continue
		}

		fmt.Println("["+ctx.Type().Field(i).Name+"]", ctx.Field(i).String())
	}
}
```