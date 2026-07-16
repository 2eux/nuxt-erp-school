# Islamic School ERP - API Specification

## Table of Contents

1. [Overview](#overview)
2. [Base URL & Versioning](#base-url--versioning)
3. [Authentication Flow](#authentication-flow)
4. [Common Patterns](#common-patterns)
5. [Error Codes](#error-codes)
6. [Rate Limiting](#rate-limiting)
7. [API Endpoints Summary](#api-endpoints-summary)
8. [Request/Response Examples](#requestresponse-examples)
9. [Real-time Events (SSE)](#real-time-events-sse)

---

## Overview

The ERP API is a **RESTful JSON API** following OpenAPI 3.0 specification. All endpoints return JSON and use standard HTTP methods (GET, POST, PUT, PATCH, DELETE) and status codes.

**Base URL (Production):** `https://api.erp-school.id/api/v1`
**Base URL (Staging):** `https://api.staging.erp-school.id/api/v1`

### Content Types

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |
| `Accept` | `application/json` |

### API Versioning

Versioning is done via URL path prefix: `/api/v1/`, `/api/v2/`. Breaking changes trigger a new major version. Non-breaking additions are made within the same version. Deprecated endpoints return a `Deprecation` header with a sunset date.

---

## Authentication Flow

### Token-Based JWT Authentication

```
┌─────────┐                          ┌──────────┐                       ┌───────────┐
│ Client  │                          │ Backend  │                       │   Redis   │
└────┬────┘                          └────┬─────┘                       └─────┬─────┘
     │  POST /api/v1/auth/login           │                                   │
     │  { email, password, device_info }  │                                   │
     │ ──────────────────────────────────►│                                   │
     │                                    │  Validate credentials             │
     │                                    │  bcrypt.Compare(password, hash)   │
     │                                    │                                   │
     │                                    │  Create session (Redis)           │
     │                                    │──────────────────────────────────►│
     │                                    │                                   │
     │  200 OK                            │                                   │
     │  { access_token (15m),             │                                   │
     │    refresh_token (7d),             │                                   │
     │    expires_in, token_type }        │                                   │
     │ ◄──────────────────────────────────│                                   │
     │                                    │                                   │
     │  GET /api/v1/students              │                                   │
     │  Authorization: Bearer <access>    │                                   │
     │ ──────────────────────────────────►│                                   │
     │                                    │  Validate JWT signature           │
     │                                    │  Check Redis session (not revoked)│
     │                                    │──────────────────────────────────►│
     │                                    │◄──────────────────────────────────│
     │                                    │  Extract: user_id, school_id,    │
     │                                    │  roles, permissions               │
     │                                    │                                   │
     │  200 OK (data)                     │                                   │
     │ ◄──────────────────────────────────│                                   │
     │                                    │                                   │
     │  [When access token expires]       │                                   │
     │  POST /api/v1/auth/refresh         │                                   │
     │  { refresh_token }                 │                                   │
     │ ──────────────────────────────────►│                                   │
     │                                    │  Validate refresh token           │
     │                                    │  Rotate refresh token             │
     │  200 OK (new token pair)           │                                   │
     │ ◄──────────────────────────────────│                                   │
```

### Token Structure

```json
// JWT Payload (decoded)
{
  "sub": "user-uuid",
  "sid": "school-uuid",
  "rol": ["admin", "teacher"],
  "prm": ["student:read", "student:create", "grade:read"],
  "typ": "access",
  "iat": 1700000000,
  "exp": 1700000900,
  "iss": "erp-school"
}
```

### Device-Based Security

- Login requests include `device_info` (user agent, IP hash)
- Refresh tokens are bound to the device that obtained them
- Suspicious device changes trigger re-authentication
- All tokens are revocable via `POST /api/v1/auth/logout`

---

## Common Patterns

### Pagination

All list endpoints support cursor-based pagination by default, with offset-based as an option.

```
GET /api/v1/students?page=1&limit=20&sort=name&order=asc
GET /api/v1/students?cursor=eyJuYW1lIjoiTGFzdCBTdHVkZW50In0&limit=20
```

**Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number (1-indexed) |
| `limit` | integer | 20 | Items per page (max: 100) |
| `cursor` | string | - | Cursor for cursor-based pagination (mutually exclusive with `page`) |
| `sort` | string | `created_at` | Field to sort by |
| `order` | string | `desc` | Sort order: `asc` or `desc` |

**Response:**

```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false,
    "next_cursor": "eyJuYW1lIjoiQWhtYWQifQ",
    "prev_cursor": null
  },
  "links": {
    "self": "/api/v1/students?page=1&limit=20",
    "first": "/api/v1/students?page=1&limit=20",
    "last": "/api/v1/students?page=8&limit=20",
    "next": "/api/v1/students?page=2&limit=20",
    "prev": null
  }
}
```

### Filtering

Filters are passed as query parameters using bracket notation:

```
GET /api/v1/students?filter[name]=Ahmad&filter[grade_id]=uuid&filter[status]=active
```

**Supported filter operators:**

| Operator | Syntax | Example |
|----------|--------|---------|
| Equality | `filter[field]=value` | `filter[status]=active` |
| Contains | `filter[field][contains]=text` | `filter[name][contains]=Ahmad` |
| Greater than | `filter[field][gt]=value` | `filter[age][gt]=18` |
| Less than | `filter[field][lt]=value` | `filter[age][lt]=18` |
| Between | `filter[field][between]=min,max` | `filter[date][between]=2026-01-01,2026-12-31` |
| In array | `filter[field][in]=a,b,c` | `filter[status][in]=active,graduated` |
| Not equal | `filter[field][neq]=value` | `filter[status][neq]=deleted` |

### Field Selection (Sparse Fieldsets)

Reduce payload size by requesting only needed fields:

```
GET /api/v1/students?fields=id,name,grade_name,photo_url
```

### Including Relations

Expand related resources:

```
GET /api/v1/students?include=grade,parents,last_attendance
```

**Response:**

```json
{
  "data": {
    "id": "uuid",
    "name": "Ahmad Hidayat",
    "grade": { "id": "uuid", "name": "Grade 10-A" },
    "parents": [{ "id": "uuid", "name": "Abdullah", "relation": "father" }],
    "last_attendance": { "date": "2026-01-15", "status": "present" }
  }
}
```

### ETag & Conditional Requests

All single-resource responses include an `ETag` header:

```
ETag: "abc123"
```

Clients can use conditional headers:
- `If-None-Match: "abc123"` → returns `304 Not Modified`
- `If-Match: "abc123"` → prevents conflicting updates (412 on mismatch)

### Request ID

Every response includes an `X-Request-ID` header for tracing:

```
X-Request-ID: req_abc123def456
```

---

## Error Codes

### HTTP Status Codes

| Code | Meaning | When |
|------|---------|------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST |
| 202 | Accepted | Async operation enqueued |
| 204 | No Content | Successful DELETE |
| 304 | Not Modified | ETag match |
| 400 | Bad Request | Validation error, malformed JSON |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate resource, stale version |
| 413 | Payload Too Large | Request body exceeds limit |
| 422 | Unprocessable Entity | Semantic validation error |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Unhandled server error |
| 503 | Service Unavailable | Maintenance or overload |

### Error Response Format

All errors follow a consistent structure:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The request contains invalid parameters.",
    "request_id": "req_abc123def456",
    "details": [
      {
        "field": "email",
        "code": "INVALID_FORMAT",
        "message": "Must be a valid email address"
      },
      {
        "field": "age",
        "code": "OUT_OF_RANGE",
        "message": "Must be between 5 and 25",
        "meta": { "min": 5, "max": 25 }
      }
    ]
  }
}
```

### Standard Error Codes

| Code | HTTP | Description |
|------|------|-------------|
| `VALIDATION_ERROR` | 400 | Input validation failed |
| `AUTHENTICATION_REQUIRED` | 401 | No valid token provided |
| `TOKEN_EXPIRED` | 401 | Access token has expired |
| `TOKEN_INVALID` | 401 | Token signature invalid |
| `INSUFFICIENT_PERMISSIONS` | 403 | Missing required role/permission |
| `SCHOOL_NOT_FOUND` | 403 | User's school context invalid |
| `RESOURCE_NOT_FOUND` | 404 | Requested resource does not exist |
| `RESOURCE_ALREADY_EXISTS` | 409 | Duplicate unique field |
| `CONCURRENT_MODIFICATION` | 409 | Resource was modified since last read |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests in time window |
| `INTERNAL_ERROR` | 500 | Unexpected server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |
| `DATABASE_ERROR` | 500 | Database operation failed |
| `FILE_TOO_LARGE` | 413 | Upload exceeds maximum size |
| `UNSUPPORTED_FILE_TYPE` | 400 | File type not allowed |
| `AI_PROVIDER_ERROR` | 502 | AI provider returned an error |
| `AI_RATE_LIMITED` | 429 | AI provider rate limit hit |
| `AI_CONTENT_FILTERED` | 400 | Content blocked by safety filter |

---

## Rate Limiting

### Limits

| Endpoint Group | Limit | Window | Scope |
|----------------|-------|--------|-------|
| General API | 100 req/s | 1 second | Per IP |
| Auth endpoints | 5 req/min | 1 minute | Per IP |
| AI endpoints | 30 req/min | 1 minute | Per IP (user) |
| File upload | 10 req/min | 1 minute | Per IP |
| Admin endpoints | 50 req/s | 1 second | Per user |

### Rate Limit Headers

Every response includes rate limit headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1700000060
Retry-After: 60
```

### Handling 429

When rate limited, the `Retry-After` header indicates seconds to wait:

```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please retry in 60 seconds.",
    "request_id": "req_abc123"
  }
}
```

---

## API Endpoints Summary

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/auth/login` | No | Login with email/password |
| POST | `/api/v1/auth/refresh` | No | Refresh access token |
| POST | `/api/v1/auth/logout` | Yes | Revoke tokens |
| POST | `/api/v1/auth/forgot-password` | No | Request password reset |
| POST | `/api/v1/auth/reset-password` | No | Reset password with token |
| GET | `/api/v1/auth/me` | Yes | Get current user profile |
| PUT | `/api/v1/auth/me` | Yes | Update own profile |
| PUT | `/api/v1/auth/change-password` | Yes | Change password |

### Students

| Method | Endpoint | Permissions | Description |
|--------|----------|-------------|-------------|
| GET | `/api/v1/students` | `student:read` | List students (paginated, filterable) |
| GET | `/api/v1/students/:id` | `student:read` | Get student details |
| POST | `/api/v1/students` | `student:create` | Create student |
| PUT | `/api/v1/students/:id` | `student:update` | Update student |
| DELETE | `/api/v1/students/:id` | `student:delete` | Soft-delete student |
| GET | `/api/v1/students/:id/grades` | `grade:read` | Get student grades |
| GET | `/api/v1/students/:id/attendances` | `attendance:read` | Get student attendance |
| GET | `/api/v1/students/:id/documents` | `document:read` | Get student documents |
| POST | `/api/v1/students/:id/documents` | `document:create` | Upload student document |
| GET | `/api/v1/students/:id/report-card` | `report:read` | Generate report card |

### Teachers / Employees

| Method | Endpoint | Permissions | Description |
|--------|----------|-------------|-------------|
| GET | `/api/v1/teachers` | `teacher:read` | List teachers |
| POST | `/api/v1/teachers` | `teacher:create` | Create teacher |
| GET | `/api/v1/teachers/:id/schedule` | `schedule:read` | Get teacher schedule |
| GET | `/api/v1/employees` | `employee:read` | List all employees |
| POST | `/api/v1/employees/:id/attendance` | `attendance:create` | Mark attendance |

### Academic

| Method | Endpoint | Permissions | Description |
|--------|----------|-------------|-------------|
| GET | `/api/v1/grades` | `grade:read` | List grades |
| GET | `/api/v1/classes` | `class:read` | List classes |
| GET | `/api/v1/subjects` | `subject:read` | List subjects |
| GET | `/api/v1/schedules` | `schedule:read` | Get schedules |
| GET | `/api/v1/attendances` | `attendance:read` | List attendance records |
| POST | `/api/v1/attendances/bulk` | `attendance:create` | Bulk mark attendance |
| GET | `/api/v1/exams` | `exam:read` | List exams |
| POST | `/api/v1/exams/:id/results` | `exam:grade` | Input exam results |
| GET | `/api/v1/gradebooks` | `gradebook:read` | View gradebook |
| GET | `/api/v1/report-cards` | `report:read` | View report cards |
| POST | `/api/v1/report-cards/generate` | `report:create` | Generate report cards |
| GET | `/api/v1/assignments` | `assignment:read` | List assignments |
| POST | `/api/v1/assignments/:id/submissions` | `assignment:submit` | Submit assignment |

### Islamic Modules

| Method | Endpoint | Permissions | Description |
|--------|----------|-------------|-------------|
| GET | `/api/v1/tahfidz/programs` | `tahfidz:read` | List Tahfidz programs |
| POST | `/api/v1/tahfidz/programs` | `tahfidz:create` | Create program |
| GET | `/api/v1/tahfidz/groups` | `tahfidz:read` | List groups |
| POST | `/api/v1/tahfidz/progress` | `tahfidz:update` | Record progress |
| GET | `/api/v1/tahfidz/tasmi` | `tahfidz:read` | Tasmi records |
| POST | `/api/v1/tahfidz/tasmi` | `tahfidz:create` | Submit Tasmi |
| GET | `/api/v1/mutabaah` | `mutabaah:read` | Daily Mutabaah |
| POST | `/api/v1/mutabaah` | `mutabaah:create` | Record Mutabaah |
| GET | `/api/v1/prayer-attendance` | `prayer:read` | Prayer attendance |
| POST | `/api/v1/prayer-attendance/bulk` | `prayer:create` | Record prayer |
| GET | `/api/v1/halaqah/groups` | `halaqah:read` | Halaqah groups |
| GET | `/api/v1/halaqah/members` | `halaqah:read` | Halaqah members |

### Finance

| Method | Endpoint | Permissions | Description |
|--------|----------|-------------|-------------|
| GET | `/api/v1/fees/types` | `fee:read` | List fee types |
| POST | `/api/v1/fees/types` | `fee:create` | Create fee type |
| GET | `/api/v1/invoices` | `invoice:read` | List invoices |
| POST | `/api/v1/invoices` | `invoice:create` | Create invoice |
| POST | `/api/v1/invoices/:id/send` | `invoice:send` | Send invoice to parent |
| GET | `/api/v1/payments` | `payment:read` | List payments |
| POST | `/api/v1/payments` | `payment:create` | Record payment |
| POST | `/api/v1/payments/gateway` | `payment:create` | Create payment gateway |
| POST | `/api/v1/payments/callback/midtrans` | No | Midtrans callback |
| GET | `/api/v1/finance/ledger` | `ledger:read` | General ledger |
| GET | `/api/v1/finance/journals` | `journal:read` | Journal entries |
| GET | `/api/v1/finance/budgets` | `budget:read` | Budget plans |
| GET | `/api/v1/finance/payroll` | `payroll:read` | Payroll records |

### AI Features

| Method | Endpoint | Permissions | Description |
|--------|----------|-------------|-------------|
| POST | `/api/v1/ai/chat` | `ai:use` | Chat with AI assistant |
| POST | `/api/v1/ai/chat/stream` | `ai:use` | Stream chat (SSE) |
| POST | `/api/v1/ai/grade/essay` | `ai:grade` | AI essay grading |
| POST | `/api/v1/ai/generate/lesson-plan` | `ai:generate` | Generate lesson plan |
| POST | `/api/v1/ai/generate/report-comment` | `ai:generate` | Generate report comments |
| POST | `/api/v1/ai/analyze/performance` | `ai:analyze` | Student performance |
| POST | `/api/v1/ai/plagiarism/check` | `ai:check` | Plagiarism detection |
| GET | `/api/v1/ai/conversations` | `ai:read` | Chat history |
| GET | `/api/v1/ai/usage` | `ai:read` | Token usage stats |

### Health & Monitoring

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/health` | No | Health check (liveness) |
| GET | `/api/v1/health/ready` | No | Readiness check |
| GET | `/api/v1/metrics` | No | Prometheus metrics |
| GET | `/swagger/index.html` | No | API documentation |

---

## Request/Response Examples

### Example 1: Create Student

**Request:**
```http
POST /api/v1/students HTTP/1.1
Host: api.erp-school.id
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json
X-Idempotency-Key: idem_abc123

{
  "first_name": "Ahmad",
  "last_name": "Hidayat",
  "nisn": "0051234567",
  "nis": "20240001",
  "gender": "male",
  "birth_date": "2010-03-15",
  "birth_place": "Jakarta",
  "grade_id": "550e8400-e29b-41d4-a716-446655440000",
  "class_id": "550e8400-e29b-41d4-a716-446655440001",
  "admission_date": "2026-01-10",
  "contact": {
    "email": "ahmad@example.com",
    "phone": "+62812345678",
    "address": "Jl. Merdeka No. 1, Jakarta"
  },
  "parents": [
    {
      "name": "Abdullah Hidayat",
      "relation": "father",
      "phone": "+62812345679",
      "email": "abdullah@example.com"
    }
  ]
}
```

**Response:**
```json
{
  "data": {
    "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "first_name": "Ahmad",
    "last_name": "Hidayat",
    "full_name": "Ahmad Hidayat",
    "nisn": "0051234567",
    "nis": "20240001",
    "gender": "male",
    "birth_date": "2010-03-15",
    "birth_place": "Jakarta",
    "grade_id": "550e8400-e29b-41d4-a716-446655440000",
    "class_id": "550e8400-e29b-41d4-a716-446655440001",
    "admission_date": "2026-01-10",
    "status": "active",
    "photo_url": null,
    "school_id": "550e8400-e29b-41d4-a716-446655440003",
    "created_at": "2026-01-15T08:30:00Z",
    "updated_at": "2026-01-15T08:30:00Z"
  }
}
```

### Example 2: Bulk Attendance

**Request:**
```http
POST /api/v1/attendances/bulk HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "date": "2026-01-15",
  "class_id": "550e8400-e29b-41d4-a716-446655440001",
  "subject_id": "550e8400-e29b-41d4-a716-446655440010",
  "records": [
    { "student_id": "uuid-1", "status": "present", "check_in": "07:15:00" },
    { "student_id": "uuid-2", "status": "sick", "note": "Has fever" },
    { "student_id": "uuid-3", "status": "permission", "note": "Family event" },
    { "student_id": "uuid-4", "status": "late", "check_in": "07:45:00" }
  ]
}
```

**Response:**
```json
{
  "data": {
    "date": "2026-01-15",
    "class_id": "550e8400-e29b-41d4-a716-446655440001",
    "total": 4,
    "summary": {
      "present": 1,
      "sick": 1,
      "permission": 1,
      "late": 1,
      "absent": 0
    },
    "records": [
      { "student_id": "uuid-1", "status": "present", "processed": true },
      { "student_id": "uuid-2", "status": "sick", "processed": true },
      { "student_id": "uuid-3", "status": "permission", "processed": true },
      { "student_id": "uuid-4", "status": "late", "processed": true }
    ]
  }
}
```

### Example 3: AI Chat (Streaming)

**Request:**
```http
POST /api/v1/ai/chat/stream HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json
Accept: text/event-stream

{
  "message": "Buatkan jadwal ujian untuk kelas 10 semester genap tahun ajaran 2025/2026",
  "context": {
    "subject": "exam_scheduling",
    "school_level": "high_school"
  }
}
```

**Response (SSE Stream):**
```
event: message
data: {"content": "Baik", "done": false}

event: message
data: {"content": ", saya akan", "done": false}

event: message
data: {"content": " membantu membuat jadwal ujian", "done": false}

...

event: message
data: {"content": "\n\nSemoga membantu!", "done": true, "usage": {"prompt_tokens": 450, "completion_tokens": 320, "total_tokens": 770}}

event: done
data: {"conversation_id": "conv_abc123"}
```

---

## Real-time Events (SSE)

### Event Stream Endpoint

```
GET /api/v1/events/stream
Authorization: Bearer <token>
Accept: text/event-stream
```

### Event Types

| Event | Payload Description |
|-------|---------------------|
| `notification.new` | New notification received |
| `student.created` | Student registered |
| `payment.received` | Payment confirmed |
| `attendance.updated` | Attendance changed |
| `exam.graded` | Exam results published |
| `message.received` | New chat message |
| `schedule.changed` | Schedule modified |
| `system.announcement` | Admin broadcast |

### Example Event

```
event: notification.new
id: evt_abc123
data: {
  "type": "payment_received",
  "title": "Pembayaran Diterima",
  "body": "SPP bulan Januari 2026 untuk Ahmad Hidayat telah lunas",
  "data": {
    "invoice_id": "uuid",
    "amount": 500000,
    "student_name": "Ahmad Hidayat"
  },
  "created_at": "2026-01-15T08:30:00Z"
}

```

---

## Webhook Callbacks

### Payment Gateway Callback

```http
POST /api/v1/payments/callback/midtrans
Content-Type: application/json

{
  "transaction_time": "2026-01-15 08:30:00",
  "transaction_status": "settlement",
  "transaction_id": "abc-123",
  "status_message": "midtrans payment notification",
  "status_code": "200",
  "signature_key": "abc...",
  "payment_type": "bank_transfer",
  "order_id": "INV-2026001",
  "merchant_id": "G123456",
  "gross_amount": "500000.00",
  "fraud_status": "accept",
  "currency": "IDR"
}
```

**Webhook Security:**
- All webhooks require signature verification
- Idempotency via `order_id` (duplicate callbacks are safe)
- IP whitelist for payment gateway servers

---

## SDK & Client Generation

The API specification is available as an OpenAPI 3.0 document at:

```
https://api.erp-school.id/swagger/doc.json
```

### Auto-generated clients:

```
openapi-generator generate \
  -i https://api.erp-school.id/swagger/doc.json \
  -g typescript-fetch \
  -o ./frontend/api-client/
```

### Frontend API client setup (ofetch):

```typescript
// composables/useApi.ts
import { ofetch } from 'ofetch'

export const $api = ofetch.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
  onRequest({ options }) {
    const token = useAuthStore().accessToken
    if (token) {
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${token}`,
      }
    }
  },
  onResponseError({ response }) {
    if (response.status === 401) {
      useAuthStore().logout()
      navigateTo('/login')
    }
  },
})
```
