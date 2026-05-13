import { format } from "date-fns"
import { cn } from "@/lib/utils"
import type { CalendarDayData } from "@/types/calendar"

interface CalendarCellProps {
  date: Date
  isCurrentMonth: boolean
  isToday: boolean
  data?: CalendarDayData
  onClick: () => void
}

export function CalendarCell({
  date,
  isCurrentMonth,
  isToday,
  data,
  onClick,
}: CalendarCellProps) {
  const dayNumber = format(date, "d")

  // Filter only stores that have someone on leave
  const storesWithLeave = data?.toko.filter((t) => t.libur.length > 0) || []

  return (
    <div
      onClick={onClick}
      className={cn(
        "min-h-[100px] cursor-pointer border-r border-b p-2 transition-colors hover:bg-muted/50",
        !isCurrentMonth && "bg-muted/30",
        isToday && "bg-primary/10"
      )}
    >
      <div className="mb-2 flex items-center justify-between">
        <span
          className={cn(
            "flex h-7 w-7 items-center justify-center rounded-full text-sm font-medium",
            !isCurrentMonth && "text-muted-foreground",
            isToday && "bg-primary text-primary-foreground"
          )}
        >
          {dayNumber}
        </span>
      </div>

      <div className="space-y-2">
        {storesWithLeave.map((t) => (
          <div key={t.toko_id} className="text-[10px] leading-tight">
            <div className="flex items-center gap-1 font-bold text-foreground">
              {t.toko_nama}
              {t.is_rawan && (
                <span className="h-1.5 w-1.5 rounded-full bg-destructive animate-pulse" title="Toko Rawan" />
              )}
            </div>
            <ul className="space-y-0.5 text-muted-foreground">
              {t.libur.map((l) => (
                <li key={l.karyawan_id} className="flex items-center gap-1">
                  <span className={cn(
                    "h-1 w-1 rounded-full",
                    l.tipe === "planned" ? "bg-primary" : "bg-warning" // Warning variable might not exist, will check
                  )} />
                  <span className="truncate">
                    {l.karyawan_nama}
                    <span className="ml-0.5 text-[8px] opacity-70">
                      ({l.tipe})
                    </span>
                  </span>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
    </div>
  )
}
