package models

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	// 可能需要额外引入的包
	//packageSql  string = "database/sql"
	packageTime string = "time"
)

// 经过解释之后的数据表定义，用于生成相应的go文件
type ParseTable struct {
	Name         string   // 数据表名，如：ad_plan，对应结构体名为：AdPlanTable，对应文件名为：gen_ad_plan.go
	PackageName  string   // 生成的程序的包名
	PrimaryType  string   // 主键的类型，如：uint32, string等
	Imports      []string // 需要import的包
	SelectFields string   // sql查询中的select fields
	Msg          string   // 表结构的说明
	Fields       []ParseField
	QueryBy      QueryBy // QueryBy函数，例如QueryById等

	// 如果该字段不为空，则返回多行记录时，是map结构（默认是list结构）
	// 关联json中的MapIndex字段
	MapIndex MapIndexField
}

// 字段的定义
type ParseField struct {
	Name string // 字段名，如：plan_id。对应到结构体的属性名就是：PlanId
	Type string // 字段类型，对应golang中的类型，如：uint32, sql.NullString
}

type MapIndexField struct {
	Name string // map下标的字段名
	Type string // map下标字段的类型
}

// 生成QueryBy函数时需要该结构
type QueryBy struct {
	FieldName string // query by函数的参数名
	FieldType string // query by函数的参数的类型，如uint32等
}

// 解释数据表的结构体
func ParseTablesStruct(tables []Table, packageName string, modelsConf *JsonConf) (parseTables []ParseTable, err error) {
	// 生成common文件
	err = GenCommonFile(packageName)
	if err != nil {
		return nil, err
	}

	// 配置的预处理
	var modelsConfMap = map[string]map[string]bool{} // 第一个下标是表名，第二个下标是字段名
	var queryByConf = map[string]string{}            // 下标是表名，值是字段名，例如id。
	var tableConfMsg = map[string]string{}           // 表的注释（在json文件中的）
	var mapIndexConf = map[string]string{}           // map下标配置
	if len(modelsConf.Tables) > 0 {
		for _, tb := range modelsConf.Tables {
			modelsConfMap[tb.Name] = map[string]bool{}
			for _, f := range tb.Fields {
				modelsConfMap[tb.Name][f] = true
			}

			if tb.QueryBy != "" {
				queryByConf[tb.Name] = tb.QueryBy
			}

			if tb.MapIndex != "" {
				mapIndexConf[tb.Name] = tb.MapIndex
			}

			tableConfMsg[tb.Name] = tb.Msg
		}
	}

	for _, table := range tables {
		if len(modelsConf.Tables) > 0 && modelsConfMap[table.Name] == nil {
			continue
		}

		ptable := ParseTable{
			Name:        table.Name,
			Msg:         tableConfMsg[table.Name],
			PackageName: packageName,
		}
		ptable.Fields, ptable.Imports, ptable.PrimaryType = ParseFieldsStruct(table.Fields, modelsConfMap[table.Name])

		// 拼接查询的字符串
		sep := ""
		for _, f := range ptable.Fields {
			ptable.SelectFields += sep + "`" + f.Name + "`"
			sep = ","
		}

		// 处理query by
		if queryByConf[table.Name] != "" {
			isMatch := false
			for _, f := range ptable.Fields {
				if f.Name == queryByConf[table.Name] {
					isMatch = true
					ptable.QueryBy.FieldName = f.Name
					ptable.QueryBy.FieldType = f.Type
				}
			}

			if isMatch == false {
				// query by的字段不在字段列表里
				return nil, errors.New("ERROR field name: " + queryByConf[table.Name] + " of QueryBy function for table name: " + table.Name)
			}
		}

		// 处理map index
		if mapIndexConf[table.Name] != "" {
			isMatch := false
			for _, f := range ptable.Fields {
				if f.Name == mapIndexConf[table.Name] {
					isMatch = true
					ptable.MapIndex.Name = f.Name
					ptable.MapIndex.Type = f.Type
				}
			}

			if isMatch == false {
				// query by的字段不在字段列表里
				return nil, errors.New("ERROR field name: " + mapIndexConf[table.Name] + " of MapIndex for table name: " + table.Name)
			}
		}

		// 生成代码文件
		err = GenFile(ptable)
		if err != nil {
			return nil, err
		}

		// 生成单元测试文件
		err = GenTestFile(ptable)
		if err != nil {
			return nil, err
		}

		parseTables = append(parseTables, ptable)
	}

	return parseTables, nil
}

