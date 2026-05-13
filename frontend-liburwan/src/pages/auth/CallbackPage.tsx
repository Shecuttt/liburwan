import { useEffect } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { useAuthStore } from "@/store/useAuthStore"
import { apiFetch } from "@/lib/api"

export default function CallbackPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const setSession = useAuthStore((state) => state.setSession)

  useEffect(() => {
    const handleCallback = async () => {
      const token = searchParams.get("token")

      if (!token) {
        console.error("No token found in callback")
        navigate("/login", { replace: true })
        return
      }

      try {
        // Fetch user profile
        const response = await apiFetch("/auth/me", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        })

        if (!response.ok) {
          throw new Error("Failed to fetch profile")
        }

        const karyawanData = await response.json()
        
        // Map backend response to the Karyawan interface
        const mappedKaryawan = {
          id: karyawanData.id,
          nama: karyawanData.nama,
          role: karyawanData.role,
          toko_id: karyawanData.toko_id,
          toko_nama: karyawanData.toko?.Nama || "Unknown Store"
        }
        
        setSession(token, mappedKaryawan)

        // Redirect to home (calendar)
        navigate("/", { replace: true })
      } catch (error) {
        console.error("Error during auth callback:", error)
        navigate("/login", { replace: true })
      }
    }

    handleCallback()
  }, [searchParams, navigate, setSession])

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="text-center">
        <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-primary border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]" />
        <p className="mt-4 text-muted-foreground font-medium">
          Mempersiapkan sesi Anda...
        </p>
      </div>
    </div>
  )
}
