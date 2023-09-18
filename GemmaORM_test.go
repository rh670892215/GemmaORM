package GemmaORM

import (
	"GemmaORM/session"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `gemmaorm:"PRIMARY KEY"`
	Age  int
}

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "bank.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}
func Test_transaction_commit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	err := s.Model(&User{}).DropTable()
	_, err = engine.Transaction(func(s *session.Session) (interface{}, error) {
		err := s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, err
	})

	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

func Test_transaction_rollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		s.Model(&User{}).CreateTable()
		s.Insert(&User{"Tom", 18})
		return nil, errors.New("Error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}
