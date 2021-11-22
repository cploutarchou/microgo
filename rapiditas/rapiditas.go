package rapiditas

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	version = "1.0.0"
)

type Rapiditas struct {
	AppName    string
	Debug      bool
	Version    string
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	WarningLog *log.Logger
	BuildLog   *log.Logger
	RootPath   string
	Routes     *chi.Mux
	config     config
}

type config struct {
	port     string
	rendered string
}

func (r *Rapiditas) New(rootPath string) error {
	pathConf := initPaths{
		rootPath:     rootPath,
		foldersNames: []string{"handles", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}
	err := r.Init(pathConf)
	if err != nil {
		return err
	}

	err = r.checkENV(rootPath)
	if err != nil {
		return err
	}
	// Read environment file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}
	// implement Logger
	infoLog, errorLog, warnLog, buildLog := r.startLogger()
	r.InfoLog = infoLog
	r.ErrorLog = errorLog
	r.WarningLog = warnLog
	r.BuildLog = buildLog
	r.RootPath = rootPath
	r.Routes = r.routes().(*chi.Mux)

	r.config = config{
		port:     os.Getenv("PORT"),
		rendered: os.Getenv("RENDERER"),
	}

	r.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	r.Version = version
	return nil
}

func (r *Rapiditas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.foldersNames {
		// create the directory if it doesn't exist'
		err := r.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListenEndServe  Starting the web server.
func (r *Rapiditas) ListenEndServe() {
	svr := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     r.ErrorLog,
		Handler:      r.routes(),
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}
	r.InfoLog.Printf("Listening on %s", os.Getenv("PORT"))
	err := svr.ListenAndServe()
	r.ErrorLog.Fatal(err)

}
func (r *Rapiditas) checkENV(path string) error {
	err := r.CreateFileIfNotExists(fmt.Sprintf("%s.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (r *Rapiditas) startLogger() (*log.Logger, *log.Logger, *log.Logger, *log.Logger) {
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
