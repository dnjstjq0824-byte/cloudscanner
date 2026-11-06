import { defineStore } from 'pinia'
import type { AuthOrg, AuthUser, LoginResponse, MeResponse, RegisterResponse } from '~/types/api'

interface AuthState {
  user: AuthUser | null
  org: AuthOrg | null
  role: string
  expiresAt: string
  hydrated: boolean
}

/** Session store. The token lives in an HTTP-visible cookie so SSR requests authenticate. */
export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    org: null,
    role: '',
    expiresAt: '',
    hydrated: false,
  }),

  getters: {
    isAuthenticated: (s) => s.user !== null,
    isAdmin: (s) => s.role === 'OWNER' || s.role === 'ADMIN',
  },

  actions: {
    token(): string | null | undefined {
      return useCookie('cs_session').value
    },

    /** Secure cookies only on HTTPS origins (SSR-safe). */
    isHttpsOrigin(): boolean {
      try {
        return useRequestURL().protocol === 'https:'
      } catch {
        return false
      }
    },

    applyLogin(res: LoginResponse) {
      this.user = res.user
      this.org = { id: res.org_id, name: '', slug: '', owner_user_id: res.user.id, created_at: '' }
      this.role = res.role
      this.expiresAt = res.expires_at
      useCookie('cs_session', { maxAge: 60 * 15, sameSite: 'lax', secure: this.isHttpsOrigin() }).value = res.token
      useCookie('cs_org', { maxAge: 60 * 60 * 24 * 7, sameSite: 'lax' }).value = res.org_id
      this.hydrated = true
    },

    applyRegister(res: RegisterResponse) {
      this.user = res.user
      this.org = res.organization
      this.role = 'OWNER'
      this.expiresAt = res.expires_at
      useCookie('cs_session', { maxAge: 60 * 15, sameSite: 'lax', secure: this.isHttpsOrigin() }).value = res.token
      useCookie('cs_org', { maxAge: 60 * 60 * 24 * 7, sameSite: 'lax' }).value = res.organization.id
      this.hydrated = true
    },

    async login(email: string, password: string) {
      const { request } = useApi()
      const res = await request<LoginResponse>('/api/v1/auth/login', {
        method: 'POST',
        body: { email, password },
        anonymous: true,
      })
      this.applyLogin(res)
    },

    async register(email: string, password: string, fullName: string) {
      const { request } = useApi()
      const res = await request<RegisterResponse>('/api/v1/auth/register', {
        method: 'POST',
        body: { email, password, full_name: fullName },
        anonymous: true,
      })
      this.applyRegister(res)
    },

    /** Rehydrate the session on page load when a token cookie exists. */
    async fetchMe(): Promise<boolean> {
      if (!this.token()) {
        this.hydrated = true
        return false
      }
      try {
        const { request } = useApi()
        const me = await request<MeResponse>('/api/v1/auth/me')
        this.user = me.user
        this.org = me.organization
        this.role = me.role
        this.expiresAt = me.expires_at
        this.hydrated = true
        return true
      } catch {
        this.clear()
        return false
      }
    },

    logout() {
      this.clear()
      navigateTo('/')
    },

    clear() {
      this.user = null
      this.org = null
      this.role = ''
      this.expiresAt = ''
      this.hydrated = true
      useCookie('cs_session').value = null
      // Keep org cookie for UX; it is only sent with a valid token.
    },
  },
})
