package services

import (
	"errors"
	"sort"
	"strings"

	"shihai/internal/dto"
	"shihai/internal/models"
)

var errPoemAnnotationNotFound = errors.New("poem annotation not found")

type PoemAnnotationService struct {
	annotationRepo poemAnnotationRepository
	poemRepo       annotationPoemRepository
}

type poemAnnotationRepository interface {
	ListByPoemID(poemID uint64) ([]models.PoemAnnotation, error)
	GetByID(id uint64) (*models.PoemAnnotation, error)
	Create(annotation *models.PoemAnnotation) error
	Update(annotation *models.PoemAnnotation) error
	Delete(id uint64) error
}

type annotationPoemRepository interface {
	GetByID(id uint64) (*models.Poem, error)
}

func NewPoemAnnotationService(annotationRepo poemAnnotationRepository, poemRepo annotationPoemRepository) *PoemAnnotationService {
	return &PoemAnnotationService{
		annotationRepo: annotationRepo,
		poemRepo:       poemRepo,
	}
}

func (s *PoemAnnotationService) ListAnnotations(poemID uint64) ([]dto.PoemAnnotationResponse, error) {
	annotations, err := s.annotationRepo.ListByPoemID(poemID)
	if err != nil {
		return nil, err
	}
	return toPoemAnnotationResponses(annotations), nil
}

