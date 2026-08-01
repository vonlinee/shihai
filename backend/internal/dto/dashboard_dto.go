package dto

import "time"

// DashboardStatResponse describes one numeric dashboard card.
type DashboardStatResponse struct {
	Key          string `json:"key"`          // Key stable frontend identifier for the statistic.
	Title        string `json:"title"`        // Title display label for the statistic.
	Value        int64  `json:"value"`        // Value current total count.
	TrendPercent int64  `json:"trendPercent"` // TrendPercent percentage change for the latest 7 days versus the previous 7 days.
}

// DashboardActivityResponse describes a recent admin dashboard activity.
type DashboardActivityResponse struct {
	ID        string    `json:"id"`        // ID stable identifier built from activity type and entity ID.
	Action    string    `json:"action"`    // Action display label.
	Detail    string    `json:"detail"`    // Detail human-readable activity summary.
	CreatedAt time.Time `json:"createdAt"` // CreatedAt activity creation timestamp.
}

// DashboardResponse contains real data for the admin dashboard.
type DashboardResponse struct {
	Stats            []DashboardStatResponse     `json:"stats"`            // Stats dashboard summary cards.
	RecentActivities []DashboardActivityResponse `json:"recentActivities"` // RecentActivities latest cross-module activity stream.
}
