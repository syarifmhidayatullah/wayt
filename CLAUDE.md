# Wayt — Queue Management via QR Code

Aplikasi antrian berbasis QR code. Monolith, solo developer, backend Go + frontend HTML embedded.

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Frontend**: HTML + Alpine.js + Tailwind CSS (via CDN, no build step)
- **QR Generator**: `skip2/go-qrcode`
- **Auth**: JWT (`golang-jwt/jwt/v5`) + bcrypt

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
├── migrations/                  # DDL SQL files (jalankan manual ke PostgreSQL)
├── pkg/
│   ├── middleware/auth.go       # JWTAuth + SuperAdminOnly middleware
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
DB_PORT=5432
DB_USER=
DB_PASSWORD=
DB_NAME=wayt

JWT_SECRET=           # Secret untuk sign JWT, wajib diganti di production
ADMIN_USERNAME=admin  # Username superadmin pertama (seed otomatis saat start)
ADMIN_PASSWORD=       # Password superadmin pertama

PUBLIC_BASE_URL=      # Base URL yang bisa diakses HP (e.g. http://192.168.0.167:8080)

QR_STORAGE_PATH=./storage/qr
QR_BASE_URL=http://localhost:8080/storage/qr
QR_EXPIRED_HOURS=24
```

## Menjalankan

```bash
# 1. Buat database
createdb wayt

# 2. Jalankan migration berurutan
psql -U wayt_svc -d wayt -f migrations/001_create_branches.sql
psql -U wayt_svc -d wayt -f migrations/002_create_qr_codes.sql
psql -U wayt_svc -d wayt -f migrations/003_create_queues.sql
psql -U wayt_svc -d wayt -f migrations/004_create_admin_users.sql
psql -U wayt_svc -d wayt -f migrations/005_add_role_to_admin_users.sql

# 3. Copy dan isi env
cp .env.example .env

# 4. Run (gunakan CGO_ENABLED=0 di macOS Sequoia/Tahoe)
CGO_ENABLED=0 go run ./cmd/api/main.go
```

Saat server pertama kali start, jika tabel `admin_users` kosong dan `ADMIN_PASSWORD` diisi, server otomatis seed user superadmin.

## API Endpoints

### Auth (public)

| Method | Endpoint | Fungsi |
|--------|----------|--------|
| POST | `/auth/login` | Login, return JWT token |

### Internal (header wajib: `Authorization: Bearer <token>`)

| Method | Endpoint | Role | Fungsi |
|--------|----------|------|--------|
| GET | `/internal/users` | superadmin | List admin users |
| POST | `/internal/users` | superadmin | Buat admin user baru |
| PUT | `/internal/users/:id` | superadmin | Update user |
| DELETE | `/internal/users/:id` | superadmin | Hapus user |
| POST | `/internal/branches` | all | Buat branch |
| GET | `/internal/branches` | all | List branch |
| PUT | `/internal/branches/:id` | all | Update branch |
| DELETE | `/internal/branches/:id` | all | Hapus branch |
| POST | `/internal/branches/:id/qr` | all | Generate QR code |
| POST | `/internal/branches/:id/next` | all | Panggil antrian berikutnya |
| GET | `/internal/branches/:id/queue` | all | List antrian aktif |
| POST | `/internal/branches/:id/reset` | all | Reset antrian & nomor |

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
Cabang/loket antrian. Punya `prefix` (e.g. "A"), `current_number` (sedang dilayani), `last_number` (total yang sudah daftar). Soft delete via kolom `deleted_at`.

### QRCode
Satu QR per generate. Punya `token` (UUID), `expired_at`, `is_active`. Bisa di-reset via endpoint reset branch.

### Queue
Satu entri per pendaftar. Status: `waiting → called → done` atau `expired` (saat reset).

### AdminUser
User untuk login ke admin panel. Role: `superadmin` atau `admin`.
- `superadmin` — akses semua fitur termasuk kelola user
- `admin` — hanya manage branches & antrian

## Auth Flow

```
POST /auth/login  →  return JWT (expire 8 jam)
JWT claims: { sub, username, role, exp }

Semua /internal/* wajib header: Authorization: Bearer <token>
/internal/users/* hanya bisa diakses role superadmin
```

JWT di-parse di frontend (Alpine.js) untuk menampilkan/menyembunyikan tab Kelola User.

## Frontend

- **`/admin`** — Login pakai username/password. JWT disimpan di `localStorage`. Tab Antrian (semua role) dan tab Kelola User (superadmin only). Auto-logout saat JWT expired (server return 401).
- **`/queue/:id`** — Halaman mobile. Tampil nomor antrian, posisi, sedang dilayani. Polling `/api/queue/id/:id/status` tiap 5 detik via Alpine.js. Polling berhenti saat status `done`/`expired`.

## Konvensi

- Response JSON selalu `{ "success": bool, "message": string, "data": ... }`
- Soft delete pada `branches` via kolom `deleted_at`
- Tidak pakai GORM AutoMigrate — DDL dikelola manual di folder `migrations/`
- Template HTML di-load via `r.LoadHTMLGlob("web/templates/*")` — path relatif dari direktori kerja saat `go run`
- Tidak bisa hapus akun sendiri (dicek di service layer)
