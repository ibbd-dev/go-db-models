package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 数据库配置
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
	Name    string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
}

func (conf *DbConf) getDb() (*sql.DB, error) {
	// Open database connection
	conn_string := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.UserName, conf.Password, conf.Host, conf.Port, conf.DbName)

	db, err := sql.Open("mysql", conn_string)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 对应sql：show tables
func (conf *DbConf) ShowTables() ([]Table, error) {
	db, err := conf.getDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("show tables")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := []Table{}
	for rows.Next() {
		t := Table{}
		err = rows.Scan(&t.Name)
		if err != nil {
			return nil, err
		}

		t.Fields, err = conf.DescTable(t.Name, db)
		if err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}

	return tables, nil

}

// 对应sql：desc table_name
func (conf *DbConf) DescTable(table_name string, db *sql.DB) ([]Field, error) {
	if db == nil {
		db, err := conf.getDb()
		if err != nil {
			return nil, err
		}
		defer db.Close()
	}

	// Execute the query
	rows, err := db.Query("desc " + table_name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var extra sql.NullString
	fields := []Field{}
	for rows.Next() {
		f := Field{}
		err = rows.Scan(&f.Name, &f.Type, &f.Null, &f.Key, &f.Default, &extra)
		if err != nil {
			return nil, err
		}

		fields = append(fields, f)
	}

	return fields, nil
}
