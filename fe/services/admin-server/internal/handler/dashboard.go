package handler

import (
	"time"

	"admin-server/internal/service"
	"shared/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// Dashboard 仪表盘统计
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	stats, err := h.dashboardService.GetDashboard()
	if err != nil {
		response.ServerError(c, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// ConsumptionReport 用量报表
func (h *DashboardHandler) ConsumptionReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	groupBy := c.DefaultQuery("group_by", "day")

	stats, err := h.dashboardService.GetConsumptionReport(startDate, endDate, groupBy)
	if err != nil {
		response.ServerError(c, "获取报表失败")
		return
	}

	response.Success(c, stats)
}

// RevenueReport 收入报表
func (h *DashboardHandler) RevenueReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	groupBy := c.DefaultQuery("group_by", "day")

	stats, err := h.dashboardService.GetRevenueReport(startDate, endDate, groupBy)
	if err != nil {
		response.ServerError(c, "获取报表失败")
		return
	}

	response.Success(c, stats)
}

// CollectionStats 采集统计
func (h *DashboardHandler) CollectionStats(c *gin.Context) {
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	stats, err := h.dashboardService.GetCollectionStats(date)
	if err != nil {
		response.ServerError(c, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}
