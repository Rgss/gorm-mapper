package gormmapper

import (
	"gorm.io/gorm"
	"log"
	"reflect"
	"strings"
)

// db instance
var db *gorm.DB

// mapper
type Mapper struct {
	ModelEntity  interface{}
	Builder      *Searcher
	gdb          *gorm.DB
	updateFields []string // 更新字å段
	selectFields []string // 查询字段

	debug          bool   // 是否debug
	databaseSource string // 数据库源
}

/**
 * 构建映射器
 * @author  imp.zhang
 * @date    2020/10/22
 * @param
 * @return
 */
func MapperBuilder() *Mapper {
	return &Mapper{}
}

/**
 * 设置实体
 * @author  imp.zhang
 * @date    2020/10/22
 * @param
 * @return
 */
func (m *Mapper) Model(entity interface{}) *Mapper {
	m.ModelEntity = entity
	return m
}

/**
 * 开启debug
 * @param
 * @return
 */
func (m *Mapper) Debug() *Mapper {
	m.debug = true
	return m
}

/**
 * 切换数据库源
 * @param
 * @return
 */
func (m *Mapper) DatabaseSource(name string) *Mapper {
	m.databaseSource = name
	return m
}

/**
 * 插入数据
 * @param	entity	struct|map
 * @return
 */
func (m *Mapper) Insert(entity interface{}) int64 {
	d := m.db().Debug().Save(entity)
	if d.Error != nil {
		return 0
	}

	return d.RowsAffected
}

/**
 * 选择性插入数据
 * @param	entity	struct|map
 * @return
 */
func (m *Mapper) InsertSelective(entity interface{}) int64 {
	d := m.db().Debug().Create(entity)
	if d.Error != nil {
		return 0
	}

	return d.RowsAffected
}

/**
 * 设置搜索条件
 * @param
 * @return
 */
func (m *Mapper) Searcher(builder *Searcher) *Mapper {
	m.Builder = builder
	return m
}

/**
 * 根据主键查询
 * @date
 * @param
 * @return
 */
func (m *Mapper) SelectByPrimaryKey(id int, entity interface{}) error {
	where := map[string]interface{}{"id": id}
	e := m.db().Model(entity).Where(where).Find(entity).Error
	m.afterBehaviourCallback()
	return e
}

/**
 * 根据Searcher查询单个数据
 * @param
 * @return
 */
func (m *Mapper) SelectOneBySearcher(builder *Searcher, entities interface{}) error {
	d := m.buildSearcher(builder)
	if d.Error != nil {
		m.afterBehaviourCallback()
		return d.Error
	}

	e := d.Limit(1).Find(entities).Error
	m.afterBehaviourCallback()
	return e
}

/**
 * 根据Searcher查询
 * @param
 * @return
 */
func (m *Mapper) SelectBySearcher(builder *Searcher, entities interface{}) error {
	d := m.buildSearcher(builder)
	d = d.Limit(builder.GetSize())

	e := d.Find(entities).Error
	m.afterBehaviourCallback()

	return e
}

/**
 * 根据Searcher分页查询
 * @param
 * @return
 */
func (m *Mapper) SelectPageBySearcher(builder *Searcher, entities interface{}) (error, *Pager) {
	d := m.buildSearcher(builder)

	pager := PagerBuilder()
	if d.Error != nil {
		m.afterBehaviourCallback()
		return d.Error, pager
	}

	// count
	var count int64
	d.Count(&count)
	pager.Page = builder.GetPage()
	pager.Size = builder.GetSize()
	pager.Total = count

	// rows
	start := (pager.GetPage() - 1) * pager.GetSize()
	d = d.Limit(builder.GetSize())
	e := d.Offset(start).Find(entities).Error
	m.afterBehaviourCallback()

	return e, pager
}

/**
 * 根据主键更新，支持struct、map。
 * @param   id
 * @param   entity	更新数据 struct | map
 * @return  int
 */
