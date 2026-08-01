package repository

import (
	"fmt"
	"sort"
	"time"

	"shihai/internal/models"

	"gorm.io/gorm"
)

// DashboardRepository provides read-only aggregate data for the admin dashboard.
type DashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a DashboardRepository.
func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// DashboardCountSet contains total and comparison-period counts for one model.
type DashboardCountSet struct {
	Total          int64
	CurrentPeriod  int64
	PreviousPeriod int64
}

// DashboardActivity contains a normalized activity item from different tables.
type DashboardActivity struct {
	ID        string
	Action    string
	Detail    string
	CreatedAt time.Time
}

// CountModel returns total, current-period and previous-period counts for model.
func (r *DashboardRepository) CountModel(model any, currentStart, previousStart time.Time) (DashboardCountSet, error) {
	var counts DashboardCountSet
	if err := r.db.Model(model).Count(&counts.Total).Error; err != nil {
		return counts, err
	}
	if err := r.db.Model(model).Where("created_at >= ?", currentStart).Count(&counts.CurrentPeriod).Error; err != nil {
		return counts, err
	}
	if err := r.db.Model(model).
		Where("created_at >= ? AND created_at < ?", previousStart, currentStart).
		Count(&counts.PreviousPeriod).Error; err != nil {
		return counts, err
	}
	return counts, nil
}

// CountCorrectionsByStatus returns correction counts filtered by status.
func (r *DashboardRepository) CountCorrectionsByStatus(statuses []string, currentStart, previousStart time.Time) (DashboardCountSet, error) {
	var counts DashboardCountSet
	baseQuery := func() *gorm.DB {
		return r.db.Model(&models.CorrectionRequest{}).Where("status IN ?", statuses)
	}
	if err := baseQuery().Count(&counts.Total).Error; err != nil {
		return counts, err
	}
	if err := baseQuery().Where("created_at >= ?", currentStart).Count(&counts.CurrentPeriod).Error; err != nil {
		return counts, err
	}
	if err := baseQuery().
		Where("created_at >= ? AND created_at < ?", previousStart, currentStart).
		Count(&counts.PreviousPeriod).Error; err != nil {
		return counts, err
	}
	return counts, nil
}

// CountActiveComments returns counts for comments that are not marked deleted.
func (r *DashboardRepository) CountActiveComments(currentStart, previousStart time.Time) (DashboardCountSet, error) {
	var counts DashboardCountSet
	baseQuery := func() *gorm.DB {
		return r.db.Model(&models.Comment{}).Where("is_deleted = ?", false)
	}
	if err := baseQuery().Count(&counts.Total).Error; err != nil {
		return counts, err
	}
	if err := baseQuery().Where("created_at >= ?", currentStart).Count(&counts.CurrentPeriod).Error; err != nil {
		return counts, err
	}
	if err := baseQuery().
		Where("created_at >= ? AND created_at < ?", previousStart, currentStart).
		Count(&counts.PreviousPeriod).Error; err != nil {
		return counts, err
	}
	return counts, nil
}

// RecentActivities returns the latest cross-module dashboard activity items.
func (r *DashboardRepository) RecentActivities(limit int) ([]DashboardActivity, error) {
	activities := make([]DashboardActivity, 0, limit*4)

	if err := r.appendRecentUsers(&activities, limit); err != nil {
		return nil, err
	}
	if err := r.appendRecentPoems(&activities, limit); err != nil {
		return nil, err
	}
	if err := r.appendRecentComments(&activities, limit); err != nil {
		return nil, err
	}
	if err := r.appendRecentCorrections(&activities, limit); err != nil {
		return nil, err
	}

	sort.SliceStable(activities, func(i, j int) bool {
		return activities[i].CreatedAt.After(activities[j].CreatedAt)
	})
	if len(activities) > limit {
		activities = activities[:limit]
	}
	return activities, nil
}

func (r *DashboardRepository) appendRecentUsers(activities *[]DashboardActivity, limit int) error {
	var users []models.User
	if err := r.db.Model(&models.User{}).Order("created_at DESC").Limit(limit).Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		name := firstNonEmpty(user.Name, user.Username, "未知用户")
		*activities = append(*activities, DashboardActivity{
			ID:        fmt.Sprintf("user-%d", user.ID),
			Action:    "用户注册",
			Detail:    fmt.Sprintf("新用户 %s 注册了账号", name),
			CreatedAt: user.CreatedAt,
		})
	}
	return nil
}

func (r *DashboardRepository) appendRecentPoems(activities *[]DashboardActivity, limit int) error {
	var poems []models.Poem
	if err := r.db.Model(&models.Poem{}).Order("created_at DESC").Limit(limit).Find(&poems).Error; err != nil {
		return err
	}
	for _, poem := range poems {
		*activities = append(*activities, DashboardActivity{
			ID:        fmt.Sprintf("poem-%d", poem.ID),
			Action:    "诗词添加",
			Detail:    fmt.Sprintf("添加了新诗词《%s》", poem.Title),
			CreatedAt: poem.CreatedAt,
		})
	}
	return nil
}

func (r *DashboardRepository) appendRecentComments(activities *[]DashboardActivity, limit int) error {
	var comments []models.Comment
	if err := r.db.Model(&models.Comment{}).
		Preload("User").
		Preload("Poem").
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Limit(limit).
		Find(&comments).Error; err != nil {
		return err
	}
	for _, comment := range comments {
		name := comment.VisitorName
		if comment.User != nil {
			name = firstNonEmpty(comment.User.Name, comment.User.Username, name)
		}
		name = firstNonEmpty(name, "用户")
		*activities = append(*activities, DashboardActivity{
			ID:        fmt.Sprintf("comment-%d", comment.ID),
			Action:    "评论发布",
			Detail:    fmt.Sprintf("%s 评论了《%s》", name, comment.Poem.Title),
			CreatedAt: comment.CreatedAt,
		})
	}
	return nil
}

func (r *DashboardRepository) appendRecentCorrections(activities *[]DashboardActivity, limit int) error {
	var corrections []models.CorrectionRequest
	if err := r.db.Model(&models.CorrectionRequest{}).
		Preload("Poem").
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&corrections).Error; err != nil {
		return err
	}
	for _, correction := range corrections {
		name := firstNonEmpty(correction.User.Name, correction.User.Username, "用户")
		*activities = append(*activities, DashboardActivity{
			ID:        fmt.Sprintf("correction-%d", correction.ID),
			Action:    "诗词纠错",
			Detail:    fmt.Sprintf("%s 提交了《%s》的纠错申请", name, correction.Poem.Title),
			CreatedAt: correction.CreatedAt,
		})
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
