# Web Profil Perusahaan (Go Fiber + Tailwind CSS)

Website profil perusahaan modern yang dibangun menggunakan **Go (Fiber)** dan **Tailwind CSS**.

## Prasyarat

- [Go](https://go.dev/dl/) (versi 1.18+)
- [Node.js & npm](https://nodejs.org/)
- [PostgreSQL](https://www.postgresql.org/download/) atau Docker

## Setup Database

1. **Install PostgreSQL** atau gunakan Docker:
   ```bash
   docker run -d --name postgres -e POSTGRES_USER=cp_user -e POSTGRES_PASSWORD=password123 -e POSTGRES_DB=company_profile -p 5432:5432 postgres
   ```

2. **Jalankan Migrasi** (opsional, tabel dibuat otomatis oleh kode):
   - Connect ke database menggunakan pgAdmin atau psql.
   - Jalankan file `migrations/001_initial_schema.sql`.

## Cara Menjalankan

1. Clone repository ini.
2. Install dependensi:
   ```bash
   go mod tidy
   npm install
   ```
3. Build Tailwind CSS:
   ```bash
   npm run build
   ```
4. Jalankan server:
   ```bash
   go run main.go
   ```
5. Buka http://localhost:3000 di browser.

## Ekspor Database

Untuk backup atau share data:
- Gunakan pgAdmin: Klik kanan database > Backup.
- Atau command line:
  ```bash
  pg_dump -h localhost -p 5432 -U cp_user -d company_profile > backup.sql
  ```

## Migrasi

File migrasi ada di folder `migrations/`. Jalankan manual di pgAdmin atau psql.