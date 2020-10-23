package gormmapper

// config
type DBConfig struct {
	User    string
	Pass    string
	Host    string
	Port    int
	DbName  string
	Charset string

	MaxIdleConns int
	MaxOpenConns int
	EnableLog    bool
}