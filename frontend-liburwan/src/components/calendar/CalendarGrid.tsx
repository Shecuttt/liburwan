import { isSameMonth, isSameDay } from "date-fns"
import { CalendarCell } from "./CalendarCell"
import type { CalendarDayData } from "@/types/calendar"
import { formatFullDate } from "@/lib/date-utils"

interface CalendarGridProps {
  days: Date[]
  currentMonth: Date
  data: CalendarDayData[]
  onDateClick: (date: Date, data?: CalendarDayData) => void
}

const WEEKDAYS = ["Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"]

export function CalendarGrid({
  days,
  currentMonth,
  data,
  onDateClick,
}: CalendarGridProps) {
  // Create a map for quick lookup of data by date string
  const dataMap = new Map(data.map((d) => [d.tanggal, d]))

  return (
    <div className="rounded-xl border border-border bg-card shadow-sm overflow-hidden">
      <div className="grid grid-cols-7 border-b bg-muted/30">
        {WEEKDAYS.map((day) => (
          <div
            key={day}
            className="py-3 text-center text-xs font-bold tracking-wider text-muted-foreground uppercase"
          >
            {day}
          </div>
        ))}
      </div>

      <div className="grid grid-cols-7">
        {days.map((date) => {
          const dateStr = formatFullDate(date)
          const dayData = dataMap.get(dateStr)
          const isToday = isSameDay(date, new Date())
          const isCurrentMonth = isSameMonth(date, currentMonth)

          return (
            <CalendarCell
              key={dateStr}
              date={date}
              isCurrentMonth={isCurrentMonth}
              isToday={isToday}
              data={dayData}
              onClick={() => onDateClick(date, dayData)}
            />
          )
        })}
      </div>
    </div>
  )
}
