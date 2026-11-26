// RETIRED 2026-11-26 (TirpScanner pivot): locale-aware 301s for removed pages.
// Unprefixed legacy URLs resolve to the default locale.
// NOTE 2026-11-26: search/calculator/rankings/providers are live pages again
// (live 44종), so RETIRED_HOME stays empty. Do not re-add without removing pages.
const LOCALES = new Set(['en', 'ko', 'ja', 'de', 'es'])
const RETIRED_HOME = new Set<string>([])
const RETIRED_CONSOLE = new Set(['connections', 'servers', 'provision', 'budgets'])

export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook('request', (event) => {
    const clean = event.path.split('?')[0].replace(/\/+$/, '')
    const segments = clean.split('/').filter(Boolean)
    let locale = 'en'
    let rest = segments
    if (segments.length > 0 && LOCALES.has(segments[0])) {
      locale = segments[0] as string
      rest = segments.slice(1)
    }
    if (rest.length === 0) {
      return
    }
    if (RETIRED_HOME.has(rest[0])) {
      sendRedirect(event, `/${locale}/`, 301)
    } else if (rest[0] === 'dashboard' && rest.length > 1 && RETIRED_CONSOLE.has(rest[1])) {
      sendRedirect(event, `/${locale}/dashboard`, 301)
    }
  })
})
