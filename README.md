# BE Payroll System Documentation

## 📘 Overview

This system is designed to manage employee payrolls in a company, handling attendance, overtime, reimbursements, and payslip generation. The backend is built with Go using the Gin framework, PostgreSQL for the database, and JWT for authentication.

---

## ⚙️ Setup Guide

### Prerequisites

* Go >= 1.20
* Docker (for PostgreSQL)
* Make or equivalent CLI tooling

### 1. Clone and Install

```bash
git clone https://github.com/kidzmyujikku/hr-system.git
cd hr-system
go mod tidy
```

### 2. Setup PostgreSQL (via Docker)

```bash
docker run --name hr-db -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres
```

### 3. Seed the Database

```bash
go run seed/seed.go seed
```

### 4. Run the Server

```bash
go run main.go
```

---

## 🔐 Authentication

* JWT-based login via `/login`
* Roles: `admin`, `employee`
* Use `Authorization: Bearer <token>` header for protected routes

---

## 📌 API Endpoints

### Auth

```
POST /login
{ "username": "admin", "password": "admin123" }
```

---

### Admin Routes

#### ➕ Add Attendance Period (PayCycle)

```
POST /admin/pay-cycle
{ "start_date": "2025-06-01", "end_date": "2025-06-30" }
```

#### 🚀 Run Payroll

```
POST /admin/payroll/run
{ "paycycle_id": 1 }
```

#### 📊 Get Summary Report

```
POST /admin/summary
{ "paycycle_id": 1 }
Returns take-home pay per employee and grand total.
```

---

### Employee Routes

#### ⏰ Submit Attendance

```
POST /employee/attendance
```

* Only Mon-Fri allowed
* Automatically detect checkin / checkout
* One submission per day

#### 🕐 Submit Overtime

```
POST /employee/overtime
{ "date": "2025-06-03", "hours": 2 }
```

* Must be after work
* Max 3 hours/day

#### 💰 Submit Reimbursement

```
POST /employee/reimbursement
{ "date": "2025-06-03", "amount": 100000, "description": "Taxi" }
```

#### 📄 Generate Payslip

```
GET /employee/payslip
Returns attendance, overtime, reimbursement breakdown
```

---

## 🧠 Business Rules Summary

* Attendance is paid pro-rata per working day
* Overtime is 2x hourly rate, capped at 3 hours/day
* Reimbursements are added directly to payslip
* Once payroll is run, data in that period is frozen

---

## 🧱 Architecture

* **MVC-style** layered structure
* **Services layer** handles logic (calculation, validation)
* **Controllers** for HTTP I/O
* **Middleware** for JWT, recovery, logging
* **PostgreSQL** via GORM ORM
* **Testing** includes unit and integration tests

---

## ✅ Testing

### Unit Test Examples

* `services/payroll_service_test.go` tests payroll calculations
* `services/auth_service_test.go` tests login logic

### Integration Test Examples

* `tests/auth_controller_test.go`
* `tests/payroll_integration_test.go`

Run all:

```bash
go test ./...
```

---

## 📄 Audit and Tracing

* All models include `created_at`, `updated_at`
* `created_by`, `ip`, and `request_id` tracked where needed
* Audit log (optional table) can store key actions

---

## ✨ Extras

* Seeded with 100 fake employees and 1 admin
* JWT secret from `.env`

---

## 👨‍💻 Contact

Maintainer: `m.izlal2003@gmail.com`
