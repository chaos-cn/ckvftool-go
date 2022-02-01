package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func InitTableAndData(datasource, database, tablename string, iNum, gNum, sNum, rowNum int) {

	InitDB(datasource)
	defer CloseDB()

	//删除可能存在的已有的表
	ExceSql("DROP TABLE IF EXISTS " + database + "." + tablename)

	//建表
	createTableSql := createTableSql(database, tablename, iNum, gNum, sNum)
	ExceSql(createTableSql)
	log.Println("建表完成:", createTableSql)

	log.Println("正在生成模拟数据....")

	//入库
	preInsertTableSql := perInsertTableSql(database, tablename, iNum, gNum, sNum)
	log.Println("准备入库sql:", preInsertTableSql)

	//生成模拟数据并入库
	moreDataAndBatchInsert(preInsertTableSql, iNum, gNum, sNum, rowNum)

	//最后
	log.Printf("用于测试的表[%s.%s]已初始化完成，数据量为:[%d]", database, tablename, rowNum)
}

//dd	用于描述数据分区，一般取当天日期
//uid	主键，不重复，用于去重和coun计数，一般可用“证件号_证件类型”
//i_type	一般的用于索引的列，用于select及放在where后面
//g_type	用于分组的列，放在group by后面，一般类型有限，尽量不超过100种，如证件类型、籍贯等
//s_type	用于求和的列，用于sum函数，类型必须是int，取值只能是0或1
//data	用于返回字段，存储除了dd、uid和i_type以外的所有内容，内容用符号|拼接
func createTableSql(database string, tableName string, iNum, gNum, sNum int) string {
	strs := strings.Builder{}
	strs.Grow(1000)
	strs.WriteString("CREATE TABLE IF NOT EXISTS ")
	strs.WriteString(database)
	strs.WriteString(".")
	strs.WriteString(tableName)
	strs.WriteString("( dd date,uid String")

	for i := 1; i <= iNum; i++ {
		strs.WriteString(", i_")
		strs.WriteString(strconv.Itoa(10000 + i))
		strs.WriteString(" String")
	}
	for i := 1; i <= gNum; i++ {
		strs.WriteString(", g_")
		strs.WriteString(strconv.Itoa(10000 + i))
		strs.WriteString(" String")
	}
	for i := 1; i <= sNum; i++ {
		strs.WriteString(", s_")
		strs.WriteString(strconv.Itoa(10000 + i))
		strs.WriteString(" UInt8")
	}
	strs.WriteString(", mergedata String) ENGINE = Memory()")
	return strs.String()
}

func perInsertTableSql(database string, tableName string, iNum, gNum, sNum int) string {
	strs := strings.Builder{}
	strs.Grow(1000)
	strs.WriteString("INSERT INTO ")
	strs.WriteString(database)
	strs.WriteString(".")
	strs.WriteString(tableName)
	strs.WriteString("( dd ,uid")
	for i := 1; i <= iNum; i++ {
		strs.WriteString(" ,i_")
		strs.WriteString(strconv.Itoa(10000 + i))
	}
	for i := 1; i <= gNum; i++ {
		strs.WriteString(" ,g_")
		strs.WriteString(strconv.Itoa(10000 + i))
	}
	for i := 1; i <= sNum; i++ {
		strs.WriteString(" ,s_")
		strs.WriteString(strconv.Itoa(10000 + i))
	}
	strs.WriteString(" ,mergedata) ")
	return strs.String()
}

func moreDataAndBatchInsert(preInsertTableSql string, iNum, gNum, sNum int, rowNum int) {

	var rowDatas [][]interface{}

	dd := time.Now()

	mergedata := strings.Builder{}
	mergedata.Grow(1000)
	for i := 0; i < rowNum; i++ {
		mergedata.Reset()
		if i%100000 == 0 && i != 0 {
			BatchInsert(preInsertTableSql, &rowDatas)
			log.Printf("生成模拟数据并入库中，当前条数[%d]", i)
			//清空切片
			rowDatas = rowDatas[0:0]
		}
		var row []interface{}
		row = append(row, dd)
		row = append(row, fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(i)))))
		var str string
		var itmp uint8
		for i := 1; i <= iNum; i++ {
			//使用一定长度的随机字符串填充
			str = getRandString(r.Intn(15) + 2)
			row = append(row, str)
		}
		for i := 1; i <= gNum; i++ {
			if i%2 == 0 {
				str = getRandomNation()
				row = append(row, str)
				mergedata.WriteString(str)
				mergedata.WriteString("|")
			} else {
				str = getRandomType()
				row = append(row, str)
				mergedata.WriteString(str)
				mergedata.WriteString("|")
			}
		}
		for i := 1; i <= sNum; i++ {
			itmp = getRandomBool()
			row = append(row, itmp)
			mergedata.WriteString(strconv.Itoa(int(itmp)))
			mergedata.WriteString("|")
		}

		row = append(row, strings.TrimRight(mergedata.String(), "|"))

		rowDatas = append(rowDatas, row)
	}

	if len(rowDatas) > 0 {
		BatchInsert(preInsertTableSql, &rowDatas)
	}
}

func getRandomBool() uint8 {
	return (uint8(r.Intn(2)))
}

func getRandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func getRandomType() string {
	randNum := r.Intn(3)
	switch randNum {
	case 0:
		return "mobile"
	case 1:
		return "adsl"
	case 2:
		return "idno"
	}
	return ""
}

var nativePlace = []int{11, 12, 13, 14, 15, 21, 22, 23, 31, 32, 33, 34, 35, 36,
	37, 41, 42, 43, 44, 45, 46, 51, 52, 53, 54, 50, 61, 62, 63, 64, 65}

func getRandomNation() string {
	return strconv.Itoa(nativePlace[r.Intn(len(nativePlace))])
}
