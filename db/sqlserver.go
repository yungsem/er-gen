package db

import (
	"github.com/jmoiron/sqlx"
)

// Sqlserver 表示 Sqlserver 数据库
type Sqlserver struct {
	db *sqlx.DB
}

// listAllTable 返回 db 中用户空间所有的表
func (r *Sqlserver) listAllTable() ([]TableInfo, error) {
	sql := `
		SELECT DISTINCT
			d.name AS TABLE_NAME,
			f.value AS TABLE_COMMENT 
		FROM
			syscolumns a
			LEFT JOIN systypes b ON a.xusertype= b.xusertype
			INNER JOIN sysobjects d ON a.id= d.id 
			AND d.xtype= 'U' 
			AND d.name != 'dtproperties'
			LEFT JOIN syscomments e ON a.cdefault= e.id
			LEFT JOIN sys.extended_properties g ON a.id= G.major_id 
			AND a.colid= g.minor_id
			LEFT JOIN sys.extended_properties f ON d.id= f.major_id 
			AND f.minor_id= 0
	`

	var tableInfos []TableInfo

	err := r.db.Select(&tableInfos, sql)
	if err != nil {
		return nil, err
	}

	return tableInfos, nil
}

// listAllColumn 返回 db 中用户空间所有表的所有列
func (r *Sqlserver) listAllColumn() ([]ColumnInfo, error) {
	sql := `
		SELECT 
			t.name AS TABLE_NAME,
			c.name AS NAME,
			ty.name AS KIND,
			c.max_length AS LENGTH,
			c.precision AS PRECISION,
				CASE WHEN c.is_nullable  = 1 THEN '是' ELSE '否' END AS NULL_FLAG,
			isnull(dc.definition, '' ) AS DEFAULT_VALUE,
			ep.value AS COMMENTS,
			CASE WHEN ic.column_id IS NULL THEN '否' ELSE '是' END AS PK_FLAG
		FROM 
			sys.tables t
		INNER JOIN 
			sys.columns c ON t.object_id = c.object_id
		INNER JOIN 
			sys.types ty ON c.system_type_id = ty.system_type_id
		LEFT JOIN 
			sys.default_constraints dc ON c.default_object_id = dc.object_id
		LEFT JOIN 
			sys.extended_properties ep ON ep.major_id = c.object_id AND ep.minor_id = c.column_id AND ep.class = 1 AND ep.name = 'MS_Description'
		LEFT JOIN 
			sys.indexes i ON t.object_id = i.object_id AND i.is_primary_key = 1
		LEFT JOIN 
			sys.index_columns ic ON i.object_id = ic.object_id AND c.column_id = ic.column_id AND i.index_id = ic.index_id
		WHERE 
    		ty.name <> 'sysname'
		ORDER BY 
			TABLE_NAME,
			c.column_id
	`

	var columnInfos []ColumnInfo

	err := r.db.Select(&columnInfos, sql)
	if err != nil {
		return nil, err
	}

	return columnInfos, nil
}

// DescribeTable 实现了 Table 接口的 TableInfos 方法
func (r *Sqlserver) DescribeTable() ([]TableInfo, error) {
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
