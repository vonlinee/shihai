import { cn } from '@/utils/cn'

interface LogoProps {
  /** Whether to render the Chinese wordmark next to the animal mark. */
  showWordmark?: boolean
  /** Use the compact size for collapsed navigation. */
  compact?: boolean
  className?: string
}

export function Logo({ showWordmark = true, compact = false, className }: LogoProps) {
  return (
    <span className={cn('inline-flex items-center gap-2', className)}>
      <svg
        aria-hidden="true"
        className={cn('shrink-0', compact ? 'h-10 w-10' : 'h-9 w-9')}
        viewBox="0 0 48 48"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M5.5 24c5.7-7.1 14.7-10.3 23.1-7.8L36 9.4l-1 10.4c2.5 1 4.7 2.4 6.8 4.2-2.1 1.8-4.3 3.2-6.8 4.2l1 10.4-7.4-6.8C20.2 34.3 11.2 31.1 5.5 24Z"
          className="fill-primary"
        />
        <circle cx="30.2" cy="20.8" r="1.45" className="fill-background" />
        <path
          d="M11.1 24c2.7-2.6 5.8-4.2 9.3-4.8M11.1 24c2.7 2.6 5.8 4.2 9.3 4.8"
          className="stroke-primary-foreground"
          strokeWidth="1.25"
          strokeLinecap="round"
          opacity="0.8"
        />
        <path
          d="M10.5 29.1c4.8-.8 9-.1 13.5 2.8 4.5-2.9 8.7-3.6 13.5-2.8v9.2c-4.8-.8-9-.1-13.5 2.8-4.5-2.9-8.7-3.6-13.5-2.8v-9.2Z"
          className="stroke-accent"
          strokeWidth="1.8"
          strokeLinejoin="round"
        />
        <path
          d="M24 31.9v9.2"
          className="stroke-accent"
          strokeWidth="1.4"
          strokeLinecap="round"
        />
      </svg>
      {showWordmark && (
        <span className="font-serif text-xl font-bold leading-none tracking-normal text-foreground">
          识海
        </span>
      )}
    </span>
  )
}
