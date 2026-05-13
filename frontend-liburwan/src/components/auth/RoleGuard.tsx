import { Navigate, Outlet } from "react-router-dom"
import { useAuth } from "@/hooks/useAuth"

interface RoleGuardProps {
  role: "admin" | "karyawan"
}

export function RoleGuard({ role }: RoleGuardProps) {
  const { karyawan } = useAuth()

  if (karyawan?.role !== role) {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}
