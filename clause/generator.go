package clause

import (
	"fmt"
	"strings"
)

type generator func(...interface{}) (string, []interface{})

var generatorMap map[AtomType]generator

func init() {
	generatorMap = make(map[AtomType]generator)
	generatorMap[INSERT] = _insert
	generatorMap[VALUE] = _value
	generatorMap[SELECT] = _select
	generatorMap[LIMIT] = _limit
	generatorMap[WHERE] = _where
	generatorMap[ORDER] = _order
	generatorMap[UPDATE] = _update
	generatorMap[DELETE] = _delete
	generatorMap[COUNT] = _count
}

// count条件子句，input : tableName
func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}

// where条件子句，input : name = ? , "bank","emma"...
func _where(values ...interface{}) (string, []interface{}) {
	// where name = ? , "bank"
	desc := values[0]
	args := values[1:]

	return fmt.Sprintf("where %s", desc), args
}

// select子句，input : tableName,[]sting{col1,col2...}
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	params := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("select %s from %s", params, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	num := values[0]
	return "limit ?", []interface{}{num}
}

func _order(values ...interface{}) (string, []interface{}) {
	orderFiled := values[0]
	return fmt.Sprintf("order by %s", orderFiled), []interface{}{}
}

// 更改记录，input : tableName map[string]interface{}{colName:colValue}
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	tableValue := values[1].(map[string]interface{})
	var valueSlice []string
	var vars []interface{}

	for k, v := range tableValue {
		valueSlice = append(valueSlice, fmt.Sprintf("%s = ?", k))
		vars = append(vars, v)
	}
	valueSliceStr := strings.Join(valueSlice, ",")
	sql := fmt.Sprintf("update %s set %s", tableName, valueSliceStr)
	return sql, vars
}

// delete子句，input : tableName
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("delete from %s", values[0]), []interface{}{}
}

// 插入记录(单条),input : tableName,[]string{"name","age"}
func _insert(values ...interface{}) (string, []interface{}) {
	// insert into user (name,age)
	tableName := values[0]
	args := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("insert into %s (%s)", tableName, args), []interface{}{}
}

// 输入values exp : []interface{"bank",25},[]interface{"emma",25}
func _value(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var vars []interface{}
	var sqlArr []string

	sql.WriteString("values ")
	for _, value := range values {
		v := value.([]interface{})
		sqlArr = append(sqlArr, fmt.Sprintf("(%s)", genBindVars(len(v))))
		vars = append(vars, v...)
	}
	sql.WriteString(strings.Join(sqlArr, ","))
	// values (?,?),(?,?)  "bank",25,"emma",25
	return sql.String(), vars
}

// 根据num生成指定数量的?,如 num = 3 ,output = ?,?,?
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}
