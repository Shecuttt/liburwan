import { create } from "zustand"
import { persist, createJSONStorage } from "zustand/middleware"

export interface Karyawan {
  id: string
  nama: string
  role: string
  toko_id: string
  toko_nama: string
}

interface AuthState {
  token: string | null
  karyawan: Karyawan | null
  setSession: (token: string, karyawan: Karyawan) => void
  clearSession: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      karyawan: null,
      setSession: (token, karyawan) => set({ token, karyawan }),
      clearSession: () => set({ token: null, karyawan: null }),
    }),
    {
      name: "auth-storage",
      storage: createJSONStorage(() => sessionStorage),
    }
  )
)
