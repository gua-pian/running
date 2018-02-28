package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/gin-gonic/gin"
	"running/handler"
)

const (
	wxAppId           = ""                // 填上自己的参数
	wxAppSecret       = ""                // 填上自己的参数
	oauth2RedirectURI = "http://xxxx.com" // 填上自己的参数
	oauth2Scope       = "snsapi_userinfo" // 填上自己的参数
)

var oauth2Endpoint oauth2.Endpoint = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Any("/running", Page2Handler)
	r.Any("/uploadImage", handler.UploadImage)
	r.Any("/uploadData", handler.UploadData)
	r.Any("/redirect", Redirect)

	r.Run(":8000")
}

func Page2Handler(c *gin.Context) {

	w := c.Writer
	r := c.Request

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	code := queryValues.Get("code")
	if code == "" {
		log.Println("用户禁止授权")
		return
	}

	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}
	log.Printf("token: %+v\r\n", token)

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}
	// Insert userinfo into redis.
	handler.NewUser(userinfo.OpenId, userinfo.HeadImageURL, userinfo.Sex)

	// Get current users and steps.
	total_users, total_steps, total_kilos, personal_total_steps, personal_total_kilos := handler.Count(token.OpenId)
	fmt.Println("in server.go, total_steps:", total_steps)

	log.Printf("userinfo: %+v\r\n", userinfo)
	sigature := handler.Running(c)
	c.HTML(http.StatusOK, "running.tmpl", gin.H{
		"AppId":              sigature.AppId,
		"NonceStr":           sigature.NonceStr,
		"TimeStamp":          sigature.TimeStamp,
		"Url":                sigature.Url,
		"Signature":          sigature.Signature,
		"OpenId":             token.OpenId,
		"TotalUsers":         total_users,
		"TotalSteps":         total_steps,
		"TotalKilos":         total_kilos,
		"PersonalTotalSteps": personal_total_steps,
		"PersonalTotalKilos": personal_total_kilos,
	})
	// http.ServeFile(w, r, "./running.html")
	return
}
