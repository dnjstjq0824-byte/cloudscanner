export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()
  if (!auth.hydrated) {
    await auth.fetchMe()
  }
  if (!auth.isAuthenticated) {
    return navigateTo('/auth/login')
  }
  if (!auth.isAdmin) {
    return navigateTo('/dashboard')
  }
})
