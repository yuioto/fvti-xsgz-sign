package config

import "net/http"

// sign gps info
const Latitude = "26.075136"
const Longitude = "119.162543"
const SingnSite = "福建省 福州市 闽侯县 闽侯县福州职业技术学院(联榕路北)"

// requeset User-Agent value
const UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)  Mobile/15E148 wxwork/4.1.30 MicroMessenger/7.0.1 Language/zh ColorScheme/Light wwmver/3.26.13.714"

// reaueset
const ExpiredToken_StatusCode = http.StatusUnauthorized
const ExpiredToken = `{"Message":"Authorization has been denied for this request."}`
