import { useState, useEffect } from "react"
import { formatInTimeZone } from "date-fns-tz"
import { Clock as ClockIcon } from "lucide-react"

export function DigitalClock() {
  const [time, setTime] = useState(new Date())

  useEffect(() => {
    const timer = setInterval(() => {
      setTime(new Date())
    }, 1000)

    return () => clearInterval(timer)
  }, [])

  return (
    <div className="flex items-center gap-2 rounded-xl bg-secondary/50 px-3 py-1.5 font-mono text-sm font-bold text-secondary-foreground shadow-inner ring-1 ring-border">
      <ClockIcon className="h-3.5 w-3.5 text-primary" />
      <span>{formatInTimeZone(time, "Asia/Jakarta", "HH:mm:ss")}</span>
      <span className="text-[10px] text-muted-foreground uppercase">WIB</span>
    </div>
  )
}
