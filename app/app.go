package app

import (
	"context"
	"fmt"
	"go-gin-swagger-test/app/controller"
	"go-gin-swagger-test/app/logger"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"go-gin-swagger-test/app/controller/docs"

	ginzap "github.com/gin-contrib/zap"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type App struct {
	Controller *controller.Controller
	Logger     *logger.ZapLogger
}

func NewApp(c *controller.Controller, l *logger.ZapLogger) *App {
	return &App{Controller: c, Logger: l}
}

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
func (a *App) Start() {

	// Set swagger info
	docs.SwaggerInfo.Title = "Account Rest API"
	docs.SwaggerInfo.Description = "Base Account System"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = "accounts.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	a.Logger.BuildLogger()

	f, err := os.Create("traces.txt")
	if err != nil {
		a.Logger.Adapter.Error(err.Error())
	}

	tp, err := initTracer(f)
	if err != nil {
		a.Logger.Adapter.Error(err.Error())
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			a.Logger.Adapter.Error(fmt.Sprintf("Error shutting down tracer provider: %v", err))
		}
	}()

	r := gin.Default()

	r.Use(otelgin.Middleware("gin-swagger-service"))
	// Using the ginzap middleware to use zap as the logger of gin instad the deffault logger
	r.Use(ginzap.Ginzap(a.Logger.Adapter, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(a.Logger.Adapter, true))

	v1 := r.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.GET("", a.Controller.GetAccounts)
			accounts.GET(":id", a.Controller.GetAccountById)
			accounts.POST("", a.Controller.CreateAccount)
			accounts.PUT(":id", a.Controller.UpdateNameById)
			accounts.DELETE(":id", a.Controller.DeleteAccountById)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}

func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),

		stdouttrace.WithPrettyPrint(),

		stdouttrace.WithoutTimestamps(),
	)
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("gin-swagger-service"),
			semconv.ServiceVersionKey.String("0.0.1v"),
			attribute.String("enviroment", "demo"),
		),
	)

	return r
}

func initTracer(w io.Writer) (*trace.TracerProvider, error) {
	exp, err := newExporter(w)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
