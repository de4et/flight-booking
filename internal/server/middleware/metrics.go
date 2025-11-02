package middleware

import (
	"flight-booking/internal/metrics"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		defer func() {
			status := fmt.Sprint(c.Writer.Status())
			metrics.RequestsTotal.WithLabelValues(method, path, status).Inc()

			duration := time.Since(start).Seconds()
			metrics.RequestDuration.WithLabelValues(method, path, status).Observe(duration)
		}()

		c.Next()
	}
}
