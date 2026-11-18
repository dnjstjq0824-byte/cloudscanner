/** Persisted light/dark theme. Defaults to dark; applied pre-paint by head script. */
export function useTheme() {
  const isDark = useState('cs-theme-dark', () => true)

  function apply(dark: boolean) {
    isDark.value = dark
    if (import.meta.client) {
      document.documentElement.classList.toggle('dark', dark)
      try {
        localStorage.setItem('cs-theme', dark ? 'dark' : 'light')
      } catch {
        // Private mode: keep the in-memory choice only.
      }
    }
  }

  function toggle() {
    apply(!isDark.value)
  }

  function init() {
    if (import.meta.client) {
      apply(document.documentElement.classList.contains('dark'))
    }
  }

  return { isDark, apply, toggle, init }
}
