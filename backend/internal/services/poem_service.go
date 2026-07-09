package services

import (
	"errors"

	"shihai/internal/dto"
	"shihai/internal/models"
)

type PoemService struct {
	poemRepo    poemRepository
	dynastyRepo dynastyRepository
	authorRepo  authorRepository
	poetRepo    poetRepository
}

type poemRepository interface {
	Create(poem *models.Poem) error
	GetByID(id uint64) (*models.Poem, error)
	Update(poem *models.Poem) error
	Delete(id uint64) error
	BatchDelete(ids []uint64) error
	List(page, pageSize int, keyword, dynasty, author, genre string) ([]models.Poem, int64, error)
	IncrementViews(id uint64) error
	IncrementLikes(id uint64) error
	GetRandom(limit int) ([]models.Poem, error)
	DistinctGenres() ([]string, error)
}

type dynastyRepository interface {
	Create(dynasty *models.Dynasty) error
	GetByID(id uint64) (*models.Dynasty, error)
	GetByName(name string) (*models.Dynasty, error)
	Update(dynasty *models.Dynasty) error
	Delete(id uint64) error
	BatchDelete(ids []uint64) error
	List() ([]models.Dynasty, error)
}

type authorRepository interface {
	Create(author *models.Author) error
	GetByID(id uint64) (*models.Author, error)
	GetByName(name string) (*models.Author, error)
	List(keyword string) ([]models.Author, error)
	Update(author *models.Author) error
	Delete(id uint64) error
}

type poetRepository interface {
	Create(poet *models.Poet) error
	GetByID(id uint64) (*models.Poet, error)
	GetByAuthorID(authorID uint64) (*models.Poet, error)
	List(keyword string, page, pageSize int) ([]models.Poet, int64, error)
	Update(poet *models.Poet) error
	Delete(id uint64) error
	BatchDelete(ids []uint64) error
}

func NewPoemService(poemRepo poemRepository, dynastyRepo dynastyRepository, authorRepo authorRepository, poetRepo poetRepository) *PoemService {
	return &PoemService{
		poemRepo:    poemRepo,
		dynastyRepo: dynastyRepo,
		authorRepo:  authorRepo,
		poetRepo:    poetRepo,
	}
}

func (s *PoemService) GetPoemList(req *dto.PoemListRequest) ([]dto.PoemResponse, int64, error) {
	poems, total, err := s.poemRepo.List(req.Page, req.PageSize, req.Keyword, req.Dynasty, req.Author, req.Genre)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.PoemResponse, 0, len(poems))
	for _, poem := range poems {
		responses = append(responses, *s.toPoemResponse(&poem))
	}
	return responses, total, nil
}

func (s *PoemService) GetPoemByID(id uint64) (*dto.PoemResponse, error) {
	poem, err := s.poemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("poem not found")
	}

	_ = s.poemRepo.IncrementViews(id)
	return s.toPoemResponse(poem), nil
}

