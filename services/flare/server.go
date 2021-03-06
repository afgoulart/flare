// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"

	"github.com/diegobernardes/flare/document"
	infraHTTP "github.com/diegobernardes/flare/infra/http"
	infraMiddleware "github.com/diegobernardes/flare/infra/http/middleware"
	"github.com/diegobernardes/flare/resource"
	"github.com/diegobernardes/flare/subscription"
)

type server struct {
	addr       string
	httpServer http.Server
	handler    struct {
		resource     *resource.Service
		subscription *subscription.Service
		document     *document.Service
	}
	middleware struct {
		timeout time.Duration
	}
	logger        log.Logger
	writeResponse func(http.ResponseWriter, interface{}, int, http.Header)
}

func (s *server) start() error {
	router, err := s.router()
	if err != nil {
		return errors.Wrap(err, "error during router initialization")
	}

	s.httpServer = http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			fmt.Println(errors.Wrap(err, "error during server initialization").Error())

			process, err := os.FindProcess(os.Getpid())
			if err != nil {
				fmt.Println(errors.Wrap(err, "could not find flare process to exit").Error())
				os.Exit(1)
			}
			if err := process.Signal(os.Interrupt); err != nil {
				fmt.Println(errors.Wrap(err, "error during signal send to flare process").Error())
				os.Exit(1)
			}
		}
	}()

	return nil
}

func (s *server) stop() error {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		return errors.Wrap(err, "error during server close")
	}
	return nil
}

func (s *server) router() (http.Handler, error) {
	r := chi.NewRouter()
	if err := s.initMiddleware(r); err != nil {
		return nil, errors.Wrap(err, "error during middleware initialization")
	}

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		s.writeResponse(w, map[string]interface{}{
			"error": map[string]interface{}{
				"status": http.StatusBadRequest,
				"title":  "method not allowed",
			},
		}, http.StatusBadRequest, nil)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.writeResponse(w, map[string]interface{}{
			"error": map[string]interface{}{
				"status": http.StatusNotFound,
				"title":  "not found",
			},
		}, http.StatusNotFound, nil)
	})

	r.Route("/resources", s.routerResource)
	r.Route("/resources/{resourceId}/subscriptions", s.routerSubscription)
	r.Route("/documents", s.routerDocument)

	return r, nil
}

func (s *server) initMiddleware(r chi.Router) error {
	logger := infraMiddleware.NewLog(s.logger)
	writer, err := infraHTTP.NewWriter(s.logger)
	if err != nil {
		return errors.New("error during writer initialization")
	}

	recoverMiddleware, err := infraMiddleware.NewRecover(s.logger, writer)
	if err != nil {
		return errors.New("error during recover middleware initialization")
	}

	r.Use(recoverMiddleware.Handler)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(s.middleware.timeout))
	r.Use(logger.Handler)

	return nil
}

func (s *server) routerResource(r chi.Router) {
	r.Get("/", s.handler.resource.HandleIndex)
	r.Post("/", s.handler.resource.HandleCreate)
	r.Get("/{id}", s.handler.resource.HandleShow)
	r.Delete("/{id}", s.handler.resource.HandleDelete)
}

func (s *server) routerSubscription(r chi.Router) {
	r.Get("/", s.handler.subscription.HandleIndex)
	r.Post("/", s.handler.subscription.HandleCreate)
	r.Get("/{id}", s.handler.subscription.HandleShow)
	r.Delete("/{id}", s.handler.subscription.HandleDelete)
}

func (s *server) routerDocument(r chi.Router) {
	r.Get("/*", s.handler.document.HandleShow)
	r.Put("/*", s.handler.document.HandleUpdate)
	r.Delete("/*", s.handler.document.HandleDelete)
}

func newServer(options ...func(*server)) (*server, error) {
	s := &server{}

	for _, option := range options {
		option(s)
	}

	if s.addr == "" {
		s.addr = ":8080"
	}

	if s.handler.resource == nil {
		return nil, errors.New("missing handler.resource")
	}

	if s.handler.subscription == nil {
		return nil, errors.New("missing handler.subscription")
	}

	if s.handler.document == nil {
		return nil, errors.New("missing handler.document")
	}

	if s.logger == nil {
		return nil, errors.New("missing logger")
	}

	s.writeResponse = infraHTTP.WriteResponse(s.logger)
	return s, nil
}

func serverAddr(addr string) func(*server) {
	return func(s *server) { s.addr = addr }
}

func serverHandlerResource(handler *resource.Service) func(*server) {
	return func(s *server) { s.handler.resource = handler }
}

func serverHandlerSubscription(handler *subscription.Service) func(*server) {
	return func(s *server) { s.handler.subscription = handler }
}

func serverHandlerDocument(handler *document.Service) func(*server) {
	return func(s *server) { s.handler.document = handler }
}

func serverLogger(logger log.Logger) func(*server) {
	return func(s *server) { s.logger = logger }
}

func serverMiddlewareTimeout(duration time.Duration) func(*server) {
	return func(s *server) { s.middleware.timeout = duration }
}
