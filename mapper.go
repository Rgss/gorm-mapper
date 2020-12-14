package gormmapper

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strings"
)

// db instance
var db *gorm.DB

// mapper
type Mapper struct {
	ModelEntity  GormMapperEntity
	Builder      *Searcher
	gdb          *gorm.DB
	updateFields []string // 更新字段
	selectFields []string // 查询字段

	debug          bool   // 是否debug
	databaseSource string // 数据库源

	aliasName string      // 表别名
	join      *MapperJoin // 连表数据
}

/**
 * 构建映射器
 * @author  imp.zhang
 * @date    2020/10/22
 * @param
 * @return
 */
func MapperBuilder() *Mapper {
	return &Mapper{
		join: &MapperJoin{},
	}
}

/**
 * 设置实体
 * @author  imp.zhang
 * @date    2020/10/22
 * @param
 * @return
 */
func (m *Mapper) Model(entity GormMapperEntity) *Mapper {
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
func (m *Mapper) Insert(entity interface{}) (int64, error) {
	d := m.db().Debug().Save(entity)
	if d.Error != nil {
		return 0, d.Error
	}

	return d.RowsAffected, nil
}

/**
 * 选择性插入数据
 * @param	entity	struct|map
 * @return
 */
func (m *Mapper) InsertSelective(entity interface{}) (int64, error) {
	d := m.db().Debug().Create(entity)
	if d.Error != nil {
		return 0, d.Error
	}

	return d.RowsAffected, nil
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
	defer m.afterBehaviourCallback()
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
		defer m.afterBehaviourCallback()
		return d.Error
	}

	e := d.Limit(1).Find(entities).Error
	defer m.afterBehaviourCallback()
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
	defer m.afterBehaviourCallback()

	return e
}

/**
 * 根据Searcher分页查询
 * @param
 * @return
 */
func (m *Mapper) SelectPageBySearcher(builder *Searcher, entities interface{}) (*Pager, error) {
	d := m.buildSearcher(builder)

	pager := PagerBuilder()
	if d.Error != nil {
		defer m.afterBehaviourCallback()
		return pager, d.Error
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
	defer m.afterBehaviourCallback()

	return pager, e
}

/**
 * 根据主键更新，支持struct、map。
 * @param   id
 * @param   entity	更新数据 struct | map
 * @return  int
 */
func (m *Mapper) UpdateByPrimaryKey(id int, entity interface{}) (int64, error) {
	where := map[string]interface{}{"id": id}
	d := m.db()
	if m.debug {
		d = d.Debug()
	}

	if len(m.updateFields) > 0 {
		if m.ModelEntity == nil {
			log.Printf("the modelentity is not initialized")
			return 0, nil
		}

		v := m.ParseUpdateValue(entity)
		if len(strings.TrimSpace(v.Update)) < 0 {
			return 0, errors.New("no data update")
		}

		vV := v.Value
		vV = append(vV, id)
		d = d.Exec(" UPDATE "+m.ModelEntity.TableName()+" SET "+v.Update+" WHERE id = ?", vV...)
	} else {
		en := toMap(entity, true)
		d = d.Model(m.ModelEntity).Where(where).Updates(en)
	}

	defer m.afterBehaviourCallback()
	if d.Error != nil {
		return 0, d.Error
	}
	return d.RowsAffected, nil
}

/**
 * 根据主键更新，选择性更新
 * @param
 * @return
 */
func (m *Mapper) UpdateSelectiveByPrimaryKey(id int, entity interface{}) (int64, error) {
	d := m.db()
	if m.debug {
		d = d.Debug()
	}

	where := map[string]interface{}{"id": id}
	d.Where(where).Updates(entity)
	defer m.afterBehaviourCallback()
	if d.Error != nil {
		return 0, d.Error
	}
	return d.RowsAffected, nil
}

/**
 * 根据SearchBuilder更新, 配合PreUpdateFields一起使用，处理为空字段问题。
 * 注：本方法不支持全局更新，无条件限制时，将返回ErrMissingWhereClause错误。
 * @param
 * @return
 */
func (m *Mapper) UpdateBySearcher(builder *Searcher, entity interface{}) (int64, error) {
	d := m.buildSearcher(builder)
	if len(m.updateFields) > 0 {
		if m.ModelEntity == nil {
			log.Printf("the modelentity is not initialized")
			return 0, nil
		}

		v := m.ParseUpdateValueBySearcher(builder, entity)
		d = d.Exec(" UPDATE "+m.ModelEntity.TableName()+" SET "+v.Update+" WHERE "+v.Where, v.Value...)
	} else {
		d = d.Updates(entity)
	}

	defer m.afterBehaviourCallback()
	if d.Error != nil {
		return 0, d.Error
	}
	return d.RowsAffected, nil
}

/**
 * 根据Searcher更新
 * 注：本方法不支持全局更新，无条件限制时，将返回ErrMissingWhereClause错误。
 * @param
 * @return
 */
func (m *Mapper) UpdateSelectiveBySearcher(builder *Searcher, entity interface{}) (int64, error) {
	d := m.buildSearcher(builder)
	d = d.UpdateColumns(entity)
	defer m.afterBehaviourCallback()
	if d.Error != nil {
		return 0, d.Error
	}
	return d.RowsAffected, nil
}

/**
 * 根据Searcher更新
 * 注：本方法支持全局更新
 * @param
 * @return
 */
func (m *Mapper) UpdateGlobalBySearcher(builder *Searcher, entity interface{}) (int64, error) {
	d := m.buildSearcher(builder)
	d = d.Session(&gorm.Session{AllowGlobalUpdate: true}).Updates(entity)
	defer m.afterBehaviourCallback()
	if d.Error != nil {
		return 0, d.Error
	}
	return d.RowsAffected, nil
}

/**
 * 根据Searcher删除
 * @param
 * @return
 */
func (m *Mapper) DeleteBySearcher(builder *Searcher) (int64, error) {
	d := m.buildSearcher(builder)
	e := d.Delete(m.ModelEntity).Error
	defer m.afterBehaviourCallback()
	if e != nil {
		return 0, e
	}
	return d.RowsAffected, nil
}

/**
 * 根据Searcher删除， 允许全局删除
 * @param
 * @return
 */
func (m *Mapper) DeleteGlobalBySearcher(builder *Searcher) (int64, error) {
	d := m.buildSearcher(builder)
	e := d.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m.ModelEntity).Error
	defer m.afterBehaviourCallback()
	if e != nil {
		return 0, e
	}
	return d.RowsAffected, nil
}

