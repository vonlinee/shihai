import { type CSSProperties, type FormEvent, useEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { X } from 'lucide-react';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

import {
  sortAnnotationDraftsByRange,
  type PoemAnnotationDraft,
} from './PoemAnnotationManager';
import { type AnnotationRange, normalizeTextareaSelectionRange } from './annotationRange';

interface PoemContentAnnotationEditorProps {
  contentValue: string;
  annotations: PoemAnnotationDraft[];
  onContentChange: (content: string) => void;
  onAnnotationsChange: (annotations: PoemAnnotationDraft[]) => void;
  isLoadingAnnotations?: boolean;
}

interface AnnotationFormState {
  title: string;
  content: string;
}

interface FloatingFormPosition {
  top: number;
  left: number;
}

const emptyForm: AnnotationFormState = {
  title: '',
  content: '',
};

const floatingFormWidth = 320;
const estimatedFloatingFormHeight = 390;
const viewportPadding = 12;
const sentenceBoundaryPattern = /[，。！？；;,.!?]/u;

export function PoemContentAnnotationEditor({
  contentValue,
  annotations,
  onContentChange,
  onAnnotationsChange,
}: PoemContentAnnotationEditorProps) {
  const textareaRef = useRef<HTMLTextAreaElement | null>(null);
  const [selectionRange, setSelectionRange] = useState<AnnotationRange | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [floatingFormPosition, setFloatingFormPosition] = useState<FloatingFormPosition | null>(null);
  const [editingIndex, setEditingIndex] = useState<number | null>(null);
  const [form, setForm] = useState<AnnotationFormState>(emptyForm);

  useEffect(() => {
    setSelectionRange(null);
    setIsFormOpen(false);
    setFloatingFormPosition(null);
    setEditingIndex(null);
    setForm(emptyForm);
  }, [contentValue]);

  useEffect(() => {
    if (!isFormOpen) return;

    const handleViewportChange = () => {
      const textarea = textareaRef.current;
      if (!textarea || textarea.selectionStart === textarea.selectionEnd) return;
      updateFloatingFormPosition(textarea.selectionStart, textarea.selectionEnd);
    };

    window.addEventListener('resize', handleViewportChange);
    window.addEventListener('scroll', handleViewportChange, true);
    return () => {
      window.removeEventListener('resize', handleViewportChange);
      window.removeEventListener('scroll', handleViewportChange, true);
    };
  }, [isFormOpen]);

  const handleSelectionChange = () => {
    const textarea = textareaRef.current;
    if (!textarea) return;

    const normalized = normalizeTextareaSelectionRange(contentValue, textarea.selectionStart, textarea.selectionEnd);
    if (!normalized) {
      if (!isFormOpen) {
        setSelectionRange(null);
      }
      return;
    }
    if (!isSingleSentenceAnnotationRange(normalized)) {
      resetForm();
      return;
    }

    setSelectionRange(normalized);
    setEditingIndex(null);
    setIsFormOpen(true);
    setForm({
      title: normalized.selectedText,
      content: '',
    });
    updateFloatingFormPosition(textarea.selectionStart, textarea.selectionEnd);
  };

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    event.stopPropagation();
    if (!selectionRange || !form.content.trim()) {
      return;
    }
    if (editingIndex === null && hasDuplicateAnnotationRange(annotations, selectionRange)) {
      toast.warning('该文本已添加过标注');
      return;
    }
    const overlappingAnnotations = editingIndex === null ? findOverlappingAnnotationRanges(annotations, selectionRange) : [];
    if (overlappingAnnotations.length > 0) {
      toast.warning(`标注文本不能与已有标注交叉：${formatConflictingAnnotationTexts(overlappingAnnotations)}`);
      return;
    }

    const existingAnnotation = editingIndex === null ? null : annotations[editingIndex];
    const nextAnnotation: PoemAnnotationDraft = {
      id: existingAnnotation?.id,
      targetField: 'content',
      ...selectionRange,
      title: form.title.trim(),
      content: form.content.trim(),
      type: 'note',
      displayOrder: existingAnnotation?.displayOrder ?? 0,
    };

    const nextAnnotations = [...annotations];
    if (editingIndex === null) {
      nextAnnotations.push(nextAnnotation);
    } else {
      nextAnnotations[editingIndex] = nextAnnotation;
    }

    onAnnotationsChange(sortAnnotationDraftsByRange(nextAnnotations));
    resetForm();
  };

  const resetForm = () => {
    setSelectionRange(null);
    setIsFormOpen(false);
    setFloatingFormPosition(null);
    setEditingIndex(null);
    setForm(emptyForm);
  };

  const updateFloatingFormPosition = (selectionStart: number, selectionEnd: number) => {
    const textarea = textareaRef.current;
    if (!textarea) {
      setFloatingFormPosition(null);
      return;
    }

    const style = window.getComputedStyle(textarea);
    const startPoint = getTextareaTextPoint(textarea, selectionStart, style);
    const endPoint = getTextareaTextPoint(textarea, selectionEnd, style);
    const textareaRect = textarea.getBoundingClientRect();
    const lineHeight = getLineHeight(style);
    const anchorTop = textareaRect.top + Math.min(startPoint.top, endPoint.top);
    const selectionStartLeft = textareaRect.left + Math.min(startPoint.left, endPoint.left);
    const selectionEndLeft = textareaRect.left + Math.max(startPoint.left, endPoint.left);
    const maxLeft = Math.max(viewportPadding, window.innerWidth - floatingFormWidth - viewportPadding);
    const preferredLeft = selectionEndLeft + 12;
    const fallbackLeft = selectionStartLeft - floatingFormWidth - 12;
    const rawLeft = preferredLeft <= maxLeft ? preferredLeft : fallbackLeft;
    const maxTop = Math.max(viewportPadding, window.innerHeight - estimatedFloatingFormHeight - viewportPadding);
    const rawTop = anchorTop + lineHeight / 2 - estimatedFloatingFormHeight / 2;

    setFloatingFormPosition({
      top: Math.min(maxTop, Math.max(viewportPadding, rawTop)),
      left: Math.min(maxLeft, Math.max(viewportPadding, rawLeft)),
    });
  };

  const floatingFormStyle: CSSProperties | undefined = floatingFormPosition
    ? {
        top: floatingFormPosition.top,
        left: floatingFormPosition.left,
        width: floatingFormWidth,
      }
    : undefined;

  return (
    <div className="space-y-4 rounded-md border p-4">
      <div className="space-y-2">
        <label className="text-sm font-medium">内容 *</label>
        <div>
          <textarea
            ref={textareaRef}
            value={contentValue}
            onChange={(event) => onContentChange(event.target.value)}
            onSelect={handleSelectionChange}
            onMouseUp={handleSelectionChange}
            onKeyUp={handleSelectionChange}
            placeholder="诗词内容"
            required
            rows={6}
            className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
          />
          {isFormOpen && selectionRange && floatingFormStyle && createPortal(
            <form
              onSubmit={handleSubmit}
              style={floatingFormStyle}
              className="fixed z-[70] space-y-3 rounded-md border bg-popover p-4 text-popover-foreground shadow-lg"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <div className="flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1">
                    <h4 className="font-medium">{editingIndex !== null ? '编辑标注' : '新增标注'}</h4>
                    <span className="rounded bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
                      {formatAnnotationTitleMeta(selectionRange)}
                    </span>
                  </div>
                </div>
                <Button type="button" variant="ghost" size="sm" onClick={resetForm}>
                  <X className="h-4 w-4" />
                </Button>
              </div>
              <div className="grid grid-cols-[3rem_minmax(0,1fr)] items-center gap-2">
                <label className="text-sm font-medium">标题</label>
                <Input
                  value={form.title}
                  onChange={(event) => setForm({ ...form, title: event.target.value })}
                  placeholder="可选，默认使用被标注文本"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">内容 *</label>
                <textarea
                  value={form.content}
                  onChange={(event) => setForm({ ...form, content: event.target.value })}
                  rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  placeholder="填写标注解释"
                />
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={!form.content.trim()}>
                  标注
                </Button>
              </div>
            </form>,
            document.body,
          )}
        </div>
      </div>

    </div>
  );
}

