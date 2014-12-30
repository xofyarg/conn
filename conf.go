package main

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func tokenConv(token string) string {
	parts := strings.Split(token, "_")
	for i, s := range parts {
		parts[i] = strings.Title(s)
	}
	return strings.Join(parts, "")
}

func LoadConfig(file string, conf interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewScanner(f)
	elem := reflect.ValueOf(conf).Elem()

	for buf.Scan() {
		line := strings.TrimLeft(buf.Text(), " ")
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.Trim(parts[0], " '\"")
		value := strings.Trim(parts[1], " '\"")

		field := elem.FieldByName(tokenConv(key))
		if field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Bool:
				field.SetBool(value == "true")
			case reflect.Int, reflect.Int8, reflect.Int16,
				reflect.Int32, reflect.Int64:
				n, _ := strconv.ParseInt(value, 10, 64)
				field.SetInt(n)
			case reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64:
				n, _ := strconv.ParseUint(value, 10, 64)
				field.SetUint(n)
			default:
			}
		} else {
			Warn("unknown option: %s", key)
		}
	}

	return nil
}