func (s *PoemAnnotationService) CreateAnnotation(poemID uint64, req *dto.PoemAnnotationCreateRequest) (*dto.PoemAnnotationResponse, error) {
	poem, err := s.poemRepo.GetByID(poemID)
	if err != nil {
		return nil, errors.New("poem not found")
	}

	annotation := &models.PoemAnnotation{
		PoemID:       poemID,
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
	if err := validatePoemAnnotation(poem.Content, annotation, 0, nil); err != nil {
		return nil, err
	}

	existing, err := s.annotationRepo.ListByPoemID(poemID)
	if err != nil {
		return nil, err
	}
	if err := validatePoemAnnotation(poem.Content, annotation, 0, existing); err != nil {
		return nil, err
	}

	if err := s.annotationRepo.Create(annotation); err != nil {
		return nil, err
	}
	annotations, err := s.annotationRepo.ListByPoemID(poemID)
	if err != nil {
		return nil, err
	}
	responses := toPoemAnnotationResponses(annotations)
	for i := range responses {
		if responses[i].ID == annotation.ID {
			return &responses[i], nil
		}
	}
	response := toPoemAnnotationResponse(*annotation, 1)
	return &response, nil
}

func (s *PoemAnnotationService) UpdateAnnotation(id uint64, req *dto.PoemAnnotationUpdateRequest) (*dto.PoemAnnotationResponse, error) {
	annotation, err := s.annotationRepo.GetByID(id)
	if err != nil {
		return nil, errPoemAnnotationNotFound
	}
	poem, err := s.poemRepo.GetByID(annotation.PoemID)
	if err != nil {
		return nil, errors.New("poem not found")
	}

	annotation.TargetField = normalizeAnnotationTarget(req.TargetField)
	annotation.StartLine = req.StartLine
	annotation.StartOffset = req.StartOffset
	annotation.EndLine = req.EndLine
	annotation.EndOffset = req.EndOffset
	annotation.SelectedText = req.SelectedText
	annotation.Title = strings.TrimSpace(req.Title)
	annotation.Content = strings.TrimSpace(req.Content)
	annotation.Type = normalizeAnnotationType(req.Type)
	annotation.DisplayOrder = req.DisplayOrder

	existing, err := s.annotationRepo.ListByPoemID(annotation.PoemID)
	if err != nil {
		return nil, err
	}
	if err := validatePoemAnnotation(poem.Content, annotation, id, existing); err != nil {
		return nil, err
	}

	if err := s.annotationRepo.Update(annotation); err != nil {
		return nil, err
	}
	annotations, err := s.annotationRepo.ListByPoemID(annotation.PoemID)
	if err != nil {
		return nil, err
	}
	responses := toPoemAnnotationResponses(annotations)
	for i := range responses {
		if responses[i].ID == id {
			return &responses[i], nil
		}
	}
	response := toPoemAnnotationResponse(*annotation, 1)
	return &response, nil
}

func (s *PoemAnnotationService) DeleteAnnotation(id uint64) error {
	if _, err := s.annotationRepo.GetByID(id); err != nil {
		return errPoemAnnotationNotFound
	}
	return s.annotationRepo.Delete(id)
}

func validatePoemAnnotation(content []string, annotation *models.PoemAnnotation, selfID uint64, existing []models.PoemAnnotation) error {
	if annotation.TargetField != models.PoemAnnotationTargetContent {
		return errors.New("unsupported annotation target field")
	}
	if annotation.Content == "" {
		return errors.New("annotation content cannot be empty")
	}
	selectedText, err := extractSelectedText(content, annotation.StartLine, annotation.StartOffset, annotation.EndLine, annotation.EndOffset)
	if err != nil {
		return err
	}
	if selectedText != annotation.SelectedText {
		return errors.New("selected text does not match poem content")
	}
	for _, item := range existing {
		if item.ID == selfID || item.TargetField != annotation.TargetField {
			continue
		}
		if rangesOverlap(annotation, &item) {
			return errors.New("annotation ranges cannot overlap")
		}
	}
	return nil
}

func extractSelectedText(content []string, startLine, startOffset, endLine, endOffset int) (string, error) {
	if len(content) == 0 {
		return "", errors.New("poem content cannot be empty")
	}
	if startLine < 0 || endLine < 0 || startLine >= len(content) || endLine >= len(content) {
		return "", errors.New("annotation range line is out of bounds")
	}
	if compareRangePoint(startLine, startOffset, endLine, endOffset) >= 0 {
		return "", errors.New("annotation range must not be empty")
	}

	startRunes := []rune(content[startLine])
	endRunes := []rune(content[endLine])
	if startOffset < 0 || startOffset > len(startRunes) || endOffset < 0 || endOffset > len(endRunes) {
		return "", errors.New("annotation range offset is out of bounds")
	}
	if startLine == endLine {
		return string(startRunes[startOffset:endOffset]), nil
	}

	parts := make([]string, 0, endLine-startLine+1)
	parts = append(parts, string(startRunes[startOffset:]))
	for line := startLine + 1; line < endLine; line++ {
		parts = append(parts, content[line])
	}
	parts = append(parts, string(endRunes[:endOffset]))
	return strings.Join(parts, "\n"), nil
}

func rangesOverlap(a, b *models.PoemAnnotation) bool {
	return compareRangePoint(a.StartLine, a.StartOffset, b.EndLine, b.EndOffset) < 0 &&
		compareRangePoint(b.StartLine, b.StartOffset, a.EndLine, a.EndOffset) < 0
}

func compareRangePoint(leftLine, leftOffset, rightLine, rightOffset int) int {
	if leftLine < rightLine {
		return -1
	}
	if leftLine > rightLine {
		return 1
	}
	if leftOffset < rightOffset {
		return -1
	}
	if leftOffset > rightOffset {
		return 1
	}
	return 0
}

func normalizeAnnotationTarget(target string) string {
	if strings.TrimSpace(target) == "" {
		return models.PoemAnnotationTargetContent
	}
	return strings.TrimSpace(target)
}

func normalizeAnnotationType(annotationType string) string {
	if strings.TrimSpace(annotationType) == "" {
		return models.PoemAnnotationTypeNote
	}
	return strings.TrimSpace(annotationType)
}

func toPoemAnnotationResponses(annotations []models.PoemAnnotation) []dto.PoemAnnotationResponse {
	sort.SliceStable(annotations, func(i, j int) bool {
		left := annotations[i]
		right := annotations[j]
		if left.StartLine != right.StartLine {
			return left.StartLine < right.StartLine
		}
		if left.StartOffset != right.StartOffset {
			return left.StartOffset < right.StartOffset
		}
		if left.EndLine != right.EndLine {
			return left.EndLine < right.EndLine
		}
		if left.EndOffset != right.EndOffset {
			return left.EndOffset < right.EndOffset
		}
		return left.ID < right.ID
	})

	responses := make([]dto.PoemAnnotationResponse, 0, len(annotations))
	for index, annotation := range annotations {
		responses = append(responses, toPoemAnnotationResponse(annotation, index+1))
	}
	return responses
}

func toPoemAnnotationResponse(annotation models.PoemAnnotation, displayNo int) dto.PoemAnnotationResponse {
	return dto.PoemAnnotationResponse{
		ID:           annotation.ID,
		PoemID:       annotation.PoemID,
		DisplayNo:    displayNo,
		TargetField:  annotation.TargetField,
		StartLine:    annotation.StartLine,
		StartOffset:  annotation.StartOffset,
		EndLine:      annotation.EndLine,
		EndOffset:    annotation.EndOffset,
		SelectedText: annotation.SelectedText,
		Title:        annotation.Title,
		Content:      annotation.Content,
		Type:         annotation.Type,
		DisplayOrder: annotation.DisplayOrder,
		CreatedAt:    annotation.CreatedAt,
		UpdatedAt:    annotation.UpdatedAt,
	}
}
