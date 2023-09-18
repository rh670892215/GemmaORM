package session

import (
	"GemmaORM/schema"
	"fmt"
	"reflect"
	"strings"
)

// Model 解析表结构
func (s *Session) Model(model interface{}) *Session {
	if s.tableSchema == nil || reflect.TypeOf(model) != reflect.TypeOf(s.tableSchema.Model) {
		s.tableSchema = schema.Parse(model, s.dialect)
	}

	return s
}

// GetRefTable 获取当前session关联表结构
func (s *Session) GetRefTable() *schema.Schema {
	return s.tableSchema
}

// CreateTable 创建表
func (s *Session) CreateTable() error {
	tableName := s.GetRefTable().Name
	var columns []string
	for _, field := range s.GetRefTable().Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	columnsStr := strings.Join(columns, ",")
	sql := fmt.Sprintf("create table %s (%s)", tableName, columnsStr)
	_, err := s.Raw(sql).Exec()
	return err
}

// DropTable 删除表
func (s *Session) DropTable() error {
	tableName := s.GetRefTable().Name
	sql := fmt.Sprintf("drop table if exists %s", tableName)
	_, err := s.Raw(sql).Exec()
	return err
}
