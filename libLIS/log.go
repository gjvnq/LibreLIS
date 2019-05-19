package libLIS

import (
	"errors"
	"os"

	"github.com/gjvnq/go-logger"
)

var TheLogger *logger.Logger
var ErrDB error = errors.New("falha ao se comunicar com o banco de dados")

func logAndPanicIfNotNil(err error) {
	if err != nil {
		TheLogger.ErrorNF(1, "%s", err.Error())
		panic(err)
	}
}

func startLogger() {
	var err error
	TheLogger, err = logger.New("test", 1, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func main() {
	startLogger()
	err := errors.New("hi")
	logAndPanicIfNotNil(err)
}
