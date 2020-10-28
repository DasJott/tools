package crconfig

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Bind binds all found values into given struct.
// Supported types are string, bool and int, float in all bitdepths
func Bind(obj interface{}) error {
	return bind(obj, false)

}

// BindExclusive binds only env tagged values into given struct.
// Supported types are string, bool and int, float in all bitdepths
func BindExclusive(obj interface{}) error {
	return bind(obj, true)
}

func bind(obj interface{}, exclusive bool) error {
	v := reflect.ValueOf(obj).Elem()
	if !v.CanSet() {
		return fmt.Errorf("can not set data to given obj")
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		fieldt := t.Field(i)
		name := fieldt.Name
		def := "" // default value?
		if tag := fieldt.Tag.Get("env"); tag != "" {
			name = tag
			if parts := strings.SplitN(tag, ",", 2); len(parts) > 1 {
				name = parts[0]
				def = parts[1]
			}
		} else if exclusive {
			continue
		}

		switch {
		case fieldt.Type.Kind() == reflect.String:
			v.Field(i).SetString(Get(name, def))
		case fieldt.Type.Kind() == reflect.Bool:
			v.Field(i).SetBool(GetBool(name, (strings.ToLower(def) == "true")))
		case strings.HasPrefix(fieldt.Type.Name(), "int"):
			val, _ := strconv.ParseInt(def, 10, 64)
			v.Field(i).SetInt(GetInt(name, val))
		case strings.HasPrefix(fieldt.Type.Name(), "float"):
			val, _ := strconv.ParseFloat(def, 64)
			v.Field(i).SetFloat(GetFloat(name, val))
		}
	}

	return nil
}

// GetWithPrefix gets all key/values found in environment and config file.
// Environment wins over config here as well.
func GetWithPrefix(prefix string) map[string]string {
	res := map[string]string{}

	for key, val := range cli {
		if strings.HasPrefix(key, prefix) {
			res[key] = val
		}
	}

	for key, val := range conf {
		if strings.HasPrefix(key, prefix) {
			res[key] = val
		}
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		if strings.HasPrefix(key, prefix) {
			res[key] = pair[1]
		}
	}

	return res
}
