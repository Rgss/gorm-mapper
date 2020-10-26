package gormmapper

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// mapper 生成器
// 主要用于根据表结构生成实体struct, 数据库映射层
type MapperGenerator struct {
	mapper                  Mapper   // 映射器
	tables                  []string // 表名
	entityPath              string   // 实体路径
	entityPackage           string   // 实体包名
	mapperPath              string   // 映射器路径
	mapperPackage           string   // 映射器包名
	mapperPathAutoSignleton bool     // 是否自动单例实例化
	overWrite               bool     // 是否覆盖文件
	allowTables             []string // 允许生成文件的表
}

// 表结构
type TableObject struct {
	Columns             []*TableColumn
	Comment             string
	Name                string
	PlainSQL            string
	PrimaryKey          string
	MaxColumnLength     int
	MaxColumnTypeLength int
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
 * 实体映射器包名
 * @param
 * @return
 */
func (mg *MapperGenerator) MapperPackage(name string) {
	mg.mapperPackage = name
}

/**
 * 实体映射器包名路径
 * @param
 * @return
 */
func (mg *MapperGenerator) MapperPath(name string) {
	mg.mapperPath = name
}

/**
 * 需要生成实体和映射器的表， 默认全部
 * @param
 * @return
 */
func (mg *MapperGenerator) AllowTables(tables []string) {
	mg.allowTables = tables
}

/**
 * 是否覆盖生成实体和映射器的文件
 * @param
 * @return
 */
func (mg *MapperGenerator) OverWrite(flag bool) {
	mg.overWrite = flag
}

/**
 * 自动单例初始化
 * @param
 * @return
 */
func (mg *MapperGenerator) MapperPathAutoSignleton(flag bool) {
	mg.mapperPathAutoSignleton = flag
}

/**
 * 启动
 * @param
 * @return
 */
func (mg *MapperGenerator) Start() {
	log.Printf("gormmapper generator start.")
	tables := mg.ShowTables()
	for _, t := range tables {
		s := mg.ShowCreateTableSql(t)
		to := mg.ParseSQL(s)

		// 创建实体
		if len(mg.entityPackage) > 0 && len(mg.entityPath) > 0 {
			err := mg.CreateEntity(to)
			if err != nil {
				log.Printf("create table ? err: %v", to.Name, err.Error())
			} else {
				log.Printf("create table %v ok", to.Name)
			}
		}

		// 创建映射器
		if len(mg.mapperPackage) > 0 && len(mg.mapperPath) > 0 {
			err2 := mg.CreateEntityMapper(to)
			if err2 != nil {
				log.Printf("create mapper ? err: %v", to.Name, err2.Error())
			} else {
				log.Printf("create mapper %v ok", to.Name)
			}
		}
	}
	log.Printf("gormmapper generator end.")
}

/**
 * 展示所有数据表
 * @param
 * @return
 */
func (mg *MapperGenerator) ShowTables() []string {
	tables := make([]string, 0)
	mg.mapper.DB().Raw("SHOW TABLES").Scan(&tables)
	return tables
}

/**
 * 数据表创建sql
 * @param
 * @return
 */
func (mg *MapperGenerator) ShowCreateTableSql(tablename string) string {
	t := make(map[string]interface{})
	mg.mapper.DB().Raw("SHOW CREATE TABLE " + tablename).Scan(&t)
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
	tableColumn.Name = strings.TrimSpace(tableColumn.Name)
	columnType := cds[1]
	tableColumn.Type = columnType

	// 字段最大长度
	l := len(strings.Replace(tableColumn.Name, "_", "", -1))
	if tableObject.MaxColumnLength < l {
		tableObject.MaxColumnLength = l
	}
	// 字段类型最大长度
	ct := mg.ConvertType(columnType)
	l2 := len(ct)
	if tableObject.MaxColumnTypeLength < l2 {
		tableObject.MaxColumnTypeLength = l2
	}

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
func (mg *MapperGenerator) CreateEntity(tableObject *TableObject) error {
	tpl := mg.TemplateFileHeader(mg.entityPackage)
	tpl += mg.TemplateEntityFile()

	structName := mg.FormatColumnAndTagName(tableObject.Name, true)
	tableName := strings.ToLower(tableObject.Name)
	tpl = strings.Replace(tpl, "{{StructName}}", structName, -1)
	tpl = strings.Replace(tpl, "{{TableName}}", tableName, -1)

	ctpl := mg.TemplateColumn(tableObject)
	tpl = strings.Replace(tpl, "{{columns}}", ctpl, -1)

	path := mg.SavePath(mg.entityPath, tableObject.Name)
	return mg.SaveFile(path, tpl)
}

/**
 * 创建实体映射器
 * @param
 * @return
 */
func (mg *MapperGenerator) CreateEntityMapper(tableObject *TableObject) error {
	tpl := mg.TemplateFileHeader(mg.mapperPackage)
	tpl += mg.TemplateMapperFile()

	structName := mg.FormatColumnAndTagName(tableObject.Name, false)
	structInstance := mg.FormatColumnAndTagName(tableObject.Name, true)
	packageName := mg.FormatColumnAndTagName(mg.mapperPackage, true)
	tpl = strings.Replace(tpl, "{{StructName}}", structName+packageName, -1)
	tpl = strings.Replace(tpl, "{{StructInstance}}", structInstance+packageName, -1)

	path := mg.SavePath(mg.mapperPath, structName+packageName)
	return mg.SaveFile(path, tpl)
}

/**
 * 保存文件
 * @param
 * @return
 */
func (mg *MapperGenerator) SaveFile(path string, content string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("create entity file error: %v", err.Error())
		return err
	}

	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

/**
 * 保存路径
 * @param
 * @return
 */
func (mg *MapperGenerator) SavePath(path string, name string) string {
	mg.CheckPath(path)
	return path + "/" + name + ".go"
}

/**
 * 检查保存目录
 * @param
 * @return
 */
func (mg *MapperGenerator) CheckPath(path string) {
	s, err := os.Stat(path)
	if err != nil || !s.IsDir() {
		err = os.MkdirAll(path, os.ModePerm)
		log.Printf("create path %v error: %v", path, err)
	}
}

/**
 * 文件头部
 * @param
 * @return
 */
func (mg *MapperGenerator) TemplateFileHeader(packageName string) string {
	comment := "package " + packageName + "\r\n"
	return comment
}

/**
 * 文件模版
 * @param
 * @return
 */
func (mg *MapperGenerator) TemplateEntityFile() string {
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
 * 文件模版
 * @param
 * @return
 */
func (mg *MapperGenerator) TemplateMapperFile() string {
	tpl := "\r\nimport gormmapper \"github.com/Rgss/gorm-mapper\"\r\n\r\n"

	if mg.mapperPathAutoSignleton {
		tpl += "// global instance\r\n"
		tpl += "var {{StructInstance}} = &{{StructName}}{}\r\n\r\n"
	}

	tpl += "// the mapper is generated by gorm-mapper automatically\r\n"
	tpl += "// {{StructName}} \r\n"
	tpl += "type {{StructName}} struct {\r\n"
	tpl += "    gormmapper.Mapper\r\n"
	tpl += "}\r\n"

	return tpl
}

/**
 * 字段模版
 * @param
 * @return
 */
func (mg *MapperGenerator) TemplateColumn(tableObject *TableObject) string {
	columns := tableObject.Columns
	tpl := ""
	for _, v := range columns {
		if v == nil || len(v.Name) <= 0 {
			continue
		}

		c := mg.FormatColumnAndTagName(v.Name, true)
		vt := mg.ConvertType(v.Type)
		ct := mg.ColumnStructTag(v)

		i := tableObject.MaxColumnLength - len(c)
		csp := mg.columnSpacePad(i)
		i2 := tableObject.MaxColumnTypeLength - len(vt)
		csp2 := mg.columnSpacePad(i2)
		tpl += fmt.Sprintf("    %s %s%s%s %s\r\n", c, csp, vt, csp2, ct)
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
 * 空白填充
 * @param
 * @return
 */
func (mg *MapperGenerator) columnSpacePad(count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(" ", count)
}

/**
 * gorm tag
 * @param
 * @return
 */
func (mg *MapperGenerator) gormColumn(tableColumn *TableColumn) string {
	tpl := "`gorm:\""
	tpl += "column:" + tableColumn.Name + ";"

	// PrimaryKey
	if tableColumn.PrimaryKey {
		tpl += " primaryKey;"
	}

	// not null
	if strings.Contains(tableColumn.PlainSQL, "NOT NULL") {
		tpl += " not null;"
	}

	if strings.Contains(tableColumn.PlainSQL, "AUTO_INCREMENT") {
		tpl += " autoIncrement;"
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

		tpl += " default:" + v + ";"
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
