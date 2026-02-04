package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/mem"
)

// SystemHandler 系统状态处理器
type SystemHandler struct {
	startTime int64
}

// SystemStatsResponse 系统状态响应
type SystemStatsResponse struct {
	MemoryUsed    uint64  `json:"memory_used"`
	MemoryTotal   uint64  `json:"memory_total"`
	MemoryPercent float64 `json:"memory_percent"`
	StartTime     int64   `json:"start_time"`
}

// NewSystemHandler 创建系统状态处理器
func NewSystemHandler(startTime int64) *SystemHandler {
	return &SystemHandler{
		startTime: startTime,
	}
}

// GetStats 获取系统状态
func (h *SystemHandler) GetStats(c *gin.Context) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		slog.Error("failed to get memory stats", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取内存信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, SystemStatsResponse{
		MemoryUsed:    vmStat.Used,
		MemoryTotal:   vmStat.Total,
		MemoryPercent: vmStat.UsedPercent,
		StartTime:     h.startTime,
	})
}
