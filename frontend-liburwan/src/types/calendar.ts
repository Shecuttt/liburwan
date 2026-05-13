export interface KaryawanLibur {
  id: string
  karyawan_id: string
  karyawan_nama: string
  tipe: "planned" | "unplanned"
  backup_assignment_id?: string | null
  backup_karyawan_nama?: string | null
}

export interface TokoCalendarInfo {
  toko_id: string
  toko_nama: string
  available_count: number
  is_rawan: boolean
  libur: KaryawanLibur[]
}

export interface CalendarDayData {
  tanggal: string
  toko: TokoCalendarInfo[]
}

export interface CalendarResponse {
  bulan: string
  hari: CalendarDayData[]
}

export interface KaryawanProfile {
  id: string
  nama: string
  toko_nama: string
}

export interface CheckAvailabilityResponse {
  available_after: number
  needs_backup: boolean
  suggested_backup: KaryawanProfile[]
}
