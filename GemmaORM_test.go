package GemmaORM

import (
	"GemmaORM/log"
	"GemmaORM/session"
	"errors"
	"fmt"
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

func Test_engine_use(t *testing.T) {
	engine, err := NewEngine("sqlite3", "bank.db")
	if err != nil {
		log.Error(err)
		return
	}
	defer engine.Close()

	s := engine.NewSession()
	s.Model(&User{}).CreateTable()
	var res []User
	s.Insert(&User{"bank", 25})
	s.Insert(&User{"emma", 25})
	s.Where("name = ?", "emma").Find(&res)
	fmt.Println(res)
	//s.Where("name = ?", "bank").Find(&res)
	//fmt.Println(res)
	//s.Where("name = ?", "bank").Update("age", 26)
	//s.Where("name = ?", "bank").Find(&res)
	//fmt.Println(res)
	//s.Where("name = ?", "bank").Update(map[string]interface{}{"age": 27})
	//s.Where("name = ?", "bank").Find(&res)
	//fmt.Println(res)
	//r, err := s.Where("name = ?", "emma").Delete()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(r)
	//count, err := s.Where("name = ?", "emma").Count()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(count)
}

func main() {
	engine, err := NewEngine("sqlite3", "bank.db")
	if err != nil {
		log.Error(err)
		return
	}
	defer engine.Close()

	s := engine.NewSession()
	s.Model(&User{}).CreateTable()
	s.Insert(&User{"bank", 25})
	var res []User
	s.Where("name = ?", "emma").Find(&res)
	s.Where("name = ?", "bank").Update("age", 26)
	s.Where("name = ?", "bank").Update(map[string]interface{}{"age": 27})
}
