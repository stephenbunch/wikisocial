package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/gocql/gocql"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/labstack/gommon/log"
	"github.com/rs/cors"
	"github.com/scylladb/gocqlx/migrate"
	"google.golang.org/grpc"

	pb "github.com/stephenbunch/wikitribe/server/_proto"
	"github.com/stephenbunch/wikitribe/server/app"
	"github.com/stephenbunch/wikitribe/server/middleware"
)

func main() {
	// Connect to database.
	log.Info("connecting to database")
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "wikitribe"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("failed connecting to db: %v", err)
		return
	}
	defer session.Close()

	// Run migrations.
	migrationsDir := filepath.Join(os.Getenv("GOPATH"), "src/github.com/stephenbunch/wikitribe/server/migrations")
	log.Infof("running migrations: %v", migrationsDir)
	err = migrate.Migrate(context.Background(), session, migrationsDir)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	} else {
		log.Info("finished migrations")
	}

	grpcServer := grpc.NewServer()
	service := app.NewWikiTribeService(session)
	pb.RegisterWikiTribeServiceServer(grpcServer, service)
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	router := chi.NewRouter()
	router.Use(
		chiMiddleware.Logger,
		chiMiddleware.Recoverer,
		middleware.NewGrpcWebMiddleware(wrappedGrpc).Handler, // Must come before general CORS handling,
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}).Handler,
	)

	// TODO: Figure out why the grpc web handler doesn't work without this.
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatalf("failed starting http2 server: %v", err)
	}
}
