package rapiditas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cploutarchou/rapiditas/render"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

// Rapiditas is the overall type for the Rapiditas package. Members that are exported in this type
// are available to any application that uses it.
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
	Render     *render.Render
	config     config
}

type config struct {
	port     string
	renderer string
}

// New reads the .env file, creates our application config, populates the Rapiditas type with settings
// based on .env values, and creates the necessary folders and files if they don't exist on the system.
func (r *Rapiditas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := r.Init(pathConfig)
	if err != nil {
		return err
	}

	err = r.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// Read values from  .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// initiate the  loggers
	infoLog, errorLog, warnLog, buildLog := r.startLoggers()
	r.InfoLog = infoLog
	r.ErrorLog = errorLog
	r.WarningLog = warnLog
	r.BuildLog = buildLog
	r.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	r.Version = version
	r.RootPath = rootPath
	r.Routes = r.routes().(*chi.Mux)

	r.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	r.createRenderer()

	return nil
}

// Init creates the necessary folders for Rapiditas application
func (r *Rapiditas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it doesn't exist
		err := r.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListenAndServe starts the application web server
func (r *Rapiditas) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     r.ErrorLog,
		Handler:      r.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	r.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	r.ErrorLog.Fatal(err)
}

func (r *Rapiditas) checkDotEnv(path string) error {
	err := r.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

//startLoggers Initializes all loggers for rapiditas application.
func (r *Rapiditas) startLoggers() (*log.Logger, *log.Logger, *log.Logger, *log.Logger) {
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

//createRenderer Create a Renderer for rapiditas application.
func (r *Rapiditas) createRenderer() {
	renderer := render.Render{
		Renderer: r.config.renderer,
		RootPath: r.RootPath,
		Port:     r.config.port,
	}
	r.Render = &renderer
}
