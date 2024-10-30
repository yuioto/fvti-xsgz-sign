package checklogin

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	cfgset "fvti-xsgz-sign/utils/config"
	"fvti-xsgz-sign/utils/savestusignin"
)

func SignLogin(id string, signjson string, authorization string) error {
	var signList SignList
	if err := json.Unmarshal([]byte(signjson), &signList); err != nil {
		return err
	}

	const formatQD = "2006/01/02"
	const formatInputDate = "2006/01/02 15:04:05"

	startDate, endDate, err := savestusignin.GetDateFromList()
	if err != nil {
		return fmt.Errorf("提取时间时出错: %v", err)
	}

	for page, time := 1, time.Now(); signList.List.Items != Items; page += 1{
		for _, item := range signList.List.Items {
			if item.QD == time.Format(formatQD) {
				return nil
			}
		}
		return fmt.Errorf("item with today date %s not fount", time.Format(formatQD))
	}

	return nil
}


func GetSignHistory(id string, page string, authorization string) (string, error) {
	signListURL := "https://xsgz.webvpn.fvti.cn/PhoneApi/api/SignIn/GetStuSignInDetailList"
	Values := url.Values{}
	Values.Set("pageIndex", page)
	Values.Set("Id", id)
	signListURL = signListURL + "?" + Values.Encode()

	req, err := http.NewRequest("GET", signListURL, nil)
	if err != nil {
		log.Fatalln("Error creating request:", err)
	}

	// Setting the request header
	// Only authorization is required
	req.Header.Set("Accept", "application/json, text/plain, */*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br") // set gzip compress
	req.Header.Set("Accept-Language", "zh-CN,zh-Hans;q=0.9")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", cfgset.Host)
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", cfgset.Referer)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("TE", "trailers")
	req.Header.Set("User-Agent", cfgset.UserAgent)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Retrieve response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response:", err)
	}

	// StatusCode should is Ok
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Failed get sign history, StatusCode:", resp.StatusCode, string(body))
	}

	return string(body), nil
}

type SignList struct {
	List List `json:"List"`
}

type List struct {
	Items []Items `json:"Items"`
}

type Items struct {
	QD        string `json:"QD"`
	InputDate string `json:"InputDate"`
}
