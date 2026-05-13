import { useState } from "react"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import {
  Settings,
  ArrowLeft,
  Save,
  X,
  Edit2,
  Loader2,
  ShieldCheck,
  CalendarDays,
  Users
} from "lucide-react"
import { Link } from "react-router-dom"

import { apiFetch } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Skeleton } from "@/components/ui/skeleton"

interface Konfigurasi {
  id: string
  key: string
  value: string
  updated_at: string
}

const CONFIG_LABELS: Record<string, { label: string; icon: React.ElementType; description: string; max: number }> = {
  maks_libur_per_bulan: {
    label: "Maksimal Hari Libur Per Bulan",
    icon: CalendarDays,
    description: "Batas maksimal pengajuan libur terencana (planned) per karyawan.",
    max: 31
  },
  min_available_per_hari: {
    label: "Minimal Karyawan Available Per Hari",
    icon: Users,
    description: "Jumlah minimal karyawan yang harus masuk agar toko tidak dianggap kritis.",
    max: 10
  }
}

export default function SettingsPage() {

  const { data: configs, isLoading } = useQuery<Konfigurasi[]>({
    queryKey: ["konfigurasi"],
    queryFn: async () => {
      const res = await apiFetch("/konfigurasi")
      if (!res.ok) throw new Error("Gagal memuat konfigurasi")
      return res.json()
    }
  })

  return (
    <div className="mx-auto max-w-4xl p-4 md:p-8 space-y-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <Link to="/" className="rounded-xl border p-2 hover:bg-muted transition-colors">
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <div className="rounded-2xl bg-primary/10 p-3 text-primary">
            <Settings className="h-6 w-6" />
          </div>
          <div>
            <h1 className="text-2xl font-bold tracking-tight text-foreground md:text-3xl">
              Pengaturan Sistem
            </h1>
            <p className="text-sm text-muted-foreground">
              Konfigurasi kebijakan libur dan ketersediaan toko.
            </p>
          </div>
        </div>
      </header>

      <div className="grid gap-6">
        {isLoading ? (
          <>
            <Skeleton className="h-32 rounded-3xl" />
            <Skeleton className="h-32 rounded-3xl" />
          </>
        ) : configs?.map(config => (
          <ConfigCard key={config.key} config={config} />
        ))}

        {!isLoading && configs?.length === 0 && (
          <div className="flex flex-col items-center justify-center h-48 border border-dashed rounded-3xl text-muted-foreground">
            <ShieldCheck className="h-10 w-10 mb-2 opacity-20" />
            <p>Tidak ada konfigurasi ditemukan.</p>
          </div>
        )}
      </div>
    </div>
  )
}

function ConfigCard({ config }: { config: Konfigurasi }) {
  const queryClient = useQueryClient()
  const [isEditing, setIsEditing] = useState(false)
  const [value, setValue] = useState(config.value)
  const [error, setError] = useState<string | null>(null)

  const info = CONFIG_LABELS[config.key] || {
    label: config.key,
    icon: Settings,
    description: "Pengaturan sistem.",
    max: 999
  }

  const mutation = useMutation({
    mutationFn: async (newValue: string) => {
      const res = await apiFetch(`/konfigurasi/${config.key}`, {
        method: "PATCH",
        body: JSON.stringify({ value: newValue })
      })
      if (res.status === 404) throw new Error("Konfigurasi tidak ditemukan")
      if (!res.ok) throw new Error("Gagal memperbarui konfigurasi")
      return res.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["konfigurasi"] })
      setIsEditing(false)
      setError(null)
    },
    onError: (err) => {
      setError(err instanceof Error ? err.message : String(err))
    }
  })

  const handleSave = () => {
    const num = parseInt(value)
    if (isNaN(num) || num <= 0) {
      setError("Value harus angka positif")
      return
    }
    if (num > info.max) {
      setError(`Maksimal value adalah ${info.max}`)
      return
    }

    mutation.mutate(value)
  }

  const handleCancel = () => {
    setIsEditing(false)
    setValue(config.value)
    setError(null)
  }

  return (
    <Card className="rounded-3xl border-none bg-muted/30 overflow-hidden">
      <CardContent className="p-6">
        <div className="flex items-start justify-between gap-4">
          <div className="flex items-center gap-4">
            <div className="rounded-2xl bg-background p-3 text-primary shadow-sm">
              <info.icon className="h-6 w-6" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-foreground">{info.label}</h3>
              <p className="text-sm text-muted-foreground">{info.description}</p>
            </div>
          </div>

          {!isEditing && (
            <div className="flex flex-col items-end gap-2">
              <div className="text-3xl font-black text-primary">{config.value}</div>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setIsEditing(true)}
                className="rounded-xl h-9 px-4 gap-2 border-primary/20 hover:bg-primary/5 hover:text-primary"
              >
                <Edit2 className="h-3.5 w-3.5" />
                Edit
              </Button>
            </div>
          )}
        </div>

        {isEditing && (
          <div className="mt-6 flex flex-col gap-4 p-4 rounded-2xl bg-background animate-in slide-in-from-top-2 duration-300">
            <div className="space-y-2">
              <label className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground ml-1">
                Value Baru
              </label>
              <div className="flex items-center gap-3">
                <Input
                  type="number"
                  value={value}
                  onChange={(e) => setValue(e.target.value)}
                  disabled={mutation.isPending}
                  className="rounded-xl h-12 flex-1 font-bold text-lg"
                  placeholder="Masukkan angka..."
                />
                <div className="flex items-center gap-2">
                  <Button
                    onClick={handleSave}
                    disabled={mutation.isPending}
                    className="rounded-xl h-12 px-6 gap-2"
                  >
                    {mutation.isPending ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Save className="h-4 w-4" />
                    )}
                    Simpan
                  </Button>
                  <Button
                    variant="ghost"
                    onClick={handleCancel}
                    disabled={mutation.isPending}
                    className="rounded-xl h-12 px-4"
                  >
                    <X className="h-5 w-5" />
                  </Button>
                </div>
              </div>
              {error && (
                <p className="text-xs text-destructive font-semibold ml-1">{error}</p>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
