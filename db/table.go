package db

import (
	"database/sql"
	"strings"
)

// TableDescriber 接口表示表信息，表信息包含表名称，表备注和表的列信息
type TableDescriber interface {
	// DescribeTable 返回所有的表信息
	DescribeTable() ([]TableInfo, error)
}

// TableInfo 表示表信息
type TableInfo struct {
	TableName    string         `db:"TABLE_NAME"`
	TableComment sql.NullString `db:"TABLE_COMMENT"`
	ColumnInfos  []ColumnInfo
}

// ColumnInfo 表示列信息
type ColumnInfo struct {
	TableName    string         `db:"TABLE_NAME"`
	Name         string         `db:"NAME"`
	Kind         string         `db:"KIND"`
	Length       sql.NullString `db:"LENGTH"`
	Precision    sql.NullString `db:"PRECISION"`
	NullFlag     string         `db:"NULL_FLAG"`
	DefaultValue sql.NullString `db:"DEFAULT_VALUE"`
	Comment      sql.NullString `db:"COMMENTS"`
	PkFlag       string         `db:"PK_FLAG"`
}

// makeTableInfo 将 tableInfos 和 columnInfos 组装起来，
// 最后返回 tableInfos
func makeTableInfo(tableInfos []TableInfo, columnInfos []ColumnInfo) []TableInfo {
	// 转换成 map
	columnInfoMap := convertToMap(columnInfos)
	// 组装 tis
	var tis []TableInfo
	for _, ti := range tableInfos {
		ti.ColumnInfos = columnInfoMap[ti.TableName]
		tis = append(tis, ti)
	}

	return tis
}

// convertToMap 对 columnInfos 按 TableName 字段分组，分组后的结果是一个 map ，
// key 为 TableName ，value 为 []ColumnInfo
func convertToMap(columnInfos []ColumnInfo) map[string][]ColumnInfo {
	result := make(map[string][]ColumnInfo)
	for _, info := range columnInfos {
		// 如果 kind 里面有括号，则截取掉，因为括号里是长度，但是长度有单独的字段表示
		index := strings.Index(info.Kind, "(")
		if index != -1 {
			info.Kind = info.Kind[:index]
		}
		result[info.TableName] = append(result[info.TableName], info)
	}
	return result
}
