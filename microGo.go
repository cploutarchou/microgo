package MicroGO

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/kataras/blocks"
	"github.com/robfig/cron/v3"

	"github.com/cploutarchou/MicroGO/cache"
	"github.com/cploutarchou/MicroGO/mailer"
	"github.com/cploutarchou/MicroGO/render"
	"github.com/cploutarchou/MicroGO/requests"
	"github.com/cploutarchou/MicroGO/session"
)

const version = "1.0.7a"

var (
	redisCache       *cache.RedisCache
	badgerCache      *cache.BadgerCache
	redisPool        *redis.Pool
	badgerConnection *badger.DB
)

// MicroGo is the overall type for the MicroGo package. Members that are exported in this type
// are available to any application that uses it.
type MicroGo struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	WarningLog    *log.Logger
	BuildLog      *log.Logger
	RootPath      string
	Routes        *chi.Mux
	Render        *render.Render
	JetView       *jet.Set
	BlocksView    *blocks.Blocks
	config        config
	Session       *scs.SessionManager
	DB            Database
	EncryptionKey string
	Cache         cache.Cache
	Scheduler     *cron.Cron
	Mailer        mailer.Mailer
	Server        Server
	Requests      *requests.Requests
}

type Server struct {
	ServerName string
	Port       string
	Secure     bool
	URL        string
}
type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	redis       redisConfig
}

// New reads the .env file, creates our application config, populates the MicroGo type with settings
// based on .env values, and creates the necessary folders and files if they don't exist on the system.
func (m *MicroGo) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "mail", "data", "public", "tmp", "logs", "middleware"},
	}

	err := m.Init(pathConfig)
	if err != nil {
		return err
	}

	err = m.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// Read values from  .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// initiate the  loggers
	infoLog, errorLog, warnLog, buildLog := m.startLoggers()
	m.InfoLog = infoLog
	m.ErrorLog = errorLog
	m.WarningLog = warnLog
	m.BuildLog = buildLog

	// Initiate database connection
	if os.Getenv("DATABASE_TYPE") != "" {
		var db *sql.DB
		switch os.Getenv("DATABASE_TYPE") {
		case "":
			m.ErrorLog.Println("DATABASE_TYPE is not set")

		case "mysql", "mariadb":
			db, err = m.OpenDB("mysql", m.BuildDataSourceName())
			if err != nil {
				errorLog.Println(err)
				os.Exit(1)
			}
		case "postgres", "postgresql":
			db, err = m.OpenDB("postgres", m.BuildDataSourceName())
			if err != nil {
				errorLog.Println(err)
				os.Exit(1)
			}

		}
		m.DB = Database{
			DatabaseType: os.Getenv("DATABASE_TYPE"),
			Pool:         db,
		}
	}
	scheduler := cron.New()
	m.Scheduler = scheduler

	if os.Getenv("CACHE") == "redis" || os.Getenv("SESSION_TYPE") == "redis" {
		redisCache = m.createRedisCacheClient()
		m.Cache = redisCache
		redisPool = redisCache.Connection
	}

	if os.Getenv("CACHE") == "badger" {
		badgerCache = m.createBadgerCacheClient()
		m.Cache = badgerCache
		badgerConnection = badgerCache.Connection

		_, err = m.Scheduler.AddFunc("@daily", func() {
			_ = badgerCache.Connection.RunValueLogGC(0.7)
		})
		if err != nil {
			return err
		}
	}
	m.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	m.Version = version
	m.RootPath = rootPath
	m.Routes = m.routes().(*chi.Mux)
	// initiate mailer
	m.Mailer = m.createMailer()
	m.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			database:       os.Getenv("DATABASE_TYPE"),
			dataSourceName: m.BuildDataSourceName(),
		},
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			port:     os.Getenv("REDIS_PORT"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}

	secure := true
	if strings.ToLower(os.Getenv("SECURE")) == "false" {
		secure = false
	}
	m.Server = Server{
		ServerName: os.Getenv("SERVER_NAME"),
		Port:       os.Getenv("PORT"),
		Secure:     secure,
		URL:        os.Getenv("APP_URL"),
	}
	// initiate session
	_session := session.Session{
		CookieLifetime: m.config.cookie.lifetime,
		CookiePersist:  m.config.cookie.persist,
		CookieName:     m.config.cookie.name,
		SessionType:    m.config.sessionType,
		CookieDomain:   m.config.cookie.domain,
		DBPool:         m.DB.Pool,
	}
	switch m.config.sessionType {
	case "redis":
		_session.RedisPool = redisCache.Connection
	case "mysql", "mariadb", "postgres", "postgresql":
		_session.DBPool = m.DB.Pool
	}

	m.Session = _session.InitializeSession()
	m.EncryptionKey = os.Getenv("ENCRYPTION_KEY")
	if m.Debug {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
			jet.InDevelopmentMode(),
		)
		m.JetView = views
	} else {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		)

		m.JetView = views
	}

	m.createRenderer()
	go m.Mailer.ListenForMessage()
	return nil
}

