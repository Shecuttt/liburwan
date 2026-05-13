import { useState, useEffect } from "react"
import { useQueryClient } from "@tanstack/react-query"
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
import { Alert, AlertDescription } from "@/components/ui/alert"
import { apiFetch } from "@/lib/api"
import type { KaryawanProfile } from "@/types/calendar"

interface AssignBackupFormProps {
  jadwalLiburId: string
  karyawanNama: string
  tokoNama: string
  onLeaveKaryawanIds: string[]
  onSuccess: () => void
  onCancel: () => void
}

export function AssignBackupForm({
  jadwalLiburId,
  karyawanNama,
  tokoNama,
  onLeaveKaryawanIds,
  onSuccess,
  onCancel,
}: AssignBackupFormProps) {
  const queryClient = useQueryClient()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isLoadingKaryawan, setIsLoadingKaryawan] = useState(true)
  const [karyawans, setKaryawans] = useState<KaryawanProfile[]>([])
  const [selectedKaryawanId, setSelectedKaryawanId] = useState<string>("")
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchKaryawan = async () => {
      try {
        setIsLoadingKaryawan(true)
        const response = await apiFetch("/karyawan")
        if (!response.ok) throw new Error("Gagal mengambil data karyawan")
        const data: KaryawanProfile[] = await response.json()
        
        // Exclude those on leave
        const available = data.filter(k => !onLeaveKaryawanIds.includes(k.id))
        setKaryawans(available)
      } catch (err) {
        setError(err instanceof Error ? err.message : String(err))
      } finally {
        setIsLoadingKaryawan(false)
      }
    }
    fetchKaryawan()
  }, [onLeaveKaryawanIds])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedKaryawanId) return

    try {
      setIsSubmitting(true)
      setError(null)

      const response = await apiFetch("/backup-assignment", {
        method: "POST",
        body: JSON.stringify({
          jadwal_libur_id: jadwalLiburId,
          backup_karyawan_id: selectedKaryawanId,
        }),
      })

      if (!response.ok) {
        const errData = await response.json()
        if (errData.code === "BACKUP_INVALID") {
          setError("Karyawan yang dipilih juga libur di tanggal ini")
        } else if (errData.code === "ALREADY_ASSIGNED") {
          setError("Jadwal libur ini sudah memiliki backup assignment")
        } else {
          setError(errData.message || "Gagal mengkonfirmasi backup")
        }
        return
      }

      await queryClient.invalidateQueries({ queryKey: ["kalender"] })
      await queryClient.invalidateQueries({ queryKey: ["jadwal-libur"] })
      onSuccess()
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-4">
        <div className="space-y-2">
          <Label className="text-xs font-bold uppercase text-muted-foreground">Jadwal Libur</Label>
          <div className="rounded-xl border bg-muted/30 p-3 space-y-1">
            <div className="flex items-center gap-2 font-medium text-foreground">
              <User className="h-3 w-3 text-primary" />
              {karyawanNama}
            </div>
            <div className="flex items-center gap-2 text-xs text-muted-foreground">
              <Store className="h-3 w-3" />
              {tokoNama}
            </div>
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="backup-karyawan" className="text-xs font-bold uppercase text-muted-foreground">Pilih Karyawan Backup</Label>
          <Select 
            value={selectedKaryawanId} 
            onValueChange={setSelectedKaryawanId} 
            disabled={isSubmitting || isLoadingKaryawan}
          >
            <SelectTrigger id="backup-karyawan" className="h-12 rounded-xl">
              <SelectValue placeholder={isLoadingKaryawan ? "Memuat data..." : "Pilih backup..."} />
            </SelectTrigger>
            <SelectContent className="rounded-xl max-h-[300px]">
              {karyawans.map((k) => (
                <SelectItem key={k.id} value={k.id}>
                  {k.nama} ({k.toko_nama})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <p className="text-[10px] text-muted-foreground px-1">
            *Karyawan yang libur di tanggal ini otomatis disembunyikan.
          </p>
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
          Konfirmasi
        </Button>
      </div>
    </form>
  )
}
