# Islamic School ERP - Architecture Documentation

## Table of Contents

1. [System Overview](#system-overview)
2. [Architecture Diagram](#architecture-diagram)
3. [Component Descriptions](#component-descriptions)
4. [Data Flow](#data-flow)
5. [Deployment Architecture](#deployment-architecture)
6. [Security Architecture](#security-architecture)
7. [AI Integration Architecture](#ai-integration-architecture)
8. [Multi-Tenancy Approach](#multi-tenancy-approach)
9. [Event-Driven Architecture](#event-driven-architecture)
10. [CQRS Patterns](#cqrs-patterns)
11. [gRPC Service Mesh Planning](#grpc-service-mesh-planning)
12. [Technology Stack](#technology-stack)
13. [Database Schema Overview](#database-schema-overview)

---

## System Overview

The Islamic School ERP is a comprehensive **multi-tenant** school management platform designed specifically for Islamic educational institutions. It covers academic, financial, administrative, Islamic studies (Tahfidz, Halaqah, Mutabaah), and operational domains with integrated AI capabilities for intelligent assistance, automated grading, student performance prediction, and knowledge retrieval (RAG).

**Key architectural principles:**
- **Domain-Driven Design (DDD)** with bounded contexts
- **CQRS** for read/write separation on complex queries
- **Event-Driven** communication between bounded contexts
- **Multi-tenancy** via `school_id` column with Row-Level Security (RLS)
- **API-first** design with versioned REST APIs
- **Observability** via OpenTelemetry (traces, metrics, logs)
- **Zero-trust security** model with JWT authentication + RBAC

---

## Architecture Diagram

```
                                        ┌──────────────────────────────────────┐
                                        │           Cloudflare CDN             │
                                        │    (DDoS Protection, WAF, Cache)      │
                                        └──────────────┬───────────────────────┘
                                                       │
                                                       ▼
┌──────────────┐    ┌──────────────────────────────────────────────────────────────┐
│   Mobile PWA │───▶│                        NGINX REVERSE PROXY                     │
│  (Nuxt SSR)  │    │  ┌─────────────┐  ┌──────────────┐  ┌─────────────────────┐  │
└──────────────┘    │  │ SSL Term.   │  │ Rate Limit   │  │ Security Headers    │  │
                    │  │ HTTP/2      │  │ WAF Rules    │  │ CSP, HSTS, CORS     │  │
┌──────────────┐    │  └─────────────┘  └──────────────┘  └─────────────────────┘  │
│  Admin Panel │───▶│                                                               │
│  (Nuxt SSR)  │    └─────┬──────────────────┬──────────────────┬───────────────────┘
└──────────────┘          │                  │                  │
                          ▼                  ▼                  ▼
               ┌──────────────────┐ ┌──────────────┐ ┌──────────────────┐
               │   FRONTEND       │ │  BACKEND     │ │  AI GATEWAY      │
               │  (Nuxt 4 SSR)   │ │  (Go/Gin)    │ │  (Go)            │
               │                  │ │              │ │                  │
               │ - SSR Rendering  │ │ - REST API   │ │ - Prompt Router  │
               │ - PWA Support    │ │ - Auth (JWT) │ │ - Rate Limiting  │
               │ - State (Pinia)  │ │ - RBAC       │ │ - Response Cache │
               │ - i18n (ID/AR/EN)│ │ - Validation │ │ - Model Fallback │
               │ - Cache (TanStack│ │ - File Mgmt  │ │ - RAG Pipeline  │
               │   Query)         │ │ - Background │ │ - Prompt Guard   │
               └──────────────────┘ │   Jobs       │ └────────┬─────────┘
                                    └──────┬───────┘          │
                                           │                  │
                      ┌────────────────────┼──────────────────┼───────────────────────┐
                      │                    │                  │                       │
                      ▼                    ▼                  ▼                       ▼
            ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐
            │  PostgreSQL  │    │    Redis      │    │   Qdrant     │    │  Ollama (Local)  │
            │     16       │    │      7        │    │  (Vectors)   │    │  LLM Server      │
            │              │    │              │    │              │    │                  │
            │ - CRUD Ops   │    │ - Sessions   │    │ - Embeddings │    │ - Local Models   │
            │ - Full-Text  │    │ - Cache      │    │ - Semantic   │    │ - llama3, gemma  │
            │ - JSONB      │    │ - Queue      │    │   Search     │    │ - nomic-embed    │
            │ - RLS        │    │ - Pub/Sub    │    │ - Knowledge  │    │                  │
            │ - Audit Logs │    │ - Rate Limit │    │   Base       │    │                  │
            └──────────────┘    └──────────────┘    └──────────────┘    └──────────────────┘
                      │
                      ▼
            ┌──────────────┐
            │    MinIO     │
            │  (S3 Object  │
            │   Storage)   │
            │              │
            │ - Documents  │
            │ - Photos     │
            │ - Reports    │
            │ - Backups    │
            └──────────────┘

┌───────────────────────────────────────────────────────────────────────────────────┐
│                            OBSERVABILITY STACK                                     │
│  ┌────────────┐  ┌──────────────┐  ┌─────────┐  ┌───────────────────────┐        │
│  │ Prometheus │  │   Grafana    │  │  Jaeger │  │ Loki (optional)       │        │
│  │  Metrics   │  │  Dashboards  │  │ Tracing │  │ Log Aggregation       │        │
│  └────────────┘  └──────────────┘  └─────────┘  └───────────────────────┘        │
└───────────────────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────────────────────┐
│                         EXTERNAL INTEGRATIONS                                      │
│  ┌────────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐     │
│  │ OpenAI API │  │ Gemini   │  │ Claude   │  │ Midtrans │  │ Firebase FCM │     │
│  │ (GPT-4o)   │  │          │  │          │  │ (Payment)│  │ (Push)       │     │
│  └────────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────────┘     │
│  ┌────────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐     │
│  │ Twilio SMS │  │ SMTP     │  │ WhatsApp │  │ Google   │  │ SSO (OIDC)   │     │
│  │ (SMS OTP)  │  │ (Email)  │  │ Business │  │ Calendar │  │              │     │
│  └────────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────────┘     │
└───────────────────────────────────────────────────────────────────────────────────┘
```

---

## Component Descriptions

### 1. Frontend (Nuxt 4)

| Aspect | Detail |
|--------|--------|
| Framework | Nuxt 4 (Vue 3.5) with Nitro server |
| Rendering | SSR with client-side hydration |
| State | Pinia stores + TanStack Query for server state |
| i18n | @nuxtjs/i18n (ID, AR, EN) |
| UI | Tailwind CSS v4 + @nuxt/ui components |
| Auth | JWT stored in httpOnly cookie + refresh token rotation |
| PWA | Offline support, push notifications |
| Charts | ApexCharts + vue3-apexcharts |
| Rich Text | TipTap editor, Vue Quill |
| Validation | Zod schemas with vee-validate |
| API Client | ofetch with auto-retry, interceptor for auth |

### 2. Backend (Go / Gin)

| Aspect | Detail |
|--------|--------|
| Framework | Gin v1.10 (HTTP router) |
| Database | pgx/v5 (direct PostgreSQL driver) |
| SQL Builder | sqlx (jmoiron/sqlx) |
| Migrations | golang-migrate (embedded SQL) |
| Auth | JWT (golang-jwt/jwt) + RBAC middleware |
| Validation | go-playground/validator |
| Config | Viper (env + file) |
| Logging | Zap (structured JSON logging) |
| Docs | Swagger (swaggo) auto-generated |
| Storage | MinIO SDK (S3-compatible) |
| AI | Qdrant Go client + custom RAG pipeline |
| Email | gomail/v2 (SMTP) |
| CQRS | Read models in separate query handlers |
| Background | Asynq (Redis-backed job queue) for async tasks |

### 3. AI Gateway

| Aspect | Detail |
|--------|--------|
| Purpose | Centralized AI request routing, caching, rate limiting |
| Routing | Provider-aware routing (OpenAI → Gemini → Claude fallback) |
| Caching | Redis-backed response cache with semantic similarity dedup |
| Guard | Prompt injection detection, content filtering |
| RAG | Knowledge document chunking + Qdrant vector store |
| Streaming | SSE stream passthrough to frontend |
| Analytics | Token usage tracking, cost monitoring |

### 4. PostgreSQL 16

- **125+ tables** across 25+ bounded contexts
- Features: JSONB, Full-Text Search (GIN), Row-Level Security, Partitioning
- Extensions: `pg_trgm`, `uuid-ossp`, `pgcrypto`, `pg_stat_statements`
- Index strategy: B-tree for lookups, GIN for FTS/JSONB, partial indexes for active records
- Backup: `pg_dump` nightly + WAL archiving for PITR

### 5. Redis 7

- Session store (hash)
- Cache layer (LRU eviction, 256MB limit)
- Rate limiting (sliding window)
- Job queue backend (Asynq)
- Pub/Sub for real-time notifications
- AI response cache (semantic dedup)

### 6. MinIO

- S3-compatible object store for:
  - Student/teacher profile photos
  - Document uploads (assignments, letters, reports)
  - System backups
  - Generated reports (PDF)
- Bucket policy: private by default, signed URLs for sharing

### 7. Qdrant

- Vector database for semantic search and RAG
- Collections:
  - `knowledge_base`: Islamic curriculum, textbooks, references
  - `student_profiles`: Learning patterns for personalized recommendations
  - `exam_questions`: Semantically similar question detection
- Embedding model: `text-embedding-3-small` (OpenAI) or local `nomic-embed-text` (Ollama)

---

## Data Flow

### Authentication Flow

```
User → Frontend (Login Form)
  → POST /api/v1/auth/login (email, password)
  → Backend validates credentials (bcrypt compare)
  → Backend generates JWT access token (15min) + refresh token (7d)
  → Backend creates session record in Redis (user_id:school_id:roles)
  → Frontend stores tokens in httpOnly cookie + memory
  → Subsequent requests: Authorization: Bearer <access_token>
  → When expired: POST /api/v1/auth/refresh (refresh_token)
  → Backend validates refresh token, issues new pair
```

### CRUD Operation Flow (CQRS Read Path)

```
User → Frontend (dashboard view)
  → GET /api/v1/students?page=1&limit=20&sort=name&filter[grade]=10
  → Backend validates JWT + RBAC (check: student:read)
  → Backend enforces school_id filter from JWT context
  → Backend query handler hits read-optimized PostgreSQL view (or Redis cache)
  → Returns paginated JSON response
  → Frontend renders with TanStack Query cache + optimistic updates
```

### CRUD Operation Flow (CQRS Write Path)

```
User → Frontend (create student form)
  → POST /api/v1/students (student data)
  → Backend validates JWT + RBAC (check: student:create)
  → Backend validates input (Zod schema on frontend, validator on backend)
  → Backend command handler inserts into PostgreSQL (with school_id)
  → Backend publishes event: student.created (Redis Pub/Sub)
  → Event consumers:
      - Notification service sends welcome email
      - Audit service logs the action
      - Cache invalidation clears student list cache
      - Search index updates (if separate)
  → Returns 201 Created with student resource
```

### AI RAG Query Flow

```
User → Frontend chat input
  → POST /api/v1/ai/chat (question, context: school_curriculum)
  → Backend → AI Gateway: POST /ai/chat/completions
  → AI Gateway checks Redis cache for similar query
  → (cache miss) → Converts query to embedding (text-embedding-3-small)
  → Queries Qdrant for relevant knowledge chunks (top_k=5, score_threshold=0.7)
  → Constructs RAG prompt: system_prompt + retrieved_context + user_question
  → Routes to AI provider (OpenAI GPT-4o or local Ollama)
  → Streams response chunks back (SSE)
  → Caches response in Redis (TTL: 1 hour)
  → Logs token usage for cost tracking
```

---

## Deployment Architecture

### Production Environment

```
┌──────────────────────────────────────────────────────────┐
│                    VPS / Bare Metal                      │
│  ┌────────────────────────────────────────────────────┐  │
│  │              Docker Compose Stack                   │  │
│  │                                                    │  │
│  │  nginx (reverse proxy)                             │  │
│  │  ├── backend (Go API) × 2 replicas                 │  │
│  │  ├── frontend (Nuxt SSR) × 2 replicas              │  │
│  │  ├── ai-gateway × 1 replica                        │  │
│  │  ├── postgres (primary)                            │  │
│  │  ├── redis                                         │  │
│  │  ├── minio                                         │  │
│  │  ├── qdrant                                        │  │
│  │  ├── prometheus                                    │  │
│  │  ├── grafana                                       │  │
│  │  ├── jaeger                                        │  │
│  │  └── ollama (optional, GPU)                        │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

### Scaling Strategy

| Stage | Approach |
|-------|----------|
| Phase 1 (Startup) | Single VPS, all containers on one host via Docker Compose |
| Phase 2 (Growth) | Database → managed PostgreSQL (Cloud SQL / RDS); CDN for static assets |
| Phase 3 (Scale) | Kubernetes (K3s or GKE); Horizontal Pod Autoscaler for backend/frontend |
| Phase 4 (Enterprise) | Multi-region deployment; read replicas; Redis Cluster; Kafka for events |

### CI/CD Pipeline

```
GitHub Push → GitHub Actions
  ├── Lint (golangci-lint, ESLint)
  ├── Test (go test -race, vitest)
  ├── Security Scan (gosec, trivy)
  ├── Docker Build & Push → ghcr.io
  └── Deploy to Staging
       ├── Deploy to Production (with manual approval)
       ├── Health Check
       └── Rollback (if failed)
```

---

## Security Architecture

### Defense in Depth Layers

```
Layer 1: Network
  ├── Cloudflare DDoS protection + WAF
  ├── Nginx rate limiting per IP/endpoint
  ├── TLS 1.2+ with strong ciphers
  └── Internal Docker network isolation

Layer 2: Authentication
  ├── JWT with RS256 (asymmetric) in production
  ├── Access token: 15min TTL
  ├── Refresh token: 7d TTL, rotation enabled
  ├── Device fingerprinting
  └── Brute force protection (Redis rate limit on /auth/*)

Layer 3: Authorization
  ├── RBAC with granular permissions (resource:action)
  ├── Multi-tenancy enforced at DB level (school_id check)
  ├── Row-Level Security (RLS) policies in PostgreSQL
  └── API scope validation per endpoint

Layer 4: Application
  ├── Input validation (Zod frontend, validator backend)
  ├── SQL injection prevention (parameterized queries)
  ├── XSS prevention (CSP headers + output encoding)
  ├── CSRF protection (SameSite cookies)
  └── File upload scanning (ClamAV for documents)

Layer 5: Data
  ├── Encryption at rest (LUKS/disk encryption)
  ├── Encryption in transit (TLS everywhere)
  ├── PII encryption (pgcrypto for sensitive fields)
  ├── Audit trail (all mutations logged)
  └── Backup encryption (GPG)

Layer 6: Monitoring
  ├── Prometheus alerts (CPU, memory, error rate)
  ├── Grafana dashboards (SLO tracking)
  ├── Jaeger distributed tracing
  └── Security incident response runbooks
```

### Security Headers (applied by Nginx)

| Header | Value | Purpose |
|--------|-------|---------|
| `Strict-Transport-Security` | `max-age=63072000; includeSubDomains; preload` | Enforce HTTPS |
| `Content-Security-Policy` | Strict CSP with nonce-based scripts | Prevent XSS |
| `X-Frame-Options` | `DENY` | Prevent clickjacking |
| `X-Content-Type-Options` | `nosniff` | Prevent MIME sniffing |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Control referrer info |
| `Permissions-Policy` | Camera/mic/geo disabled | Limit browser features |
| `Cross-Origin-Opener-Policy` | `same-origin-allow-popups` | Process isolation |

---

## AI Integration Architecture

### AI Capabilities

```
┌─────────────────────────────────────────────────────────────────────┐
│                        AI Feature Set                               │
├───────────────┬──────────────────┬─────────────────────────────────┤
│ Chatbot       │ Auto-Grading     │ Student Performance Prediction   │
│ (RAG + LLM)   │ (Essay scoring)  │ (ML model on grades data)        │
├───────────────┼──────────────────┼─────────────────────────────────┤
│ Quran Tutor   │ Lesson Plan Gen  │ Plagiarism Detection             │
│ (Tahfidz AI)  │ (Curriculum AI)  │ (Semantic similarity)            │
├───────────────┼──────────────────┼─────────────────────────────────┤
│ Smart Reports │ Recommendation   │ Attendance Anomaly Detection      │
│ (NLP summary) │ (Adaptive learn) │ (Statistical analysis)            │
└───────────────┴──────────────────┴─────────────────────────────────┘
```

### AI Gateway Architecture

```
                       ┌─────────────────────┐
                       │    AI Gateway        │
                       │                     │
Request ──────────────▶│ ┌─────────────────┐ │
                       │ │ Prompt Router   │ │──▶ OpenAI
                       │ │                 │ │──▶ Gemini (fallback)
                       │ │ Selects best    │ │──▶ Claude (fallback)
                       │ │ provider/model  │ │──▶ Ollama (local)
                       │ └─────────────────┘ │
                       │                     │
                       │ ┌─────────────────┐ │
                       │ │ Response Cache  │ │──▶ Redis
                       │ │ (semantic dup)  │ │    (TTL: 1h)
                       │ └─────────────────┘ │
                       │                     │
                       │ ┌─────────────────┐ │
                       │ │ RAG Pipeline    │ │──▶ Qdrant
                       │ │ - Query embed   │ │    (top_k=5)
                       │ │ - Vector search │ │
                       │ │ - Context build │ │
                       │ └─────────────────┘ │
                       │                     │
                       │ ┌─────────────────┐ │
                       │ │ Prompt Guard    │ │
                       │ │ - Injection det │ │
                       │ │ - Content filter│ │
                       │ │ - Hallucination │ │
                       │ │   check         │ │
                       │ └─────────────────┘ │
                       │                     │
                       │ ┌─────────────────┐ │
                       │ │ Usage Tracker   │ │
                       │ │ - Token count   │ │──▶ PostgreSQL
                       │ │ - Cost calc     │ │    (analytics)
                       │ │ - Rate limit    │ │
                       │ └─────────────────┘ │
                       └─────────────────────┘
```

### Prompt Strategy

All AI interactions follow a structured prompt template:

```
[System] You are an Islamic School ERP assistant. You help with academic,
         administrative, and Islamic education tasks. Follow these guidelines:
         1. Always cite sources from the provided context
         2. Be respectful and culturally appropriate
         3. Do not generate harmful or inappropriate content
         4. For religious questions, refer to standard Islamic sources

[Context] (from Qdrant vector search results):
         - Document 1: "..."
         - Document 2: "..."

[Conversation History] (last 10 messages, summarized)

[User Question] {user_input}
```

---

## Multi-Tenancy Approach

### Design

This system uses a **shared database, shared schema** multi-tenancy model with **Row-Level Security (RLS)** in PostgreSQL.

```
┌─────────────────────────────────────────────────────────┐
│                  Single PostgreSQL Database              │
│                                                         │
│  school_id = uuid_1    school_id = uuid_2              │
│  ┌──────────────┐     ┌──────────────┐                 │
│  │ School A      │     │ School B      │                 │
│  │  - 50 students│     │  - 30 students│                 │
│  │  - 5 teachers │     │  - 3 teachers │                 │
│  │  - 20 classes │     │  - 10 classes │                 │
│  └──────────────┘     └──────────────┘                 │
│                                                         │
│  Row-Level Security:                                    │
│  CREATE POLICY school_isolation ON students             │
│    USING (school_id = current_setting('app.school_id'));│
└─────────────────────────────────────────────────────────┘
```

### Tenant Isolation Implementation

1. **App Layer**: JWT contains `school_id` claim; middleware sets PostgreSQL session variable
2. **Database Layer**: RLS policies on all tenant-scoped tables (125+ tables)
3. **Cache Layer**: Redis keys prefixed with `{school_id}:`
4. **Storage Layer**: MinIO bucket path: `{school_id}/documents/...`
5. **Super Admin**: Bypasses RLS for cross-school reporting

### Per-Tenant Configuration

Each school has its own settings row in the `school_settings` table:
- Academic year configuration
- Grading scales (A-F, 0-100, 4.0 scale)
- Fee structures
- Curriculum mapping
- UI branding (logo, colors)
- Feature flags

---

## Event-Driven Architecture

### Event Types

```
┌──────────────────────┬────────────────────────────────────┐
│ Event                │ Producer          │ Consumers      │
├──────────────────────┼──────────────────┼────────────────┤
│ student.created      │ Student Service   │ Notification,  │
│                      │                   │ Audit, Cache   │
│ student.updated      │ Student Service   │ Cache Inval.   │
│ invoice.paid         │ Finance Service   │ Notification,  │
│                      │                   │ Ledger         │
│ exam.graded          │ Exam Service      │ Report Gen,    │
│                      │                   │ Gradebook      │
│ attendance.marked    │ Attendance Svc    │ Analytics       │
│ tahfidz.progressed   │ Tahfidz Service   │ Notification,  │
│                      │                   │ Achievement    │
│ user.login           │ Auth Service      │ Login Log       │
│ document.uploaded    │ Document Service  │ Virus Scan     │
└──────────────────────┴──────────────────┼────────────────┘
```

### Event Bus

- **Development**: Redis Pub/Sub (simple, single-node)
- **Production**: Redis Streams with consumer groups → Migration path to Kafka/Redpanda for large-scale deployments

### Event Schema (CloudEvents Format)

```json
{
  "specversion": "1.0",
  "type": "erp.school.student.created",
  "source": "/api/v1/students",
  "subject": "student-123",
  "id": "evt-abc-123",
  "time": "2026-01-15T08:30:00Z",
  "datacontenttype": "application/json",
  "data": {
    "student_id": "uuid",
    "school_id": "uuid",
    "name": "Ahmad Hidayat",
    "grade_id": "uuid",
    "changes": {}
  }
}
```

---

## CQRS Patterns

### Overview

CQRS (Command Query Responsibility Segregation) is applied selectively to bounded contexts with high read complexity or reporting needs.

### Example: Gradebook Context

```
┌─────────────────────────────────────────────────────────────┐
│                    GRADEBOOK BOUNDED CONTEXT                 │
├─────────────────────────┬───────────────────────────────────┤
│ COMMAND SIDE (Write)    │ QUERY SIDE (Read)                 │
├─────────────────────────┼───────────────────────────────────┤
│ GradeCommandHandler     │ GradeQueryHandler                 │
│ - AddGrade(cmd)         │ - GetStudentGrades(student_id)    │
│ - UpdateGrade(cmd)      │ - GetClassAverage(class_id)       │
│ - CalculateGPA(cmd)     │ - GetGradeDistribution(grade_id)  │
│                         │ - GetReportCard(student_id, term) │
├─────────────────────────┼───────────────────────────────────┤
│ PostgreSQL (normalized) │ PostgreSQL View / Redis Cache     │
│ grades table            │ v_student_report_card (view)      │
│ grade_components table  │ student_gpa_cache (Redis)         │
└─────────────────────────┴───────────────────────────────────┘
```

### When to Use CQRS

| Use CQRS | Use Simple CRUD |
|----------|-----------------|
| Report cards with complex aggregations | Simple user profile update |
| Dashboard with multi-table calculations | Creating a single record |
| Analytics queries over millions of rows | Basic CRUD for settings |
| Finance ledger with balance calculations | Master data management |
| Attendance reports with summaries | Notification preference update |

---

## gRPC Service Mesh Planning

### Future Architecture (Phase 3+)

As the system scales, individual bounded contexts will be extracted into microservices communicating via gRPC:

```
                        ┌──────────────────┐
                        │   API Gateway    │
                        │  (REST → gRPC)   │
                        └────────┬─────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              │                  │                  │
    ┌─────────▼──────┐  ┌───────▼──────┐  ┌────────▼────────┐
    │ Student Service │  │Finance Svc   │  │ Academic Svc    │
    │ (gRPC server)   │  │(gRPC server) │  │ (gRPC server)   │
    └────────┬────────┘  └───────┬──────┘  └────────┬────────┘
             │                   │                  │
    ┌────────▼────────┐  ┌───────▼──────┐  ┌────────▼────────┐
    │ Student DB      │  │ Finance DB   │  │ Academic DB     │
    │ (separate schema)│  │(separate)    │  │ (separate)      │
    └─────────────────┘  └──────────────┘  └─────────────────┘
```

### gRPC Service Definitions (Example)

```protobuf
service StudentService {
  rpc GetStudent(GetStudentRequest) returns (Student);
  rpc ListStudents(ListStudentsRequest) returns (ListStudentsResponse);
  rpc CreateStudent(CreateStudentRequest) returns (Student);
  rpc GetStudentGrades(GetStudentGradesRequest) returns (StudentGradesResponse);
}

service FinanceService {
  rpc GetInvoice(GetInvoiceRequest) returns (Invoice);
  rpc ProcessPayment(ProcessPaymentRequest) returns (PaymentResult);
  rpc GetOutstandingBalance(GetBalanceRequest) returns (BalanceResponse);
}
```

### Service Mesh Components (Planned)

| Component | Technology |
|-----------|------------|
| Service Registry | Consul / Kubernetes DNS |
| Load Balancing | Envoy (sidecar) + gRPC client LB |
| mTLS | Istio / Linkerd |
| Circuit Breaking | Envoy circuit breaker |
| Distributed Tracing | OpenTelemetry → Jaeger |
| API Gateway | Envoy / custom Go gateway |

---

## Technology Stack

| Layer          | Technology                           | Version   |
|----------------|--------------------------------------|-----------|
| Frontend       | Nuxt 4 (Vue 3.5)                     | 4.x       |
| UI Framework   | Tailwind CSS v4 + @nuxt/ui           | v2.20     |
| State Mgmt     | Pinia + TanStack Query               | latest    |
| Backend        | Go + Gin Framework                   | 1.23      |
| Database       | PostgreSQL                           | 16.3      |
| Cache          | Redis                                | 7.4       |
| Object Storage | MinIO (S3-compatible)                | latest    |
| Vector DB      | Qdrant                               | 1.12      |
| AI/LLM         | OpenAI API / Ollama (local)          | -         |
| Reverse Proxy  | Nginx                                | 1.27      |
| Metrics        | Prometheus                           | 2.55      |
| Dashboards     | Grafana                              | 11.3      |
| Tracing        | Jaeger (OpenTelemetry)               | 1.63      |
| CI/CD          | GitHub Actions                       | -         |
| Container      | Docker + Docker Compose              | 27+       |
| Orchestration  | K3s / GKE (planned)                  | -         |
| IaC            | Terraform (planned)                  | -         |

---

## Database Schema Overview

The complete database consists of approximately **125 tables** organized into **25+ bounded contexts**. See `docs/architecture/ERD.md` for the full entity relationship diagram.

### Core Table Categories

| # | Module | Tables | Key Entities |
|---|--------|--------|--------------|
| 1 | Core/Master | 7 | schools, academic_years, semesters, grades, classes, subjects, curriculum |
| 2 | Users/RBAC | 7 | users, roles, permissions, role_permissions, user_profiles, user_sessions, password_resets |
| 3 | Students | 5 | students, student_parents, student_documents, achievements, student_behavior |
| 4 | Teachers | 7 | teachers, employees, teacher_subjects, employee_attendances, leave_requests |
| 5 | Academic | 14 | schedules, attendances, exams, exam_results, gradebooks, report_cards, assignments |
| 6 | Islamic | 12 | tahfidz_programs, halaqah_groups, mutabaah_yaumiyah, prayer_attendance, tasmi_records |
| 7 | Finance | 18 | fee_types, invoices, payments, journals, budget_plans, payroll |
| 8 | Inventory | 5 | assets, inventory_items, procurements |
| 9 | Library | 4 | library_books, book_borrowings, book_reservations |
| 10 | Medical | 3 | medical_records, immunization_records, student_health_info |
| 11 | Operational | 10 | transport_routes, dormitories, canteen_products |
| 12 | Communication | 8 | notifications, documents, letters, approval_workflows, meetings |
| 13 | System | 5 | settings, audit_logs, activity_logs, login_logs |
| 14 | Admissions | 4 | admission_batches, admission_applicants |
| 15 | AI/Knowledge | 6 | knowledge_documents, knowledge_chunks, ai_conversations |
| | **Total** | **~125** | |

### Naming Conventions

| Convention | Example |
|------------|---------|
| Table names | `academic_years` (plural, snake_case) |
| Primary keys | `id` (UUID, all tables) |
| Foreign keys | `school_id`, `student_id` |
| Timestamps | `created_at`, `updated_at` (TIMESTAMPTZ) |
| Soft delete | `deleted_at` (nullable TIMESTAMPTZ) |
| Audit columns | `created_by`, `updated_by` (nullable UUID) |
| Status enums | VARCHAR with CHECK constraints |
| Flexible data | `jsonb` for settings, metadata, etc. |
