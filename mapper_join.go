package gormmapper

const JoinTypeInner = 1

const JoinTypeLeft = 2

const JoinTypeRight = 3

const JoinTypeOuter = 4

// join struct
type MapperJoin struct {
	JoinId    string
	JoinType  int
	TableName string
	AliasName string
	OnWhere   map[string][]Operator
}

// On
type MapperJoinOn struct {
}

/**
 * join
 * @date    2020/11/23
 * @param
 * @return
 */
func (m *Mapper) Join(names ...string) *Mapper {
	j := newMapperJoin(JoinTypeInner, names)
	m.join = j
	return m
}

/**
 * left join
 * @date    2020/11/24
 * @param
 * @return
 */
func (m *Mapper) LeftJoin(names ...string) *Mapper {
	j := newMapperJoin(JoinTypeLeft, names)
	m.join = j
	return m
}

/**
 * right join
 * @date    2020/11/24
 * @param
 * @return
 */
func (m *Mapper) RightJoin(names ...string) *Mapper {
	j := newMapperJoin(JoinTypeRight, names)
	m.join = j
	return m
}

/**
 * out join
 * @date    2020/11/24
 * @param
 * @return
 */
func (m *Mapper) OuterJoin(names ...string) *Mapper {
	j := newMapperJoin(JoinTypeOuter, names)
	m.join = j
	return m
}

/**
 * on
 * @date    2020/11/24
 * @param
 * @return
 */
func (m *Mapper) On(joinName string, operator Operator) *Mapper {
	if _, e := m.join.OnWhere[joinName]; !e {
		m.join.OnWhere[joinName] = make([]Operator, 0)
	}

	m.join.OnWhere[joinName] = append(m.join.OnWhere[joinName], operator)
	return m
}

/**
 * newMapperJoin
 * @date    2020/11/24
 * @param
 * @return
 */
func newMapperJoin(joinType int, args []string) *MapperJoin {
	t, a := parseJoinTableNameAndAliasName(args)
	j := &MapperJoin{
		JoinType:  joinType,
		TableName: t,
		AliasName: a,
		OnWhere:   make(map[string][]Operator),
	}
	return j
}

/**
 * 解析表名和别名
 * @date    2020/11/24
 * @param
 * @return
 */
func parseJoinTableNameAndAliasName(args []string) (string, string) {
	tn := ""
	an := ""
	for k, v := range args {
		if k == 0 {
			tn = v
		}
		if k == 1 {
			an = v
		}
	}
	return tn, an
}
