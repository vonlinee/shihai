import { useState, useRef, useEffect } from 'react'
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
}

export function Combobox({
  options,
  value,
  onChange,
  placeholder = '请选择...',
  allowCustom = true,
  customPlaceholder = '或直接输入...',
  className = '',
}: ComboboxProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [inputValue, setInputValue] = useState('')
  const [isCustom, setIsCustom] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  // Find the selected option label
  const selectedOption = options.find((o) => String(o.value) === String(value))
  const displayValue = isCustom ? inputValue : (selectedOption?.label || '')

  // Filter options by input
  const filteredOptions = options.filter((o) =>
    o.label.toLowerCase().includes(inputValue.toLowerCase())
  )

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

  const handleSelect = (option: ComboboxOption) => {
    setIsCustom(false)
    setInputValue('')
    onChange(option.value, option)
    setIsOpen(false)
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value
    setInputValue(val)
    setIsCustom(true)
    onChange(val) // pass raw string as custom value
    if (!isOpen) setIsOpen(true)
  }

  const handleInputFocus = () => {
    setIsOpen(true)
    if (!isCustom && selectedOption) {
      setInputValue('')
    }
  }

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation()
    setIsCustom(false)
    setInputValue('')
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
            onClick={() => setIsOpen(!isOpen)}
            className="p-0.5 hover:bg-muted rounded"
          >
            <ChevronDown className={`h-3.5 w-3.5 text-muted-foreground transition-transform ${isOpen ? 'rotate-180' : ''}`} />
          </button>
        </div>
      </div>

      {isOpen && (
        <div className="absolute z-50 mt-1 w-full rounded-md border bg-popover shadow-md max-h-60 overflow-y-auto">
          {filteredOptions.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-6 px-3 text-muted-foreground">
              <SearchX className="h-8 w-8 mb-2 opacity-40" />
              <p className="text-sm">{inputValue ? '无匹配结果' : '暂无数据'}</p>
              {inputValue && (
                <p className="text-xs mt-1">尝试其他关键词{allowCustom ? '，或直接输入' : ''}</p>
              )}
            </div>
          ) : (
            filteredOptions.map((option) => (
              <button
                key={String(option.value)}
                type="button"
                onClick={() => handleSelect(option)}
                className={`w-full text-left px-3 py-2 text-sm hover:bg-muted/50 transition-colors ${
                  String(option.value) === String(value) && !isCustom
                    ? 'bg-primary/5 font-medium'
                    : ''
                }`}
              >
                <span>{option.label}</span>
                {option.description && (
                  <span className="ml-2 text-xs text-muted-foreground">{option.description}</span>
                )}
              </button>
            ))
          )}
          {allowCustom && inputValue.trim() && filteredOptions.length > 0 && !filteredOptions.find(
            (o) => o.label.toLowerCase() === inputValue.toLowerCase()
          ) && (
            <div className="px-3 py-2 text-xs text-muted-foreground border-t">
              回车确认输入「{inputValue}」
            </div>
          )}
          {allowCustom && inputValue.trim() && filteredOptions.length === 0 && (
            <div className="px-3 py-2.5 text-xs text-primary border-t font-medium">
              回车确认创建「{inputValue}」
            </div>
          )}
        </div>
      )}
    </div>
  )
}
