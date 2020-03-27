package chaojiying_test

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin-com/chaojiying"
)

// 使用单账号的方式示例
func Example_oneAccount() {
	// 加载验证码图片
	pic, err := os.Open("./example/pin.png")
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

	// Output:
	// 破解结果: myeee

}