func (m *Mapper) UpdateByPrimaryKey(id int, entity interface{}) int64 {
	where := map[string]interface{}{"id": id}
	d := m.db().Debug()

	if len(m.updateFields) > 0 {
		value := m.ParseUpdateValue(entity)
		v := value.Value
		v = append(v, id)
		d = d.Exec(" UPDATE user SET "+value.Query+" WHERE id = ?", v...)
	} else {
		d = d.Where(where).Updates(entity)
	}

	m.afterBehaviourCallback()
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 转换成map数据结构
 * @param
 * @return
 */
func (m *Mapper) toMap(entity interface{}) map[string]interface{} {
	mm := make(map[string]interface{})
	valueOf := reflect.ValueOf(entity)
	elem := valueOf.Elem()
	switch valueOf.Kind() {
	case reflect.Struct:
		break
	case reflect.Map:
		break
	case reflect.Ptr:
		for i := 0; i < elem.NumField(); i++ {
			vN := elem.Type().Field(i).Name
			vV := elem.Field(i).Interface()
			mm[vN] = vV

			if vV == nil {
				log.Printf("%v empty.", vN)
			}
		}
		break
	default:
	}
	//log.Printf("mm: %v", mm)
	return mm
}

/**
 * 根据主键更新，选择性更新
 * @param
 * @return
 */
func (m *Mapper) UpdateSelectiveByPrimaryKey(id int, entity interface{}) int64 {
	where := map[string]interface{}{"id": id}
	d := m.db().Debug().Where(where).Updates(entity)
	m.afterBehaviourCallback()
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据SearchBuilder更新, 配合PreUpdateFields一起使用，处理为空字段问题。
 * 注：本方法不支持全局更新，无条件限制时，将返回ErrMissingWhereClause错误。
 * @param
 * @return
 */
func (m *Mapper) UpdateBySearcher(builder *Searcher, entity interface{}) int64 {
	d := m.buildSearcher(builder)
	d = d.Updates(entity)
	m.afterBehaviourCallback()
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据Searcher更新
 * 注：本方法不支持全局更新，无条件限制时，将返回ErrMissingWhereClause错误。
 * @param
 * @return
 */
func (m *Mapper) UpdateSelectiveBySearcher(builder *Searcher, entity interface{}) int64 {
	d := m.buildSearcher(builder)
	d = d.UpdateColumns(entity)
	m.afterBehaviourCallback()
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据Searcher更新
 * 注：本方法支持全局更新
 * @param
 * @return
 */
func (m *Mapper) UpdateGlobalBySearcher(builder *Searcher, entity interface{}) int64 {
	d := m.buildSearcher(builder)
	d = d.Session(&gorm.Session{AllowGlobalUpdate: true}).Updates(entity)
	m.afterBehaviourCallback()
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据Searcher删除
 * @param
 * @return
 */
func (m *Mapper) DeleteBySearcher(builder *Searcher) int64 {
	d := m.buildSearcher(builder)
	e := d.Delete(m.ModelEntity).Error
	m.afterBehaviourCallback()
	if e != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据Searcher删除， 允许全局删除
 * @param
 * @return
 */
func (m *Mapper) DeleteGlobalBySearcher(builder *Searcher) int64 {
	d := m.buildSearcher(builder)
	e := d.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m.ModelEntity).Error
	m.afterBehaviourCallback()
	if e != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据主键删除
 * @param
 * @return
 */
func (m *Mapper) DeleteByPrimaryKey(id int) int64 {
	where := map[string]interface{}{"id": id}
	d := m.db().Where(where).Delete(m.ModelEntity)
	m.afterBehaviourCallback()
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据Searcher统计数量
 * @param
 * @return
 */
func (m *Mapper) CountBySearcher(builder *Searcher) int64 {
	var count int64
	d := m.buildSearcher(builder)
	e := d.Count(&count).Error
	m.afterBehaviourCallback()
	if e != nil {
		return 0
	}
	return count
}

/**
 * 自增
 * @param
 * @return
 */
func (m *Mapper) IncrBySearcher(column string, step int, where interface{}) {

}

/**
 * 自减
 * @param
 * @return
 */
func (m *Mapper) DecrBySearcher(column string, step int, where interface{}) {

}

/**
 * 执行原生查询sql
 * @param
 * @return
 */
func (m *Mapper) RawQuery(sqlQuery string, entity interface{}) error {
	return m.db().Raw(sqlQuery).Scan(&entity).Error
}

/**
 * 执行原生sql
 * @param
 * @return
 */
func (m *Mapper) RawExec(sqlQuery string, args []interface{}) int64 {
	d := m.db().Debug().Exec(sqlQuery, args...)
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 返回db实例
 * @param
 * @return
 */
func (m *Mapper) db() *gorm.DB {
	return GDB(m.databaseSource)
}

/**
 * 返回db实例
 * @param
 * @return
 */
func (m *Mapper) DB() *gorm.DB {
	return m.db()
}

/**
 * 设置写入更新字段，主要用于更新空字符串、0值等问题
 * @param	fields 更新字段
 * @return
 */
func (m *Mapper) PreUpdateFields(args ...interface{}) *Mapper {
	fields := make([]string, 0)
	m.updateFields = fields
	return m
}

/**
 * 设置查询返回条件
 * @param   fields 读取字段
 * @return
 */
func (m *Mapper) PreSelectFields(fields []string) *Mapper {
	m.selectFields = fields
	return m
}

/**
 * 构建排序
 * @param
 * @return
 */
func (m *Mapper) buildOrder(builder *Searcher, d *gorm.DB) *gorm.DB {
	sortQuery := m.ParseSortBySearcher(builder)
	if len(sortQuery) > 0 {
		d.Order(sortQuery)
	}
	return d
}

/**
 * 构建条件
 * @param
 * @return
 */
func (m *Mapper) buildSearcher(builder *Searcher) *gorm.DB {
	d := m.db()

	if builder.Entity != nil {
		m.ModelEntity = builder.Entity
		d = d.Model(builder.Entity)
	}

	// where
	queryValue := m.ParseQueryAndValueBySearcher(builder)
	//log.Printf("queryValue: %v", queryValue)
	if len(queryValue.Query) > 0 {
		d = d.Where(queryValue.Query, queryValue.Value...)
	}

	// order
	sortQuery := m.ParseSortBySearcher(builder)
	if len(sortQuery) > 0 {
		d = d.Order(sortQuery)
	}

	// debug
	if builder.getDebug() {
		d = d.Debug()
	}

	return d
}

// 搜索条件
type SearcherQueryValue struct {
	Query string
	Value []interface{}
}

// 更新数据
type SearchBuilderUpdateValue struct {
}

/**
 * 解析条件
 * @param
 * @return
 */
func (m *Mapper) ParseQueryAndValueBySearcher(builder *Searcher) *SearcherQueryValue {
	parsedWhere := builder.GetParsedWhere()
	if len(parsedWhere) <= 0 {
		return &SearcherQueryValue{}
	}

	parsedValue := builder.GetParsedValue()
	whereValue := make([]interface{}, 0)
	whereQuery := ""
	for k, v := range parsedWhere {
		vs := v.(string)
		whereQuery += vs + " AND "

		value := parsedValue[k]
		whereValue = append(whereValue, value)
	}

	whereQuery = strings.Trim(whereQuery, " AND ")
	queryWhere := &SearcherQueryValue{
		Query: whereQuery,
		Value: whereValue,
	}
	return queryWhere
}

/**
 * 解析排序
 * @param
 * @return
 */
func (m *Mapper) ParseSortBySearcher(builder *Searcher) string {
	if len(builder.GetSort()) <= 0 {
		return ""
	}

	sortQuery := ""
	sorts := builder.GetSort()
	for k, v := range sorts {
		sortQuery += k + " " + v
	}
	return sortQuery
}

/**
 * 解析更新字段
 * @param
 * @return
 */
func (m *Mapper) ParseUpdateValue(entity interface{}) *SearcherQueryValue {
	vData := make([]interface{}, 0)
	kData := ""
	updateFields := m.updateFields
	mm := m.toMap(entity)
	log.Printf("mm: %v", mm)
	for _, v := range updateFields {
		if _, exists := mm[v]; exists {
			kData += v + " = ?,"
			vData = append(vData, mm[v])
		}
	}

	value := &SearcherQueryValue{
		Query: strings.Trim(kData, ","),
		Value: vData,
	}
	log.Printf("updateFields: %v", updateFields)
	log.Printf("mm: %v", mm)
	return value
}

/**
 * 执行完sql后回调
 * @param
 * @return
 */
func (m *Mapper) afterBehaviourCallback() {
	m.release()
}

/**
 * sql执行完之后操作
 * @param
 * @return
 */
func (m *Mapper) release() {
	m.updateFields = nil
	m.selectFields = nil
	m.debug = false
}

/**
 * db 实例
 * @param
 * @return
 */
func GDB(name string) *gorm.DB {
	return Connection(name).DB()
}
