import { type FormEvent, useEffect, useRef, useState } from 'react';
import { Edit2, Plus, Trash2, X } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { useConfirmDialog } from '@/components/ui/ConfirmDialog';
import {
  useAdminCreatePoemAnnotation,
  useAdminDeletePoemAnnotation,
  useAdminPoemAnnotations,
  useAdminUpdatePoemAnnotation,
} from '@/hooks/useAdmin';
import { cn } from '@/utils/cn';
import type { Poem, PoemAnnotation } from '@/types';

import { type AnnotationRange, normalizeSelectionRange } from './annotationRange';

export interface PoemAnnotationDraft {
  id?: string;
  targetField: 'content';
  startLine: number;
  startOffset: number;
  endLine: number;
  endOffset: number;
  selectedText: string;
  title?: string;
  content: string;
  type: 'note';
  displayOrder: number;
}

interface PoemAnnotationManagerProps {
  poem: Poem | null;
  open: boolean;
  onClose: () => void;
}

interface PoemAnnotationPanelProps {
  poem: Poem;
  className?: string;
}

interface PoemAnnotationDraftPanelProps {
  content: string[];
  annotations: PoemAnnotationDraft[];
  onChange: (annotations: PoemAnnotationDraft[]) => void;
  className?: string;
  isLoading?: boolean;
}

interface PoemAnnotationEditorProps {
  content: string[];
  annotations: PoemAnnotationDraft[];
  className?: string;
  isLoading?: boolean;
  isSaving?: boolean;
  isDeleting?: boolean;
  onSaveAnnotation: (annotation: PoemAnnotationDraft, index: number | null) => void;
  onDeleteAnnotation: (annotation: PoemAnnotationDraft, index: number) => void;
}

interface AnnotationFormState {
  title: string;
  content: string;
}

const emptyForm: AnnotationFormState = {
  title: '',
  content: '',
};

export function toAnnotationDrafts(annotations: PoemAnnotation[] = []): PoemAnnotationDraft[] {
  return sortAnnotationDraftsByRange(
    annotations.map((annotation) => ({
      id: annotation.id,
      targetField: 'content',
      startLine: annotation.startLine,
      startOffset: annotation.startOffset,
      endLine: annotation.endLine,
      endOffset: annotation.endOffset,
      selectedText: annotation.selectedText,
      title: annotation.title,
      content: annotation.content,
      type: 'note',
      displayOrder: annotation.displayOrder,
    })),
  );
}

export function sortAnnotationDraftsByRange(annotations: PoemAnnotationDraft[]): PoemAnnotationDraft[] {
  return [...annotations].sort((left, right) => {
    if (left.startLine !== right.startLine) return left.startLine - right.startLine;
    if (left.startOffset !== right.startOffset) return left.startOffset - right.startOffset;
    if (left.endLine !== right.endLine) return left.endLine - right.endLine;
    if (left.endOffset !== right.endOffset) return left.endOffset - right.endOffset;
    return (left.id ?? '').localeCompare(right.id ?? '');
  });
}

export function PoemAnnotationDraftPanel({
  content,
  annotations,
  onChange,
  className,
  isLoading,
}: PoemAnnotationDraftPanelProps) {
  const handleSaveAnnotation = (annotation: PoemAnnotationDraft, index: number | null) => {
    const nextAnnotations = [...annotations];
    if (index === null) {
      nextAnnotations.push(annotation);
    } else {
      nextAnnotations[index] = annotation;
    }
    onChange(sortAnnotationDraftsByRange(nextAnnotations));
  };

  const handleDeleteAnnotation = (_annotation: PoemAnnotationDraft, index: number) => {
    onChange(annotations.filter((_, itemIndex) => itemIndex !== index));
  };

  return (
    <PoemAnnotationEditor
      className={className}
      content={content}
      annotations={annotations}
      isLoading={isLoading}
      onSaveAnnotation={handleSaveAnnotation}
      onDeleteAnnotation={handleDeleteAnnotation}
    />
  );
}

export function PoemAnnotationPanel({ poem, className }: PoemAnnotationPanelProps) {
  const poemId = String(poem.id);
  const { data: annotations = [], isLoading } = useAdminPoemAnnotations(poemId);
  const createMutation = useAdminCreatePoemAnnotation();
  const updateMutation = useAdminUpdatePoemAnnotation();
  const deleteMutation = useAdminDeletePoemAnnotation();
  const drafts = toAnnotationDrafts(annotations);

  const handleSaveAnnotation = (annotation: PoemAnnotationDraft, index: number | null) => {
    if (index !== null && annotation.id) {
      updateMutation.mutate({ id: annotation.id, poemId, data: annotation });
      return;
    }
    createMutation.mutate({ poemId, data: annotation });
  };

  const handleDeleteAnnotation = (annotation: PoemAnnotationDraft) => {
    if (!annotation.id) return;
    deleteMutation.mutate({ id: annotation.id, poemId });
  };

  return (
    <PoemAnnotationEditor
      className={className}
      content={poem.content}
      annotations={drafts}
      isLoading={isLoading}
      isSaving={createMutation.isPending || updateMutation.isPending}
      isDeleting={deleteMutation.isPending}
      onSaveAnnotation={handleSaveAnnotation}
      onDeleteAnnotation={handleDeleteAnnotation}
    />
  );
}

