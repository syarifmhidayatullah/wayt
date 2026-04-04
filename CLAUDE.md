# Wayt — Queue Management via QR Code

Aplikasi antrian berbasis QR code. Monolith, solo developer, backend Go + frontend HTML embedded.

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL
- **Frontend**: HTML + Alpine.js + Tailwind CSS (via CDN, no build step)
- **QR Generator**: `skip2/go-qrcode`

## Struktur Folder

```
wayt/
├── cmd/api/main.go              # Entry point, DI manual, router setup
├── config/config.go             # Load .env via godotenv
├── internal/
│   ├── handler/                 # HTTP layer (Gin handlers)
│   ├── service/                 # Business logic
│   ├── repository/              # Database queries (GORM)
│   └── model/                   # Struct + TableName()
├── migrations/                  # DDL SQL files (jalankan manual ke MySQL)
├── pkg/
│   ├── middleware/auth.go       # InternalAuth: cek header X-Internal-Key
│   └── response/response.go    # Helper JSON response standar
├── web/templates/
│   ├── admin.html               # Internal tools dashboard
│   └── queue.html               # Halaman mobile status antrian
├── storage/qr/                  # Gambar QR yang di-generate (di-gitignore)
├── .env                         # Tidak di-commit
└── .env.example                 # Template env
```

## Arsitektur

Layered architecture: `handler → service → repository`

Dependency injection manual di `main.go` — tidak pakai DI framework.

## Environment Variables (.env)

```
APP_PORT=8080
APP_ENV=development

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=
DB_PASSWORD=
DB_NAME=wayt

INTERNAL_API_KEY=        # Key untuk semua endpoint /internal/*
PUBLIC_BASE_URL=         # Base URL yang bisa diakses HP (e.g. http://192.168.0.167:8080)

QR_STORAGE_PATH=./storage/qr
QR_BASE_URL=http://localhost:8080/storage/qr
QR_EXPIRED_HOURS=24
```

## Menjalankan

```bash
# 1. Setup database
mysql -u root -p < migrations/001_create_branches.sql
mysql -u root -p < migrations/002_create_qr_codes.sql
mysql -u root -p < migrations/003_create_queues.sql

# 2. Copy dan isi env
cp .env.example .env

# 3. Run
go run ./cmd/api/main.go
```

## API Endpoints

### Internal (header wajib: `X-Internal-Key: <value>`)

| Method | Endpoint | Fungsi |
|--------|----------|--------|
| POST | `/internal/branches` | Buat branch |
| GET | `/internal/branches` | List branch |
| PUT | `/internal/branches/:id` | Update branch |
| DELETE | `/internal/branches/:id` | Hapus branch |
| POST | `/internal/branches/:id/qr` | Generate QR code |
| POST | `/internal/branches/:id/next` | Panggil antrian berikutnya |
| GET | `/internal/branches/:id/queue` | List antrian aktif |
| POST | `/internal/branches/:id/reset` | Reset antrian & nomor |

### Public

| Method | Endpoint | Fungsi |
|--------|----------|--------|
| POST | `/api/queue/register` | Daftar antrian via token QR (JSON) |
| GET | `/api/queue/:token/status` | Status antrian by QR token (JSON) |
| GET | `/api/queue/id/:id/status` | Status antrian by queue ID (JSON) |
| GET | `/q/:token` | Scan QR → auto register → redirect ke `/queue/:id` |
| GET | `/queue/:id` | Halaman HTML status antrian per orang |
| GET | `/admin` | Dashboard internal tools (HTML) |

## Alur QR Code

```
Admin generate QR  →  QR berisi URL: {PUBLIC_BASE_URL}/q/{token}
User scan QR       →  browser buka GET /q/{token}
                   →  server register antrian baru
                   →  redirect ke /queue/{queue_id}
                   →  halaman auto-refresh tiap 5 detik
```

**1 QR bisa dipakai berkali-kali** — setiap scan membuat entri antrian baru dengan ID unik.
Setiap orang mendapat URL `/queue/{id}` yang berbeda untuk memantau posisinya.

## Model Data

### Branch
Cabang/loket antrian. Punya `prefix` (e.g. "A"), `current_number` (sedang dilayani), `last_number` (total yang sudah daftar).

### QRCode
Satu QR per generate. Punya `token` (UUID), `expired_at`, `is_active`. Bisa di-reset via endpoint reset branch.

### Queue
Satu entri per pendaftar. Status: `waiting → called → done` atau `expired` (saat reset).

## Frontend

- **`/admin`** — Login pakai API key (disimpan di `localStorage`). Fitur: CRUD branch, generate QR (tampil preview + download), panggil next, reset, lihat list antrian.
- **`/queue/:id`** — Halaman mobile. Tampil nomor antrian, posisi, sedang dilayani. Polling `/api/queue/id/:id/status` tiap 5 detik via Alpine.js.

## Konvensi

- Response JSON selalu `{ "success": bool, "message": string, "data": ... }`
- Soft delete pada `branches` via kolom `deleted_at`
- Tidak pakai GORM AutoMigrate — DDL dikelola manual di folder `migrations/`
- Template HTML di-load via `r.LoadHTMLGlob("web/templates/*")` — path relatif dari direktori kerja saat `go run`
