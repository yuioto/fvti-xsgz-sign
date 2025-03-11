package set

import "net/http"

// sign gps info
const Latitude = "26.075136"
const Longitude = "119.162543"
const SingnSite = "福建省 福州市 闽侯县 闽侯县福州职业技术学院(联榕路北)"

// requeset User-Agent value
const UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)  Mobile/15E148 wxwork/4.1.30 MicroMessenger/7.0.1 Language/zh ColorScheme/Light wwmver/3.26.13.714"

// reauest header
const Host = "zhxg.fvti.edu.cn"
const Origin = "https://" + Host
const Referer = "http://" + Host + "/Phone/index.html"

// encode password publickey
const PublicKeyBase64 = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCK4n2xrbtnRyBqMJ2iiDeDRdJ/F8EVmzcjSGy/vVNfEVahl6sQOjQXZTc8AEbiZdyLnP9QwX3ZkIsEGUz1VMaPUJeHLHQC5uVljRWR0ORt4oiU7mtN5ZsEl8gPQBzSbC7IpnXVRN1Mx7s/RlFsWZgkuZKbPjxcfgoA9zXyhmcHywIDAQAB"

const NotyId = "fvti-xsgz-sign-task-default-status"

// reaueset
const ExpiredToken_StatusCode = http.StatusUnauthorized
const ExpiredToken = `{"Message":"Authorization has been denied for this request."}`

// server mang
const StatusInternalServer_StatusCode = http.StatusInternalServerError
const StatusInternalServer = `<h1>error</h1><span>dial tcp 10.1.2.243:80: connect: cannot assign requested address</span>`

const StatusGatewayTimeout_StatusCode = http.StatusGatewayTimeout
const StatusGatewayTimeout = `<html><body><h1>504 Gateway Time-out</h1>
The server didn't respond in time.
</body></html>`

const StatusSignOk_StatusCode = http.StatusOK
const StatusSignOk = `{"errcode":0,"errmsg":"签到成功！"}`

const StatusSignSuccessfullyOk = "是"
