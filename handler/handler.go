package handler

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	rand "math/rand"
	// "net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	wxAppId     = "" // 填写 appid
	wxAppSecret = "" // 填写 appsecret
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	redisClient *redis.Client
)

type (
	TokenResponse struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpireIn    int64  `json:"expires_in"`
	}
	JsApiTicketResponse struct {
		ErrCode  int    `json:"errcode"`
		ErrMsg   string `json:"errmsg"`
		Ticket   string `json:"ticket"`
		ExpireIn int64  `json:"expires_in"`
	}
	Signature struct {
		AppId     string `json:"appId"`
		NonceStr  string `json:"nonceStr"`
		TimeStamp int    `json:"timestamp"`
		Url       string `json:"url"`
		Signature string `json:"signature"`
	}
)

const (
	wechatAPI = "https://api.weixin.qq.com"
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rand.Seed(time.Now().UnixNano())
}

func Running(c *gin.Context) Signature {
	// Get AccessToken
	accessToken, err := getAccesToken()
	if err != nil {
		fmt.Println("can not get accesstoken.")
	}
	fmt.Println(accessToken)

	// Get JsApiTicket
	jsApiTicket, err := getJsApiTicket(accessToken)
	if err != nil {
		fmt.Println("can not get JsApiTicket.")
	}
	fmt.Println(jsApiTicket)

	// Generate url
	request := c.Request
	url := "http://" + request.Host + request.RequestURI

	// TimeStamp
	timestamp := int(time.Now().Unix())
	str_timestamp := strconv.Itoa(timestamp)

	// Nonce
	nonceStr := RandStringRunes(16)

	// raw_string
	raw_string := "jsapi_ticket=" + jsApiTicket + "&noncestr=" + nonceStr + "&timestamp=" + str_timestamp + "&url=" + url

	// sha1kkk
	sign := fmt.Sprintf("%x", sha1.Sum([]byte(raw_string)))

	return Signature{AppId: wxAppId, NonceStr: nonceStr, TimeStamp: timestamp, Url: url, Signature: sign}

}

func getAccesToken() (string, error) {
	// Check if exist in redis.
	accessToken, err := redisClient.Get("accessToken").Result()
	if err != nil {
		fmt.Println("Can't get accessToken:", err)
		// Del JsApiTicket
		redisClient.Del("jsApiTicket")
		// Get accessToken
		params := url.Values{}
		params.Set("grant_type", "client_credential")
		params.Set("appid", wxAppId)
		params.Set("secret", wxAppSecret)
		url := wechatAPI + "/cgi-bin/token?" + params.Encode()
		token := &TokenResponse{}
		err := httplib.Get(url).ToJSON(token)
		if err != nil {
			return "", err
		}
		if token.AccessToken == "" || token.ExpireIn <= 0 {
			return "", errors.New("token not exist or expired")
		}
		// Store in redis.
		redisClient.Set("accessToken", token.AccessToken, time.Duration(token.ExpireIn)*time.Second)
		return token.AccessToken, nil
	}
	return accessToken, nil
}

func getJsApiTicket(accessToken string) (string, error) {
	// Check if exist in redis.
	JsApiTicket, err := redisClient.Get("jsApiTicket").Result()
	if err != nil {
		fmt.Println("Can't get JsApiTicket:", err)
		// Get JsApiTicket
		params := url.Values{}
		params.Set("type", "jsapi")
		params.Set("access_token", accessToken)
		url := wechatAPI + "/cgi-bin/ticket/getticket?" + params.Encode()
		jsApiTicket := &JsApiTicketResponse{}
		err := httplib.Get(url).ToJSON(jsApiTicket)
		if err != nil {
			return "", err
		}
		if jsApiTicket.Ticket == "" || jsApiTicket.ExpireIn <= 0 {
			return "", errors.New("jsApiTicket not exist or expired")
		}
		// Store in redis.
		redisClient.Set("jsApiTicket", jsApiTicket.Ticket, time.Duration(jsApiTicket.ExpireIn)*time.Second)
		return jsApiTicket.Ticket, nil
	}
	return JsApiTicket, nil
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
