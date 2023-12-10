package dm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils"
)

func Update(config *callbacks.Config) func(db *gorm.DB) {
	supportReturning := utils.Contains(config.UpdateClauses, "RETURNING")
	return func(db *gorm.DB) {
		if db.Error != nil {
			return
		}

		if db.Statement.Schema != nil {
			for _, c := range db.Statement.Schema.UpdateClauses {
				db.Statement.AddClause(c)
			}
		}

		if db.Statement.SQL.Len() == 0 {
			db.Statement.SQL.Grow(180)
			db.Statement.AddClauseIfNotExists(clause.Update{})
			if set := callbacks.ConvertToAssignments(db.Statement); len(set) != 0 {
				db.Statement.AddClause(set)
			} else if _, ok := db.Statement.Clauses["SET"]; !ok {
				return
			}

			db.Statement.Build(db.Statement.BuildClauses...)
		}

		checkMissingWhereConditions(db)

		if !db.DryRun && db.Error == nil {
			if ok, mode := hasReturning(db, supportReturning); ok {
				// SQL 适配mysql
				sql := ConvertMysqlSql(db.Statement.SQL.String())
				if rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, sql, db.Statement.Vars...); db.AddError(err) == nil {
					dest := db.Statement.Dest
					db.Statement.Dest = db.Statement.ReflectValue.Addr().Interface()
					gorm.Scan(rows, db, mode)
					db.Statement.Dest = dest
					db.AddError(rows.Close())
				}
			} else {
				// SQL 适配mysql
				sql := ConvertMysqlSql(db.Statement.SQL.String())
				result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, sql, db.Statement.Vars...)

				if db.AddError(err) == nil {
					db.RowsAffected, _ = result.RowsAffected()
				}
			}
		}
	}
}

func checkMissingWhereConditions(db *gorm.DB) {
	if !db.AllowGlobalUpdate && db.Error == nil {
		where, withCondition := db.Statement.Clauses["WHERE"]
		if withCondition {
			if _, withSoftDelete := db.Statement.Clauses["soft_delete_enabled"]; withSoftDelete {
				whereClause, _ := where.Expression.(clause.Where)
				withCondition = len(whereClause.Exprs) > 1
			}
		}
		if !withCondition {
			db.AddError(gorm.ErrMissingWhereClause)
		}
		return
	}
}

func hasReturning(tx *gorm.DB, supportReturning bool) (bool, gorm.ScanMode) {
	if supportReturning {
		if c, ok := tx.Statement.Clauses["RETURNING"]; ok {
			returning, _ := c.Expression.(clause.Returning)
			if len(returning.Columns) == 0 || (len(returning.Columns) == 1 && returning.Columns[0].Name == "*") {
				return true, 0
			}
			return true, gorm.ScanUpdate
		}
	}
	return false, 0
}