function PoemAnnotationEditor({
  content,
  annotations,
  className,
  isLoading,
  isSaving,
  isDeleting,
  onSaveAnnotation,
  onDeleteAnnotation,
}: PoemAnnotationEditorProps) {
  const contentRef = useRef<HTMLDivElement | null>(null);
  const [selectionRange, setSelectionRange] = useState<AnnotationRange | null>(null);
  const [activeLineIndex, setActiveLineIndex] = useState<number | null>(null);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [editingIndex, setEditingIndex] = useState<number | null>(null);
  const [form, setForm] = useState<AnnotationFormState>(emptyForm);
  const { confirm, ConfirmDialog } = useConfirmDialog();
  const contentKey = content.join('\n');

  useEffect(() => {
    resetEditor();
  }, [contentKey]);

  const handleContentMouseUp = () => {
    const selection = window.getSelection();
    if (!selection || selection.isCollapsed || selection.rangeCount === 0 || !contentRef.current) {
      return;
    }
    const range = selection.getRangeAt(0);
    if (!contentRef.current.contains(range.commonAncestorContainer)) {
      return;
    }

    const startLineElement = getLineElement(range.startContainer);
    const endLineElement = getLineElement(range.endContainer);
    if (!startLineElement || !endLineElement) {
      return;
    }

    const normalized = normalizeSelectionRange(content, {
      startLine: Number(startLineElement.dataset.lineIndex),
      startOffset: getOffsetInLine(startLineElement, range.startContainer, range.startOffset),
      endLine: Number(endLineElement.dataset.lineIndex),
      endOffset: getOffsetInLine(endLineElement, range.endContainer, range.endOffset),
    });
    if (normalized) {
      setActiveLineIndex(normalized.startLine);
      setSelectionRange(normalized);
      setEditingIndex(null);
      setForm(emptyForm);
      setIsPopoverOpen(true);
    }
    selection.removeAllRanges();
  };

  const handleEdit = (annotation: PoemAnnotationDraft, index: number) => {
    setActiveLineIndex(annotation.startLine);
    setEditingIndex(index);
    setSelectionRange({
      startLine: annotation.startLine,
      startOffset: annotation.startOffset,
      endLine: annotation.endLine,
      endOffset: annotation.endOffset,
      selectedText: annotation.selectedText,
    });
    setForm({
      title: annotation.title ?? '',
      content: annotation.content,
    });
    setIsPopoverOpen(true);
  };

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    if (!selectionRange || !form.content.trim()) {
      return;
    }
    const existingAnnotation = editingIndex === null ? null : annotations[editingIndex];
    onSaveAnnotation(
      {
        id: existingAnnotation?.id,
        targetField: 'content',
        ...selectionRange,
        title: form.title.trim(),
        content: form.content.trim(),
        type: 'note',
        displayOrder: existingAnnotation?.displayOrder ?? 0,
      },
      editingIndex,
    );
    resetEditor();
  };

  const handleDelete = async (annotation: PoemAnnotationDraft, index: number) => {
    const confirmed = await confirm({
      title: '删除标注',
      description: `确定要删除“${annotation.selectedText}”的标注吗？此操作不可撤销。`,
      confirmText: '删除',
      destructive: true,
    });
    if (!confirmed) return;
    onDeleteAnnotation(annotation, index);
    resetEditor();
  };

  function resetEditor() {
    setActiveLineIndex(null);
    setSelectionRange(null);
    setIsPopoverOpen(false);
    setEditingIndex(null);
    setForm(emptyForm);
  }

  const renderAnnotationForm = () => (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <h4 className="font-medium">{editingIndex !== null ? '编辑标注' : '新增标注'}</h4>
          <p className="mt-1 truncate text-sm text-muted-foreground">{selectionRange?.selectedText}</p>
        </div>
        <Button type="button" variant="ghost" size="sm" onClick={resetEditor}>
          <X className="h-4 w-4" />
        </Button>
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium">标题</label>
        <Input
          value={form.title}
          onChange={(event) => setForm({ ...form, title: event.target.value })}
          placeholder="可选，默认使用被标注文本"
        />
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium">内容 *</label>
        <Input
          value={form.content}
          onChange={(event) => setForm({ ...form, content: event.target.value })}
          placeholder="填写标注解释"
        />
      </div>
      <Button type="submit" disabled={!selectionRange || !form.content.trim() || isSaving}>
        {editingIndex !== null ? '保存标注' : '添加标注'}
      </Button>
    </form>
  );

  return (
    <div className={cn('space-y-4 rounded-md border p-4', className)}>
      <div className="flex items-center justify-between">
        <h3 className="font-serif text-base font-semibold">诗词标注</h3>
        {editingIndex !== null && (
          <Button type="button" variant="ghost" size="sm" onClick={resetEditor}>
            <Plus className="mr-1 h-4 w-4" />
            新增
          </Button>
        )}
      </div>

      <div className="grid gap-4 lg:grid-cols-[1.15fr_0.85fr]">
        <div className="min-w-0 space-y-3">
          <div
            ref={contentRef}
            className="select-text space-y-2 rounded-md border bg-muted/20 p-4 font-serif text-lg leading-loose"
            onMouseUp={handleContentMouseUp}
          >
            {content.length === 0 ? (
              <p className="text-sm text-muted-foreground">请先填写诗词正文</p>
            ) : (
              content.map((line, index) => (
                <Popover
                  key={`${line}-${index}`}
                  open={activeLineIndex === index && Boolean(selectionRange) && isPopoverOpen}
                  onOpenChange={(open) => {
                    if (!open) resetEditor();
                  }}
                >
                  <PopoverTrigger asChild>
                    <p data-line-index={index} className="min-h-7 cursor-text">
                      {line}
                    </p>
                  </PopoverTrigger>
                  <PopoverContent align="start" side="top" className="w-80 font-sans">
                    {renderAnnotationForm()}
                  </PopoverContent>
                </Popover>
              ))
            )}
          </div>
          <div className="rounded bg-muted/40 p-3 text-sm">
            {selectionRange ? (
              <>
                <div className="text-muted-foreground">已选文本</div>
                <div className="mt-1 whitespace-pre-wrap font-medium">{selectionRange.selectedText}</div>
              </>
            ) : (
              <span className="text-muted-foreground">未选择文本</span>
            )}
          </div>
        </div>

        <div className="min-w-0 space-y-4">
          <div className="space-y-3">
            <h4 className="font-medium">已有标注</h4>
            {isLoading ? (
              <div className="py-4 text-sm text-muted-foreground">加载中...</div>
            ) : annotations.length === 0 ? (
              <div className="py-4 text-sm text-muted-foreground">暂无标注</div>
            ) : (
              annotations.map((annotation, index) => (
                <div key={annotation.id ?? `${annotation.startLine}-${annotation.startOffset}-${index}`} className="rounded-md border p-4">
                  <div className="mb-2 flex items-start justify-between gap-3">
                    <div>
                      <div className="font-medium">
                        {index + 1}. {annotation.title || annotation.selectedText}
                      </div>
                      <div className="mt-1 whitespace-pre-wrap text-sm text-muted-foreground">{annotation.selectedText}</div>
                    </div>
                    <div className="flex shrink-0 gap-1">
                      <Button type="button" variant="ghost" size="sm" onClick={() => handleEdit(annotation, index)}>
                        <Edit2 className="h-4 w-4" />
                      </Button>
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        disabled={isDeleting}
                        onClick={() => handleDelete(annotation, index)}
                      >
                        <Trash2 className="h-4 w-4 text-cinnabar" />
                      </Button>
                    </div>
                  </div>
                  <p className="text-sm leading-6">{annotation.content}</p>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
      <ConfirmDialog />
    </div>
  );
}

export function PoemAnnotationManager({ poem, open, onClose }: PoemAnnotationManagerProps) {
  if (!open || !poem) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="flex max-h-[92vh] w-full max-w-6xl flex-col rounded-lg bg-background shadow-lg">
        <div className="flex items-center justify-between border-b px-6 py-4">
          <div>
            <h2 className="font-serif text-xl font-bold">诗词标注</h2>
            <p className="mt-1 text-sm text-muted-foreground">{poem.title}</p>
          </div>
          <button type="button" className="rounded p-1 hover:bg-muted" onClick={onClose}>
            <X className="h-5 w-5" />
          </button>
        </div>
        <div className="min-h-0 overflow-y-auto p-6">
          <PoemAnnotationPanel poem={poem} />
        </div>
      </div>
    </div>
  );
}

function getLineElement(node: Node): HTMLElement | null {
  const element = node.nodeType === Node.ELEMENT_NODE ? (node as Element) : node.parentElement;
  return element?.closest<HTMLElement>('[data-line-index]') ?? null;
}

function getOffsetInLine(lineElement: HTMLElement, node: Node, offset: number): number {
  const range = document.createRange();
  range.selectNodeContents(lineElement);
  range.setEnd(node, offset);
  return Array.from(range.toString()).length;
}
