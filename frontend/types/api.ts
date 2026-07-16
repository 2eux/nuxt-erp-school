import type { ApiResponse, PaginatedResponse } from './index'

export interface RequestConfig {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  params?: Record<string, unknown>
  body?: unknown
  headers?: Record<string, string>
  ssr?: boolean
  cache?: boolean
}

export interface ApiError {
  statusCode: number
  message: string
  error: string
  details?: Record<string, string[]>
}

export interface ApiClient {
  get<T>(url: string, params?: Record<string, unknown>): Promise<T>
  post<T>(url: string, body?: unknown): Promise<T>
  put<T>(url: string, body?: unknown): Promise<T>
  patch<T>(url: string, body?: unknown): Promise<T>
  delete<T>(url: string): Promise<T>
  upload<T>(url: string, formData: FormData, onProgress?: (progress: number) => void): Promise<T>
  paginate<T>(url: string, params?: Record<string, unknown>): Promise<PaginatedResponse<T>>
}

export type ApiEndpoint = string

export interface QueryOptions {
  enabled?: boolean
  staleTime?: number
  cacheTime?: number
  retry?: number
  onSuccess?: (data: unknown) => void
  onError?: (error: ApiError) => void
}

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
