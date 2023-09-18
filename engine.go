package GemmaORM

import (
	"GemmaORM/dialect"
	"GemmaORM/log"
	"GemmaORM/session"
	"database/sql"
	"errors"
	"fmt"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine input : sqlite3,bank.db
func NewEngine(driver, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		return nil, errors.New(fmt.Sprintf("dialect %s not found", driver))
	}
	return &Engine{db: db, dialect: dial}, nil
}

// NewSession 新建session
func (e *Engine) NewSession() *session.Session {
	return session.NewSession(e.db, e.dialect)
}

// Close 关闭数据库连接
func (e *Engine) Close() {
	e.db.Close()
}

type TxFunc func(*session.Session) (interface{}, error)

func (e *Engine) Transaction(f TxFunc) (res interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			log.Info("come in recover rollback")
			log.Info(p)
			_ = s.RollBack()
		} else if err != nil {
			log.Info("come in err rollback")
			_ = s.RollBack()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}
