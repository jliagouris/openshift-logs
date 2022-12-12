package main

/*
func main() {
	db, err := sql.Open("mysql", "push_app:JesusS980721!@/log_database")
	checkErr(err)

	//插入数据
	//stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	//checkErr(err)

	//res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	//checkErr(err)

	//id, err := res.LastInsertId()
	//checkErr(err)

	//fmt.Println(id)
	////更新数据
	//stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	//checkErr(err)

	//res, err = stmt.Exec("astaxieupdate", id)
	//checkErr(err)

	//affect, err := res.RowsAffected()
	//checkErr(err)

	//fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM test_table")
	checkErr(err)

	for rows.Next() {
		var id int
		var timestamp int
		var value int
		err = rows.Scan(&id, &timestamp, &value)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(timestamp)
		fmt.Println(value)
	}

	//删除数据
	//stmt, err = db.Prepare("delete from userinfo where uid=?")
	//checkErr(err)

	//res, err = stmt.Exec(id)
	//checkErr(err)

	//affect, err = res.RowsAffected()
	//checkErr(err)

	//fmt.Println(affect)

	err = db.Close()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
*/
