import { useState, useEffect } from "react"
import { useQueryClient } from "@tanstack/react-query"
import { formatInTimeZone } from "date-fns-tz"
import { id as localeID } from "date-fns/locale"
import { Loader2, AlertCircle, User, Store } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { apiFetch } from "@/lib/api"
import { formatFullDate } from "@/lib/date-utils"
import type { KaryawanProfile } from "@/types/calendar"

interface UnplannedLeaveFormProps {
  date: Date
  onSuccess: () => void
  onCancel: () => void
}

interface UnplannedResponse {
  jadwal_libur: { id?: string; tanggal?: string; tipe?: string }
  availability_after: number
  suggested_backup: KaryawanProfile[]
}

export function UnplannedLeaveForm({ date, onSuccess, onCancel }: UnplannedLeaveFormProps) {
  const queryClient = useQueryClient()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isLoadingKaryawan, setIsLoadingKaryawan] = useState(true)
  const [karyawans, setKaryawans] = useState<KaryawanProfile[]>([])
  const [selectedKaryawanId, setSelectedKaryawanId] = useState<string>("")
  const [error, setError] = useState<string | null>(null)
  const [result, setResult] = useState<UnplannedResponse | null>(null)

  const dateStr = formatFullDate(date)
  const displayDate = formatInTimeZone(date, "Asia/Jakarta", "EEEE, d MMMM yyyy", { locale: localeID })

  useEffect(() => {
    const fetchKaryawan = async () => {
      try {
        setIsLoadingKaryawan(true)
        const response = await apiFetch("/karyawan")
        if (!response.ok) throw new Error("Gagal mengambil data karyawan")
        const data = await response.json()
        setKaryawans(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : String(err))
      } finally {
        setIsLoadingKaryawan(false)
      }
    }
    fetchKaryawan()
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedKaryawanId) return

    try {
      setIsSubmitting(true)
      setError(null)

      const response = await apiFetch("/jadwal-libur/unplanned", {
        method: "POST",
        body: JSON.stringify({
          karyawan_id: selectedKaryawanId,
          tanggal: dateStr,
        }),
      })

      if (!response.ok) {
        const errData = await response.json()
        setError(errData.message || "Gagal mencatat unplanned leave")
        return
      }

      const data: UnplannedResponse = await response.json()

      // If availability is safe (>= 2), we can just finish
      if (data.availability_after >= 2) {
        await queryClient.invalidateQueries({ queryKey: ["kalender"] })
        await queryClient.invalidateQueries({ queryKey: ["jadwal-libur"] })
        onSuccess()
      } else {
        // Show warning/result view
        setResult(data)
        await queryClient.invalidateQueries({ queryKey: ["kalender"] })
        await queryClient.invalidateQueries({ queryKey: ["jadwal-libur"] })
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setIsSubmitting(false)
    }
  }

  if (result) {
    return (
      <div className="space-y-6">
        <Alert className="rounded-2xl border-warning/20 bg-warning/5 text-warning-foreground ring-1 ring-warning/10">
          <AlertCircle className="h-4 w-4 text-warning" />
          <AlertTitle className="font-bold">Perhatian: Toko Kekurangan Karyawan</AlertTitle>
          <AlertDescription className="text-xs">
            Unplanned leave berhasil dicatat, namun toko ini sekarang hanya memiliki {result.availability_after} karyawan yang bertugas.
          </AlertDescription>
        </Alert>

        <div className="space-y-4">
          <h4 className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">Karyawan Tersedia (Backup):</h4>
          {result.suggested_backup.length > 0 ? (
            <div className="space-y-2">
              {result.suggested_backup.map((b) => (
                <div key={b.id} className="flex items-center justify-between rounded-xl border bg-muted/30 p-3 text-xs">
                  <div className="flex items-center gap-2">
                    <User className="h-3 w-3 text-muted-foreground" />
                    <span className="font-medium">{b.nama}</span>
                  </div>
                  <div className="flex items-center gap-1 text-muted-foreground">
                    <Store className="h-3 w-3" />
                    {b.toko_nama}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-xs text-muted-foreground italic">Tidak ada karyawan lain yang tersedia sebagai backup hari ini.</p>
          )}
        </div>

        <Button onClick={onSuccess} className="w-full rounded-xl h-12 font-bold">
          Tutup
        </Button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-4">
        <div className="space-y-2">
          <Label className="text-xs font-bold uppercase text-muted-foreground">Tanggal Unplanned Leave</Label>
          <div className="rounded-xl border bg-muted/30 p-3 font-medium text-foreground">
            {displayDate}
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="karyawan" className="text-xs font-bold uppercase text-muted-foreground">Pilih Karyawan</Label>
          <Select
            value={selectedKaryawanId}
            onValueChange={setSelectedKaryawanId}
            disabled={isSubmitting || isLoadingKaryawan}
          >
            <SelectTrigger id="karyawan" className="h-12 rounded-xl">
              <SelectValue placeholder={isLoadingKaryawan ? "Memuat data..." : "Pilih karyawan..."} />
            </SelectTrigger>
            <SelectContent className="rounded-xl max-h-[300px]">
              {karyawans.map((k) => (
                <SelectItem key={k.id} value={k.id}>
                  {k.nama} ({k.toko_nama})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {error && (
          <Alert variant="destructive" className="rounded-2xl border-destructive/20 bg-destructive/5">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription className="text-xs font-medium">
              {error}
            </AlertDescription>
          </Alert>
        )}
      </div>

      <div className="flex gap-3 pt-4">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isSubmitting}
          className="flex-1 rounded-xl h-12"
        >
          Batal
        </Button>
        <Button
          type="submit"
          disabled={isSubmitting || !selectedKaryawanId}
          className="flex-1 rounded-xl h-12 gap-2"
        >
          {isSubmitting && <Loader2 className="h-4 w-4 animate-spin" />}
          Simpan
        </Button>
      </div>
    </form>
  )
}
