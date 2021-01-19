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
type Searcher struct {
	page        int
	size        int
	limit       int
	group       interface{}
	sort        map[string]string
	where       *Where
	parsedWhere []*KV
	parsedValue []*KV
	fields      interface{}
	entity      GormMapperEntity
	debug       bool
	maxPage     int
	maxSize     int
}

type KV struct {
	k string
	v interface{}
}

/**
 * 创建搜索构建器
 * @param args	实体对象|缺省项
 * @return
 */
func SearcherBuilder(args ...interface{}) *Searcher {
	var entity GormMapperEntity
	if len(args) > 0 {
		entity = args[0].(GormMapperEntity)
	}

	return &Searcher{
		entity:      entity,
		where:       &Where{},
		sort:        make(map[string]string),
		parsedWhere: make([]*KV, 0),
		parsedValue: make([]*KV, 0),
	}
}

/**
 * 设置条件
 * @param where
 * @return
 */
func (sb *Searcher) Where(where *Where) *Searcher {
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
func (sb *Searcher) Sort(column string, sortType string) *Searcher {
	sb.sort[column] = strings.ToUpper(sortType)
	return sb
}

/**
 * 设置字段
 * @date    2020/10/22
 * @param
 * @return
 */
func (sb *Searcher) Field(fields interface{}) *Searcher {
	sb.fields = fields
	return sb
}

/**
 * 设置页数
 * @date    2020/10/22
 * @param
 * @return
 */
func (sb *Searcher) Page(page int) *Searcher {
	sb.page = page
	return sb
}

/**
 * 设置一页数量
 * @param
 * @return
 */
func (sb *Searcher) Size(size int) *Searcher {
	sb.size = size
	return sb
}

/**
 * 设置debug
 * @param where
 * @return
 */
func (sb *Searcher) Debug() *Searcher {
	sb.debug = true
	return sb
}

/**
 * 设置limit
 * @param where
 * @return
 */
func (sb *Searcher) Limit(limit int) *Searcher {
	sb.limit = limit
	return sb
}

/**
 * 设置实体类型
 * @param
 * @return
 */
func (sb *Searcher) Entity(entity GormMapperEntity) *Searcher {
	sb.entity = entity
	return sb
}

/**
 * 构建完整searcher
 * @param where
 * @return
 */
func (sb *Searcher) Build() *Searcher {
	sb.ParseWhere()
	return sb
}

func (sb *Searcher) GetEntity() GormMapperEntity {
	return sb.entity
}

/**
 * 获取页码
 * @param
 * @return
 */
func (sb *Searcher) GetSize() int {
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
func (sb *Searcher) GetPage() int {
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
func (sb *Searcher) GetMaxPage() int {
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
func (sb *Searcher) GetMaxSize() int {
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
func (sb *Searcher) GetSort() map[string]string {
	return sb.sort
}

/**
 * 获取分组
 * @param
 * @return
 */
func (sb *Searcher) GetGroup() string {
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
func (sb *Searcher) getDebug() bool {
	return sb.debug
}

/**
 * 获取条件
 * @param
 * @return
 */
func (sb *Searcher) GetWhere() interface{} {
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
func (sb *Searcher) GetParsedWhere() []*KV {
	return sb.parsedWhere
}

/**
 * 获取解析之后的条件
 * @param
 * @return
 */
func (sb *Searcher) GetParsedValue() []*KV {
	return sb.parsedValue
}

/**
 * 构建完整searcher
 * @param where
 * @return
 */
func (sb *Searcher) ParseWhere() {
	where := sb.where.Iterator()
	for k, v := range where {
		sb.parseOperator(k, v)
	}
}

/**
 * 构建完整searcher
 * @param where
 * @return
 */
func (sb *Searcher) parseOperator(key string, val interface{}) {
	var op string
	var rKey string
	index := strings.LastIndex(key, "_")
	if index < 0 {
		op = OPERATE_EQ
		rKey = key
	} else {
		op = key[index+1:]
		op = strings.ToLower(op)
		rKey = key[0:index]
	}

	wv := ""
	switch op {
	case OPERATE_GT:
		wv = rKey + " > ? "
	case OPERATE_GTE:
		wv = rKey + " >= ? "
	case OPERATE_LT:
		wv = rKey + " < ? "
	case OPERATE_LTE:
		wv = rKey + " <= ? "
	case OPERATE_EQ:
		wv = rKey + " = ? "
	case OPERATE_NE:
		wv = rKey + " <> ? "
	case OPERATE_IN:
		wv = rKey + " IN (?) "
	case OPERATE_NOT_IN:
		wv = rKey + " NOT IN (?) "
	case OPERATE_LIKE:
		wv = rKey + " LIKE ? "
	case OPERATE_NOT_LIKE:
		wv = rKey + " NOT LIKE ? "
	case OPERATE_EXIST:
		wv = rKey + " EXIST (?) "
	default:
		rKey = key
		//wv = rKey + " = ? "
	}

	kv := &KV{
		k: rKey,
		v: wv,
	}
	sb.parsedWhere = append(sb.parsedWhere, kv)

	vv := &KV{
		k: rKey,
		v: val,
	}
	sb.parsedValue = append(sb.parsedValue, vv)
}
