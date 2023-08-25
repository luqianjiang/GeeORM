package schema

import (
	"day6-transaction/dialect"
	"go/ast"
	"reflect"
)

// Field :表示数据库的一列
type Field struct {
	// 字段名
	Name string

	// 字段类型
	Type string

	// 约束条件
	Tag string
}

// Schema 表示一个数据库表
type Schema struct {
	// 被映射的对象
	Model interface{}

	// 表名
	Name string

	// 字段信息
	Fields []*Field

	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// Parse 解析对象
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	// 遍历对象结果进行映射
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i) // 获取第i个属性
		// 非匿名（匿名：无字段名） 且可导出（首字母大写）
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// RecordValues 将对象u1 := &User{Name: "Tom", Age: 18} 转化成 ——>["Tom",18]
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
