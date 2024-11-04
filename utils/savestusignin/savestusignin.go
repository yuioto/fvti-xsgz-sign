package savestusignin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	cfgset "fvti-xsgz-sign/utils/config"
)

func PostStuSignIn(studentid string, id string, authorization string) error {
	signURL := "https://xsgz.webvpn.fvti.cn/PhoneApi/api/SignIn/SaveStuSignIn"

	data := url.Values{}
	data.Add("ApplyInfo[Id]", "00000000-0000-0000-0000-000000000000")
	data.Add("ApplyInfo[OrderId]", id)
	data.Add("ApplyInfo[StudentId]", studentid)
	data.Add("ApplyInfo[SignWayText]", "定位")
	data.Add("ApplyInfo[IsPhoto]", "false")
	data.Add("ApplyInfo[IsLocal]", "true")
	data.Add("ApplyInfo[IsQrCode]", "false")
	data.Add("ApplyInfo[QrCodeContent]", "")
	data.Add("ApplyInfo[IsDWQDW]", "0")
	data.Add("ApplyInfo[SingnScope]", "")
	data.Add("ApplyInfo[Latitude]", cfgset.Latitude)
	data.Add("ApplyInfo[Longitude]", cfgset.Longitude)
	data.Add("ApplyInfo[SingnSite]", cfgset.SingnSite)
	data.Add("ApplyInfo[InputUser]", "")
	data.Add("ApplyInfo[InputDate]", "")
	data.Add("ApplyInfo[collegeNo]", "")
	data.Add("ApplyInfo[classNo]", "")
	data.Add("ApplyInfo[qdType]", "")
	data.Add("ApplyInfo[qdTime]", "")
	data.Add("ApplyInfo[InsertUserId]", "")
	data.Add("ApplyInfo[InsertUserName]", "")
	data.Add("ApplyInfo[InsertDate]", "")

	req, err := http.NewRequest("POST", signURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	// Setting the request header
	// Only authorization is required
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br") // set gzip compress
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Retrieve response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	// StatusCode should is Ok
	if resp.StatusCode != cfgset.StatusSignOk_StatusCode {
		return fmt.Errorf("failed post sign, StatusCode: %d, %s", resp.StatusCode, string(body))
	}
	if QD, err := GetTaskQD(id, authorization); QD != cfgset.StatusSignSuccessfullyOk {
		return fmt.Errorf("server ok, but sign failed: %w", err)
	}

	return nil
}
