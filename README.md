# Islamic School ERP

**AI-Native Enterprise School Management Platform for Islamic Education Institutions**

A next-generation, modular, AI-first ERP platform designed specifically for Islamic schools. Supports single-school to multi-school deployments with comprehensive modules for academic, financial, HR, Islamic studies (Tahfidz, Halaqah, Mutaba'ah), and operations management.

---

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Tech Stack](#tech-stack)
5. [Project Structure](#project-structure)
6. [Prerequisites](#prerequisites)
7. [Installation Guide](#installation-guide)
   - [Quick Start (Docker)](#quick-start-docker)
   - [Manual Setup - Backend](#manual-setup---backend)
   - [Manual Setup - Frontend](#manual-setup---frontend)
   - [Manual Setup - AI Gateway](#manual-setup---ai-gateway)
8. [Database Setup](#database-setup)
9. [Configuration](#configuration)
10. [API Documentation](#api-documentation)
11. [User Roles & Permissions](#user-roles--permissions)
12. [Module Reference](#module-reference)
13. [AI Features](#ai-features)
14. [Security](#security)
15. [Deployment](#deployment)
16. [Monitoring & Observability](#monitoring--observability)
17. [Development Guide](#development-guide)
18. [Testing](#testing)
19. [Roadmap](#roadmap)
20. [Contributing](#contributing)
21. [License](#license)

---

## Overview

The Islamic School ERP is a comprehensive **multi-tenant** school management platform built specifically for Islamic educational institutions. It covers:

| Domain | Capabilities |
|--------|-------------|
| **Academic** | Classes, Subjects, Curriculum, Schedules, Attendance, Exams, Gradebooks, Report Cards |
| **Islamic** | Tahfidz/Tahsin/Tilawah, Mutaba'ah Yaumiyah, Prayer Tracking, Halaqah, Islamic Character |
| **Finance** | SPP, Invoices, Payments, GL, Journals, Budget (RKAS), Payroll, Cashflow |
| **HR** | Employees, Attendance, Leave, Payroll, Training, Performance |
| **Students** | Enrollment, Parents, Promotion, Graduation, Documents, Medical |
| **Admissions** | PPDB (Registration), Exams, Selection, Enrollment |
| **Library** | Books, Borrowing, Returns, Reservations |
| **Inventory** | Assets, Items, Stock, Procurement |
| **AI** | Lesson Plan Generator, Quiz Generator, RAG Document Q&A, Tahfidz Planner, Financial Analysis |
| **Analytics** | Executive Dashboards, KPI Tracking, Predictive Insights |

### Key Principles

- **Domain-Driven Design** with bounded contexts
- **CQRS-ready** for read/write separation
- **Event-driven** communication between bounded contexts
- **Multi-tenancy** via `school_id` with Row-Level Security
- **API-first** design with versioned REST APIs + gRPC-ready architecture
- **MCP-compatible** AI integration (Model Context Protocol)
- **Zero-trust** security with JWT + RBAC

---

## Features

### Core Platform
- Multi-school, multi-campus, multi-academic-year, multi-branch
- Multi-language (i18n: English, Bahasa Indonesia, Arabic)
- Multi-timezone support
- Dark mode / Light mode
- Progressive Web App (PWA) with offline support
- Responsive, mobile-first design
- Audit trail and activity logs on all mutations
- Soft delete with versioning
- Approval workflows for documents, leave, procurement
- Dynamic settings per school

### Islamic-Specific
- **Tahfidz Management** — Programs, groups, memorization targets, progress tracking
- **Tasmi' Records** — Oral Quran recitation exam records
- **Mutaba'ah Yaumiyah** — Daily worship tracking (5 prayers, sunnah, dhuha, tahajjud, dhikr, charity)
- **Prayer Attendance** — Per-student daily prayer logging
- **Halaqah Groups** — Quran study circles with schedule and member management
- **Islamic Character** — Adab & Akhlaq assessment, behavior notes
- **Quranic Competencies** — Tajwid, Tilawah, Tahfidz, Tafsir scoring
- **ZISWAF** — Zakat, Infaq, Shadaqah, Waqf management
- **Islamic Events** — Ramadhan programs, Islamic calendar events

### AI Capabilities
- **RAG** (Retrieval Augmented Generation) — Upload school documents and query them via chat
- **Multiple LLM Providers** — OpenAI, Google Gemini, Anthropic Claude, local Ollama
- **25+ ERP Tools** via MCP — Lesson plans, quizzes, exams, worksheets, rubrics
- **Semantic Cache** — Similar question deduplication in Redis
- **Smart Routing** — Automatic provider fallback and cost optimization
- **Prompt Guard** — Injection detection, content filtering

---

## Architecture

```
                        ┌──────────────────────────────────────┐
                        │           Cloudflare CDN             │
                        └──────────────┬───────────────────────┘
                                       │
                                       ▼
┌──────────────┐    ┌──────────────────────────────────────────────┐
│   Mobile PWA │───▶│                NGINX REVERSE PROXY            │
│  (Nuxt SSR)  │    └─────┬──────────────────┬───────────────────┬──┘
└──────────────┘          │                  │                   │
                          ▼                  ▼                   ▼
               ┌──────────────────┐ ┌──────────────┐ ┌──────────────────┐
               │   FRONTEND       │ │  BACKEND     │ │  AI GATEWAY      │
               │  Nuxt 4 (SSR)    │ │  Go / Gin    │ │  Go              │
               │                  │ │              │ │                  │
               │ - SSR + SPA      │ │ - REST APIs  │ │ - Prompt Router  │
               │ - PWA + Offline  │ │ - JWT + RBAC │ │ - Rate Limiting  │
               │ - Pinia + TanQ   │ │ - Validation │ │ - Response Cache │
               │ - Tailwind CSS 4 │ │ - File Mgmt  │ │ - RAG Pipeline   │
               └──────────────────┘ └──────┬───────┘ └────────┬─────────┘
                                           │                  │
                      ┌────────────────────┼──────────────────┼───────────────────────┐
                      │                    │                  │                       │
                      ▼                    ▼                  ▼                       ▼
            ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐
            │ PostgreSQL 16│    │   Redis 7    │    │  Qdrant 1.12 │    │  Ollama (Local)  │
            │  125 tables  │    │ Cache+Queue  │    │ Vector DB    │    │  LLM Server      │
            └──────────────┘    └──────────────┘    └──────────────┘    └──────────────────┘
                      │
                      ▼
            ┌──────────────┐     ┌──────────────────────────────────────────────┐
            │  MinIO (S3)  │     │              OBSERVABILITY                    │
            │  Object Store│     │  Prometheus + Grafana + Jaeger                │
            └──────────────┘     └──────────────────────────────────────────────┘
```

### Data Flow

```
Login:  User → Frontend → POST /api/v1/auth/login → Backend → JWT (15m) + Refresh (7d)
CRUD:   User → Frontend → GET/POST/PUT/DELETE /api/v1/* → Backend → PostgreSQL
RAG:    User → Frontend → POST /api/v1/ai/chat → Backend → AI Gateway → Qdrant → LLM
Events: Backend → Redis Pub/Sub → Notification Service → Email/SMS/Push/In-App
```

---

## Tech Stack

| Layer           | Technology                                             |
|-----------------|--------------------------------------------------------|
| **Frontend**    | Nuxt 4, Vue 3.5, TypeScript, Tailwind CSS 4, Nuxt UI   |
| **State**       | Pinia, TanStack Query, VueUse                          |
| **Charts**      | ApexCharts, vue3-apexcharts                            |
| **Validation**  | Zod, vee-validate                                      |
| **Backend**     | Go 1.23+, Gin Framework, sqlx, pgx                     |
| **Database**    | PostgreSQL 16 (JSONB, FTS, RLS, 125 tables)            |
| **Cache**       | Redis 7 (Sessions, Cache, Queue, Pub/Sub)              |
| **Storage**     | MinIO (S3-compatible object storage)                   |
| **Vector DB**   | Qdrant 1.12 (Semantic search, RAG)                     |
| **AI/LLM**      | OpenAI GPT-4o, Google Gemini, Anthropic Claude, Ollama |
| **Auth**        | JWT (RS256), Refresh Token Rotation, RBAC              |
| **API Docs**    | Swagger / OpenAPI 3.0                                  |
| **Proxy**       | Nginx 1.27 (SSL, CSP, HSTS, Rate Limiting)             |
| **CI/CD**       | GitHub Actions (Lint, Test, Build, Deploy)             |
| **Monitoring**  | Prometheus, Grafana, Jaeger (OpenTelemetry)            |
| **Container**   | Docker, Docker Compose                                 |
| **Messaging**   | Redis Pub/Sub (dev) → Kafka/Redpanda (scale)           |

---

## Project Structure

```
nuxt-erp-school/
├── backend/                          # Go REST API server
│   ├── cmd/server/main.go            # Entry point with DI and graceful shutdown
│   ├── internal/
│   │   ├── config/                   # Configuration (env, viper)
│   │   ├── domain/                   # Domain entities and errors
│   │   ├── dto/                      # Request/Response DTOs
│   │   ├── repository/               # Data access layer (4 repositories)
│   │   ├── service/                  # Business logic layer (17 services)
│   │   ├── handler/                  # HTTP handlers (14 handlers + router)
│   │   ├── middleware/               # Auth, RBAC, CORS, Rate Limit, Audit
│   │   ├── infrastructure/           # PostgreSQL, Redis, MinIO clients
│   │   └── ai/                       # AI Gateway client, RAG, prompts
│   ├── migrations/                   # SQL migrations (125 tables)
│   ├── pkg/                          # Shared utilities (hash, validator, pagination)
│   └── docs/                         # Swagger annotations
│
├── frontend/                         # Nuxt 4 application
│   ├── pages/                        # 55+ page modules
│   ├── components/
│   │   ├── app/                      # App shell (Sidebar, Header, Notifications)
│   │   └── ui/                       # Reusable components (DataTable, FormDialog, etc.)
│   ├── composables/                  # 6 composables (useApi, useAuth, useSchool, etc.)
│   ├── stores/                       # 3 Pinia stores (Auth, School, App)
│   ├── layouts/                      # 3 layouts (Default, Auth, Blank)
│   ├── middleware/                    # Auth guard middleware
│   ├── plugins/                      # API client, Day.js, ApexCharts
│   ├── types/                        # 54 TypeScript interfaces
│   ├── i18n/locales/                 # en.ts, id.ts
│   └── server/api/                   # Nuxt server proxy to Go backend
│
├── ai-gateway/                       # AI orchestration service
│   ├── cmd/server/main.go            # Entry point
│   ├── internal/
│   │   ├── providers/                # OpenAI, Gemini, Claude, Ollama clients
│   │   ├── mcp/                      # MCP server, tools, resources, prompts
│   │   ├── rag/                      # Qdrant client, document processing, embeddings
│   │   ├── handlers/                 # Chat, MCP, RAG endpoints
│   │   ├── router/                   # Provider routing with fallback
│   │   ├── cache/                    # Redis semantic cache
│   │   └── middleware/               # Auth, logging
│   └── .env.example
│
├── docker/                           # Docker configurations
│   ├── backend.Dockerfile            # Multi-stage Go build
│   ├── frontend.Dockerfile           # Multi-stage Node + Nginx
│   ├── ai-gateway.Dockerfile         # Multi-stage Go build
│   ├── nginx.conf                    # Reverse proxy with SSL, CSP, HSTS
│   ├── prometheus.yml                # Scrape configs + alert rules
│   └── grafana/                      # Datasources, dashboards
│
├── .github/workflows/                # CI/CD pipelines
│   ├── ci.yml                        # Lint, test, build, security scan
│   └── deploy.yml                    # Staging + production deploy
│
├── docs/architecture/                # Architecture documentation
│   ├── ARCHITECTURE.md               # Full system architecture
│   ├── API_SPEC.md                   # API specification
│   └── ERD.md                        # Entity relationship diagram
│
├── docker-compose.yml                # 11-service orchestration
└── .env.example                      # Environment variables reference
```

---

## Prerequisites

### Required
- **Docker** 24+ and **Docker Compose** v2+
- **Git**

### For Manual Development
- **Go** 1.23+
- **Node.js** 22+
- **PostgreSQL** 16+
- **Redis** 7+
- **pnpm** 9+ (recommended) or npm

### AI Features (Optional)
- OpenAI API key (for GPT-4o)
- Google Gemini API key
- Anthropic Claude API key
- Or: Run Ollama locally (requires 8GB+ RAM, GPU recommended)

---

## Installation Guide

### Quick Start (Docker)

The fastest way to get the entire platform running:

```bash
# 1. Clone the repository
git clone https://github.com/opencode/nuxt-erp-school.git
cd nuxt-erp-school

# 2. Copy and configure environment
cp .env.example .env
# Edit .env — set your secrets (JWT keys, DB passwords, AI keys)

# 3. Start all services
docker compose up -d

# 4. Check service status
docker compose ps

# 5. View logs
docker compose logs -f

# 6. Access the platform
# Frontend:  http://localhost:3000
# API:       http://localhost:8080
# AI Gateway: http://localhost:8081
# MinIO:      http://localhost:9001
# Grafana:    http://localhost:3001
# Qdrant:     http://localhost:6333
# Jaeger:     http://localhost:16686
```

### Services Started

| Service      | Port  | Description                          |
|-------------|-------|--------------------------------------|
| Frontend    | 3000  | Nuxt 4 SSR web application           |
| Backend     | 8080  | Go REST API server                   |
| AI Gateway  | 8081  | AI orchestration service             |
| PostgreSQL  | 5432  | Primary relational database          |
| Redis       | 6379  | Cache, sessions, job queue           |
| MinIO       | 9000  | S3-compatible object storage         |
| MinIO Console | 9001 | MinIO web console                    |
| Qdrant      | 6333  | Vector database for AI/RAG           |
| Nginx       | 80/443 | Reverse proxy with SSL termination   |
| Ollama      | 11434 | Local LLM (requires `--profile ai-full`) |
| Prometheus  | 9090  | Metrics collection                   |
| Grafana     | 3001  | Monitoring dashboards                |
| Jaeger      | 16686 | Distributed tracing                  |

### Optional: Enable Local AI (Ollama)

```bash
# Start with Ollama profile
docker compose --profile ai-full up -d ollama

# Pull models for local AI
docker exec -it erp-ollama ollama pull llama3.2
docker exec -it erp-ollama ollama pull nomic-embed-text
```

### Database Initialization

The database is automatically initialized when the `postgres` container starts:
- Extensions enabled: `uuid-ossp`, `pg_trgm`, `pgcrypto`
- Schema created from `backend/migrations/000001_init_schema.up.sql`
- All 125 tables, 307 indexes, and constraints applied

To manually run migrations:

```bash
# Using docker compose
docker compose exec backend go run cmd/migrate/main.go up

# Or connect directly
psql -h localhost -p 5432 -U erp_admin -d erp_school
```

### Default Admin Account

On first startup, seed the super admin account:

```bash
docker compose exec backend go run cmd/seed/main.go
```

Default credentials (change immediately after first login):
- Email: `admin@erp-school.id`
- Password: `Admin@2026!`

---

### Manual Setup — Backend

```bash
# Navigate to backend
cd backend

# Copy environment
cp .env.example .env
# Edit .env with your database credentials

# Install dependencies
go mod download
go mod tidy

# Run database migrations
go run cmd/migrate/main.go up

# Seed initial data
go run cmd/seed/main.go

# Start development server
go run cmd/server/main.go

# Or with hot-reload (install air first: go install github.com/air-verse/air@latest)
air
```

### Manual Setup — Frontend

```bash
cd frontend

# Copy environment
cp .env.example .env

# Install dependencies
pnpm install
# or: npm install

# Start development server
pnpm dev
# or: npm run dev

# Build for production
pnpm build

# Preview production build
pnpm preview
```

### Manual Setup — AI Gateway

```bash
cd ai-gateway

# Copy environment
cp .env.example .env
# Set your AI provider API keys

# Install dependencies
go mod download
go mod tidy

# Start server
go run cmd/server/main.go
```

---

## Database Setup

### Schema Overview

The database contains **125 tables** across **24 modules** with **307 indexes**.

```
Module          Tables   Key Entities
─────────────── ──────── ───────────────────────────────────────────
Core/Master      7       schools, academic_years, semesters, grades, classes
Users/RBAC       7       users, roles, permissions, sessions
Students         5       students, parents, documents, behavior
Teachers/HR      7       teachers, employees, attendances, leave
Academic        14       schedules, attendances, exams, gradebooks, report cards
Islamic         12       tahfidz, halaqah, mutabaah, prayer, tasmi
Finance         18       fees, invoices, payments, journals, ledger, payroll
Inventory        5       assets, items, procurements
Library          4       books, borrowings, reservations
Medical          3       records, immunization, health info
Counseling       1       counseling sessions
Transport        3       routes, students, attendance
Dormitory        4       dormitories, rooms, residents
Canteen          3       products, orders, items
Notifications    3       notifications, templates, preferences
Documents        8       documents, letters, workflows, approvals
Meetings         5       meetings, events, tasks, task boards
Announcements    1       announcements
Settings         2       school settings, system settings
Audit            3       audit logs, activity logs, login logs
Islamic Events   6       mosque, islamic events, ramadhan, zakat, infaq, waqf
Admissions       4       batches, applicants, exams, documents
Graduation       2       batches, candidates
AI/Knowledge     6       knowledge docs, chunks, conversations, messages, prompts
──────────────────────
Total         ~125
```

### Key Design Decisions

- **UUID Primary Keys** — All tables use `UUID DEFAULT gen_random_uuid()`
- **Multi-Tenancy** — Every tenant table has `school_id` FK with NOT NULL
- **Soft Delete** — Uses `deleted_at TIMESTAMPTZ` for logical deletion
- **Audit Trail** — `created_by`, `updated_by`, `created_at`, `updated_at` on all tables
- **Full-Text Search** — GIN indexes on name/description columns
- **JSONB** — Flexible metadata, settings, and permissions storage
- **CHECK Constraints** — All status/enum fields use CHECK constraints
- **Row-Level Security** — PostgreSQL RLS policies enforce tenant isolation

---

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

#### Critical (Must Change in Production)

```env
# JWT Secrets (generate: openssl rand -hex 64)
JWT_ACCESS_SECRET=your_access_secret_min_64_chars
JWT_REFRESH_SECRET=your_refresh_secret_min_64_chars

# Database
DB_PASSWORD=strong_db_password

# Redis
REDIS_PASSWORD=strong_redis_password

# MinIO
MINIO_SECRET_KEY=strong_minio_password

# AI Providers (at least one required for AI features)
AI_OPENAI_KEY=sk-your-openai-api-key
# or
AI_GEMINI_KEY=your-gemini-api-key
# or
AI_CLAUDE_KEY=your-claude-api-key
```

#### All Variables

| Category       | Variable              | Description                    | Default              |
|---------------|----------------------|-------------------------------|----------------------|
| App           | `APP_ENV`             | Environment                    | `production`         |
| App           | `APP_PORT`            | Backend port                   | `8080`               |
| Database      | `DB_HOST`             | PostgreSQL host                | `postgres`           |
| Database      | `DB_PORT`             | PostgreSQL port                | `5432`               |
| Database      | `DB_NAME`             | Database name                  | `erp_school`         |
| Redis         | `REDIS_HOST`          | Redis host                     | `redis`              |
| Redis         | `REDIS_PASSWORD`      | Redis password                 | (required)           |
| MinIO         | `MINIO_ENDPOINT`      | MinIO endpoint                 | `minio:9000`         |
| JWT           | `JWT_ACCESS_SECRET`   | Access token signing key       | (required)           |
| JWT           | `JWT_ACCESS_TTL`      | Access token lifetime          | `15m`                |
| JWT           | `JWT_REFRESH_TTL`     | Refresh token lifetime         | `168h` (7 days)      |
| AI            | `AI_PROVIDER`         | Default AI provider            | `openai`             |
| AI            | `AI_DEFAULT_MODEL`    | Default LLM model              | `gpt-4o`             |
| Qdrant        | `QDRANT_API_KEY`      | Qdrant API key                 | (required)           |
| SMTP          | `SMTP_HOST`           | SMTP server                    | `smtp.gmail.com`     |
| Payment       | `PAYMENT_GATEWAY_KEY` | Midtrans server key            | (optional)           |
| Feature Flags | `FEATURE_AI_ENABLED`  | Enable AI features             | `true`               |

---

## API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

All protected endpoints require a JWT Bearer token:

```http
Authorization: Bearer <access_token>
```

### Key Endpoints

| Method   | Endpoint                      | Auth | Description              |
|----------|-------------------------------|------|--------------------------|
| POST     | `/api/v1/auth/login`          | No   | Login with email/password |
| POST     | `/api/v1/auth/refresh`        | No   | Refresh access token      |
| POST     | `/api/v1/auth/logout`         | Yes  | Revoke tokens             |
| GET      | `/api/v1/auth/me`             | Yes  | Current user profile      |
| GET      | `/api/v1/students`            | Yes  | List students (paginated) |
| POST     | `/api/v1/students`            | Yes  | Create student             |
| GET      | `/api/v1/teachers`            | Yes  | List teachers              |
| GET      | `/api/v1/schedules`           | Yes  | Get schedules              |
| POST     | `/api/v1/attendances/bulk`    | Yes  | Bulk attendance entry      |
| GET      | `/api/v1/gradebooks`          | Yes  | View gradebook             |
| POST     | `/api/v1/report-cards/generate` | Yes | Generate report cards     |
| GET      | `/api/v1/tahfidz/progress`    | Yes  | Tahfidz progress           |
| POST     | `/api/v1/mutabaah`            | Yes  | Record daily worship       |
| GET      | `/api/v1/invoices`            | Yes  | List invoices              |
| GET      | `/api/v1/analytics/dashboard` | Yes  | Executive dashboard        |
| POST     | `/api/v1/ai/chat`             | Yes  | AI chat (RAG)              |
| POST     | `/api/v1/ai/generate`         | Yes  | AI content generation      |

### Pagination

All list endpoints support:

```
GET /api/v1/students?page=1&limit=20&sort=name&order=asc
GET /api/v1/students?filter[grade_id]=uuid&filter[status]=active
```

Response includes `meta` with `total`, `total_pages`, `has_next`, `has_prev`, and `links`.

### Swagger UI

```
http://localhost:8080/swagger/index.html
```

### Full API Specification

See [`docs/architecture/API_SPEC.md`](docs/architecture/API_SPEC.md) for complete endpoint listing, request/response examples, error codes, and SSE event streams.

---

## User Roles & Permissions

| Role                | Description                                          |
|---------------------|------------------------------------------------------|
| **Super Admin**     | System-wide access, multi-school management          |
| **Foundation**      | Yayasan-level oversight                              |
| **Principal**       | Full school management                               |
| **Vice Principal**  | Academic and administrative support                   |
| **Academic Staff**  | Curriculum, scheduling, gradebook                     |
| **Finance**         | SPP, invoices, payments, journals, payroll            |
| **HR**              | Employee management, leave, payroll                   |
| **Teacher**         | Teaching, attendance, assignments, grading            |
| **Homeroom Teacher**| Class management + teacher permissions                |
| **Quran Teacher**   | Tahfidz, Tahsin, Tilawah, Halaqah management          |
| **Student**         | View own records, submit assignments                  |
| **Parent**          | View children's records, payments, messages           |
| **Library Staff**   | Book catalog, borrowing management                    |
| **Clinic Staff**    | Medical records, health info                          |
| **Security**        | Gate check, visitor management                        |
| **Guest**           | Limited read-only access                              |

Permissions follow `resource:action` format (e.g., `student:read`, `student:create`, `grade:update`, `invoice:delete`). All permissions are fully configurable per role through the admin panel.

---

## Module Reference

### Dashboard
Real-time KPI cards, enrollment trends, revenue charts, prayer attendance summary, Tahfidz progress overview, recent activities, upcoming events.

### Student Management
Full student lifecycle: enrollment → class assignment → academic tracking → promotion → graduation. Parent linking, document management, behavior tracking.

### Academic
Classes, subjects, curriculum, weekly schedules (timetable grid), attendance tracking (individual + bulk), exam management, gradebook with calculation, report card generation (PDF), assignments with submission review, lesson plans, teaching journals.

### Islamic Modules
- **Tahfidz**: Programs, groups, daily progress (new memorization + muroja'ah), memorization targets, Tasmi' recording, analytics charts
- **Mutaba'ah Yaumiyah**: Daily 5 prayers, sunnah, dhuha, tahajjud, Quran reading, dhikr, charity
- **Prayer Attendance**: Per-student daily log, class summaries, compliance charts
- **Halaqah**: Quran study circles, schedule, member management
- **Islamic Character**: Adab & Akhlaq assessment with radar charts
- **Quranic Competencies**: Tajwid, Tilawah, Tahfidz, Tafsir scoring
- **ZISWAF**: Zakat, Infaq, Shadaqah, Waqf records
- **Islamic Events**: Maulid, Isra Mi'raj, Ramadhan programs, Eid management

### Finance
Fee type configuration (SPP, registration, development, uniform, books, exams, events), fee assignment, invoice generation (single + bulk), payment recording and verification, payment gateway integration (Midtrans), general ledger, journal entries (double entry), budget planning (RKAS), payroll processing with BPJS/tax calculation, cashflow tracking, bank reconciliation.

### HR
Employee database, attendance tracking, leave management (submit → approve → reject workflow), payroll integration, training records, performance evaluation.

### Admissions (PPDB)
Batch management, applicant registration, document verification, admission exams (written, interview, Quran test), selection workflow, automatic enrollment.

### Additional Modules
- **Library**: Book catalog, borrowing, returns with late fees, reservations
- **Medical**: Health records, immunization tracking, allergy/condition alerts
- **Counseling**: Session records, action plans, follow-ups
- **Inventory**: Item stock tracking, asset management, procurement workflow
- **Transportation**: Routes, student assignments, pickup/dropoff attendance
- **Dormitory**: Room management, resident tracking, attendance
- **Canteen**: Product catalog, order management
- **Announcements**: Rich text editor, target audience, pinning
- **Meetings**: Scheduling, participant tracking, minutes
- **Documents**: Upload, versioning, folders, approval workflow
- **Letters**: Template-based generation, numbering, approval
- **Calendar**: Monthly/weekly view, event management
- **Messages**: Internal messaging, inbox/compose

---

## AI Features

### AI Gateway

The AI Gateway is a dedicated service that orchestrates all AI interactions. It supports:

| Provider  | Models                                  |
|-----------|----------------------------------------|
| OpenAI    | GPT-4o, GPT-4o-mini, GPT-3.5-turbo     |
| Gemini    | Gemini 1.5 Pro, Gemini 1.5 Flash       |
| Claude    | Claude 3.5 Sonnet, Claude 3 Opus       |
| Ollama    | llama3.2, mistral, codellama (local)   |

### MCP Tools (25+ registered tools)

The platform implements the **Model Context Protocol (MCP)** with the following tools:

**Academic AI:**
- `lesson_plan_generator` — Generate structured lesson plans
- `quiz_generator` — Create quizzes with answer keys
- `exam_generator` — Generate exams with marking schemes
- `worksheet_generator` — Produce student worksheets
- `rubric_generator` — Create assessment rubrics
- `homework_generator` — Generate homework assignments
- `student_feedback_generator` — Provide personalized feedback
- `report_card_comment_generator` — Auto-generate report comments

**Islamic AI:**
- `tahfidz_planner` — Create memorization schedules
- `memorization_recommendation` — Suggest targets based on progress
- `tilawah_assistant` — Help with Quran recitation
- `tajwid_assistant` — Tajwid rule explanations
- `arabic_vocabulary_trainer` — Interactive Arabic learning
- `islamic_quiz_generator` — Islamic knowledge quizzes
- `islamic_story_generator` — Generate Islamic stories
- `akhlaq_recommendation` — Character development suggestions
- `mutabaah_analyzer` — Worship consistency analysis

**Finance AI:**
- `rkas_assistant` — Budget planning assistant
- `budget_recommendation` — Budget allocation suggestions
- `cashflow_predictor` — Predictive cashflow analysis
- `payroll_assistant` — Payroll calculation helper
- `financial_health_analyzer` — Financial KPI analysis
- `donation_analytics` — Zakat/Infaq/Shadaqah analytics

**Administration AI:**
- `meeting_minutes_generator` — Meeting minutes from notes
- `letter_generator` — Official letter generation
- `sop_generator` — Standard operating procedure docs
- `policy_assistant` — Policy recommendation
- `hr_assistant` — HR document generation
- `recruitment_assistant` — Teacher recruitment support

### RAG (Document AI)

Upload school documents (PDF, DOCX, Excel, PowerPoint, images) and query them via natural language:

```bash
# Upload document
POST /api/v1/ai/knowledge/upload

# Query documents
POST /api/v1/ai/knowledge/query
{
  "query": "Apa kebijakan sekolah tentang seragam?",
  "filters": { "document_type": "school_policy" }
}
```

Pipeline: Upload → OCR/Text Extraction → Chunking → Embedding → Store in Qdrant → Semantic Search → LLM Answer

---

## Security

### Defense in Depth

```
Layer 1: Network
  ├── Cloudflare DDoS + WAF (production)
  ├── Nginx rate limiting
  ├── TLS 1.2+ with strong ciphers
  └── Internal Docker network isolation

Layer 2: Authentication
  ├── JWT RS256 (asymmetric)
  ├── Access token: 15min TTL
  ├── Refresh token: 7d TTL + rotation
  ├── Device fingerprinting
  └── Brute force protection

Layer 3: Authorization
  ├── RBAC (resource:action granularity)
  ├── Multi-tenancy (school_id isolation)
  ├── PostgreSQL Row-Level Security
  └── API scope validation

Layer 4: Application
  ├── Input validation (Zod + go-validator)
  ├── SQL injection prevention (parameterized queries)
  ├── XSS prevention (CSP + output encoding)
  ├── CSRF protection (SameSite cookies)
  └── File upload scanning

Layer 5: Data
  ├── Encryption at rest
  ├── Encryption in transit (TLS)
  ├── PII encryption (pgcrypto)
  ├── Audit trail on all mutations
  └── Encrypted backups

Layer 6: Monitoring
  ├── Prometheus alerts
  ├── Grafana dashboards
  ├── Jaeger tracing
  └── Incident response runbooks
```

### Security Headers (Applied by Nginx)

| Header                         | Value                                        |
|-------------------------------|----------------------------------------------|
| Content-Security-Policy        | Strict CSP with nonce-based scripts          |
| Strict-Transport-Security      | `max-age=63072000; includeSubDomains; preload` |
| X-Frame-Options               | `DENY`                                        |
| X-Content-Type-Options        | `nosniff`                                     |
| Referrer-Policy               | `strict-origin-when-cross-origin`             |
| Permissions-Policy            | Camera/mic/geo disabled                       |

---

## Deployment

### Docker Compose (Recommended)

```bash
# Production startup
docker compose -f docker-compose.yml up -d

# With local AI (GPU recommended)
docker compose --profile ai-full up -d

# Scale backend instances
docker compose up -d --scale backend=3
```

### Production Checklist

- [ ] Change all default passwords and secrets
- [ ] Generate strong JWT secrets (`openssl rand -hex 64`)
- [ ] Configure valid SSL certificates
- [ ] Set up DNS records for your domain
- [ ] Configure SMTP for email notifications
- [ ] Set up Cloudflare or similar CDN/WAF
- [ ] Configure automated database backups
- [ ] Set `APP_ENV=production` and `APP_DEBUG=false`
- [ ] Enable Prometheus alerting to Slack/Email
- [ ] Review and restrict CORS origins
- [ ] Test disaster recovery procedures

### CI/CD Pipeline

```yaml
GitHub Push → GitHub Actions
  ├── Lint (golangci-lint, ESLint)
  ├── Test (go test -race, vitest)
  ├── Security Scan (gosec, trivy)
  ├── Docker Build & Push → ghcr.io
  └── Deploy to Staging
       ├── Deploy to Production (manual approval)
       ├── Health Check
       └── Rollback (if failed)
```

### Kubernetes (Planned)

For large-scale deployments, migration path to K3s/GKE with:
- Horizontal Pod Autoscaler for backend/frontend
- PostgreSQL with read replicas
- Redis Cluster for session management
- Kafka for event streaming
- gRPC service mesh (Istio/Linkerd)

---

## Monitoring & Observability

### Prometheus Metrics

Available at `http://localhost:9090`:
- HTTP request rate, latency, error rate
- Database connection pool status
- Redis hit/miss ratio
- AI token usage and cost
- Active user sessions

### Grafana Dashboards

Available at `http://localhost:3001`:
- **ERP Overview**: 20 panels covering system health, API performance, DB metrics, AI usage, and business KPIs
- Login: `admin` / `change_me_grafana_pass` (default)

### Jaeger Tracing

Available at `http://localhost:16686`:
- Distributed tracing across all services
- Request flow visualization
- Latency breakdown by service

### Alert Rules (Pre-configured)

| Alert                    | Condition                    | Severity |
|-------------------------|------------------------------|----------|
| Service Down            | Instance unreachable > 1min  | Critical |
| High API Error Rate     | 5xx > 5% for 5min           | Warning  |
| Database Pool Exhausted | Active > 80% of max          | Critical |
| Redis Unreachable       | Connection failures > 1min   | Critical |
| High CPU Usage          | CPU > 80% for 10min           | Warning  |
| High Memory Usage       | Memory > 85% for 5min        | Warning  |
| Slow Response Time      | P95 latency > 2s for 5min    | Warning  |

---

## Development Guide

### Backend (Go)

```bash
cd backend

# Run tests
go test ./... -race -cover

# Run linter
golangci-lint run

# Generate Swagger docs
swag init -g cmd/server/main.go -o docs/

# Build
go build -o bin/server cmd/server/main.go
```

**Key Packages:**
- `internal/domain/` — Entity definitions
- `internal/repository/` — Data access (sqlx)
- `internal/service/` — Business logic
- `internal/handler/` — HTTP controllers (Gin)
- `internal/middleware/` — Auth, RBAC, CORS, Rate Limit

**Patterns:**
- Clean Architecture (Domain → Repository → Service → Handler)
- Dependency Injection (manual, in `cmd/server/main.go`)
- Repository pattern with interfaces
- Service layer for business logic
- Standard `dto.APIResponse` wrapper for all responses

### Frontend (Nuxt 4)

```bash
cd frontend

# Install dependencies
pnpm install

# Development
pnpm dev

# Type checking
pnpm typecheck

# Linting
pnpm lint

# Testing
pnpm test

# Build
pnpm build
```

**Key Conventions:**
- Pages in `pages/` use `definePageMeta` for layout + middleware
- API calls through `useApi` composable (auto-injects auth token)
- State managed in Pinia stores (`auth`, `school`, `app`)
- Server state cached with TanStack Query
- Components imported automatically (Nuxt auto-imports)
- i18n via `$t()` function (locales in `i18n/locales/`)

### AI Gateway (Go)

```bash
cd ai-gateway

# Run tests
go test ./... -v

# Build
go build -o bin/ai-gateway cmd/server/main.go
```

---

## Testing

### Backend Tests
```bash
cd backend
go test ./internal/... -race -cover
go test ./internal/service/... -v
```

### Frontend Tests
```bash
cd frontend
pnpm test          # Run vitest
pnpm test:coverage # With coverage
```

### AI Gateway Tests
```bash
cd ai-gateway
go test ./... -race
```

---

## Roadmap

### Phase 1 — Foundation (Completed)
- [x] Project architecture and structure
- [x] PostgreSQL schema (125 tables)
- [x] Go backend (54 files, 140+ endpoints)
- [x] Nuxt 4 frontend (104 files, 55 pages)
- [x] AI Gateway with MCP (21 Go files)
- [x] Docker Compose with 11 services
- [x] Documentation (ERD, API Spec, Architecture)
- [x] CI/CD pipelines (lint, test, build, deploy)
- [x] Monitoring stack (Prometheus, Grafana, Jaeger)

### Phase 2 — Stabilization (In Progress)
- [ ] Data seeder scripts (sample school data)
- [ ] Unit tests for all services
- [ ] Integration tests for API endpoints
- [ ] E2E tests for critical user journeys
- [ ] Performance benchmarking
- [ ] Load testing with k6
- [ ] Payment gateway live integration testing
- [ ] Mobile PWA testing across devices

### Phase 3 — AI Enhancement
- [ ] Tahfidz voice recognition (speech-to-text for Quran recitation)
- [ ] Automated essay grading with rubric matching
- [ ] Plagiarism detection for assignments
- [ ] Student performance prediction models
- [ ] Personalized learning path recommendations
- [ ] AI-powered parent communication assistant
- [ ] Quran recitation error detection (tajwid rules)

### Phase 4 — Scale
- [ ] gRPC service extraction (student, finance, academic microservices)
- [ ] Kafka/Redpanda event streaming
- [ ] Kubernetes deployment (K3s)
- [ ] Read replicas for analytics queries
- [ ] Redis Cluster for session management
- [ ] CDN for static assets and uploaded documents
- [ ] Multi-region deployment

### Phase 5 — Enterprise
- [ ] SSO integration (OIDC, SAML)
- [ ] Advanced analytics with data warehouse (ClickHouse)
- [ ] Custom report builder (drag-and-drop)
- [ ] Marketplace for school templates and extensions
- [ ] White-label deployment for education foundations (Yayasan)
- [ ] Blockchain-based certificate verification
- [ ] IoT integration (attendance kiosks, smart cards)

---

## Contributing

Contributions are welcome. Please follow the established conventions:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Follow the Clean Architecture patterns
4. Write tests for new functionality
5. Run the lint checks (`golangci-lint run` / `pnpm lint`)
6. Commit with conventional commit messages
7. Push and create a Pull Request

### Commit Convention

```
feat: add tahfidz progress analytics chart
fix: resolve invoice generation race condition
docs: update API documentation for payment endpoints
test: add integration tests for gradebook service
refactor: extract notification dispatch to event handler
```

---

## License

Proprietary. All rights reserved.

---

## Support

- Documentation: See `docs/` directory
- Architecture: See [`docs/architecture/ARCHITECTURE.md`](docs/architecture/ARCHITECTURE.md)
- API Reference: See [`docs/architecture/API_SPEC.md`](docs/architecture/API_SPEC.md)
- Database ERD: See [`docs/architecture/ERD.md`](docs/architecture/ERD.md)

---

**Built with Go, Nuxt 4, PostgreSQL, Redis, MinIO, Qdrant, and AI.**

```
                  بسم الله الرحمن الرحيم
        Islamic School ERP — AI-Native Enterprise Platform
```
"# nuxt-erp-school" 
