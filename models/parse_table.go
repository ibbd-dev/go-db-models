package models

import (
	"strings"
)

const (
	// 可能需要额外引入的包
	packageSql  string = "database/sql"
	packageTime string = "time"
)

type ParseTable struct {
	Name        string   // 数据表名，如：ad_plan，对应结构体名为：AdPlanTable，对应文件名为：gen_ad_plan.go
	PackageName string   // 生成的程序的包名
	PrimaryType string   // 主键的类型，如：uint32, sql.NullString等
	Imports     []string // 需要import的包
	Fields      []ParseField
}

type ParseField struct {
	Name string // 字段名，如：plan_id。对应到结构体的属性名就是：PlanId
	Type string // 字段类型，对应golang中的类型，如：uint32, sql.NullString
}

// 解释数据表的结构体
func ParseTablesStruct(tables []Table, package_name string, models_conf *JsonConf) (parse_tables []ParseTable) {
	// 配置的预处理
	var modelsConfMap = map[string]map[string]bool{}
	if len(models_conf.Tables) > 0 {
		for _, tb := range models_conf.Tables {
			modelsConfMap[tb.Name] = map[string]bool{}
			for _, f := range tb.Fields {
				modelsConfMap[tb.Name][f] = true
			}
		}
	}

	for _, table := range tables {
		if len(models_conf.Tables) > 0 && modelsConfMap[table.Name] == nil {
			continue
		}

		ptable := ParseTable{
			Name:        table.Name,
			PackageName: package_name,
		}
		ptable.Fields, ptable.Imports, ptable.PrimaryType = ParseFieldsStruct(table.Fields, modelsConfMap[table.Name])

		// 生成代码文件
		GenFile(ptable)

		parse_tables = append(parse_tables, ptable)
	}

	return parse_tables
}

// 解释一个数据表的所有字段
func ParseFieldsStruct(fields []Field, fields_conf map[string]bool) (pfields []ParseField, imports []string, primary_type string) {
	for _, f := range fields {
		if fields_conf != nil && fields_conf[f.Name] == false {
			continue
		}

		pf := ParseField{
			Name: f.Name,
		}

		if isString(f.Type) {
			// 字符串
			if f.Null == "YES" {
				pf.Type = "sql.NullString"
				imports = importsPush(imports, packageSql)
			} else {
				pf.Type = "string"
			}
		} else if strings.Contains(f.Type, "int") {
			// 整型
			if f.Null == "YES" {
				pf.Type = "sql.NullInt64"
				imports = importsPush(imports, packageSql)
			} else {
				prefix := ""
				if strings.Contains(f.Type, "unsigned") {
					prefix = "u"
				}

				int_type := "int32"
				if strings.Contains(f.Type, "tinyint") {
					int_type = "int8"
				} else if strings.Contains(f.Type, "smallint") {
					int_type = "int16"
				} else if strings.Contains(f.Type, "bigint") {
					int_type = "int64"
				}
				pf.Type = prefix + int_type
			}
		} else if strings.Contains(f.Type, "float") {
			// 浮点数
			if f.Null == "YES" {
				pf.Type = "sql.NullFloat64"
				imports = importsPush(imports, packageSql)
			} else {
				pf.Type = "float"
			}
		} else if strings.Contains(f.Type, "double") || strings.Contains(f.Type, "decimal") {
			// 高精度浮点数
			if f.Null == "YES" {
				pf.Type = "sql.NullFloat64"
				imports = importsPush(imports, packageSql)
			} else {
				pf.Type = "float64"
			}
		} else if strings.Contains(f.Type, "year") {
			// 年份
			if f.Null == "YES" {
				pf.Type = "sql.NullInt64"
				imports = importsPush(imports, packageSql)
			} else {
				pf.Type = "uint16"
			}
		} else if strings.Contains(f.Type, "datetime") || strings.Contains(f.Type, "timestamp") || strings.Contains(f.Type, "date") {
			// 日期时间
			pf.Type = "time.Time"
			imports = importsPush(imports, packageTime)
		}

		pfields = append(pfields, pf)
	}

	return pfields, imports, primary_type
}

// 判断是否是字符串
func isString(field_type string) bool {
	return strings.Contains(field_type, "text") || strings.Contains(field_type, "char") || strings.Contains(field_type, "binary") || strings.Contains(field_type, "blob")
}

func importsPush(imports []string, packagename string) []string {
	is_exists := false
	for _, pn := range imports {
		if pn == packagename {
			is_exists = true
			break
		}
	}

	if !is_exists {
		imports = append(imports, packagename)
	}

	return imports
}
