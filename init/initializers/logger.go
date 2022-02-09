package initializers

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// LoggerConfiguration represents configuration for logs.
type LoggerConfiguration struct {
	// Level of logging, can be DebugLevel, InfoLevel, WarnLevel, ErrorLevel, DPanicLevel, PanicLevel
	Level string `yaml:"level"`
}

var loggerConfig LoggerConfiguration

var logg *zap.Logger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func createDirectoryIfDoesntExist() {
	crrPath, _ := crrFSGetter.getwd()
	if _, err := os.Stat(crrPath + "/logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}
}

func getLogWriter() zapcore.WriteSyncer {
	path, err :=  os.Getwd()
	if err != nil {
		panic(err)
	}
	timeString := time.Now().Format("02-01-2006")
	file, err := os.OpenFile(path + "/logs/" + timeString +".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return zapcore.AddSync(file)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.UTC().Format("02-01-2006T15:04:05Z0700"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLevel() zapcore.Level {
	err := LoadConfigSection("logger", &loggerConfig)
	if err != nil {
		panic(err)
	}
	if level, err := levelMap[loggerConfig.Level]; err {
		return level
	}
	return zapcore.DebugLevel
}

func ChiLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logg.Info("Request:", zap.String("path", r.URL.Path), zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()), zap.String("requestId", middleware.GetReqID(r.Context())))
		}
		return http.HandlerFunc(fn)
	}
}

func LoggerInitializer() {
	createDirectoryIfDoesntExist()
	writerSync := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writerSync, getLevel()),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevel()),
	)
	logg = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)
}