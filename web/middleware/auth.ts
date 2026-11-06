export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()
  if (!auth.hydrated) {
    await auth.fetchMe()
  } else if (!auth.isAuthenticated && auth.token()) {
    await auth.fetchMe()
  }
  if (!auth.isAuthenticated) {
    return navigateTo('/auth/login')
  }
})
