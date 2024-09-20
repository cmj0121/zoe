package monitor

import (
	"context"
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zoe/pkg/monitor/routes"
)

// The HTTP server that show the records of the honeypot.
type Server struct {
	Bind *string `name:"bind" help:"The address to bind the HTTP server"`

	*gin.Engine `kong:"-"`
}

func (s *Server) Run(ctx context.Context) {
	if s.Bind == nil {
		log.Info().Msg("no bind address, skip the HTTP server")
		return
	}

	srv := s.serve()
	// gracefun shutdown the server
	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown the HTTP server")
		return
	}

	log.Info().Msg("successfully shutdown the HTTP server")
}

func (s *Server) serve() *http.Server {
	gin.SetMode(gin.ReleaseMode)

	s.Engine = gin.New()
	// setup the middleware
	s.Engine.Use(gin.Recovery())
	s.Engine.Use(logger.SetLogger())
	s.register()

	srv := &http.Server{
		Addr:    *s.Bind,
		Handler: s.Engine,
	}

	go func() {
		switch err := srv.ListenAndServe(); err {
		case nil:
			log.Info().Msg("start the HTTP server")
		case http.ErrServerClosed:
			log.Info().Msg("shutdown the HTTP server")
		default:
			log.Fatal().Err(err).Msg("failed to start the HTTP server")
		}
	}()

	log.Info().Str("bind", *s.Bind).Msg("starting the HTTP server")
	return srv
}

// register the routes of the HTTP server.
func (s *Server) register() {
	s.Engine.GET("/", routes.APIIndex)
	s.Engine.GET("/messages/daily-popular", routes.MessagePopular)
	s.Engine.GET("/messages/daily-popular/:field", routes.APIMessagePopular)
}
