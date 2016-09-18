# go-db-models

根据配置参数自动创建所有数据表的struct结构，每个结构独立为一个文件。

## INSTALL

```sh
go get github.com/ibbd-dev/go-db-models
```

## USAGE

```sh
# 生成struct
go-db-models -h host -d database -u username -p password -n packageName {json-file.json}

# 帮助
go-db-models
```

生成的文件:

- 公共配置文件`common_gen_go`
- 表结构及查询文件`*_tb_gen.go`, 主要包含以下几部分内容：
  - 表结构定义
  - 单行结果查询函数
  - 根据某个字段查询单行记录的函数，例如根据id字段进行查询
  - 多行结果查询函数

`*_tb_gen.go`文件样例如下：

```go

// DON'T EDIT *** generated by go-db-models *** DON'T EDIT //
package test

import "strings"
import "database/sql"

// 对应数据表：ad_project
type AdProjectTable struct {
	Id     uint32         `name:"id"`
	PlanId uint32         `name:"plan_id"`
	Name   string         `name:"name"`
	Status uint8          `name:"status"`
	Remark sql.NullString `name:"remark"`
}

const adProjectSelectFields string = "`id`,`plan_id`,`name`,`status`,`remark`"

// 查询单行记录（根据某个字段）
func AdProjectQueryById(db *sql.DB, idVal uint32) (ad_project *AdProjectTable, err error) {
	queryString := "SELECT " + adProjectSelectFields + " FROM `ad_project` WHERE `id` = ?"
	ad_project = &AdProjectTable{}
	err = db.QueryRow(queryString, idVal).Scan(
		&ad_project.Id,
		&ad_project.PlanId,
		&ad_project.Name,
		&ad_project.Status,
		&ad_project.Remark,
	)

	if err != nil {
		return nil, err
	}

	return ad_project, nil
}

// 查询单行记录
func AdProjectQueryRow(db *sql.DB, queryString string) (ad_project *AdProjectTable, err error) {
	queryString = strings.Replace(queryString, SelectFieldsTemp, adProjectSelectFields, 1)
	ad_project = &AdProjectTable{}
	err = db.QueryRow(queryString).Scan(
		&ad_project.Id,
		&ad_project.PlanId,
		&ad_project.Name,
		&ad_project.Status,
		&ad_project.Remark,
	)

	if err != nil {
		return nil, err
	}

	return ad_project, nil
}

// 查询多行记录
func AdProjectQuery(db *sql.DB, queryString string) (ad_project []*AdProjectTable, err error) {
	queryString = strings.Replace(queryString, SelectFieldsTemp, adProjectSelectFields, 1)
	rows, err := db.Query(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oneRow = &AdProjectTable{}
		err = rows.Scan(
			&oneRow.Id,
			&oneRow.PlanId,
			&oneRow.Name,
			&oneRow.Status,
			&oneRow.Remark,
		)
		if err != nil {
			return nil, err
		}

		ad_project = append(ad_project, oneRow)
	}

	return ad_project, nil
}
```

## Json数据表定义文件

如果只是希望生成部分的数据表，则可以使用该文件。该文件的结构如：

```json
{
    // 需要生成的数据表
    "tables": [
        {
            // 数据表的名字
            "name": "ad_plan",
            // 需要生成的字段
            "fields": ["id", "name", "status", "daily_budget", "start_date", "created_at"],
            // 这是可选项，如果定义了该字段，则会自动生成一个类似QueryById的查询函数。
            "queryBy": "id"
        },
        {
            "name": "ad_project",
            "fields": ["id", "plan_id", "name", "status", "remark"]
        }
    ]
}
```

## TODO

- 查询缓存
- 连表查询

