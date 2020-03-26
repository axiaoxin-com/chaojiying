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
