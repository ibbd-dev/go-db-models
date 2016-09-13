package models

import (
	"os"
	"strings"
	"text/template"
)

func GenFile(table ParseTable) {
	out_file := "gen_" + table.Name + ".go"
	fout, err := os.Create(out_file)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	code, err := template.New("gomodels").Parse(codeTemplate)
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
