import type { ThemePreference } from '~/types'

export const useTheme = () => {
  const colorMode = useColorMode()
  const appStore = useAppStore()

  const currentMode = computed(() => colorMode.preference as 'light' | 'dark' | 'system')
  const isDark = computed(() => colorMode.value === 'dark')

  const themePreference = computed<ThemePreference>(() => ({
    mode: colorMode.preference as 'light' | 'dark' | 'system',
    primaryColor: appStore.settings.primaryColor || 'emerald',
    accentColor: appStore.settings.accentColor || 'gold',
    borderRadius: appStore.settings.borderRadius || 'md',
  }))

  const setMode = (mode: 'light' | 'dark' | 'system'): void => {
    colorMode.preference = mode
    if (import.meta.client) {
      localStorage.setItem('color_mode', mode)
    }
  }

  const toggleMode = (): void => {
    const next = isDark.value ? 'light' : 'dark'
    setMode(next)
  }

  const setPrimaryColor = (color: string): void => {
    appStore.updateSettings({ primaryColor: color })
    if (import.meta.client) {
      localStorage.setItem('primary_color', color)
    }
  }

  const setAccentColor = (color: string): void => {
    appStore.updateSettings({ accentColor: color })
    if (import.meta.client) {
      localStorage.setItem('accent_color', color)
    }
  }

  const setBorderRadius = (radius: 'none' | 'sm' | 'md' | 'lg'): void => {
    appStore.updateSettings({ borderRadius: radius })
  }

  const isSystemDark = (): boolean => {
    if (!import.meta.client) return false
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }

  const initTheme = (): void => {
    if (import.meta.client) {
      const savedMode = localStorage.getItem('color_mode')
      if (savedMode && ['light', 'dark', 'system'].includes(savedMode)) {
        colorMode.preference = savedMode as 'light' | 'dark' | 'system'
      }

      const savedPrimary = localStorage.getItem('primary_color')
      if (savedPrimary) {
        appStore.updateSettings({ primaryColor: savedPrimary })
      }

      const savedAccent = localStorage.getItem('accent_color')
      if (savedAccent) {
        appStore.updateSettings({ accentColor: savedAccent })
      }
    }
  }

  const applyColorVariables = (): void => {
    if (!import.meta.client) return

    const root = document.documentElement
    const primary = themePreference.value.primaryColor

    const colorMap: Record<string, { h: string; s: string; l: string }> = {
      emerald: { h: '160', s: '84%', l: '39%' },
      blue: { h: '217', s: '91%', l: '60%' },
      indigo: { h: '239', s: '84%', l: '67%' },
      purple: { h: '271', s: '81%', l: '56%' },
      rose: { h: '350', s: '89%', l: '60%' },
    }

    const c = colorMap[primary] || colorMap.emerald
    root.style.setProperty('--color-primary', `${c.h} ${c.s} ${c.l}`)
  }

  return {
    currentMode,
    isDark,
    themePreference,
    setMode,
    toggleMode,
    setPrimaryColor,
    setAccentColor,
    setBorderRadius,
    isSystemDark,
    initTheme,
    applyColorVariables,
  }
}
