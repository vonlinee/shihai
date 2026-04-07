package services

import (
	"errors"
	"shihai/internal/dto"
	"shihai/internal/models"
	"shihai/internal/repository"
)

type PoemService struct {
	poemRepo   *repository.PoemRepository
	dynastyRepo *repository.DynastyRepository
}

func NewPoemService(poemRepo *repository.PoemRepository, dynastyRepo *repository.DynastyRepository) *PoemService {
	return &PoemService{
		poemRepo:    poemRepo,
		dynastyRepo: dynastyRepo,
	}
}

// GetPoemList 获取诗词列表
func (s *PoemService) GetPoemList(req *dto.PoemListRequest) ([]dto.PoemResponse, int64, error) {
	poems, total, err := s.poemRepo.List(req.Page, req.PageSize, req.Keyword, req.Dynasty, req.Author, req.Genre)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.PoemResponse
	for _, poem := range poems {
		responses = append(responses, *s.toPoemResponse(&poem))
	}

	return responses, total, nil
}

// GetPoemByID 根据ID获取诗词
func (s *PoemService) GetPoemByID(id uint64) (*dto.PoemResponse, error) {
	poem, err := s.poemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("poem not found")
	}

	// 增加浏览量
	s.poemRepo.IncrementViews(id)

	return s.toPoemResponse(poem), nil
}

// CreatePoem 创建诗词
func (s *PoemService) CreatePoem(req *dto.PoemCreateRequest) (*dto.PoemResponse, error) {
	poem := &models.Poem{
		Title:        req.Title,
		Content:      req.Content,
		AuthorID:     req.AuthorID,
		DynastyID:    req.DynastyID,
		Genre:        req.Genre,
		Translation:  req.Translation,
		Appreciation: req.Appreciation,
		Annotation:   req.Annotation,
		AudioURL:     req.AudioURL,
		CoverImage:   req.CoverImage,
	}

	err := s.poemRepo.Create(poem)
	if err != nil {
		return nil, err
	}

	// 重新获取以加载关联数据
	return s.GetPoemByID(poem.ID)
}

// UpdatePoem 更新诗词
func (s *PoemService) UpdatePoem(id uint64, req *dto.PoemUpdateRequest) (*dto.PoemResponse, error) {
	poem, err := s.poemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("poem not found")
	}

	// 更新字段
	if req.Title != "" {
		poem.Title = req.Title
	}
	if req.Content != "" {
		poem.Content = req.Content
	}
	if req.AuthorID > 0 {
		poem.AuthorID = req.AuthorID
	}
	if req.DynastyID > 0 {
		poem.DynastyID = req.DynastyID
	}
	if req.Genre != "" {
		poem.Genre = req.Genre
	}
	if req.Translation != "" {
		poem.Translation = req.Translation
	}
	if req.Appreciation != "" {
		poem.Appreciation = req.Appreciation
	}
	if req.Annotation != "" {
		poem.Annotation = req.Annotation
	}
	if req.AudioURL != "" {
		poem.AudioURL = req.AudioURL
	}
	if req.CoverImage != "" {
		poem.CoverImage = req.CoverImage
	}

	err = s.poemRepo.Update(poem)
	if err != nil {
		return nil, err
	}

	return s.toPoemResponse(poem), nil
}

// DeletePoem 删除诗词
func (s *PoemService) DeletePoem(id uint64) error {
	return s.poemRepo.Delete(id)
}

// LikePoem 点赞诗词
func (s *PoemService) LikePoem(id uint64) error {
	return s.poemRepo.IncrementLikes(id)
}

// GetRandomPoems 随机获取诗词
func (s *PoemService) GetRandomPoems(limit int) ([]dto.PoemResponse, error) {
	poems, err := s.poemRepo.GetRandom(limit)
	if err != nil {
		return nil, err
	}

	var responses []dto.PoemResponse
	for _, poem := range poems {
		responses = append(responses, *s.toPoemResponse(&poem))
	}

	return responses, nil
}

// GetDynastyList 获取朝代列表
func (s *PoemService) GetDynastyList() ([]dto.DynastyResponse, error) {
	dynasties, err := s.dynastyRepo.List()
	if err != nil {
		return nil, err
	}

	var responses []dto.DynastyResponse
	for _, dynasty := range dynasties {
		responses = append(responses, dto.DynastyResponse{
			ID:          dynasty.ID,
			Name:        dynasty.Name,
			Period:      dynasty.Period,
			Description: dynasty.Description,
		})
	}

	return responses, nil
}

// toPoemResponse 转换为响应格式
func (s *PoemService) toPoemResponse(poem *models.Poem) *dto.PoemResponse {
	resp := &dto.PoemResponse{
		ID:           poem.ID,
		Title:        poem.Title,
		Content:      poem.Content,
		AuthorID:     poem.AuthorID,
		DynastyID:    poem.DynastyID,
		Genre:        poem.Genre,
		Translation:  poem.Translation,
		Appreciation: poem.Appreciation,
		Annotation:   poem.Annotation,
		AudioURL:     poem.AudioURL,
		CoverImage:   poem.CoverImage,
		Views:        poem.Views,
		Likes:        poem.Likes,
		Dislikes:     poem.Dislikes,
		Favorites:    poem.Favorites,
		CreatedAt:    poem.CreatedAt,
		UpdatedAt:    poem.UpdatedAt,
	}

	if poem.Author.ID > 0 {
		resp.Author = dto.PoetResponse{
			ID:        poem.Author.ID,
			Name:      poem.Author.Name,
			DynastyID: poem.Author.DynastyID,
			Biography: poem.Author.Biography,
			Avatar:    poem.Author.Avatar,
		}
	}

	if poem.Dynasty.ID > 0 {
		resp.Dynasty = dto.DynastyResponse{
			ID:          poem.Dynasty.ID,
			Name:        poem.Dynasty.Name,
			Period:      poem.Dynasty.Period,
			Description: poem.Dynasty.Description,
		}
	}

	return resp
}
