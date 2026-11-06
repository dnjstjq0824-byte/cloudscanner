export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()
  if (auth.token() && !auth.isAuthenticated && !auth.hydrated) {
    await auth.fetchMe()
  }
  if (auth.isAuthenticated || auth.token()) {
    return navigateTo('/dashboard')
  }
})
