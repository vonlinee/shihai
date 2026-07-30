package services

import (
	"errors"
	"sort"
	"strings"
	"unicode/utf8"

	"shihai/internal/dto"
	"shihai/internal/models"
	"shihai/internal/poetry"
)

type PoemService struct {
	poemRepo           poemRepository
	dynastyRepo        dynastyRepository
	authorRepo         authorRepository
	poetRepo           poetRepository
	poemAnnotationRepo poemAnnotationReader
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
	ListPoemTypes() ([]models.PoemType, error)
	CreatePoemType(poemType *models.PoemType) error
	GetPoemTypeByID(id uint64) (*models.PoemType, error)
	UpdatePoemType(poemType *models.PoemType) error
	DeletePoemType(id uint64) error
	BatchDeletePoemTypes(ids []uint64) error
	ListCiTunes() ([]models.CiTune, error)
	CreateCiTune(ciTune *models.CiTune) error
	GetCiTuneByID(id uint64) (*models.CiTune, error)
	UpdateCiTune(ciTune *models.CiTune) error
	DeleteCiTune(id uint64) error
	BatchDeleteCiTunes(ids []uint64) error
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
	List(keyword string, dynastyID uint64, page, pageSize int) ([]models.Poet, int64, error)
	Update(poet *models.Poet) error
	Delete(id uint64) error
	BatchDelete(ids []uint64) error
}

type poemAnnotationReader interface {
	ListByPoemID(poemID uint64) ([]models.PoemAnnotation, error)
	GetByID(id uint64) (*models.PoemAnnotation, error)
	Create(annotation *models.PoemAnnotation) error
	Update(annotation *models.PoemAnnotation) error
	Delete(id uint64) error
}

func NewPoemService(poemRepo poemRepository, dynastyRepo dynastyRepository, authorRepo authorRepository, poetRepo poetRepository) *PoemService {
	return &PoemService{
		poemRepo:    poemRepo,
		dynastyRepo: dynastyRepo,
		authorRepo:  authorRepo,
		poetRepo:    poetRepo,
	}
}

func (s *PoemService) SetPoemAnnotationRepository(annotationRepo poemAnnotationReader) {
	s.poemAnnotationRepo = annotationRepo
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
	resp := s.toPoemResponse(poem)
	if s.poemAnnotationRepo != nil {
		annotations, err := s.poemAnnotationRepo.ListByPoemID(id)
		if err != nil {
			return nil, err
		}
		resp.Annotations = toPoemAnnotationResponses(annotations)
	}
	return resp, nil
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

	pingze, err := poetry.PreparePingzeLinesForSave(req.Content, req.Pingze)
	if err != nil {
		return nil, err
	}

	genreCategory := s.resolveGenreCategory(req.GenreCategory, req.Genre)
	ciTuneID, err := s.resolveCiTuneID(uint64(req.CiTuneID), req.Title, genreCategory)
	if err != nil {
		return nil, err
	}
	poem := &models.Poem{
		Title:         req.Title,
		Content:       req.Content,
		Pingze:        pingze,
		AuthorID:      authorID,
		DynastyID:     dynastyID,
		GenreCategory: genreCategory,
		Genre:         req.Genre,
		CiTuneID:      ciTuneID,
		Translation:   req.Translation,
		Appreciation:  req.Appreciation,
		Annotation:    req.Annotation,
		AudioURL:      req.AudioURL,
		CoverImage:    req.CoverImage,
	}

	if err := s.poemRepo.Create(poem); err != nil {
		return nil, err
	}
	if req.Annotations != nil {
		if err := s.syncPoemAnnotations(poem, *req.Annotations); err != nil {
			return nil, err
		}
	}
	return s.toPoemResponseWithAnnotations(poem)
}

