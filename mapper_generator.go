package gormmapper

import (
	"log"
	"os"
	"regexp"
	"strings"
)

// mapper 生成器
// 主要用于根据表结构生成实体struct, 数据库映射层
type MapperGenerator struct {
	mapper        Mapper
	tables        []string
	entityPath    string
	entityPackage string
}

// 表结构
type TableObject struct {
	Columns    []*TableColumn
	Comment    string
	Name       string
	PlainSQL   string
	PrimaryKey string
}

// 表字段
type TableColumn struct {
	Name          string
	Type          string
	Default       string
	NotNull       bool
	PrimaryKey    bool
	AutoIncrement int64
	Index         bool
	IndexName     string
	IndexType     string
	Comment       string
	Charset       string
	PlainSQL      string
}

/**
 * 实例化
 * @param
 * @return
 */
func MapperGeneratorBuilder(mapper Mapper) *MapperGenerator {
	return &MapperGenerator{
		mapper: mapper,
	}
}

/**
 * 实体包名
 * @param
 * @return
 */
func (mg *MapperGenerator) EntityPackage(name string) {
	mg.entityPackage = name
}

/**
 * 实体包名路径
 * @param
 * @return
 */
func (mg *MapperGenerator) EntityPath(name string) {
	mg.entityPath = name
}

/**
 * 启动
 * @param
 * @return
 */
func (mg *MapperGenerator) Start() {
	tables := mg.ShowTables()
	for _, t := range tables {
		log.Printf("t: %v", t)
		createTableSql := mg.ShowCreateTableSql(t)
		to := mg.ParseSQL(createTableSql)
		mg.CreateEntity(to)

	}
}

/**
 * 展示所有数据表
 * @param
 * @return
 */
func (mg *MapperGenerator) ShowTables() []string {
	tables := make([]string, 0)
	mg.mapper.DB().Debug().Raw("SHOW TABLES").Scan(&tables)

	log.Printf("tables: %v", tables)

	return tables
}

/**
 * 数据表创建sql
 * @param
 * @return
 */
func (mg *MapperGenerator) ShowCreateTableSql(tablename string) string {
	t := make(map[string]interface{})
	mg.mapper.DB().Debug().Raw("SHOW CREATE TABLE " + tablename).Scan(&t)
	s := ""
	if _, ok := t["Create Table"]; ok {
		s = t["Create Table"].(string)
	}
	return s
}

/**
 * 解析sql
 * @param
 * @return
 */
func (mg *MapperGenerator) ParseSQL(s string) *TableObject {
	cds := strings.Split(s, "\n")

	// primary key
	primaryKey := ""
	reg := regexp.MustCompile("PRIMARY KEY \\(`(.*?)`\\)")
	data := reg.FindAllStringSubmatch(s, -1)
	if len(data) > 0 {
		primaryKey = data[0][1]
	}

	tableObject := &TableObject{
		PlainSQL:   s,
		PrimaryKey: primaryKey,
	}
	tableColumns := make([]*TableColumn, 0)

	for k, v := range cds {
		if k == 0 {
			tableName := mg.ParseTableName(v)
			tableObject.Name = tableName
		} else {
			tc := mg.ParseColumn(v, tableObject)
			if tc != nil {
				tableColumns = append(tableColumns, tc)
			}

		}
	}

	tableObject.Columns = tableColumns
	return tableObject
}

/**
 * 解析表名
 * @param
 * @return
 */
func (mg *MapperGenerator) ParseTableName(s string) string {
	s = strings.Replace(s, "`", "", -1)
	reg := regexp.MustCompile(`CREATE TABLE (.*?) \(`)
	data := reg.FindAllStringSubmatch(s, -1)
	if len(data) < 0 {
		return ""
	}

	tableName := strings.TrimSpace(data[0][1])
	return tableName
}

/**
 * 解析字段
 * @param
 * @return
 */
func (mg *MapperGenerator) ParseColumn(s string, tableObject *TableObject) *TableColumn {
	s = strings.TrimSpace(s)
	// 非·字符开头，则非表字段
	if strings.Index(s, "`") != 0 {
		return nil
	}

	tableColumn := &TableColumn{}
	tableColumn.PlainSQL = s

	cds := strings.Split(s, " ")
	if len(cds) <= 1 {
		return nil
	}

	columnName := cds[0]
	tableColumn.Name = strings.Replace(columnName, "`", "", -1)
	columnType := cds[1]
	tableColumn.Type = columnType

	if strings.Compare(tableObject.PrimaryKey, tableColumn.Name) == 0 {
		tableColumn.PrimaryKey = true
	}

	return tableColumn
}

/**
 * 创建实体
 * @param
 * @return
 */
func (mg *MapperGenerator) CreateEntity(tableObject *TableObject) {
	tpl := ""
	tpl += mg.TemplateFileHeader()
	tpl += mg.TemplateFile()

	structName := mg.FormatColumnAndTagName(tableObject.Name, true)
	tableName := strings.ToLower(tableObject.Name)
	tpl = strings.Replace(tpl, "{{StructName}}", structName, -1)
	tpl = strings.Replace(tpl, "{{TableName}}", tableName, -1)

	columnTpl := mg.TemplateColumn(tableObject.Columns)
	tpl = strings.Replace(tpl, "{{columns}}", columnTpl, -1)

	log.Printf("\r\n\r\ntpl:\r\n%v", tpl)

	path := mg.SavePath(tableObject.Name)
	mg.SaveFile(path, tpl)
}

