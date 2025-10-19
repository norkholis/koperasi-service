# Koperasi Service API Documentation

Complete API documentation for the Cooperative Management System with role-based access control.

## Base URL
```
http://localhost:8080/api
```

## CORS Support
The API includes CORS (Cross-Origin Resource Sharing) support for frontend applications. All origins are currently allowed for development purposes.

## Authentication

All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer {your_jwt_token}
```

## Roles
- **super_admin**: Full access to all resources
- **admin**: Access to own data and user management
- **member**: Access to own data only

---

## Authentication Endpoints

### Register User
```http
POST /api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "role_id": 3
}
```

**Response:**
```json
{
  "message": "User registered"
}
```

### Login
```http
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Get Current User Info
```http
GET /api/me
Authorization: Bearer {token}
```

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "role": {
    "id": 3,
    "name": "member"
  }
}
```

---

## User Management (Admin/Super Admin Only)

### List All Users (Super Admin Only)
```http
GET /api/users
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "email": "admin@example.com",
      "role_id": 2
    }
  ]
}
```

### Get User Detail (Super Admin or Self)
```http
GET /api/users/{id}
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "email": "admin@example.com",
    "role_id": 2
  }
}
```

### Create User (Super Admin Only)
```http
POST /api/users
Authorization: Bearer {token}
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password123",
  "role_id": 3
}
```

### Update User (Super Admin or Self)
```http
PUT /api/users/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "email": "updated@example.com",
  "password": "newpassword123",
  "role_id": 2
}
```

**Note:** Only super_admin can change role_id

### Delete User (Super Admin or Self)
```http
DELETE /api/users/{id}
Authorization: Bearer {token}
```

---

## Simpanan (Savings) Management

### Create Simpanan
```http
POST /api/simpanan
Authorization: Bearer {token}
Content-Type: application/json

{
  "type": "wajib",
  "amount": 100000,
  "description": "Simpanan wajib bulan Januari"
}
```

**Types:** `wajib` or `sukarela`

### List Simpanan
```http
GET /api/simpanan
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": [
    {
      "ID": 1,
      "CreatedAt": "2024-01-15T10:00:00Z",
      "UpdatedAt": "2024-01-15T10:00:00Z",
      "UserID": 1,
      "Type": "wajib",
      "Amount": 100000,
      "Description": "Simpanan wajib bulan Januari"
    }
  ]
}
```

### Get Simpanan Detail
```http
GET /api/simpanan/{id}
Authorization: Bearer {token}
```

### Update Simpanan
```http
PUT /api/simpanan/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "type": "sukarela",
  "amount": 150000,
  "description": "Updated description"
}
```

### Delete Simpanan
```http
DELETE /api/simpanan/{id}
Authorization: Bearer {token}
```

---

## Pinjaman (Loan) Management

### Create Pinjaman
```http
POST /api/pinjaman
Authorization: Bearer {token}
Content-Type: application/json

{
  "jumlah_pinjaman": 5000000,
  "bunga_persen": 2.5,
  "lama_bulan": 12,
  "jumlah_angsuran": 450000,
  "user_id": 1
}
```

**Optional fields:**
- `kode_pinjaman`: Auto-generated if not provided
- `status`: Defaults to "proses"

### List Pinjaman
```http
GET /api/pinjaman
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": [
    {
      "ID": 1,
      "CreatedAt": "2024-01-15T10:00:00Z",
      "KodePinjaman": "PJM1705123456",
      "UserID": 1,
      "TanggalPinjam": "2024-01-15T10:00:00Z",
      "JumlahPinjaman": 5000000,
      "BungaPersen": 2.5,
      "LamaBulan": 12,
      "JumlahAngsuran": 450000,
      "SisaAngsuran": 12,
      "Status": "proses",
      "User": {
        "ID": 1,
        "Email": "user@example.com"
      }
    }
  ]
}
```

### Get Pinjaman Detail
```http
GET /api/pinjaman/{id}
Authorization: Bearer {token}
```

### Update Pinjaman
```http
PUT /api/pinjaman/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": "disetujui",
  "sisa_angsuran": 11
}
```

**Valid statuses:** `proses`, `disetujui`, `lunas`, `macet`

### Delete Pinjaman
```http
DELETE /api/pinjaman/{id}
Authorization: Bearer {token}
```

---

## Angsuran (Installment) Management

### Create Angsuran Payment
```http
POST /api/angsuran
Authorization: Bearer {token}
Content-Type: application/json

