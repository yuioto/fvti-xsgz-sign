package savestusignin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cfgset "fvti-xsgz-sign/pkg/set"
)

func GetSignId(id string, authorization string) (string, error) {
	taskList, _ := GetTaskList(authorization)
	qd, err := GetSignIdFromList(taskList, id)
	if err != nil {
		return "", err
	}
	return qd, nil
}

func GetTaskQD(id string, authorization string) (string, error) {
	taskList, _ := GetTaskList(authorization)
	qd, err := GetQDFromList(taskList, id)
	if err != nil {
		return "", err
	}
	return qd, nil
}

func GetTaskId(name string, authorization string) (string, error) {
	taskList, err := GetTaskList(authorization)
	if err != nil {
		return "", err
	}
	id, err := GetIdFromList(taskList, name)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetSignIdFromList(taskjson string, id string) (string, error) {
	var taskList TaskList
	if err := json.Unmarshal([]byte(taskjson), &taskList); err != nil {
		return "", err
	}

	for _, item := range taskList.List.Items {
		if item.Id == id {
			return item.SignId, nil
		}
	}
	return "", fmt.Errorf("item with id %s not fount", id)
}

func GetQDFromList(taskjson string, id string) (string, error) {
	var taskList TaskList
	if err := json.Unmarshal([]byte(taskjson), &taskList); err != nil {
		return "", err
	}

	for _, item := range taskList.List.Items {
		if item.Id == id {
			return item.QD, nil
		}
	}
	return "", fmt.Errorf("item with id %s not fount", id)
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

func GetTaskList(authorization string) (string, error) {
	// must use http, if use https, will get timeout
	url := "http://" + cfgset.Host + "/PhoneApi/api/SignIn/GetStuSignInList"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
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
		return "", fmt.Errorf("error sending request:  %w", err)
	}
	defer resp.Body.Close()

	// Retrieve response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// StatusCode should is Ok
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed get task list, StatusCode: %d, %s", resp.StatusCode, string(body))
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
	QD         string `json:"QD"`
	SignId     string `json:"SignID"`
	QDTimeText string `json:"QDTimeText"`
}
