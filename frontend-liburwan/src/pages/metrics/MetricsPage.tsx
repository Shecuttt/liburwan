import { useState, useMemo } from "react"
import { useQuery } from "@tanstack/react-query"
import { format, subMonths, startOfMonth, parseISO, isValid } from "date-fns"
import { id as localeID } from "date-fns/locale"
import {
  BarChart3,
  Users,
  Store as StoreIcon,
  Calendar as CalendarIcon,
  Search,
  User as UserIcon,
  AlertTriangle,
  HeartHandshake,
  Clock,
  ArrowLeft
} from "lucide-react"
import { Link } from "react-router-dom"

import { apiFetch } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import type { KaryawanProfile } from "@/types/calendar"

export default function MetricsPage() {
  const [activeTab, setActiveTab] = useState("karyawan")

  return (
    <div className="mx-auto max-w-7xl p-4 md:p-8 space-y-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <Link to="/" className="rounded-xl border p-2 hover:bg-muted transition-colors">
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <div className="rounded-2xl bg-primary/10 p-3 text-primary">
            <BarChart3 className="h-6 w-6" />
          </div>
          <div>
            <h1 className="text-2xl font-bold tracking-tight text-foreground md:text-3xl">
              Metrik & Statistik
            </h1>
            <p className="text-sm text-muted-foreground">
              Pantau rekapitulasi libur dan ketersediaan toko.
            </p>
          </div>
        </div>
      </header>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="mb-4">
          <TabsTrigger value="karyawan" className="gap-2">
            <Users className="h-4 w-4" />
            Per Karyawan
          </TabsTrigger>
          <TabsTrigger value="toko" className="gap-2">
            <StoreIcon className="h-4 w-4" />
            Per Toko
          </TabsTrigger>
        </TabsList>

        <TabsContent value="karyawan">
          <KaryawanMetricsTab />
        </TabsContent>
        <TabsContent value="toko">
          <TokoMetricsTab />
        </TabsContent>
      </Tabs>
    </div>
  )
}

