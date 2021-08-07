package main

import (
        "database/sql"
        "fmt"

        _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() (err error) {
        dsn := "root:@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
        // 不会校验账号密码是否正确
        // 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
        db, err = sql.Open("mysql", dsn)
        if err != nil {
                panic(err)
        }
        // 尝试与DB建立连接（校验dsn是否正确）
        err = db.Ping()
        if err != nil {
                return err
        }
        return err
}

func main() {
        err := initDB()
        if err != nil {
                fmt.Printf("Init db failed, err: %v\n", err)
                return
        }
        defer db.Close()
        fmt.Printf("Connect to db success...")
}
