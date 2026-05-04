package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

// GetByID 根据ID获取评论
func (r *CommentRepository) GetByID(id uint64) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").Preload("Replies.User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListByPoem 获取诗词的评论列表
func (r *CommentRepository) ListByPoem(poemID uint64, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).
		Where("poem_id = ? AND parent_id IS NULL", poemID).
		Preload("User").
		Preload("Replies.User")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// Update 更新评论
func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论（软删除）
func (r *CommentRepository) Delete(id uint64) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// IncrementLikes 增加点赞
func (r *CommentRepository) IncrementLikes(id uint64) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + 1")).Error
}

// IncrementDislikes 增加点踩
func (r *CommentRepository) IncrementDislikes(id uint64) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn("dislikes", gorm.Expr("dislikes + 1")).Error
}

// IncrementReplyCount 增加回复数
func (r *CommentRepository) IncrementReplyCount(id uint64) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn("reply_count", gorm.Expr("reply_count + 1")).Error
}

// GetVote 获取用户投票
func (r *CommentRepository) GetVote(commentID uint64, userID *uint64, visitorID string) (*models.CommentVote, error) {
	var vote models.CommentVote
	query := r.db.Where("comment_id = ?", commentID)
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	} else {
		query = query.Where("visitor_id = ?", visitorID)
	}
	err := query.First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

// CreateVote 创建投票
func (r *CommentRepository) CreateVote(vote *models.CommentVote) error {
	return r.db.Create(vote).Error
}

// UpdateVote 更新投票
func (r *CommentRepository) UpdateVote(vote *models.CommentVote) error {
	return r.db.Save(vote).Error
}
