package services

import (
	"errors"
	"fmt"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"gorm.io/gorm"
)

// PopupStatsService 弹窗统计服务
type PopupStatsService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewPopupStatsService 创建新的弹窗统计服务
func NewPopupStatsService(db *gorm.DB, logger logger.Logger) *PopupStatsService {
	return &PopupStatsService{
		db:     db,
		logger: logger,
	}
}

// GetPopupStats 获取弹窗统计信息
func (s *PopupStatsService) GetPopupStats(id uint) (map[string]interface{}, error) {
	// 检查弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, fmt.Errorf("failed to check popup: %w", err)
	}

	// 获取统计信息
	stats := make(map[string]interface{})

	// 总提交数
	var totalSubmissions int64
	s.db.Model(&models.Submission{}).Where("popup_id = ?", id).Count(&totalSubmissions)
	stats["total_submissions"] = totalSubmissions

	// 今日提交数
	var todaySubmissions int64
	today := time.Now().Truncate(24 * time.Hour)
	s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ?", id, today).Count(&todaySubmissions)
	stats["today_submissions"] = todaySubmissions

	// 本周提交数
	var weekSubmissions int64
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())).Truncate(24 * time.Hour)
	s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ?", id, weekStart).Count(&weekSubmissions)
	stats["week_submissions"] = weekSubmissions

	// 本月提交数
	var monthSubmissions int64
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ?", id, monthStart).Count(&monthSubmissions)
	stats["month_submissions"] = monthSubmissions

	// 最近提交时间
	var lastSubmission models.Submission
	if err := s.db.Where("popup_id = ?", id).Order("created_at DESC").First(&lastSubmission).Error; err == nil {
		stats["last_submission_at"] = lastSubmission.CreatedAt
	} else {
		stats["last_submission_at"] = nil
	}

	return stats, nil
}

// GetOverallStats 获取总体统计信息
func (s *PopupStatsService) GetOverallStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总弹窗数
	var totalPopups int64
	if err := s.db.Model(&models.Popup{}).Count(&totalPopups).Error; err != nil {
		return nil, fmt.Errorf("failed to count total popups: %w", err)
	}
	stats["total_popups"] = totalPopups

	// 活跃弹窗数
	var activePopups int64
	if err := s.db.Model(&models.Popup{}).Where("is_active = ?", true).Count(&activePopups).Error; err != nil {
		return nil, fmt.Errorf("failed to count active popups: %w", err)
	}
	stats["active_popups"] = activePopups

	// 总提交数
	var totalSubmissions int64
	if err := s.db.Model(&models.Submission{}).Count(&totalSubmissions).Error; err != nil {
		return nil, fmt.Errorf("failed to count total submissions: %w", err)
	}
	stats["total_submissions"] = totalSubmissions

	// 今日提交数
	var todaySubmissions int64
	today := time.Now().Truncate(24 * time.Hour)
	if err := s.db.Model(&models.Submission{}).Where("created_at >= ?", today).Count(&todaySubmissions).Error; err != nil {
		return nil, fmt.Errorf("failed to count today submissions: %w", err)
	}
	stats["today_submissions"] = todaySubmissions

	// 计算转化率（如果有总访问数据的话）
	if totalPopups > 0 && totalSubmissions > 0 {
		conversionRate := float64(totalSubmissions) / float64(totalPopups)
		stats["conversion_rate"] = conversionRate
	} else {
		stats["conversion_rate"] = 0.0
	}

	return stats, nil
}

// GetPopupPerformance 获取弹窗性能数据
func (s *PopupStatsService) GetPopupPerformance(id uint, days int) (map[string]interface{}, error) {
	// 检查弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, fmt.Errorf("failed to check popup: %w", err)
	}

	performance := make(map[string]interface{})

	// 计算开始时间
	startTime := time.Now().AddDate(0, 0, -days).Truncate(24 * time.Hour)

	// 按天统计提交数
	type DailyStats struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	var dailyStats []DailyStats
	for i := 0; i < days; i++ {
		date := startTime.AddDate(0, 0, i)
		nextDate := date.AddDate(0, 0, 1)

		var count int64
		s.db.Model(&models.Submission{}).
			Where("popup_id = ? AND created_at >= ? AND created_at < ?", id, date, nextDate).
			Count(&count)

		dailyStats = append(dailyStats, DailyStats{
			Date:  date.Format("2006-01-02"),
			Count: count,
		})
	}

	performance["daily_stats"] = dailyStats

	// 计算总提交数（指定时间范围内）
	var totalInPeriod int64
	s.db.Model(&models.Submission{}).
		Where("popup_id = ? AND created_at >= ?", id, startTime).
		Count(&totalInPeriod)
	performance["total_in_period"] = totalInPeriod

	// 计算平均每日提交数
	avgDaily := float64(totalInPeriod) / float64(days)
	performance["avg_daily_submissions"] = avgDaily

	return performance, nil
}
