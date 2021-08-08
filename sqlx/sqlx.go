package main

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type user struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

// sqlx事务操作示例
// 对于事务操作，我们可以使用sqlx中提供的db.Beginx()和tx.Exec()方法
func transactionDemo2() (err error) {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Printf("Begin trans failed, err: %v\n", err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("32: Rollback")
			tx.Rollback()
		} else {
			err := tx.Commit()
			if err != nil {
				fmt.Println("37: Rollback")
				tx.Rollback()
			}
			fmt.Println("commit")
		}
	}()
	sqlStr1 := "UPDATE user SET age = 20 WHERE id = ?"
	result, err := db.Exec(sqlStr1, 1)
	if err != nil {
		fmt.Printf("Update failed, err: %v\n", err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected() failed, err: %v\n", err)
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}

	sqlStr2 := "UPDATE user SET age = 50 WHERE id = ?"
	result, err = db.Exec(sqlStr2, 5)
	if err != nil {
		return err
	}
	n, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr2 failed")
	}
	return err
}

// NamedQuery()示例
// NamedExec(), NamedQuery()方法用来绑定SQL语句中与结构体或map中的同名字段
func namedQuery() {
	sqlStr := "SELECT * FROM user WHERE name = :name"
	// 使用map做命名查询
	rows, err := db.NamedQuery(sqlStr,
		map[string]interface{}{
			"name": "hello",
		})
	if err != nil {
		fmt.Printf("Query failed, err: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("struct scan failed, err: %v\n", err)
			return
		}
		fmt.Printf("user: %#v\n", u)
	}

	// 使用结构体命名查询，根绝结构体字段的db tag进行映射
	u := user{
		Name: "hello",
	}
	rows, err = db.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("db.NamedQuery() failed, err: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("rows.StructScan() failed, err: %v\n", err)
			continue
		}
		fmt.Printf("user: %#v\n", u)
	}
}

// NamedExec(), NamedQuery()方法用来绑定SQL语句中与结构体或map中的同名字段
// NamedExec()示例
func insertUserDemo() {
	sqlStr := "INSERT INTO user(name, age) VALUES(:name, :age)"
	result, err := db.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "hello",
			"age":  18,
		})
	if err != nil {
		fmt.Printf("Insert failed, err: %v\n", err)
		return
	}
	theId, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get last insert id failed, err: %v\n", err)
		return
	}
	fmt.Printf("Insert success, lasteInsertId: %d\n", theId)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "DELETE FROM user WHERE id = ?"
	result, err := db.Exec(sqlStr, 4)
	if err != nil {
		fmt.Printf("Delete failed, err: %v\n", err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get rowsaffected failed, err: %v\n", err)
		return
	}
	fmt.Printf("Delete success, rowsAffected is: %d\n", n)
}

// 数据更新示例
func updateRowDemo() {
	sqlStr := "UPDATE user SET age = ? WHERE id = ?"
	result, err := db.Exec(sqlStr, 22, 4)
	if err != nil {
		fmt.Printf("Update failed, err: %v\n", err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed, err: %v\n", err)
		return
	}
	fmt.Printf("Update success, affected rows: %d\n", n)
}

// 数据插入示例
func insertRowDemo() {
	sqlStr := "INSERT INTO user(name, age) VALUES(?, ?)"
	result, err := db.Exec(sqlStr, "hhh", 26)
	if err != nil {
		fmt.Printf("Insert into failed, err: %v\n", err)
		return
	}
	theId, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get lastInsertId failed, err: %v\n", err)
		return
	}
	fmt.Printf("Insert Success, the insert id is: %d\n", theId)
}

// 多行查询示例
func queryMultiRowDemo() {
	sqlStr := "SELECT id, name, age FROM user WHERE id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("Query failed, err: %v\n", err)
		return
	}
	fmt.Printf("users: %#v\n", users)
}

// 单行数据查询示例
func queryRowDemo() {
	sqlStr := "SELECT id, name, age FROM user WHERE id = ?"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("db.Get() failed, err: %v\n", err)
		return
	}
	fmt.Printf("id: %d\nname: %s\nage: %d\n", u.Id, u.Name, u.Age)
}

// 连接数据库示例
func initDB() (err error) {
	dsn := "root:@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("Connect DB failed, err: %v\n", err)
		return
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	return
}

func main() {
	if err := initDB(); err != nil {
		fmt.Printf("InitDB failed, err: %v\n", err)
		return
	} else {
		fmt.Println("Connect DB success!!!")
	}
	// queryRowDemo()
	// insertRowDemo()
	updateRowDemo()
	// queryMultiRowDemo()
	// deleteRowDemo()
	// insertUserDemo()
	queryMultiRowDemo()
	namedQuery()
	transactionDemo2()
}
