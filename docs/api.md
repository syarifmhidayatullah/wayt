# Wayt API Documentation

Wayt adalah sistem manajemen antrian berbasis QR code. Dokumen ini menjelaskan semua endpoint API, alur penggunaan, dan referensi teknis.

---

## Base URL

```
http://localhost:8080
```

Ganti dengan domain/IP sesuai environment.

---

## Authentication

Semua endpoint `/internal/*` membutuhkan JWT token di header:

```
Authorization: Bearer <token>
```

Token didapat dari endpoint `POST /auth/login`. Token berlaku selama **8 jam**.

---

## Response Format

Semua response menggunakan format JSON yang konsisten:

**Success:**
```json
{
  "success": true,
  "message": "pesan sukses",
  "data": {}
}
```

**Error:**
```json
{
  "success": false,
  "message": "pesan error",
  "error": "detail error"
}
```

**HTTP Status Codes:**

| Code | Keterangan |
|------|------------|
| `200` | OK |
| `201` | Created |
| `400` | Bad Request — validasi gagal atau business logic error |
| `401` | Unauthorized — token tidak ada atau tidak valid |
| `403` | Forbidden — role tidak cukup |
| `404` | Not Found |
| `500` | Internal Server Error |

---

## Role

| Role | Deskripsi |
|------|-----------|
| `superadmin` | Akses penuh termasuk kelola user admin |
| `admin` | Manage branches dan antrian |
| `public` | Tidak perlu login |

---

## Endpoints

### Auth

#### POST /auth/login

Login dan dapatkan JWT token.

**Request:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**JWT Claims:**
```json
{
  "sub": 1,
  "username": "admin",
  "role": "superadmin",
  "exp": 1234567890
}
```

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | username dan password wajib diisi | Field kosong |
| 401 | unauthorized | Username atau password salah |

---

### User Management

> Semua endpoint ini hanya bisa diakses oleh `superadmin`.

#### GET /internal/users

List semua admin user.

**Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": [
    {
      "id": 1,
      "username": "admin",
      "role": "superadmin",
      "created_at": "2026-04-01T10:00:00Z",
      "updated_at": "2026-04-01T10:00:00Z"
    }
  ]
}
```

---

#### POST /internal/users

Buat admin user baru.

**Request:**
```json
{
  "username": "staff1",
  "password": "password123",
  "role": "admin"
}
```

> `role` opsional, default `admin`. Nilai valid: `admin`, `superadmin`.

**Response (201):**
```json
{
  "success": true,
  "message": "user created",
  "data": {
    "id": 2,
    "username": "staff1",
    "role": "admin",
    "created_at": "2026-04-01T10:00:00Z",
    "updated_at": "2026-04-01T10:00:00Z"
  }
}
```

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | username dan password wajib diisi | Field kosong |
| 400 | username sudah digunakan | Username duplikat |

---

#### PUT /internal/users/:id

Update user (username, password, atau role). Field yang kosong tidak diupdate.

**Request:**
```json
{
  "username": "staff1_updated",
  "password": "newpassword",
  "role": "superadmin"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "user updated",
  "data": {
    "id": 2,
    "username": "staff1_updated",
    "role": "superadmin",
    "created_at": "2026-04-01T10:00:00Z",
    "updated_at": "2026-04-02T08:00:00Z"
  }
}
```

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | invalid id | ID bukan angka |
| 400 | user tidak ditemukan | ID tidak ada di database |

---

#### DELETE /internal/users/:id

Hapus admin user.

**Response (200):**
```json
{
  "success": true,
  "message": "user deleted",
  "data": null
}
```

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | tidak bisa menghapus akun sendiri | ID = ID diri sendiri |
| 400 | user tidak ditemukan | ID tidak ada |

---

### Branch Management

> Dapat diakses oleh `admin` dan `superadmin`.

#### GET /internal/branches

List semua branch (tidak termasuk yang sudah dihapus).

**Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Kasir 1",
      "prefix": "A",
      "is_active": true,
      "current_number": 5,
      "last_number": 10,
      "created_at": "2026-04-01T08:00:00Z",
      "updated_at": "2026-04-01T09:00:00Z"
    }
  ]
}
```

---

#### POST /internal/branches

Buat branch baru.

**Request:**
```json
{
  "name": "Kasir 1",
  "prefix": "A"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "branch created",
  "data": {
    "id": 1,
    "name": "Kasir 1",
    "prefix": "A",
    "is_active": true,
    "current_number": 0,
    "last_number": 0,
    "created_at": "2026-04-01T08:00:00Z",
    "updated_at": "2026-04-01T08:00:00Z"
  }
}
```

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | invalid request | Field `name` atau `prefix` kosong |

---

#### PUT /internal/branches/:id

Update branch. Field yang kosong tidak diupdate, kecuali `is_active`.

**Request:**
```json
{
  "name": "Kasir 1 Updated",
  "prefix": "A",
  "is_active": false
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "branch updated",
  "data": { ... }
}
```

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | branch not found | ID tidak ada |

---

#### DELETE /internal/branches/:id

Soft delete branch (data tidak benar-benar dihapus dari database).

**Response (200):**
```json
{
  "success": true,
  "message": "branch deleted",
  "data": null
}
```

---

