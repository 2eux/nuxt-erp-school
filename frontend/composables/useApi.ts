import type { ApiClient, RequestConfig, ApiError } from '~/types/api'

const createApiClient = (): ApiClient => {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()
  const { t } = useI18n()
  const toast = useToast()

  const baseURL = config.public.apiBase as string

  const getHeaders = (): Record<string, string> => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    }

    const token = authStore.accessToken
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    return headers
  }

  const handleError = (error: unknown): never => {
    const apiError = error as ApiError & { data?: ApiError; response?: { status: number; _data?: ApiError } }

    if (apiError.response) {
      const status = apiError.response.status
      const data = apiError.response._data || apiError.data

      if (status === 401) {
        authStore.clearAuth()
        navigateTo('/auth/login')
        throw new Error('Unauthorized')
      }

      if (status === 403) {
        toast.add({
          title: t('error.forbidden'),
          description: t('error.forbidden_description'),
          color: 'error',
        })
        throw new Error('Forbidden')
      }

      if (status === 422 && data?.details) {
        throw { statusCode: 422, message: data.message || 'Validation Error', details: data.details }
      }

      if (status >= 500) {
        toast.add({
          title: t('error.server_error'),
          description: data?.message || t('error.server_error_description'),
          color: 'error',
        })
        throw error
      }

      throw error
    }

    if (apiError instanceof TypeError && apiError.message === 'Failed to fetch') {
      toast.add({
        title: t('error.network_error'),
        description: t('error.network_error_description'),
        color: 'error',
      })
    }

    throw error
  }

  const request = async <T>(url: string, config: RequestConfig = {}): Promise<T> => {
    const { method = 'GET', params, body, headers: extraHeaders } = config

    try {
      const queryString = params
        ? '?' + new URLSearchParams(
            Object.entries(params)
              .filter(([, v]) => v !== undefined && v !== null && v !== '')
              .map(([k, v]) => [k, String(v)])
          ).toString()
        : ''

      const response = await $fetch<T>(`${baseURL}${url}${queryString}`, {
        method,
        headers: { ...getHeaders(), ...extraHeaders },
        body: body ? JSON.stringify(body) : undefined,
        onResponseError({ response }) {
          handleError({ response, data: response._data })
        },
      })

      return response
    } catch (error) {
      return handleError(error) as never
    }
  }

  return {
    get: <T>(url: string, params?: Record<string, unknown>) =>
      request<T>(url, { method: 'GET', params }),

    post: <T>(url: string, body?: unknown) =>
      request<T>(url, { method: 'POST', body }),

    put: <T>(url: string, body?: unknown) =>
      request<T>(url, { method: 'PUT', body }),

    patch: <T>(url: string, body?: unknown) =>
      request<T>(url, { method: 'PATCH', body }),

    delete: <T>(url: string) =>
      request<T>(url, { method: 'DELETE' }),

    upload: async <T>(url: string, formData: FormData, onProgress?: (progress: number) => void): Promise<T> => {
      const token = authStore.accessToken
      const headers: Record<string, string> = {}

      if (token) {
        headers['Authorization'] = `Bearer ${token}`
      }

      try {
        const response = await $fetch<T>(`${baseURL}${url}`, {
          method: 'POST',
          headers,
          body: formData,
          onRequest({ options }) {
            if (onProgress) {
              const total = (formData.get('file') as File)?.size || 0
              let loaded = 0
              const originalOnResponse = options.onResponse
              options.onResponse = (ctx) => {
                loaded += 4096
                onProgress(Math.min((loaded / total) * 100, 100))
                if (originalOnResponse) originalOnResponse(ctx)
              }
            }
          },
        })

        return response
      } catch (error) {
        return handleError(error) as never
      }
    },

    paginate: async <T>(url: string, params?: Record<string, unknown>): Promise<{ data: T[]; pagination: { page: number; limit: number; total: number; totalPages: number; hasNextPage: boolean; hasPreviousPage: boolean } }> => {
      const queryParams = { ...params, page: params?.page || 1, limit: params?.limit || 10 }
      return request(url, { method: 'GET', params: queryParams })
    },
  }
}

export const useApi = () => {
  return createApiClient()
}
