import { useState, useMemo } from "react"
import { useQuery } from "@tanstack/react-query"
import {
  addMonths,
  subMonths,
  startOfMonth,
  isBefore,
  isAfter,
} from "date-fns"
import { ChevronLeft, ChevronRight, Loader2, Calendar as CalendarIcon, RefreshCw, BarChart3, Settings } from "lucide-react"
import { Link } from "react-router-dom"

import { useAuth } from "@/hooks/useAuth"
import { apiFetch } from "@/lib/api"
import { getCalendarDays, formatYearMonth, getMonthYearDisplay } from "@/lib/date-utils"
import { CalendarGrid } from "@/components/calendar/CalendarGrid"
import { DayDetailSheet } from "@/components/calendar/DayDetailSheet"
import { DigitalClock } from "@/components/calendar/DigitalClock"
import type { CalendarResponse, CalendarDayData } from "@/types/calendar"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"

export default function CalendarPage() {
  const { karyawan, isAdmin } = useAuth()

  // State for current displayed month
  const [currentMonth, setCurrentMonth] = useState(startOfMonth(new Date()))
  const [selectedDate, setSelectedDate] = useState<Date | null>(null)
  const [selectedData, setSelectedData] = useState<CalendarDayData | undefined>()
  const [isSheetOpen, setIsSheetOpen] = useState(false)

  const bulanStr = useMemo(() => formatYearMonth(currentMonth), [currentMonth])
  const days = useMemo(() => getCalendarDays(currentMonth), [currentMonth])

  // TanStack Query to fetch calendar data
  const { data, isLoading, error, refetch } = useQuery<CalendarResponse>({
    queryKey: ["kalender", bulanStr],
    queryFn: async () => {
      const response = await apiFetch(`/kalender?bulan=${bulanStr}`)
      if (!response.ok) throw new Error("Gagal mengambil data kalender")
      return response.json()
    },
  })

  // Navigation handlers
  const today = startOfMonth(new Date())
  const nextMonth = addMonths(today, 1)

  const handlePrevMonth = () => {
    if (!isBefore(currentMonth, today)) {
      setCurrentMonth((prev) => subMonths(prev, 1))
    }
  }

  const handleNextMonth = () => {
    if (!isAfter(currentMonth, nextMonth)) {
      setCurrentMonth((prev) => addMonths(prev, 1))
    }
  }

  const handleDateClick = (date: Date, dayData?: CalendarDayData) => {
    setSelectedDate(date)
    setSelectedData(dayData)
    setIsSheetOpen(true)
  }

  const isPrevDisabled = !isAfter(currentMonth, today)
  const isNextDisabled = !isBefore(currentMonth, nextMonth)

  return (
    <div className="mx-auto max-w-7xl p-4 md:p-8">
      {/* Header with Navigation */}
      <header className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <div className="rounded-2xl bg-primary/10 p-3 text-primary">
            <CalendarIcon className="h-6 w-6" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-bold tracking-tight text-foreground md:text-3xl">
                Shared Calendar
              </h1>
              {isAdmin() && (
                <div className="flex items-center gap-1">
                  <Link 
                    to="/metrik" 
                    title="Lihat Metrik & Statistik"
                    className="rounded-lg p-1.5 hover:bg-muted transition-colors text-muted-foreground hover:text-primary"
                  >
                    <BarChart3 className="h-5 w-5" />
                  </Link>
                  <Link 
                    to="/konfigurasi" 
                    title="Pengaturan Sistem"
                    className="rounded-lg p-1.5 hover:bg-muted transition-colors text-muted-foreground hover:text-primary"
                  >
                    <Settings className="h-5 w-5" />
                  </Link>
                </div>
              )}
            </div>
            <p className="text-sm text-muted-foreground">
              Selamat datang, <span className="font-semibold text-foreground">{karyawan?.nama}</span>
            </p>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-4">
          <DigitalClock />
          
          <div className="flex items-center justify-between gap-4 rounded-2xl bg-card p-2 shadow-sm ring-1 ring-border sm:justify-end">
            <Button
              variant="ghost"
              size="icon"
              onClick={handlePrevMonth}
              disabled={isPrevDisabled}
              className="rounded-xl"
            >
              <ChevronLeft className="h-5 w-5" />
            </Button>

            <div className="min-w-[140px] text-center font-bold text-foreground">
              {getMonthYearDisplay(currentMonth)}
            </div>

            <Button
              variant="ghost"
              size="icon"
              onClick={handleNextMonth}
              disabled={isNextDisabled}
              className="rounded-xl"
            >
              <ChevronRight className="h-5 w-5" />
            </Button>
          </div>
        </div>
      </header>

      {/* Main Grid Section */}
      <div className="relative">
        {error ? (
          <div className="flex min-h-[400px] flex-col items-center justify-center gap-4 rounded-3xl border border-dashed border-destructive/50 bg-destructive/5 text-center p-8">
            <div className="rounded-full bg-destructive/10 p-3 text-destructive">
              <RefreshCw className="h-6 w-6" />
            </div>
            <div>
              <h3 className="text-lg font-bold">Terjadi Kesalahan</h3>
              <p className="max-w-xs text-sm text-muted-foreground">
                Gagal memuat data kalender. Silakan coba lagi nanti.
              </p>
            </div>
            <Button variant="outline" onClick={() => refetch()} className="rounded-xl">
              Coba Lagi
            </Button>
          </div>
        ) : isLoading ? (
          <div className="space-y-4">
            <div className="grid grid-cols-7 gap-px overflow-hidden rounded-xl border bg-muted/50">
              {Array.from({ length: 35 }).map((_, i) => (
                <Skeleton key={i} className="h-24 w-full rounded-none" />
              ))}
            </div>
            <div className="flex items-center justify-center gap-2 text-sm text-muted-foreground animate-pulse">
              <Loader2 className="h-4 w-4 animate-spin" />
              Sinkronisasi data kalender...
            </div>
          </div>
        ) : (
          <CalendarGrid
            days={days}
            currentMonth={currentMonth}
            data={data?.hari || []}
            onDateClick={handleDateClick}
          />
        )}
      </div>

      <footer className="mt-8 flex flex-wrap items-center gap-6 text-xs text-muted-foreground">
        <div className="flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-primary" />
          Planned Leave
        </div>
        <div className="flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-warning" />
          Unplanned Leave
        </div>
        <div className="flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-destructive animate-pulse" />
          Toko Rawan (Tersisa 2 orang)
        </div>
      </footer>

      {/* Day Detail Interaction */}
      <DayDetailSheet
        isOpen={isSheetOpen}
        onClose={() => setIsSheetOpen(false)}
        date={selectedDate}
        data={selectedData}
      />
    </div>
  )
}
