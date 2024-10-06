package savestusignin

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)  Mobile/15E148 wxwork/4.1.30 MicroMessenger/7.0.1 Language/zh ColorScheme/Light wwmver/3.26.13.714"

func GetTaskId(name string, authorization string) (string, error) {
	taskList, _ := GetTaskList(authorization)
	id, err := GetIdFromList(taskList, name)
	if err != nil {
		log.Fatalln(err)
	}
	return id, nil
}

func GetIdFromList(taskjson string, name string) (string, error) {
	var taskList TaskList
	if err := json.Unmarshal([]byte(taskjson), &taskList); err != nil {
		return "", err
	}

	for _, item := range taskList.List.Items {
		if item.Name == name {
			return item.Id, nil
		}
	}
	return "", fmt.Errorf("Item with name %s not fount", name)
}

func GetTaskList(authorization string) (string, error) {
	url := "https://xsgz.webvpn.fvti.cn/PhoneApi/api/SignIn/GetStuSignInList"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error creating request:", err)
	}

	// Setting the request header
	// Only authorization is required
	req.Header.Set("Accept", "application/json, text/plain, */*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br") // set gzip compress
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Referer", "https://xsgz.webvpn.fvti.cn/Phone/index.html")
	req.Header.Set("Host", "xsgz.webvpn.fvti.cn")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Accept-Language", "zh-CN,zh-Hans;q=0.9")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Sec-Fetch-Dest", "empty")

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

	// StatusCode should is ok(200)
	if resp.StatusCode != 200 {
		log.Fatalln("Failed get task list, StatusCode:", resp.StatusCode, string(body))
	}

	return string(body), nil
}

type TaskList struct {
	List List `json:"List"`
}

type List struct {
	Items []Items `json:"Items"`
}

type Items struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}
