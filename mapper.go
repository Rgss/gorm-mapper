package gormmapper

import (
	"gorm.io/gorm"
	"log"
	"reflect"
	"strings"
)

type Mapper struct {
	ModelEntity  interface{}
	Builder      *SearchBuilder
	gdb          *gorm.DB
	upFields     []string // 更新字段
	selectFields []string // 查询字段
}

/**
 * 构建映射器
 * @author  imp.zhang
 * @date    2020/10/22
 * @param
 * @return
 */
func NMapper() *Mapper {
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
 * 插入数据
 * @param	entity	struct|map
 * @return
 */
func (m *Mapper) Insert(entity interface{}) int64 {
	d := m.DB().Debug().Create(entity)
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
	d := m.DB().Debug().Create(entity)
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
func (m *Mapper) SearchBuilder(builder *SearchBuilder) *Mapper {
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
	e := m.DB().Model(entity).Where(where).Find(entity).Error
	return e
}

/**
 * 根据SearchBuilder查询单个数据
 * @param
 * @return
 */
func (m *Mapper) SelectOneBySearchBuilder(builder *SearchBuilder, entities interface{}) error {
	d := m.buildSearchBuilder(builder)
	if d.Error != nil {
		return d.Error
	}

	e := d.Limit(1).Find(entities).Error
	return e
}

/**
 * 根据SearchBuilder查询
 * @param
 * @return
 */
func (m *Mapper) SelectBySearchBuilder(builder *SearchBuilder, entities interface{}) error {
	d := m.buildSearchBuilder(builder)
	d = d.Limit(builder.GetSize())
	if d.Error != nil {
		return d.Error
	}

	return d.Find(entities).Error
}

/**
 * 根据SearchBuilder分页查询
 * @param
 * @return
 */
func (m *Mapper) SelectPageBySearchBuilder(builder *SearchBuilder, entities interface{}) (error, *Pager) {
	d := m.buildSearchBuilder(builder)
	d = d.Limit(builder.GetSize())

	pager := NPager()
	if d.Error != nil {
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
	e := d.Offset(start).Find(entities).Error
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
	m.toMap(entity)
	d := m.DB().Debug().Where(where).Save(entity)
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
	log.Printf("valueOf.Kind: %v", valueOf.Kind())
	log.Printf("valueOf.MapKeys: %v", valueOf.Type())
	log.Printf("valueOf.elem: %v", elem)
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

			log.Printf("name: %v, value: %v, %T", vN, vV, vV)
			if vV == nil {
				log.Printf("%v empty.", vN)
			}
		}
		break
	default:
	}
	log.Printf("mm: %v", mm)
	return mm
}

/**
 * 根据主键更新，选择性更新
 * @param
 * @return
 */
func (m *Mapper) UpdateSelectiveByPrimaryKey(id int, entity interface{}) int64 {
	where := map[string]interface{}{"id": id}
	d := m.DB().Debug().Where(where).Updates(entity)
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据SearchBuilder更新
 * 注：本方法不支持全局更新，无条件限制时，将返回ErrMissingWhereClause错误。
 * @param
 * @return
 */
func (m *Mapper) UpdateBySearchBuilder(builder *SearchBuilder, entity interface{}) int64 {
	d := m.buildSearchBuilder(builder)
	d = d.Updates(entity)
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据SearchBuilder更新
 * 注：本方法不支持全局更新，无条件限制时，将返回ErrMissingWhereClause错误。
 * @param
 * @return
 */
func (m *Mapper) UpdateSelectiveBySearchBuilder(builder *SearchBuilder, entity interface{}) int64 {
	d := m.buildSearchBuilder(builder)
	d = d.UpdateColumns(entity)
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据SearchBuilder更新
 * 注：本方法支持全局更新
 * @param
 * @return
 */
func (m *Mapper) UpdateGlobalBySearchBuilder(builder *SearchBuilder, entity interface{}) int64 {
	d := m.buildSearchBuilder(builder)
	d = d.Session(&gorm.Session{AllowGlobalUpdate: true}).Updates(entity)
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据SearchBuilder删除
 * @param
 * @return
 */
func (m *Mapper) DeleteBySearchBuilder(builder *SearchBuilder) int64 {
	d := m.buildSearchBuilder(builder)
	e := d.Delete(m.ModelEntity).Error
	if e != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据SearchBuilder删除， 允许全局删除
 * @param
 * @return
 */
func (m *Mapper) DeleteGlobalBySearchBuilder(builder *SearchBuilder) int64 {
	d := m.buildSearchBuilder(builder)
	e := d.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m.ModelEntity).Error
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
	d := m.DB().Where(where).Delete(m.ModelEntity)
	if d.Error != nil {
		return 0
	}
	return d.RowsAffected
}

/**
 * 根据searchbuilder统计数量
 * @param
 * @return
 */
func (m *Mapper) CountBySearchBuilder(builder *SearchBuilder) int64 {
	var count int64
	d := m.buildSearchBuilder(builder)
	e := d.Count(&count).Error
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
func (m *Mapper) IncrBySearchBuilder(column string, step int, where interface{}) {

}

/**
 * 自减
 * @param
 * @return
 */
func (m *Mapper) DecrBySearchBuilder(column string, step int, where interface{}) {

}

/**
 * 执行原生查询sql
 * @param
 * @return
 */
func (m *Mapper) RawQuery(sqlQuery string, entity interface{}) error {
	return m.DB().Raw(sqlQuery).Scan(&entity).Error
}

/**
 * 执行原生sql
 * @param
 * @return
 */
func (m *Mapper) RawExec(sqlQuery string) int64 {
	d := m.DB().Exec(sqlQuery)
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
	return GDB()
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
func (m *Mapper) PreUpdateFields(fields []string) *Mapper {
	m.upFields = fields
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
func (m *Mapper) buildOrder(builder *SearchBuilder, d *gorm.DB) *gorm.DB {
	sortQuery := m.ParseSortBySearchBuilder(builder)
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
func (m *Mapper) buildSearchBuilder(builder *SearchBuilder) *gorm.DB {
	d := m.DB()

	if builder.Entity != nil {
		m.ModelEntity = builder.Entity
		d = d.Model(builder.Entity)
	}

	// where
	queryValue := m.ParseQueryAndValueBySearchBuilder(builder)
	log.Printf("queryValue: %v", queryValue)
	if len(queryValue.Query) > 0 {
		d = d.Where(queryValue.Query, queryValue.Value)
	}

	// order
	sortQuery := m.ParseSortBySearchBuilder(builder)
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
type SearchBuilderQueryValue struct {
	Query string
	Value interface{}
}

/**
 * 解析条件
 * @param
 * @return
 */
func (m *Mapper) ParseQueryAndValueBySearchBuilder(builder *SearchBuilder) *SearchBuilderQueryValue {
	parsedWhere := builder.GetParsedWhere()
	if len(parsedWhere) <= 0 {
		return &SearchBuilderQueryValue{}
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
	queryWhere := &SearchBuilderQueryValue{
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
func (m *Mapper) ParseSortBySearchBuilder(builder *SearchBuilder) string {
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
	m.upFields = nil
	m.selectFields = nil
}

/**
 * db 实例
 * @param
 * @return
 */
func GDB() *gorm.DB {
	return Connection().DB()
}

func (m *Mapper) M() {
	log.Printf("mapper.m")
}
