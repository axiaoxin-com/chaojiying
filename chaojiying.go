// Package chaojiying 超级鹰封装
package chaojiying

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/axiaoxin-com/logging"

	"github.com/pkg/errors"
)

// GetScoreResp 超级鹰GetScore接口返回结构
type GetScoreResp struct {
	ErrNo     int    `json:"err_no"`
	ErrStr    string `json:"err_str"`
	Tifen     int    `json:"tifen"`
	TifenLock int    `json:"tifen_lock"`
}

// ProcessingResp 超级鹰Processing接口返回结构
type ProcessingResp struct {
	ErrNo  int    `json:"err_no"`
	ErrStr string `json:"err_str"`
	PicID  string `json:"pic_id"`
	PicStr string `json:"pic_str"`
	Md5    string `json:"md5"`
}

// Account 超级鹰账号
type Account struct {
	User  string `json:"user"`
	Pass  string `json:"pass"`
	Score int    `json:"score"`
}

// Client 超级鹰客户端结构体对象
type Client struct {
	accounts []Account
}

// New 创建超级鹰Client对象
// 传入所有可用账号备用
func New(accounts []Account) (*Client, error) {
	c := &Client{}
	// accounts 传nil则不添加账号
	if accounts == nil {
		return c, nil
	}
	// 过滤不可用账号
	avaliableAccounts := []Account{}
	for _, account := range accounts {
		resp, err := c.GetScore(account.User, account.Pass)
		if err != nil {
			logging.Errorw(nil, "chaojiying New GetScore error", "user", account.User, "pass", account.Pass, "error", err, "resp", resp)
			continue
		}
		account.Score = resp.Tifen
		avaliableAccounts = append(avaliableAccounts, account)
		logging.Debugs(nil, "chaojiying New add account:", account)
	}
	if len(avaliableAccounts) == 0 {
		return nil, errors.New("No avaliable accounts")
	}
	c.accounts = avaliableAccounts
	return c, nil
}

// PickOneAccount 随机选择一个可用账号
// 遍历随机账号列表直到选出一个可用的账号
func (c *Client) PickOneAccount() (Account, error) {
	// get account from shuffled accounts
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var account Account
	for _, i := range r.Perm(len(c.accounts)) {
		// 获取账号
		account = c.accounts[i]
		// 验证账号是否可用,可用直接返回，不可用则选择下一个
		resp, err := c.GetScore(account.User, account.Pass)
		if err == nil && resp.Tifen >= 10 { // 最少需要10题分
			return account, nil
		}
	}
	return account, errors.New("No avaliable accounts")
}

// GetScore 超级鹰GetScore接口
func (c *Client) GetScore(user, pass string) (*GetScoreResp, error) {
	/* 接口说明:查询用户的题分信息

	   接口网址:http://upload.chaojiying.net/Upload/GetScore.php

	   POST发送模式：application/x-www-form-urlencoded, multipart/form-data, application/json

	   返回格式:json

	   返回编码:utf-8

	   返回汉字编码:Unicode

	   发送说明:
	   user=用户账号
	   pass=用户密码 或 pass2=用户密码的md5值(32位小写)

	   返回说明:
	   err_no(数值) 返回代码;
	   err_str(字符串) 中文描述的返回信息;
	   tifen(数值) 题分;
	   tifen_lock(数值) 锁定题分

	   返回json字符串示例:
	   {"err_no":0, "err_str":"OK", "tifen":821690, "tifen_lock":0}
	*/
	apiURL := "http://upload.chaojiying.net/Upload/GetScore.php"
	data := url.Values{
		"user": {user},
		"pass": {pass},
	}
	r := &GetScoreResp{}
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return r, errors.Wrap(err, "chaojiying GetScore PostForm error")
	}
	if resp.StatusCode != 200 {
		return r, errors.New("chaojiying GetScore resp.Status error:" + resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, errors.Wrap(err, "chaojiying GetScore ReadAll error")
	}
	err = json.Unmarshal(body, r)
	if err != nil {
		return r, errors.Wrap(err, "chaojiying GetScore Unmarshal error")
	}
	if r.ErrNo != 0 {
		return r, errors.New("chaojiying New GetScore ErrNo error:" + string(body))
	}
	return r, nil
}

