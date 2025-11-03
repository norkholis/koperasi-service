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

## Roles & Permissions
- **super_admin**: Full unrestricted access to all resources and users
- **admin**: Access to own data + users they registered (admin user hierarchy)
- **member**: Access to own data only

### Admin User Hierarchy
- When **admin** users create new users, those users are assigned to them via `admin_id`
- **admin** users can only access/manage users where `admin_id` matches their own ID
- **super_admin** users have no restrictions and can access all users regardless of `admin_id`
- **member** users can only access their own profile

---

## Authentication Endpoints

### Register User
```http
POST /api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "address": "123 Main Street, City",
  "phone_number": "081234567890",
  "nik": "1234567890123456",
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
  "name": "John Doe",
  "address": "123 Main Street, City",
  "phone_number": "081234567890",
  "nik": "1234567890123456",
  "role": {
    "id": 3,
    "name": "member"
  },
  "admin_id": 1
}
```

### Change Password
```http
POST /api/change-password
Authorization: Bearer {token}
Content-Type: application/json

{
  "current_password": "oldpassword123",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Response:**
```json
{
  "message": "Password changed successfully"
}
```

**Error Response (400):**
```json
{
  "error": "Current password is incorrect"
}
```

### Forgot Password
```http
POST /api/forgot-password
Content-Type: application/json

{
  "email": "user@example.com",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Response:**
```json
{
  "message": "Password reset successfully"
}
```

**Error Response (400):**
```json
{
  "error": "Password confirmation does not match"
}
```

---

## User Management

### List All Users
```http
GET /api/users
Authorization: Bearer {token}
```

**Access Control:**
- **super_admin**: Gets all users
- **admin**: Gets only users they registered (where `admin_id` matches their ID)
- **member**: Forbidden

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "email": "admin@example.com",
      "name": "Admin User",
      "address": "Admin Address",
      "phone_number": "081234567890",
      "nik": "1234567890123456",
      "role_id": 2,
      "admin_id": null
    },
    {
      "id": 2,
      "email": "member@example.com",
      "name": "Member User",
      "address": "Member Address",
      "phone_number": "087654321098",
      "nik": "6543210987654321",
      "role_id": 3,
      "admin_id": 1
    }
  ]
}
```

### Get User Detail
```http
GET /api/users/{id}
Authorization: Bearer {token}
```

**Access Control:**
- **super_admin**: Can get any user
- **admin**: Can get users they registered OR themselves
- **member**: Can only get themselves

**Response:**
```json
{
  "data": {
    "id": 2,
    "email": "member@example.com",
    "name": "Member User",
    "address": "Member Address",
    "phone_number": "087654321098",
    "nik": "6543210987654321",
    "role_id": 3,
    "admin_id": 1
  }
}
```

### Create User
```http
POST /api/users
Authorization: Bearer {token}
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "New User",
  "address": "New User Address",
  "phone_number": "089876543210",
  "nik": "9876543210987654",
  "role_id": 3
}
```

**Access Control:**
- **super_admin**: Can create any user (no `admin_id` assigned)
- **admin**: Can create users (automatically assigns their ID as `admin_id`)
- **member**: Forbidden

**Response:**
```json
{
  "data": {
    "id": 3,
    "email": "newuser@example.com",
    "name": "New User",
    "address": "New User Address",
    "phone_number": "089876543210",
    "nik": "9876543210987654",
    "role_id": 3,
    "admin_id": 1
  }
}
```

### Update User
```http
PUT /api/users/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "email": "updated@example.com",
  "name": "Updated Name",
  "address": "Updated Address",
  "phone_number": "081111111111",
  "nik": "1111111111111111",
  "password": "newpassword123",
  "role_id": 2
}
```

**Access Control:**
- **super_admin**: Can update any user + change roles
- **admin**: Can update users they registered OR themselves (cannot change roles)
- **member**: Can only update themselves (cannot change roles)

**Note:** Only super_admin can change role_id

### Delete User
```http
DELETE /api/users/{id}
Authorization: Bearer {token}
```

**Access Control:**
- **super_admin**: Can delete any user
- **admin**: Can delete users they registered OR themselves
- **member**: Can only delete themselves

---

## Simpanan (Wallet) Management

The new simpanan system works as a **wallet** with 3 automatic wallet types for each user:
- **pokok**: Basic capital savings
- **wajib**: Mandatory savings  
- **sukarela**: Voluntary savings

### Business Model:
1. **Automatic Wallet Creation**: Each user automatically gets 3 wallet types when registered
2. **User Top-up Flow**: Users can request top-ups → Creates pending transactions → Admin verifies → Balance updated
3. **Admin Management**: Admin/Super Admin can verify top-up requests and directly adjust balances
4. **Transaction Tracking**: All transactions are tracked with approval workflow and history
5. **Role-based Access**: Members see own wallets, Admins manage all wallets

### Key Features:
- **Wallet Types**: `pokok`, `wajib`, `sukarela` (automatically created for each user)
- **Top-up Requests**: User-initiated, admin-verified
- **Balance Adjustments**: Admin-only direct balance modifications
- **Transaction History**: Complete audit trail for all wallet activities
- **Verification Workflow**: Pending → Verified/Rejected status for top-ups

---

### Get User Wallets
```http
GET /api/simpanan/wallets
Authorization: Bearer {token}
```

**Query Parameters:**
- `user_id` (optional): For admin to view specific user's wallets

**Access Control:**
- **Member**: Can only see their own wallets
- **Admin/Super Admin**: Can see any user's wallets with `user_id` parameter

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "type": "pokok",
      "balance": 500000,
      "description": "Wallet pokok",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "type": "wajib",
      "balance": 1200000,
      "description": "Wallet wajib",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": 3,
      "user_id": 1,
      "type": "sukarela",
      "balance": 800000,
      "description": "Wallet sukarela",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

### Get All Wallets (Admin Only)
```http
GET /api/simpanan/wallets/all
Authorization: Bearer {token}
```

**Access Control:** Admin/Super Admin only

### Top-up Wallet
```http
POST /api/simpanan/topup
Authorization: Bearer {token}
Content-Type: application/json

