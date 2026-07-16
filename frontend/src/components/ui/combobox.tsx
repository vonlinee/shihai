import { useState, useRef, useEffect, useMemo } from 'react'
import { ChevronDown, X, SearchX } from 'lucide-react'

export interface ComboboxOption {
  value: string | number
  label: string
  description?: string
}

interface ComboboxProps {
  options: ComboboxOption[]
  value: string | number | undefined
  onChange: (value: string | number | undefined, option?: ComboboxOption) => void
  placeholder?: string
  allowCustom?: boolean
  customPlaceholder?: string
  className?: string
  isLoading?: boolean
  virtualItemHeight?: number
  virtualListHeight?: number
  virtualOverscan?: number
  onSearchChange?: (keyword: string) => void
  groupOptions?: ComboboxOption[]
  groupValue?: string | number | undefined
  onGroupChange?: (value: string | number, option: ComboboxOption) => void
  groupPlaceholder?: string
}

export function Combobox({
  options,
  value,
  onChange,
  placeholder = '请选择...',
  allowCustom = true,
  customPlaceholder = '或直接输入...',
  className = '',
  isLoading = false,
  virtualItemHeight = 40,
  virtualListHeight = 240,
  virtualOverscan = 4,
  onSearchChange,
  groupOptions,
  groupValue,
  onGroupChange,
  groupPlaceholder = '选择分类',
}: ComboboxProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [inputValue, setInputValue] = useState('')
  const [isCustom, setIsCustom] = useState(false)
  const [listScrollTop, setListScrollTop] = useState(0)
  const containerRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const listRef = useRef<HTMLDivElement>(null)
  const hasGroups = Array.isArray(groupOptions)
  const selectedGroupOption = groupOptions?.find((o) => String(o.value) === String(groupValue))

  // Find the selected option label
  const selectedOption = options.find((o) => String(o.value) === String(value))
  const displayValue = isCustom ? inputValue : (selectedOption?.label || '')

  // Filter options by input
  const filteredOptions = useMemo(
    () => options.filter((o) =>
      o.label.toLowerCase().includes(inputValue.toLowerCase())
    ),
    [inputValue, options]
  )
  const listHeight = Math.min(virtualListHeight, filteredOptions.length * virtualItemHeight)
  const optionListViewportHeight = hasGroups ? Math.max(0, virtualListHeight - virtualItemHeight) : listHeight
  const virtualStartIndex = Math.max(0, Math.floor(listScrollTop / virtualItemHeight) - virtualOverscan)
  const virtualEndIndex = Math.min(
    filteredOptions.length,
    Math.ceil((listScrollTop + optionListViewportHeight) / virtualItemHeight) + virtualOverscan
  )
  const virtualOptions = useMemo(
    () => filteredOptions.slice(virtualStartIndex, virtualEndIndex).map((option, offset) => ({
      option,
      index: virtualStartIndex + offset,
    })),
    [filteredOptions, virtualEndIndex, virtualStartIndex]
  )
  const virtualListTotalHeight = filteredOptions.length * virtualItemHeight

  // Close on click outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  useEffect(() => {
    setListScrollTop(0)
    listRef.current?.scrollTo({ top: 0 })
  }, [inputValue, options])

  const handleSelect = (option: ComboboxOption) => {
    setIsCustom(false)
    setInputValue('')
    onSearchChange?.('')
    onChange(option.value, option)
    setIsOpen(false)
  }

  const handleGroupSelect = (option: ComboboxOption) => {
    setIsCustom(false)
    setInputValue('')
    onSearchChange?.('')
    onGroupChange?.(option.value, option)
    setListScrollTop(0)
    listRef.current?.scrollTo({ top: 0 })
    inputRef.current?.focus()
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value
    setInputValue(val)
    setIsCustom(true)
    onSearchChange?.(val)
    onChange(val) // pass raw string as custom value
    if (!isOpen) setIsOpen(true)
  }

  const handleInputFocus = () => {
    setIsOpen(true)
    onSearchChange?.(inputValue)
    if (!isCustom && selectedOption) {
      setInputValue('')
      onSearchChange?.('')
    }
  }

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation()
    setIsCustom(false)
    setInputValue('')
    onSearchChange?.('')
    onChange(undefined)
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      setIsOpen(false)
    }
    if (e.key === 'Enter' && isOpen) {
      // If exact match, select it
      const exactMatch = filteredOptions.find(
        (o) => o.label.toLowerCase() === inputValue.toLowerCase()
      )
      if (exactMatch) {
        handleSelect(exactMatch)
      } else if (allowCustom && inputValue.trim()) {
        // Custom value already set via onChange
        setIsOpen(false)
      }
    }
  }

  return (
    <div className={`relative ${className}`} ref={containerRef}>
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          value={isOpen && !isCustom ? inputValue : displayValue}
          onChange={handleInputChange}
          onFocus={handleInputFocus}
          onKeyDown={handleKeyDown}
          placeholder={isCustom ? customPlaceholder : placeholder}
          className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm pr-8 focus:outline-none focus:ring-1 focus:ring-ring"
        />
        <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-0.5">
          {value !== undefined && value !== '' && (
            <button
              type="button"
              onClick={handleClear}
              className="p-0.5 hover:bg-muted rounded"
            >
              <X className="h-3 w-3 text-muted-foreground" />
            </button>
          )}
          <button
            type="button"
            onClick={() => {
              const nextOpen = !isOpen
              setIsOpen(nextOpen)
              if (nextOpen && hasGroups) {
                setInputValue('')
                onSearchChange?.('')
              }
            }}
            className="p-0.5 hover:bg-muted rounded"
          >
            <ChevronDown className={`h-3.5 w-3.5 text-muted-foreground transition-transform ${isOpen ? 'rotate-180' : ''}`} />
          </button>
        </div>
      </div>

      {isOpen && (
        <div className="absolute z-50 mt-1 w-full rounded-md border bg-popover shadow-md">
          {hasGroups ? (
            <div className="grid grid-cols-[minmax(104px,38%)_minmax(0,62%)]" style={{ height: virtualListHeight }}>
              <div className="min-w-0 border-r">
                <div className="border-b px-3 py-2 text-xs text-muted-foreground">
                  {groupPlaceholder}
                </div>
                <div className="overflow-y-auto" style={{ height: virtualListHeight - virtualItemHeight }}>
                  {(groupOptions ?? []).length === 0 ? (
                    <div className="px-3 py-6 text-center text-sm text-muted-foreground">暂无数据</div>
                  ) : (groupOptions ?? []).map((option) => (
                  <button
                    key={String(option.value)}
                    type="button"
                    onClick={() => handleGroupSelect(option)}
                    className={`flex w-full items-center px-3 py-2 text-left text-sm hover:bg-muted/50 transition-colors ${
                      String(option.value) === String(groupValue)
                        ? 'bg-primary/5 font-medium'
                        : ''
                    }`}
                  >
                    <span className="truncate">{option.label}</span>
                    {option.description && (
                      <span className="ml-2 shrink-0 text-xs text-muted-foreground">{option.description}</span>
                    )}
                  </button>
                  ))}
                </div>
              </div>
              <div className="min-w-0">
                <div className="border-b px-3 py-2 text-xs text-muted-foreground">
                  {selectedGroupOption?.label ?? '请选择朝代'}
                </div>
                {!groupValue ? (
                  <div className="flex flex-col items-center justify-center px-3 text-center text-sm text-muted-foreground" style={{ height: virtualListHeight - virtualItemHeight }}>
                    <SearchX className="mb-2 h-8 w-8 opacity-40" />
                    <p>请先选择朝代</p>
                  </div>
                ) : isLoading && filteredOptions.length === 0 ? (
                  <div className="px-3 py-6 text-center text-sm text-muted-foreground">
                    加载中...
                  </div>
                ) : filteredOptions.length === 0 ? (
                  <div className="flex flex-col items-center justify-center px-3 text-muted-foreground" style={{ height: virtualListHeight - virtualItemHeight }}>
                    <SearchX className="h-8 w-8 mb-2 opacity-40" />
                    <p className="text-sm">{inputValue ? '无匹配结果' : '暂无数据'}</p>
                  </div>
                ) : (
                  <div
                    ref={listRef}
                    className="overflow-y-auto"
                    style={{ height: optionListViewportHeight }}
                    onScroll={(event) => setListScrollTop(event.currentTarget.scrollTop)}
                  >
                    <div className="relative" style={{ height: virtualListTotalHeight }}>
                      {virtualOptions.map(({ option, index }) => (
                        <button
                          key={String(option.value)}
                          type="button"
                          onClick={() => handleSelect(option)}
                          className={`absolute left-0 right-0 flex w-full items-center text-left px-3 text-sm hover:bg-muted/50 transition-colors ${
                            String(option.value) === String(value) && !isCustom
                              ? 'bg-primary/5 font-medium'
                              : ''
                          }`}
                          style={{
                            height: virtualItemHeight,
                            transform: `translateY(${index * virtualItemHeight}px)`,
                          }}
                        >
                          <span className="truncate">{option.label}</span>
                          {option.description && (
                            <span className="ml-2 shrink-0 text-xs text-muted-foreground">{option.description}</span>
                          )}
                        </button>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </div>
          ) : (
          <>
          {isLoading && filteredOptions.length === 0 ? (
            <div className="px-3 py-6 text-center text-sm text-muted-foreground">
              加载中...
            </div>
          ) : filteredOptions.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-6 px-3 text-muted-foreground">
              <SearchX className="h-8 w-8 mb-2 opacity-40" />
              <p className="text-sm">{inputValue ? '无匹配结果' : '暂无数据'}</p>
              {inputValue && (
                <p className="text-xs mt-1">尝试其他关键词{allowCustom ? '，或直接输入' : ''}</p>
              )}
            </div>
          ) : (
            <div
              ref={listRef}
              className="overflow-y-auto"
              style={{ height: listHeight }}
              onScroll={(event) => setListScrollTop(event.currentTarget.scrollTop)}
            >
              <div className="relative" style={{ height: virtualListTotalHeight }}>
                {virtualOptions.map(({ option, index }) => (
                  <button
                    key={String(option.value)}
                    type="button"
                    onClick={() => handleSelect(option)}
                    className={`absolute left-0 right-0 flex w-full items-center text-left px-3 text-sm hover:bg-muted/50 transition-colors ${
                      String(option.value) === String(value) && !isCustom
                        ? 'bg-primary/5 font-medium'
                        : ''
                    }`}
                    style={{
                      height: virtualItemHeight,
                      transform: `translateY(${index * virtualItemHeight}px)`,
                    }}
                  >
                    <span className="truncate">{option.label}</span>
                    {option.description && (
                      <span className="ml-2 shrink-0 text-xs text-muted-foreground">{option.description}</span>
                    )}
                  </button>
                ))}
              </div>
            </div>
          )}
          </>
          )}
          {!hasGroups && allowCustom && inputValue.trim() && filteredOptions.length > 0 && !filteredOptions.find(
            (o) => o.label.toLowerCase() === inputValue.toLowerCase()
          ) && (
            <div className="px-3 py-2 text-xs text-muted-foreground border-t">
              回车确认输入「{inputValue}」
            </div>
          )}
          {!hasGroups && allowCustom && inputValue.trim() && filteredOptions.length === 0 && (
            <div className="px-3 py-2.5 text-xs text-primary border-t font-medium">
              回车确认创建「{inputValue}」
            </div>
          )}
        </div>
      )}
    </div>
  )
}
