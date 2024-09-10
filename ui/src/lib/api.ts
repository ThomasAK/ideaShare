let csrfToken: string = ''

type Listener = () => Promise<void>
const unauthorizedListeners: Listener[] = []
const notFoundListeners: Listener[] = []
export const addUnauthorizedListener = (listener: Listener): void => { unauthorizedListeners.push(listener) }
export const addNotFoundListener = (listener: Listener): void => { notFoundListeners.push(listener) }

const getCsrfToken = (): string => {
  if (csrfToken) return csrfToken
  csrfToken = document.cookie.split(';').map(c => c.trim()).find(c => c.startsWith('csrf-token='))?.split('=')[1] ?? ''
  return csrfToken
}

export function addDefaultHeaders (headers?: Record<string, string>): Record<string, string> {
  return {
    accept: 'application/json',
    'content-type': 'application/json',
    'csrf-token': getCsrfToken(),
    ...(headers ?? {})
  }
}

export class ApiError extends Error {
  status: number
  statusText: string
  body: any
  constructor (err: string, status: number, statusText: string, body: any) {
    super(err)
    this.status = status
    this.statusText = statusText
    this.body = body
  }
}

async function makeApiError (res: Response): Promise<ApiError> {
  const txt = await (await res.blob()).text()
  let body: any = txt
  try {
    body = JSON.parse(txt)
  } catch (e) {
    console.warn(`failed to parse error body: ${String(e)}`)
  }
  console.dir(body)
  throw new ApiError(`Non-OK request(${res.status}) : ${res.statusText}`, res.status, res.statusText, body)
}

export async function apiRequest<T> (path: string, method: string, body: any, headers: Record<string, string>): Promise<T | null> {
  const res = await fetch(path, {
    method,
    body: body ? JSON.stringify(body) : body,
    headers: headers || {}
  })

  if (Math.floor(res.status / 100) !== 2) {
    if (res.status === 401) {
      for (const listener of unauthorizedListeners) {
        await listener()
      }
      throw await makeApiError(res)
    }

    if (res.status === 404) {
      for (const listener of notFoundListeners) {
        await listener()
      }
      return null
    }

    throw await makeApiError(res)
  } else {
    return await res.json()
  }
}

export async function apiGet<T> (path: string, headers?: Record<string, string>): Promise<T | null> {
  return await apiRequest(path, 'GET', null, addDefaultHeaders(headers))
}

export async function apiDelete<T> (path: string, headers?: Record<string, string>): Promise<T | null> {
  return await apiRequest(path, 'DELETE', null, addDefaultHeaders(headers))
}

export async function apiPost<T> (path: string, body: any, headers?: Record<string, string>): Promise<T | null> {
  return await apiRequest(path, 'POST', body, addDefaultHeaders(headers))
}

export async function apiPut<T> (path: string, body: any, headers?: Record<string, string>): Promise<T | null> {
  return await apiRequest(path, 'PUT', body, addDefaultHeaders(headers))
}
