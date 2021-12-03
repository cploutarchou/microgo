package microGo

import (
	"crypto/rand"
	"os"
)

const (
	rndString = "sd54qw54dfg2365swe$445dfg765SDFDff"
)

// CreateRandomString  A Random String Generator function based on n value length.
// From the values in the rndString const
func (m *MicroGo) CreateRandomString(n int) string {
	s, r := make([]rune, n), []rune(rndString)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

// CreateDirIfNotExist  creates the necessary folder if not exist.
func (m *MicroGo) CreateDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateFileIfNotExists  creates the necessary files if not exist.
func (m *MicroGo) CreateFileIfNotExists(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}
