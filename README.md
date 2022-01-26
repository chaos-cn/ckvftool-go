# ckvftool-go
go语言编写的clickhouse垂直筛选场景模拟工具。

## 参数说明：

```
  -db string
    	database name (default "default")
  -gnum int
    	用于分组字段个数, >0 (default 8)
  -h string
    	clickhouse-server ip or host (default "127.0.0.1")
  -inum int
    	普通索引字段个数, >0 (default 5)
  -p int
    	clickhouse-server port (default 9000)
  -password string
    	clickhouse-server password
  -ptnum int
    	模拟并发用户数量, >0 (default 1)
  -rownum int
    	表中模拟的数据量, >0 (default 1000)
  -snum int
    	用于求和字段个数, >0 (default 10)
  -step string
    	执行步骤, all/init/test
  -tb string
    	table name (default "bigtable")
```

注意配置"-step"参数，用于指定是初始化表还是执行测试过程，默认为空，必须手动指定

## 问题说明

在Linux或者mac平台，如果显示无运行权限，请用chmod 命令进行添加权限
```bash
 # Linux amd64平台
 chmod 0755 ckvftool-linux-amd64
 # Mac darwin arm64平台
 chmod 0755 ckvftool-mac-arm64
 # Mac darwin amd64平台
 chmod 0755 ckvftool-mac-amd64
```

## 依赖项

[**ClickHouse**](https://clickhouse.com/docs/en/)

[**clickhouse-go**](https://github.com/ClickHouse/clickhouse-go)



## License

MIT