func (s *PoemService) CreatePoem(req *dto.PoemCreateRequest) (*dto.PoemResponse, error) {
	dynastyID, err := s.resolveDynastyID(uint64(req.DynastyID), req.DynastyName)
	if err != nil {
		return nil, err
	}

	authorID, err := s.resolveAuthorID(uint64(req.AuthorID), req.AuthorName, dynastyID)
	if err != nil {
		return nil, err
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

	if err := s.poemRepo.Create(poem); err != nil {
		return nil, err
	}
	return s.GetPoemByID(poem.ID)
}

func (s *PoemService) UpdatePoem(id uint64, req *dto.PoemUpdateRequest) (*dto.PoemResponse, error) {
	poem, err := s.poemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("poem not found")
	}

	if req.Title != "" {
		poem.Title = req.Title
	}
	if len(req.Content) > 0 {
		poem.Content = req.Content
	}
	if req.AuthorID > 0 {
		poem.AuthorID = uint64(req.AuthorID)
	}
	if req.DynastyID > 0 {
		poem.DynastyID = uint64(req.DynastyID)
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

	if err := s.poemRepo.Update(poem); err != nil {
		return nil, err
	}
	return s.toPoemResponse(poem), nil
}

func (s *PoemService) DeletePoem(id uint64) error {
	return s.poemRepo.Delete(id)
}

func (s *PoemService) BatchDeletePoems(ids []uint64) error {
	if err := validateBatchDeleteIDs(ids); err != nil {
		return err
	}
	return s.poemRepo.BatchDelete(ids)
}

func (s *PoemService) LikePoem(id uint64) error {
	return s.poemRepo.IncrementLikes(id)
}

func (s *PoemService) GetRandomPoems(limit int) ([]dto.PoemResponse, error) {
	poems, err := s.poemRepo.GetRandom(limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PoemResponse, 0, len(poems))
	for _, poem := range poems {
		responses = append(responses, *s.toPoemResponse(&poem))
	}
	return responses, nil
}

func (s *PoemService) GetDynastyList() ([]dto.DynastyResponse, error) {
	dynasties, err := s.dynastyRepo.List()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DynastyResponse, 0, len(dynasties))
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

func (s *PoemService) GetPoetList(keyword string, page, pageSize int) ([]dto.PoetResponse, int64, error) {
	poets, total, err := s.poetRepo.List(keyword, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.PoetResponse, 0, len(poets))
	for _, poet := range poets {
		responses = append(responses, s.toPoetResponse(&poet))
	}
	return responses, total, nil
}

func (s *PoemService) GetGenreList() ([]string, error) {
	return s.poemRepo.DistinctGenres()
}

func (s *PoemService) CreateDynasty(req *dto.DynastyCreateRequest) (*dto.DynastyResponse, error) {
	dynasty := &models.Dynasty{
		Name:        req.Name,
		Period:      req.Period,
		Description: req.Description,
	}
	if err := s.dynastyRepo.Create(dynasty); err != nil {
		return nil, err
	}
	return &dto.DynastyResponse{
		ID:          dynasty.ID,
		Name:        dynasty.Name,
		Period:      dynasty.Period,
		Description: dynasty.Description,
	}, nil
}

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
	if err := s.dynastyRepo.Update(dynasty); err != nil {
		return nil, err
	}
	return &dto.DynastyResponse{
		ID:          dynasty.ID,
		Name:        dynasty.Name,
		Period:      dynasty.Period,
		Description: dynasty.Description,
	}, nil
}

func (s *PoemService) DeleteDynasty(id uint64) error {
	return s.dynastyRepo.Delete(id)
}

func (s *PoemService) BatchDeleteDynasties(ids []uint64) error {
	if err := validateBatchDeleteIDs(ids); err != nil {
		return err
	}
	return s.dynastyRepo.BatchDelete(ids)
}

func (s *PoemService) CreatePoet(req *dto.PoetCreateRequest) (*dto.PoetResponse, error) {
	author := &models.Author{
		Name:      req.Name,
		Biography: req.Biography,
		Avatar:    req.Avatar,
	}
	if err := s.authorRepo.Create(author); err != nil {
		return nil, err
	}

	poet := &models.Poet{
		AuthorID:  author.ID,
		DynastyID: uint64(req.DynastyID),
		BirthYear: req.BirthYear,
		DeathYear: req.DeathYear,
	}
	if err := s.poetRepo.Create(poet); err != nil {
		return nil, err
	}
	poet.Author = *author
	return responsePtr(s.toPoetResponse(poet)), nil
}

