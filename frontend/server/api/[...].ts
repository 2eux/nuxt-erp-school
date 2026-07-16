export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string

  const path = event.path.replace(/^\/api/, '')
  const method = event.method
  const headers = getHeaders(event)
  const body = method !== 'GET' && method !== 'HEAD' ? await readBody(event).catch(() => null) : null
  const query = getQuery(event)

  const forwardHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  }

  if (headers['authorization']) {
    forwardHeaders['Authorization'] = headers['authorization']
  }

  if (headers['x-school-id']) {
    forwardHeaders['X-School-Id'] = headers['x-school-id']
  }

  const url = new URL(`${apiBase}${path}`)

  for (const [key, value] of Object.entries(query)) {
    if (value !== undefined && value !== null) {
      url.searchParams.append(key, String(value))
    }
  }

  try {
    const response = await $fetch(url.toString(), {
      method: method as 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE',
      headers: forwardHeaders,
      body: body ? JSON.stringify(body) : undefined,
      timeout: 30000,
      ignoreResponseError: true,
    })

    return response
  } catch (error: unknown) {
    const err = error as { statusCode?: number; message?: string; data?: unknown }

    setResponseStatus(event, err.statusCode || 502)

    return {
      success: false,
      error: err.message || 'Proxy Error',
      data: err.data || null,
    }
  }
})
