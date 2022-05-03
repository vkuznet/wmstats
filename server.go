package main

// server module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dmwm/cmsauth"
	"github.com/gorilla/mux"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logging "github.com/vkuznet/http-logging"
)

// global variables
var _top, _bottom, _search string

// GitVersion defines git version of the server
var GitVersion string

// ServerInfo defines server info
var ServerInfo string

// StartTime represents initial time when we started the server
var StartTime time.Time

// CMSAuth structure to create CMS Auth headers
var CMSAuth cmsauth.CMSAuth

func basePath(api string) string {
	base := Config.Base
	if base != "" {
		if strings.HasPrefix(api, "/") {
			api = strings.Replace(api, "/", "", 1)
		}
		if strings.HasPrefix(base, "/") {
			return fmt.Sprintf("%s/%s", base, api)
		}
		return fmt.Sprintf("/%s/%s", base, api)
	}
	return api
}

func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true) // to allow /route and /route/ end-points

	// aux APIs used by server
	router.HandleFunc(basePath("/healthz"), StatusHandler).Methods("GET")
	router.HandleFunc(basePath("/metrics"), MetricsHandler).Methods("GET")

	// main page
	router.HandleFunc(basePath("/"), MainHandler).Methods("GET")

	// for all requests
	router.Use(logging.LoggingMiddleware)
	// for all requests perform first auth/authz action
	router.Use(authMiddleware)

	// use limiter middleware to slow down clients
	router.Use(limitMiddleware)
	return router
}

// Server represents main web server for service
//gocyclo:ignore
func Server(configFile string) {
	StartTime = time.Now()
	err := ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(0)
	if Config.Verbose > 0 {
		log.SetFlags(log.Lshortfile)
	}
	log.SetOutput(new(logging.LogWriter))
	if Config.LogFile != "" {
		logName := Config.LogFile
		hostname := os.Getenv("HOSTNAME")
		if hostname == "" {
			hostname, err = os.Hostname()
			if err != nil {
				hostname = "localhost"
			}
		}
		if strings.HasSuffix(logName, ".log") {
			logName = fmt.Sprintf("%s-%s.log", strings.Split(logName, ".log")[0], hostname)
		} else {
			// it is log dir
			logName = fmt.Sprintf("%s/%s.log", logName, hostname)
		}
		logName = strings.Replace(logName, "//", "/", -1)
		//         rl, err := rotatelogs.New(Config.LogFile + "-%Y%m%d")
		rl, err := rotatelogs.New(logName + "-%Y%m%d")
		if err == nil {
			rotlogs := logging.RotateLogWriter{RotateLogs: rl}
			log.SetOutput(rotlogs)
		} else {
			log.Println("ERROR: unable to get rotatelogs", err)
		}
	}
	if err != nil {
		log.Printf("Unable to parse, time: %v, config: %v\n", time.Now(), configFile)
	}
	log.Println("Configuration:", Config.String())

	// initialize cmsauth layer
	CMSAuth.Init(Config.Hmac)

	// initialize limiter
	initLimiter(Config.LimiterPeriod)

	// initialize templates
	tmplData := make(map[string]interface{})
	tmplData["Time"] = time.Now()
	var templates Templates
	_top = templates.Tmpl(Config.Templates, "top.tmpl", tmplData)
	_bottom = templates.Tmpl(Config.Templates, "bottom.tmpl", tmplData)

	// static handlers
	for _, dir := range []string{"js", "css", "images"} {
		m := fmt.Sprintf("%s/%s/", Config.Base, dir)
		d := fmt.Sprintf("%s/%s", Config.StaticDir, dir)
		http.Handle(m, http.StripPrefix(m, http.FileServer(http.Dir(d))))
	}

	// setup WMStatsManager to handle our cache
	wmgr := NewWMStatsManager(Config.AccessURI)
	wmgr.update()

	// define our HTTP server
	addr := fmt.Sprintf(":%d", Config.Port)
	server := &http.Server{
		Addr: addr,
	}

	// make extra channel for graceful shutdown
	// https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	httpDone := make(chan os.Signal, 1)
	signal.Notify(httpDone, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start necessary HTTP servers

	go func() {
		// Start either HTTPs or HTTP web server
		_, e1 := os.Stat(Config.ServerCrt)
		_, e2 := os.Stat(Config.ServerKey)
		if e1 == nil && e2 == nil {
			//start HTTPS server which require user certificates
			rootCA := x509.NewCertPool()
			caCert, _ := ioutil.ReadFile(Config.RootCA)
			rootCA.AppendCertsFromPEM(caCert)
			server = &http.Server{
				Addr: addr,
				TLSConfig: &tls.Config{
					//                 ClientAuth: tls.RequestClientCert,
					RootCAs: rootCA,
				},
			}
			log.Printf("Starting HTTPs server at %v", addr)
			err = server.ListenAndServeTLS(Config.ServerCrt, Config.ServerKey)
		} else {
			// Start server without user certificates
			log.Printf("Starting HTTP server at %s", addr)
			err = server.ListenAndServe()
		}
		if err != nil {
			log.Printf("Fail to start server %v", err)
		}
	}()

	// properly stop our HTTP and Migration Servers
	<-httpDone
	log.Print("HTTP server stopped")

	// add extra timeout for shutdown service stuff
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("HTTP server exited properly")
}
