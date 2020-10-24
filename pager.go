package gormmapper

import "math"

// 分页器
type Pager struct {
	Page      int
	Size      int
	Total     int64
	TotalPage int
}

/**
 * new pager
 * @param
 * @return
 */
func NPager() *Pager {
	return &Pager{}
}

///**
// * 设置页数
// * @param
// * @return
// */
//func (p *Pager) Page(page int) *Pager {
//	p.page = page
//	return p
//}
//
///**
// * 设置一页数量
// * @param
// * @return
// */
//func (p *Pager) Size(size int) *Pager {
//	p.size = size
//	return p
//}
//
///**
// * 设置总记录数
// * @param
// * @return
// */
//func (p *Pager) Total(total int64) *Pager {
//	p.total = total
//	return p
//}

/**
 * 当前页码
 * @param
 * @return
 */
func (p *Pager) GetPage() int {
	return p.Page
}

/**
 * 一页数量
 * @param
 * @return
 */
func (p *Pager) GetSize() int {
	return p.Size
}

/**
 * 总条数
 * @param
 * @return
 */
func (p *Pager) GetTotal() int64 {
	return p.Total
}

/**
 * 总页数
 * @param
 * @return
 */
func (p *Pager) GetTotalPage() int {
	totalPage := float64(p.Total / int64(p.Size))
	totalPage = math.Ceil(totalPage)
	return int(totalPage)
}
