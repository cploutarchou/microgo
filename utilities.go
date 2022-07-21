package MicroGO

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

func (m *MicroGo) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	programCaller, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(programCaller)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	m.InfoLog.Println(fmt.Sprintf("Load Time: %s took %s", name, elapsed))
}
