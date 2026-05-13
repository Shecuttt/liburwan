import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { ProtectedRoute } from "@/components/auth/ProtectedRoute"
import { RoleGuard } from "@/components/auth/RoleGuard"
import LoginPage from "@/pages/login/LoginPage"
import CallbackPage from "@/pages/auth/CallbackPage"
import CalendarPage from "@/pages/calendar/CalendarPage"
import MetricsPage from "@/pages/metrics/MetricsPage"
import SettingsPage from "@/pages/settings/SettingsPage"

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public Routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/auth/callback" element={<CallbackPage />} />

        {/* Protected Routes */}
        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<CalendarPage />} />
          <Route path="/calendar" element={<Navigate to="/" replace />} />
          
          {/* Admin Only Routes */}
          <Route element={<RoleGuard role="admin" />}>
            <Route path="/metrik" element={<MetricsPage />} />
            <Route path="/konfigurasi" element={<SettingsPage />} />
          </Route>
        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
