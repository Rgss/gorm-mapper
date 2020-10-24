package gormmapper

import (
	"strings"
)

// 最大页数
const maxPage = 10000

// 单页最大条数
const maxSize = 10000

// 常规单页条数
const defaultSize = 10

// 搜索条件构建器
type SearchBuilder struct {
	page        int
	size        int
	limit       int
	group       interface{}
	sort        map[string]string
	where       map[string]interface{}
	parsedWhere map[string]interface{}
	parsedValue map[string]interface{}
	fields      interface{}
	Entity      interface{}
	debug       bool
	maxPage     int
	maxSize     int
}

/**
 * 创建搜索构建器
 * @param
 * @return
 */
func Builder(entity interface{}) *SearchBuilder {
	return &SearchBuilder{
		Entity:      entity,
		sort:        make(map[string]string),
		parsedWhere: make(map[string]interface{}),
		parsedValue: make(map[string]interface{}),
	}
}

/**
 * 设置条件
 * @param where
 * @return
 */
func (sb *SearchBuilder) Where(where map[string]interface{}) *SearchBuilder {
	sb.where = where
	return sb
}

/**
 * 条件
 * @author  imp.zhang
 * @date    2020/10/22
 * @param
 * @return
 */
func (sb *SearchBuilder) Sort(column string, sortType string) *SearchBuilder {
	sb.sort[column] = sortType
	return sb
}

/**
 * 设置字段
 * @date    2020/10/22
 * @param
 * @return
 */
func (sb *SearchBuilder) Field(fields interface{}) *SearchBuilder {
	sb.fields = fields
	return sb
}

/**
 * 设置页数
 * @date    2020/10/22
 * @param
 * @return
 */
func (sb *SearchBuilder) Page(page int) *SearchBuilder {
	sb.page = page
	return sb
}

/**
 * 设置一页数量
 * @param
 * @return
 */
func (sb *SearchBuilder) Size(size int) *SearchBuilder {
	sb.size = size
	return sb
}

/**
 * 设置debug
 * @param where
 * @return
 */
func (sb *SearchBuilder) Debug() *SearchBuilder {
	sb.debug = true
	return sb
}

/**
 * 设置limit
 * @param where
 * @return
 */
func (sb *SearchBuilder) Limit(limit int) *SearchBuilder {
	sb.limit = limit
	return sb
}

/**
 * 构建完整searcher
 * @param where
 * @return
 */
func (sb *SearchBuilder) Build() *SearchBuilder {
	sb.ParseWhere()
	return sb
}

/**
 * 获取页码
 * @param
 * @return
 */
func (sb *SearchBuilder) GetSize() int {
	if sb.size <= 0 {
		return defaultSize
	}

	mSize := sb.GetMaxSize()
	if sb.size > mSize {
		return mSize
	}
	return sb.size
}

/**
 * 获取页数
 * @param
 * @return
 */
func (sb *SearchBuilder) GetPage() int {
	if sb.page <= 0 {
		return 1
	}

	mPage := sb.GetMaxPage()
	if sb.page > mPage {
		return mPage
	}
	return sb.page
}

/**
 * 允许最大页数
 * @param
 * @return
 */
func (sb *SearchBuilder) GetMaxPage() int {
	if sb.maxPage > 0 {
		return sb.maxPage
	}
	return maxPage
}

/**
 * 允许最大页码
 * @param
 * @return
 */
func (sb *SearchBuilder) GetMaxSize() int {
	if sb.maxSize > 0 {
		return sb.maxSize
	}
	return maxSize
}

/**
 * 获取排序
 * @param
 * @return
 */
func (sb *SearchBuilder) GetSort() map[string]string {
	return sb.sort
}

/**
 * 获取分组
 * @param
 * @return
 */
func (sb *SearchBuilder) GetGroup() string {
	if sb.group == nil {
		return ""
	}
	return sb.group.(string)
}

/**
 * 是否debug
 * @param
 * @return
 */
func (sb *SearchBuilder) getDebug() bool {
	return sb.debug
}

/**
 * 获取条件
 * @param
 * @return
 */
func (sb *SearchBuilder) GetWhere() interface{} {
	if sb.where == nil {
		return make(map[string]interface{})
	}
	return sb.where
}

/**
 * 获取解析之后的条件
 * @param
 * @return
 */
func (sb *SearchBuilder) GetParsedWhere() map[string]interface{} {
	return sb.parsedWhere
}

/**
 * 获取解析之后的条件
 * @param
 * @return
 */
func (sb *SearchBuilder) GetParsedValue() map[string]interface{} {
	return sb.parsedValue
}

/**
 * 构建完整searcher
 * @param where
 * @return
 */
func (sb *SearchBuilder) ParseWhere() {
	for k, v := range sb.where {
		sb.parseOperator(k, v)
	}
}

/**
 * 构建完整searcher
 * @param where
 * @return
 */
func (sb *SearchBuilder) parseOperator(key string, val interface{}) string {
	index := strings.LastIndex(key, "_")
	if index <= 0 {
		return OPERATE_EQ
	}

	op := key[index+1:]
	op = strings.ToLower(op)
	rKey := key[0:index]

	// 设置解析值
	sb.parsedValue[rKey] = val

	switch op {
	case OPERATE_GT:
		sb.parsedWhere[rKey] = rKey + " > ? "
	case OPERATE_GTE:
		sb.parsedWhere[rKey] = rKey + " >= ? "
	case OPERATE_LT:
		sb.parsedWhere[rKey] = rKey + " < ? "
	case OPERATE_LTE:
		sb.parsedWhere[rKey] = rKey + " <= ? "
	case OPERATE_EQ:
		sb.parsedWhere[rKey] = rKey + " = ? "
	case OPERATE_NE:
		sb.parsedWhere[rKey] = rKey + " <> ? "
	case OPERATE_IN:
		sb.parsedWhere[rKey] = rKey + " IN (?) "
	case OPERATE_NOT_IN:
		sb.parsedWhere[rKey] = rKey + " NOT IN (?) "
	case OPERATE_LIKE:
		sb.parsedWhere[rKey] = rKey + " LIKE ? "
	case OPERATE_NOT_LIKE:
		sb.parsedWhere[rKey] = rKey + " NOT LIKE ? "
	case OPERATE_EXIST:
		sb.parsedWhere[rKey] = rKey + " EXIST (?) "
	}

	return op
}