function formatAnnotationTitleMeta(range: AnnotationRange): string {
  const startLine = range.startLine + 1;
  const endLine = range.endLine + 1;
  const startColumn = range.startOffset + 1;
  const endColumn = Math.max(1, range.endOffset);

  if (range.startLine === range.endLine) {
    return `第 ${startLine} 行 · ${startColumn}-${endColumn} 列`;
  }

  return `第 ${startLine}-${endLine} 行 · ${range.endLine - range.startLine + 1} 行`;
}

function isSingleSentenceAnnotationRange(range: AnnotationRange): boolean {
  return range.startLine === range.endLine && !sentenceBoundaryPattern.test(range.selectedText);
}

function hasDuplicateAnnotationRange(annotations: PoemAnnotationDraft[], range: AnnotationRange): boolean {
  return annotations.some((annotation) => (
    annotation.startLine === range.startLine
    && annotation.startOffset === range.startOffset
    && annotation.endLine === range.endLine
    && annotation.endOffset === range.endOffset
  ));
}

function findOverlappingAnnotationRanges(annotations: PoemAnnotationDraft[], range: AnnotationRange): PoemAnnotationDraft[] {
  return annotations.filter((annotation) => {
    if (annotation.startLine !== range.startLine || annotation.endLine !== range.endLine) {
      return false;
    }
    return annotation.startOffset < range.endOffset && range.startOffset < annotation.endOffset;
  });
}

