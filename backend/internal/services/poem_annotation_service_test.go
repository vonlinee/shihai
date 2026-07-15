package services

import (
	"encoding/json"
	"strings"
	"testing"

	"shihai/internal/dto"
	"shihai/internal/models"
)

func TestPoemAnnotationService_CreateRejectsSelectedTextMismatch(t *testing.T) {
	service := NewPoemAnnotationService(
		&fakePoemAnnotationRepository{},
		&fakeAnnotationPoemRepository{
			poem: &models.Poem{
				BaseModel: models.BaseModel{ID: 1},
				Content:   []string{"床前明月光", "疑是地上霜"},
			},
		},
	)

	_, err := service.CreateAnnotation(1, &dto.PoemAnnotationCreateRequest{
		StartLine:    0,
		StartOffset:  0,
		EndLine:      0,
		EndOffset:    2,
		SelectedText: "明月",
		Content:      "注释内容",
	})

	if err == nil || !strings.Contains(err.Error(), "selected text does not match") {
		t.Fatalf("CreateAnnotation error = %v, want selected text mismatch", err)
	}
}

func TestPoemAnnotationService_CreateRejectsOverlappingRanges(t *testing.T) {
	annotationRepo := &fakePoemAnnotationRepository{
		annotations: []models.PoemAnnotation{
			{
				BaseModel:    models.BaseModel{ID: 11},
				PoemID:       1,
				TargetField:  "content",
				StartLine:    0,
				StartOffset:  1,
				EndLine:      0,
				EndOffset:    3,
				SelectedText: "前明",
				Content:      "已有注释",
			},
		},
	}
	service := NewPoemAnnotationService(
		annotationRepo,
		&fakeAnnotationPoemRepository{
			poem: &models.Poem{
				BaseModel: models.BaseModel{ID: 1},
				Content:   []string{"床前明月光"},
			},
		},
	)

	_, err := service.CreateAnnotation(1, &dto.PoemAnnotationCreateRequest{
		StartLine:    0,
		StartOffset:  2,
		EndLine:      0,
		EndOffset:    4,
		SelectedText: "明月",
		Content:      "新增注释",
	})

	if err == nil || !strings.Contains(err.Error(), "annotation ranges cannot overlap") {
		t.Fatalf("CreateAnnotation error = %v, want overlap error", err)
	}
	if annotationRepo.created != nil {
		t.Fatalf("created annotation = %#v, want nil", annotationRepo.created)
	}
}

func TestPoemAnnotationService_CreateTreatsHiddenStatusAsPersistentAnnotation(t *testing.T) {
	annotationRepo := &fakePoemAnnotationRepository{
		annotations: []models.PoemAnnotation{
			{
				BaseModel:    models.BaseModel{ID: 11},
				PoemID:       1,
				TargetField:  models.PoemAnnotationTargetContent,
				StartLine:    0,
				StartOffset:  1,
				EndLine:      0,
				EndOffset:    3,
				SelectedText: "bc",
				Content:      "existing annotation",
			},
		},
	}
	service := NewPoemAnnotationService(
		annotationRepo,
		&fakeAnnotationPoemRepository{
			poem: &models.Poem{
				BaseModel: models.BaseModel{ID: 1},
				Content:   []string{"abcdef"},
			},
		},
	)

	var req dto.PoemAnnotationCreateRequest
	if err := json.Unmarshal([]byte(`{
		"startLine": 0,
		"startOffset": 2,
		"endLine": 0,
		"endOffset": 4,
		"selectedText": "cd",
		"content": "new annotation",
		"status": "hidden"
	}`), &req); err != nil {
		t.Fatalf("Unmarshal request error = %v, want nil", err)
	}

	_, err := service.CreateAnnotation(1, &req)

	if err == nil || !strings.Contains(err.Error(), "annotation ranges cannot overlap") {
		t.Fatalf("CreateAnnotation error = %v, want overlap error", err)
	}
	if annotationRepo.created != nil {
		t.Fatalf("created annotation = %#v, want nil", annotationRepo.created)
	}
}