// 解释一个数据表的所有字段
func ParseFieldsStruct(fields []Field, fieldsConf map[string]bool) (pfields []ParseField, imports []string, primaryType string) {
	for _, f := range fields {
		if fieldsConf != nil && fieldsConf[f.Name] == false {
			continue
		}

		pf := ParseField{
			Name: f.Name,
		}

		if isString(f.Type) {
			// 字符串
			if f.Null == "YES" {
				pf.Type = "sql.NullString"
				//imports = importsPush(imports, packageSql)
			} else {
				pf.Type = "string"
			}
		} else if strings.Contains(f.Type, "int") {
			// 整型
			if f.Null == "YES" {
				pf.Type = "sql.NullInt64"
				//imports = importsPush(imports, packageSql)
			} else {
				prefix := ""
				if strings.Contains(f.Type, "unsigned") {
					prefix = "u"
				}

				intType := "int32"
				if strings.Contains(f.Type, "tinyint") {
					intType = "int8"
				} else if strings.Contains(f.Type, "smallint") {
					intType = "int16"
				} else if strings.Contains(f.Type, "bigint") {
					intType = "int64"
				}
				pf.Type = prefix + intType
			}
		} else if strings.Contains(f.Type, "float") {
			// 浮点数
			if f.Null == "YES" {
				pf.Type = "sql.NullFloat64"
				//imports = importsPush(imports, packageSql)
			} else {
				pf.Type = "float32"
			}
			//println(f.Name, ": ", f.Type, "===>", pf.Type)
		} else if strings.Contains(f.Type, "double") {
			// 高精度浮点数
			if f.Null == "YES" {
				pf.Type = "sql.NullFloat64"
				//imports = importsPush(imports, packageSql)
			} else {
				pf.Type = "float64"
			}
			//println(f.Name, ": ", f.Type, "===>", pf.Type)
		} else if strings.Contains(f.Type, "decimal") {
			// 高精度浮点数
			if f.Null == "YES" {
				pf.Type = "sql.NullFloat64"
				//imports = importsPush(imports, packageSql)
			} else {
				// 这个时候还需要根据有效位数进行判断
				re := regexp.MustCompile("decimal.*\\(([0-9]+),")
				strs := re.FindStringSubmatch(f.Type)
				if len(strs) != 2 {
					pf.Type = "float64"
				} else {
					// 判断有效位数
					num, err := strconv.Atoi(strings.TrimSpace(strs[1]))
					if err != nil {
						panic(err)
					}

					if num > 6 {
						pf.Type = "float64"
					} else {
						pf.Type = "float32"
					}
				}
				//println(f.Name, ": ", f.Type, "===>", pf.Type)
			}
		} else if strings.Contains(f.Type, "year") {
			// 年份
			if f.Null == "YES" {
				pf.Type = "sql.NullInt64"
				//imports = importsPush(imports, packageSql)
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

	return pfields, imports, primaryType
}

// 判断是否是字符串
func isString(fieldType string) bool {
	return strings.Contains(fieldType, "text") || strings.Contains(fieldType, "char") || strings.Contains(fieldType, "binary") || strings.Contains(fieldType, "blob")
}

// 增加需要import的package
func importsPush(imports []string, packagename string) []string {
	isExists := false
	for _, pn := range imports {
		if pn == packagename {
			isExists = true
			break
		}
	}

	if !isExists {
		imports = append(imports, packagename)
	}

	return imports
}
