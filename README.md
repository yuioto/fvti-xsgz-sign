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
}
```
