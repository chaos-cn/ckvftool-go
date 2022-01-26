package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"log"
	"time"
)

const DrivierName = "clickhouse"

var conn *sql.DB
var ctx context.Context

// InitDB
/**
"clickhouse://127.0.0.1:9000?dial_timeout=1s&compress=true&max_execution_time=60"
*/
func InitDB(datasource string) {
	var err error

	conn, err = sql.Open(DrivierName, datasource)

	if err != nil {
		log.Fatal(err)
	}

	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	ctx = clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		fmt.Println("progress: ", p)
	}))
	if err := conn.PingContext(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		log.Fatal(err)
	}
	log.Println("初始化连接 ", datasource, "成功!")
}

func ExceSql(sql string) {
	_, err := conn.ExecContext(ctx, sql)

	if err != nil {
		log.Fatal(err)
	}
}

func QuerySql(querySql string) []map[string]string {
	rows, err := conn.QueryContext(ctx, querySql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var mapString []map[string]string
	var count = 0
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		rowMap := make(map[string]string)
		var value string
		for i, col := range values {
			if col != nil {
				value = string(col)
				rowMap[columns[i]] = value
			}
		}
		mapString = append(mapString, rowMap)
		count++
	}
	return mapString
}

// BatchInsert
/**
INSERT INTO example (Col1, Col2, Col3)
*/
func BatchInsert(preSql string, rowDatas *[][]interface{}) {
	if rowDatas == nil || *rowDatas == nil || len(*rowDatas) == 0 {
		log.Println("待入库的数据为空…………")
	}

	scope, err := conn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	batch, err := scope.PrepareContext(ctx, preSql)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(*rowDatas); i++ {
		if _, err := batch.Exec((*rowDatas)[i]...); err != nil {
			log.Fatal(err)
		}

		if i%10000 == 0 && i != 0 {
			if err := scope.Commit(); err != nil {
				log.Fatal(err)
			}
			batch.Close()

			scope, err = conn.Begin()
			if err != nil {
				log.Fatal(err)
			}
			batch, err = scope.PrepareContext(ctx, preSql)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := scope.Commit(); err != nil {
		log.Fatal(err)
	}
	batch.Close()

}

func CloseDB() {
	if conn != nil {
		conn.Close()
	}
	log.Println("关闭数据库连接")
}
