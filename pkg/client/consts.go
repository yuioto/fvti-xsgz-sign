package client

import "net/http"

const (
	// DefaultHost is the default host for the API.
	DefaultHost = "zhxg.fvti.edu.cn"
	// DefaultUserAgent is the default user agent string.
	DefaultUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)  Mobile/15E148 wxwork/4.1.30 MicroMessenger/7.0.1 Language/zh ColorScheme/Light wwmver/3.26.13.714"
	// PublicKeyBase64 is the public key for password encryption.
	PublicKeyBase64 = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCK4n2xrbtnRyBqMJ2iiDeDRdJ/F8EVmzcjSGy/vVNfEVahl6sQOjQXZTc8AEbiZdyLnP9QwX3ZkIsEGUz1VMaPUJeHLHQC5uVljRWR0ORt4oiU7mtN5ZsEl8gPQBzSbC7IpnXVRN1Mx7s/RlFsWZgkuZKbPjxcfgoA9zXyhmcHywIDAQAB"
)

// Default location (can be overridden in Client)
const (
	DefaultLatitude  = "26.075136"
	DefaultLongitude = "119.162543"
	DefaultSignSite  = "福建省 福州市 闽侯县 闽侯县福州职业技术学院(联榕路北)"
)

const (
	// StatusSignSuccessfullyOk indicates a successful sign-in status string.
	StatusSignSuccessfullyOk = "是"
	// StatusSignOkStatusCode is the HTTP status code for a successful sign-in.
	StatusSignOkStatusCode = http.StatusOK
)

const (
	pathLogin       = "/PhoneApi/api/Account/Login"
	pathGetTaskList = "/PhoneApi/api/SignIn/GetStuSignInList"
	pathSign        = "/PhoneApi/api/SignIn/SaveStuSignIn"
)
