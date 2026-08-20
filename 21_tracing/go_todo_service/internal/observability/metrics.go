package observability

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/akozadaev/go_todo_service/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Metrics struct {
	Registry     *prometheus.Registry
	httpRequests *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
	httpInFlight *prometheus.GaugeVec
	grpcRequests *prometheus.CounterVec
	grpcDuration *prometheus.HistogramVec
	grpcInFlight *prometheus.GaugeVec
}

func NewMetrics(sqlDB *sql.DB) *Metrics {
	m := &Metrics{
		Registry: prometheus.NewRegistry(),
		httpRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "todo", Subsystem: "http", Name: "requests_total",
			Help: "Total number of HTTP requests.",
		}, []string{"method", "route", "status"}),
		httpDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "todo", Subsystem: "http", Name: "request_duration_seconds",
			Help: "HTTP request latency.", Buckets: prometheus.DefBuckets,
		}, []string{"method", "route", "status"}),
		httpInFlight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "todo", Subsystem: "http", Name: "requests_in_flight",
			Help: "Current number of HTTP requests.",
		}, []string{"method"}),
		grpcRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "todo", Subsystem: "grpc", Name: "requests_total",
			Help: "Total number of gRPC requests.",
		}, []string{"method", "type", "code"}),
		grpcDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "todo", Subsystem: "grpc", Name: "request_duration_seconds",
			Help: "gRPC request latency.", Buckets: prometheus.DefBuckets,
		}, []string{"method", "type", "code"}),
		grpcInFlight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "todo", Subsystem: "grpc", Name: "requests_in_flight",
			Help: "Current number of gRPC requests.",
		}, []string{"method", "type"}),
	}

	m.Registry.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		m.httpRequests, m.httpDuration, m.httpInFlight,
		m.grpcRequests, m.grpcDuration, m.grpcInFlight,
	)
	if sqlDB != nil {
		m.Registry.MustRegister(newDBStatsCollector(sqlDB))
	}
	return m
}

func (m *Metrics) HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		m.httpInFlight.WithLabelValues(method).Inc()
		defer m.httpInFlight.WithLabelValues(method).Dec()
		started := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		statusCode := strconv.Itoa(c.Writer.Status())
		m.httpRequests.WithLabelValues(method, route, statusCode).Inc()
		observeWithTrace(c.Request.Context(), m.httpDuration.WithLabelValues(method, route, statusCode), time.Since(started).Seconds())
	}
}

func (m *Metrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		const rpcType = "unary"
		started := time.Now()
		m.grpcInFlight.WithLabelValues(info.FullMethod, rpcType).Inc()
		defer m.grpcInFlight.WithLabelValues(info.FullMethod, rpcType).Dec()
		log := logger.FromContext(ctx).With(zap.String("rpc.method", info.FullMethod), zap.String("rpc.type", rpcType))
		log.Info("gRPC request started")
		resp, err := handler(ctx, req)
		code := status.Code(err).String()
		duration := time.Since(started)
		m.grpcRequests.WithLabelValues(info.FullMethod, rpcType, code).Inc()
		observeWithTrace(ctx, m.grpcDuration.WithLabelValues(info.FullMethod, rpcType, code), duration.Seconds())
		log.Info("gRPC request completed", zap.String("rpc.code", code), zap.Duration("duration", duration), zap.Error(err))
		return resp, err
	}
}

func (m *Metrics) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		const rpcType = "stream"
		ctx := stream.Context()
		started := time.Now()
		m.grpcInFlight.WithLabelValues(info.FullMethod, rpcType).Inc()
		defer m.grpcInFlight.WithLabelValues(info.FullMethod, rpcType).Dec()
		log := logger.FromContext(ctx).With(zap.String("rpc.method", info.FullMethod), zap.String("rpc.type", rpcType))
		log.Info("gRPC request started")
		err := handler(srv, stream)
		code := status.Code(err).String()
		duration := time.Since(started)
		m.grpcRequests.WithLabelValues(info.FullMethod, rpcType, code).Inc()
		observeWithTrace(ctx, m.grpcDuration.WithLabelValues(info.FullMethod, rpcType, code), duration.Seconds())
		log.Info("gRPC request completed", zap.String("rpc.code", code), zap.Duration("duration", duration), zap.Error(err))
		return err
	}
}

func observeWithTrace(ctx context.Context, observer prometheus.Observer, value float64) {
	spanContext := trace.SpanContextFromContext(ctx)
	if exemplarObserver, ok := observer.(prometheus.ExemplarObserver); ok && spanContext.IsSampled() {
		exemplarObserver.ObserveWithExemplar(value, prometheus.Labels{"trace_id": spanContext.TraceID().String()})
		return
	}
	observer.Observe(value)
}

func newDBStatsCollector(db *sql.DB) prometheus.Collector {
	return collectors.NewDBStatsCollector(db, "todo")
}
