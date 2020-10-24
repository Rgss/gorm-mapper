package gormmapper

/**
 * 搜索条件
 */
type Where struct {
	Map
}

/**
 * 条件初始化
 * @date    2020/10/22
 * @param
 * @return
 */
func WhereBuilder() *Where {
	return &Where{
		Map{dict: make(map[string]interface{})},
	}
}

/**
 * 添加操作符
 * @param
 * @return
 */
func (w *Where) AddOperator(op Operator) *Where {
	key := op.Key + "_" + op.Op
	w.Put(key, op.Value)
	return w
}