### QR Code

#### POST /internal/branches/:id/qr

Generate QR code untuk branch. QR berisi URL yang saat di-scan akan otomatis mendaftarkan antrian.

**Response (201):**
```json
{
  "success": true,
  "message": "QR code generated",
  "data": {
    "token": "550e8400-e29b-41d4-a716-446655440000",
    "qr_image_url": "http://localhost:8080/storage/qr/550e8400.png",
    "expired_at": "2026-04-05T08:00:00Z"
  }
}
```

> QR image berisi URL: `{PUBLIC_BASE_URL}/q/{token}`. QR yang sama bisa dipakai berkali-kali oleh banyak orang.

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | branch not found | Branch tidak ada |
| 400 | branch is not active | Branch sedang nonaktif |

---

### Queue — Internal

#### GET /internal/branches/:id/queue

List antrian aktif (status `waiting` dan `called`) di branch.

**Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": [
    {
      "id": 10,
      "branch_id": 1,
      "qr_token": "550e8400-...",
      "queue_number": "A-001",
      "status": "waiting",
      "created_at": "2026-04-04T09:00:00Z",
      "updated_at": "2026-04-04T09:00:00Z"
    }
  ]
}
```

---

#### POST /internal/branches/:id/next

Panggil antrian berikutnya. Status antrian berubah dari `waiting` ke `called`, `current_number` branch bertambah 1.

**Response (200):**
```json
{
  "success": true,
  "message": "next queue called",
  "data": {
    "id": 10,
    "branch_id": 1,
    "queue_number": "A-001",
    "status": "called",
    "created_at": "2026-04-04T09:00:00Z",
    "updated_at": "2026-04-04T09:05:00Z"
  }
}
```

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | no waiting queue found | Tidak ada antrian yang menunggu |

---

#### POST /internal/branches/:id/reset

Reset antrian branch. Semua antrian `waiting` di-expire, semua QR dinonaktifkan, counter `current_number` dan `last_number` direset ke 0.

**Response (200):**
```json
{
  "success": true,
  "message": "queue reset successfully",
  "data": null
}
```

---

### Queue — Public

#### GET /q/:token

Endpoint yang dibuka saat HP scan QR code. Sistem otomatis mendaftarkan antrian baru lalu redirect ke halaman status.

**Flow:**
```
Scan QR  →  GET /q/{token}  →  register antrian  →  redirect ke /queue/{id}
```

**Error:** Menampilkan halaman HTML dengan pesan error jika token tidak valid/expired.

---

#### GET /queue/:id *(HTML)*

Halaman status antrian untuk pengunjung. Auto-refresh setiap 5 detik.

**Menampilkan:**
- Nomor antrian
- Status (Menunggu / Dipanggil / Selesai)
- Nomor yang sedang dilayani
- Jumlah orang di depan

---

#### POST /api/queue/register

Daftar antrian via API (alternatif selain scan QR).

**Request:**
```json
{
  "qr_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "queue registered",
  "data": {
    "queue_id": 10,
    "queue_number": "A-007",
    "branch_name": "Kasir 1",
    "position": 3,
    "people_ahead": 2
  }
}
```

**Errors:**
| Status | Message | Penyebab |
|--------|---------|----------|
| 400 | QR code not found | Token tidak ada |
| 400 | QR code is no longer active | QR sudah direset |
| 400 | QR code has expired | QR sudah lewat expired_at |
| 400 | branch not found | Branch dihapus |

---

#### GET /api/queue/:token/status

Cek status antrian berdasarkan QR token.

**Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": {
    "queue_number": "A-007",
    "status": "waiting",
    "current_serving": "A-005",
    "people_ahead": 2
  }
}
```

---

#### GET /api/queue/id/:id/status

Cek status antrian berdasarkan queue ID. Dipakai oleh halaman `/queue/:id` untuk polling.

**Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": {
    "queue_number": "A-007",
    "status": "called",
    "current_serving": "A-007",
    "people_ahead": 0
  }
}
```

---

## Queue Status Flow

```
waiting  →  called  →  done
   └──────────────────→  expired  (saat branch di-reset)
```

| Status | Keterangan |
|--------|------------|
| `waiting` | Menunggu dipanggil |
| `called` | Dipanggil, sedang dilayani |
| `done` | Selesai |
| `expired` | Kadaluarsa (branch di-reset) |

---

## Alur Lengkap

### Alur Admin

```
1. Login           →  POST /auth/login
2. Buat branch     →  POST /internal/branches
3. Generate QR     →  POST /internal/branches/:id/qr
4. Print/display QR code
5. Panggil antrian →  POST /internal/branches/:id/next
6. Lihat antrian   →  GET  /internal/branches/:id/queue
7. Reset (akhir hari) → POST /internal/branches/:id/reset
```

### Alur Pengunjung

```
1. Scan QR code    →  GET /q/:token  (otomatis)
2. Terima nomor antrian + redirect ke /queue/:id
3. Pantau status   →  halaman auto-refresh tiap 5 detik
4. Tunggu notifikasi "Dipanggil"
```

---

## Pages

| URL | Keterangan |
|-----|------------|
| `/admin` | Dashboard admin (login required) |
| `/queue/:id` | Halaman status antrian pengunjung |
