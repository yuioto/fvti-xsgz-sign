# fvti-xsgz-sign

Sending sign-in data

## build

build fvti-xsgz-sign

```shell
go build -o out/ ./cmd/sign
```

build make config file tool

```shell
go build -o out/ ./cmd/mkconfig
```

## config

use [kdl v1](https://kdl.dev/spec-v1)

```kdl
login {
    student_id "id"
    password "cardid[len(cardid)-8:]"
}

notify {
    ntfy {
        topic "fvti-xsgz-sign-task-default-status"
    }
    email {
        host "smtp.example.com"
        port "587"
        username "user@example.com"
        password "email_password"
        from "user@example.com"
        from_name "定时签到"  # 可选，邮件 From 头显示昵称（默认定时签到）
        to "recipient@example.com"
        cc "cc@example.com"  # 可选，抄送地址，支持逗号分隔
    }
}

log {
    console true
}
```
