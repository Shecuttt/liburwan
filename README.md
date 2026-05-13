# 🏖️ Liburwan (Libur Karyawan)

**Liburwan** adalah sistem manajemen cuti dan shift karyawan yang dirancang untuk memudahkan proses administrasi kehadiran dalam satu platform terintegrasi. Proyek ini dibangun menggunakan arsitektur monorepo untuk menyatukan ekosistem frontend dan backend.

---

## 🚀 Fitur Utama
- **Manajemen Shift**: Pengaturan jadwal kerja yang fleksibel.
- **Pengajuan Cuti**: Sistem approval cuti yang transparan.
- **Manajemen Karyawan**: Database karyawan yang terorganisir.
- **Role-Based Access Control**: Keamanan akses berdasarkan peran (Admin/Karyawan).

---

## 🛠️ Tech Stack

### Backend
- **Bahasa**: Go (Golang)
- **Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: PostgreSQL with [GORM](https://gorm.io/)
- **Migration**: golang-migrate
- **Auth**: JWT (JSON Web Token)

### Frontend
- **Framework**: [React](https://reactjs.org/) + [Vite](https://vitejs.dev/)
- **Bahasa**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: [shadcn/ui](https://ui.shadcn.com/)

---

## 📁 Struktur Proyek
```text
liburwan/
├── backend-liburwan/  # Kode sumber Backend (Go)
├── frontend-liburwan/ # Kode sumber Frontend (React)
└── README.md          # Dokumentasi utama
```

---

## 🏁 Cara Menjalankan

### Prerequisites
- Go 1.26+
- Node.js & npm
- PostgreSQL

### 1. Jalankan Backend
```bash
cd backend-liburwan
# Copy dan sesuaikan env
cp .env.example .env
# Jalankan aplikasi
go run cmd/api/main.go
```

### 2. Jalankan Frontend
```bash
cd frontend-liburwan
# Install dependencies
npm install
# Jalankan mode development
npm run dev
```

---

## 📄 Lisensi
Distribusi di bawah lisensi MIT. Lihat `LICENSE` untuk informasi lebih lanjut.

---

Dibuat pakai [Antigravity](https://github.com/google-deepmind)

AOWKAOWKAWOKAWKAOWKAOWK