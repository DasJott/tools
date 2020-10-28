package crconfig

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

var conf map[string]string
var cli map[string]string

// Read parses file for valid config, if you have one
func Read(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	conf = map[string]string{"-d": "DEBUG"}
	cli = map[string]string{}
	args := map[string]string{}

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		if line[0] == '-' {
			parts := strings.SplitN(line, " ", 2)
			args[parts[0]] = parts[1]
		} else {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) > 1 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				conf[key] = val
			}
		}
	}

	parseSwitches(args)

	// find variables
	for k, v := range conf {
		ref := Get(v, v)
		if _, is := conf[v]; is {
			for is {
				ref = Get(ref, ref)
				_, is = conf[ref]
			}
			conf[k] = ref
		}
	}

	return nil
}

// Get gets the value according to the given key.
// if key is not found, def is returned
func Get(key, def string) string {
	if val, ok := cli[key]; ok {
		return val
	}
	if val := os.Getenv(key); val != "" {
		return val
	}
	if val, ok := conf[key]; ok {
		return val
	}
	return def
}

// GetBool gets the value as bool, according to the given key.
// if key is not found, def is used
func GetBool(key string, def bool) bool {
	if val := Get(key, ""); val != "" {
		return strings.ToLower(val) == "true"
	}
	return def
}

// GetInt gets the value as int64, according to the given key.
// if key is not found, def is used
func GetInt(key string, def int64) int64 {
	if val := Get(key, ""); val != "" {
		n, _ := strconv.ParseInt(val, 10, 64)
		return n
	}
	return def
}

// GetFloat gets the value as float64, according to the given key.
// if key is not found, def is used
func GetFloat(key string, def float64) float64 {
	if val := Get(key, ""); val != "" {
		n, _ := strconv.ParseFloat(val, 64)
		return n
	}
	return def
}

// GetDuration gets the value as int64, according to the given key.
// if key is not found, def is used
func GetDuration(key string, def time.Duration) time.Duration {
	if val := GetInt(key, 0); val != 0 {
		return time.Duration(val)
	}
	return def
}
