package session

import (
	"day3-save-query/clause"
	"reflect"
)

// Insert 插入 入参数 ：参入对象
func (s *Session) Insert(values ...interface{}) (int64, error) {
	// recordValues 保存每个对象的值
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// reflect.New(destType) 将type类型转化为value类型
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			// dest.FieldByName(name).Addr() 返回字段的地址，然后通过rows.Scan(values...)赋值给values，同时也修改了dest字段地址的值
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		// destSlice 间接指向入参values，因此修改destSlice也会修改入参values，
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