function formatConflictingAnnotationTexts(annotations: PoemAnnotationDraft[]): string {
  const texts = annotations.map((annotation) => `“${annotation.selectedText}”`);
  if (texts.length <= 3) {
    return texts.join('、');
  }
  return `${texts.slice(0, 3).join('、')} 等 ${texts.length} 条`;
}

function getTextareaTextPoint(
  textarea: HTMLTextAreaElement,
  position: number,
  style: CSSStyleDeclaration,
): { top: number; left: number } {
  const mirror = document.createElement('div');
  const marker = document.createElement('span');
  const properties = [
    'box-sizing',
    'width',
    'border-top-width',
    'border-right-width',
    'border-bottom-width',
    'border-left-width',
    'padding-top',
    'padding-right',
    'padding-bottom',
    'padding-left',
    'font-family',
    'font-size',
    'font-weight',
    'font-style',
    'letter-spacing',
    'line-height',
    'text-transform',
    'text-indent',
    'text-align',
    'word-spacing',
  ];

  properties.forEach((property) => {
    mirror.style.setProperty(property, style.getPropertyValue(property));
  });
  mirror.style.position = 'absolute';
  mirror.style.left = '-9999px';
  mirror.style.top = '0';
  mirror.style.visibility = 'hidden';
  mirror.style.whiteSpace = 'pre-wrap';
  mirror.style.overflowWrap = 'break-word';
  mirror.style.wordBreak = 'break-word';

  mirror.textContent = textarea.value.slice(0, position);
  marker.textContent = textarea.value.slice(position, position + 1) || '.';
  mirror.appendChild(marker);
  document.body.appendChild(mirror);

  const markerRect = marker.getBoundingClientRect();
  const mirrorRect = mirror.getBoundingClientRect();
  const point = {
    top: markerRect.top - mirrorRect.top - textarea.scrollTop,
    left: markerRect.left - mirrorRect.left - textarea.scrollLeft,
  };

  document.body.removeChild(mirror);
  return point;
}

function getLineHeight(style: CSSStyleDeclaration): number {
  const parsedLineHeight = Number.parseFloat(style.lineHeight);
  if (Number.isFinite(parsedLineHeight)) {
    return parsedLineHeight;
  }
  const parsedFontSize = Number.parseFloat(style.fontSize);
  return Number.isFinite(parsedFontSize) ? parsedFontSize * 1.2 : 20;
}