func (s *PoemService) UpdatePoem(id uint64, req *dto.PoemUpdateRequest) (*dto.PoemResponse, error) {
	poem, err := s.poemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("poem not found")
	}

	contentChanged := false
	if req.Title != "" {
		poem.Title = req.Title
	}
	if len(req.Content) > 0 {
		poem.Content = req.Content
		contentChanged = true
	}
	if req.Pingze != nil {
		pingze, err := poetry.PreparePingzeLinesForSave(poem.Content, *req.Pingze)
		if err != nil {
			return nil, err
		}
		poem.Pingze = pingze
	} else if contentChanged {
		poem.Pingze = poetry.RecognizePingzeLines(poem.Content)
	}
	if req.AuthorID > 0 {
		poem.AuthorID = uint64(req.AuthorID)
	}
	if req.DynastyID > 0 {
		poem.DynastyID = uint64(req.DynastyID)
	}
	if req.GenreCategory != "" {
		poem.GenreCategory = req.GenreCategory
	} else if req.Genre != "" || (poem.GenreCategory == "" && poem.Genre != "") {
		poem.GenreCategory = s.resolveGenreCategory("", firstNonEmpty(req.Genre, poem.Genre))
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
	if isCiCategory(poem.GenreCategory) {
		if req.CiTuneID != nil {
			ciTuneID, err := s.resolveCiTuneID(uint64(*req.CiTuneID), poem.Title, poem.GenreCategory)
			if err != nil {
				return nil, err
			}
			poem.CiTuneID = ciTuneID
			poem.CiTune = nil
		} else if poem.CiTuneID == nil {
			ciTuneID, err := s.resolveCiTuneID(0, poem.Title, poem.GenreCategory)
			if err != nil {
				return nil, err
			}
			poem.CiTuneID = ciTuneID
			poem.CiTune = nil
		}
	} else {
		poem.CiTuneID = nil
		poem.CiTune = nil
	}

	if err := s.poemRepo.Update(poem); err != nil {
		return nil, err
	}
	if req.Annotations != nil {
		if err := s.syncPoemAnnotations(poem, *req.Annotations); err != nil {
			return nil, err
		}
	}
	return s.toPoemResponseWithAnnotations(poem)
}

