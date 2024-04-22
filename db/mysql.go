package db

import (
	"github.com/jmoiron/sqlx"
)

// Mysql 表示 mysql 数据库
type Mysql struct {
	schema string
	db     *sqlx.DB
}

// listAllTable 返回 db 中用户空间所有的表
func (r *Mysql) listAllTable() ([]TableInfo, error) {
	sql := `
		SELECT
			TABLE_NAME,
			TABLE_COMMENT 
		FROM
			information_schema.TABLES 
		WHERE
			TABLE_SCHEMA = ?
			AND table_type = 'BASE TABLE'
	`

	var tableInfos []TableInfo

	err := r.db.Select(&tableInfos, sql, r.schema)
	if err != nil {
		return nil, err
	}

	return tableInfos, nil
}

// listAllColumn 返回 db 中用户空间所有表的所有列
func (r *Mysql) listAllColumn() ([]ColumnInfo, error) {
	sql := `
		SELECT
			TABLE_NAME,
			COLUMN_NAME AS 'NAME',
			COLUMN_TYPE AS 'KIND',
			IF(CHARACTER_MAXIMUM_LENGTH IS NULL, NUMERIC_PRECISION, CHARACTER_MAXIMUM_LENGTH) AS 'LENGTH',
			NUMERIC_SCALE AS 'PRECISION',
			IF(IS_NULLABLE = 'YES', '是', '否') AS 'NULL_FLAG',
			COLUMN_DEFAULT AS 'DEFAULT_VALUE',
			IF(COLUMN_NAME = 'ID', '主键ID', COLUMN_COMMENT) AS 'COMMENTS',
			IF(COLUMN_KEY = 'PRI', '是', '否') AS 'PK_FLAG'
		FROM
			information_schema.COLUMNS
		WHERE
			TABLE_SCHEMA = ?
		ORDER BY
			ORDINAL_POSITION
	`

	var columnInfos []ColumnInfo

	err := r.db.Select(&columnInfos, sql, r.schema)
	if err != nil {
		return nil, err
	}

	return columnInfos, nil
}

// DescribeTable 实现了 Table 接口的 TableInfos 方法
func (r *Mysql) DescribeTable() ([]TableInfo, error) {
	tableInfos, err := r.listAllTable()
	if err != nil {
		return nil, err
	}

	columnInfos, err := r.listAllColumn()
	if err != nil {
		return nil, err
	}

	return makeTableInfo(tableInfos, columnInfos), nil
}
