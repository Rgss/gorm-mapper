package gormmapper

import (
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

type Mapper struct {
	ModelEntity interface{}
	Builder     *SearchBuilder
	db          *gorm.DB
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
 * 设置搜索条件
 * @param
 * @return
 */
func (m *Mapper) SearchBuilder(builder *SearchBuilder) *Mapper {
	m.Builder = builder
	return m;
}


type valueType interface {
}

/**
 * 根据SearchBuilder查询单个数据
 * @param
 * @return
 */
func (m *Mapper) SelectOneBySearchBuilder(builder *SearchBuilder, entities interface{}) error {
	d := m.buildSearchBuilder(builder)
	d  = d.Limit(builder.GetSize())
	if d.Error != nil {
		return d.Error
	}

	d = d.Limit(1).Find(entities);
	return d.Error;
}

/**
 * 根据SearchBuilder查询
 * @param
 * @return
 */
func (m *Mapper) SelectBySearchBuilder(builder *SearchBuilder, entities interface{}) error {
	d := m.buildSearchBuilder(builder)
	d  = d.Limit(builder.GetSize())
	if d.Error != nil {
		return d.Error
	}

	d = d.Find(entities)
	return d.Error
}

/**
 * 根据SearchBuilder分页查询
 * @param
 * @return
 */
func (m *Mapper) SelectPageBySearchBuilder(builder *SearchBuilder, entities interface{}, pager Pager) error {
	d := m.buildSearchBuilder(builder)
	d  = d.Limit(builder.GetSize())
	if d.Error != nil {
		return d.Error
	}

	d = d.Find(entities);
	return d.Error;
}

/**
 * 根据SearchBuilder更新
 * @param
 * @return
 */
func (m *Mapper) UpdateBySearchBuilder(builder *SearchBuilder) {

}

/**
 * 根据SearchBuilder删除
 * @param
 * @return
 */
func (m *Mapper) DeleteBySearchBuilder(builder *SearchBuilder) {

}

/**
 *
 * @param   
 * @return  
 */
func (m *Mapper) buildSort(builder *SearchBuilder, d *gorm.DB) *gorm.DB {
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
	d := GDB()

	if builder.Entity != nil {
		d = d.Model(builder.Entity)
	}

	// where
	queryValue := m.ParseQueryAndValueBySearchBuilder(builder)
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
	for k, v :=  range parsedWhere {
		vs := v.(string)
		whereQuery += vs + " AND "

		value := parsedValue[k]
		whereValue = append(whereValue, value)
	}

	whereQuery  = strings.Trim(whereQuery, " AND ")
	queryWhere := &SearchBuilderQueryValue{
		Query: whereQuery,
		Value: whereValue,
	}
	return queryWhere;
}

/**
 * 解析值
 * @param
 * @return
 */
func (m *Mapper) ParseQueryValueBySearchBuilder(builder *SearchBuilder) {

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
	for k, v := range sorts  {
		sortQuery +=  k + " " + v
	}
	return sortQuery
}


/**
 * 
 * @date    2020/10/22
 * @param   
 * @return  
 */
func GDB() *gorm.DB {
	return Connection().DB()
}

// create
func (m *Mapper) Insert(entity interface{}) (interface{}, error) {
	err := Connection().DB().Create(entity).Error
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// create
func (m *Mapper) MultiInsert(entity []interface{}) (interface{}, error) {
	err := Connection().DB().Create(entity).Error
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// Save
func (m *Mapper) Save(entity interface{}) (bool, error) {
	res := Connection().DB().Debug().Save(entity)
	if res.Error != nil {
		return false, res.Error
	}

	if res.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

// delete
func (m *Mapper) Delete(id int) (bool, error) {
	res := Connection().DB().Where("id = ?", strconv.Itoa(id)).Limit(1).Delete(m.ModelEntity)
	if res.RowsAffected >= 1 {
		return true, nil
	}

	return false, res.Error
}

// DeleteBy
func (m *Mapper) DeleteBy(where interface{}) (bool, error) {
	res := Connection().DB().Where(where).Delete(m.ModelEntity)
	if res.RowsAffected >= 1 {
		return true, nil
	}

	return false, res.Error
}



// find
func (m *Mapper) Find(id int, entity interface{}) (interface{}, error) {
	res := Connection().DB().Where("id = ?", strconv.Itoa(id)).Limit(1).Find(entity)
	if res.Error != nil {
		return entity, res.Error
	}

	return entity, nil
}

// FindBy
func (m *Mapper) FindBy(where interface{}, entity interface{}) (interface{}, error) {
	res := Connection().DB().Where(where).Find(entity)
	if res.Error != nil {
		return entity, res.Error
	}

	return entity, nil
}

// update
func (m *Mapper) Update(id int, entity interface{}) (bool, error) {
	res := Connection().DB().Model(entity).Where("id = ?", id).Updates(entity)
	if res.Error != nil {
		return false, res.Error
	}

	//log.Printf("Update.RowsAffected: %v\n", res.RowsAffected)
	if res.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

// update
func (m *Mapper) UpdateColumn(id int, value interface{}) (bool, error) {
	//res := Connection().DB().Model(this.ModelEntity).Where("id = ?", id).UpdateColumn(value)
	//if res.Error != nil {
	//	return false, res.Error
	//}
	//
	//if res.RowsAffected > 0 {
	//	return true, nil
	//}

	return false, nil
}

// update
func (m *Mapper) UpdateColumns(id int, value interface{}, model interface{}) (bool, error) {
	res := Connection().DB().Model(model).Where("id = ?", id).UpdateColumns(value)
	if res.Error != nil {
		return false, res.Error
	}

	if res.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

// update by
func (m *Mapper) UpdateBy(where interface{}, entity interface{}) (bool, error) {
	res := Connection().DB().Model(m.ModelEntity).Where(where).Updates(entity)
	if res.Error != nil {
		return false, res.Error
	}

	log.Printf("UpdateBy.RowsAffected: %v\n", res.RowsAffected)
	if res.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func (m *Mapper) Exists(id int, entity interface{}) (bool, error) {
	_, err := m.Find(id, entity)
	if err != nil {
		return true, err
	}

	return false, nil
}

func (m *Mapper) ExistsBy(where interface{}, entity interface{}) (bool, error) {
	_, err := m.FindBy(where, entity)
	if err != nil {
		return true, err
	}

	return false, nil
}

// Count
func (m *Mapper) Count(where interface{}) (int, error) {
	var count int64
	res := Connection().DB().Where(where).Model(m.ModelEntity).Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}

	return int(count), nil
}

func (m *Mapper) IncreaseBy(column string, step int, where interface{}) {
}

func (m *Mapper) DecreaseBy(column string, step int, where interface{}) {
}
