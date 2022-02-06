package cmd

import (
	"crypto/md5"
	"encoding/json"
	"et/config"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

type Tbsign struct {
}

type TbsBasket struct {
	Tbs string
}
type TbLikes struct {
	Time       int
	Error_code string
	Forum_list ForumList
}
type ForumList struct {
	Non_gconforum []NonGconforum `json:"non-gconforum"`
}
type NonGconforum struct {
	Id   string
	Name string
}

func init() {
	rootCmd.AddCommand(tbsignCmd)
}

var tbsignCmd = &cobra.Command{
	Use:   "tbsign",
	Short: "贴吧签到",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		var cls Tbsign
		cls.run(cls)
	},
}

var wg sync.WaitGroup

func (Tbsign) run(cls Tbsign) {
	info := config.BaseInfo{}
	conf := info.GetConf()
	client := resty.New()

	// 获取tbs
	tbs := cls.getTbs(conf.TbConfig.Bduss)
	// 获取贴吧列表
	nonGconforums := cls.getLikes(conf.TbConfig.Bduss, client)
	// 签到
	wg.Add(len(nonGconforums))
	for _, v := range nonGconforums {
		go cls.sign(v, conf.TbConfig.Bduss, client, tbs)
	}
	wg.Wait()
}

func (Tbsign) sign(item NonGconforum, bduss string, client *resty.Client, tbs string) {
	headers := map[string]string{
		"Host":            "c.tieba.baidu.com",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Charset":         "UTF-8",
		"net":             "3",
		"User-Agent":      "bdtb for Android 8.4.0.1",
		"Connection":      "Keep-Alive",
		"Accept-Encoding": "gzip",
		"cookie":          "ca=open",
	}
	query := client.R().SetHeaders(headers)

	str := "BDUSS=" + bduss + "fid=" + item.Id + "kw=" + item.Name + "tbs=" + tbs + "tiebaclient!!!"
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	url := "http://c.tieba.baidu.com/c/c/forum/sign?BDUSS=" + url.QueryEscape(bduss) + "&fid=" + item.Id + "&kw=" + url.QueryEscape(item.Name) + "&sign=" + sign + "&tbs=" + tbs
	resp, _ := query.Get(url)

	type Res struct {
		Error_code string
		Error_msg  string
	}
	var res Res
	_ = json.Unmarshal(resp.Body(), &res)
	fmt.Println("[" + item.Name + "]" + res.Error_msg)
	defer wg.Done()
}

func (Tbsign) getLikes(bduss string, client *resty.Client) []NonGconforum {

	headers := map[string]string{
		"Host":            "c.tieba.baidu.com",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Charset":         "UTF-8",
		"net":             "3",
		"User-Agent":      "bdtb for Android 8.4.0.1",
		"Connection":      "Keep-Alive",
		"Accept-Encoding": "gzip",
		"cookie":          "ca=open",
	}

	query := client.R().SetHeaders(headers)

	// 获取关注吧列表，如果超过100个需要分页
	str := "BDUSS=" + bduss + "_client_version=8.1.0.4" + "page_no=1" + "page_size=100" + "tiebaclient!!!"
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	resp, err := query.Get("http://c.tieba.baidu.com/c/f/forum/like?BDUSS=" + bduss + "&_client_version=8.1.0.4&page_no=1&page_size=100&sign=" + sign)

	if err != nil {
		log.Fatal(err)
	}
	var tbLikes TbLikes
	err = json.Unmarshal(resp.Body(), &tbLikes)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("本次获取到%d个关注贴吧\n", len(tbLikes.Forum_list.Non_gconforum))

	return tbLikes.Forum_list.Non_gconforum

}

func (Tbsign) getTbs(bduss string) string {

	client := resty.New()

	headers := map[string]string{
		"Host":            "tieba.baidu.com",
		"User-Agent":      "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
		"Referer":         "http://tieba.baidu.com/",
		"Connection":      "keep-alive",
		"cookie":          "BDUSS=" + bduss,
	}

	resp, err := client.R().
		SetHeaders(headers).
		Post("http://tieba.baidu.com/dc/common/tbs")

	if err != nil {
		log.Fatal(err)
	}
	var tbsBasket TbsBasket
	err = json.Unmarshal(resp.Body(), &tbsBasket)

	if err != nil {
		log.Fatal(err)
	}
	return tbsBasket.Tbs
}
