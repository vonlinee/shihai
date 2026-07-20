package repository

import (
	"shihai/internal/models"

	"gorm.io/gorm"
)

type ForumRepository struct {
	db *gorm.DB
}

func NewForumRepository(db *gorm.DB) *ForumRepository {
	return &ForumRepository{db: db}
}

func (r *ForumRepository) CreatePost(post *models.ForumPost) error {
	return r.db.Create(post).Error
}

func (r *ForumRepository) GetPostByID(id uint64, includeDeleted bool) (*models.ForumPost, error) {
	var post models.ForumPost
	query := r.db.Preload("User")
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	if err := query.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *ForumRepository) ListPosts(page, pageSize int, keyword string, includeDeleted bool) ([]models.ForumPost, int64, error) {
	var posts []models.ForumPost
	var total int64

	query := r.db.Model(&models.ForumPost{}).Preload("User")
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("is_pinned DESC, updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *ForumRepository) UpdatePost(post *models.ForumPost) error {
	return r.db.Save(post).Error
}

func (r *ForumRepository) DeletePost(id uint64) error {
	return r.db.Model(&models.ForumPost{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

func (r *ForumRepository) SetPostPinned(id uint64, isPinned bool) error {
	return r.db.Model(&models.ForumPost{}).
		Where("id = ?", id).
		Update("is_pinned", isPinned).Error
}

func (r *ForumRepository) IncrementPostViews(id uint64) error {
	return r.db.Model(&models.ForumPost{}).
		Where("id = ? AND is_deleted = ?", id, false).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
}

func (r *ForumRepository) IncrementPostReplyCount(id uint64) error {
	return r.db.Model(&models.ForumPost{}).
		Where("id = ? AND is_deleted = ?", id, false).
		UpdateColumn("reply_count", gorm.Expr("reply_count + 1")).Error
}

func (r *ForumRepository) DecrementPostReplyCount(id uint64) error {
	return r.db.Model(&models.ForumPost{}).
		Where("id = ? AND reply_count > 0", id).
		UpdateColumn("reply_count", gorm.Expr("reply_count - 1")).Error
}

func (r *ForumRepository) CreateReply(reply *models.ForumReply) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reply).Error; err != nil {
			return err
		}
		return tx.Model(&models.ForumPost{}).
			Where("id = ? AND is_deleted = ?", reply.PostID, false).
			UpdateColumn("reply_count", gorm.Expr("reply_count + 1")).Error
	})
}

func (r *ForumRepository) GetReplyByID(id uint64, includeDeleted bool) (*models.ForumReply, error) {
	var reply models.ForumReply
	query := r.db.Preload("User").Preload("Parent.User")
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	if err := query.First(&reply, id).Error; err != nil {
		return nil, err
	}
	return &reply, nil
}

func (r *ForumRepository) ListReplies(postID uint64, page, pageSize int) ([]models.ForumReply, int64, error) {
	var replies []models.ForumReply
	var total int64

	query := r.db.Model(&models.ForumReply{}).
		Where("post_id = ? AND is_deleted = ?", postID, false).
		Preload("User").
		Preload("Parent.User")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&replies).Error; err != nil {
		return nil, 0, err
	}

	return replies, total, nil
}

func (r *ForumRepository) DeleteReply(id uint64, postID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ForumReply{}).
			Where("id = ?", id).
			Update("is_deleted", true).Error; err != nil {
			return err
		}
		return tx.Model(&models.ForumPost{}).
			Where("id = ? AND reply_count > 0", postID).
			UpdateColumn("reply_count", gorm.Expr("reply_count - 1")).Error
	})
}
