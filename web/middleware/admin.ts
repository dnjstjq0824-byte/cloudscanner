export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()
  if (!auth.hydrated) {
    await auth.fetchMe()
  }
  if (!auth.isAuthenticated) {
    return navigateTo(useLocalePath()('/auth/login'))
  }
  if (!auth.isAdmin) {
    return navigateTo(useLocalePath()('/dashboard'))
  }
})
