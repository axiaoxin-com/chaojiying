package chaojiying

import (
	"encoding/json"
	"os"
	"testing"
)

func TestAndDumpAccountsJSON(t *testing.T) {
	// 加载账号信息
	jsonFilename := "../personal-data/chaojiying_accounts.json"
	accounts, err := LoadAccountsFromJSONFile(jsonFilename)
	if err != nil {
		t.Error(err)
	}
	// 创建客户端
	cli, err := New(accounts)
	if err != nil {
		t.Error(err)
	}
	/*
		// 测试破解
		pic, err := os.Open("./example/pin.png")
		if err != nil {
			t.Error(err)
		}
		defer pic.Close()
		pin, err := cli.Cr4ck(pic)
		if err != nil {
			t.Error(err)
		}
		t.Log("pin:", pin)
	*/

	// 更新可用账号到json
	b, err := json.Marshal(cli.accounts)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Create(jsonFilename)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	f.Write(b)

}
