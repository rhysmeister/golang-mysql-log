package main

import (
	"database/sql"
	"fmt"
	"time"

	mysql "github.com/go-sql-driver/mysql"

	"flag"
	"golang-mysql-log/mysqllib"

	"bufio"
	"os"
)

var mysqlUser string
var mysqlPasswd string
var mysqlHost string
var mysqlPort int
var mysqlDB string
var createDB bool
var log string
var pingDB bool
var debug bool

func init() {
	flag.StringVar(&mysqlUser, "u", "go", "MySQL username")
	flag.StringVar(&mysqlPasswd, "p", "", "MySQL password")
	flag.StringVar(&mysqlHost, "h", "127.0.0.1", "MySQL host")
	flag.IntVar(&mysqlPort, "P", 3306, "MySQL Port")
	flag.StringVar(&mysqlDB, "D", "golang", "MySQL Database to log to")
	flag.BoolVar(&createDB, "createDB", false, "Create schema in MySQL DB and exit")
	flag.StringVar(&log, "log", "", "Message to log in database")
	flag.BoolVar(&pingDB, "pingDB", false, "Test the MySQL DB connection and exit")
	flag.BoolVar(&debug, "debug", false, "Outpuzt debug info")
}

func main() {
	flag.Parse()
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = mysqlUser
	mysqlConfig.Passwd = mysqlPasswd
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = fmt.Sprintf("%s:%d", mysqlHost, mysqlPort)

	if createDB {
		db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
		if debug {
			fmt.Println(mysqlConfig.FormatDSN())
		}
		if err != nil {
			panic(err)
		}
		defer db.Close()
		mysqllib.CreateDatabase(db, mysqlDB)
		os.Exit(0)
	} else if pingDB {
		db, _ := sql.Open("mysql", mysqlConfig.FormatDSN())
		mysqllib.TestDatabaseConnection(db)
		os.Exit(0)
	}
	mysqlConfig.DBName = mysqlDB
	// sql.Open does not actually open a connection. It just returns a handle
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Check if there is something to read on stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin := []string{}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin = append(stdin, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		if debug {
			fmt.Printf("stdin = %s\n", stdin)
		}
		mysqllib.InsertLog(db, mysqlDB, stdin)
	} else {
		if log != "" {
			if debug {
				fmt.Printf("log = %s\n", log)
			}
			mysqllib.InsertLog(db, mysqlDB, []string{log})
		} else {
			panic("A string message to log should be provided on stdin or with the -log flag")
		}
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}
