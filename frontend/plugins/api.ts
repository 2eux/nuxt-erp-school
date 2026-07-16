import { toast } from 'vue-sonner'
import type { ApiClient } from '~/types/api'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const $api = $fetch.create({
    baseURL: config.public.apiBase as string,
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
    onRequest({ options }) {
      const token = authStore.accessToken
      if (token) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${token}`,
        }
      }
    },
    onResponseError({ response }) {
      const status = response.status
      const data = response._data as { message?: string; error?: string }

      switch (status) {
        case 401: {
          const currentPath = window.location.pathname
          if (!currentPath.startsWith('/auth/')) {
            authStore.clearAuth()
            navigateTo('/auth/login')
          }
          break
        }
        case 403:
          toast.error(data?.message || 'Forbidden access', {
            description: 'You do not have permission to perform this action.',
          })
          break
        case 422:
          toast.error(data?.message || 'Validation Error', {
            description: data?.error || 'Please check your input.',
          })
          break
        case 429:
          toast.error('Too many requests', {
            description: 'Please wait a moment before trying again.',
          })
          break
        case 500:
        case 502:
        case 503:
          toast.error('Server Error', {
            description: data?.message || 'An unexpected error occurred. Please try again later.',
          })
          break
        default:
          break
      }
    },
    onRequestError() {
      toast.error('Network Error', {
        description: 'Unable to connect to the server. Please check your connection.',
      })
    },
  })

  return {
    provide: {
      apiFetch: $api,
    },
  }
})
