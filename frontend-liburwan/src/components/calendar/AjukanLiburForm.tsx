import { useState, useEffect } from "react"
import { useQueryClient } from "@tanstack/react-query"
import { formatInTimeZone } from "date-fns-tz"
import { id as localeID } from "date-fns/locale"
import { Loader2, AlertCircle, CheckCircle2 } from "lucide-react"

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
import type { CheckAvailabilityResponse } from "@/types/calendar"

interface AjukanLiburFormProps {
  date: Date
  onSuccess: () => void
  onCancel: () => void
}

export function AjukanLiburForm({ date, onSuccess, onCancel }: AjukanLiburFormProps) {
  const queryClient = useQueryClient()
  const [isChecking, setIsChecking] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [checkData, setCheckData] = useState<CheckAvailabilityResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [backupId, setBackupId] = useState<string>("")

  const dateStr = formatFullDate(date)
  const displayDate = formatInTimeZone(date, "Asia/Jakarta", "EEEE, d MMMM yyyy", { locale: localeID })

  useEffect(() => {
    const performCheck = async () => {
      try {
        setIsChecking(true)
        setError(null)
        const response = await apiFetch(`/jadwal-libur/check?tanggal=${dateStr}`)
        
        if (!response.ok) {
          const errData = await response.json()
          if (errData.code === "KUOTA_HABIS") {
            setError("KUOTA_HABIS")
          } else if (errData.code === "OUT_OF_WINDOW") {
            setError("OUT_OF_WINDOW")
          } else {
            setError("Gagal melakukan pengecekan ketersediaan.")
          }
          return
        }

        const data = await response.json()
        setCheckData(data)
      } catch {
        setError("Terjadi kesalahan jaringan.")
      } finally {
        setIsChecking(false)
      }
    }

    performCheck()
  }, [dateStr])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (checkData?.needs_backup && !backupId) {
      setError("BACKUP_REQUIRED")
      return
    }

    try {
      setIsSubmitting(true)
      setError(null)

      const response = await apiFetch("/jadwal-libur", {
        method: "POST",
        body: JSON.stringify({
          tanggal: dateStr,
          backup_karyawan_id: backupId || null,
        }),
      })

      if (!response.ok) {
        const errData = await response.json()
        setError(errData.code || "Gagal mengajukan libur.")
        return
      }

      // Success
      const bulanKey = formatInTimeZone(date, "Asia/Jakarta", "yyyy-MM")
      await queryClient.invalidateQueries({ queryKey: ["kalender", bulanKey] })
      await queryClient.invalidateQueries({ queryKey: ["jadwal-libur"] })
      onSuccess()
    } catch {
      setError("Terjadi kesalahan saat mengirim pengajuan.")
    } finally {
      setIsSubmitting(false)
    }
  }

  const getErrorMessage = (code: string) => {
    switch (code) {
      case "KUOTA_HABIS":
        return "Kuota libur Anda untuk bulan ini sudah habis."
      case "OUT_OF_WINDOW":
        return "Tanggal ini di luar jendela pengajuan (hanya bulan berjalan dan depan)."
      case "BACKUP_REQUIRED":
        return "Pilih karyawan backup terlebih dahulu."
      case "BACKUP_INVALID":
        return "Karyawan backup yang dipilih juga sedang libur di tanggal ini."
      case "NO_BACKUP_AVAILABLE":
        return "Tidak ada karyawan lain yang tersedia untuk menjadi backup di tanggal ini."
      default:
        return code
    }
  }

  if (isChecking) {
    return (
      <div className="flex flex-col items-center justify-center py-12 gap-4">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p className="text-sm text-muted-foreground italic">Mengecek kuota dan ketersediaan toko...</p>
      </div>
    )
  }

  if (error && (error === "KUOTA_HABIS" || error === "OUT_OF_WINDOW")) {
    return (
      <div className="space-y-6">
        <Alert variant="destructive" className="rounded-2xl border-destructive/20 bg-destructive/5">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle className="font-bold">Pengajuan Tidak Memungkinkan</AlertTitle>
          <AlertDescription>
            {getErrorMessage(error)}
          </AlertDescription>
        </Alert>
        <Button onClick={onCancel} variant="outline" className="w-full rounded-xl">Kembali</Button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="tanggal" className="text-xs font-bold uppercase text-muted-foreground">Tanggal Libur</Label>
          <div className="rounded-xl border bg-muted/30 p-3 font-medium text-foreground">
            {displayDate}
          </div>
        </div>

        {checkData?.needs_backup && (
          <div className="space-y-2">
            <Label htmlFor="backup" className="text-xs font-bold uppercase text-muted-foreground">Pilih Karyawan Backup</Label>
            <Select value={backupId} onValueChange={setBackupId} disabled={isSubmitting}>
              <SelectTrigger id="backup" className="h-12 rounded-xl ring-offset-background focus:ring-primary">
                <SelectValue placeholder="Pilih karyawan backup..." />
              </SelectTrigger>
              <SelectContent className="rounded-xl">
                {checkData.suggested_backup.map((k) => (
                  <SelectItem key={k.id} value={k.id}>
                    {k.nama} ({k.toko_nama})
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <p className="text-[10px] text-muted-foreground flex items-center gap-1">
              <AlertCircle className="h-3 w-3" />
              Toko rawan libur, diperlukan 1 backup dari toko lain.
            </p>
          </div>
        )}

        {!checkData?.needs_backup && checkData && (
          <Alert className="rounded-2xl border-success/20 bg-success/5 text-success-foreground ring-1 ring-success/10">
            <CheckCircle2 className="h-4 w-4 text-success" />
            <AlertTitle className="font-bold text-success">Kondisi Aman</AlertTitle>
            <AlertDescription className="text-xs">
              Kuota toko masih mencukupi. Anda bisa mengajukan libur tanpa backup.
            </AlertDescription>
          </Alert>
        )}

        {error && !["KUOTA_HABIS", "OUT_OF_WINDOW"].includes(error) && (
          <Alert variant="destructive" className="rounded-2xl border-destructive/20 bg-destructive/5">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription className="text-xs font-medium">
              {getErrorMessage(error)}
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
          disabled={isSubmitting || (checkData?.needs_backup && !backupId)}
          className="flex-1 rounded-xl h-12 gap-2"
        >
          {isSubmitting && <Loader2 className="h-4 w-4 animate-spin" />}
          Ajukan
        </Button>
      </div>
    </form>
  )
}
