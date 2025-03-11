package savestusignin

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	cfgset "fvti-xsgz-sign/pkg/set"
	"net/http"
	"net/url"
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

func GetEncodePassword(password string) (string, error) {
	publicKeyBase64 := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCK4n2xrbtnRyBqMJ2iiDeDRdJ/F8EVmzcjSGy/vVNfEVahl6sQOjQXZTc8AEbiZdyLnP9QwX3ZkIsEGUz1VMaPUJeHLHQC5uVljRWR0ORt4oiU7mtN5ZsEl8gPQBzSbC7IpnXVRN1Mx7s/RlFsWZgkuZKbPjxcfgoA9zXyhmcHywIDAQAB"

	publicKey, err := parsePublicKey(publicKeyBase64)
	if err != nil {
		return "", fmt.Errorf("公钥解析失败: %v", err)
	}

	// 确保 publicKey 不是 nil
	if publicKey == nil {
		return "", errors.New("公钥为空")
	}

	encrypted, err := encryptPassword(password, publicKey)
	if err != nil {
		return "", fmt.Errorf("加密失败: %v", err)
	}

	return encrypted, nil
}

func login(studentId string, password string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", studentId)
	data.Set("password", password)

	// must use http, if use https, will get timeout
	req, err := http.NewRequest("POST", "http://"+cfgset.Host+"/PhoneApi/api/Account/Login?OpenId=", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("User-Agent", cfgset.UserAgent)
	req.Header.Set("Host", cfgset.Host)
	req.Header.Set("Referer", cfgset.Referer)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", cfgset.Origin)
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
		/* debug: print response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		fmt.Println("Response body:", string(body))
		*/
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
