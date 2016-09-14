package models

import (
	"encoding/json"
	"io/ioutil"
)

// json配置文件的结构定义
type JsonConf struct {
	Tables []struct {
		Name   string
		Fields []string
	}
}

// 将json文件decode成结构体
func JsonUnmarshal(filename string) (json_conf *JsonConf) {
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonBlob, &json_conf)
	if err != nil {
		panic(err)
	}

	return json_conf
}
