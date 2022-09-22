package MicroGO

import "gorm.io/gorm"

type initPaths struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}

type databaseConfig struct {
	dataSourceName string
	database       string
}

type Database struct {
	DatabaseType string
	Client       *gorm.DB
}

type redisConfig struct {
	host     string
	port     string
	username string
	password string
	prefix   string
}
