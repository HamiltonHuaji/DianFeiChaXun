package main

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type LoginRespData struct {
	Token string `json:"token"`
}

type LoginResp struct {
	MessageCode string        `json:"messageCode"`
	Data        LoginRespData `json:"data"`
}

type QueryRespData struct {
	Oddl string `json:"oddl"`
}

type QueryResp struct {
	MessageCode string          `json:"messageCode"`
	Data        []QueryRespData `json:"data"`
}

func login(accountNum string, password string) string {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	type Payload struct {
		AccountNum  string `json:"accountNum"`  // 学号
		Password    string `json:"password"`    // md5(密码)
		DataSource  string `json:"dataSource"`  // "PAY"
		SchoolCode  string `json:"schoolCode"`  // "datalook"
		AccountType string `json:"accountType"` // "BASE"
	}

	hash := md5.Sum([]byte(password))

	data := Payload{
		// fill struct
		AccountNum:  accountNum,
		Password:    hex.EncodeToString(hash[:]),
		DataSource:  "PAY",
		SchoolCode:  "datalook",
		AccountType: "BASE",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
		panic(err)
	}
	payload := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://ucaspay.ucas.ac.cn/api/pay/web/login/login", payload)
	if err != nil {
		// handle err
		panic(err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "http://ucaspay.ucas.ac.cn")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Referer", "http://ucaspay.ucas.ac.cn/mobile/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.26")

	resp, err := client.Do(req)
	if err != nil {
		// handle err
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	r := LoginResp{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err)
	}
	return r.Data.Token
}

func query(token string, projectId, roomId string) float64 {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	type Payload struct {
		ProjectID  string `json:"projectId"`
		Kddzroomid string `json:"kddzroomid"`
	}

	data := Payload{
		// fill struct
		ProjectID:  projectId,
		Kddzroomid: roomId,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
		panic(err)
	}
	payload := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://ucaspay.ucas.ac.cn/api/pay/web/payEleCostController/querySydl", payload)
	if err != nil {
		// handle err
		panic(err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("datalook_reimbursement_token=%s; datalook_login_status=false", token))
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "http://ucaspay.ucas.ac.cn")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Referer", "http://ucaspay.ucas.ac.cn/mobile/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.26")
	req.Header.Set("X-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		// handle err
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	r := QueryResp{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err)
	}
	if len(r.Data) != 1 {
		panic(errors.New("data wrong!!!"))
	}
	oddl, err := strconv.ParseFloat(r.Data[0].Oddl, 64)
	if err != nil {
		panic(err)
	}
	return oddl
}

type ftResp struct {
	Errno   int    `json:"errno"`
	Errmsg  string `json:"errmsg"`
	Dataset string `json:"dataset"`
}

var Extra string

func Ft(key string, title string, username string) {
	fmt.Println(title)
	if key == "" {
		return
	}
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	// curl -X POST https://sctapi.ftqq.com/XXXXXXXXXXXXXXXXXXXXXXXX.send \
	//  --data-raw title=messagetitle&desp=message@content

	params := url.Values{}
	params.Add("title", title)
	if Extra != "" {
		params.Add("desp", username+"\n\n"+Extra)
	} else {
		params.Add("desp", username)
	}
	body := strings.NewReader(params.Encode())

	req, _ := http.NewRequest("POST", fmt.Sprintf("https://sctapi.ftqq.com/%s.send", key), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
}

var username = flag.String("username", "", "学号")
var password = flag.String("password", "", "密码")
var projectId = flag.String("projectId", "", "校区id")
var roomId = flag.String("roomId", "", "公寓和房间id")
var ftKey string

func errorNotifier(title string) {
	if err := recover(); err != nil {
		Ft(ftKey, title, *username)
		panic(title)
	}
}

func initFlags() {
	ftKey = os.Getenv("FT_SCKEY")
	defer errorNotifier("查询失败：参数配置错误")
	flag.Parse()
	if *username == "" {
		*username = os.Getenv("USERNAME")
	}
	if *password == "" {
		*password = os.Getenv("PASSWORD")
	}
	if *projectId == "" {
		*projectId = os.Getenv("PROJECT")
	}
	if *roomId == "" {
		*roomId = os.Getenv("ROOM")
	}
	if *username == "" || *password == "" || *projectId == "" || *roomId == "" {
		panic("Must specify arguments")
	}
}

func main() {
	initFlags()
	token := login(*username, *password)
	oddl := query(token, *projectId, *roomId)
	fmt.Printf("剩余电量: %.2f度", oddl)
	if oddl < 10 {
		Ft(ftKey, "电量已不足10度", *username)
	}
}
