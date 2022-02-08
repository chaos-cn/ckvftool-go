package main

import (
	"flag"
	"fmt"
	ckvftool "github.com/chaos-cn/ckvftool-go/cmd"
	"log"
	"os"
	"strconv"
	"strings"
)

// VERSION
/**
1.0.0:实现基础的工鞥
1.1.0:线上调试完毕
1.2.0:增加sql别名功能
2.0.0:切换成github.com/ClickHouse/clickhouse-go/v2
*/
const VERSION = "1.2.0"

var (
	//clickhouse服务器相关信息
	ckip       string
	ckport     int
	ckpassword string
	database   string
	//表结构相关信息
	tablename string
	iNum      int
	gNum      int
	sNum      int
	rowNum    int
	//测试用户数
	ptNum int
	//程序执行步骤
	step string
)

func init() {
	flag.StringVar(&ckip, "h", "127.0.0.1", "clickhouse-server ip or host")
	flag.IntVar(&ckport, "p", 9000, "clickhouse-server port")
	flag.StringVar(&ckpassword, "password", "", "clickhouse-server password")
	flag.StringVar(&database, "db", "default", "database name")
	flag.StringVar(&tablename, "tb", "bigtable", "table name")

	flag.IntVar(&iNum, "inum", 5, "普通索引字段个数, >0")
	flag.IntVar(&gNum, "gnum", 8, "用于分组字段个数, >0")
	flag.IntVar(&sNum, "snum", 10, "用于求和字段个数, >0")
	flag.IntVar(&rowNum, "rownum", 1000, "表中模拟的数据量, >0")
	//遇到[Too many simultaneous queries] 需要修改[/etc/clickhouse-server/config.xml<max_concurrent_queries>]
	flag.IntVar(&ptNum, "ptnum", 1, "模拟并发用户数量, >0")

	flag.StringVar(&step, "step", "", "执行步骤, all/init/test")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	_, err := fmt.Fprintf(os.Stderr, "ckutl version: ckvftool-go/%s\n\nOptions:\n", VERSION)
	if err != nil {
		log.Panic(err)
	}
	flag.PrintDefaults()
}

func paramCheck() {
	if iNum < 1 {
		log.Fatal("iNum 至少为 1")
	}
	if gNum < 1 {
		log.Fatal("gNum 至少为 1")
	}
	if sNum < 1 {
		log.Fatal("sNum 至少为 1")
	}
	if rowNum < 1 {
		log.Fatal("rowNum 至少为 1")
	}
	if ptNum < 1 {
		log.Fatal("ptNum 至少为 1")
	}

	if step == "" {
		usage()
		log.Fatal("[step] not specified!")
	}
}

func datasourceJoin() string {
	strs := strings.Builder{}
	strs.WriteString("clickhouse://")
	if ckpassword != "" {
		strs.WriteString(":")
		strs.WriteString(ckpassword)
		strs.WriteString("@")
	}
	strs.WriteString(ckip)
	strs.WriteString(":")
	strs.WriteString(strconv.Itoa(ckport))
	strs.WriteString("?")

	strs.WriteString("dial_timeout=1s&compress=true&max_execution_time=120")
	return strs.String()
}

func main() {

	flag.Parse()

	paramCheck()

	datasource := datasourceJoin()

	log.Printf("datasource:[%s] iNum:[%d] gNum:[%d] sNum:[%d] rowNum:[%d] ptNum:[%d] step:[%s]", datasource, iNum, gNum, sNum, rowNum, ptNum, step)

	switch step {
	case "init":
		//初始化表，并塞入测试数据
		ckvftool.InitTableAndData(datasource, database, tablename, iNum, gNum, sNum, rowNum)
	case "test":
		ckvftool.PerformanceTest(datasource, database, tablename, iNum, gNum, sNum, ptNum)
	case "all":
		ckvftool.InitTableAndData(datasource, database, tablename, iNum, gNum, sNum, rowNum)
		ckvftool.PerformanceTest(datasource, database, tablename, iNum, gNum, sNum, ptNum)
	default:
		log.Fatal("[step] unknown, 执行步骤可选 all/init/performance")
	}
}
