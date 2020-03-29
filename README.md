# chaojiying

验证码破解平台超级鹰 Golang 版 SDK

## 特性

### 支持普通的单账号破解

调用 `Processing` 方法，传入账号密码和验证码图片返回识别结构结构体

### 支持加载多个账号破解

初始化客户端时设置可用账号列表，调用 `Cr4ck` 方法传入图片直接返回结果。

`Cr4ck` 随机选择出一个可用的账号进行调用，不可用账号将被跳过。

### 支持账号题分查询

调用 `GetScore` 方法，传入账号信息查询可用题分


## 安装

```
go get -u github.com/axiaoxin-com/chaojiying
```

## 在线文档

<https://godoc.org/github.com/axiaoxin-com/chaojiying>

## 用法示例

### 单账号示例

[example/one_account_crack.go](https://github.com/axiaoxin-com/chaojiying/blob/master/example/one_account_crack.go)

```go
// 使用单账号的方式示例

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin-com/chaojiying"
)

func main() {
	// 加载验证码图片
	pic, err := os.Open("./pin.png")
	if err != nil {
		log.Println(err)
	}
	defer pic.Close()

	// 创建客户端
	cli, err := chaojiying.New(nil)
	if err != nil {
		log.Println(err)
	}

	// 从环境变量获取账号信息
	user := os.Getenv("user")
	pass := os.Getenv("pass")

	// 指定账号破解验证码
	resp, err := cli.Processing(user, pass, pic)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("破解结果:", resp.PicStr)

}
```

### 多账号示例

[example/one_account_crack.go](https://github.com/axiaoxin-com/chaojiying/blob/master/example/random_account_crack.go)

```go
// 使用随机选择账号的方式示例
// 从json文件中加载所有可用账号
// json格式： [{"user":"xx", "pass":"yy"}]
// 破解验证码时随机选择其中一个可用账号进行调用

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin-com/chaojiying"
)

func main() {
	// 加载账号json信息
	jsonFilename := "../../personal-data/chaojiying_accounts.json"
	accounts, err := chaojiying.LoadAccountsFromJSONFile(jsonFilename)
	if err != nil {
		log.Fatal(err)
	}

	// 加载验证码图片
	pic, err := os.Open("./pin.png")
	if err != nil {
		log.Fatal(err)
	}
	defer pic.Close()

	// 创建多账号客户端
	cli, err := chaojiying.New(accounts)
	if err != nil {
		log.Fatal(err)
	}

	// 随机账号破解
	pinCode, err := cli.Cr4ck(pic)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("破解结果:", pinCode)

}
```
