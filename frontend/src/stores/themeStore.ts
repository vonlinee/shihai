import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type ThemeName = 'light' | 'dark' | 'ink-wash' | 'bamboo' | 'cinnabar'

interface ThemeState {
  theme: ThemeName
  setTheme: (theme: ThemeName) => void
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      theme: 'light',
      setTheme: (theme) => {
        set({ theme })
        applyTheme(theme)
      },
    }),
    {
      name: 'theme-storage',
      onRehydrateStorage: () => {
        return (state) => {
          if (state) {
            applyTheme(state.theme)
          }
        }
      },
    }
  )
)

// Apply theme CSS variables to :root
export function applyTheme(theme: ThemeName) {
  const root = document.documentElement
  const vars = themes[theme].variables
  Object.entries(vars).forEach(([key, value]) => {
    root.style.setProperty(key, value)
  })
  // Toggle dark class for components that rely on it
  root.classList.toggle('dark', theme === 'dark')
}

export interface ThemePreset {
  name: ThemeName
  label: string
  description: string
  icon: string
  variables: Record<string, string>
}

export const themes: Record<ThemeName, ThemePreset> = {
  light: {
    name: 'light',
    label: '素雅',
    description: '淡雅明亮的经典风格',
    icon: '☀️',
    variables: {
      '--background': '0 0% 98%',
      '--foreground': '0 0% 17%',
      '--card': '30 30% 97%',
      '--card-foreground': '0 0% 17%',
      '--popover': '0 0% 100%',
      '--popover-foreground': '0 0% 17%',
      '--primary': '25 76% 31%',
      '--primary-foreground': '30 30% 97%',
      '--secondary': '30 35% 71%',
      '--secondary-foreground': '25 76% 20%',
      '--muted': '30 20% 92%',
      '--muted-foreground': '0 0% 40%',
      '--accent': '350 82% 44%',
      '--accent-foreground': '0 0% 100%',
      '--destructive': '0 84% 60%',
      '--destructive-foreground': '0 0% 100%',
      '--border': '30 20% 85%',
      '--input': '30 20% 85%',
      '--ring': '25 76% 31%',
    },
  },
  dark: {
    name: 'dark',
    label: '暗夜',
    description: '深色护眼模式',
    icon: '🌙',
    variables: {
      '--background': '0 0% 7%',
      '--foreground': '0 0% 92%',
      '--card': '0 0% 10%',
      '--card-foreground': '0 0% 92%',
      '--popover': '0 0% 10%',
      '--popover-foreground': '0 0% 92%',
      '--primary': '25 60% 50%',
      '--primary-foreground': '0 0% 100%',
      '--secondary': '0 0% 18%',
      '--secondary-foreground': '0 0% 85%',
      '--muted': '0 0% 15%',
      '--muted-foreground': '0 0% 55%',
      '--accent': '350 70% 50%',
      '--accent-foreground': '0 0% 100%',
      '--destructive': '0 72% 51%',
      '--destructive-foreground': '0 0% 100%',
      '--border': '0 0% 20%',
      '--input': '0 0% 20%',
      '--ring': '25 60% 50%',
    },
  },
  'ink-wash': {
    name: 'ink-wash',
    label: '水墨',
    description: '黑白灰的古典水墨风',
    icon: '🖌️',
    variables: {
      '--background': '210 10% 95%',
      '--foreground': '0 0% 12%',
      '--card': '210 8% 93%',
      '--card-foreground': '0 0% 12%',
      '--popover': '210 10% 96%',
      '--popover-foreground': '0 0% 12%',
      '--primary': '0 0% 25%',
      '--primary-foreground': '0 0% 95%',
      '--secondary': '0 0% 75%',
      '--secondary-foreground': '0 0% 15%',
      '--muted': '0 0% 88%',
      '--muted-foreground': '0 0% 40%',
      '--accent': '0 0% 30%',
      '--accent-foreground': '0 0% 98%',
      '--destructive': '0 72% 45%',
      '--destructive-foreground': '0 0% 98%',
      '--border': '0 0% 80%',
      '--input': '0 0% 80%',
      '--ring': '0 0% 25%',
    },
  },
  bamboo: {
    name: 'bamboo',
    label: '竹青',
    description: '清新翠竹的自然风',
    icon: '🎋',
    variables: {
      '--background': '140 20% 96%',
      '--foreground': '140 20% 12%',
      '--card': '140 15% 94%',
      '--card-foreground': '140 20% 12%',
      '--popover': '140 18% 97%',
      '--popover-foreground': '140 20% 12%',
      '--primary': '150 50% 30%',
      '--primary-foreground': '140 20% 97%',
      '--secondary': '140 20% 70%',
      '--secondary-foreground': '150 50% 18%',
      '--muted': '140 15% 90%',
      '--muted-foreground': '140 10% 38%',
      '--accent': '80 50% 38%',
      '--accent-foreground': '0 0% 100%',
      '--destructive': '0 72% 45%',
      '--destructive-foreground': '0 0% 100%',
      '--border': '140 15% 82%',
      '--input': '140 15% 82%',
      '--ring': '150 50% 30%',
    },
  },
  cinnabar: {
    name: 'cinnabar',
    label: '丹砂',
    description: '朱砂红的古典华丽风',
    icon: '🏯',
    variables: {
      '--background': '20 30% 96%',
      '--foreground': '10 30% 12%',
      '--card': '20 25% 93%',
      '--card-foreground': '10 30% 12%',
      '--popover': '20 28% 97%',
      '--popover-foreground': '10 30% 12%',
      '--primary': '350 75% 40%',
      '--primary-foreground': '20 30% 98%',
      '--secondary': '20 25% 72%',
      '--secondary-foreground': '350 60% 20%',
      '--muted': '20 15% 90%',
      '--muted-foreground': '10 15% 40%',
      '--accent': '40 70% 45%',
      '--accent-foreground': '20 30% 98%',
      '--destructive': '0 72% 45%',
      '--destructive-foreground': '0 0% 100%',
      '--border': '20 18% 82%',
      '--input': '20 18% 82%',
      '--ring': '350 75% 40%',
    },
  },
}
