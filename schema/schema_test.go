package schema

import (
	"GemmaORM/dialect"
	"testing"
)

type User struct {
	Name string `gemmaorm:"PRIMARY KEY"`
	Age  int
}

func TestParse(t *testing.T) {
	user := &User{}
	d, _ := dialect.GetDialect("sqlite3")
	schema := Parse(user, d)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetFieldByName("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}
