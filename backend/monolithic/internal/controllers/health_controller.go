package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type HealthController struct {
	session *gocql.Session
}

func NewHealthController(session *gocql.Session) *HealthController {
	return &HealthController{session: session}
}

// is the process up? a failing liveness probe makes k8s restart the pod
func (h *HealthController) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// readiness
func (h *HealthController) Readyz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var t gocql.UUID
	if err := h.session.Query("SELECT now() FROM system.local").
		WithContext(ctx).
		Scan(&t); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
