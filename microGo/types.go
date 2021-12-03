package microGo

import "database/sql"

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
	Pool         *sql.DB
}
