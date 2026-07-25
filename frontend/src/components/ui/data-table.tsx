import {
  type MouseEvent as ReactMouseEvent,
  type ReactNode,
  type TouchEvent as ReactTouchEvent,
  useMemo,
  useRef,
  useState,
} from 'react'
import {
  type ColumnDef,
  type ColumnFiltersState,
  type Column,
  type ColumnOrderState,
  type ColumnSizingState,
  type OnChangeFn,
  type RowData,
  type SortingState,
  type VisibilityState,
  flexRender,
  getFacetedRowModel,
  getFacetedUniqueValues,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { Check, ChevronDown, ChevronUp, ChevronsUpDown, Clipboard, Filter, GripVertical, Info, Loader2, X } from 'lucide-react'

import { cn } from '@/utils/cn'

import { Button } from './button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './dropdown-menu'
import { Input } from './input'
import { Skeleton } from './skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './table'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './tooltip'

interface DataTableFilterOption {
  label: string
  value: string
}

interface DataTableDateRangeFilter {
  from?: string
  to?: string
}

type DataTableFilterVariant = 'text' | 'select' | 'multiSelect' | 'dateRange'

declare module '@tanstack/react-table' {
  interface ColumnMeta<TData extends RowData, TValue> {
    align?: 'left' | 'center' | 'right'
    className?: string
    copyable?: boolean | ((value: TValue, row: TData) => string)
    draggable?: boolean
    emptyText?: ReactNode
    filterOptions?: DataTableFilterOption[]
    filterPlaceholder?: string
    filterVariant?: DataTableFilterVariant
    fixed?: 'left' | 'right'
    headerAlign?: 'left' | 'center' | 'right'
    headerClassName?: string
    hidden?: boolean
    loading?: boolean
    minWidth?: number | string
    resizable?: boolean
    tooltip?: ReactNode
    width?: number | string
  }
}

export type DataTableColumn<TData, TValue = unknown> = ColumnDef<TData, TValue>

interface DataTableProps<TData> {
  className?: string
  columnFilters?: ColumnFiltersState
  columnOrder?: ColumnOrderState
  columnSizing?: ColumnSizingState
  columnVisibility?: VisibilityState
  columns: DataTableColumn<TData, unknown>[]
  data: TData[]
  emptyText?: string
  enableColumnDragging?: boolean
  enableColumnFilters?: boolean
  enableColumnResizing?: boolean
  enableSorting?: boolean
  getRowId?: (originalRow: TData, index: number) => string
  loading?: boolean
  loadingText?: string
  onColumnFiltersChange?: OnChangeFn<ColumnFiltersState>
  onColumnOrderChange?: OnChangeFn<ColumnOrderState>
  onColumnSizingChange?: OnChangeFn<ColumnSizingState>
  onColumnVisibilityChange?: OnChangeFn<VisibilityState>
  onSortingChange?: OnChangeFn<SortingState>
  sorting?: SortingState
  tableClassName?: string
}

function toCssSize(value: number | string | undefined): string | undefined {
  if (typeof value === 'number') return `${value}px`
  return value
}

function toNumericSize(value: number | string | undefined): number | undefined {
  if (typeof value === 'number') return value
  return undefined
}

function getTextAlignClass(align: 'left' | 'center' | 'right' | undefined) {
  if (align === 'center') return 'text-center'
  if (align === 'right') return 'text-right'
  return 'text-left'
}

function normalizeDateRangeFilterValue(value: unknown): DataTableDateRangeFilter {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return {}
  const range = value as DataTableDateRangeFilter
  return {
    from: typeof range.from === 'string' ? range.from : undefined,
    to: typeof range.to === 'string' ? range.to : undefined,
  }
}

function isDateRangeFilterActive(value: DataTableDateRangeFilter) {
  return Boolean(value.from?.trim() || value.to?.trim())
}

function parseDateTime(value: unknown): number | undefined {
  if (value instanceof Date) {
    const timestamp = value.getTime()
    return Number.isNaN(timestamp) ? undefined : timestamp
  }
  if (typeof value === 'number') return Number.isFinite(value) ? value : undefined
  if (typeof value !== 'string' || !value.trim()) return undefined

  const timestamp = new Date(value).getTime()
  return Number.isNaN(timestamp) ? undefined : timestamp
}

function dateRangeFilterFn(row: { getValue: (columnId: string) => unknown }, columnId: string, filterValue: unknown) {
  const range = normalizeDateRangeFilterValue(filterValue)
  if (!isDateRangeFilterActive(range)) return true

  const cellTimestamp = parseDateTime(row.getValue(columnId))
  if (cellTimestamp === undefined) return false

  const fromTimestamp = parseDateTime(range.from)
  const toTimestamp = parseDateTime(range.to)

  if (fromTimestamp !== undefined && cellTimestamp < fromTimestamp) return false
  if (toTimestamp !== undefined && cellTimestamp > toTimestamp) return false
  return true
}

function getColumnId<TData>(column: DataTableColumn<TData, unknown>, index: number): string {
  if (column.id) return column.id
  if ('accessorKey' in column && typeof column.accessorKey === 'string') return column.accessorKey
  return `column-${index}`
}

function getInitialColumnOrder<TData>(columns: DataTableColumn<TData, unknown>[]) {
  return columns.map((column, index) => getColumnId(column, index))
}

function orderFixedColumns<TData>(
  columns: DataTableColumn<TData, unknown>[],
  columnOrder: ColumnOrderState,
) {
  const knownColumnIds = columns.map((column, index) => getColumnId(column, index))
  const normalizedOrder = [
    ...columnOrder.filter((columnId) => knownColumnIds.includes(columnId)),
    ...knownColumnIds.filter((columnId) => !columnOrder.includes(columnId)),
  ]
  const columnById = new Map(columns.map((column, index) => [getColumnId(column, index), column]))
  const leftFixedColumnIds = normalizedOrder.filter((columnId) => columnById.get(columnId)?.meta?.fixed === 'left')
  const rightFixedColumnIds = normalizedOrder.filter((columnId) => columnById.get(columnId)?.meta?.fixed === 'right')
  const normalColumnIds = normalizedOrder.filter((columnId) => {
    const fixed = columnById.get(columnId)?.meta?.fixed
    return fixed !== 'left' && fixed !== 'right'
  })

  return [...leftFixedColumnIds, ...normalColumnIds, ...rightFixedColumnIds]
}

function getInitialColumnVisibility<TData>(columns: DataTableColumn<TData, unknown>[]) {
  return columns.reduce<VisibilityState>((visibility, column, index) => {
    if (column.meta?.hidden) visibility[getColumnId(column, index)] = false
    return visibility
  }, {})
}

function normalizeColumns<TData>(columns: DataTableColumn<TData, unknown>[]) {
  return columns.map((column) => {
    const filterVariant = column.meta?.filterVariant
    const filterFn = column.filterFn ?? (
      filterVariant === 'select'
        ? (row, columnId, filterValue) => String(row.getValue(columnId) ?? '') === String(filterValue)
        : filterVariant === 'multiSelect'
          ? (row, columnId, filterValue) => {
            if (!Array.isArray(filterValue) || filterValue.length === 0) return true
            return filterValue.includes(String(row.getValue(columnId) ?? ''))
          }
          : filterVariant === 'dateRange'
            ? dateRangeFilterFn
            : undefined
    )

    return {
      ...column,
      enableResizing: column.enableResizing ?? column.meta?.resizable !== false,
      filterFn,
      minSize: column.minSize ?? toNumericSize(column.meta?.minWidth),
      size: column.size ?? toNumericSize(column.meta?.width),
    }
  })
}

function getPinnedOffset(
  orderedColumns: { id: string; fixed?: 'left' | 'right'; width: number }[],
  columnId: string,
  fixed: 'left' | 'right' | undefined,
) {
  if (!fixed) return undefined

  const columnIndex = orderedColumns.findIndex((column) => column.id === columnId)
  if (columnIndex < 0) return undefined

  const pinnedColumns = fixed === 'left'
    ? orderedColumns.slice(0, columnIndex)
    : orderedColumns.slice(columnIndex + 1)

  const offset = pinnedColumns
    .filter((column) => column.fixed === fixed)
    .reduce((sum, column) => sum + column.width, 0)

  return offset ? `${offset}px` : '0px'
}

function getColumnFilterOptions<TData>(column: Column<TData, unknown>): DataTableFilterOption[] {
  const customOptions = column.columnDef.meta?.filterOptions
  if (customOptions) return customOptions

  return Array.from(column.getFacetedUniqueValues().keys())
    .filter((value) => value !== undefined && value !== null && String(value).trim() !== '')
    .slice(0, 50)
    .map((value) => {
      const text = String(value)
      return { label: text, value: text }
    })
}

function isEmptyCellValue(value: unknown) {
  return value === undefined || value === null || value === ''
}

function getCopyText<TData, TValue>(copyable: boolean | ((value: TValue, row: TData) => string), value: TValue, row: TData) {
  if (typeof copyable === 'function') return copyable(value, row)
  return isEmptyCellValue(value) ? '' : String(value)
}

interface ColumnFilterMenuProps<TData> {
  column: Column<TData, unknown>
  placeholder?: string
}

function ColumnFilterMenu<TData>({ column, placeholder }: ColumnFilterMenuProps<TData>) {
  const filterVariant = column.columnDef.meta?.filterVariant ?? 'text'
  const rawFilterValue = column.getFilterValue()
  const textFilterValue = (rawFilterValue as string | undefined) ?? ''
  const selectedValues = Array.isArray(rawFilterValue) ? rawFilterValue.map(String) : []
  const dateRangeValue = normalizeDateRangeFilterValue(rawFilterValue)
  const filterOptions = filterVariant === 'dateRange' ? [] : getColumnFilterOptions(column)
  const hasFilter = filterVariant === 'multiSelect'
    ? selectedValues.length > 0
    : filterVariant === 'dateRange'
      ? isDateRangeFilterActive(dateRangeValue)
      : textFilterValue.trim().length > 0

  const toggleMultiSelectFilter = (value: string) => {
    const nextValue = selectedValues.includes(value)
      ? selectedValues.filter((item) => item !== value)
      : [...selectedValues, value]
    column.setFilterValue(nextValue.length > 0 ? nextValue : undefined)
  }

  const updateDateRangeFilter = (key: keyof DataTableDateRangeFilter, value: string) => {
    const nextValue = {
      ...dateRangeValue,
      [key]: value || undefined,
    }
    column.setFilterValue(isDateRangeFilterActive(nextValue) ? nextValue : undefined)
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          type="button"
          aria-label="筛选"
          variant="ghost"
          size="icon"
          className={cn(
            'inline-flex h-7 w-7 items-center justify-center rounded-sm text-muted-foreground hover:bg-background/80 hover:text-foreground',
            hasFilter && 'bg-background text-primary shadow-sm',
          )}
        >
          <Filter className="h-3.5 w-3.5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align="end"
        sideOffset={8}
        className={cn('p-2', filterVariant === 'dateRange' ? 'w-[32rem] max-w-[calc(100vw-2rem)]' : 'w-56')}
      >
        {filterVariant === 'dateRange' ? (
          <div className="space-y-3">
            <div className="flex items-center gap-2">
              <div className="min-w-0 flex-1 text-xs font-medium text-muted-foreground">
                {placeholder ?? '选择时间范围'}
              </div>
              {hasFilter && (
                <Button
                  type="button"
                  aria-label="清除筛选"
                  onClick={() => column.setFilterValue(undefined)}
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8 shrink-0 rounded-sm text-muted-foreground hover:bg-muted hover:text-foreground"
                >
                  <X className="h-3.5 w-3.5" />
                </Button>
              )}
            </div>
            <div className="flex items-center rounded-md border border-input bg-background focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2">
              <Input
                aria-label="开始时间"
                type="datetime-local"
                step={1}
                value={dateRangeValue.from ?? ''}
                onChange={(event) => updateDateRangeFilter('from', event.target.value)}
                className="h-9 min-w-0 flex-1 border-0 bg-transparent text-xs shadow-none focus-visible:ring-0 focus-visible:ring-offset-0"
              />
              <span className="shrink-0 px-2 text-sm text-muted-foreground">→</span>
              <Input
                aria-label="结束时间"
                type="datetime-local"
                step={1}
                value={dateRangeValue.to ?? ''}
                onChange={(event) => updateDateRangeFilter('to', event.target.value)}
                className="h-9 min-w-0 flex-1 border-0 bg-transparent text-xs shadow-none focus-visible:ring-0 focus-visible:ring-offset-0"
              />
            </div>
          </div>
        ) : (
          <>
          <div className="flex items-center gap-2">
            {filterVariant === 'text' && (
              <Input
                value={textFilterValue}
                onChange={(event) => column.setFilterValue(event.target.value)}
                placeholder={placeholder ?? '筛选...'}
                className="h-8 text-xs"
              />
            )}
            {filterVariant !== 'text' && (
              <div className="min-w-0 flex-1 text-xs font-medium text-muted-foreground">
                {placeholder ?? '选择筛选项'}
              </div>
            )}
            {hasFilter && (
              <Button
                type="button"
                aria-label="清除筛选"
                onClick={() => column.setFilterValue(undefined)}
                variant="ghost"
                size="icon"
                className="h-8 w-8 shrink-0 rounded-sm text-muted-foreground hover:bg-muted hover:text-foreground"
              >
                <X className="h-3.5 w-3.5" />
              </Button>
            )}
          </div>
          {filterOptions.length > 0 && (
            <div className="mt-2 max-h-56 overflow-y-auto">
              {filterOptions.map((option) => {
                const selected = filterVariant === 'multiSelect'
                  ? selectedValues.includes(option.value)
                  : textFilterValue === option.value
                return (
                  <DropdownMenuItem
                    key={option.value}
                    onSelect={(event) => {
                      if (filterVariant === 'multiSelect') {
                        event.preventDefault()
                        toggleMultiSelectFilter(option.value)
                        return
                      }
                      column.setFilterValue(selected ? undefined : option.value)
                    }}
                    className="cursor-pointer justify-between text-xs"
                  >
                    <span className="truncate">{option.label}</span>
                    {selected && <Check className="ml-2 h-3.5 w-3.5 shrink-0 text-primary" />}
                  </DropdownMenuItem>
                )
              })}
            </div>
          )}
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

export function DataTable<TData>({
  className,
  columnFilters,
  columnOrder,
  columnSizing,
  columnVisibility,
  columns,
  data,
  emptyText = '暂无数据',
  enableColumnDragging = false,
  enableColumnFilters = false,
  enableColumnResizing = true,
  enableSorting = true,
  getRowId,
  loading = false,
  loadingText = '加载中...',
  onColumnFiltersChange,
  onColumnOrderChange,
  onColumnSizingChange,
  onColumnVisibilityChange,
  onSortingChange,
  sorting,
  tableClassName,
}: DataTableProps<TData>) {
  const tableColumns = useMemo(() => normalizeColumns(columns), [columns])
  const initialColumnOrder = useMemo(() => getInitialColumnOrder(tableColumns), [tableColumns])
  const initialColumnVisibility = useMemo(() => getInitialColumnVisibility(tableColumns), [tableColumns])
  const [internalColumnFilters, setInternalColumnFilters] = useState<ColumnFiltersState>([])
  const [internalColumnOrder, setInternalColumnOrder] = useState<ColumnOrderState>(initialColumnOrder)
  const [internalColumnSizing, setInternalColumnSizing] = useState<ColumnSizingState>({})
  const [internalColumnVisibility, setInternalColumnVisibility] = useState<VisibilityState>(initialColumnVisibility)
  const [internalSorting, setInternalSorting] = useState<SortingState>([])
  const [draggingColumnId, setDraggingColumnId] = useState<string | null>(null)
  const [resizeGuideX, setResizeGuideX] = useState<number | null>(null)
  const tableWrapperRef = useRef<HTMLDivElement>(null)

  const activeColumnOrder = columnOrder ?? internalColumnOrder
  const pinnedColumnOrder = useMemo(
    () => orderFixedColumns(tableColumns, activeColumnOrder),
    [activeColumnOrder, tableColumns],
  )
  const activeColumnFilters = columnFilters ?? internalColumnFilters
  const activeColumnSizing = columnSizing ?? internalColumnSizing
  const activeColumnVisibility = columnVisibility ?? internalColumnVisibility
  const activeSorting = sorting ?? internalSorting

  const handleColumnOrderChange: OnChangeFn<ColumnOrderState> = (updater) => {
    if (onColumnOrderChange) {
      onColumnOrderChange(updater)
      return
    }
    setInternalColumnOrder(updater)
  }

  const table = useReactTable({
    columnResizeMode: 'onEnd',
    columns: tableColumns,
    data,
    enableColumnResizing,
    enableSorting,
    getCoreRowModel: getCoreRowModel(),
    getFacetedRowModel: getFacetedRowModel(),
    getFacetedUniqueValues: getFacetedUniqueValues(),
    getFilteredRowModel: getFilteredRowModel(),
    getRowId,
    getSortedRowModel: getSortedRowModel(),
    onColumnFiltersChange: onColumnFiltersChange ?? setInternalColumnFilters,
    onColumnOrderChange: handleColumnOrderChange,
    onColumnSizingChange: onColumnSizingChange ?? setInternalColumnSizing,
    onColumnVisibilityChange: onColumnVisibilityChange ?? setInternalColumnVisibility,
    onSortingChange: onSortingChange ?? setInternalSorting,
    state: {
      columnFilters: activeColumnFilters,
      columnOrder: pinnedColumnOrder,
      columnSizing: activeColumnSizing,
      columnVisibility: activeColumnVisibility,
      sorting: activeSorting,
    },
  })

  const leafColumns = table.getVisibleLeafColumns()
  const orderedColumns = leafColumns.map((column) => ({
    fixed: column.columnDef.meta?.fixed,
    id: column.id,
    width: column.getSize(),
  }))
  const tableMinWidth = leafColumns.reduce((total, column) => total + column.getSize(), 0)
  const rows = table.getRowModel().rows
  const shouldShowLoadingRows = loading && rows.length === 0
  const shouldShowLoadingOverlay = loading && rows.length > 0
  const moveColumn = (fromColumnId: string, toColumnId: string) => {
    if (fromColumnId === toColumnId) return

    const nextOrder = [...activeColumnOrder]
    const fromIndex = nextOrder.indexOf(fromColumnId)
    const toIndex = nextOrder.indexOf(toColumnId)
    if (fromIndex < 0 || toIndex < 0) return

    nextOrder.splice(fromIndex, 1)
    nextOrder.splice(toIndex, 0, fromColumnId)
    handleColumnOrderChange(nextOrder)
  }

  const updateResizeGuide = (clientX: number) => {
    const rect = tableWrapperRef.current?.getBoundingClientRect()
    if (!rect) return

    const nextX = Math.max(0, Math.min(rect.width, clientX - rect.left))
    setResizeGuideX(nextX)
  }

  const beginColumnResize = (
    event: ReactMouseEvent<HTMLButtonElement> | ReactTouchEvent<HTMLButtonElement>,
    resizeHandler: (event: unknown) => void,
  ) => {
    event.preventDefault()
    event.stopPropagation()

    if ('touches' in event) {
      const firstTouch = event.touches[0]
      if (firstTouch) updateResizeGuide(firstTouch.clientX)
    } else {
      updateResizeGuide(event.clientX)
    }

    const handleMouseMove = (moveEvent: MouseEvent) => updateResizeGuide(moveEvent.clientX)
    const handleTouchMove = (moveEvent: TouchEvent) => {
      const touch = moveEvent.touches[0]
      if (touch) updateResizeGuide(touch.clientX)
    }
    const clearGuide = () => {
      setResizeGuideX(null)
      window.removeEventListener('mousemove', handleMouseMove)
      window.removeEventListener('mouseup', clearGuide)
      window.removeEventListener('dragend', clearGuide)
      window.removeEventListener('blur', clearGuide)
      window.removeEventListener('contextmenu', clearGuide)
      window.removeEventListener('touchmove', handleTouchMove)
      window.removeEventListener('touchend', clearGuide)
      window.removeEventListener('touchcancel', clearGuide)
    }

    window.addEventListener('mousemove', handleMouseMove)
    window.addEventListener('mouseup', clearGuide)
    window.addEventListener('dragend', clearGuide)
    window.addEventListener('blur', clearGuide)
    window.addEventListener('contextmenu', clearGuide)
    window.addEventListener('touchmove', handleTouchMove)
    window.addEventListener('touchend', clearGuide)
    window.addEventListener('touchcancel', clearGuide)

    resizeHandler(event)
  }

  return (
    <TooltipProvider delayDuration={200}>
      <div
        ref={tableWrapperRef}
        aria-busy={loading}
        className={cn('relative w-full min-w-0 max-w-full overflow-hidden', className)}
      >
      {resizeGuideX !== null && (
        <div
          aria-hidden="true"
          className="pointer-events-none absolute top-0 bottom-0 z-30 w-px bg-primary shadow-[0_0_0_1px_hsl(var(--primary)/0.25)]"
          style={{ left: `${resizeGuideX}px` }}
        />
      )}
      <div className="w-full min-w-0 overflow-x-auto overflow-y-hidden">
        <Table
          className={cn('table-fixed', tableClassName)}
          style={{ minWidth: `${tableMinWidth}px`, width: `${tableMinWidth}px` }}
        >
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                const meta = header.column.columnDef.meta
                const fixedOffset = getPinnedOffset(orderedColumns, header.column.id, meta?.fixed)
                const canDrag = enableColumnDragging && meta?.draggable !== false
                const canResize = enableColumnResizing && header.column.getCanResize()
                const canSort = header.column.getCanSort()
                const sortDirection = header.column.getIsSorted()
                const headerAlign = meta?.headerAlign ?? meta?.align
                const renderedHeader = header.isPlaceholder
                  ? null
                  : flexRender(header.column.columnDef.header, header.getContext())
                const width = toCssSize(enableColumnResizing ? header.getSize() : meta?.width ?? header.getSize())
                const minWidth = toCssSize(meta?.minWidth)

                return (
                  <TableHead
                    key={header.id}
                    draggable={canDrag && resizeGuideX === null}
                    onDragStart={(event) => {
                      if (!canDrag || resizeGuideX !== null) {
                        event.preventDefault()
                        return
                      }
                      setDraggingColumnId(header.column.id)
                    }}
                    onDragOver={(event) => {
                      if (!canDrag) return
                      event.preventDefault()
                    }}
                    onDrop={() => {
                      if (draggingColumnId && canDrag) moveColumn(draggingColumnId, header.column.id)
                      setDraggingColumnId(null)
                    }}
                    onDragEnd={() => setDraggingColumnId(null)}
                    style={{
                      left: meta?.fixed === 'left' ? fixedOffset : undefined,
                      minWidth,
                      right: meta?.fixed === 'right' ? fixedOffset : undefined,
                      width,
                    }}
                    className={cn(
                      'relative select-none',
                      getTextAlignClass(headerAlign),
                      meta?.fixed && 'sticky z-20 bg-muted',
                      meta?.fixed === 'left' && 'left-0 shadow-sm',
                      meta?.fixed === 'right' && 'right-0 shadow-[-8px_0_12px_-12px_rgba(0,0,0,0.45)]',
                      draggingColumnId === header.column.id && 'opacity-60',
                      meta?.headerClassName,
                    )}
                  >
                    <div className={cn('flex items-center gap-2', headerAlign === 'right' && 'justify-end', headerAlign === 'center' && 'justify-center')}>
                      {canDrag && <GripVertical className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />}
                      {canSort ? (
                        <button
                          type="button"
                          onClick={header.column.getToggleSortingHandler()}
                          className="inline-flex min-w-0 items-center gap-1 rounded-sm text-left font-semibold hover:text-primary"
                        >
                          <span className="truncate">{renderedHeader}</span>
                          {sortDirection === 'asc'
                            ? <ChevronUp className="h-3.5 w-3.5 shrink-0" />
                            : sortDirection === 'desc'
                              ? <ChevronDown className="h-3.5 w-3.5 shrink-0" />
                              : <ChevronsUpDown className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />}
                        </button>
                      ) : (
                        <div className="min-w-0 truncate font-semibold">{renderedHeader}</div>
                      )}
                      {meta?.tooltip && (
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <span className="inline-flex h-5 w-5 items-center justify-center text-muted-foreground">
                              <Info className="h-3.5 w-3.5" />
                            </span>
                          </TooltipTrigger>
                          <TooltipContent>{meta.tooltip}</TooltipContent>
                        </Tooltip>
                      )}
                      {enableColumnFilters && header.column.getCanFilter() && (
                        <ColumnFilterMenu
                          column={header.column}
                          placeholder={meta?.filterPlaceholder}
                        />
                      )}
                    </div>
                    {canResize && (
                      <button
                        type="button"
                        aria-label="调整列宽"
                        draggable={false}
                        onDoubleClick={() => header.column.resetSize()}
                        onDragStart={(event) => event.preventDefault()}
                        onMouseDown={(event) => beginColumnResize(event, header.getResizeHandler())}
                        onTouchStart={(event) => beginColumnResize(event, header.getResizeHandler())}
                        className={cn(
                          'absolute right-0 top-0 h-full w-2 cursor-col-resize touch-none select-none',
                          'after:absolute after:right-0 after:top-1/2 after:h-6 after:w-px after:-translate-y-1/2 after:bg-border',
                          'hover:after:bg-primary',
                          header.column.getIsResizing() && 'after:bg-primary after:opacity-0',
                        )}
                      />
                    )}
                  </TableHead>
                )
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {shouldShowLoadingRows ? (
            Array.from({ length: 5 }).map((_, rowIndex) => (
              <TableRow key={`loading-${rowIndex}`}>
                {leafColumns.map((column, columnIndex) => {
                  const meta = column.columnDef.meta
                  const fixedOffset = getPinnedOffset(orderedColumns, column.id, meta?.fixed)
                  const width = toCssSize(enableColumnResizing ? column.getSize() : meta?.width ?? column.getSize())
                  const minWidth = toCssSize(meta?.minWidth)

                  return (
                    <TableCell
                      key={`${column.id}-${rowIndex}`}
                      style={{
                        left: meta?.fixed === 'left' ? fixedOffset : undefined,
                        minWidth,
                        right: meta?.fixed === 'right' ? fixedOffset : undefined,
                        width,
                      }}
                      className={cn(
                        getTextAlignClass(meta?.align),
                        meta?.fixed && 'sticky z-10 bg-background',
                        meta?.fixed === 'left' && 'left-0 shadow-sm',
                        meta?.fixed === 'right' && 'right-0 shadow-[-8px_0_12px_-12px_rgba(0,0,0,0.45)]',
                        meta?.className,
                      )}
                    >
                      <Skeleton className={cn('h-4', columnIndex === 0 ? 'w-4' : 'w-full')} />
                    </TableCell>
                  )
                })}
              </TableRow>
            ))
          ) : rows.length > 0 ? (
            rows.map((row) => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => {
                  const meta = cell.column.columnDef.meta
                  const fixedOffset = getPinnedOffset(orderedColumns, cell.column.id, meta?.fixed)
                  const width = toCssSize(enableColumnResizing ? cell.column.getSize() : meta?.width ?? cell.column.getSize())
                  const minWidth = toCssSize(meta?.minWidth)
                  const cellValue = cell.getValue()
                  const isEmpty = isEmptyCellValue(cellValue)
                  const copyable = meta?.copyable
                  const shouldCopy = Boolean(copyable) && !isEmpty && !meta?.loading
                  const renderedCell = meta?.loading
                    ? <Skeleton className="h-4 w-20" />
                    : isEmpty && meta?.emptyText !== undefined
                      ? meta.emptyText
                      : flexRender(cell.column.columnDef.cell, cell.getContext())

                  return (
                    <TableCell
                      key={cell.id}
                      style={{
                        left: meta?.fixed === 'left' ? fixedOffset : undefined,
                        minWidth,
                        right: meta?.fixed === 'right' ? fixedOffset : undefined,
                        width,
                      }}
                      className={cn(
                        getTextAlignClass(meta?.align),
                        meta?.fixed && 'sticky z-10 bg-background',
                        meta?.fixed === 'left' && 'left-0 shadow-sm',
                        meta?.fixed === 'right' && 'right-0 shadow-[-8px_0_12px_-12px_rgba(0,0,0,0.45)]',
                        meta?.className,
                      )}
                    >
                      <div className={cn('flex min-w-0 items-center gap-2', meta?.align === 'right' && 'justify-end', meta?.align === 'center' && 'justify-center')}>
                        <div className="min-w-0 truncate">{renderedCell}</div>
                        {shouldCopy && (
                          <Button
                            type="button"
                            aria-label="复制"
                            variant="ghost"
                            size="icon"
                            onClick={(event) => {
                              event.stopPropagation()
                              if (copyable) void navigator.clipboard?.writeText(getCopyText(copyable, cellValue, row.original))
                            }}
                            className="h-6 w-6 shrink-0 rounded-sm text-muted-foreground hover:bg-muted hover:text-foreground"
                          >
                            <Clipboard className="h-3.5 w-3.5" />
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  )
                })}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={leafColumns.length} className="h-24 text-center text-muted-foreground">
                {emptyText}
              </TableCell>
            </TableRow>
          )}
        </TableBody>
        </Table>
      </div>
      {shouldShowLoadingOverlay && (
        <div className="pointer-events-none absolute inset-0 z-40 flex items-center justify-center bg-background/60 backdrop-blur-[1px]">
          <div
            role="status"
            aria-live="polite"
            className="inline-flex items-center gap-2 rounded-md border bg-background px-3 py-2 text-sm text-muted-foreground shadow-sm"
          >
            <Loader2 className="h-4 w-4 animate-spin text-primary" />
            {loadingText}
          </div>
        </div>
      )}
      </div>
    </TooltipProvider>
  )
}