/**
 * 根据主键删除
 * @param
 * @return
 */
func (m *Mapper) DeleteByPrimaryKey(id int) (int64, error) {
	where := map[string]interface{}{"id": id}
	d := m.db().Where(where).Delete(m.ModelEntity)
	defer m.afterBehaviourCallback()
	if d.Error != nil {
		return 0, d.Error
	}
	return d.RowsAffected, nil
}

/**
 * 根据Searcher统计数量
 * @param
 * @return
 */
func (m *Mapper) CountBySearcher(builder *Searcher) (int64, error) {
	var count int64
	d := m.buildSearcher(builder)
	e := d.Count(&count).Error
	defer m.afterBehaviourCallback()
	if e != nil {
		return 0, e
	}
	return count, nil
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
func (m *Mapper) PreUpdateFields(entity GormMapperEntity, args ...interface{}) *Mapper {
	m.ModelEntity = entity
	fields := make([]string, 0)
	for _, v := range args {
		t := reflect.TypeOf(v).Kind()
		if t == reflect.Slice {
			fields = v.([]string)
		} else {
			fields = append(fields, v.(string))
		}
	}
	m.updateFields = fields
	return m
}

/**
 * 设置查询返回条件
 * @param   fields 读取字段
 * @return
 */
func (m *Mapper) PreSelectFields(entity GormMapperEntity, args ...interface{}) *Mapper {
	m.ModelEntity = entity
	fields := make([]string, 0)
	for _, v := range args {
		t := reflect.TypeOf(v).Kind()
		if t == reflect.Slice {
			fields = v.([]string)
		} else {
			fields = append(fields, v.(string))
		}
	}
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

	m.ModelEntity = builder.GetEntity()
	d = d.Model(m.ModelEntity)

	// where
	queryValue := m.ParseQueryAndValueBySearcher(builder)
	if len(queryValue.Query) > 0 {
		d = d.Where(queryValue.Query, queryValue.Value...)
	}

	// order
	sortQuery := m.ParseSortBySearcher(builder)
	if len(sortQuery) > 0 {
		d = d.Order(sortQuery)
	}

	// debug
	if builder.getDebug() || m.debug {
		d = d.Debug()
	}

	return d
}

// 搜索条件
type searcherQueryValue struct {
	Query string
	Value []interface{}
}

// 更新数据
type searcherUpdateValue struct {
	Update string        // 更新字段
	Where  string        // 更新条件
	Value  []interface{} // 更新值
}

/**
 * 解析条件
 * @param
 * @return
 */
func (m *Mapper) ParseQueryAndValueBySearcher(builder *Searcher) *searcherQueryValue {
	parsedWhere := builder.GetParsedWhere()
	if len(parsedWhere) <= 0 {
		return &searcherQueryValue{}
	}

	parsedValue := builder.GetParsedValue()
	whereValue := make([]interface{}, 0)
	whereQuery := ""
	for k, v := range parsedWhere {
		w := v.(string)
		w = toDBColumnName(w)
		whereQuery += "" + w + " AND "

		value := parsedValue[k]
		whereValue = append(whereValue, value)
	}

	whereQuery = strings.Trim(whereQuery, " AND ")
	queryWhere := &searcherQueryValue{
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
func (m *Mapper) ParseUpdateValue(entity interface{}) *searcherUpdateValue {
	vData := make([]interface{}, 0)
	kData := ""
	updateFields := m.updateFields
	mm := toMap(entity, true)
	for _, v := range updateFields {
		k := toDBColumnName(v)
		if _, exists := mm[k]; exists {
			kData += k + " = ?,"
			vData = append(vData, mm[k])
		}
	}

	value := &searcherUpdateValue{
		Update: strings.Trim(kData, ","),
		Value:  vData,
	}

	return value
}

/**
 * 解析更新字段
 * @param
 * @return
 */
func (m *Mapper) ParseUpdateValueBySearcher(searcher *Searcher, entity interface{}) *searcherUpdateValue {
	sqv := m.ParseQueryAndValueBySearcher(searcher)
	uv := m.ParseUpdateValue(entity)

	v := uv.Value
	v = append(v, sqv.Value...)
	suv := &searcherUpdateValue{
		Update: uv.Update,
		Where:  sqv.Query,
		Value:  v,
	}
	return suv
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
	m.ModelEntity = nil
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
