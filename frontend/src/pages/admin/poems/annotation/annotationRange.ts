export interface AnnotationRange {
  startLine: number;
  startOffset: number;
  endLine: number;
  endOffset: number;
  selectedText: string;
}

export function extractSelectedText(content: string[], range: Omit<AnnotationRange, 'selectedText'>): string {
  const startLine = content[range.startLine] ?? '';
  const endLine = content[range.endLine] ?? '';

  if (range.startLine === range.endLine) {
    return Array.from(startLine).slice(range.startOffset, range.endOffset).join('');
  }

  const parts = [
    Array.from(startLine).slice(range.startOffset).join(''),
    ...content.slice(range.startLine + 1, range.endLine),
    Array.from(endLine).slice(0, range.endOffset).join(''),
  ];
  return parts.join('\n');
}

export function normalizeSelectionRange(content: string[], range: Omit<AnnotationRange, 'selectedText'>): AnnotationRange | null {
  if (range.startLine > range.endLine || (range.startLine === range.endLine && range.startOffset >= range.endOffset)) {
    return null;
  }
  if (range.startLine < 0 || range.endLine >= content.length) {
    return null;
  }
  const selectedText = extractSelectedText(content, range);
  if (!selectedText.trim()) {
    return null;
  }
  return { ...range, selectedText };
}

export function normalizeTextareaSelectionRange(content: string, selectionStart: number, selectionEnd: number): AnnotationRange | null {
  if (selectionStart >= selectionEnd) {
    return null;
  }

  const lines = content.split(/\r?\n/);
  const startPoint = getTextareaPoint(content, selectionStart);
  const endPoint = getTextareaPoint(content, selectionEnd);

  return normalizeSelectionRange(lines, {
    startLine: startPoint.line,
    startOffset: startPoint.offset,
    endLine: endPoint.line,
    endOffset: endPoint.offset,
  });
}

function getTextareaPoint(content: string, position: number): { line: number; offset: number } {
  const precedingLines = content.slice(0, position).split(/\r?\n/);
  const currentLine = precedingLines[precedingLines.length - 1] ?? '';
  return {
    line: precedingLines.length - 1,
    offset: Array.from(currentLine).length,
  };
}
