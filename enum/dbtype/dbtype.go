package dbtype

// DbType --
type DbType string

const (
	MYSQL     DbType = "mysql"
	POSTGRES  DbType = "postgres"
	SQLITE3   DbType = "sqlite3"
	SQLSERVER DbType = "mssql"
)