func TestPoemAnnotationService_ListSortsAndAssignsDisplayNumbers(t *testing.T) {
	service := NewPoemAnnotationService(
		&fakePoemAnnotationRepository{
			annotations: []models.PoemAnnotation{
				{
					BaseModel:    models.BaseModel{ID: 22},
					PoemID:       1,
					TargetField:  "content",
					StartLine:    1,
					StartOffset:  0,
					EndLine:      1,
					EndOffset:    1,
					SelectedText: "疑",
					Content:      "第二条",
				},
				{
					BaseModel:    models.BaseModel{ID: 11},
					PoemID:       1,
					TargetField:  "content",
					StartLine:    0,
					StartOffset:  2,
					EndLine:      0,
					EndOffset:    4,
					SelectedText: "明月",
					Content:      "第一条",
				},
			},
		},
		&fakeAnnotationPoemRepository{
			poem: &models.Poem{
				BaseModel: models.BaseModel{ID: 1},
				Content:   []string{"床前明月光", "疑是地上霜"},
			},
		},
	)

	annotations, err := service.ListAnnotations(1)

	if err != nil {
		t.Fatalf("ListAnnotations error = %v, want nil", err)
	}
	if len(annotations) != 2 {
		t.Fatalf("len(annotations) = %d, want 2", len(annotations))
	}
	if annotations[0].ID != 11 || annotations[0].DisplayNo != 1 {
		t.Fatalf("first annotation = %#v, want ID 11 displayNo 1", annotations[0])
	}
	if annotations[1].ID != 22 || annotations[1].DisplayNo != 2 {
		t.Fatalf("second annotation = %#v, want ID 22 displayNo 2", annotations[1])
	}
}

type fakeAnnotationPoemRepository struct {
	poem *models.Poem
}

func (r *fakeAnnotationPoemRepository) GetByID(id uint64) (*models.Poem, error) {
	if r.poem != nil {
		return r.poem, nil
	}
	return &models.Poem{BaseModel: models.BaseModel{ID: id}}, nil
}

type fakePoemAnnotationRepository struct {
	annotations []models.PoemAnnotation
	created     *models.PoemAnnotation
	updated     *models.PoemAnnotation
	deletedID   uint64
}

func (r *fakePoemAnnotationRepository) ListByPoemID(poemID uint64) ([]models.PoemAnnotation, error) {
	items := make([]models.PoemAnnotation, 0, len(r.annotations))
	for _, annotation := range r.annotations {
		if annotation.PoemID != poemID {
			continue
		}
		items = append(items, annotation)
	}
	return items, nil
}

func (r *fakePoemAnnotationRepository) GetByID(id uint64) (*models.PoemAnnotation, error) {
	for _, annotation := range r.annotations {
		if annotation.ID == id {
			copied := annotation
			return &copied, nil
		}
	}
	return nil, errPoemAnnotationNotFound
}

func (r *fakePoemAnnotationRepository) Create(annotation *models.PoemAnnotation) error {
	copied := *annotation
	copied.ID = 100
	r.created = &copied
	r.annotations = append(r.annotations, copied)
	annotation.ID = copied.ID
	return nil
}

func (r *fakePoemAnnotationRepository) Update(annotation *models.PoemAnnotation) error {
	copied := *annotation
	r.updated = &copied
	for i := range r.annotations {
		if r.annotations[i].ID == copied.ID {
			r.annotations[i] = copied
			return nil
		}
	}
	return nil
}

func (r *fakePoemAnnotationRepository) Delete(id uint64) error {
	r.deletedID = id
	for i := range r.annotations {
		if r.annotations[i].ID == id {
			r.annotations = append(r.annotations[:i], r.annotations[i+1:]...)
			return nil
		}
	}
	return nil
}
