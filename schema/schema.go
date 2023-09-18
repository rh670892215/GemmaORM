package schema

import (
	"GemmaORM/dialect"
	"reflect"
)

// Field 表结构字段定义
type Field struct {
	Type string
	Name string
	Tag  string
}

// Schema 表结构
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FiledMap   map[string]*Field
	FieldNames []string
}

// Parse 根据model解析表结构
func Parse(model interface{}, dialect dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	schema := &Schema{
		Model:    model,
		Name:     modelType.Name(),
		FiledMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		// 若字段是匿名字段或不可导出，则不进行解析
		if modelType.Field(i).Anonymous || !modelType.Field(i).IsExported() {
			continue
		}
		filedName := modelType.Field(i).Name
		filedType := modelType.Field(i).Type
		fieldTag := modelType.Field(i).Tag.Get("gemmaorm")

		field := &Field{
			Name: filedName,
			Type: dialect.DataType(reflect.New(filedType)),
			Tag:  fieldTag,
		}
		schema.FiledMap[filedName] = field
		schema.Fields = append(schema.Fields, field)
		schema.FieldNames = append(schema.FieldNames, filedName)
	}

	return schema
}

// GetFieldByName 根据字段名获取field
func (s *Schema) GetFieldByName(name string) *Field {
	return s.FiledMap[name]
}

// RecordValues 将dest中各字段的值打平，如User{"bank",25} -> []interface{}{"bank",25}
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	value := reflect.Indirect(reflect.ValueOf(dest))
	var res []interface{}
	for _, field := range s.Fields {
		res = append(res, value.FieldByName(field.Name).Interface())
	}
	return res
}
