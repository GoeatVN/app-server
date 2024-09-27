package middleware

import (
	"app-server/pkg/cache"
	"app-server/pkg/response"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CachingMiddleware là middleware lưu cache kết quả của HTTP request trong Redis và trả về theo định dạng APIResponse
func CachingMiddleware(redisCache *cache.RedisCache, expiration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sử dụng URL request làm cache key
		cacheKey := "cache:" + c.Request.URL.String()

		// Kiểm tra nếu dữ liệu đã có trong cache
		cacheValue, err := redisCache.Get(cacheKey)
		if err == nil && cacheValue != "" {
			// Nếu có dữ liệu trong cache, trả về dưới định dạng APIResponse
			var cachedData interface{}
			if err := json.Unmarshal([]byte(cacheValue), &cachedData); err == nil {
				// Trả về dữ liệu đã cache theo APIResponse
				response.Success(c, cachedData)
				c.Abort() // Dừng xử lý request thêm
				return
			}
		}

		// Nếu không có dữ liệu cache, tiếp tục xử lý request
		writer := &responseBuffer{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = writer
		c.Next()

		// Nếu response thành công, lưu kết quả vào Redis Cache
		if writer.statusCode == http.StatusOK {
			var data interface{}
			if err := json.Unmarshal(writer.body.Bytes(), &data); err == nil {
				// Chuyển dữ liệu thành JSON và lưu vào cache
				cachedData, err := json.Marshal(data)
				if err == nil {
					redisCache.Set(cacheKey, cachedData, expiration)
				}
			}
		}
	}
}

// responseBuffer ghi lại response để lưu vào cache sau khi request hoàn tất
type responseBuffer struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseBuffer) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseBuffer) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
