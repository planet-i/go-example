package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   uint64
	Name string
}

func main() {
	r := gin.New()

	// 单个路由
	users := []User{{ID: 123, Name: "张三"}, {ID: 256, Name: "李四"}}
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, users)
	})
	// 路由组
	user := r.Group("/user")
	{
		// 实现对应的HTTP Method方法注册。
		user.POST("/", func(c *gin.Context) { //创建一个用户
			c.JSON(http.StatusOK, gin.H{"message": "User created"})
		})
		user.DELETE("/123", func(c *gin.Context) { //删除ID为123的用户
			c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
		})
		user.PUT("/123", func(c *gin.Context) { //更新ID为123的用户
			c.JSON(http.StatusOK, gin.H{"message": "User updated"})
		})
		user.PATCH("/123", func(c *gin.Context) { //更新ID为123用户部分信息
			c.JSON(http.StatusOK, gin.H{"message": "User partially updated"})
		})
	}

	// 应用中间件：根据参数字段决定返回的字段
	r.Use(FieldFilterMiddleware())

	// 测试中间件路由
	r.GET("/users", func(c *gin.Context) {
		users := []map[string]interface{}{
			{
				"id":    1,
				"name":  "Alice",
				"email": "alice@example.com",
				"age":   30,
				"address": map[string]string{
					"city":    "New York",
					"country": "USA",
				},
			},
			{
				"id":    2,
				"name":  "Bob",
				"email": "bob@example.com",
				"age":   25,
				"address": map[string]string{
					"city":    "London",
					"country": "UK",
				},
			},
		}
		c.JSON(http.StatusOK, users)
	})

	r.Run(":8080")
}

// 字段过滤中间件
func FieldFilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建自定义响应写入器
		writer := &responseWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = writer

		// 处理请求
		c.Next()

		var changed bool
		defer func() {
			if !changed {
				// 如果没有改变，则直接将缓存的内容写出
				writer.ResponseWriter.Write(writer.body.Bytes())
			}
			writer.body.Reset() // 重置缓冲区以备后续请求使用
		}()
		// 获取客户端请求的字段
		fieldsParam := c.Query("fields")
		if fieldsParam == "" {
			return // 未指定字段，返回全部数据
		}

		// 解析字段列表
		fields := strings.Split(fieldsParam, ",")
		if len(fields) == 0 {
			return
		}

		// 获取响应状态码
		// statusCode := writer.Status()

		// 检查响应是否为JSON
		contentType := writer.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return
		}

		// 解析响应数据
		var resp interface{}
		if err := json.Unmarshal(writer.body.Bytes(), &resp); err != nil {
			return
		}

		// 过滤数据
		filteredData := filterFields(resp, fields)

		// 重新序列化并返回过滤后的数据
		filteredJSON, err := json.Marshal(filteredData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize response"})
			return
		}

		changed = true
		// 重置响应
		//c.Writer = c.Writer.(*responseWriter).ResponseWriter
		//c.Writer.Header().Set("Content-Type", "application/json")
		//c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(filteredJSON)))
		////c.Writer.WriteHeader(400)
		//fmt.Println("___---wirete")
		//fmt.Println(c.Writer.Size())
		//c.Writer.Write(filteredJSON)
		_, _ = writer.ResponseWriter.Write(filteredJSON)
		// 阻止后续处理
		c.Abort()
	}
}

// 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 重写Write方法，捕获响应内容
func (w *responseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
	//fmt.Println("-------", w.body.String())
	//return w.ResponseWriter.Write(b)
}

// 递归过滤字段
func filterFields(data interface{}, fields []string) interface{} {
	fmt.Println("filterFields:", fields)
	fmt.Println("data:", data)
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for _, field := range fields {
			// 处理嵌套字段，如 "user.name"
			parts := strings.SplitN(field, ".", 2)
			if len(parts) == 1 {
				if val, exists := v[field]; exists {
					result[field] = val
				}
			} else {
				// 处理嵌套字段
				nestedField := parts[0]
				remainingFields := []string{parts[1]}
				if nestedVal, exists := v[nestedField]; exists {
					result[nestedField] = filterFields(nestedVal, remainingFields)
				}
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = filterFields(item, fields)
		}
		return result
	default:
		return v
	}
}
