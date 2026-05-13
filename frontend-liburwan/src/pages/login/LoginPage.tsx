import { Button } from "@/components/ui/button"
import { FcGoogle } from "react-icons/fc"

export default function LoginPage() {
  const handleLogin = () => {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
    window.location.href = `${apiBaseUrl}/auth/google`
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-linear-to-br from-background to-muted p-4">
      <div className="w-full max-w-md space-y-8 rounded-2xl border bg-card p-8 shadow-xl">
        <div className="text-center">
          <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
            <span className="text-3xl font-bold text-primary">L</span>
          </div>
          <h1 className="mt-6 text-3xl font-extrabold tracking-tight text-foreground">
            Liburwan
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Sistem Manajemen Libur & Penugasan Karyawan
          </p>
        </div>

        <div className="mt-8">
          <Button
            onClick={handleLogin}
            className="group relative flex w-full justify-center gap-3 rounded-xl bg-card px-4 py-6 text-sm font-semibold text-foreground shadow-sm ring-1 ring-inset ring-border hover:bg-muted/50"
          >
            <FcGoogle className="h-5 w-5" />
            <span>Login dengan Google</span>
          </Button>
        </div>

        <div className="mt-6 text-center text-xs text-muted-foreground">
          <p>Akses terbatas hanya untuk karyawan resmi.</p>
        </div>
      </div>
    </div>
  )
}
