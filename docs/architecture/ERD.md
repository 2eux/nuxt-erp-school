# Islamic School ERP - Entity Relationship Diagram

## Overview

This document describes the complete database schema for the Islamic School ERP system. The system supports multi-tenancy (multiple schools), RBAC, soft deletes, audit trails, and comprehensive coverage of academic, financial, Islamic, and administrative domains.

---

## Entity Relationship Diagrams

### Core / Master Data

```
┌──────────────┐       ┌────────────────┐       ┌──────────────┐
│   schools    │──1:N──│ academic_years │──1:N──│  semesters   │
└──────────────┘       └────────────────┘       └──────────────┘
       │                                               │
       │ 1:N                                  ┌────────┘
       ▼                                      ▼
┌──────────┐   ┌──────────┐   ┌──────────┐  ┌──────────┐
│  grades  │   │ subjects │   │ curriculum│  │ subjects │
└──────────┘   └──────────┘   └──────────┘  └──────────┘
    │ 1:N                                        │
    ▼                          ┌─────────────────┘
┌──────────┐                   │
│  classes │───────────────────┘
└──────────┘
```

### Users / RBAC

```
┌─────────┐    ┌──────────────┐    ┌────────────────┐
│  roles  │─── │role_permiss. │────│  permissions   │
└─────────┘    └──────────────┘    └────────────────┘
     │
     │ 1:N
     ▼
┌─────────┐───1:1───┌───────────────┐
│  users  │         │ user_profiles │
└─────────┘         └───────────────┘
     │ 1:N
     ├─── user_sessions
     ├─── password_resets
     ├─── notifications
     └─── activity_logs / login_logs
```

### Students

```
┌──────────┐    ┌──────────────────┐    ┌────────────────────┐
│  users   │─── │    students      │─── │  student_parents   │
└──────────┘    └──────────────────┘    └────────────────────┘
                       │ 1:N
                       ├── student_documents
                       ├── achievements
                       ├── student_behavior
                       └── student_health_info
```

### Teachers / Employees

```
┌──────────┐    ┌───────────┐    ┌──────────────────┐
│  users   │─── │ teachers  │─── │ teacher_subjects │─── subjects / classes
└──────────┘    └───────────┘    └──────────────────┘
     │
     ├── employees ─── employee_attendances / leave_requests / trainings
     └── payroll_details
```

### Academic

```
┌───────────┐
│ schedules │─── classes / subjects / teachers
└───────────┘
      │
      ▼
┌──────────────┐
│ attendances  │─── students
└──────────────┘

┌─────────┐    ┌────────┐    ┌────────────────┐
│exams_type│───│ exams  │───│  exam_results  │─── students
└─────────┘    └────────┘    └────────────────┘

┌──────────────┐    ┌──────────────┐    ┌───────────────┐
│ gradebooks   │─── │ report_cards │─── │  assignments  │─── submissions
└──────────────┘    └──────────────┘    └───────────────┘
```

### Islamic Modules

```
┌──────────────────┐    ┌─────────────────┐
│ tahfidz_programs │─── │ tahfidz_groups  │─── teachers
└──────────────────┘    └─────────────────┘
                              │
                    ┌─────────┴──────────┐
                    ▼                    ▼
          ┌──────────────┐    ┌───────────────────┐
          │tahfidz_members│─── │ tahfidz_progress  │
          └──────────────┘    └───────────────────┘
                    │
                    ├── tahfidz_targets
                    ├── tasmi_records
                    └── quranic_competencies

┌─────────────────────┐    ┌──────────────────┐    ┌──────────────────────┐
│ mutabaah_yaumiyah   │    │ prayer_attendance │    │ islamic_char_notes  │
└─────────────────────┘    └──────────────────┘    └──────────────────────┘

┌────────────────┐    ┌─────────────────┐
│halaqah_groups  │─── │ halaqah_members │
└────────────────┘    └─────────────────┘
```

### Finance

```
┌────────────┐    ┌─────────────────┐    ┌───────────┐
│ fee_types  │─── │ fee_assignments │─── │ invoices  │─── invoice_items
└────────────┘    └─────────────────┘    └───────────┘
                                                │
                                                ▼
                                          ┌──────────┐
                                          │ payments │─── payment_gateway_requests
                                          └──────────┘

bank_accounts    cash_transactions    general_ledger    journals──journal_entries

budget_plans──budget_items    payroll_periods──payroll_details    tax_reports
```

### Inventory / Assets / Library / Medical / Counseling

```
assets──asset_maintenance
inventory_items──inventory_transactions
procurements──procurement_items

library_books──book_borrowings──library_members
            └──book_reservations──library_members

medical_records / immunization_records / student_health_info
counseling_sessions──students
```

