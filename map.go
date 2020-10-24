package gormmapper

// map
type Map struct {
	dict  map[string]interface{} // 数据字典
	value interface{}            // 当前值
	len   int64                  // 字典大小
	kType string                 // key 类型
	vType string                 // value 类型
}

/**
 * 实例化
 * @param
 * @return
 */
func NMap() *Map {
	hm := &Map{}
	hm.dict = make(map[string]interface{})
	return hm
}

/**
 * put
 * @param
 * @return
 */
func (m *Map) Put(key string, value interface{}) *Map {
	m.dict[key] = value
	m.len++
	return m
}

/**
 * get
 * @param
 * @return
 */
func (m *Map) Get(key string) interface{} {
	m.value = m.dict[key]
	return m.value
}

/**
 * key
 * @param
 * @return
 */
func (m *Map) Key(key string) *Map {
	m.value = m.dict[key]
	return m
}

/**
 *
 * @param
 * @return
 */
func (m *Map) String() string {
	if m.value == nil {
		return ""
	}
	return m.value.(string)
}

/**
 *
 * @param
 * @return
 */
func (m *Map) Int() int {
	if m.value == nil {
		return 0
	}
	return m.value.(int)
}

/**
 *
 * @param
 * @return
 */
func (m *Map) Int64() int64 {
	if m.value == nil {
		return 0
	}
	return m.value.(int64)
}

/**
 *
 * @param
 * @return
 */
func (m *Map) Float32() float32 {
	if m.value == nil {
		return 0
	}
	return m.value.(float32)
}

/**
 *
 * @param
 * @return
 */
func (m *Map) Float64() float64 {
	if m.value == nil {
		return 0
	}
	return m.value.(float64)
}

/**
 *
 * @param
 * @return
 */
func (m *Map) Keys() []string {
	keys := make([]string, 0)
	for k, _ := range m.dict {
		keys = append(keys, k)
	}
	return keys
}

/**
 *
 * @param
 * @return
 */
func (m *Map) Values() map[string]interface{} {
	return m.dict
}

/**
 *
 * @param
 * @return
 */
func (m *Map) Iterator() map[string]interface{} {
	return m.Values()
}
