# erajaya-test-backend
Erajaya Technical Test - Backend

Proyek ini menggunakan pendekatan **Clean Architecture**, yang memisahkan kode berdasarkan tanggung jawabnya ke dalam beberapa layer utama:

**1. Delivery Layer (controller)**
Berfungsi sebagai titik masuk aplikasi. Di sinilah request dari client diterima dan response dikembalikan.

**2. Usecase Layer (service)**
Berisi logika bisnis utama aplikasi. Layer ini mengatur alur proses dan menggabungkan berbagai komponen agar bekerja sama.

**3. Repository Layer (repository)**
Bertanggung jawab atas akses ke data (misalnya: database, cache). Layer ini menjadi jembatan antara usecase dan sumber data aktual.

**4. Entity/Model Layer**
Mendefinisikan struktur data dan aturan bisnis dasar (entity) yang digunakan di seluruh lapisan.

**5. Config & Middleware**
Mengelola konfigurasi environment, koneksi database, Redis, serta middleware seperti authentication, logging, dan recovery.

**Why Clean Architecture?**
Beberapa alasan memilih Clean Architecture:

**1. Separation of Concerns**
Memisahkan logika bisnis dari teknologi (seperti database dan framework), sehingga lebih mudah untuk mengembangkan dan menguji bagian-bagian aplikasi secara terpisah.

**2. Maintainability**
Struktur yang jelas memudahkan dalam menambahkan fitur baru atau memperbaiki bug tanpa mengganggu bagian lain dari aplikasi.

**3. Scalability**
Aplikasi lebih siap untuk tumbuh dan menangani kompleksitas lebih besar seiring waktu, karena arsitektur sudah modular dari awal.

**4. Framework-agnostic**
Arsitektur ini tidak bergantung pada framework tertentu seperti Gin, jadi jika ingin mengganti framework, tidak perlu mengubah logika bisnis utama.

## How To Use
1. Clone the repository
  ```bash
  git clone https://github.com/kevinlau2102/erajaya-test-backend.git
  ```
2. Navigate to the project directory:
  ```bash
  cd erajaya-test-backend
  ```
3. Copy the example environment file and configure it:
  ```bash 
  cp .env.example .env
  ```
There are 2 ways to do running
### With Docker
1. Build Docker
  ```bash
  make up
  ```
2. Run Initial UUID V4 for Auto Generate UUID
  ```bash
  make init-uuid
  ```
3. Run Migration and Seeder
  ```bash
  make migrate-seed
  ```

## List API

You can try the API using this Postman collection: https://api.postman.com/collections/25822863-2feafd6d-19ff-40ce-a9d2-85d9024a6015?access_key=PMAT-01K0DD79WE6XQNMCJ1V119SSSP

### 1. **POST** `/api/v1/user`

Register a new user.

**Request Body:**
```json
{
    "name": "Kevin Laurence",
    "telp_number": "08123456789",
    "email": "kevinlaurence@gmail.com",
    "password": "12345678"
}
```

**Response:**
```json
{
    "status": 200,
    "message": "success create user",
    "data": {
        "id": "56c281ca-abf4-40d2-8149-f58662b5d9fa",
        "name": "Kevin Laurence",
        "email": "kevinlaurence@gmail.com",
        "telp_number": "08123456789",
        "role": "user"
    }
}
```

### 2. **POST** `/api/v1/user/login`

Login to get access token and refresh token.

**Request Body:**
```json
{
    "email": "kevinlaurence@gmail.com",
    "password": "12345678"
}
```

**Response:**
```json
{
    "status": 200,
    "message": "success login",
    "data": {
        "access_token": "...",
        "refresh_token": "...",
        "role": "user"
    }
}
```

### 3. **POST** `/api/v1/user/refresh`

API Refresh token to get a new access token when expired.

**Request Body:**
```json
{
    "refresh_token": "..."
}
```

**Response:**
```json
{
    "status": 200,
    "message": "successfully refreshed token",
    "data": {
        "access_token": "...",
        "refresh_token": "...",
        "role": "user"
    }
}
```

### 4. **POST** `/api/v1/product`

Create a new product.

**Headers:**
```
Authorization: Bearer <your_token_here>
Content-Type: application/json
```

**Request Body:**
```json
{
    "name": "Wireless Mouse",
    "description": "Ergonomic wireless mouse with 3 adjustable DPI levels and silent click.",
    "price": 149000,
    "quantity": 25
}
```

**Response:**
```json
{
    "status": 200,
    "message": "success create product",
    "data": {
        "id": "a3fffe22-21bb-47de-a87e-3371ad12606a",
        "name": "Wireless Mouse",
        "description": "Ergonomic wireless mouse with 3 adjustable DPI levels and silent click.",
        "price": 149000,
        "quantity": 25
    }
}
```

### 5. **GET** `/api/v1/product`

Get all products with pagination, sorting, and search.

**Headers:**
```
Authorization: Bearer <your_token_here>
Content-Type: application/json
```

**Query Parameters:**
- `page` (int, optional) - default: 1
- `limit` (int, optional) - default: 10
- `search` (string, optional)
- `sortBy` (string: `name`, `price`, `created_at`)
- `order` (string: `asc`, `desc`)

**Response:**
```json
{
    "status": 200,
    "message": "success get list product",
    "data": [...],
    "meta": {
        "page": 1,
        "limit": 10,
        "max_page": 1,
        "count": 3
    }
}
```
