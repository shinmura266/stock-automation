package database

import "reflect"

// GetStringField リフレクションでstring型フィールドを取得
func GetStringField(v reflect.Value, fieldName string) string {
	field := v.FieldByName(fieldName)
	if !field.IsValid() || field.Kind() != reflect.String {
		return ""
	}
	return field.String()
}

// GetFloatField リフレクションでfloat64型フィールドを取得
func GetFloatField(v reflect.Value, fieldName string) float64 {
	field := v.FieldByName(fieldName)
	if !field.IsValid() || field.Kind() != reflect.Float64 {
		return 0
	}
	return field.Float()
}
