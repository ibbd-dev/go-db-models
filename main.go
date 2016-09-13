package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/ibbd-dev/go-db-models/models"
)

func main() {
	log.SetFlags(0)

	// 命令参数设置
	database := flag.String("d", "", "")
	charset := flag.String("c", "utf8", "")
	username := flag.String("u", "", "")
	password := flag.String("p", "", "")
	packagename := flag.String("n", "", "")
	host := flag.String("o", "127.0.0.1", "")
	port := flag.Int("r", 3306, "")
	version := flag.Bool("v", false, "")
	help := flag.Bool("h", false, "")
	flag.StringVar(database, "database", "", "")
	flag.StringVar(charset, "charset", "utf8", "")
	flag.StringVar(username, "username", "", "")
	flag.StringVar(password, "password", "", "")
	flag.StringVar(host, "host", "127.0.0.1", "")
	flag.StringVar(packagename, "package", "", "")
	flag.IntVar(port, "port", 3306, "")
	flag.BoolVar(version, "version", false, "")
	flag.BoolVar(help, "help", false, "")
	flag.Usage = func() { log.Println(usageText) } // call on flag error
	flag.Parse()

	//if debug {
	fmt.Println("**********************************")
	fmt.Printf("database: %s\n", *database)
	fmt.Printf("charset: %s\n", *charset)
	fmt.Printf("username: %s\n", *username)
	fmt.Printf("password: %s\n", *password)
	fmt.Printf("host: %s\n", *host)
	fmt.Printf("port: %d\n", *port)
	fmt.Println("**********************************")
	//}

	if *help {
		// not an error, send to stdout
		// that way people can: scaneo -h | less
		fmt.Println(usageText)
		return
	}

	if len(*packagename) == 0 {
		fmt.Println("packagename is empty! use -n or --package")
		fmt.Println(usageText)
		return
	}

	if *version {
		fmt.Println(versionText)
		return
	}

	db_conf := &models.DbConf{
		Host:     *host,
		Port:     *port,
		DbName:   *database,
		UserName: *username,
		Password: *password,
	}

	tables := db_conf.ShowTables()
	//fmt.Println(tables)

	_ = models.ParseTablesStruct(tables, *packagename)
	//fmt.Println(ptables)

	runFmt()
	fmt.Println("\nAll is ok!")
}

func runFmt() {
	in := bytes.NewBuffer(nil)
	cmd := exec.Command("sh")
	cmd.Stdin = in
	go func() {
		in.WriteString("go fmt\n")
		in.WriteString("exit\n")
	}()

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