// Processing 超级鹰验证码识别接口
func (c *Client) Processing(user, pass string, pic io.Reader) (*ProcessingResp, error) {
	/*
		接口说明:识别核心接口，一步识别，发送图片和其他相关信息给服务器，服务器返回识别结果和一些其他信息 POST
		接口网址:http://upload.chaojiying.net/Upload/Processing.php

		发送说明:
		user=用户账号
		pass=用户密码 //或 pass2=用户密码的md5值(32位小写)
		softid=软件ID  在用户中心，软件ID处可以生成
		codetype=验证码类型   在价格体系中选用一个适合的类型 点击这里进入
		len_min=最小位数 //默认0为不启用,图片类型为可变位长时可启用这个参数

		以下两个参数选其一  图片文件的宽推荐不超过460px,高不超过310px

		userfile=图片文件二进制流(或是称之为内存流,文件流,字节流的概念)
		file_base64=图片文件base64字符串

		返回说明:
		err_no,(数值) 返回代码
		err_str,(字符串) 中文描述的返回信息
		pic_id,(字符串) 图片标识号，或图片id号
		pic_str,(字符串) 识别出的结果
		md5,(字符串) md5校验值,用来校验此条数据返回是否真实有效 点击这里查看md5校验算法

		返回json字符串示例:{"err_no":0,"err_str":"OK","pic_id":"1662228516102","pic_str":"8vka","md5":"35d5c7f6f53223fbdc5b72783db0c2c0"}
		推荐逻辑处理流程：if (err_no == 0) {识别结果 = pic_str} else {错误代码 = err_no}
	*/
	apiURL := "http://upload.chaojiying.net/Upload/Processing.php"
	picContent, err := ioutil.ReadAll(pic)
	if err != nil {
		return nil, errors.Wrap(err, "chaojiying Processing ReadAll pic error")
	}
	b64 := base64.StdEncoding.EncodeToString(picContent)
	data := url.Values{
		"user":        {user},
		"pass":        {pass},
		"softid":      {"96002"},
		"codetype":    {"1902"},
		"len_min":     {"0"},
		"file_base64": {b64},
	}
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return nil, errors.Wrap(err, "chaojiying Processing PostForm error")
	}
	if resp.StatusCode != 200 {
		return nil, errors.Wrap(err, "chaojiying Processing resp.Status error:"+resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "chaojiying Processing ReadAll resp.Body error")
	}
	r := &ProcessingResp{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, errors.Wrap(err, "chaojiying Processing Unmarshal error")
	}
	return r, nil
}

// Cr4ck 破解图片验证码
// 封装Proccessing接口，判断是否成功，成功返回破解出的验证码
func (c *Client) Cr4ck(pic io.Reader) (string, error) {
	account, err := c.PickOneAccount()
	logging.Infos(nil, "chaojiying Cr4ck picked account:", account)
	if err != nil {
		return "", errors.Wrap(err, "chaojiying Cr4ck PickOneAccount error")
	}
	resp, err := c.Processing(account.User, account.Pass, pic)
	if err != nil {
		return "", errors.Wrap(err, "chaojiying Cr4ck Processing error")
	}
	if resp.ErrNo != 0 {
		return "", errors.New("chaojiying Cr4ck Proccessing resp.ErrNo != 0:" + resp.ErrStr)
	}
	if resp.PicStr == "" {
		return "", errors.New("chaojiying can't Cr4ck this pic")
	}
	logging.Infos(nil, "chaojiying Cr4ck result:", resp)
	return resp.PicStr, nil
}

// LoadAccountsFromJSONFile 从指定位置的json文件中加载账号
//
// json格式: [{"user": "", "pass": ""}, ...]
func LoadAccountsFromJSONFile(filePath string) ([]Account, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "chaojiying LoadAccountsFromJSONFile ReadFile error")
	}
	accounts := []Account{}
	if err := json.Unmarshal(b, &accounts); err != nil {
		return nil, errors.Wrap(err, "chaojiying LoadAccountsFromJSONFile Unmarshal error")
	}
	return accounts, nil
}
