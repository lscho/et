package cmd

import (
	"et/config"
	"fmt"
	"log"
	"strings"

	"regexp"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

type v2sign struct {
}

func init() {
	rootCmd.AddCommand(v2signCmd)
}

var v2signCmd = &cobra.Command{
	Use:   "v2sign",
	Short: "v2ex签到",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		var cls v2sign
		cls.run()
	},
}

func (v2sign) run() {
	info := config.BaseInfo{}
	conf := info.GetConf()
	client := resty.New()

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.57 Safari/537.36 OPR/40.0.2308.15 (Edition beta)",
		"Referer":    "",
		"Origin":     "https://www.v2ex.com",
		"cookie":     conf.V2Config.Cookie,
	}

	query := client

	if conf.V2Config.Proxy != "" {
		query.SetProxy(conf.V2Config.Proxy)
	}

	resp, err := query.R().
		SetHeaders(headers).
		Get("https://www.v2ex.com/mission/daily")

	if err != nil {
		log.Fatal(err)
	}

	body := string(resp.Body())

	signReg := regexp.MustCompile(`\/mission\/daily\/redeem\?once=\d+`)
	if signReg == nil {
		fmt.Println("regexp err")
		return
	}
	sign := signReg.FindString(body)

	if sign == "" {
		labelReg := regexp.MustCompile(`已连续登录\s\d+\s天`)
		if labelReg == nil {
			fmt.Println("regexp err")
			return
		}
		label := labelReg.FindString(body)
		fmt.Println(label)
		return
	}
	fmt.Println(sign)

	signResp, errs := query.R().
		SetHeaders(headers).
		Get("https://www.v2ex.com" + sign)

	if errs != nil {
		log.Fatal(err)
	}
	signBody := string(signResp.Body())

	if strings.Contains(signBody, "已成功领取每日登录奖励") {
		fmt.Println("签到成功")
	} else {
		fmt.Println("签到失败")
	}
}
