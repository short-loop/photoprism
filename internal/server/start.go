package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/config"

	"github.com/short-loop/shortloop-go/shortloopgin"
)

// Start the REST API server using the configuration provided
func Start(ctx context.Context, conf *config.Config) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	start := time.Now()

	// Set HTTP server mode.
	if conf.HttpMode() != "" {
		gin.SetMode(conf.HttpMode())
	} else if conf.Debug() == false {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create new HTTP router engine without standard middleware.
	router := gin.New()
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           os.Getenv("SENTRY_DSN"),
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		Debug:            true,
		AttachStacktrace: true,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	shortloopSdk, err := shortloopgin.Init(shortloopgin.Options{
		ShortloopEndpoint: os.Getenv("ShortloopEndpoint"),
		ApplicationName:   os.Getenv("ApplicationName"),
		LoggingEnabled:    true,
		LogLevel:          "INFO",
	})
	if err != nil {
		fmt.Println("Error initializing shortloopgin: ", err)
	} else {
		router.Use(shortloopSdk.Filter())
	}

	// Set proxy addresses from which headers related to the client and protocol can be trusted
	if err = router.SetTrustedProxies(conf.TrustedProxies()); err != nil {
		log.Warnf("server: %s", err)
	}

	router.GET(conf.BaseUri(config.ApiUri)+"/sentry-test", func(c *gin.Context) {
		panic("test panic for sentry")
	})

	// Register common middleware.
	router.Use(Recovery(), Security(conf), Logger())

	router.GET(conf.BaseUri(config.ApiUri)+"/panic1", func(c *gin.Context) {
		panic("test panic")
	})

	router.GET(conf.BaseUri(config.ApiUri)+"/panic2", func(c *gin.Context) {
		var p *int = nil
		fmt.Println(*p)
	})

	// Initialize package extensions.
	Ext().Init(router, conf)

	// Enable HTTP compression?
	switch conf.HttpCompression() {
	case "gzip":
		log.Infof("server: enabling gzip compression")
		router.Use(gzip.Gzip(
			gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{
				conf.BaseUri(config.ApiUri + "/t"),
				conf.BaseUri(config.ApiUri + "/folders/t"),
				conf.BaseUri(config.ApiUri + "/zip"),
				conf.BaseUri(config.ApiUri + "/albums"),
				conf.BaseUri(config.ApiUri + "/labels"),
				conf.BaseUri(config.ApiUri + "/videos"),
			})))
	}

	// Find and load templates.
	router.LoadHTMLFiles(conf.TemplateFiles()...)

	// Register HTTP route handlers.
	registerRoutes(router, conf)

	var tlsErr error
	var tlsManager *autocert.Manager
	var server *http.Server

	// Enable TLS?
	if tlsManager, tlsErr = AutoTLS(conf); tlsErr == nil {
		server = &http.Server{
			Addr:      fmt.Sprintf("%s:%d", conf.HttpHost(), conf.HttpPort()),
			TLSConfig: tlsManager.TLSConfig(),
			Handler:   router,
		}
		log.Infof("server: starting in auto tls mode on %s [%s]", server.Addr, time.Since(start))
		go StartAutoTLS(server, tlsManager, conf)
	} else if publicCert, privateKey := conf.TLS(); publicCert != "" && privateKey != "" {
		log.Infof("server: starting in tls mode")
		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", conf.HttpHost(), conf.HttpPort()),
			Handler: router,
		}
		log.Infof("server: listening on %s [%s]", server.Addr, time.Since(start))
		go StartTLS(server, publicCert, privateKey)
	} else {
		log.Infof("server: %s", tlsErr)
		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", conf.HttpHost(), conf.HttpPort()),
			Handler: router,
		}
		log.Infof("server: listening on %s [%s]", server.Addr, time.Since(start))
		go StartHttp(server)
	}

	// Graceful HTTP server shutdown.
	<-ctx.Done()
	log.Info("server: shutting down")
	err = server.Close()
	if err != nil {
		log.Errorf("server: shutdown failed (%s)", err)
	}
}

// StartHttp starts the web server in http mode.
func StartHttp(s *http.Server) {
	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Info("server: shutdown complete")
		} else {
			log.Errorf("server: %s", err)
		}
	}
}

// StartTLS starts the web server in https mode.
func StartTLS(s *http.Server, httpsCert, privateKey string) {
	if err := s.ListenAndServeTLS(httpsCert, privateKey); err != nil {
		if err == http.ErrServerClosed {
			log.Info("server: shutdown complete")
		} else {
			log.Errorf("server: %s", err)
		}
	}
}

// StartAutoTLS starts the web server with auto tls enabled.
func StartAutoTLS(s *http.Server, m *autocert.Manager, conf *config.Config) {
	var g errgroup.Group

	g.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf("%s:%d", conf.HttpHost(), conf.HttpPort()), m.HTTPHandler(http.HandlerFunc(redirect)))
	})

	g.Go(func() error {
		return s.ListenAndServeTLS("", "")
	})

	if err := g.Wait(); err != nil {
		if err == http.ErrServerClosed {
			log.Info("server: shutdown complete")
		} else {
			log.Errorf("server: %s", err)
		}
	}
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.RequestURI

	http.Redirect(w, req, target, httpsRedirect)
}
