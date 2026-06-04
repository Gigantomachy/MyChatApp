const API_BASE = 'http://localhost:8080'

export interface BackendUser {
  user_id: string
  username: string
  email: string
  first_name: string
  last_name: string
  created_at: string
}

export interface AuthResponse {
  token: string
  user: BackendUser
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  email: string
  first_name: string
  last_name: string
}

export interface ApiError {
  error: string
}

async function apiPost<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  const data = await res.json()

  if (!res.ok) {
    const message = (data as ApiError).error ?? res.statusText
    throw new Error(message)
  }

  return data as T
}

export function login(body: LoginRequest): Promise<AuthResponse> {
  return apiPost<AuthResponse>('/login', body)
}

export function register(body: RegisterRequest): Promise<AuthResponse> {
  return apiPost<AuthResponse>('/register', body)
}