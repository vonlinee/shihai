package dto

import "time"

type PoemAnnotationCreateRequest struct {
	TargetField  string `json:"targetField"`
	StartLine    int    `json:"startLine" binding:"min=0"`
	StartOffset  int    `json:"startOffset" binding:"min=0"`
	EndLine      int    `json:"endLine" binding:"min=0"`
	EndOffset    int    `json:"endOffset" binding:"min=0"`
	SelectedText string `json:"selectedText" binding:"required"`
	Title        string `json:"title" binding:"max=100"`
	Content      string `json:"content" binding:"required"`
	Type         string `json:"type"`
	DisplayOrder int    `json:"displayOrder"`
}

type PoemAnnotationUpdateRequest struct {
	TargetField  string `json:"targetField"`
	StartLine    int    `json:"startLine" binding:"min=0"`
	StartOffset  int    `json:"startOffset" binding:"min=0"`
	EndLine      int    `json:"endLine" binding:"min=0"`
	EndOffset    int    `json:"endOffset" binding:"min=0"`
	SelectedText string `json:"selectedText" binding:"required"`
	Title        string `json:"title" binding:"max=100"`
	Content      string `json:"content" binding:"required"`
	Type         string `json:"type"`
	DisplayOrder int    `json:"displayOrder"`
}

type PoemAnnotationUpsertRequest struct {
	ID           RequestID `json:"id"`
	TargetField  string    `json:"targetField"`
	StartLine    int       `json:"startLine" binding:"min=0"`
	StartOffset  int       `json:"startOffset" binding:"min=0"`
	EndLine      int       `json:"endLine" binding:"min=0"`
	EndOffset    int       `json:"endOffset" binding:"min=0"`
	SelectedText string    `json:"selectedText" binding:"required"`
	Title        string    `json:"title" binding:"max=100"`
	Content      string    `json:"content" binding:"required"`
	Type         string    `json:"type"`
	DisplayOrder int       `json:"displayOrder"`
}

type PoemAnnotationResponse struct {
	ID           uint64    `json:"id,string"`
	PoemID       uint64    `json:"poemId,string"`
	DisplayNo    int       `json:"displayNo"`
	TargetField  string    `json:"targetField"`
	StartLine    int       `json:"startLine"`
	StartOffset  int       `json:"startOffset"`
	EndLine      int       `json:"endLine"`
	EndOffset    int       `json:"endOffset"`
	SelectedText string    `json:"selectedText"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Type         string    `json:"type"`
	DisplayOrder int       `json:"displayOrder"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
