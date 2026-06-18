const API_BASE = 'http://localhost:8080'

let authToken: string | null = null

export function setAuthToken(token: string) {
  authToken = token
}

export function clearAuthToken() {
  authToken = null
}

function authHeaders(): Record<string, string> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`
  }
  return headers
}

// --- Types ---

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

export interface FriendRequestItem {
  sender_id: string
  sender_username: string
  sender_first_name: string
  sender_last_name: string
  created_at: string
}

export interface FriendRequest {
  recipient_id: string
  sender_id: string
  status: string
  created_at: string
}

export interface SendFriendRequestResponse {
  recipient_id: string
  status: string
  created_at: string
}

export interface RemoveFriendResponse {
  message: string
  id: string
}

interface ApiError {
  error: string
}

// --- Generic helpers ---

async function request<T>(
  method: string,
  path: string,
  body?: unknown
): Promise<T> {
  const init: RequestInit = {
    method,
    headers: authHeaders(),
  }
  if (body !== undefined) {
    init.body = JSON.stringify(body)
  }
  const res = await fetch(`${API_BASE}${path}`, init)
  const text = await res.text()
  let data: unknown = undefined
  if (text) {
    try { data = JSON.parse(text) } catch { data = text }
  }
  if (!res.ok) {
    const message = (data as ApiError)?.error ?? res.statusText
    throw new Error(message)
  }
  return data as T
}

// --- Auth ---

export function login(body: LoginRequest): Promise<AuthResponse> {
  return request<AuthResponse>('POST', '/api/login', body)
}

export function register(body: RegisterRequest): Promise<AuthResponse> {
  return request<AuthResponse>('POST', '/api/register', body)
}

// --- Users ---

export function searchUsers(query?: string): Promise<BackendUser[]> {
  const qs = query ? `?q=${encodeURIComponent(query)}` : ''
  return request<BackendUser[]>('GET', `/api/users${qs}`)
}

// --- Friends ---

export function getFriends(): Promise<BackendUser[]> {
  return request<BackendUser[]>('GET', '/api/friends')
}

export function removeFriend(friendId: string): Promise<RemoveFriendResponse> {
  return request<RemoveFriendResponse>('DELETE', `/api/friends/${friendId}`)
}

// --- Friend Requests ---

export function getFriendRequests(): Promise<FriendRequestItem[]> {
  return request<FriendRequestItem[]>('GET', '/api/friend-requests')
}

export function getFriendRequestByID(recipientId: string): Promise<FriendRequest> {
  return request<FriendRequest>('GET', `/api/friend-requests/${recipientId}`)
}

export function sendFriendRequest(recipientId: string): Promise<SendFriendRequestResponse> {
  return request<SendFriendRequestResponse>('POST', `/api/friend-requests/${recipientId}`)
}

export function acceptFriendRequest(senderId: string): Promise<void> {
  return request<void>('PUT', `/api/friend-requests/${senderId}`)
}

export function cancelFriendRequest(recipientId: string): Promise<void> {
  return request<void>('DELETE', `/api/friend-requests/outgoing/${recipientId}`)
}

export function rejectFriendRequest(senderId: string): Promise<void> {
  return request<void>('DELETE', `/api/friend-requests/incoming/${senderId}`)
}