{
  "type": "wajib",
  "amount": 200000,
  "description": "Monthly mandatory savings"
}
```

**Valid Types:** `pokok`, `wajib`, `sukarela`

**Response:**
```json
{
  "message": "Top-up request created, waiting for admin verification"
}
```

**Note:** Creates a pending transaction that requires admin verification.

### Get Wallet Detail
```http
GET /api/simpanan/{wallet_id}
Authorization: Bearer {token}
```

**Access Control:**
- **Member**: Can only view their own wallets
- **Admin/Super Admin**: Can view any wallet

**Response:**
```json
{
  "data": {
    "id": 1,
    "user_id": 1,
    "type": "pokok",
    "balance": 500000,
    "description": "Wallet pokok",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

### Get Wallet Transaction History
```http
GET /api/simpanan/{wallet_id}/transactions
Authorization: Bearer {token}
```

**Access Control:**
- **Member**: Can only view their own wallet transactions
- **Admin/Super Admin**: Can view any wallet transactions

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "simpanan_id": 1,
      "type": "topup",
      "amount": 200000,
      "description": "Monthly mandatory savings",
      "status": "verified",
      "verified_by_id": 2,
      "verified_at": "2024-01-15T11:00:00Z",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Verify Transaction (Admin Only)
```http
PUT /api/simpanan/transactions/{transaction_id}/verify
Authorization: Bearer {token}
Content-Type: application/json

{
  "approve": true
}
```

**Access Control:** Admin/Super Admin only

**Response:**
```json
{
  "message": "Transaction approved"
}
```

**Note:** When approved, the wallet balance is automatically updated.

### Adjust Wallet Balance (Admin Only)
```http
PUT /api/simpanan/{wallet_id}/adjust
Authorization: Bearer {token}
Content-Type: application/json

{
  "amount": -50000,
  "description": "Administrative fee deduction"
}
```

**Access Control:** Admin/Super Admin only

**Note:** 
- Positive amounts increase balance
- Negative amounts decrease balance
- Creates verified transaction immediately

### Get Pending Transactions (Admin Only)
```http
GET /api/simpanan/transactions/pending
Authorization: Bearer {token}
```

**Access Control:** Admin/Super Admin only

**Response:**
```json
{
  "data": [
    {
      "id": 5,
      "simpanan_id": 2,
      "simpanan": {
        "id": 2,
        "user_id": 3,
        "type": "wajib"
      },
      "type": "topup",
      "amount": 150000,
      "description": "Monthly savings",
      "status": "pending",
      "created_at": "2024-01-15T14:30:00Z"
    }
  ]
}
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

**Role-based Access:**
- **Regular Users**: Can only see their own loans
- **Admin Users**: Can see loans from users they registered (admin hierarchy)
- **Super Admin**: Can see all loans in the system

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

**Role-based Access:**
- **Regular Users**: Can only access their own loan details
- **Admin Users**: Can access loan details for users they registered
- **Super Admin**: Can access any loan details

**Response:** Same structure as List Pinjaman but for a single loan.

### Update Pinjaman
```http
PUT /api/pinjaman/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": "disetujui",
  "jumlah_pinjaman": 5000000,
  "bunga_persen": 2.5,
  "lama_bulan": 12,
  "jumlah_angsuran": 450000
}
```

**Important Notes:**
- `sisa_angsuran` cannot be directly updated via API
- `sisa_angsuran` is automatically managed by the system:
  - Initialized to `lama_bulan` when loan is created
  - Decremented by 1 when installment payments are verified
  - When reaches 0, loan status automatically becomes "lunas"
- `bunga_persen` is only updated when explicitly provided with value > 0
- Fields with 0 values are ignored to prevent accidental resets

**Example - Approve loan without changing other fields:**
```json
{
  "status": "disetujui"
}
```

**Role-based Access:**
- **Regular Users**: Can only update their own loan details
- **Admin Users**: Can update loans for users they registered
- **Super Admin**: Can update any loan

**Valid statuses:** `proses`, `disetujui`, `lunas`, `macet`

### Delete Pinjaman
```http
DELETE /api/pinjaman/{id}
Authorization: Bearer {token}
```

**Role-based Access:**
- **Regular Users**: Can only delete their own loans
- **Admin Users**: Can delete loans for users they registered
- **Super Admin**: Can delete any loan

---

## Loan Workflow & Business Logic

### Loan Lifecycle
1. **Application**: User creates loan with status "proses" and `sisa_angsuran = lama_bulan`
2. **Approval**: Admin changes status to "disetujui" (approved) - `sisa_angsuran` remains unchanged
3. **Payments**: User makes installment payments (angsuran) with status "proses"  
4. **Verification**: Admin verifies payments - `sisa_angsuran` decrements by 1 for each verified payment
5. **Completion**: When `sisa_angsuran = 0`, loan status automatically becomes "lunas" (paid off)

### Important Rules
- `sisa_angsuran` is **system-managed** and cannot be directly updated via API
- Only verified installment payments can reduce `sisa_angsuran`
- Loan approval does NOT affect the remaining installment count
- Each verified payment reduces `sisa_angsuran` by exactly 1

---

## Angsuran (Installment) Management

### Create Angsuran Payment
```http
POST /api/angsuran
Authorization: Bearer {token}
Content-Type: application/json

{
  "pinjaman_id": 1,
  "pokok": 400000,
  "bunga": 50000,
  "denda": 0
}
```

**Auto-generated fields:**
- `angsuran_ke`: Automatically incremented based on existing payments for the loan
- `total_bayar`: Auto-calculated if not provided (pokok + bunga + denda)

**Optional fields:**
- `angsuran_ke`: Can be manually specified if needed (otherwise auto-generated)
- `user_id`: Defaults to loan owner
- `denda`: Defaults to 0

### List Angsuran
```http
GET /api/angsuran
Authorization: Bearer {token}

# Optional filter by loan
GET /api/angsuran?pinjaman_id=1
```

**Role-based Access:**
- **Regular Users**: Can only see their own installment payments
- **Admin Users**: Can see installment payments from users they registered (admin hierarchy)
- **Super Admin**: Can see all installment payments in the system

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

## SHU (Annual Profit Sharing) Management

The SHU (Sisa Hasil Usaha) system provides three types of report generation:

1. **Full SHU Report** (`/api/shu/generate`) - **Admin/Super Admin Only**: Manual SHU calculation - admin inputs total SHU amount
2. **Automated SHU Report** (`/api/shu/generate-auto`) - **Admin/Super Admin Only**: Automated SHU calculation - system calculates total SHU from income and expenses  
3. **Individual User SHU** (`/api/shu/user/{user_id}/generate`) - **All Users**: Generates SHU calculation for a specific user using existing SHU records

**Access Control:**
- **Admin/Super Admin**: Can access all SHU management features and view any user's SHU data
- **Regular Users**: Can only view their own SHU data using their own user ID

**Workflow:**
1. Admin first generates a SHU report using either manual (`/generate`) or automated (`/generate-auto`) method
2. System saves the SHU record with total SHU amount (status: "draft")
3. Admin can review and update the SHU record status from "draft" to "final" using `/api/shu/{id}` PUT endpoint
4. Users can then generate their own SHU calculations which automatically use the saved total SHU amount

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
    "total_shu_koperasi": 97141305,
    "persen_jasa_modal": 30,
    "persen_jasa_usaha": 70,
    "total_simpanan_all": 390350000,
    "total_penjualan_all": 842458574,
    "tanggal_hitung": "2024-01-15T10:00:00Z",
    "detail_anggota": [
      {
        "user_id": 1,
        "email": "user@example.com",
        "total_simpanan": 700000,
        "total_penjualan": 1700000,
        "jasa_modal": 26098,
        "jasa_usaha": 68606,
        "total_shu_anggota": 94704
      }
    ]
  }
}
```

### Generate User SHU Report
```http
POST /api/shu/user/{user_id}/generate
Authorization: Bearer {token}
Content-Type: application/json

{
  "tahun": 2024
}
```

**Description:** Generate SHU calculation for a specific user. The `total_shu_koperasi` value is automatically retrieved from existing SHU records for the specified year. 

**Access Control:**
- **Users** can only access their own SHU data (user_id must match their own ID)
- **Admin/Super Admin** can access any user's SHU data
- Admin/Super Admin must generate the SHU report first using `/api/shu/generate` or `/api/shu/generate-auto` before using this endpoint

**Error Response (if SHU record not found):**
```json
{
  "error": "SHU record not found for the specified year. Please generate SHU report first"
}
```

**Error Response (if user tries to access other user's data):**
```json
{
  "error": "forbidden"
}
```

**Response:**
```json
{
  "message": "User SHU calculated successfully",
  "data": {
    "user_id": 17,
    "email": "user7@gmail.com",
    "total_simpanan": 700000,
    "total_penjualan": 1700000,
    "jasa_modal": 26098,
    "jasa_usaha": 68606,
    "total_shu_anggota": 94704
  }
}
```

### Generate Automated SHU Report
```http
POST /api/shu/generate-auto
Authorization: Bearer {token}
Content-Type: application/json

{
  "tahun": 2024,
  "beban_operasional": 15000000,
  "beban_non_operasional": 5000000,
  "beban_pajak": 2000000
}
```

**Description:** Automated SHU calculation where the system calculates income automatically from loan interest and other sources. Admin only needs to input the expenses.

**Formula Used:** 
`SHU Total = (Pendapatan Operasional + Pendapatan Non-Operasional) - (Beban Operasional + Beban Non-Operasional + Beban Pajak)`

**Response:**
```json
{
  "message": "Automated SHU report generated successfully",
  "data": {
    "tahun": 2024,
    "pendapatan_operasional": 35000000,
    "pendapatan_non_operasional": 7000000,
    "beban_operasional": 15000000,
    "beban_non_operasional": 5000000,
    "beban_pajak": 2000000,
    "total_shu_koperasi": 20000000,
    "persen_jasa_modal": 30,
    "persen_jasa_usaha": 70,
    "total_simpanan_all": 390350000,
    "total_penjualan_all": 842458574,
    "tanggal_hitung": "2024-01-15T10:00:00Z",
    "detail_anggota": [
      {
        "user_id": 1,
        "email": "user@example.com",
        "total_simpanan": 700000,
        "total_penjualan": 1700000,
        "jasa_modal": 26754,
        "jasa_usaha": 70819,
        "total_shu_anggota": 97573
      }
    ]
  }
}
```

### Save Automated SHU Record
```http
POST /api/shu/save-auto
Authorization: Bearer {token}
Content-Type: application/json

{
  "tahun": 2024,
  "pendapatan_operasional": 35000000,
  "pendapatan_non_operasional": 7000000,
  "beban_operasional": 15000000,
  "beban_non_operasional": 5000000,
  "beban_pajak": 2000000,
  "total_shu": 20000000,
  "status": "draft"
}
```

**Description:** Save the automated SHU calculation with detailed income and expense breakdown.

**Valid statuses:** `draft`, `final`

**Response:**
```json
{
  "message": "Automated SHU record saved successfully",
  "data": {
    "id": 1,
    "tahun": 2024,
    "pendapatan_operasional": 35000000,
    "pendapatan_non_operasional": 7000000,
    "beban_operasional": 15000000,
    "beban_non_operasional": 5000000,
    "beban_pajak": 2000000,
    "total_shu": 20000000,
    "tanggal_hitung": "2024-01-15T10:00:00Z",
    "status": "draft"
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

**Description:** Update an existing SHU record. Common use case is updating status from "draft" to "final" once the SHU calculation is approved.

**Valid statuses:** `draft`, `final`

**Response:**
```json
{
  "data": {
    "id": 1,
    "tahun": 2024,
    "total_shu": 42000000,
    "status": "final",
    "tanggal_hitung": "2024-01-15T10:00:00Z"
  }
}
```

### Delete SHU Record
```http
DELETE /api/shu/{id}
Authorization: Bearer {token}
```

---

## SHU Calculation Formulas

The system implements two types of SHU calculation methods:

### 1. Manual SHU Calculation
Admin inputs the total SHU amount directly, and the system distributes it among members.

### 2. Automated SHU Calculation  
System calculates the total SHU automatically based on income and expenses:

```
SHU Total = (Pendapatan Operasional + Pendapatan Non-Operasional) - (Beban Operasional + Beban Non-Operasional + Beban Pajak)
```

**Components:**
- **Pendapatan Operasional**: Automatically calculated from loan interest and other operational income
- **Pendapatan Non-Operasional**: Automatically calculated from investments and other non-operational income  
- **Beban Operasional**: Input by admin (operational expenses)
- **Beban Non-Operasional**: Input by admin (non-operational expenses)
- **Beban Pajak**: Input by admin (tax expenses)

### Member Distribution Formulas (Both Methods)

**Step 1: Calculate SHU for Members**
```
SHU untuk Anggota = 50% × Total SHU Koperasi
```

**Step 2: Allocate Member SHU**
```
Alokasi Jasa Modal = 30% × SHU untuk Anggota
Alokasi Jasa Usaha = 70% × SHU untuk Anggota
```

**Step 3: Calculate Individual Member SHU**

### Jasa Modal Anggota (JMA)
```
JMA = (Simpanan anggota / Total simpanan koperasi) × Alokasi Jasa Modal
```

### Jasa Usaha Anggota (JUA)  
```
JUA = (Pinjaman anggota / Total pinjaman koperasi) × Alokasi Jasa Usaha
```

### Total SHU Anggota
```
SHU Anggota = JMA + JUA
```

**Example Calculation (Pak Abdul):**
- SHU Total: Rp97,141,305
- SHU untuk Anggota: 50% × Rp97,141,305 = Rp48,570,652
- Alokasi Jasa Modal: 30% × Rp48,570,652 = Rp14,571,196
- Alokasi Jasa Usaha: 70% × Rp48,570,652 = Rp33,999,456
- Total Simpanan: Rp390,350,000
- Total Pinjaman: Rp842,458,574
- Simpanan Pak Abdul: Rp700,000
- Pinjaman Pak Abdul: Rp1,700,000

**Result:**
- JMA: (700,000 / 390,350,000) × 14,571,196 = Rp26,098
- JUA: (1,700,000 / 842,458,574) × 33,999,456 = Rp68,606
- **Total SHU Anggota: Rp94,704**

---

## Complete SHU Management Workflow Example

**Scenario:** Admin needs to calculate and finalize SHU for year 2024

**Step 1:** Generate SHU Report (Automated)
```http
POST /api/shu/generate-auto
{
  "tahun": 2024,
  "beban_operasional": 15000000,
  "beban_non_operasional": 5000000,
  "beban_pajak": 2000000
}
```
*System calculates total SHU = Rp20,000,000*

**Step 2:** Save SHU Record
```http
POST /api/shu/save-auto
{
  "tahun": 2024,
  "pendapatan_operasional": 35000000,
  "pendapatan_non_operasional": 7000000,
  "beban_operasional": 15000000,
  "beban_non_operasional": 5000000,
  "beban_pajak": 2000000,
  "total_shu": 20000000,
  "status": "draft"
}
```
*SHU record saved with status "draft"*

**Step 3:** Review and Finalize
```http
PUT /api/shu/1
{
  "status": "final"
}
```
*Status updated from "draft" to "final"*

**Step 4:** Generate Individual Member SHU

*Admin can access any user's SHU:*
```http
POST /api/shu/user/17/generate
{
  "tahun": 2024
}
```

*Regular user can access their own SHU (if their ID is 17):*
```http
POST /api/shu/user/17/generate
{
  "tahun": 2024
}
```
*Uses saved total SHU (Rp20,000,000) automatically*

**Step 5:** Save Individual Member SHU (Optional)

*Users can save their calculated SHU:*
```http
POST /api/shu-anggota/user/17/save
{
  "tahun": 2024
}
```
*Saves the calculated SHU data to shu_anggota table*

---

## SHU Anggota (Individual Member SHU Records) Management

### Save User SHU
```http
POST /api/shu-anggota/user/{user_id}/save
Authorization: Bearer {token}
Content-Type: application/json

{
  "tahun": 2024
}
```

**Description:** Save the calculated SHU for a specific user to the database. This calculates and permanently stores the user's SHU breakdown.

**Access Control:**
- **Users** can save their own SHU data (user_id must match their own ID)
- **Admin/Super Admin** can save any user's SHU data

**Response:**
```json
{
  "message": "User SHU saved successfully",
  "data": {
    "id_shu_anggota": 1,
    "id_shu": 5,
    "id_anggota": 17,
    "jumlah_modal": 26098,
    "jumlah_usaha": 68606,
    "shu_diterima": 94704,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

### Get User SHU by Year
```http
GET /api/shu-anggota/user/{user_id}/{tahun}
Authorization: Bearer {token}

# Example
GET /api/shu-anggota/user/17/2024
```

**Description:** Retrieve saved SHU data for a specific user and year.

**Access Control:**
- **Users** can access their own saved SHU data
- **Admin/Super Admin** can access any user's saved SHU data

**Response:**
```json
{
  "data": {
    "id_shu_anggota": 1,
    "id_shu": 5,
    "id_anggota": 17,
    "jumlah_modal": 26098,
    "jumlah_usaha": 68606,
    "shu_diterima": 94704,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "shu": {
      "id_shu": 5,
      "tahun": 2024,
      "total_shu": 20000000,
      "status": "final"
    },
    "user": {
      "id": 17,
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

### Get User SHU History
```http
GET /api/shu-anggota/user/{user_id}/history
Authorization: Bearer {token}

# Example
GET /api/shu-anggota/user/17/history
```

**Description:** Retrieve all saved SHU history for a specific user across all years.

**Access Control:**
- **Users** can access their own SHU history
- **Admin/Super Admin** can access any user's SHU history

**Response:**
```json
{
  "data": [
    {
      "id_shu_anggota": 1,
      "id_shu": 5,
      "id_anggota": 17,
      "jumlah_modal": 26098,
      "jumlah_usaha": 68606,
      "shu_diterima": 94704,
      "created_at": "2024-01-15T10:00:00Z",
      "shu": {
        "tahun": 2024,
        "total_shu": 20000000,
        "status": "final"
      }
    },
    {
      "id_shu_anggota": 2,
      "id_shu": 3,
      "id_anggota": 17,
      "jumlah_modal": 45000,
      "jumlah_usaha": 85000,
      "shu_diterima": 130000,
      "created_at": "2023-01-15T10:00:00Z",
      "shu": {
        "tahun": 2023,
        "total_shu": 25000000,
        "status": "final"
      }
    }
  ]
}
```

### List All SHU Anggota Records (Admin Only)
```http
GET /api/shu-anggota
Authorization: Bearer {token}
```

**Description:** Retrieve all saved SHU records for all users (Admin/Super Admin only).

### Get SHU Records by SHU ID (Admin Only)
```http
GET /api/shu-anggota/shu/{shu_id}
Authorization: Bearer {token}

# Example
GET /api/shu-anggota/shu/5
```

**Description:** Retrieve all user SHU records for a specific SHU year (Admin/Super Admin only).

### Delete SHU Anggota Record (Admin Only)
```http
DELETE /api/shu-anggota/{id}
Authorization: Bearer {token}

# Example
DELETE /api/shu-anggota/1
```

**Description:** Delete a saved SHU record (Admin/Super Admin only).

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
| User Management | Self only | Self + registered users | All users |
| Simpanan Wallets | Own wallets | All wallets | All wallets |
| Simpanan Top-up | Own wallets | ❌ | ❌ |
| Simpanan Verify | ❌ | ✅ | ✅ |
| Simpanan Adjust | ❌ | ✅ | ✅ |
| Pinjaman | Own only | Own + registered users | All records |
| Angsuran | Own only | Own + registered users | All records |
| Angsuran Verify | ❌ | Own + registered users | All records |
| SHU Management | ❌ | ✅ | ✅ |

### User Management Access Details:
- **Member**: Can only view/edit/delete their own profile
- **Admin**: Can view/edit/delete themselves + users they registered (where `admin_id` matches their ID)
- **Super Admin**: Can view/edit/delete any user without restrictions

### Simpanan (Wallet) Access Details:
- **Member**: Can view own wallets, request top-ups, view own transaction history
- **Admin/Super Admin**: Can view all wallets, verify/reject top-up requests, adjust balances directly
- **Top-up Flow**: Member creates request → Admin verifies → Balance updated automatically
- **Adjust Flow**: Admin can directly modify wallet balances (creates verified transaction)

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

## Data Models

### User Model
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "address": "123 Main Street, City",
  "phone_number": "081234567890",
  "nik": "1234567890123456",
  "role_id": 3,
  "admin_id": 1,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Fields:**
- `id`: Unique user identifier
- `email`: User's email address (unique)
- `name`: User's full name (required)
- `address`: User's address (optional)
- `phone_number`: User's phone number (optional)
- `nik`: National Identity Number - unique identifier (required)
- `role_id`: Reference to role (1=super_admin, 2=admin, 3=member)
- `admin_id`: ID of the admin who registered this user (null for super_admin users)
- `created_at`: Registration timestamp
- `updated_at`: Last modification timestamp

**Admin Hierarchy:**
- When `admin_id` is `null`: User was created by super_admin or is super_admin
- When `admin_id` has value: User was registered by the admin with that ID
- Admin users can only access users where `admin_id` matches their own ID

### Simpanan (Wallet) Model
```json
{
  "id": 1,
  "user_id": 1,
  "type": "wajib",
  "balance": 1200000,
  "description": "Wallet wajib",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Fields:**
- `id`: Unique wallet identifier
- `user_id`: Owner of the wallet
- `type`: Wallet type (`pokok`, `wajib`, `sukarela`)
- `balance`: Current wallet balance
- `description`: Wallet description
- `created_at`: Wallet creation timestamp
- `updated_at`: Last modification timestamp

**Wallet Types:**
- `pokok`: Basic capital savings
- `wajib`: Mandatory monthly savings
- `sukarela`: Voluntary savings

### Simpanan Transaction Model
```json
{
  "id": 1,
  "simpanan_id": 1,
  "type": "topup",
  "amount": 200000,
  "description": "Monthly mandatory savings",
  "status": "verified",
  "verified_by_id": 2,
  "verified_at": "2024-01-15T11:00:00Z",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

**Fields:**
- `id`: Unique transaction identifier
- `simpanan_id`: Reference to the wallet
- `type`: Transaction type (`topup`, `adjustment`)
- `amount`: Transaction amount (positive/negative)
- `description`: Transaction description
- `status`: Transaction status (`pending`, `verified`, `rejected`)
- `verified_by_id`: Admin who verified the transaction
- `verified_at`: Verification timestamp
- `created_at`: Transaction creation timestamp
- `updated_at`: Last modification timestamp

**Transaction Types:**
- `topup`: User-initiated wallet top-up (requires verification)
- `adjustment`: Admin-initiated balance adjustment (auto-verified)

**Transaction Status:**
- `pending`: Waiting for admin verification
- `verified`: Approved and balance updated
- `rejected`: Rejected by admin

---

## Testing Examples

### Admin User Hierarchy Example

1. **Login as Admin:**
   ```bash
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@example.com","password":"password123"}'
   ```

2. **Admin Creates New User (gets admin_id assigned):**
   ```bash
   curl -X POST http://localhost:8080/api/users \
     -H "Authorization: Bearer {admin_token}" \
     -H "Content-Type: application/json" \
     -d '{"email":"member@example.com","password":"password123","name":"New Member","address":"Member Address","phone_number":"081234567890","nik":"1234567890123456","role_id":3}'
   
   # Response includes admin_id:
   # {"data":{"id":5,"email":"member@example.com","name":"New Member","address":"Member Address","phone_number":"081234567890","nik":"1234567890123456","role_id":3,"admin_id":2}}
   ```

3. **Admin Lists Their Users Only:**
   ```bash
   curl -X GET http://localhost:8080/api/users \
     -H "Authorization: Bearer {admin_token}"
   
   # Only returns users where admin_id matches admin's ID
   ```

4. **Super Admin Sees All Users:**
   ```bash
   curl -X GET http://localhost:8080/api/users \
     -H "Authorization: Bearer {super_admin_token}"
   
   # Returns all users regardless of admin_id
   ```

### Complete Workflow Example

1. **Register and Login:**
   ```bash
   # Register
   curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123","name":"Test User","address":"Test Address","phone_number":"081234567890","nik":"1234567890123456","role_id":3}'
   
   # Login
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. **View Wallets:**
   ```bash
   curl -X GET http://localhost:8080/api/simpanan/wallets \
     -H "Authorization: Bearer {token}"
   ```

3. **Top-up Wallet:**
   ```bash
   curl -X POST http://localhost:8080/api/simpanan/topup \
     -H "Authorization: Bearer {token}" \
     -H "Content-Type: application/json" \
     -d '{"type":"wajib","amount":500000,"description":"Monthly savings"}'
   ```

4. **Admin Verify Transaction:**
   ```bash
   curl -X PUT http://localhost:8080/api/simpanan/transactions/1/verify \
     -H "Authorization: Bearer {admin_token}" \
     -H "Content-Type: application/json" \
     -d '{"approve":true}'
   ```

5. **Apply for Loan:**
   ```bash
   curl -X POST http://localhost:8080/api/pinjaman \
     -H "Authorization: Bearer {token}" \
     -H "Content-Type: application/json" \
     -d '{"jumlah_pinjaman":5000000,"bunga_persen":2.5,"lama_bulan":12,"jumlah_angsuran":450000}'
   ```

6. **Make Payment:**
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