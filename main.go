package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type user struct {
	id   int
	age  int
	name string
}

//  事务示例
func transactionDemo() {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		fmt.Printf("Begin trans failed, err: %v\n", err)
		return
	}
	sqlStr1 := "UPDATE user SET age = 30 WHERE id = ?"
	result1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Exec sql1 failed, err: %v\n", err)
		return
	}
	rowAffected1, err := result1.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Printf("Exec result1.RowsAffected() failed, err: %v\n", err)
		return
	}
	sqlStr2 := "UPDATE user SET age = 40 WHERE id = ?"
	result2, err := tx.Exec(sqlStr2, 4)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Exec sql2 failed, err: %v\n", err)
		return
	}
	rowAffedted2, err := result2.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Printf("Exec result2.RowsAffected() failed, err: %v\n", err)
		return
	}
	fmt.Println(rowAffected1, rowAffedted2)
	if rowAffected1 == 1 && rowAffedted2 == 1 {
		fmt.Println("提交事务！！！")
		tx.Commit()
	} else {
		tx.Rollback()
		fmt.Println("事务回滚！！！")
	}
	fmt.Println("Exec trans success！！！")

}

func prepareInsertDemo() {
	sqlStr := "INSERT INTO user(name, age) values(?, ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("Prepare failed, err: %v\n", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("小王子", 19)
	if err != nil {
		fmt.Printf("Insert into failed, err: %v\n", err)
		return
	}
	_, err = stmt.Exec("anlinawang", 27)
	if err != nil {
		fmt.Printf("Insert into failed, err: %v\n", err)
		return
	}
	fmt.Println("Insert success...")
}

// 普通SQL语句执行过程：
// 客户端对SQL语句进行占位符替换得到完整的SQL语句。
// 客户端发送完整SQL语句到MySQL服务端
// MySQL服务端执行完整的SQL语句并将结果返回给客户端。

// 预处理执行过程：
// 把SQL语句分成两部分，命令部分与数据部分。
// 先把命令部分发送给MySQL服务端，MySQL服务端进行SQL预处理。
// 然后把数据部分发送给MySQL服务端，MySQL服务端对SQL语句进行占位符替换。
// MySQL服务端执行完整的SQL语句并将结果返回给客户端。

// 为什么要预处理？
// 优化MySQL服务器重复执行SQL的方法，可以提升服务器性能，提前让服务器编译，一次编译多次执行，节省后续编译的成本。
// 避免SQL注入问题。

// 预处理：查询示例
func prepareQueryDemo() {
	sqlStr := "SELECT id, name, age FROM user WHERE id > ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err: %v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("Query failed, err: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("Scan failed, err: %v\n", err)
			return
		}
		fmt.Printf("id: %d\nname: %s\nage: %d\n", u.id, u.name, u.age)
	}
}

func deleteRowDemo() {
	sqlStr := "DELETE FROM user WHERE id = ?"
	result, err := db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Printf("Delete failed, err: %v\n", err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get rowsAffected failed, err: %v\n", err)
		return
	}
	fmt.Printf("Delete success, affected rows: %d\n", n)
}

// 更新数据示例
func updateRowDemo() {
	sqlStr := "UPDATE user SET age = ? WHERE id = ?"
	result, err := db.Exec(sqlStr, 39, 3)
	if err != nil {
		fmt.Printf("Update failed, err: %v\n", err)
		return
	}
	// 操作影响的行数
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed, err: %v\n", err)
		return
	}
	fmt.Printf("Update success, affected rows: %d\n", n)
}

// 插入数据示例
func insertRowDemo() {
	sqlStr := "INSERT INTO user(name, age) VALUES(?, ?)"
	result, err := db.Exec(sqlStr, "王五", 38)
	if err != nil {
		fmt.Printf("Insert failed, err: %v\n", err)
		return
	}
	theId, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get lastInsertId failed, err: %v\n", err)
		return
	}
	fmt.Printf("Insert success, the id is %d\n", theId)
}

// 多行查询示例
func queryMultiRowDemo() {
	sqlStr := "SELECT id, name, age FROM user WHERE id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("Query failed, err: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err: %v\n", err)
			return
		}
		fmt.Printf("id: %d\nname: %s\nage:%d\n", u.id, u.name, u.age)
	}
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "SELECT id, name, age FROM user WHERE id=?"
	var u user
	// ！！！确保QueryRow之后调用Scan方法，否则持有的数据库连接不会被释放
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err: %v\n", err)
		return
	}
	fmt.Printf("id: %d\nname: %s\nage: %d\n", u.id, u.name, u.age)
}

// DB初始化示例
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
	// 设置数值大小，需要根据业务场景而定
	db.SetConnMaxLifetime(time.Second * 10)
	// 最大连接数
	db.SetMaxOpenConns(200)
	// 最大空闲连接数
	db.SetMaxIdleConns(10)

	return
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("Init db failed, err: %v\n", err)
		return
	}
	defer db.Close()
	fmt.Printf("Connect to db success...\n")
	queryRowDemo()
	queryMultiRowDemo()
	updateRowDemo()
	deleteRowDemo()
	prepareQueryDemo()
	transactionDemo()
}
