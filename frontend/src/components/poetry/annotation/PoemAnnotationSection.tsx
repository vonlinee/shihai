import { useRef, useState } from 'react';
import type { ReactNode } from 'react';
import { ChevronDown, ChevronRight } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { cn } from '@/utils/cn';
import { normalizePoemContent } from '@/utils/poemContent';
import type { PoemAnnotation } from '@/types';

interface PoemAnnotationSectionProps {
  content: string[] | string | null | undefined;
  annotations: PoemAnnotation[];
  className?: string;
  contentClassName?: string;
}

export function PoemAnnotationSection({
  content,
  annotations,
  className,
  contentClassName,
}: PoemAnnotationSectionProps) {
  const [showInlineMarkers, setShowInlineMarkers] = useState(true);
  const [activeAnnotationId, setActiveAnnotationId] = useState<string | null>(null);
  const itemRefs = useRef<Record<string, HTMLLIElement | null>>({});
  const contentLines = normalizePoemContent(content);

  const handleMarkerClick = (annotation: PoemAnnotation) => {
    setActiveAnnotationId(annotation.id);
    requestAnimationFrame(() => {
      itemRefs.current[annotation.id]?.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    });
  };

  if (contentLines.length === 0) {
    return null;
  }

  return (
    <section className={cn('space-y-6', className)}>
      <div className="space-y-4 text-center">
        {annotations.length > 0 && (
          <div className="flex justify-center">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => setShowInlineMarkers((current) => !current)}
            >
              {showInlineMarkers ? <ChevronDown className="mr-1 h-4 w-4" /> : <ChevronRight className="mr-1 h-4 w-4" />}
              {showInlineMarkers ? '隐藏正文标注序号' : '显示正文标注序号'}
            </Button>
          </div>
        )}
        <div className={cn('poetry-text space-y-2 text-xl leading-loose', contentClassName)}>
          {contentLines.map((line, index) => (
            <p key={`${line}-${index}`}>
              {renderLineWithAnnotations(line, index, annotations, showInlineMarkers, handleMarkerClick)}
            </p>
          ))}
        </div>
      </div>

      {annotations.length > 0 && (
        <div className="border-t pt-6">
          <h3 className="font-serif font-semibold">标注</h3>
          <ol className="mt-4 overflow-hidden rounded-md border bg-muted/30">
            {annotations.map((annotation) => (
              <li
                key={annotation.id}
                ref={(node) => {
                  itemRefs.current[annotation.id] = node;
                }}
                className={cn(
                  'border-b p-4 transition-colors last:border-b-0',
                  activeAnnotationId === annotation.id && 'bg-primary/5',
                )}
              >
                <div className="flex items-start gap-2">
                  <span className="inline-flex h-6 min-w-6 items-center justify-center rounded-full bg-primary px-2 text-xs font-medium text-primary-foreground">
                    {annotation.displayNo}
                  </span>
                  <p className="min-w-0 whitespace-pre-wrap text-sm leading-6 text-foreground/80">
                    <span className="font-medium text-foreground">{annotation.selectedText}</span>
                    <span>：</span>
                    <span>{annotation.content}</span>
                  </p>
                </div>
              </li>
            ))}
          </ol>
        </div>
      )}
    </section>
  );
}

function renderLineWithAnnotations(
  line: string,
  lineIndex: number,
  annotations: PoemAnnotation[],
  showInlineMarkers: boolean,
  onMarkerClick: (annotation: PoemAnnotation) => void,
) {
  if (!showInlineMarkers) {
    return line;
  }

  const chars = Array.from(line);
  const lineAnnotations = annotations
    .filter((annotation) => annotation.targetField === 'content' && annotation.startLine === lineIndex)
    .sort((left, right) => left.startOffset - right.startOffset);

  if (lineAnnotations.length === 0) {
    return line;
  }

  const parts: ReactNode[] = [];
  let cursor = 0;
  lineAnnotations.forEach((annotation) => {
    const startOffset = Math.max(0, Math.min(annotation.startOffset, chars.length));
    const endOffset = annotation.endLine === lineIndex
      ? Math.max(startOffset, Math.min(annotation.endOffset, chars.length))
      : chars.length;

    if (startOffset > cursor) {
      parts.push(chars.slice(cursor, startOffset).join(''));
    }

    parts.push(
      <span key={annotation.id} className="inline-flex items-start">
        <span className="bg-primary/10">{chars.slice(startOffset, endOffset).join('')}</span>
        <button
          type="button"
          className="ml-0.5 align-super text-xs font-medium text-primary hover:underline"
          onClick={() => onMarkerClick(annotation)}
        >
          {annotation.displayNo}
        </button>
      </span>,
    );
    cursor = endOffset;
  });

  if (cursor < chars.length) {
    parts.push(chars.slice(cursor).join(''));
  }

  return parts;
}