func (s *PoemService) UpdatePoet(id uint64, req *dto.PoetUpdateRequest) (*dto.PoetResponse, error) {
	poet, err := s.poetRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("poet not found")
	}

	author := &poet.Author
	if author.ID == 0 {
		loadedAuthor, err := s.authorRepo.GetByID(poet.AuthorID)
		if err != nil {
			return nil, errors.New("author not found")
		}
		author = loadedAuthor
	}

	if req.Name != "" {
		author.Name = req.Name
	}
	if req.DynastyID > 0 {
		poet.DynastyID = uint64(req.DynastyID)
	}
	if req.Biography != "" {
		author.Biography = req.Biography
	}
	if req.Avatar != "" {
		author.Avatar = req.Avatar
	}
	if req.BirthYear > 0 {
		poet.BirthYear = req.BirthYear
	}
	if req.DeathYear > 0 {
		poet.DeathYear = req.DeathYear
	}

	if err := s.authorRepo.Update(author); err != nil {
		return nil, err
	}
	if err := s.poetRepo.Update(poet); err != nil {
		return nil, err
	}
	poet.Author = *author
	return responsePtr(s.toPoetResponse(poet)), nil
}

func (s *PoemService) DeletePoet(id uint64) error {
	return s.poetRepo.Delete(id)
}

func (s *PoemService) BatchDeletePoets(ids []uint64) error {
	if err := validateBatchDeleteIDs(ids); err != nil {
		return err
	}
	return s.poetRepo.BatchDelete(ids)
}

func (s *PoemService) resolveDynastyID(reqDynastyID uint64, dynastyName string) (uint64, error) {
	dynastyID := reqDynastyID
	dynastyResolved := false
	if dynastyName != "" {
		dynasty, err := s.dynastyRepo.GetByName(dynastyName)
		if err != nil {
			dynasty = &models.Dynasty{Name: dynastyName}
			if createErr := s.dynastyRepo.Create(dynasty); createErr != nil {
				return 0, createErr
			}
		}
		dynastyID = dynasty.ID
		dynastyResolved = true
	}

	if dynastyID > 0 && !dynastyResolved {
		if _, err := s.dynastyRepo.GetByID(dynastyID); err != nil {
			return 0, errors.New("dynasty not found")
		}
	}
	return dynastyID, nil
}

func (s *PoemService) resolveAuthorID(reqAuthorID uint64, authorName string, dynastyID uint64) (uint64, error) {
	if reqAuthorID > 0 {
		if _, err := s.authorRepo.GetByID(reqAuthorID); err != nil {
			return 0, errors.New("author not found")
		}
		return reqAuthorID, nil
	}
	if authorName == "" {
		return 0, nil
	}

	author, err := s.authorRepo.GetByName(authorName)
	if err != nil {
		author = &models.Author{Name: authorName}
		if createErr := s.authorRepo.Create(author); createErr != nil {
			return 0, createErr
		}
	}

	if _, err := s.poetRepo.GetByAuthorID(author.ID); err != nil {
		poet := &models.Poet{AuthorID: author.ID, DynastyID: dynastyID}
		if createErr := s.poetRepo.Create(poet); createErr != nil {
			return 0, createErr
		}
	}
	return author.ID, nil
}

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
		resp.Author = dto.AuthorResponse{
			ID:        poem.Author.ID,
			Name:      poem.Author.Name,
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

func (s *PoemService) toPoetResponse(poet *models.Poet) dto.PoetResponse {
	return dto.PoetResponse{
		ID:        poet.ID,
		AuthorID:  poet.AuthorID,
		Name:      poet.Author.Name,
		DynastyID: poet.DynastyID,
		Biography: poet.Author.Biography,
		Avatar:    poet.Author.Avatar,
		BirthYear: poet.BirthYear,
		DeathYear: poet.DeathYear,
	}
}

func responsePtr[T any](value T) *T {
	return &value
}

func validateBatchDeleteIDs(ids []uint64) error {
	if len(ids) == 0 {
		return errors.New("ids cannot be empty")
	}
	for _, id := range ids {
		if id == 0 {
			return errors.New("ids must be greater than 0")
		}
	}
	return nil
}
