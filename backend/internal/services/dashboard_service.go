package services

import (
	"fmt"
	"math"
	"time"

	"shihai/internal/dto"
	"shihai/internal/models"
	"shihai/internal/repository"
)

// DashboardService builds admin dashboard data.
type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
	now           func() time.Time
}

// NewDashboardService creates a DashboardService.
func NewDashboardService(dashboardRepo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		now:           time.Now,
	}
}

// GetDashboard returns real summary stats and recent activities.
func (s *DashboardService) GetDashboard() (*dto.DashboardResponse, error) {
	now := s.now()
	currentStart := now.AddDate(0, 0, -7)
	previousStart := now.AddDate(0, 0, -14)

	userCounts, err := s.dashboardRepo.CountModel(&models.User{}, currentStart, previousStart)
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}
	poemCounts, err := s.dashboardRepo.CountModel(&models.Poem{}, currentStart, previousStart)
	if err != nil {
		return nil, fmt.Errorf("count poems: %w", err)
	}
	commentCounts, err := s.dashboardRepo.CountActiveComments(currentStart, previousStart)
	if err != nil {
		return nil, fmt.Errorf("count comments: %w", err)
	}
	correctionCounts, err := s.dashboardRepo.CountCorrectionsByStatus([]string{"pending", "voting"}, currentStart, previousStart)
	if err != nil {
		return nil, fmt.Errorf("count pending corrections: %w", err)
	}

	activities, err := s.dashboardRepo.RecentActivities(6)
	if err != nil {
		return nil, fmt.Errorf("list recent activities: %w", err)
	}

	return &dto.DashboardResponse{
		Stats: []dto.DashboardStatResponse{
			toDashboardStat("users", "用户总数", userCounts),
			toDashboardStat("poems", "诗词总数", poemCounts),
			toDashboardStat("comments", "评论总数", commentCounts),
			toDashboardStat("pendingCorrections", "待处理纠错", correctionCounts),
		},
		RecentActivities: toDashboardActivities(activities),
	}, nil
}

func toDashboardStat(key, title string, counts repository.DashboardCountSet) dto.DashboardStatResponse {
	return dto.DashboardStatResponse{
		Key:          key,
		Title:        title,
		Value:        counts.Total,
		TrendPercent: calculateTrendPercent(counts.CurrentPeriod, counts.PreviousPeriod),
	}
}

func toDashboardActivities(activities []repository.DashboardActivity) []dto.DashboardActivityResponse {
	responses := make([]dto.DashboardActivityResponse, 0, len(activities))
	for _, activity := range activities {
		responses = append(responses, dto.DashboardActivityResponse{
			ID:        activity.ID,
			Action:    activity.Action,
			Detail:    activity.Detail,
			CreatedAt: activity.CreatedAt,
		})
	}
	return responses
}

func calculateTrendPercent(current, previous int64) int64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return int64(math.Round(float64(current-previous) / float64(previous) * 100))
}
