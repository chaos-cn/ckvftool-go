package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func PerformanceTest(datasource, database, tableName string, iNum, gNum, sNum, ptNum int) {

	InitDB(datasource)
	defer CloseDB()

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}
	beginTime := time.Now()

	var wgtop sync.WaitGroup
	for i := 0; i < ptNum; i++ {
		wgtop.Add(1)
		go TestDbWorkerGroup(&wgtop, database, tableName, iNum, gNum, sNum)
	}
	wgtop.Wait()

	fmt.Println("PerformanceTest time total used:", time.Since(beginTime))
}

func TestDbWorkerGroup(wgtop *sync.WaitGroup, database, tableName string, iNum, gNum, sNum int) {
	defer wgtop.Done()

	var wg sync.WaitGroup

	for i := 1; i <= gNum; i++ {
		wg.Add(1)
		go TestDbWorker(&wg, countSql(database, tableName, i))
	}

	wg.Add(1)
	go TestDbWorker(&wg, sumSql(database, tableName, sNum))

	wg.Add(1)
	go TestDbWorker(&wg, "select count(*) as \"cnt\" from "+database+"."+tableName)

	wg.Add(1)
	go TestDbWorker(&wg, selectSql(database, tableName, iNum))

	//业务逻辑结束
	wg.Wait()
}

func TestDbWorker(wg *sync.WaitGroup, sql string) {
	defer wg.Done()

	//fmt.Println(sql)
	result := QuerySql(sql)
	log.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓\n",
		sql, "\n", result,
		"\n↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑\n")
}

func selectSql(database, tableName string, iNum int) string {
	sql := strings.Builder{}
	sql.WriteString("select dd, uid ")
	var columnName string
	for i := 1; i <= iNum; i++ {
		columnName = "i_" + strconv.Itoa(10000+i)
		sql.WriteString(", ")
		sql.WriteString(columnName)
	}
	sql.WriteString(", mergedata from ")
	sql.WriteString(database)
	sql.WriteString(".")
	sql.WriteString(tableName)
	sql.WriteString(" limit 0,10")

	return sql.String()
}

func countSql(database, tableName string, gNumIndex int) string {
	sql := strings.Builder{}
	sql.WriteString("select ")
	columnName := "g_" + strconv.Itoa(10000+gNumIndex)
	sql.WriteString(columnName)
	sql.WriteString(",count(uid) as \"cnt\" from ")
	sql.WriteString(database)
	sql.WriteString(".")
	sql.WriteString(tableName)
	sql.WriteString(" group by ")
	sql.WriteString(columnName)
	return sql.String()
}

func sumSql(database, tableName string, sNum int) string {
	sql := strings.Builder{}
	sql.WriteString("select ")
	var columnName string
	for i := 1; i <= sNum; i++ {
		columnName = "s_" + strconv.Itoa(10000+i)
		sql.WriteString("sum(")
		sql.WriteString(columnName)
		sql.WriteString(")")

		if i == sNum {
			sql.WriteString(" as \"sum_")
			sql.WriteString(strconv.Itoa(10000 + i))
			sql.WriteString("\" ")
		} else {
			sql.WriteString(" as \"sum_")
			sql.WriteString(strconv.Itoa(10000 + i))
			sql.WriteString("\" ,")
		}
	}
	sql.WriteString(" from ")
	sql.WriteString(database)
	sql.WriteString(".")
	sql.WriteString(tableName)

	return sql.String()
}
