package savestusignin

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	cfgset "fvti-xsgz-sign/utils/config"
	"net/http"
	"net/url"

	"github.com/dop251/goja"
)

func GetAuthorization(studentId string, password string) (string, error) {
	encodePass, err := GetEncodePassword(password)
	if err != nil {
		return "", err
	}

	bearer, err := login(studentId, encodePass)
	if err != nil {
		return "", err
	}

	bearer = "Bearer" + " " + bearer

	return bearer, nil
}

//go:embed encode.js
var jsScript string

func GetEncodePassword(password string) (string, error) {
	// 创建一个新的JavaScript运行时
	vm := goja.New()

	// 执行嵌入的JavaScript代码
	_, err := vm.RunString(jsScript)
	if err != nil {
		return "", err
	}

	// 调用JavaScript函数
	vm.Set("password", password)
	result, err := vm.RunString(`encode(password)`)
	if err != nil {
		return "", err
	}

	// 获取返回值并打印
	encodePass := result.ToString().String()
	return encodePass, nil
}

func login(studentId string, password string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", studentId)
	data.Set("password", password)

	req, err := http.NewRequest("POST", "https://xsgz.webvpn.fvti.cn/PhoneApi/api/Account/Login?OpenId=", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("User-Agent", cfgset.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode == http.StatusOK {
		var loginResp LoginResponse
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			return "", err
		}
		return loginResp.AccessToken, nil
	}
	return "", fmt.Errorf("login failed with status: %s", resp.Status)
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	UserType     string `json:"UserType"`
	IsActive     bool   `json:"IsActive"`
	Msg          string `json:"Msg"`
}
