package services

import (
	"errors"
	"shihai/internal/dto"
	"shihai/internal/models"
	"shihai/internal/repository"
)

type PoemService struct {
	poemRepo    *repository.PoemRepository
	dynastyRepo *repository.DynastyRepository
	poetRepo    *repository.PoetRepository
}

func NewPoemService(poemRepo *repository.PoemRepository, dynastyRepo *repository.DynastyRepository, poetRepo *repository.PoetRepository) *PoemService {
	return &PoemService{
		poemRepo:    poemRepo,
		dynastyRepo: dynastyRepo,
		poetRepo:    poetRepo,
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
	// 如果没有朝代ID但有朝代名称，自动查找或创建朝代
	dynastyID := req.DynastyID
	if dynastyID == 0 && req.DynastyName != "" {
		dynasty, err := s.dynastyRepo.GetByName(req.DynastyName)
		if err != nil {
			// 朝代不存在，创建新的
			dynasty = &models.Dynasty{Name: req.DynastyName}
			if createErr := s.dynastyRepo.Create(dynasty); createErr != nil {
				return nil, createErr
			}
		}
		dynastyID = dynasty.ID
	}

	// 如果没有作者ID但有作者名称，自动查找或创建作者
	authorID := req.AuthorID
	if authorID == 0 && req.AuthorName != "" {
		poet, err := s.poetRepo.GetByName(req.AuthorName)
		if err != nil {
			// 作者不存在，创建新的
			poet = &models.Poet{Name: req.AuthorName, DynastyID: dynastyID}
			if createErr := s.poetRepo.Create(poet); createErr != nil {
				return nil, createErr
			}
		}
		authorID = poet.ID
	}

	poem := &models.Poem{
		Title:        req.Title,
		Content:      req.Content,
		AuthorID:     authorID,
		DynastyID:    dynastyID,
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

// GetPoetList 获取诗人列表
func (s *PoemService) GetPoetList(keyword string) ([]dto.PoetResponse, error) {
	poets, err := s.poetRepo.List(keyword)
	if err != nil {
		return nil, err
	}

	var responses []dto.PoetResponse
	for _, poet := range poets {
		responses = append(responses, dto.PoetResponse{
			ID:        poet.ID,
			Name:      poet.Name,
			DynastyID: poet.DynastyID,
			Biography: poet.Biography,
			Avatar:    poet.Avatar,
		})
	}

	return responses, nil
}

// GetGenreList 获取体裁列表
func (s *PoemService) GetGenreList() ([]string, error) {
	var genres []string
	err := s.poemRepo.DB().Model(&models.Poem{}).
		Where("genre != '' AND genre IS NOT NULL").
		Distinct("genre").
		Pluck("genre", &genres).Error
	return genres, err
}

// CreateDynasty 创建朝代（如果不存在则创建）
func (s *PoemService) CreateDynasty(req *dto.DynastyCreateRequest) (*dto.DynastyResponse, error) {
	dynasty := &models.Dynasty{
		Name:        req.Name,
		Period:      req.Period,
		Description: req.Description,
	}
	err := s.dynastyRepo.Create(dynasty)
	if err != nil {
		return nil, err
	}
	return &dto.DynastyResponse{
		ID:          dynasty.ID,
		Name:        dynasty.Name,
		Period:      dynasty.Period,
		Description: dynasty.Description,
	}, nil
}

// UpdateDynasty 更新朝代
func (s *PoemService) UpdateDynasty(id uint64, req *dto.DynastyUpdateRequest) (*dto.DynastyResponse, error) {
	dynasty, err := s.dynastyRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("dynasty not found")
	}
	if req.Name != "" {
		dynasty.Name = req.Name
	}
	if req.Period != "" {
		dynasty.Period = req.Period
	}
	if req.Description != "" {
		dynasty.Description = req.Description
	}
	err = s.dynastyRepo.Update(dynasty)
	if err != nil {
		return nil, err
	}
	return &dto.DynastyResponse{
		ID:          dynasty.ID,
		Name:        dynasty.Name,
		Period:      dynasty.Period,
		Description: dynasty.Description,
	}, nil
}

// DeleteDynasty 删除朝代
func (s *PoemService) DeleteDynasty(id uint64) error {
	return s.dynastyRepo.Delete(id)
}

// CreatePoet 创建诗人（如果不存在则创建）
func (s *PoemService) CreatePoet(req *dto.PoetCreateRequest) (*dto.PoetResponse, error) {
	poet := &models.Poet{
		Name:      req.Name,
		DynastyID: req.DynastyID,
		Biography: req.Biography,
		Avatar:    req.Avatar,
		BirthYear: req.BirthYear,
		DeathYear: req.DeathYear,
	}
	err := s.poetRepo.Create(poet)
	if err != nil {
		return nil, err
	}
	return &dto.PoetResponse{
		ID:        poet.ID,
		Name:      poet.Name,
		DynastyID: poet.DynastyID,
		Biography: poet.Biography,
		Avatar:    poet.Avatar,
	}, nil
}

// UpdatePoet 更新诗人
func (s *PoemService) UpdatePoet(id uint64, req *dto.PoetUpdateRequest) (*dto.PoetResponse, error) {
	poet, err := s.poetRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("poet not found")
	}
	if req.Name != "" {
		poet.Name = req.Name
	}
	if req.DynastyID > 0 {
		poet.DynastyID = req.DynastyID
	}
	if req.Biography != "" {
		poet.Biography = req.Biography
	}
	if req.Avatar != "" {
		poet.Avatar = req.Avatar
	}
	if req.BirthYear > 0 {
		poet.BirthYear = req.BirthYear
	}
	if req.DeathYear > 0 {
		poet.DeathYear = req.DeathYear
	}
	err = s.poetRepo.Update(poet)
	if err != nil {
		return nil, err
	}
	return &dto.PoetResponse{
		ID:        poet.ID,
		Name:      poet.Name,
		DynastyID: poet.DynastyID,
		Biography: poet.Biography,
		Avatar:    poet.Avatar,
	}, nil
}

// DeletePoet 删除诗人
func (s *PoemService) DeletePoet(id uint64) error {
	return s.poetRepo.Delete(id)
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
