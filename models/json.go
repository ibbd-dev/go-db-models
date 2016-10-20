package models

import (
	"encoding/json"
	"io/ioutil"
)

// json配置文件的结构定义
type JsonConf struct {
	Tables []JsonTableConf
}

// 单个数据表的配置
type JsonTableConf struct {
	Name    string   // 表名
	Fields  []string // 字段名
	QueryBy string   // QueryBy函数定义
	Msg     string   // 表结构说明
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
