package session

import (
	"GemmaORM/clause"
	"errors"
	"reflect"
)

// Insert 插入记录，exp: &User{"bank",25}
func (s *Session) Insert(values ...interface{}) (int64, error) {
	// insert into user (f1,f2) values ("bank",25),("emma",25)
	var recordValues []interface{}
	for _, value := range values {
		// 这里先调用model进行表schema解析
		table := s.Model(value).GetRefTable()
		s.CallMethod(BeforeInsert, value)
		recordValue := table.RecordValues(value)
		recordValues = append(recordValues, recordValue)
		// 添加insert语句
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
	}

	s.clause.Set(clause.VALUE, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUE)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return res.RowsAffected()
}

// Update 更改记录，支持输入 : 1."age",26; 2. map[string]interface{}{"age":26};
func (s *Session) Update(values ...interface{}) (int64, error) {
	// update user set col1 = val1,col2 = val2,... where [condition]
	valuesTable, ok := values[0].(map[string]interface{})
	if !ok {
		valuesTable = make(map[string]interface{})
		for i := 0; i < len(values); i += 2 {
			valuesTable[values[i].(string)] = values[i+1]
		}
	}

	s.CallMethod(BeforeUpdate, nil)
	s.clause.Set(clause.UPDATE, s.GetRefTable().Name, valuesTable)
	sql, args := s.clause.Build(clause.UPDATE, clause.WHERE)
	res, err := s.Raw(sql, args...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return res.RowsAffected()
}

// Delete 删除指定记录
func (s *Session) Delete() (int64, error) {
	// delete from user where [condition]
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.tableSchema.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return res.RowsAffected()
}

// Find 查找匹配的记录，并写入value，value需传入指针
func (s *Session) Find(value interface{}) error {
	s.CallMethod(BeforeQuery, nil)
	// select * from user
	valueSlice := reflect.Indirect(reflect.ValueOf(value))
	// 第一次Elem()调用，获取slice元素类型 - struct(kind 25)
	elementType := valueSlice.Type().Elem()
	// 第二次Elem()调用，获取slice的struct元素的具体类型
	table := s.Model(reflect.New(elementType).Elem().Interface()).GetRefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDER, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).Query()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(elementType).Elem()
		var values []interface{}
		for _, name := range s.GetRefTable().FieldNames {
			// 收集dest中每个字段的地址，填充到values中
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}

		// 将rows中的值写入values，即dest的每个字段的地址
		if err := rows.Scan(values...); err != nil {
			return err
		}

		s.CallMethod(AfterQuery, dest.Addr().Interface())
		valueSlice.Set(reflect.Append(valueSlice, dest))
	}

	return rows.Close()
}

// Where 条件语句，支持链式调用，desc : name = ?,args : "bank","emma"
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var params []interface{}
	// 将desc和args数组打平，统一成数组传入set
	s.clause.Set(clause.WHERE, append(append(params, desc), args...)...)
	return s
}

// Limit Limit子句，支持链式调用
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Order Order子句，支持链式调用
func (s *Session) Order(field string) *Session {
	s.clause.Set(clause.ORDER, field)
	return s
}

// First 获取符合条件的第一条记录
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

// Count 获取符合条件的记录数量
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.GetRefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