// Init creates the necessary folders for MicroGo application
func (m *MicroGo) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it doesn't exist
		err := m.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListenAndServe starts the application web server
func (m *MicroGo) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     m.ErrorLog,
		Handler:      m.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}
	if m.DB.Pool != nil {
		defer func(Pool *sql.DB) {
			err := Pool.Close()
			if err != nil {
				m.WarningLog.Println(err)
			}
		}(m.DB.Pool)
	}

	if redisPool != nil {
		defer func(redisPool *redis.Pool) {
			err := redisPool.Close()
			if err != nil {
				m.WarningLog.Println(err)
			}
		}(redisPool)
	}
	if badgerConnection != nil {
		defer func(badgerConnection *badger.DB) {
			err := badgerConnection.Close()
			if err != nil {
				m.WarningLog.Println(err)
			}
		}(badgerConnection)
	}
	m.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	m.ErrorLog.Fatal(err)
}

func (m *MicroGo) checkDotEnv(path string) error {
	err := m.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

// startLoggers Initializes all loggers for microGo application.
func (m *MicroGo) startLoggers() (*log.Logger, *log.Logger, *log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger
	var warnLog *log.Logger
	var buildLog *log.Logger
	warnLog = log.New(os.Stderr, "[ WARNING ] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog = log.New(os.Stderr, "[ INFO ] ", log.Ldate|log.Ltime|log.Lshortfile)
	buildLog = log.New(os.Stderr, "[ BUILD ] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(os.Stderr, "[ ERROR ] ", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errorLog, warnLog, buildLog
}

// createRenderer Create a Renderer for microGo application.
func (m *MicroGo) createRenderer() {
	renderer := render.Render{
		Renderer:    m.config.renderer,
		RootPath:    m.RootPath,
		Port:        m.config.port,
		JetViews:    m.JetView,
		BlocksViews: m.BlocksView,
		Session:     m.Session,
	}
	m.Render = &renderer
}

// createRenderer Create a Renderer for microGo application.
func (m *MicroGo) createMailer() mailer.Mailer {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	_mailer := mailer.Mailer{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Templates:   m.RootPath + "/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		FromName:    os.Getenv("FROM_NAME"),
		Jobs:        make(chan mailer.Message, 20),
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		ApiKey:      os.Getenv("MAILER_KEY"),
		ApiUrl:      os.Getenv("MAILER_URL"),
	}
	return _mailer
}

// BuildDataSourceName builds the datasource name for our database, and returns it as a string
func (m *MicroGo) BuildDataSourceName() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
		return dsn
	case "mysql", "mariadb":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASS"),
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_NAME"))
	default:

	}
	return ""
}

func (m *MicroGo) createRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", m.config.redis.host, m.config.redis.port),
				redis.DialPassword(m.config.redis.password),
				redis.DialUsername(m.config.redis.username),
			)
		},
		DialContext: nil,
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err

		},
		MaxIdle:         50,
		MaxActive:       10000,
		IdleTimeout:     240 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}

func (m *MicroGo) createRedisCacheClient() *cache.RedisCache {
	_client := cache.RedisCache{
		Connection: m.createRedisPool(),
		Prefix:     m.config.redis.prefix,
	}
	return &_client
}
func (m *MicroGo) createBadgerCacheClient() *cache.BadgerCache {
	cacheClient := cache.BadgerCache{
		Connection: m.connectToBadgerCache(),
	}
	return &cacheClient
}
func (m *MicroGo) connectToBadgerCache() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(m.RootPath + "/tmp/badger"))
	if err != nil {
		return nil
	}
	return db
}
