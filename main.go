package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"oauth-server/app/controller"
	"oauth-server/app/helper"
	"oauth-server/app/middleware"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/app/service"
	"oauth-server/config"
	"oauth-server/graphql/generated"
	"oauth-server/graphql/resolver"
	"oauth-server/grpc/user"
	"oauth-server/package/database"
	logger "oauth-server/package/log"
	_validator "oauth-server/package/validator"
	"oauth-server/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"

	_ "ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	conf := config.GetConfiguration()

	// Register repositories
	postgresDB := database.GetPostgres()
	postgresRepo := postgres_repository.RegisterPostgresRepositories(postgresDB)
	mws := middleware.RegisterMiddleware()

	// Register Others
	helpers := helper.RegisterHelpers(postgresRepo)
	services := service.RegisterServices(helpers, postgresRepo)

	// Start HTTP Server
	srv := initHTTPServer(conf, services)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()

	// Start GRPC Server
	lis, grpcServer := initGRPCServer(conf, postgresRepo, helpers, mws)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Panicf("GRPC Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	grpcServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 3 seconds.")
	}
	log.Println("Server exiting")
}

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	database.InitPostgres()
	logger.Init()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           config.GetConfiguration().Server.SentryDNS,
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Panicf("Sentry initialization failed: %v\n", err)
	}
}

func initHTTPServer(conf config.Configuration, services service.ServiceCollections) *http.Server {
	// Run gin server
	gin.SetMode(conf.Server.Mode)
	app := gin.Default()
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	// Register validator
	v := binding.Validator.Engine().(*validator.Validate)
	v.SetTagName("validate")
	_validator.RegisterCustomValidators(v)

	// Register controllers
	app.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})
	controller.RegisterControllers(app, services)
	app.POST("/query", registerGraphQL(services))

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: app,
	}

	return srv
}

func initGRPCServer(
	conf config.Configuration,
	postgresRepo repository.RepositoryCollections,
	helpers helper.HelperCollections,
	mws middleware.MiddlewareCollections,
) (net.Listener, *grpc.Server) {
	// Start GRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", conf.GRPC.Port))
	if err != nil {
		log.Panicf("GRPC Failed to listen: %v", err)
		return nil, nil
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(mws.RPCMiddleware.Handler),
	}
	grpcServer := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(grpcServer, user.NewUserServiceServer(postgresRepo, helpers))

	return lis, grpcServer
}

func registerGraphQL(
	services service.ServiceCollections,
) gin.HandlerFunc {
	h := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		Services: services,
	}}))

	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
