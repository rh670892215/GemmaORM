package session

import (
	"GemmaORM/clause"
	"GemmaORM/dialect"
	"GemmaORM/schema"
	"database/sql"
	"strings"
)

// Session 一次与数据库交互的会话
type Session struct {
	db   *sql.DB
	tx   *sql.Tx
	sql  strings.Builder
	vars []interface{}

	clause      *clause.Clause
	tableSchema *schema.Schema
	dialect     dialect.Dialect
}

// CommonDB 数据库最小功能集合
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// NewSession new session
func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
		clause:  &clause.Clause{},
	}
}

// Clear 重新初始化session
func (s *Session) Clear() {
	s.sql.Reset()
	s.vars = []interface{}{}
}

// GetDB 获取本session连接的数据库
func (s *Session) GetDB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// Raw 添加sql语句和参数
func (s *Session) Raw(sql string, vars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.vars = append(s.vars, vars...)
	return s
}

// Query 获取符合查询条件的所有记录
func (s *Session) Query() (*sql.Rows, error) {
	defer s.Clear()
	return s.GetDB().Query(s.sql.String(), s.vars...)
}

// Exec 执行session中的sql语句
func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	return s.GetDB().Exec(s.sql.String(), s.vars...)
}

// QueryRow 获取符合查询条件的第一条记录
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	return s.GetDB().QueryRow(s.sql.String(), s.vars...)
}
