package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//解析yml文件
type BaseInfo struct {
	V2Config V2Config `yaml:"v2config"`
	TbConfig TbConfig `yaml:"tbconfig"`
}

type TbConfig struct {
	Bduss string `yaml:"bduss"`
}

type V2Config struct {
	Cookie string `yaml:"cookie"`
	Proxy  string `yaml:"proxy"`
}

func (c *BaseInfo) GetConf() *BaseInfo {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
	}
	yamlFile, err := ioutil.ReadFile(dir + "/.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}
