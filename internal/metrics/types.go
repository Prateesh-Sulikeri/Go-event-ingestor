package metrics

import "github.com/gin-gonic/gin"

// Shared extractor type used by limiter + metrics middleware
type ClientIDExtractor func(*gin.Context) string
