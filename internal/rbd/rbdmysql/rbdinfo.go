package rbdmysql

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// executeSQLFile 函数读取并执行指定的 SQL 文件，并将结果导出到 CSV 文件
func executeSQLFile(db *sql.DB, filePath, outputFileName string) error {
	// 读取 SQL 文件内容
	query, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 执行 SQL 查询
	rows, err := db.Query(string(query))
	if err != nil {
		return fmt.Errorf("执行查询失败: %v", err)
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("获取列名失败: %v", err)
	}

	// 创建 CSV 文件
	file, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("创建 CSV 文件失败: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入列名到 CSV 文件
	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("写入列名到 CSV 文件失败: %v", err)
	}

	// 写入查询结果到 CSV 文件
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("扫描行失败: %v", err)
		}

		record := make([]string, len(columns))
		for i, val := range values {
			if b, ok := val.([]byte); ok {
				record[i] = string(b)
			} else {
				record[i] = fmt.Sprintf("%v", val)
			}
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("写入记录到 CSV 文件失败: %v", err)
		}
	}

	return nil
}

// RunSQL 函数根据传入的参数连接数据库并执行相应的 SQL 文件
func RunSQL(sip, port, user, password, env string) {
	// 构建数据源名称 (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/console", user, password, sip, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 根据环境参数选择要执行的 SQL 文件
	var filePath string
	switch env {
	case "store-ent3":
		filePath = "internal/rbd/rbdmysql/store_ent3.sql"
	case "store-ent4":
		filePath = "internal/rbd/rbdmysql/store_ent4.sql"
	default:
		log.Fatalf("无效的环境参数")
	}

	// 设置输出文件名
	outputFileName := fmt.Sprintf("output_%s.csv", env)

	// 执行 SQL 文件并导出结果到 CSV 文件
	if err := executeSQLFile(db, filePath, outputFileName); err != nil {
		log.Fatalf("执行 SQL 文件失败: %v", err)
	}

	fmt.Printf("SQL 文件执行成功，结果已导出到 %s\n", outputFileName)
}
