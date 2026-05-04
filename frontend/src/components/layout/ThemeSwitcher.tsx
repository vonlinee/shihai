import { useState } from 'react'
import { Palette, Check, ChevronDown, ChevronUp } from 'lucide-react'
import { useThemeStore, themes } from '@/stores/themeStore'

export function ThemeSwitcher() {
  const { theme: currentTheme, setTheme } = useThemeStore()
  const [isExpanded, setIsExpanded] = useState(false)

  const themeList = Object.values(themes)

  return (
    <div>
      <button
        onClick={(e) => {
          e.stopPropagation()
          setIsExpanded(!isExpanded)
        }}
        className="flex items-center gap-3 px-4 py-2 text-sm w-full text-left hover:bg-muted/50 transition-colors"
      >
        <Palette className="h-4 w-4 text-muted-foreground" />
        <span className="flex-1">主题风格</span>
        <span className="text-xs">{themes[currentTheme].icon} {themes[currentTheme].label}</span>
        {isExpanded ? (
          <ChevronUp className="h-3 w-3 text-muted-foreground" />
        ) : (
          <ChevronDown className="h-3 w-3 text-muted-foreground" />
        )}
      </button>

      {isExpanded && (
        <div className="pb-1">
          {themeList.map((t) => (
            <button
              key={t.name}
              onClick={(e) => {
                e.stopPropagation()
                setTheme(t.name)
              }}
              className={`flex items-center gap-3 px-4 py-2 text-sm w-full text-left transition-colors ${
                currentTheme === t.name ? 'bg-primary/5' : 'hover:bg-muted/50'
              }`}
            >
              {/* Color preview swatches */}
              <div className="flex -space-x-1">
                <span
                  className="inline-block h-4 w-4 rounded-full border border-border"
                  style={{ backgroundColor: `hsl(${t.variables['--primary']})` }}
                />
                <span
                  className="inline-block h-4 w-4 rounded-full border border-border"
                  style={{ backgroundColor: `hsl(${t.variables['--background']})` }}
                />
                <span
                  className="inline-block h-4 w-4 rounded-full border border-border"
                  style={{ backgroundColor: `hsl(${t.variables['--accent']})` }}
                />
              </div>

              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-1.5">
                  <span className="text-sm">{t.icon}</span>
                  <span className="font-medium">{t.label}</span>
                </div>
                <p className="text-xs text-muted-foreground truncate">{t.description}</p>
              </div>

              {currentTheme === t.name && (
                <Check className="h-4 w-4 text-primary shrink-0" />
              )}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
