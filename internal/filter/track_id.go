package filter

import "github.com/gin-gonic/gin"
import "github.com/google/uuid"

func TrackIdFilter(c *gin.Context) {
	trackID := uuid.New().String()
	c.Set("TrackID", trackID)
	c.Writer.Header().Set("X-Track-ID", trackID)
	c.Next()
}