{
  "pinjaman_id": 1,
  "angsuran_ke": 1,
  "pokok": 400000,
  "bunga": 50000,
  "denda": 0
}
```

**Optional fields:**
- `total_bayar`: Auto-calculated if not provided
- `user_id`: Defaults to loan owner

### List Angsuran
```http
GET /api/angsuran
Authorization: Bearer {token}

# Optional filter by loan
GET /api/angsuran?pinjaman_id=1
```

**Response:**
```json
{
  "data": [
    {
      "ID": 1,
      "PinjamanID": 1,
      "AngsuranKe": 1,
      "TanggalBayar": "2024-01-15T10:00:00Z",
      "Pokok": 400000,
      "Bunga": 50000,
      "Denda": 0,
      "TotalBayar": 450000,
      "UserID": 1,
      "Status": "proses",
      "Pinjaman": { /* loan details */ },
      "User": { /* user details */ }
    }
  ]
}
```

### Get Angsuran Detail
```http
GET /api/angsuran/{id}
Authorization: Bearer {token}
```

### Update Angsuran
```http
PUT /api/angsuran/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "pokok": 420000,
  "bunga": 52000,
  "denda": 10000,
  "status": "proses"
}
```

### Verify Payment (Admin/Super Admin Only)
```http
PUT /api/angsuran/{id}/verify
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": "verified"
}
```

**Valid verification statuses:** `verified`, `kurang`, `lebih`

**Note:** When status is `verified`, the system automatically:
- Decrements `sisa_angsuran` in the related pinjaman
- Changes pinjaman status to `lunas` when all installments are paid

### Get Pending Payments (Admin/Super Admin Only)
```http
GET /api/angsuran/pending
Authorization: Bearer {token}
```

### Delete Angsuran
```http
DELETE /api/angsuran/{id}
Authorization: Bearer {token}
```

---

## SHU (Annual Profit Sharing) Management (Admin/Super Admin Only)

### Generate SHU Report
```http
POST /api/shu/generate
Authorization: Bearer {token}
Content-Type: application/json

{
  "tahun": 2024,
  "total_shu_koperasi": 40000000
}
```

**Response:**
```json
{
  "message": "SHU report generated successfully",
  "data": {
    "tahun": 2024,
    "total_shu_koperasi": 40000000,
    "persen_jasa_modal": 25,
    "persen_jasa_usaha": 30,
    "total_simpanan_all": 60000000,
    "total_penjualan_all": 100000000,
    "tanggal_hitung": "2024-01-15T10:00:00Z",
    "detail_anggota": [
      {
        "user_id": 1,
        "email": "user@example.com",
        "total_simpanan": 3000000,
        "total_penjualan": 1000000,
        "jasa_modal": 500000,
        "jasa_usaha": 120000,
        "total_shu_anggota": 620000
      }
    ]
  }
}
```

### Save SHU Record
```http
POST /api/shu
Authorization: Bearer {token}
Content-Type: application/json

{
  "tahun": 2024,
  "total_shu": 40000000,
  "status": "draft"
}
```

**Valid statuses:** `draft`, `final`

### List SHU Records
```http
GET /api/shu
Authorization: Bearer {token}
```

### Get SHU Detail
```http
GET /api/shu/{id}
Authorization: Bearer {token}
```

### Get SHU by Year
```http
GET /api/shu/year/{tahun}
Authorization: Bearer {token}

