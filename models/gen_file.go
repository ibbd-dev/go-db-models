package models

import (
	"os"
	"strings"
	"text/template"
)

func GenCommonFile(package_name string) {
	out_file := "common_gen.go"
	fout, err := os.Create(out_file)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	code, err := template.New("gomodels-common").Parse(commonCodeTemplate)
	if err != nil {
		panic(err)
	}

	var data = struct {
		PackageName string
	}{
		PackageName: package_name,
	}

	if err := code.Execute(fout, data); err != nil {
		panic(err)
	}
}

func GenFile(table ParseTable) {
	out_file := table.Name + "_tb_gen.go"
	fout, err := os.Create(out_file)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	fnMap := template.FuncMap{
		"Format2StructName": Format2StructName,
		"Format2StructTag":  Format2StructTag,
		"Format2Title":      Format2Title,
	}
	code, err := template.New("gomodels").Funcs(fnMap).Parse(codeTemplate)
	if err != nil {
		panic(err)
	}

	if err := code.Execute(fout, table); err != nil {
		panic(err)
	}
}

// 将下划线分割的字符串改为驼峰格式的字符串
// 如：hello_world => HelloWorld
func Format2StructName(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	title_str := strings.Title(str)
	return strings.Replace(title_str, " ", "", -1)
}

// 首字母小写
func Format2Title(str string) string {
	str = Format2StructName(str)
	return strings.ToLower(str[0:1]) + str[1:]
}

// 增加重音符在字符串的前后
// sql的字段等都需要
func AddBackquote(str string) string {
	return "`" + str + "`"
}

// 格式化输出的结构体的Tag标签
// 增加重音符及name标签
func Format2StructTag(str string) string {
	return "`name:\"" + str + "\"`"
}
