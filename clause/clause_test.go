package clause

import (
	"reflect"
	"testing"
)

func TestClause_Set(t *testing.T) {
	var clause Clause
	clause.Set(INSERT, "User", []string{"Name", "Age"})
	sql := clause.sql[INSERT]
	vars := clause.sqlVars[INSERT]
	t.Log(sql, vars)
	if sql != "insert into User (Name,Age)" || len(vars) != 0 {
		t.Fatal("failed to get clause")
	}
}

func TestSelect(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(ORDER, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDER, LIMIT)
	t.Log(sql, vars)
	if sql != "select * from User where Name = ? order by Age ASC limit ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}
