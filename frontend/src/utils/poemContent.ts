export function normalizePoemContent(content: string[] | string | null | undefined): string[] {
  if (Array.isArray(content)) {
    return content.filter((line) => line.trim().length > 0);
  }
  if (typeof content === 'string') {
    return content
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line.length > 0);
  }
  return [];
}

export function splitPoemContentInput(content: string): string[] {
  return content
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter((line) => line.length > 0);
}
