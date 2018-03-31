package zbcommon

import (
	"fmt"
	"io/ioutil"
	"time"

	"os"

	diodes "code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"strings"
)

// ZBL is pointer to client logger
var ZBL *zerolog.Logger

var zbLogPrefix = "zbc"
var zbLogpath *string

func init() {
	if ZBL == nil {
		InitLogger()
		SetLogLevel()
	}
}

// DestroyLogger will destroy current package logger and reset log prefix to 'zbc'.
func DestroyLogger() {
	ZBL = nil
	zbLogPrefix = "zbc"
}

// SetLogOutput will determine where to flush the logs.
func SetLogOutput(t string) {
	os.Setenv("ZBC_LOG", t)
}

// SetLogPrefix will prefix the current snapshot log with the given name.
func SetLogPrefix(prefix string) {
	zbLogPrefix = prefix
}

// SetLogLevel will determine how much logging we want.
func SetLogLevel() {
	zerolog.SetGlobalLevel(getLogLevelFromEnv())
}

func getLogLevelFromEnv() zerolog.Level {
	logLevel := strings.ToLower(os.Getenv("ZBC_LOG_LEVEL"))
	switch logLevel {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		// default to info
		return zerolog.InfoLevel
	}
}

// InitLogger will create package logger based on ZBC_LOG env variable.
func InitLogger() {
	zerolog.TimeFieldFormat = time.StampNano
	var ll zerolog.Logger

	d := diodes.NewManyToOne(100000, diodes.AlertFunc(func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	}))

	switch os.Getenv("ZBC_LOG") {

	case "stdout": // export ZBC_LOG=stdout
		ll = zerolog.New(diode.NewWriter(zerolog.ConsoleWriter{Out: os.Stdout}, d, 10*time.Millisecond)).With().Timestamp().Logger()
		//ll = ll.With().Str("component", string(GetGID())).Logger()
		break

	case "snapshot": // export ZBC_LOG=snapshot://$GOPATH/src/github.com/zbc-go/.snapshots
		projectPath := fmt.Sprintf("%s/src/github.com/zeebe-io/zbc-go/.snapshots", os.Getenv("GOPATH"))
		os.MkdirAll(projectPath, os.FileMode(0777))
		filename := fmt.Sprintf("%s-%s.log", zbLogPrefix, time.Now().Format("20060102150405"))
		logpath := fmt.Sprintf("%s/%s", projectPath, filename)

		f, err := os.Create(logpath)
		if err != nil {
			panic(err)
		}
		zbLogpath = &logpath
		ll = zerolog.New(diode.NewWriter(f, d, 10*time.Millisecond)).With().Timestamp().Logger()

		break

	case "file": // export ZBC_LOG=file://tmp/zbc/zbc-%s.log
		filename := fmt.Sprintf("zbc-%s.log", time.Now().Format("20060102150405"))
		logpath := fmt.Sprintf("/tmp/zbc/%s", filename)

		f, err := os.Create(logpath)
		if err != nil {
			panic("cannot create log file, unset ZBC_LOG or make sure user running this process can create log file")
		}
		zbLogpath = &logpath

		ll = zerolog.New(diode.NewWriter(f, d, 10*time.Millisecond))
		break

	default:
		ll = zerolog.New(diode.NewWriter(ioutil.Discard, d, 10*time.Millisecond))
		break

	}

	ZBL = &ll
}