# Example
GET /api/shu/year/2024
```

### Update SHU Record
```http
PUT /api/shu/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "total_shu": 42000000,
  "status": "final"
}
```

### Delete SHU Record
```http
DELETE /api/shu/{id}
Authorization: Bearer {token}
```

---

## SHU Calculation Formulas

The system implements the exact cooperative SHU calculation formulas:

### Jasa Modal Anggota (JMA)
```
JMA = (Simpanan anggota / Total simpanan koperasi) × 25% × Total SHU Koperasi
```

### Jasa Usaha Anggota (JUA)  
```
JUA = (Penjualan anggota / Total penjualan koperasi) × 30% × Total SHU Koperasi
```

### Total SHU Anggota
```
SHU Anggota = JMA + JUA
```

**Example Calculation:**
- SHU Koperasi: Rp40,000,000
- Total Simpanan: Rp60,000,000  
- Total Penjualan: Rp100,000,000
- Simpanan Anggota: Rp3,000,000
- Penjualan Anggota: Rp1,000,000

**Result:**
- JMA: (3,000,000 / 60,000,000) × 25% × 40,000,000 = Rp500,000
- JUA: (1,000,000 / 100,000,000) × 30% × 40,000,000 = Rp120,000
- **Total SHU Anggota: Rp620,000**

---

## Error Responses

### Common Error Codes

**400 Bad Request:**
```json
{
  "error": "Invalid input data"
}
```

**401 Unauthorized:**
```json
{
  "error": "Invalid token"
}
```

**403 Forbidden:**
```json
{
  "error": "forbidden"
}
```

**404 Not Found:**
```json
{
  "error": "not found"
}
```

**409 Conflict:**
```json
{
  "error": "SHU for this year already exists"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error"
}
```

---

## Access Control Summary

| Endpoint | Member | Admin | Super Admin |
|----------|--------|-------|-------------|
| Authentication | ✅ | ✅ | ✅ |
| User Management | Self only | Self only | All users |
| Simpanan | Own only | Own only | All records |
| Pinjaman | Own only | Own only | All records |
| Angsuran | Own only | Own only | All records |
| Angsuran Verify | ❌ | Own users | All records |
| SHU Management | ❌ | ✅ | ✅ |

---

## Setup and Running

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set environment variables:**
   ```bash
   export DB_HOST=localhost
   export DB_USER=postgres
   export DB_PASS=password
   export DB_NAME=koperasi
   export DB_PORT=5432
   export JWT_SECRET=your-secret-key
   ```

3. **Run the server:**
   ```bash
   go run ./cmd/main.go
   ```

4. **Server starts on port 8080**

The system automatically creates database tables and seeds default roles:
- `super_admin` (ID: 1)
- `admin` (ID: 2)  
- `member` (ID: 3)

---

## Testing Examples

### Complete Workflow Example

1. **Register and Login:**
   ```bash
   # Register
   curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123","role_id":3}'
   
   # Login
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. **Create Savings:**
   ```bash
   curl -X POST http://localhost:8080/api/simpanan \
     -H "Authorization: Bearer {token}" \
     -H "Content-Type: application/json" \
     -d '{"type":"wajib","amount":1000000,"description":"Initial savings"}'
   ```

3. **Apply for Loan:**
   ```bash
   curl -X POST http://localhost:8080/api/pinjaman \
     -H "Authorization: Bearer {token}" \
     -H "Content-Type: application/json" \
     -d '{"jumlah_pinjaman":5000000,"bunga_persen":2.5,"lama_bulan":12,"jumlah_angsuran":450000}'
   ```

4. **Make Payment:**
   ```bash
   curl -X POST http://localhost:8080/api/angsuran \
     -H "Authorization: Bearer {token}" \
     -H "Content-Type: application/json" \
     -d '{"pinjaman_id":1,"angsuran_ke":1,"pokok":400000,"bunga":50000}'
   ```

5. **Generate SHU (Admin only):**
   ```bash
   curl -X POST http://localhost:8080/api/shu/generate \
     -H "Authorization: Bearer {admin_token}" \
     -H "Content-Type: application/json" \
     -d '{"tahun":2024,"total_shu_koperasi":40000000}'
   ```