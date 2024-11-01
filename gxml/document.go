package gxml

import (
	"engine/gxml/etree"
	"engine/log"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

var logger = log.GetLogger()

// Initialize 初始化读取配置信
// path xml路径
// c 配置结构体
// m 存储的map 和 MapValue 中的参数一致
func Initialize(path string, c interface{}, m interface{}) {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Errorf("配置表读取失败:%s,%s", path, err.Error())
		return
	}

	doc := etree.NewDocument()
	err = doc.ReadFromBytes(content)
	if err != nil {
		logger.Errorf("配置表解析错误:%s,%s", path, err.Error())
		return
	}
	elements := doc.FindElements("/def/unit")
	if len(elements) == 0 {
		return
	}
	for _, el := range elements {
		singleElement(el, c, m, path)
	}
}

// singleElement 单个元素数据填充
func singleElement(el *etree.Element, c interface{}, m interface{}, path string) {
	defer func() {
		if r := recover(); r != nil {
			logger.TraceErr(r)
		}
	}()
	st := reflect.TypeOf(c).Elem()
	vl := reflect.New(st)
	elem := vl.Elem()
	var key int
	for i := 0; i < elem.NumField(); i++ {
		fieldInfo := elem.Type().Field(i)
		v := elem.FieldByName(fieldInfo.Name)
		tag := fieldInfo.Tag
		if len(tag) == 0 {
			continue
		}
		ret := tag.Get("xml")

		if len(ret) == 0 {
			continue
		}
		ell := el.SelectAttr(ret)
		if ell == nil || len(ell.Value) <= 0 {
			continue
		}
		fieldKind := v.Kind()
		switch fieldKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value, _ := strconv.ParseInt(ell.Value, 10, 64)
			v.SetInt(value)
			if ret == "id" {
				key = int(value)
			}
		case reflect.Float32, reflect.Float64:
			value, _ := strconv.ParseFloat(ell.Value, 64)
			v.SetFloat(value)
		case reflect.String:
			v.SetString(ell.Value)
		case reflect.Slice:
			va := v.Interface()
			if _, ok := va.([]int); ok {
				str := strings.Split(ell.Value, ",")
				var a []int
				for _, s1 := range str {
					value, _ := strconv.Atoi(s1)
					a = append(a, value)
				}
				v.Set(reflect.ValueOf(a))
				continue
			}
			if _, ok := va.([]int64); ok {
				str := strings.Split(ell.Value, ",")
				var a []int64
				for _, s1 := range str {
					value, _ := strconv.ParseInt(s1, 10, 64)
					a = append(a, value)
				}
				v.Set(reflect.ValueOf(a))
				continue
			}
			if _, ok := va.([]float32); ok {
				str := strings.Split(ell.Value, ",")
				var a []float32
				for _, s1 := range str {
					value, _ := strconv.ParseFloat(s1, 32)
					a = append(a, float32(value))
				}
				v.Set(reflect.ValueOf(a))
				continue
			}
			if _, ok := va.([]float64); ok {
				str := strings.Split(ell.Value, ",")
				var a []float64
				for _, s1 := range str {
					value, _ := strconv.ParseFloat(s1, 64)
					a = append(a, value)
				}
				v.Set(reflect.ValueOf(a))
				continue
			}
			if _, ok := va.([]string); ok {
				str := strings.Split(ell.Value, ",")
				v.Set(reflect.ValueOf(str))
				continue
			}
			ret1 := tag.Get("func")
			mo := vl.MethodByName(ret1)
			args := []reflect.Value{reflect.ValueOf(ell.Value)}
			r := mo.Call(args)
			if r != nil {
				v.Set(r[0])
			}

		default:
			ret1 := tag.Get("func")
			mo := vl.MethodByName(ret1)
			args := []reflect.Value{reflect.ValueOf(ell.Value)}
			r := mo.Call(args)
			if r != nil {
				v.Set(r[0])
			}
		}
	}
	mp := reflect.ValueOf(m)
	v := mp.MapIndex(reflect.ValueOf(key))
	if v.IsValid() {
		logger.Errorf("配置表出现相同的key:%s,%s", path, key)
	}
	mp.SetMapIndex(reflect.ValueOf(key), vl)
	mo := vl.MethodByName("AfterLoad")
	if mo.IsValid() {
		args := make([]reflect.Value, 0)
		mo.Call(args)
	}
}