### Operational Modules

```
transport_routes──transport_students──transport_attendance
dormitories──dormitory_rooms──dormitory_residents──dormitory_attendance
canteen_products──canteen_orders──canteen_order_items
```

### Communication & Workflow

```
notifications / notification_templates / notification_preferences
documents──document_versions
letters──letter_templates
approval_workflows──approval_requests──approval_actions
meetings──meeting_participants
events / task_boards──tasks──task_comments
announcements
```

### System

```
settings / system_settings
audit_logs / activity_logs / login_logs
```

### Admissions & Graduation

```
admission_batches──admission_applicants──admission_exams / admission_documents
graduation_batches──graduation_candidates
```

### AI / Knowledge Base

```
knowledge_documents──knowledge_chunks
ai_conversations──ai_messages
ai_prompts──ai_generation_history
```

---

## Naming Conventions

| Convention       | Example                  |
|-----------------|--------------------------|
| Table names     | `academic_years` (plural, snake_case) |
| Primary keys    | `id` (UUID, all tables)  |
| Foreign keys    | `school_id`, `student_id`|
| Timestamps      | `created_at`, `updated_at` (TIMESTAMPTZ) |
| Soft delete     | `deleted_at` (nullable TIMESTAMPTZ) |
| Audit columns   | `created_by`, `updated_by` (nullable UUID) |
| Status enums    | VARCHAR with CHECK constraints |
| Flexible data   | `jsonb` for settings, metadata, etc. |

## Key Design Decisions

1. **UUID Primary Keys** - All tables use `UUID DEFAULT gen_random_uuid()` for distributed safety
2. **Multi-Tenancy** - All tenant tables include `school_id` FK to `schools(id)` with NOT NULL
3. **Soft Delete** - Non-immutable tables use `deleted_at TIMESTAMPTZ` for logical deletion
4. **Audit Trail** - `created_by`, `updated_by` reference `users(id)`, with `created_at`, `updated_at`
5. **Full-Text Search** - GIN indexes on text columns using `to_tsvector('simple', column)`
6. **JSONB** - Settings, metadata, permissions stored as `jsonb` for schema flexibility
7. **Enum via CHECK** - Status fields use `VARCHAR` with CHECK constraints (not ENUM type)
8. **Foreign Keys** - CASCADE where logical, SET NULL where reference may be removed, RESTRICT for critical
9. **Indexes** - FK columns, status columns, date-range columns, and search columns all indexed
10. **Comments** - Every table and column has `COMMENT ON` for documentation

## Index Strategy

| Index Type              | Use Case                          |
|------------------------|-----------------------------------|
| B-tree (default)       | PKs, FKs, status, dates, lookups  |
| GIN (jsonb)            | JSONB settings/metadata columns   |
| GIN (full-text)        | name, title, description columns  |
| Unique                 | code fields, composite business keys |
| Partial                | `WHERE deleted_at IS NULL` for active records |

## Relationships Summary Table

| Module         | Tables | Major Relationships                                      |
|---------------|--------|----------------------------------------------------------|
| Core/Master    | 7      | schools → academic_years → semesters → grades → classes  |
| Users/RBAC     | 7      | roles → permissions, users → profiles/sessions            |
| Students       | 5      | users → students → parents/documents                     |
| Teachers       | 7      | users → teachers → subjects, employees → attendances     |
| Academic       | 14     | schedules → attendances, exams → results, gradebooks → reports |
| Islamic        | 12     | tahfidz → groups → members → progress, halaqah           |
| Finance        | 18     | fees → invoices → payments, journal → entries, payroll   |
| Inventory      | 5      | assets → maintenance, inventory → transactions           |
| Library        | 4      | books → borrowings/reservations                          |
| Medical        | 3      | records, immunization, health info                       |
| Counseling     | 1      | counseling_sessions                                      |
| Transportation | 3      | routes → students → attendance                           |
| Dormitory      | 4      | dormitories → rooms → residents → attendance             |
| Canteen        | 3      | products → orders → items                                |
| Notifications  | 3      | notifications, templates, preferences                    |
| Documents      | 8      | documents → versions, letters, workflow → approvals       |
| Meetings       | 5      | meetings → participants, events, tasks → comments        |
| Announcements  | 1      | announcements                                            |
| Settings       | 2      | school settings, system settings                         |
| Audit          | 3      | audit_logs, activity_logs, login_logs                    |
| Islamic Events | 6      | mosque activities, islamic events, ramadhan programs     |
| Admissions     | 4      | batches → applicants → exams → documents                 |
| Graduation     | 2      | batches → candidates                                     |
| AI/Knowledge   | 6      | knowledge docs → chunks, AI conversations → messages     |
| **Total**      | **~125** |                                                       |