function KaryawanMetricsTab() {
  const [selectedKaryawanId, setSelectedKaryawanId] = useState<string>("")
  const [selectedMonth, setSelectedMonth] = useState<string>("all")
  const [filter, setFilter] = useState<{ id: string; month: string } | null>(null)

  const months = useMemo(() => {
    const list = []
    const now = startOfMonth(new Date())
    for (let i = 0; i < 7; i++) {
      const d = subMonths(now, i)
      list.push({
        value: format(d, "yyyy-MM"),
        label: format(d, "MMMM yyyy", { locale: localeID })
      })
    }
    return list
  }, [])

  const { data: karyawans } = useQuery<KaryawanProfile[]>({
    queryKey: ["karyawan-list"],
    queryFn: async () => {
      const res = await apiFetch("/karyawan")
      return res.json()
    }
  })

  const { data: metrics, isLoading: isLoadingMetrics, isFetching } = useQuery({
    queryKey: ["metrik-karyawan", filter],
    queryFn: async () => {
      if (!filter?.id) return null
      const url = `/metrik/karyawan/${filter.id}${filter.month !== "all" ? `?bulan=${filter.month}` : ""}`
      const res = await apiFetch(url)
      if (!res.ok) throw new Error("Gagal memuat metrik")
      return res.json()
    },
    enabled: !!filter?.id
  })

  const handleSearch = () => {
    if (selectedKaryawanId) {
      setFilter({ id: selectedKaryawanId, month: selectedMonth })
    }
  }

  return (
    <div className="space-y-6">
      <Card className="rounded-3xl border-none bg-muted/30">
        <CardContent className="p-6">
          <div className="flex flex-wrap items-end gap-4">
            <div className="space-y-2 flex-1 min-w-[240px]">
              <label className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground ml-1">
                Pilih Karyawan
              </label>
              <Select value={selectedKaryawanId} onValueChange={setSelectedKaryawanId}>
                <SelectTrigger className="rounded-xl h-12 bg-background">
                  <SelectValue placeholder="Cari karyawan..." />
                </SelectTrigger>
                <SelectContent className="rounded-xl max-h-[300px]">
                  {karyawans?.map((k, idx) => (
                    <SelectItem key={k.id || `karyawan-${idx}`} value={k.id}>
                      {k.nama} ({k.toko_nama})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2 w-[180px]">
              <label className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground ml-1">
                Bulan
              </label>
              <Select value={selectedMonth} onValueChange={setSelectedMonth}>
                <SelectTrigger className="rounded-xl h-12 bg-background">
                  <SelectValue placeholder="Pilih bulan" />
                </SelectTrigger>
                <SelectContent className="rounded-xl">
                  <SelectItem value="all">Semua</SelectItem>
                  {months.map(m => (
                    <SelectItem key={m.value} value={m.value}>{m.label}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <Button
              onClick={handleSearch}
              disabled={!selectedKaryawanId || isFetching}
              className="rounded-xl h-12 px-8 gap-2 font-bold"
            >
              <Search className="h-4 w-4" />
              Lihat
            </Button>
          </div>
        </CardContent>
      </Card>

      {isLoadingMetrics || isFetching ? (
        <div className="space-y-6">
          <div className="grid gap-4 sm:grid-cols-3">
            {[1, 2, 3].map(i => <Skeleton key={i} className="h-32 rounded-3xl" />)}
          </div>
          <Skeleton className="h-[400px] rounded-3xl" />
        </div>
      ) : metrics ? (
        <div className="space-y-6 animate-in fade-in duration-500">
          <div className="grid gap-4 sm:grid-cols-3">
            <Card className="rounded-3xl border-none bg-primary/5 ring-1 ring-primary/10">
              <CardContent className="p-6">
                <div className="flex items-center gap-3 mb-2 text-primary">
                  <CalendarIcon className="h-5 w-5" />
                  <span className="text-sm font-bold uppercase tracking-tight">Planned Leave</span>
                </div>
                <div className="text-4xl font-black">{metrics.ringkasan.total_planned}</div>
                <p className="text-xs text-muted-foreground mt-1">Hari libur terjadwal</p>
              </CardContent>
            </Card>
            <Card className="rounded-3xl border-none bg-warning/5 ring-1 ring-warning/10">
              <CardContent className="p-6">
                <div className="flex items-center gap-3 mb-2 text-warning">
                  <Clock className="h-5 w-5" />
                  <span className="text-sm font-bold uppercase tracking-tight">Unplanned</span>
                </div>
                <div className="text-4xl font-black">{metrics.ringkasan.total_unplanned}</div>
                <p className="text-xs text-muted-foreground mt-1">Tidak masuk dadakan</p>
              </CardContent>
            </Card>
            <Card className="rounded-3xl border-none bg-success/5 ring-1 ring-success/10">
              <CardContent className="p-6">
                <div className="flex items-center gap-3 mb-2 text-success">
                  <HeartHandshake className="h-5 w-5" />
                  <span className="text-sm font-bold uppercase tracking-tight">Jadi Backup</span>
                </div>
                <div className="text-4xl font-black">{metrics.ringkasan.total_jadi_backup}</div>
                <p className="text-xs text-muted-foreground mt-1">Membantu cabang lain</p>
              </CardContent>
            </Card>
          </div>

          <Card className="rounded-3xl border shadow-none overflow-hidden">
            <CardHeader className="border-b bg-muted/10">
              <CardTitle className="text-lg">Histori Libur</CardTitle>
              <CardDescription>Menampilkan semua catatan libur dan absen.</CardDescription>
            </CardHeader>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Tanggal</TableHead>
                  <TableHead>Tipe</TableHead>
                  <TableHead>Backup Assignment</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(!metrics.histori || metrics.histori.length === 0) ? (
                  <TableRow>
                    <TableCell colSpan={3} className="h-32 text-center text-muted-foreground">
                      Tidak ada data histori.
                    </TableCell>
                  </TableRow>
                ) : (
                  metrics.histori.map((h: { id?: string; tanggal?: string; tipe?: string; backup_karyawan_nama?: string }, idx: number) => (
                    <TableRow key={h.id || `histori-${idx}`}>
                      <TableCell className="font-medium">
                        {(() => {
                          if (!h.tanggal) return "Tanggal tidak tersedia"
                          const date = parseISO(h.tanggal)
                          return isValid(date) 
                            ? format(date, "EEEE, d MMMM yyyy", { locale: localeID })
                            : "Format tanggal salah"
                        })()}
                      </TableCell>
                      <TableCell>
                        <Badge variant={h.tipe === "planned" ? "outline" : "secondary"} className="uppercase text-[10px]">
                          {h.tipe}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {h.backup_karyawan_nama ? (
                          <div className="flex items-center gap-2">
                            <UserIcon className="h-3 w-3 text-success" />
                            {h.backup_karyawan_nama}
                          </div>
                        ) : (
                          <span className="text-muted-foreground italic">-</span>
                        )}
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </Card>
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center h-64 border border-dashed rounded-3xl bg-muted/5 text-muted-foreground">
          <Search className="h-12 w-12 mb-4 opacity-20" />
          <p className="text-sm">Silakan pilih karyawan dan bulan untuk melihat metrik.</p>
        </div>
      )}
    </div>
  )
}

function TokoMetricsTab() {
  const [selectedTokoId, setSelectedTokoId] = useState<string>("")
  const [selectedMonth, setSelectedMonth] = useState<string>("all")
  const [filter, setFilter] = useState<{ id: string; month: string } | null>(null)

  const months = useMemo(() => {
    const list = []
    const now = startOfMonth(new Date())
    for (let i = 0; i < 7; i++) {
      const d = subMonths(now, i)
      list.push({
        value: format(d, "yyyy-MM"),
        label: format(d, "MMMM yyyy", { locale: localeID })
      })
    }
    return list
  }, [])

  const { data: tokos } = useQuery<{ id?: string; ID?: string; nama?: string; Nama?: string; is_pusat?: boolean; IsPusat?: boolean }[]>({
    queryKey: ["toko-list"],
    queryFn: async () => {
      const res = await apiFetch("/toko")
      return res.json()
    }
  })

  const { data: metrics, isFetching } = useQuery({
    queryKey: ["metrik-toko", filter],
    queryFn: async () => {
      if (!filter?.id) return null
      const url = `/metrik/toko/${filter.id}${filter.month !== "all" ? `?bulan=${filter.month}` : ""}`
      const res = await apiFetch(url)
      if (!res.ok) throw new Error("Gagal memuat metrik")
      return res.json()
    },
    enabled: !!filter?.id
  })

  const handleSearch = () => {
    if (selectedTokoId) {
      setFilter({ id: selectedTokoId, month: selectedMonth })
    }
  }

  return (
    <div className="space-y-6">
      <Card className="rounded-3xl border-none bg-muted/30">
        <CardContent className="p-6">
          <div className="flex flex-wrap items-end gap-4">
            <div className="space-y-2 flex-1 min-w-[240px]">
              <label className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground ml-1">
                Pilih Toko
              </label>
              <Select value={selectedTokoId} onValueChange={setSelectedTokoId}>
                <SelectTrigger className="rounded-xl h-12 bg-background">
                  <SelectValue placeholder="Pilih cabang toko..." />
                </SelectTrigger>
                <SelectContent className="rounded-xl">
                  {tokos?.map((t) => {
                    const id = t.id || t.ID || "";
                    if (!id) return null;
                    const nama = t.nama || t.Nama;
                    const isPusat = t.is_pusat || t.IsPusat;
                    return (
                      <SelectItem key={id} value={id}>
                        {nama} {isPusat ? "(Pusat)" : ""}
                      </SelectItem>
                    );
                  })}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2 w-[180px]">
              <label className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground ml-1">
                Bulan
              </label>
              <Select value={selectedMonth} onValueChange={setSelectedMonth}>
                <SelectTrigger className="rounded-xl h-12 bg-background">
                  <SelectValue placeholder="Pilih bulan" />
                </SelectTrigger>
                <SelectContent className="rounded-xl">
                  <SelectItem value="all">Semua</SelectItem>
                  {months.map(m => (
                    <SelectItem key={m.value} value={m.value}>{m.label}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <Button
              onClick={handleSearch}
              disabled={!selectedTokoId || isFetching}
              className="rounded-xl h-12 px-8 gap-2 font-bold"
            >
              <Search className="h-4 w-4" />
              Lihat
            </Button>
          </div>
        </CardContent>
      </Card>

      {isFetching ? (
        <div className="space-y-6">
          <div className="grid gap-4 sm:grid-cols-3">
            {[1, 2, 3].map(i => <Skeleton key={i} className="h-32 rounded-3xl" />)}
          </div>
          <Skeleton className="h-[400px] rounded-3xl" />
        </div>
      ) : metrics ? (
        <div className="space-y-6 animate-in fade-in duration-500">
          <div className="grid gap-4 sm:grid-cols-3">
            <Card className="rounded-3xl border-none bg-destructive/5 ring-1 ring-destructive/10">
              <CardContent className="p-6">
                <div className="flex items-center gap-3 mb-2 text-destructive">
                  <AlertTriangle className="h-5 w-5" />
                  <span className="text-sm font-bold uppercase tracking-tight">Hari Kritis</span>
                </div>
                <div className="text-4xl font-black">{metrics.ringkasan.total_hari_kritis}</div>
                <p className="text-xs text-muted-foreground mt-1">Availability ≤ 2 karyawan</p>
              </CardContent>
            </Card>
            <Card className="rounded-3xl border-none bg-primary/5 ring-1 ring-primary/10">
              <CardContent className="p-6">
                <div className="flex items-center gap-3 mb-2 text-primary">
                  <HeartHandshake className="h-5 w-5" />
                  <span className="text-sm font-bold uppercase tracking-tight">Backup Masuk</span>
                </div>
                <div className="text-4xl font-black">{metrics.ringkasan.total_backup_dari_luar}</div>
                <p className="text-xs text-muted-foreground mt-1">Bantuan dari cabang lain</p>
              </CardContent>
            </Card>
            <Card className="rounded-3xl border-none bg-warning/5 ring-1 ring-warning/10">
              <CardContent className="p-6">
                <div className="flex items-center gap-3 mb-2 text-warning">
                  <Clock className="h-5 w-5" />
                  <span className="text-sm font-bold uppercase tracking-tight">Unplanned</span>
                </div>
                <div className="text-4xl font-black">{metrics.ringkasan.total_unplanned}</div>
                <p className="text-xs text-muted-foreground mt-1">Total absen dadakan di toko</p>
              </CardContent>
            </Card>
          </div>

          <Card className="rounded-3xl border shadow-none overflow-hidden">
            <CardHeader className="border-b bg-muted/10">
              <CardTitle className="text-lg">Log Hari Kritis</CardTitle>
              <CardDescription>Daftar tanggal di mana ketersediaan karyawan minim.</CardDescription>
            </CardHeader>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Tanggal</TableHead>
                  <TableHead>Jumlah Tersedia</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(!metrics.hari_kritis || metrics.hari_kritis.length === 0) ? (
                  <TableRow>
                    <TableCell colSpan={3} className="h-32 text-center text-muted-foreground">
                      Tidak ada data hari kritis.
                    </TableCell>
                  </TableRow>
                ) : (
                  metrics.hari_kritis?.map((hk: { tanggal?: string; available_count?: number }, i: number) => (
                    <TableRow key={i}>
                      <TableCell className="font-medium">
                        {(() => {
                          if (!hk.tanggal) return "Tanggal tidak tersedia"
                          const date = parseISO(hk.tanggal)
                          return isValid(date) 
                            ? format(date, "EEEE, d MMMM yyyy", { locale: localeID })
                            : "Format tanggal salah"
                        })()}
                      </TableCell>
                      <TableCell className="font-bold">
                        {hk.available_count} Karyawan
                      </TableCell>
                      <TableCell>
                        <Badge variant={hk.available_count === 1 ? "destructive" : "secondary"}>
                          {hk.available_count === 1 ? "KRITIS" : "RAWAN"}
                        </Badge>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </Card>
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center h-64 border border-dashed rounded-3xl bg-muted/5 text-muted-foreground">
          <Search className="h-12 w-12 mb-4 opacity-20" />
          <p className="text-sm">Silakan pilih toko dan bulan untuk melihat metrik.</p>
        </div>
      )}
    </div>
  )
}
