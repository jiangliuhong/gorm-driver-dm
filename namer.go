package dm

import (
	"gorm.io/gorm/schema"
	"strings"
)

type Namer struct {
	schema.NamingStrategy
}

func ConvertNameToFormat(x string) string {
	//name := strings.ToUpper(x)
	name := x
	// 对关键字进行处理
	//if IsReservedWord(name) {
	//	name = fmt.Sprintf(`"%s"`, name)
	//}
	return name
}

func (n Namer) TableName(table string) (name string) {
	//return ConvertNameToFormat(n.NamingStrategy.TableName(table))
	var tn string
	if n.NamingStrategy.SingularTable {
		tn = n.NamingStrategy.TablePrefix + n.NamingStrategy.ColumnName("", table)
	}
	tn = n.NamingStrategy.TablePrefix + n.NamingStrategy.ColumnName("", table)
	return ConvertNameToFormat(tn)
}

func (n Namer) ColumnName(table, column string) (name string) {
	return ConvertNameToFormat(n.NamingStrategy.ColumnName(table, column))
}

func (n Namer) JoinTableName(table string) (name string) {
	return ConvertNameToFormat(n.NamingStrategy.JoinTableName(table))
}

func (n Namer) RelationshipFKName(relationship schema.Relationship) (name string) {
	return ConvertNameToFormat(n.NamingStrategy.RelationshipFKName(relationship))
}

func (n Namer) CheckerName(table, column string) (name string) {
	return ConvertNameToFormat(n.NamingStrategy.CheckerName(table, column))
}

func (n Namer) IndexName(table, column string) (name string) {
	tlc := strings.ToLower(column)

	cl := n.NamingStrategy.IndexName(table, column)
	if strings.Contains(tlc, "idx_"+strings.ToLower(table)) && strings.Contains(tlc, strings.ToLower(column)) {
		cl = column
	}

	return ConvertNameToFormat(cl)
}
