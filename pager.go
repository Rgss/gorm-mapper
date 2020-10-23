package gormmapper

import "math"

// 分页器
type Pager struct {
	page int
	size int
	total int
	totalPage int
}

/**
 * 实力化
 * @param
 * @return 
 */
func NPager() *Pager {
	return &Pager{}
}

/**
 * 设置当前页码
 * @param
 * @return
 */
func (p *Pager) Page(page int) *Pager {
	p.page = page
	return p
}

/**
 * 当前页码
 * @param   
 * @return  
 */
func (p *Pager) getPage() int {
	return p.page
}

/**
 * 一页数量
 * @param   
 * @return  
 */
func (p *Pager) getSize() int {
	return p.size
}

/**
 * 总条数
 * @param   
 * @return  
 */
func (p *Pager) getTotal() int {
	return p.total
}

/**
 * 总页数
 * @param
 * @return
 */
func (p *Pager) getTotalPage() int {
	totalPage := float64(p.total / p.size)
	totalPage = math.Ceil(totalPage)
	return int(totalPage)
}