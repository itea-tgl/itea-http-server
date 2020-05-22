package itea_http_server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Rule struct {
	Key 	string
	Type 	string
	Rule 	string
	Msg 	string
	Default interface{}
	Enum	[]interface{}
}

type res struct {
	Key 	string
	Value 	interface{}
	Err 	string
}

type ret struct {
	data map[string]interface{}
}

func (r *ret) set(k string, v interface{}) {
	r.data[k] = v
}

func (r *ret) Int(k string) int {
	if v, ok := r.data[k].(int); ok {
		return v
	}
	return 0
}

func (r *ret) String(k string) string {
	if v, ok := r.data[k].(string); ok {
		return v
	}
	return ""
}

func (r *ret) Bool(k string) bool {
	if v, ok := r.data[k].(bool); ok {
		return v
	}
	return false
}

func Validate(r *http.Request, rules []Rule) (*ret, error) {
	data := &ret{make(map[string]interface{})}
	l := len(rules)
	ch := make(chan res, l)
	defer close(ch)

	for _, rule := range rules {
		go func(rule Rule) {
			if strings.EqualFold(rule.Key, "") {
				ch <- res{}
				return
			}

			key, value := getValue(rule.Key, r)
			if strings.EqualFold(value.(string), "") {
				if rule.Default != nil {
					value = rule.Default
				}
			} else {
				value = switchValue(value.(string), rule.Type)
			}

			if strings.EqualFold(rule.Rule, "") {
				ch <- res{key, value, ""}
				return
			}
			if checkRule(value, rule) {
				ch <- res{key, value, ""}
				return
			}

			if !strings.EqualFold(rule.Msg, "") {
				ch <- res{key, "", rule.Msg}
				return
			}

			ch <- res{key, "", fmt.Sprintf("Parameter [%s] validate error", key)}
			return
		}(rule)
	}

	hasError := false
	var errMsg []string

	for i := 0; i < l; i++ {
		res := <-ch
		if !strings.EqualFold(res.Err, "") {
			hasError = true
			errMsg = append(errMsg, res.Err)
			continue
		}
		if strings.EqualFold(res.Key, "") {
			continue
		}
		data.set(res.Key, res.Value)
	}

	if hasError {
		return nil, errors.New(strings.Join(errMsg, "; "))
	}

	return data, nil
}

func getValue(key string, r *http.Request) (string, interface{}) {
	keyArr := strings.Split(key, "|")
	if len(keyArr) == 1 {
		return keyArr[0], r.FormValue(key)
	}
	if strings.EqualFold(keyArr[1], "header") {
		return keyArr[0], r.Header.Get(keyArr[0])
	}
	return keyArr[0], ""
}

func switchValue(value string, typ string) interface{} {
	switch strings.ToLower(typ) {
	case "int":
		v, _ := strconv.Atoi(value)
		return v
	case "bool":
		v, _ := strconv.ParseBool(value)
		return v
	case "string":
		return value
	default:
		return value
	}
}

func checkRule(value interface{}, rule Rule) bool {
	switch rule.Rule {
	case "":
		return true
	case "required":
		if value == nil {
			return false
		}
		if v, ok := value.(string); ok && strings.EqualFold(v, "") {
			return false
		}
		return true
	case "include":
		for _, v := range rule.Enum {
			if value == v {
				return true
			}
		}
		return false
	default:
		return false
	}
}