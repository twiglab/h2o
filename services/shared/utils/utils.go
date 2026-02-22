package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateNo 生成业务编号
// 格式: 前缀 + YYYYMMDD + 6位随机数
func GenerateNo(prefix string) string {
	dateStr := time.Now().Format("20060102")
	randomStr := fmt.Sprintf("%06d", rand.Intn(1000000))
	return prefix + dateStr + randomStr
}

// GetPage 获取分页参数
func GetPage(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

// GetOffset 计算偏移量
func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

// ParseInt64 解析 int64
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseInt64Param 从路由参数解析 int64
func ParseInt64Param(c *gin.Context, key string) (int64, error) {
	return ParseInt64(c.Param(key))
}

// ParseInt64Query 从查询参数解析 int64，如果为空或无效返回 nil
func ParseInt64Query(c *gin.Context, key string) *int64 {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &v
}

// ParseInt16Query 从查询参数解析 int16，如果为空或无效返回 nil
func ParseInt16Query(c *gin.Context, key string) *int16 {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return nil
	}
	result := int16(v)
	return &result
}

// ParseInt8Query 从查询参数解析 int8，如果为空或无效返回 nil
func ParseInt8Query(c *gin.Context, key string) *int8 {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return nil
	}
	result := int8(v)
	return &result
}
