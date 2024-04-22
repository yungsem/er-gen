package db

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
	"net/url"
)

const (
	TypeMysql     = "mysql"
	TypeSqlserver = "sqlserver"
	TypeOracle    = "oracle"

	// DNSMysql 格式：user:password@tcp(ip:port)/dbname
	DNSMysql = "%s:%s@tcp(%s:%d)/%s"
	// DNSSqlserver 格式：sqlserver://sa:mypass@localhost:1234?database=master&connection+timeout=30
	DNSSqlserver = "sqlserver://%s:%s@%s:%d?database=%s&connection+timeout=30"
	// DNSOracle 格式：oracle://user:password@[::1]:1521/service
	DNSOracle = "oracle://%s:%s@%s:%d/%s"
)

// NewTableDescriber 创建一个 TableDescriber
// 如果 dbType = mysql ，返回的是 Mysql 实例
// 如果 dbType = sqlserver ，返回的是 Sqlserver 实例
// 如果 dbType = oracle ，返回的是 Oracle 实例
// 默认返回 Mysql 实例
func NewTableDescriber(dbType string, host string, port int,
	username string, password string, schema string) (TableDescriber, error) {
	db, err := newDB(dbType, host, port, username, password, schema)
	if err != nil {
		return nil, err
	}

	switch dbType {
	case TypeMysql:
		return &Mysql{schema , db, }, nil
	case TypeSqlserver:
		return &Sqlserver{db}, nil
	case TypeOracle:
		return &Oracle{db}, nil
	default:

		return &Mysql{schema, db}, nil
	}
}

// newDB 创建一个 *sqlx.DB 实例
func newDB(dbType string, host string, port int, username string, password string, schema string) (*sqlx.DB, error) {
	var dnsF string
	switch dbType {
	case TypeMysql:
		dnsF = DNSMysql
	case TypeSqlserver:
		dnsF = DNSSqlserver
	case TypeOracle:
		dnsF = DNSOracle
	default:
		dnsF = DNSMysql
	}

	dns := fmt.Sprintf(dnsF,
		url.QueryEscape(username),
		url.QueryEscape(password),
		host,
		port,
		schema)

	db, err := sqlx.Open(dbType, dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
