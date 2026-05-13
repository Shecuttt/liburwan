import { useAuthStore } from "@/store/useAuthStore"

const BASE_URL = import.meta.env.VITE_API_BASE_URL

export async function apiFetch(endpoint: string, options: RequestInit = {}) {
  const { token, clearSession } = useAuthStore.getState()

  const headers = new Headers(options.headers)
  if (token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers,
  })

  if (response.status === 401) {
    clearSession()
    window.location.href = "/login"
    throw new Error("Unauthorized")
  }

  return response
}
