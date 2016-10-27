package main

const (
	versionText = "go-db-models v1.1"

	usageText = `GOMODELSSCAN
    根据定义从数据库生成定义结构，每个数据表生成一个对应go文件，包括数据结构定义，及两种基本的查询功能，文件名和数据表名对应。如果json-filepath为空，则生成所有的数据表的结构体。

USAGE
    go-db-models [options] [json-filepath]

OPTIONS
    -d, -database
        数据库名字
    -c, -charset
        数据库的编码，默认为：utf8 (暂时用不上，默认都是utf8)
    -u, -user
        数据库的用户名
    -p, -password
        数据库密码
    -h, -host
        数据库的host，默认为：127.0.0.1
    -r, -port
        数据库的端口号，默认为：3306
    -n, -package
        生成的go文件的包名
    -v, -version
        Print version and exit.
    -e, -help
        Print help and exit.

EXAMPLES
    go-db-models -h host -d database_name -u username -p password 

INSTALL
    go get -u github.com/ibbd-dev/go-db-models
`
)
