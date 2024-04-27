package monitor

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/url"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

//go:embed web/template/*.htm
var templates embed.FS

// The type instance of Monitor that show the records from the honeypot
type Monitor struct {
	// The database connection
	*sql.DB `-`

	Bind     string // The bind address of the monitor
	Database string // The database URI
}

// Serve the monitor service
func (m *Monitor) Serve() error {
	err := m.prologue()
	if err != nil {
		log.Error().Err(err).Msg("failed to setup the monitor service")
		return err
	}
	defer m.epilogue()

	return m.serve()
}

// setup the HTTP server and run the monitor service
func (m *Monitor) serve() error {
	routes := gin.New()
	// setup the middleware
	routes.Use(gin.Recovery())
	routes.Use(logger.SetLogger())
	// register the routes
	m.register(routes)

	// setup HTML templates
	switch tmpl, err := template.ParseFS(templates, "web/template/*.htm"); err {
	case nil:
		routes.SetHTMLTemplate(tmpl)
	default:
		log.Warn().Err(err).Msg("failed to parse the HTML templates")
		return err
	}

	// run the monitor service
	return routes.Run(m.Bind)
}

// setup everything before running the monitor service
func (m *Monitor) prologue() error {
	// parse the database URI
	uri, err := url.Parse(m.Database)
	if err != nil {
		log.Warn().Err(err).Str("database", m.Database).Msg("failed to parse the database URI")
		return err
	}

	switch uri.Scheme {
	case "sqlite", "sqlite3":
		path := fmt.Sprintf("%s%s", uri.Host, uri.Path)
		// open the database connection
		db, err := sql.Open("sqlite3", path)
		if err != nil {
			log.Warn().Err(err).Str("database", m.Database).Msg("failed to open the database")
			return err
		}

		m.DB = db
	default:
		err := errors.New("unsupported database")
		log.Warn().Err(err).Str("database", m.Database).Msg("unsupported database")
		return err
	}

	// set the gin to release mode
	gin.SetMode(gin.ReleaseMode)

	return nil
}

// close all the allocated resources after running the monitor service
func (m *Monitor) epilogue() {
	m.DB.Close()
}

// register the routes for the monitor service
func (m *Monitor) register(r *gin.Engine) {
	r.GET("/", m.index)
	r.GET("/view/group_by/:field", m.group_by)
	r.GET("/view/chart/:field", m.chart)
	r.GET("/static/:filepath", m.static)
	r.GET("/livez", m.livez)
	r.GET("/readyz", m.readyz)
}
