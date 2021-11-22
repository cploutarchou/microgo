package microGo

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cploutarchou/microGo/render"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

// MicroGo is the overall type for the MicroGo package. Members that are exported in this type
// are available to any application that uses it.
type MicroGo struct {
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
	JetView    *jet.Set
	config     config
}

type config struct {
	port     string
	renderer string
}

// New reads the .env file, creates our application config, populates the MicroGo type with settings
// based on .env values, and creates the necessary folders and files if they don't exist on the system.
func (m *MicroGo) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
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
	m.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	m.Version = version
	m.RootPath = rootPath
	m.Routes = m.routes().(*chi.Mux)

	m.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	m.JetView = views
	m.createRenderer()

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

//startLoggers Initializes all loggers for microGo application.
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

//createRenderer Create a Renderer for microGo application.
func (m *MicroGo) createRenderer() {
	renderer := render.Render{
		Renderer: m.config.renderer,
		RootPath: m.RootPath,
		Port:     m.config.port,
		JetViews: m.JetView,
	}
	m.Render = &renderer
}
