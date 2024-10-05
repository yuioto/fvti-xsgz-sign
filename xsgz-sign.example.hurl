GET https://xsgz.webvpn.fvti.cn/PhoneApi/api/SignIn/GetStuSignInList

    [Headers]
    Accept: application/json, text/plain, */*
    Accept-Encoding: gzip, deflate, br
    Sec-Fetch-Mode: cors
    Referer: https://xsgz.webvpn.fvti.cn/Phone/index.html
    Host: xsgz.webvpn.fvti.cn
    Authorization: {{Bearer}}

    HTTP 200

POST https://xsgz.webvpn.fvti.cn/PhoneApi/api/SignIn/SaveStuSignIn

    [Headers]
    Sec-Fetch-Dest: empty
    Connection: keep-alive
    Accept-Encoding: gzip, deflate, br
    Content-Type: application/x-www-form-urlencoded
    Sec-Fetch-Site: same-origin
    Origin: https://xsgz.webvpn.fvti.cn
    User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)  Mobile/15E148 wxwork/4.1.30 MicroMessenger/7.0.1 Language/zh ColorScheme/Light wwmver/3.26.13.714
    Authorization: {{Bearer}}
    Sec-Fetch-Mode: cors
    Host: xsgz.webvpn.fvti.cn
    Referer: https://xsgz.webvpn.fvti.cn/Phone/index.html
    Accept-Language: zh-CN,zh-Hans;q=0.9
    Accept: application/json, text/plain, */*
    
    [FormParams]
    ApplyInfo[Id]: 00000000-0000-0000-0000-000000000000
    ApplyInfo[OrderId]: {{Id}}          # This is sign task id
    ApplyInfo[StudentId]: {{StudentId}} # This is student id
    ApplyInfo[SignWayText]: 定位
    ApplyInfo[IsPhoto]: false
    ApplyInfo[IsLocal]: true
    ApplyInfo[IsQrCode]: false
    ApplyInfo[QrCodeContent]: 
    ApplyInfo[IsDWQDW]: 0
    ApplyInfo[SingnScope]: 
    ApplyInfo[Latitude]: 26.075136
    ApplyInfo[Longitude]: 119.162543
    ApplyInfo[SingnSite]: 福建省 福州市 闽侯县 福州职业技术学院
    ApplyInfo[InputUser]: 
    ApplyInfo[InputDate]: 
    ApplyInfo[collegeNo]: 
    ApplyInfo[classNo]: 
    ApplyInfo[qdType]: 
    ApplyInfo[qdTime]: 
    ApplyInfo[InsertUserId]: 
    ApplyInfo[InsertUserName]: 
    ApplyInfo[InsertDate]: 

    HTTP 200
