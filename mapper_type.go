package gormmapper

var MapperMysqlType = map[string]string{
	"TINYINT":   "int8",
	"SMALLINT":  "int16",
	"MEDIUMINT": "int32",
	"INT":       "int",
	"BIGINT":    "int64",
	"FLOAT":     "float",
	"DOUBLE":    "float64",
	"DECIMAL":   "float",

	"DATE":      "string",
	"TIME":      "string",
	"YEAR":      "string",
	"DATETIME":  "string",
	"TIMESTAMP": "string",

	"CHAR":       "string",
	"VARCHAR":    "string",
	"TINYBLOB":   "string",
	"TINYTEXT":   "string",
	"BLOB":       "string",
	"TEXT":       "string",
	"MEDIUMBLOB": "string",
	"MEDIUMTEXT": "string",
	"LONGBLOB":   "string",
	"LONGTEXT":   "string",
}
