package endrLogger

import (
	"errors"
	"log"
	"os"
	"runtime/debug"
)

const (
	LevelNone = iota
	LevelInfo
	LevelError
	LevelCritical
	LevelDebug
)

type Config struct {
	//ErrorLogFile string `json:"error_log_file"`
	//DebugLogFile string `json:"debug_log_file"`
	//InfoLogFile  string `json:"info_log_file"`
	Dir   string
	Level int
}

type Logger struct {
	defaultLogger  *log.Logger
	errLogger      *log.Logger
	debugLogger    *log.Logger
	infoLogger     *log.Logger
	criticalLogger *log.Logger
	files          []*os.File
	config         *Config
}

func (logger *Logger) Println(msg ...interface{}) {
	logger.defaultLogger.Println(msg...)
	logger.Error(msg...)
}

func (logger *Logger) Info(msg ...interface{}) {
	if logger.infoLogger != nil {
		logger.infoLogger.Println(msg...)
	} else {
		logger.Error(msg...)
	}
	logger.checkDebug(msg...)
}

func (logger *Logger) Error(msg ...interface{}) {
	if logger.errLogger != nil {
		logger.errLogger.Println(msg...)
	} else {
		logger.Critical(msg...)
	}
	logger.checkDebug(msg...)
}

func (logger *Logger) Critical(msg ...interface{}) {
	if logger.criticalLogger != nil {
		logger.criticalLogger.Println(msg...)
	}
	logger.checkDebug(msg...)
	log.Fatal(msg)
}

func (logger *Logger) Debug(msg ...interface{}) {
	if logger.debugLogger != nil {
		logger.debugLogger.Println(msg...)
		logger.debugLogger.Println(string(debug.Stack()))
	}
}

func (logger *Logger) checkDebug(msg ...interface{}) {
	if logger.debugLogger != nil {
		logger.Debug(msg...)
	}
}

func (logger *Logger) Close() {
	for _, file := range logger.files {
		file.Close()
	}
}

func (logger *Logger) GetConfig() *Config {
	return logger.config
}

func NewLogger(config Config) (logger *Logger, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			logger = nil
			err = errors.New(rec.(string))
		}
	}()
	if config.Dir == "" {
		config.Dir = "./_log/"
	}
	if _, err := os.Stat(config.Dir); os.IsNotExist(err) {
		os.MkdirAll(config.Dir, os.ModePerm)
	}
	if config.Dir[len(config.Dir)-1] != '/' {
		config.Dir += "/"
	}
	defaultLogger := log.New(log.Writer(), "", log.LstdFlags)
	var files []*os.File
	var debugLogger *log.Logger
	if config.Level >= LevelDebug {
		l, debugFile := createLogger(config.Dir+"debug.log", "DEBUG: ")
		if debugFile != nil {
			files = append(files, debugFile)
			debugLogger = l
		}
	}
	var criticalLogger *log.Logger
	if config.Level >= LevelCritical {
		l, crFile := createLogger(config.Dir+"critical.log", "ERROR: ")
		if crFile != nil {
			files = append(files, crFile)
			criticalLogger = l
		}
	}
	var errLogger *log.Logger
	if config.Level >= LevelError {
		l, errFile := createLogger(config.Dir+"error.log", "ERROR: ")
		if errFile != nil {
			files = append(files, errFile)
			errLogger = l
		}
	}
	var infoLogger *log.Logger
	if config.Level >= LevelInfo {
		l, infoFile := createLogger(config.Dir+"info.log", "INFO: ")
		if infoFile != nil {
			files = append(files, infoFile)
			infoLogger = l
		}
	}
	logger = &Logger{
		defaultLogger:  defaultLogger,
		errLogger:      errLogger,
		debugLogger:    debugLogger,
		infoLogger:     infoLogger,
		criticalLogger: criticalLogger,
		files:          files,
		config:         &config,
	}
	return
}

func createLogger(filePath string, prefix string) (*log.Logger, *os.File) {
	if filePath == "" {
		return log.New(log.Writer(), prefix, log.LstdFlags), nil
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if os.IsNotExist(err) {
		file, err = os.Create(filePath)
	}
	if err != nil {
		panic(errors.New("log file error: " + err.Error()))
	}
	return log.New(file, prefix, log.LstdFlags), file
}
