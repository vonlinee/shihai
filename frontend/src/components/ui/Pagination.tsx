import { FormEvent, useEffect, useMemo, useState } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/utils/cn'

interface PaginationProps {
  page: number
  pageSize: number
  total: number
  onPageChange: (page: number) => void
  onPageSizeChange?: (pageSize: number) => void
  pageSizeOptions?: number[]
  disabled?: boolean
  className?: string
}

const defaultPageSizeOptions = [10, 20, 50, 100]

export function Pagination({
  page,
  pageSize,
  total,
  onPageChange,
  onPageSizeChange,
  pageSizeOptions = defaultPageSizeOptions,
  disabled = false,
  className,
}: PaginationProps) {
  const totalPages = useMemo(() => Math.max(1, Math.ceil(total / pageSize)), [pageSize, total])
  const currentPage = Math.min(Math.max(page, 1), totalPages)
  const [jumpPage, setJumpPage] = useState(String(currentPage))

  useEffect(() => {
    setJumpPage(String(currentPage))
  }, [currentPage])

  const goToPage = (nextPage: number) => {
    onPageChange(Math.min(Math.max(nextPage, 1), totalPages))
  }

  const handleJump = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const nextPage = Number(jumpPage)
    if (!Number.isFinite(nextPage)) return
    goToPage(Math.trunc(nextPage))
  }

  if (total <= 0) {
    return null
  }

  return (
    <div className={cn('flex flex-col gap-3 border-t pt-4 text-sm text-muted-foreground md:flex-row md:items-center md:justify-between', className)}>
      <div className="flex flex-wrap items-center gap-3">
        <span>共 {total} 条</span>
        <label className="flex items-center gap-2">
          <span>每页</span>
          <select
            className="h-9 rounded-md border border-input bg-background px-2 text-sm text-foreground outline-none focus:ring-2 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
            value={pageSize}
            disabled={disabled || !onPageSizeChange}
            onChange={(event) => onPageSizeChange?.(Number(event.target.value))}
          >
            {pageSizeOptions.map((option) => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
          </select>
          <span>条</span>
        </label>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={disabled || currentPage <= 1}
          onClick={() => goToPage(currentPage - 1)}
          title="上一页"
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <span className="min-w-24 text-center">
          第 {currentPage} / {totalPages} 页
        </span>
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={disabled || currentPage >= totalPages}
          onClick={() => goToPage(currentPage + 1)}
          title="下一页"
        >
          <ChevronRight className="h-4 w-4" />
        </Button>
        <form className="flex items-center gap-2" onSubmit={handleJump}>
          <span>跳至</span>
          <Input
            className="h-9 w-20"
            type="number"
            min={1}
            max={totalPages}
            value={jumpPage}
            disabled={disabled}
            onChange={(event) => setJumpPage(event.target.value)}
          />
          <span>页</span>
          <Button type="submit" variant="outline" size="sm" disabled={disabled}>
            跳转
          </Button>
        </form>
      </div>
    </div>
  )
}
