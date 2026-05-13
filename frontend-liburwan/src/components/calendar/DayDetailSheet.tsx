import { useState, useEffect } from "react"
import { useQueryClient } from "@tanstack/react-query"
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet"
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
} from "@/components/ui/drawer"
import { Alert, AlertDescription } from "@/components/ui/alert"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { formatInTimeZone } from "date-fns-tz"
import { id as localeID } from "date-fns/locale"
import type { CalendarDayData } from "@/types/calendar"
import { useAuth } from "@/hooks/useAuth"
import { isBookingWindow } from "@/lib/date-utils"
import { Calendar, User, Store, AlertCircle, PlusCircle, ArrowLeft, Loader2, ShieldCheck, Trash2 } from "lucide-react"
import { useMediaQuery } from "@/hooks/useMediaQuery"
import { AjukanLiburForm } from "./AjukanLiburForm"
import { UnplannedLeaveForm } from "./UnplannedLeaveForm"
import { AssignBackupForm } from "./AssignBackupForm"
import { apiFetch } from "@/lib/api"

interface DayDetailSheetProps {
  isOpen: boolean
  onClose: () => void
  date: Date | null
  data?: CalendarDayData
}

export function DayDetailSheet({
  isOpen,
  onClose,
  date,
  data,
}: DayDetailSheetProps) {
  const queryClient = useQueryClient()
  const { karyawan, isAdmin } = useAuth()
  const isDesktop = useMediaQuery("(min-width: 768px)")
  const [view, setView] = useState<"detail" | "create" | "unplanned" | "assign-backup">("detail")
  const [selectedJadwal, setSelectedJadwal] = useState<{ id: string; nama: string; toko: string } | null>(null)
  const [isCancelDialogOpen, setIsCancelDialogOpen] = useState(false)
  const [isDeleteBackupDialogOpen, setIsDeleteBackupDialogOpen] = useState(false)
  const [selectedBackupId, setSelectedBackupId] = useState<string | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)
  const [isDeletingBackup, setIsDeletingBackup] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [isDeleteUnplannedDialogOpen, setIsDeleteUnplannedDialogOpen] = useState(false)
  const [selectedUnplannedId, setSelectedUnplannedId] = useState<string | null>(null)

  // Reset view when sheet opens/closes or date changes
  useEffect(() => {
    let timer: ReturnType<typeof setTimeout>
    if (!isOpen) {
      timer = setTimeout(() => {
        setView("detail")
        setSelectedJadwal(null)
        setSelectedBackupId(null)
        setSelectedUnplannedId(null)
        setError(null)
      }, 300)
    }
    return () => clearTimeout(timer)
  }, [isOpen])

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setView("detail")
    setError(null)
  }, [date])

  if (!date) return null

  const dateDisplay = formatInTimeZone(date, "Asia/Jakarta", "EEEE, d MMMM yyyy", { locale: localeID })
  const canAction = isBookingWindow(date)

  const userLeave = data?.toko
    .flatMap((t) => t.libur)
    .find((l) => l.karyawan_id === karyawan?.id)

  const allLeaves = data?.toko.flatMap(t => t.libur.map(l => ({ ...l, toko_nama: t.toko_nama, is_rawan: t.is_rawan }))) || []
  const leavesNeedingBackup = allLeaves.filter(l => !l.backup_assignment_id && l.is_rawan)
  const onLeaveKaryawanIds = allLeaves.map(l => l.karyawan_id)

  const handleCancelLeave = async () => {
    if (!userLeave) return
    try {
      setIsDeleting(true)
      setError(null)
      const response = await apiFetch(`/jadwal-libur/${userLeave.id}`, {
        method: "DELETE",
      })

      if (!response.ok) {
        const errData = await response.json()
        setError(`Gagal membatalkan: ${errData.message || response.statusText}`)
        return
      }

      await queryClient.invalidateQueries({ queryKey: ["kalender"] })
      await queryClient.invalidateQueries({ queryKey: ["jadwal-libur"] })
      onClose()
    } catch (err) {
      setError(`Terjadi kesalahan: ${err instanceof Error ? err.message : String(err)}`)
      console.error(err)
    } finally {
      setIsDeleting(false)
      setIsCancelDialogOpen(false)
    }
  }

  const handleDeleteUnplannedLeave = async () => {
    if (!selectedUnplannedId) return
    try {
      setIsDeleting(true)
      setError(null)
      const response = await apiFetch(`/jadwal-libur/${selectedUnplannedId}`, {
        method: "DELETE",
      })

      if (!response.ok) {
        const errData = await response.json()
        setError(`Gagal menghapus unplanned leave: ${errData.message || response.statusText}`)
        return
      }

      await queryClient.invalidateQueries({ queryKey: ["kalender"] })
      await queryClient.invalidateQueries({ queryKey: ["jadwal-libur"] })
      setIsDeleteUnplannedDialogOpen(false)
      onClose()
    } catch (err) {
      setError(`Terjadi kesalahan: ${err instanceof Error ? err.message : String(err)}`)
    } finally {
      setIsDeleting(false)
    }
  }

  const handleDeleteBackup = async () => {
    if (!selectedBackupId) return
    try {
      setIsDeletingBackup(true)
      setError(null)
      const response = await apiFetch(`/backup-assignment/${selectedBackupId}`, {
        method: "DELETE",
      })

      if (!response.ok) {
        const errData = await response.json()
        setError(`Gagal menghapus backup: ${errData.message || response.statusText}`)
        return
      }

      await queryClient.invalidateQueries({ queryKey: ["kalender"] })
      await queryClient.invalidateQueries({ queryKey: ["jadwal-libur"] })
      setIsDeleteBackupDialogOpen(false)
    } catch (err) {
      setError(`Terjadi kesalahan: ${err instanceof Error ? err.message : String(err)}`)
    } finally {
      setIsDeletingBackup(false)
    }
  }

  const DetailContent = (
    <div className="space-y-8">
      {error && (
        <Alert variant="destructive" className="rounded-2xl border-destructive/20 bg-destructive/5">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription className="text-xs font-medium">
            {error}
          </AlertDescription>
        </Alert>
      )}
      <div className="space-y-4">
        <h3 className="flex items-center gap-2 text-sm font-bold uppercase tracking-wider text-muted-foreground">
          <User className="h-4 w-4" />
          Karyawan Libur
        </h3>

        {(!data || data.toko.every(t => t.libur.length === 0)) ? (
          <div className="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">
            Tidak ada karyawan yang libur di hari ini.
          </div>
        ) : (
          <div className="space-y-6">
            {data.toko.map((t) => (
              t.libur.length > 0 && (
                <div key={t.toko_id} className="space-y-2">
                  <div className="flex items-center justify-between border-b pb-1">
                    <div className="flex items-center gap-2 text-sm font-bold">
                      <Store className="h-4 w-4 text-primary/70" />
                      {t.toko_nama}
                    </div>
                    {t.is_rawan && (
                      <Badge variant="destructive" className="h-5 text-[10px] uppercase">
                        Rawan
                      </Badge>
                    )}
                  </div>
                  <ul className="space-y-3 pl-6">
                    {t.libur.map((l) => (
                      <li key={l.karyawan_id} className="flex flex-col gap-1">
                        <div className="flex items-center justify-between text-sm">
                          <span className="font-medium text-foreground">{l.karyawan_nama}</span>
                          <div className="flex items-center gap-2">
                            <Badge variant={l.tipe === "planned" ? "outline" : "secondary"} className="h-5 text-[9px] uppercase">
                              {l.tipe}
                            </Badge>
                            {isAdmin() && l.tipe === "unplanned" && (
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => {
                                  setSelectedUnplannedId(l.id)
                                  setIsDeleteUnplannedDialogOpen(true)
                                }}
                                className="h-6 w-6 text-destructive/50 hover:text-destructive hover:bg-destructive/5"
                              >
                                <Trash2 className="h-3 w-3" />
                              </Button>
                            )}
                          </div>
                        </div>
                        {l.backup_karyawan_nama && (
                          <div className="flex items-center justify-between">
                            <div className="text-xs text-muted-foreground flex items-center gap-1">
                              <span className="h-1 w-1 rounded-full bg-success" />
                              Backup: {l.backup_karyawan_nama}
                            </div>
                            {isAdmin() && l.backup_assignment_id && (
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => {
                                  setSelectedBackupId(l.backup_assignment_id!)
                                  setIsDeleteBackupDialogOpen(true)
                                }}
                                className="h-6 w-6 text-destructive/50 hover:text-destructive hover:bg-destructive/5"
                              >
                                <Trash2 className="h-3 w-3" />
                              </Button>
                            )}
                          </div>
                        )}
                      </li>
                    ))}
                  </ul>
                </div>
              )
            ))}
          </div>
        )}
      </div>

      {isAdmin() && (
        <div className="space-y-4 border-t pt-6">
          <h3 className="flex items-center gap-2 text-sm font-bold uppercase tracking-wider text-muted-foreground">
            <ShieldCheck className="h-4 w-4" />
            Backup Assignment
          </h3>

          {leavesNeedingBackup.length === 0 ? (
            <div className="rounded-lg bg-success/5 p-4 text-center text-xs text-success font-medium border border-success/10 border-dashed">
              Semua jadwal libur kritis sudah memiliki backup.
            </div>
          ) : (
            <div className="space-y-2">
              {leavesNeedingBackup.map((l) => (
                <div key={l.id} className="flex items-center justify-between rounded-xl border bg-muted/30 p-3 text-xs">
                  <div className="space-y-1">
                    <div className="font-medium">{l.karyawan_nama}</div>
                    <div className="flex items-center gap-1 text-muted-foreground text-[10px]">
                      <Store className="h-3 w-3" />
                      {l.toko_nama} • {l.tipe}
                    </div>
                  </div>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => {
                      setSelectedJadwal({ id: l.id, nama: l.karyawan_nama, toko: l.toko_nama })
                      setView("assign-backup")
                    }}
                    className="h-8 rounded-lg text-[10px] font-bold"
                  >
                    Assign Backup
                  </Button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      <div className="space-y-4 border-t pt-6">
        <h3 className="flex items-center gap-2 text-sm font-bold uppercase tracking-wider text-muted-foreground">
          <AlertCircle className="h-4 w-4" />
          Aksi
        </h3>

        <div className="grid gap-3">
          {canAction && (
            <>
              {!userLeave ? (
                <Button onClick={() => setView("create")} className="w-full gap-2 rounded-xl h-12 font-bold shadow-lg shadow-primary/20">
                  <PlusCircle className="h-4 w-4" />
                  Ajukan Libur
                </Button>
              ) : (
                <Button variant="destructive" onClick={() => setIsCancelDialogOpen(true)} className="w-full h-12 rounded-xl shadow-lg shadow-destructive/10">
                  Batalkan Jadwal Libur
                </Button>
              )}
            </>
          )}

          {!canAction && !isAdmin() && (
            <div className="rounded-lg bg-muted/50 p-4 text-center text-xs text-muted-foreground border border-dashed">
              Window pengajuan untuk tanggal ini sudah ditutup atau belum dibuka.
            </div>
          )}

          {isAdmin() && (
            <Button
              variant="secondary"
              onClick={() => setView("unplanned")}
              className="w-full gap-2 rounded-xl h-12 mt-2 font-medium border-primary/10 hover:bg-primary/5"
            >
              <User className="h-4 w-4" />
              Input Unplanned Leave
            </Button>
          )}
        </div>
      </div>
    </div>
  )

  const FormContent = (
    <div>
      {view === "create" && (
        <AjukanLiburForm
          date={date}
          onSuccess={onClose}
          onCancel={() => setView("detail")}
        />
      )}
      {view === "unplanned" && (
        <UnplannedLeaveForm
          date={date}
          onSuccess={onClose}
          onCancel={() => setView("detail")}
        />
      )}
      {view === "assign-backup" && selectedJadwal && (
        <AssignBackupForm
          jadwalLiburId={selectedJadwal.id}
          karyawanNama={selectedJadwal.nama}
          tokoNama={selectedJadwal.toko}
          onLeaveKaryawanIds={onLeaveKaryawanIds}
          onSuccess={onClose}
          onCancel={() => setView("detail")}
        />
      )}
    </div>
  )

  const titleMap = {
    detail: "Detail Jadwal",
    create: "Form Ajukan Libur",
    unplanned: "Input Unplanned Leave",
    "assign-backup": "Assign Backup",
  }

  const Header = (
    <>
      <SheetTitle className="flex items-center gap-2">
        {view !== "detail" && (
          <Button variant="ghost" size="icon" onClick={() => setView("detail")} className="-ml-2 h-8 w-8 rounded-full">
            <ArrowLeft className="h-4 w-4" />
          </Button>
        )}
        <Calendar className="h-5 w-5 text-primary" />
        {titleMap[view]}
      </SheetTitle>
      <SheetDescription className="font-medium text-foreground/70">{dateDisplay}</SheetDescription>
    </>
  )

  return (
    <>
      {isDesktop ? (
        <Sheet open={isOpen} onOpenChange={onClose}>
          <SheetContent className="sm:max-w-md flex flex-col h-full p-0">
            <SheetHeader className="p-6 border-b bg-background/80 backdrop-blur-md sticky top-0 z-10">
              {Header}
            </SheetHeader>
            <div className="flex-1 overflow-y-auto p-6 scrollbar-thin scrollbar-thumb-muted-foreground/20 hover:scrollbar-thumb-muted-foreground/40">
              {view === "detail" ? DetailContent : FormContent}
            </div>
          </SheetContent>
        </Sheet>
      ) : (
        <Drawer open={isOpen} onOpenChange={onClose}>
          <DrawerContent>
            <DrawerHeader className="text-left">
              {Header}
            </DrawerHeader>
            <div className="max-h-[80vh] overflow-y-auto px-6 pb-8 scrollbar-thin scrollbar-thumb-muted-foreground/20 hover:scrollbar-thumb-muted-foreground/40">
              {view === "detail" ? DetailContent : FormContent}
            </div>
          </DrawerContent>
        </Drawer>
      )}

      {/* Dialog Batal Libur */}
      <AlertDialog open={isCancelDialogOpen} onOpenChange={setIsCancelDialogOpen}>
        <AlertDialogContent className="rounded-3xl border-none shadow-2xl">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-xl font-bold">Batalkan Jadwal Libur?</AlertDialogTitle>
            <AlertDialogDescription className="text-sm">
              {userLeave?.backup_karyawan_nama
                ? "Jadwal ini memiliki backup assignment. Membatalkan akan menghapus backup tersebut secara otomatis."
                : "Yakin ingin membatalkan jadwal libur Anda di tanggal ini? Tindakan ini tidak dapat dibatalkan."}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="gap-3">
            <AlertDialogCancel className="rounded-xl h-11 flex-1">Tidak, Kembali</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleCancelLeave}
              disabled={isDeleting}
              className="rounded-xl h-11 flex-1 bg-destructive hover:bg-destructive/90 text-white gap-2"
            >
              {isDeleting && <Loader2 className="h-4 w-4 animate-spin" />}
              Ya, Batalkan
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Dialog Hapus Backup */}
      <AlertDialog open={isDeleteBackupDialogOpen} onOpenChange={setIsDeleteBackupDialogOpen}>
        <AlertDialogContent className="rounded-3xl border-none shadow-2xl">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-xl font-bold">Hapus Backup Assignment?</AlertDialogTitle>
            <AlertDialogDescription className="text-sm">
              Tindakan ini akan menghapus tugas backup untuk jadwal libur ini. Karyawan yang bersangkutan akan kembali ke status 'Butuh Backup'.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="gap-3">
            <AlertDialogCancel className="rounded-xl h-11 flex-1">Batal</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteBackup}
              disabled={isDeletingBackup}
              className="rounded-xl h-11 flex-1 bg-destructive hover:bg-destructive/90 text-white gap-2"
            >
              {isDeletingBackup && <Loader2 className="h-4 w-4 animate-spin" />}
              Ya, Hapus
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Dialog Hapus Unplanned */}
      <AlertDialog open={isDeleteUnplannedDialogOpen} onOpenChange={setIsDeleteUnplannedDialogOpen}>
        <AlertDialogContent className="rounded-3xl border-none shadow-2xl">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-xl font-bold">Hapus Unplanned Leave?</AlertDialogTitle>
            <AlertDialogDescription className="text-sm">
              Yakin ingin menghapus unplanned leave ini? Tindakan ini tidak dapat dibatalkan.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="gap-3">
            <AlertDialogCancel className="rounded-xl h-11 flex-1">Batal</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteUnplannedLeave}
              disabled={isDeleting}
              className="rounded-xl h-11 flex-1 bg-destructive hover:bg-destructive/90 text-white gap-2"
            >
              {isDeleting && <Loader2 className="h-4 w-4 animate-spin" />}
              Ya, Hapus
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
