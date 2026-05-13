import { useAuthStore } from "@/store/useAuthStore"

export function useAuth() {
  const { karyawan, token } = useAuthStore()

  const isAdmin = () => {
    return karyawan?.role === "admin"
  }

  return {
    karyawan,
    token,
    isAdmin,
    isAuthenticated: !!token,
  }
}
