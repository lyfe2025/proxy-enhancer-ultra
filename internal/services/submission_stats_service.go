package services

import (
	"errors"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubmissionStatsService 提交统计服务
type SubmissionStatsService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewSubmissionStatsService 创建新的提交统计服务
func NewSubmissionStatsService(db *gorm.DB, logger logger.Logger) *SubmissionStatsService {
	return &SubmissionStatsService{
		db:     db,
		logger: logger,
	}
}

// GetSubmissionStats 获取提交统计信息
func (s *SubmissionStatsService) GetSubmissionStats(popupID uuid.UUID) (*SubmissionStats, error) {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, popupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, err
	}

	stats := &SubmissionStats{
		PopupID: popupID,
	}

	// 获取总提交数
	if err := s.db.Model(&models.Submission{}).Where("popup_id = ?", popupID).Count(&stats.TotalSubmissions).Error; err != nil {
		return nil, err
	}

	// 获取今日提交数
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	if err := s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ? AND created_at < ?", popupID, today, tomorrow).Count(&stats.TodaySubmissions).Error; err != nil {
		return nil, err
	}

	// 获取本周提交数
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	weekEnd := weekStart.Add(7 * 24 * time.Hour)
	if err := s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ? AND created_at < ?", popupID, weekStart, weekEnd).Count(&stats.WeekSubmissions).Error; err != nil {
		return nil, err
	}

	// 获取本月提交数
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	if err := s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ? AND created_at < ?", popupID, monthStart, monthEnd).Count(&stats.MonthSubmissions).Error; err != nil {
		return nil, err
	}

	// 获取最近提交时间
	var lastSubmission models.Submission
	if err := s.db.Where("popup_id = ?", popupID).Order("created_at DESC").First(&lastSubmission).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// 没有提交记录，lastSubmission 保持为 nil
	} else {
		stats.LastSubmission = &lastSubmission.CreatedAt
	}

	return stats, nil
}

// GetOverallStats 获取总体统计信息
func (s *SubmissionStatsService) GetOverallStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总提交数
	var totalSubmissions int64
	if err := s.db.Model(&models.Submission{}).Count(&totalSubmissions).Error; err != nil {
		return nil, err
	}
	stats["total_submissions"] = totalSubmissions

	// 今日提交数
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	var todaySubmissions int64
	if err := s.db.Model(&models.Submission{}).Where("created_at >= ? AND created_at < ?", today, tomorrow).Count(&todaySubmissions).Error; err != nil {
		return nil, err
	}
	stats["today_submissions"] = todaySubmissions

	// 本周提交数
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	weekEnd := weekStart.Add(7 * 24 * time.Hour)
	var weekSubmissions int64
	if err := s.db.Model(&models.Submission{}).Where("created_at >= ? AND created_at < ?", weekStart, weekEnd).Count(&weekSubmissions).Error; err != nil {
		return nil, err
	}
	stats["week_submissions"] = weekSubmissions

	// 本月提交数
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	var monthSubmissions int64
	if err := s.db.Model(&models.Submission{}).Where("created_at >= ? AND created_at < ?", monthStart, monthEnd).Count(&monthSubmissions).Error; err != nil {
		return nil, err
	}
	stats["month_submissions"] = monthSubmissions

	// 活跃弹窗数（有提交记录的弹窗）
	var activePopups int64
	if err := s.db.Model(&models.Submission{}).Distinct("popup_id").Count(&activePopups).Error; err != nil {
		return nil, err
	}
	stats["active_popups"] = activePopups

	return stats, nil
}

// GetSubmissionTrends 获取提交趋势数据
func (s *SubmissionStatsService) GetSubmissionTrends(popupID uuid.UUID, days int) ([]map[string]interface{}, error) {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, popupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, err
	}

	var trends []map[string]interface{}
	today := time.Now().Truncate(24 * time.Hour)

	for i := days - 1; i >= 0; i-- {
		date := today.AddDate(0, 0, -i)
		nextDate := date.Add(24 * time.Hour)

		var count int64
		if err := s.db.Model(&models.Submission{}).
			Where("popup_id = ? AND created_at >= ? AND created_at < ?", popupID, date, nextDate).
			Count(&count).Error; err != nil {
			return nil, err
		}

		trends = append(trends, map[string]interface{}{
			"date":  date.Format("2006-01-02"),
			"count": count,
		})
	}

	return trends, nil
}
