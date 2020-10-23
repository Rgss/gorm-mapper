package gormmapper

/**
 * 搜索条件
 */
type Where Map

/**
 * 条件初始化
 * @date    2020/10/22
 * @param
 * @return
 */
func NWhere() *Where {
	return &Where{}
}