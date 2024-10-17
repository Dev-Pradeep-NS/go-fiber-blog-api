package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func HealthCheckMiddleware() fiber.Handler {
	isReady := true
	startTime := time.Now()

	return healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			uptime := time.Since(startTime)
			if uptime > 24*time.Hour {
				log.Println("Liveness probe failed: app has been running for too long")
				return false
			}
			return true
		},
		ReadinessProbe: func(c *fiber.Ctx) bool {
			if !isReady {
				log.Println("Readiness probe failed: app is not ready")
				return false
			}
			return true
		},
		LivenessEndpoint:  "/live",
		ReadinessEndpoint: "/ready",
	})
}
