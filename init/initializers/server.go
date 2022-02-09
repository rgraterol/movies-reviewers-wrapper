package initializers

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var serverConfig ServerConfiguration

// ServerConfiguration represents a server configuration.
type ServerConfiguration struct {
	// Address is where the Server will listen
	Address string `yaml:"address"`
	// Timeout for all requests.
	Timeout int `yaml:"timeout"`
}

func ServerInitializer() {
	err := LoadConfigSection("server", &serverConfig)
	if err != nil {
		panic(errors.New("failed to read the server config"))
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(serverConfig.Timeout) * time.Second))
	r.Use(ChiLogger())

	zap.S().Info("Application running on address ", serverConfig.Address, " and enviroment ", Env())
	http.ListenAndServe(serverConfig.Address, r)
}