func (s *PoemService) syncPoemAnnotations(poem *models.Poem, requests []dto.PoemAnnotationUpsertRequest) error {
	if s.poemAnnotationRepo == nil {
		return errors.New("poem annotation repository not configured")
	}

	existing, err := s.poemAnnotationRepo.ListByPoemID(poem.ID)
	if err != nil {
		return err
	}
	existingByID := make(map[uint64]models.PoemAnnotation, len(existing))
	for _, annotation := range existing {
		existingByID[annotation.ID] = annotation
	}

	desired := make([]models.PoemAnnotation, 0, len(requests))
	seenIDs := make(map[uint64]struct{}, len(requests))
	for _, req := range requests {
		annotation := models.PoemAnnotation{
			PoemID:       poem.ID,
			TargetField:  normalizeAnnotationTarget(req.TargetField),
			StartLine:    req.StartLine,
			StartOffset:  req.StartOffset,
			EndLine:      req.EndLine,
			EndOffset:    req.EndOffset,
			SelectedText: req.SelectedText,
			Title:        strings.TrimSpace(req.Title),
			Content:      strings.TrimSpace(req.Content),
			Type:         normalizeAnnotationType(req.Type),
			DisplayOrder: req.DisplayOrder,
		}
		if req.ID > 0 {
			id := uint64(req.ID)
			current, ok := existingByID[id]
			if !ok || current.PoemID != poem.ID {
				return errPoemAnnotationNotFound
			}
			annotation.BaseModel = current.BaseModel
			seenIDs[id] = struct{}{}
		}
		desired = append(desired, annotation)
	}
	if err := validatePoemAnnotationSet(poem.Content, desired); err != nil {
		return err
	}

	for i := range desired {
		if desired[i].ID > 0 {
			if err := s.poemAnnotationRepo.Update(&desired[i]); err != nil {
				return err
			}
			continue
		}
		if err := s.poemAnnotationRepo.Create(&desired[i]); err != nil {
			return err
		}
	}
	for _, annotation := range existing {
		if _, ok := seenIDs[annotation.ID]; !ok {
			if err := s.poemAnnotationRepo.Delete(annotation.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func validatePoemAnnotationSet(content []string, annotations []models.PoemAnnotation) error {
	validated := make([]models.PoemAnnotation, 0, len(annotations))
	for i := range annotations {
		if err := validatePoemAnnotation(content, &annotations[i], annotations[i].ID, validated); err != nil {
			return err
		}
		validated = append(validated, annotations[i])
	}
	return nil
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

func (s *PoemService) GetPoetList(keyword string, dynastyID uint64, page, pageSize int) ([]dto.PoetResponse, int64, error) {
	poets, total, err := s.poetRepo.List(keyword, dynastyID, page, pageSize)
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

func (s *PoemService) resolveGenreCategory(genreCategory, genre string) string {
	genreCategory = strings.TrimSpace(genreCategory)
	if genreCategory != "" {
		return genreCategory
	}
	genre = strings.TrimSpace(genre)
	if genre == "" {
		return ""
	}
	poemTypes, err := s.poemRepo.ListPoemTypes()
	if err != nil {
		return ""
	}
	for _, poemType := range poemTypes {
		if poemType.Name == genre {
			category := strings.TrimSpace(poemType.Category)
			if category == "" {
				return "其他"
			}
			return category
		}
	}
	return "其他"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// GetGenreCategories returns poem genres grouped by top-level category.
func (s *PoemService) GetGenreCategories() ([]dto.GenreCategoryResponse, error) {
	poemTypes, err := s.poemRepo.ListPoemTypes()
	if err != nil {
		return nil, err
	}
	return toGenreCategoryResponses(poemTypes), nil
}

func toGenreCategoryResponses(poemTypes []models.PoemType) []dto.GenreCategoryResponse {
	categories := make([]dto.GenreCategoryResponse, 0)
	categoryIndexes := make(map[string]int)
	for _, poemType := range poemTypes {
		categoryName := strings.TrimSpace(poemType.Category)
		if categoryName == "" {
			categoryName = "其他"
		}
		categoryIndex, ok := categoryIndexes[categoryName]
		if !ok {
			categoryIndex = len(categories)
			categoryIndexes[categoryName] = categoryIndex
			categories = append(categories, dto.GenreCategoryResponse{Name: categoryName})
		}
		categories[categoryIndex].Genres = append(categories[categoryIndex].Genres, dto.GenreResponse{
			Name:         poemType.Name,
			Lines:        poemType.Lines,
			CharsPerLine: poemType.CharsPerLine,
			Description:  poemType.Description,
		})
	}
	return categories
}

// GetPoemTypes returns poem type reference data for admin management.
func (s *PoemService) GetPoemTypes() ([]dto.PoemTypeResponse, error) {
	poemTypes, err := s.poemRepo.ListPoemTypes()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.PoemTypeResponse, 0, len(poemTypes))
	for _, poemType := range poemTypes {
		responses = append(responses, toPoemTypeResponse(&poemType))
	}
	return responses, nil
}

// CreatePoemType creates poem type reference data.
func (s *PoemService) CreatePoemType(req *dto.PoemTypeCreateRequest) (*dto.PoemTypeResponse, error) {
	poemType := &models.PoemType{
		Name:         strings.TrimSpace(req.Name),
		Category:     strings.TrimSpace(req.Category),
		Lines:        req.Lines,
		CharsPerLine: req.CharsPerLine,
		Description:  req.Description,
	}
	if poemType.Name == "" {
		return nil, errors.New("poem type name is required")
	}
	if poemType.Category == "" {
		return nil, errors.New("poem type category is required")
	}
	if err := s.poemRepo.CreatePoemType(poemType); err != nil {
		return nil, err
	}
	return responsePtr(toPoemTypeResponse(poemType)), nil
}

// UpdatePoemType updates poem type reference data.
func (s *PoemService) UpdatePoemType(id uint64, req *dto.PoemTypeUpdateRequest) (*dto.PoemTypeResponse, error) {
	poemType, err := s.poemRepo.GetPoemTypeByID(id)
	if err != nil {
		return nil, errors.New("poem type not found")
	}
	poemType.Name = strings.TrimSpace(req.Name)
	poemType.Category = strings.TrimSpace(req.Category)
	poemType.Lines = req.Lines
	poemType.CharsPerLine = req.CharsPerLine
	poemType.Description = req.Description
	if poemType.Name == "" {
		return nil, errors.New("poem type name is required")
	}
	if poemType.Category == "" {
		return nil, errors.New("poem type category is required")
	}
	if err := s.poemRepo.UpdatePoemType(poemType); err != nil {
		return nil, err
	}
	return responsePtr(toPoemTypeResponse(poemType)), nil
}

// DeletePoemType deletes poem type reference data.
func (s *PoemService) DeletePoemType(id uint64) error {
	return s.poemRepo.DeletePoemType(id)
}

// BatchDeletePoemTypes deletes poem type reference data by IDs.
func (s *PoemService) BatchDeletePoemTypes(ids []uint64) error {
	if err := validateBatchDeleteIDs(ids); err != nil {
		return err
	}
	return s.poemRepo.BatchDeletePoemTypes(ids)
}

func toPoemTypeResponse(poemType *models.PoemType) dto.PoemTypeResponse {
	return dto.PoemTypeResponse{
		ID:           poemType.ID,
		Name:         poemType.Name,
		Category:     poemType.Category,
		Lines:        poemType.Lines,
		CharsPerLine: poemType.CharsPerLine,
		Description:  poemType.Description,
		CreatedAt:    poemType.CreatedAt,
		UpdatedAt:    poemType.UpdatedAt,
	}
}

// GetCiTunes returns ci tune reference data for admin management.
func (s *PoemService) GetCiTunes() ([]dto.CiTuneResponse, error) {
	ciTunes, err := s.poemRepo.ListCiTunes()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.CiTuneResponse, 0, len(ciTunes))
	for _, ciTune := range ciTunes {
		responses = append(responses, toCiTuneResponse(&ciTune))
	}
	return responses, nil
}

// CreateCiTune creates ci tune reference data.
func (s *PoemService) CreateCiTune(req *dto.CiTuneCreateRequest) (*dto.CiTuneResponse, error) {
	poemTypeID, err := s.resolveCiTunePoemTypeID(uint64(req.PoemTypeID))
	if err != nil {
		return nil, err
	}
	ciTune := &models.CiTune{
		Name:        strings.TrimSpace(req.Name),
		Aliases:     normalizeCiTuneAliases(req.Aliases),
		PoemTypeID:  poemTypeID,
		Description: req.Description,
	}
	if ciTune.Name == "" {
		return nil, errors.New("ci tune name is required")
	}
	if err := s.poemRepo.CreateCiTune(ciTune); err != nil {
		return nil, err
	}
	if ciTune.PoemTypeID != nil {
		if poemType, err := s.poemRepo.GetPoemTypeByID(*ciTune.PoemTypeID); err == nil {
			ciTune.PoemType = poemType
		}
	}
	return responsePtr(toCiTuneResponse(ciTune)), nil
}

// UpdateCiTune updates ci tune reference data.
func (s *PoemService) UpdateCiTune(id uint64, req *dto.CiTuneUpdateRequest) (*dto.CiTuneResponse, error) {
	ciTune, err := s.poemRepo.GetCiTuneByID(id)
	if err != nil {
		return nil, errors.New("ci tune not found")
	}
	poemTypeID, err := s.resolveCiTunePoemTypeID(uint64(req.PoemTypeID))
	if err != nil {
		return nil, err
	}
	ciTune.Name = strings.TrimSpace(req.Name)
	ciTune.Aliases = normalizeCiTuneAliases(req.Aliases)
	ciTune.PoemTypeID = poemTypeID
	ciTune.PoemType = nil
	ciTune.Description = req.Description
	if ciTune.Name == "" {
		return nil, errors.New("ci tune name is required")
	}
	if err := s.poemRepo.UpdateCiTune(ciTune); err != nil {
		return nil, err
	}
	if ciTune.PoemTypeID != nil {
		if poemType, err := s.poemRepo.GetPoemTypeByID(*ciTune.PoemTypeID); err == nil {
			ciTune.PoemType = poemType
		}
	}
	return responsePtr(toCiTuneResponse(ciTune)), nil
}

// DeleteCiTune deletes ci tune reference data.
func (s *PoemService) DeleteCiTune(id uint64) error {
	return s.poemRepo.DeleteCiTune(id)
}

// BatchDeleteCiTunes deletes ci tune reference data by IDs.
func (s *PoemService) BatchDeleteCiTunes(ids []uint64) error {
	if err := validateBatchDeleteIDs(ids); err != nil {
		return err
	}
	return s.poemRepo.BatchDeleteCiTunes(ids)
}

// ParseCiTuneTitle parses a poem title and returns the matched ci tune.
func (s *PoemService) ParseCiTuneTitle(title string) (*dto.CiTuneTitleParseResponse, error) {
	ciTunes, err := s.poemRepo.ListCiTunes()
	if err != nil {
		return nil, err
	}
	ciTune, matchedName := findCiTuneByTitle(title, ciTunes)
	if ciTune == nil {
		return &dto.CiTuneTitleParseResponse{}, nil
	}
	response := toCiTuneResponse(ciTune)
	return &dto.CiTuneTitleParseResponse{
		CiTune:      &response,
		MatchedName: matchedName,
	}, nil
}

func (s *PoemService) resolveCiTuneID(reqCiTuneID uint64, title, genreCategory string) (*uint64, error) {
	if !isCiCategory(genreCategory) {
		return nil, nil
	}
	if reqCiTuneID > 0 {
		if _, err := s.poemRepo.GetCiTuneByID(reqCiTuneID); err != nil {
			return nil, errors.New("ci tune not found")
		}
		return &reqCiTuneID, nil
	}
	ciTunes, err := s.poemRepo.ListCiTunes()
	if err != nil {
		return nil, err
	}
	ciTune, _ := findCiTuneByTitle(title, ciTunes)
	if ciTune == nil {
		return nil, nil
	}
	return &ciTune.ID, nil
}

func (s *PoemService) resolveCiTunePoemTypeID(reqPoemTypeID uint64) (*uint64, error) {
	if reqPoemTypeID == 0 {
		return nil, nil
	}
	poemType, err := s.poemRepo.GetPoemTypeByID(reqPoemTypeID)
	if err != nil {
		return nil, errors.New("poem type not found")
	}
	if !isCiCategory(poemType.Category) {
		return nil, errors.New("ci tune poem type must belong to ci category")
	}
	return &reqPoemTypeID, nil
}

func toCiTuneResponse(ciTune *models.CiTune) dto.CiTuneResponse {
	aliases := append([]string(nil), ciTune.Aliases...)
	if aliases == nil {
		aliases = []string{}
	}
	response := dto.CiTuneResponse{
		ID:          ciTune.ID,
		Name:        ciTune.Name,
		Aliases:     aliases,
		PoemTypeID:  ciTune.PoemTypeID,
		Description: ciTune.Description,
		CreatedAt:   ciTune.CreatedAt,
		UpdatedAt:   ciTune.UpdatedAt,
	}
	if ciTune.PoemType != nil && ciTune.PoemType.ID > 0 {
		poemType := toPoemTypeResponse(ciTune.PoemType)
		response.PoemType = &poemType
	}
	return response
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
	var dynasty *models.Dynasty
	if req.DynastyID > 0 {
		loadedDynasty, err := s.dynastyRepo.GetByID(uint64(req.DynastyID))
		if err != nil {
			return nil, errors.New("dynasty not found")
		}
		dynasty = loadedDynasty
	}

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
	if dynasty != nil {
		poet.Dynasty = *dynasty
	}
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
		dynasty, err := s.dynastyRepo.GetByID(uint64(req.DynastyID))
		if err != nil {
			return nil, errors.New("dynasty not found")
		}
		poet.DynastyID = dynasty.ID
		poet.Dynasty = *dynasty
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
		ID:            poem.ID,
		Title:         poem.Title,
		Content:       poem.Content,
		Pingze:        poetry.AlignPingzeLines(poem.Content, poem.Pingze),
		AuthorID:      poem.AuthorID,
		DynastyID:     poem.DynastyID,
		GenreCategory: poem.GenreCategory,
		Genre:         poem.Genre,
		CiTuneID:      poem.CiTuneID,
		Translation:   poem.Translation,
		Appreciation:  poem.Appreciation,
		Annotation:    poem.Annotation,
		AudioURL:      poem.AudioURL,
		CoverImage:    poem.CoverImage,
		Views:         poem.Views,
		Likes:         poem.Likes,
		Dislikes:      poem.Dislikes,
		Favorites:     poem.Favorites,
		CreatedAt:     poem.CreatedAt,
		UpdatedAt:     poem.UpdatedAt,
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
	if poem.CiTune != nil && poem.CiTune.ID > 0 {
		ciTune := toCiTuneResponse(poem.CiTune)
		resp.CiTune = &ciTune
	}
	return resp
}

func (s *PoemService) toPoemResponseWithAnnotations(poem *models.Poem) (*dto.PoemResponse, error) {
	resp := s.toPoemResponse(poem)
	if s.poemAnnotationRepo == nil {
		return resp, nil
	}
	annotations, err := s.poemAnnotationRepo.ListByPoemID(poem.ID)
	if err != nil {
		return nil, err
	}
	resp.Annotations = toPoemAnnotationResponses(annotations)
	return resp, nil
}

func (s *PoemService) toPoetResponse(poet *models.Poet) dto.PoetResponse {
	resp := dto.PoetResponse{
		ID:        poet.ID,
		AuthorID:  poet.AuthorID,
		Name:      poet.Author.Name,
		DynastyID: poet.DynastyID,
		Biography: poet.Author.Biography,
		Avatar:    poet.Author.Avatar,
		BirthYear: poet.BirthYear,
		DeathYear: poet.DeathYear,
	}
	if poet.Dynasty.ID > 0 {
		resp.Dynasty = dto.DynastyResponse{
			ID:          poet.Dynasty.ID,
			Name:        poet.Dynasty.Name,
			Period:      poet.Dynasty.Period,
			Description: poet.Dynasty.Description,
		}
	}
	return resp
}

func responsePtr[T any](value T) *T {
	return &value
}

func isCiCategory(category string) bool {
	return strings.TrimSpace(category) == "词"
}

func normalizeCiTuneAliases(aliases []string) []string {
	normalized := make([]string, 0, len(aliases))
	seen := make(map[string]struct{}, len(aliases))
	for _, alias := range aliases {
		alias = strings.TrimSpace(alias)
		if alias == "" {
			continue
		}
		if _, ok := seen[alias]; ok {
			continue
		}
		seen[alias] = struct{}{}
		normalized = append(normalized, alias)
	}
	return normalized
}

func findCiTuneByTitle(title string, ciTunes []models.CiTune) (*models.CiTune, string) {
	normalizedTitle := normalizeCiTuneTitle(title)
	if normalizedTitle == "" {
		return nil, ""
	}

	type candidate struct {
		ciTune *models.CiTune
		name   string
	}
	candidates := make([]candidate, 0, len(ciTunes))
	for i := range ciTunes {
		names := append([]string{ciTunes[i].Name}, ciTunes[i].Aliases...)
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			candidates = append(candidates, candidate{ciTune: &ciTunes[i], name: name})
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return utf8.RuneCountInString(candidates[i].name) > utf8.RuneCountInString(candidates[j].name)
	})
	for _, item := range candidates {
		if hasCiTuneTitlePrefix(normalizedTitle, item.name) {
			return item.ciTune, item.name
		}
	}
	return nil, ""
}

func normalizeCiTuneTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.Trim(title, "《》〈〉「」『』“”\"'")
	return strings.TrimSpace(title)
}

func hasCiTuneTitlePrefix(title, name string) bool {
	if !strings.HasPrefix(title, name) {
		return false
	}
	remainder := strings.TrimPrefix(title, name)
	if remainder == "" {
		return true
	}
	firstRune, _ := utf8.DecodeRuneInString(remainder)
	return strings.ContainsRune("·・ 　:：-—_，,。.", firstRune)
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
