import {
  startOfMonth,
  endOfMonth,
  startOfWeek,
  endOfWeek,
  eachDayOfInterval,
  addMonths,
  isBefore,
  startOfDay,
} from "date-fns"
import { formatInTimeZone, toDate } from "date-fns-tz"
import { id as localeID } from "date-fns/locale"

const TIMEZONE = "Asia/Jakarta"

export const getCalendarDays = (date: Date) => {
  const start = startOfWeek(startOfMonth(date))
  const end = endOfWeek(endOfMonth(date))

  return eachDayOfInterval({ start, end })
}

export const formatYearMonth = (date: Date) => formatInTimeZone(date, TIMEZONE, "yyyy-MM")
export const formatFullDate = (date: Date) => formatInTimeZone(date, TIMEZONE, "yyyy-MM-dd")

export const isBookingWindow = (date: Date) => {
  // Use toDate to ensure we are working with the correct point in time in Jakarta
  const nowInJakarta = toDate(new Date(), { timeZone: TIMEZONE })
  const today = startOfDay(nowInJakarta)

  const startOfCurrentMonth = startOfMonth(today)
  const endOfNextMonth = endOfMonth(addMonths(today, 1))

  // Normalize target date to Jakarta
  const targetDate = toDate(date, { timeZone: TIMEZONE })
  const targetDay = startOfDay(targetDate)

  // Past dates are not booking window
  if (isBefore(targetDay, today)) return false

  // Must be in current or next month
  return (
    (targetDay >= startOfCurrentMonth && targetDay <= endOfNextMonth)
  )
}

export const getMonthYearDisplay = (date: Date) =>
  formatInTimeZone(date, TIMEZONE, "MMMM yyyy", { locale: localeID })
