package query

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func CreateQueryBuilder(entity TableInterface) string {
	entityValue := reflect.ValueOf(entity).Elem()
	entityType := entityValue.Type()
	tableName := entity.TableName()

	var columns, values []string
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)

		columns = append(columns, field.Tag.Get("db"))
		values = append(values, fmt.Sprintf(":%v", field.Tag.Get("db")))
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))
}

func ReadQueryBuilder(entity TableInterface, key string) string {

	query := "SELECT "

	entityValue := reflect.ValueOf(entity).Elem()
	entityType := entityValue.Type()
	tableName := entity.TableName()

	var columns []string
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)

		columns = append(columns, fmt.Sprintf("IFNULL(%s, '') AS %s", field.Tag.Get("db"), field.Tag.Get("db")))
	}

	query += fmt.Sprintf("%s FROM %s ", strings.Join(columns, ", "), tableName)
	if key != "" {
		query += fmt.Sprintf("WHERE %s = :%s", key, key)
	}

	return query
}

func UpdateQueryBuilder(entity TableInterface) string {
	entityValue := reflect.ValueOf(entity).Elem()
	entityType := entityValue.Type()
	tableName := entity.TableName()

	var columns []string
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)

		columnName := field.Tag.Get("db")
		columns = append(columns, fmt.Sprintf("%s = :%s", columnName, columnName))
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE id = :id", tableName, strings.Join(columns, ", "))
}

func DeleteQueryBuilder(entity TableInterface) string {
	tableName := entity.TableName()
	return fmt.Sprintf("DELETE FROM %s WHERE id = :id", tableName)
}

func GetValueByKey(entity interface{}, key string) (interface{}, error) {

	key = ConvertToCamelCase(key)
	v := reflect.ValueOf(entity)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("entity must be a pointer to a struct")
	}

	v = v.Elem()
	field := v.FieldByName(key)

	if !field.IsValid() {
		return nil, errors.New("key not found")
	}

	if !field.CanInterface() {
		return nil, errors.New("field is not accessible (make sure it's public)")
	}

	return field.Interface(), nil
}

func ConvertToCamelCase(input string) string {
	if strings.Contains(input, "_") {
		words := strings.Split(input, "_")

		for i := range words {
			words[i] = strings.Title(words[i])
		}

		return strings.Join(words, "")
	}

	runes := []rune(input)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func ConvertToSnakeCase(input string) string {
	var result strings.Builder

	for i, r := range input {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}

		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

type TableInterface interface {
	TableName() string
}
