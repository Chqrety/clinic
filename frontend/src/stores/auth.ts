import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface AuthUser {
  id: number
  employee_number: string
  name: string
  role: 'SUPERADMIN_OWNER' | 'DOCTOR' | 'GRO_CASHIER' | 'PHARMACIST'
  clinic_unit: 'KIDINEST' | 'SSC' | 'ALL'
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('clinic_auth_token') || '')
  const user = ref<AuthUser | null>(
    localStorage.getItem('clinic_auth_user')
      ? JSON.parse(localStorage.getItem('clinic_auth_user')!)
      : null
  )
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isOwner = computed(() => user.value?.role === 'SUPERADMIN_OWNER')
  const isDoctor = computed(() => user.value?.role === 'DOCTOR')
  const isCashier = computed(() => user.value?.role === 'GRO_CASHIER')
  const isPharmacist = computed(() => user.value?.role === 'PHARMACIST')

  function setAuth(newToken: string, newUser: AuthUser) {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('clinic_auth_token', newToken)
    localStorage.setItem('clinic_auth_user', JSON.stringify(newUser))
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    localStorage.removeItem('clinic_auth_token')
    localStorage.removeItem('clinic_auth_user')
  }

  async function login(employeeNumber: string, password: string): Promise<boolean> {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ employee_number: employeeNumber, password }),
      })

      const data = await res.json()
      if (!res.ok || !data.success) {
        error.value = data.error || 'Login gagal, periksa NIK Karyawan dan password.'
        return false
      }

      setAuth(data.data.token, data.data.user)
      return true
    } catch (err: any) {
      error.value = 'Gagal terhubung ke server backend: ' + err.message
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMe(): Promise<boolean> {
    if (!token.value) return false
    try {
      const res = await fetch('/api/auth/me', {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
      const data = await res.json()
      if (res.ok && data.success) {
        user.value = data.data
        localStorage.setItem('clinic_auth_user', JSON.stringify(data.data))
        return true
      }
      clearAuth()
      return false
    } catch {
      return false
    }
  }

  async function logout() {
    try {
      if (token.value) {
        await fetch('/api/auth/logout', {
          method: 'POST',
          headers: { Authorization: `Bearer ${token.value}` },
        })
      }
    } catch {
      // Ignored
    } finally {
      clearAuth()
    }
  }

  return {
    token,
    user,
    isLoading,
    error,
    isAuthenticated,
    isOwner,
    isDoctor,
    isCashier,
    isPharmacist,
    setAuth,
    clearAuth,
    login,
    fetchMe,
    logout,
  }
})
