package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DbConf struct {
	Host     string
	Port     int
	DbName   string
	UserName string
	Password string
	Charset  string
}

// 数据表
type Table struct {
	Name   string
	Fields []Field
}

// 字段定义
type Field struct {
	Field   string
	Type    string
	Null    string
	Default sql.NullString
}

func (conf *DbConf) getDb() *sql.DB {
	// Open database connection
	conn_string := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.UserName, conf.Password, conf.Host, conf.Port, conf.DbName)

	db, err := sql.Open("mysql", conn_string)
	if err != nil {
		panic(err.Error())
	}

	return db
}

// 对应sql：show tables
func (conf *DbConf) ShowTables() []Table {
	db := conf.getDb()
	defer db.Close()

	// Execute the query
	rows, err := db.Query("show tables")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	tables := []Table{}
	for rows.Next() {
		t := Table{}
		err = rows.Scan(&t.Name)
		if err != nil {
			panic(err.Error())
		}

		t.Fields = conf.DescTable(t.Name, db)
		tables = append(tables, t)
	}

	return tables

}

// 对应sql：desc table_name
func (conf *DbConf) DescTable(table_name string, db *sql.DB) []Field {
	if db == nil {
		db := conf.getDb()
		defer db.Close()
	}

	// Execute the query
	rows, err := db.Query("desc " + table_name)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var tmp sql.NullString
	fields := []Field{}
	for rows.Next() {
		f := Field{}
		err = rows.Scan(&f.Field, &f.Type, &f.Null, &tmp, &f.Default, &tmp)
		if err != nil {
			panic(err.Error())
		}

		fields = append(fields, f)
	}

	return fields
}
