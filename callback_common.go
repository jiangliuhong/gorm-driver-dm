package dm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"strings"
)

// ConvertMysqlSql 兼容mysql的sql
// 将 ` 转为 "
func ConvertMysqlSql(sql string) string {
	// 替换 `` 符号
	sql = strings.ReplaceAll(sql, "`", "\"")
	return sql
}

func BuildQuerySQL(db *gorm.DB) {
	callbacks.BuildQuerySQL(db)
	sql := db.Statement.SQL.String()
	sql = ConvertMysqlSql(sql)
	db.Statement.SQL.Reset()
	db.Statement.SQL.WriteString(sql)
}
