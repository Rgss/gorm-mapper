package gormmapper

import (
	"crypto/md5"
	"fmt"
	"io"
	"reflect"
	"strings"
)

/**
 * md5
 * @param
 * @return
 */
func Md5(in string) string {
	w := md5.New()
	io.WriteString(w, in)
	m := fmt.Sprintf("%x", w.Sum(nil))
	return m
}

/**
 * 转换成map数据结构
 * @param
 * @return
 */
func toMap(entity interface{}, dbFormat bool) map[string]interface{} {
	nm := make(map[string]interface{})
	vo := reflect.ValueOf(entity)
	switch vo.Kind() {
	case reflect.Map:
		nm = entity.(map[string]interface{})
		break
	case reflect.Struct:
	case reflect.Ptr:
		elem := vo.Elem()
		for i := 0; i < elem.NumField(); i++ {
			// 包含gorm tag
			tag := elem.Type().Field(i).Tag.Get("gorm")
			tag = strings.TrimSpace(tag)
			if tag == "-" {
				continue
			}

			vN := parseGormTagColumn(tag)
			if len(vN) <= 0 {
				vN = elem.Type().Field(i).Name
			}

			nm[vN] = elem.Field(i).Interface()
		}
		break
	default:
	}

	if !dbFormat {
		return nm
	}

	nnm := make(map[string]interface{})
	for k, v := range nm {
		nk := toDBColumnName(k)
		nnm[nk] = v
	}

	return nnm
}

/**
 * 解析tag 字段名设置
 * @param
 * @return
 */
func parseGormTagColumn(tag string) string {
	c := ""
	cs := "column:"
	p1 := strings.Index(tag, cs)
	if p1 < 0 {
		return c
	}

	s := tag[p1+7:]
	p2 := strings.Index(s, ";")
	if p2 < 0 {
		c = s
	} else {
		c = s[:p2]
	}

	return strings.TrimSpace(c)
}

/**
 * struct属性名转为数据库名
 * @param
 * @return
 */
func toDBColumnName(name string) string {
	var c string
	n := []rune(name)
	for i := 0; i < len(n); i++ {
		if n[i] >= 65 && n[i] <= 90 {
			n[i] += 32
			if i == 0 {
				c += string(n[i])
			} else {
				c += "_" + string(n[i])
			}
		} else {
			c += string(n[i])
		}
	}

	return c
}