/**
 * 保存文件
 * @param
 * @return
 */
func (mg *MapperGenerator) SaveFile(path string, content string) bool {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("create entity file error: %v", err.Error())
		return false
	}

	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return false
	}

	return true
}

/**
 * 保存路径
 * @param
 * @return
 */
func (mg *MapperGenerator) SavePath(name string) string {
	mg.CheckPath()
	return mg.entityPath + "/" + name + ".go"
}

/**
 * 检查保存目录
 * @param
 * @return
 */
func (mg *MapperGenerator) CheckPath() {
	s, err := os.Stat(mg.entityPath)
	if err != nil || !s.IsDir() {
		err = os.MkdirAll(mg.entityPath, os.ModePerm)
		log.Printf("create path %v error: %v", mg.entityPath, err)
	}
}

/**
 * 文件头部
 * @param
 * @return
 */
func (mg *MapperGenerator) TemplateFileHeader() string {
	comment := "package " + mg.entityPackage + "\r\n"
	return comment
}

/**
 * 文件模版
 * @param
 * @return
 */
func (mg *MapperGenerator) TemplateFile() string {
	tpl := ""
	tpl += "\r\n// gorm-mapper auto generate entity\r\n"
	tpl += "// struct \r\n"
	tpl += "type {{StructName}} struct {\r\n"
	tpl += "{{columns}}"
	tpl += "}"

	tpl += "\r\n\r\n// tablename"
	tpl += "\r\nfunc (e {{StructName}}) TableName() string {\r\n"
	tpl += "    return \"{{TableName}}\"\r\n"
	tpl += "}\r\n\r\n"
	return tpl
}

/**
 * 字段模版
 * @param
 * @return
 */
func (mg *MapperGenerator) TemplateColumn(columns []*TableColumn) string {
	tpl := ""
	for _, v := range columns {
		if v == nil || len(v.Name) <= 0 {
			continue
		}

		varType := mg.ConvertType(v.Type)
		columnTag := mg.ColumnStructTag(v)
		column := mg.FormatColumnAndTagName(v.Name, true)
		tpl += "    " + column + "\t" + varType + columnTag + "\r\n"
	}
	return tpl
}

/**
 * struct tag
 * @param
 * @return
 */
func (mg *MapperGenerator) ColumnStructTag(tableColumn *TableColumn) string {
	columnTag := mg.gormColumn(tableColumn)

	jsonName := mg.FormatColumnAndTagName(tableColumn.Name, false)
	columnTag += " json:\"" + jsonName + "\""
	columnTag += " form:\"" + jsonName + "\""
	return columnTag + "`"
}

/**
 * gorm tag
 * @param
 * @return
 */
func (mg *MapperGenerator) gormColumn(tableColumn *TableColumn) string {
	tpl := "\t`gorm:\""

	if tableColumn.PrimaryKey {
		tpl += "primary_key:" + tableColumn.Name + ";"
	} else {
		tpl += "column:" + tableColumn.Name + ";"
	}

	// log.Printf("tableColumn.PlainSQL: %v", tableColumn.PlainSQL)

	// not null
	if strings.Contains(tableColumn.PlainSQL, "NOT NULL") {
		tpl += " NOT NULL;"
	}

	if strings.Contains(tableColumn.PlainSQL, "AUTO_INCREMENT") {
		tpl += " AUTO_INCREMENT;"
	}

	// default value
	if strings.Contains(tableColumn.PlainSQL, "DEFAULT ") {
		v := ""
		reg := regexp.MustCompile(`DEFAULT '(.*?)'`)
		data := reg.FindAllStringSubmatch(tableColumn.PlainSQL, -1)
		if len(data) <= 0 {
			if strings.Compare(tableColumn.Type, "string") == 0 {
				v = "''"
			} else {
				v = "0"
			}
		} else {
			if strings.Compare(tableColumn.Type, "string") == 0 {
				v = "'" + data[0][1] + "'"
			} else {
				v = data[0][1]
			}
		}

		tpl += " DEFAULT " + v + ";"
	}

	tpl += "\""
	return tpl
}

/**
 * mysql类型转为go类型
 * @param
 * @return
 */
func (mg *MapperGenerator) ConvertType(s string) string {
	gType := "string"
	s = strings.ToUpper(s)
	for k, v := range MapperMysqlType {
		key := strings.ToUpper(k)
		if strings.Contains(s, key) {
			gType = v
			break
		}
	}
	return gType
}

/**
 * 下划线转驼峰 首字母大写
 * @param
 * @return
 */
func (mg *MapperGenerator) FormatColumnAndTagName(str string, first bool) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		if first || (!first && y != 0) {
			for i := 0; i < len(vv); i++ {
				if i == 0 && vv[i] >= 97 && vv[i] <= 122 {
					vv[i] -= 32
					upperStr += string(vv[i])
				} else {
					upperStr += string(vv[i])
				}
			}
		}
	}

	if !first {
		upperStr = temp[0] + upperStr
	}

	return upperStr
}
