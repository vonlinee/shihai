import { cn } from '@/utils/cn';
import { normalizePoemContent } from '@/utils/poemContent';

interface PoemContentProps {
  content: string[] | string | null | undefined;
  className?: string;
  lineClassName?: string;
}

export const PoemContent = ({ content, className, lineClassName }: PoemContentProps) => {
  const contentLines = normalizePoemContent(content);

  if (contentLines.length === 0) {
    return null;
  }

  return (
    <div className={cn('poetry-text space-y-1', className)}>
      {contentLines.map((line, index) => (
        <p key={`${line}-${index}`} className={lineClassName}>
          {line}
        </p>
      ))}
    </div>
  );
};
