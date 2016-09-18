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
func JsonUnmarshal(filename string) (jsonConf *JsonConf, err error) {
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBlob, &jsonConf)
	if err != nil {
		return nil, err
	}

	return jsonConf, nil
}
