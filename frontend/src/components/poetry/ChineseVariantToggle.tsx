import { Languages } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { cn } from '@/utils/cn';

export type ChineseVariant = 'simplified' | 'traditional';
export type ChineseVariantConversionMode = 't2s' | 's2t';

interface ChineseVariantToggleProps {
  value?: ChineseVariant;
  disabled?: boolean;
  isLoading?: boolean;
  className?: string;
  onChange: (variant: ChineseVariant, mode: ChineseVariantConversionMode) => void;
}

const variantOptions: { value: ChineseVariant; label: string; mode: ChineseVariantConversionMode }[] = [
  { value: 'simplified', label: '\u7b80', mode: 't2s' },
  { value: 'traditional', label: '\u7e41', mode: 's2t' },
];

export function ChineseVariantToggle({
  value,
  disabled = false,
  isLoading = false,
  className,
  onChange,
}: ChineseVariantToggleProps) {
  return (
    <div className={cn('inline-flex items-center rounded-md border bg-background p-0.5', className)}>
      {variantOptions.map((option) => {
        const isActive = value === option.value;
        return (
          <Button
            key={option.value}
            type="button"
            variant={isActive ? 'secondary' : 'ghost'}
            size="sm"
            className="h-7 gap-1 rounded-sm px-2"
            disabled={disabled || isLoading}
            onClick={() => onChange(option.value, option.mode)}
          >
            <Languages className="h-3.5 w-3.5" />
            {option.label}
          </Button>
        );
      })}
    </div>
  );
}
