package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"

	"github.com/axiaoxin-com/chaojiying"
)

func main() {
	// 加载账号json信息
	jsonFilename := "./accounts.json"
	accounts, err := chaojiying.LoadAccountsFromJSONFile(jsonFilename)
	if err != nil {
		log.Fatal(err)
	}

	src := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAAAmCAYAAAAycj4zAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAPfSURBVGhD7Zo7buMwEIZ9HecKewFVOYCrnGArHWDVJsWWSSEELpw6SOEiQBZuHMC9UizgxnAdBEi33azkkNKQ4mNIjhTn8QEsRJEzw/lNDhVkAkxUVXVo74ktBkpslDFjMMGBpASVMpcLWwyU2ChjxoBNEBtD2PzI+PLhPbJSE5o6P4ZQn0PFaLLr82UUBE/yGRiCVJ+h81P92Yix6xUkhNh5Ohx2KDY4/HDDdstqGHuBLn+UWChjGqjjOGAVhJ1VXkdYh9i2+hkRmqiYxMo5eJ7tWe+P4YgFWWli9AVx8rKB7OICJs42h/JFjLdgSrLeJ5/1/hhIglAccQSjgHdHXosTSqQgrnWwr9HA6IJQx2FB/hX37RzyfBvbZSfIcis6O1z2D+8ezmE5ncLNyUnXZgt4FWNSCRbEFXCD6T3uM71X2QFkYmfordwR5rvYwlm7O5bNodjit7uHp5kmhNYe12KoB5evepU0pBGXsQb8fnU7gcmvrmXrHezWmXjOlYR0uAVJYbeZt7sj1zaHe11+MW6mP2HtyAvG5ateJQ13wDo7KC9VMfrNJogA15BEId7YQi53x3zTyE7mdTFDiZ/B0168QITlxz6eLIiE4rjbBXW7LNvFqzvGLsjBx/UZqyCu3eFmDY+yZpyes9WKEQXBuyOD8ll0H1hBjgS5NtiS9imC+GORvEA59+8Oo739oivi9bFk2BysBAviByUd7Y43sFhuQShHljGBJtDNKtvYPzzM9tAOsbWCWM0JKIKQF+iELsgYNaRZz/3dlRDE/yFoQqkhtmaoLTH5/AKCbKC4EseVdtVtoK+ZcNPSvkdC8inHDiCII+nPJWRSrNEEeei+PQI/BEkoNeYU/jzE2ZJxtIIkB4b5m4uk1+1Wph0LNaIgmzv4IQRx1Y8U8JG2XByjIL3km9o4glALuhH5658WjtuVepRRv9Z1eoIMgfrdcVbfqu7pNcRC6A8n/vtDK+aGm1Sv2FuEC4l5UEH6BBR1C3Jx+iJNfTao42BdqAn3tOXCvI/I/mqcgoQYolDVO6T4nSaIRI9NPlPidY3V+/eF52YlWmzt0LEKogetBxqDKkhzhKn2OHxQsfky9uObFG6Wv2ul4BVEYgw0GHxkva8gOu/pGxNVQ44leE5MawpdZ+h4E8p/LlJJdZw6fyxC4+RYF1kQDmcSTlufDfKR9Z1EN1z56QkiDXM5+CpQ8+Z7HyWI7Z1rzlfBlwPfe9Zbls/ZN36sgny25H6U9SiC4KBjFiDnhM4zwWVHEmMvdE6MDx2rIDHI+alBNeh2uOyGEOqTI8ZWEA5jQ3Ls8XFQVRX8BxZDC9OmPWGyAAAAAElFTkSuQmCC"
	data, err := base64.StdEncoding.DecodeString(src[22:]) // 去掉前缀 data:image/png;base64,
	if err != nil {
		// 处理错误
	}

	im, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		panic("Bad png")
	}

	// 创建缓冲区
	buff := new(bytes.Buffer)
	// 将 im 编码到缓冲区
	err = png.Encode(buff, im)
	if err != nil {
		fmt.Println("failed to create buffer", err)
	}
	// 将缓冲区转换为 reader
	pic := bytes.NewReader(buff.Bytes())

	cli, err := chaojiying.New(accounts)
	if err != nil {
		log.Fatal(err)
	}

	pinCode, err := cli.Cr4ck(pic)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("破解结果:", pinCode)
}
