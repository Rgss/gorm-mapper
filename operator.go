package gormmapper

const OPERATE_EQ = "eq"

const OPERATE_NE = "ne"

const OPERATE_GT = "gt"

const OPERATE_GTE = "gte"

const OPERATE_LT = "lt"

const OPERATE_LTE = "lte"

const OPERATE_IN = "in"

const OPERATE_NOT_IN = "notIn"

const OPERATE_LIKE = "like"

const OPERATE_NOT_LIKE = "notLike"

const OPERATE_EXIST = "exist"

// 操作符
type Operator struct {
	Key   string
	Value interface{}
	Op    string
}

/**
 * =
 * @param
 * @return
 */
func OperatorEQ(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_EQ,
	}
}

/**
 * <>
 * @param
 * @return
 */
func OperatorNE(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_NE,
	}
}

/**
 * >
 * @param
 * @return
 */
func OperatorGT(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_GT,
	}
}

/**
 * >=
 * @param
 * @return
 */
func OperatorGTE(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_GTE,
	}
}

/**
 * <
 * @param
 * @return
 */
func OperatorLT(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_LT,
	}
}

/**
 * <=
 * @param
 * @return
 */
func OperatorLTE(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_LTE,
	}
}

/**
 * in
 * @param
 * @return
 */
func OperatorIN(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_IN,
	}
}

/**
 * not in
 * @param
 * @return
 */
func OperatorNOTIN(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_NOT_IN,
	}
}

/**
 * like
 * @param
 * @return
 */
func OperatorLIKE(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_LIKE,
	}
}

/**
 * not like
 * @param
 * @return
 */
func OperatorNOTLIKE(key string, val interface{}) Operator {
	return Operator{
		Key:   key,
		Value: val,
		Op:    OPERATE_NOT_LIKE,
	}
}
