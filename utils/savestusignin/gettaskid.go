package savestusignin

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	cfgset "fvti-xsgz-sign/utils/config"
)

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
	return "", fmt.Errorf("item with name %s not fount", name)
}

func GetDateFromList(taskjson string, name string) (startTime string, endTime string, err error) {
	var taskList TaskList
	if err := json.Unmarshal([]byte(taskjson), &taskList); err != nil {
		return "", "", err
	}

	for _, item := range taskList.List.Items {
		if item.Name == name {
			if startTime, endTime, err := extractTimes(item.QDTimeText); err != nil {
				return "", "", fmt.Errorf("提取时间范围错误: %v", err)
			} else {
				return startTime, endTime, nil
			}
		}
	}
	return "", "", fmt.Errorf("item with name %s not fount date", name)
}

// 函数提取时间范围
func extractTimes(timeRange string) (startTime string, endTime string, err error) {
	// 按“至”分割字符串
	parts := strings.Split(timeRange, "至")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("时间范围格式错误")
	}

	// 去除多余的空格
	startTime = strings.TrimSpace(parts[0])
	endTime = strings.TrimSpace(parts[1])

	// 验证时间格式
	_, err1 := time.Parse("15:04", startTime)
	_, err2 := time.Parse("15:04", endTime)
	if err1 != nil || err2 != nil {
		return "", "", fmt.Errorf("时间格式错误: %v %v", err1, err2)
	}

	return startTime, endTime, nil
}

func GetTaskList(authorization string) (string, error) {
	url := "https://xsgz.webvpn.fvti.cn/PhoneApi/api/SignIn/GetStuSignInList"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error creating request:", err)
	}

	// Setting the request header
	// Only authorization is required
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br") // set gzip compress
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", cfgset.UserAgent)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Host", cfgset.Host)
	req.Header.Set("Referer", cfgset.Referer)
	req.Header.Set("Accept-Language", "zh-CN,zh-Hans;q=0.9")
	req.Header.Set("Accept", "application/json, text/plain, */*")

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
	Id         string `json:"Id"`
	Name       string `json:"Name"`
	QDTimeText string `json:"QDTimeText"`
}
