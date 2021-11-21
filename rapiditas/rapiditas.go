package rapiditas

import "fmt"

const (
	version = "1.0.0"
)

type Rapiditas struct {
	AppName string
	Debug   bool
	Version string
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
	return nil
}

func (r *Rapiditas) Init(p initPaths) error {
	fmt.Print(version)
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
