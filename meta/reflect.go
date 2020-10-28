package meta

import (
	"reflect"
	"strings"
)

// GetJSONName returns the name and omitempty of given structfield
func GetJSONName(field reflect.StructField) (string, bool) {
	return GetFieldName("json", field)
}

// GetSQLName returns the name and omitempty of given structfield
func GetSQLName(field reflect.StructField) (string, bool) {
	return GetFieldName("sql", field)
}

// GetBSONName returns the name and omitempty of given structfield
func GetBSONName(field reflect.StructField) (string, bool) {
	return GetFieldName("bson", field)
}

// GetFieldName returns the name and omitempty of the tag of the given structfield
func GetFieldName(tag string, field reflect.StructField) (string, bool) {
	name, omitempty := field.Tag.Get(tag), false
	if parts := strings.Split(name, ","); name != "" && len(parts) > 0 {
		name = parts[0]
		omitempty = len(parts) > 1 && parts[1] == "omitempty"
	} else {
		name = field.Name
	}
	return name, omitempty
}

// GetFieldNames returns a slice of names of the given structfield
func GetFieldNames(obj interface{}, tag string) []string {
	t := reflect.TypeOf(obj)
	c := t.NumField()

	names := make([]string, 0, c)

	for i := 0; i < c; i++ {
		if name, omit := GetFieldName(tag, t.Field(i)); !omit {
			names = append(names, name)
		}
	}

	return names
}

// RangeFields iterates the fields of given struct while calling the given function with index and the value as interface{}.
func RangeFields(obj interface{}, cb func(int, interface{})) {
	v := reflect.ValueOf(obj)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		cb(i, v.Field(i).Interface())
	}
}

// GetFieldValueMap gets a map from given tag name value to interace{}
func GetFieldValueMap(obj interface{}, tag string) map[string]interface{} {
	themap := map[string]interface{}{}

	v := reflect.ValueOf(obj)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if name, omit := GetFieldName(tag, t.Field(i)); !omit {
			themap[name] = v.Field(i).Interface()
		}
	}

	return themap
}
