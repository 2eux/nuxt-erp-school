-- ============================================================================
-- Islamic School ERP - Rollback Migration
-- Migration: 000001_init_schema
-- Description: Drops all tables created in the up migration (reverse order)
-- ============================================================================

BEGIN;

-- Drop triggers first
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN
        SELECT trigger_name, event_object_table
        FROM information_schema.triggers
        WHERE trigger_schema = 'public'
          AND trigger_name LIKE 'trg_%_updated_at'
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS %I ON %I', r.trigger_name, r.event_object_table);
    END LOOP;
END
$$;

DROP FUNCTION IF EXISTS update_updated_at_column();

-- ============================================================================
-- Drop tables in reverse dependency order (children before parents)
-- ============================================================================

-- AI / Knowledge Base
DROP TABLE IF EXISTS ai_generation_history CASCADE;
DROP TABLE IF EXISTS ai_prompts CASCADE;
DROP TABLE IF EXISTS ai_messages CASCADE;
DROP TABLE IF EXISTS ai_conversations CASCADE;
DROP TABLE IF EXISTS knowledge_chunks CASCADE;
DROP TABLE IF EXISTS knowledge_documents CASCADE;

-- Graduation
DROP TABLE IF EXISTS graduation_candidates CASCADE;
DROP TABLE IF EXISTS graduation_batches CASCADE;

-- Admissions
DROP TABLE IF EXISTS admission_documents CASCADE;
DROP TABLE IF EXISTS admission_exams CASCADE;
DROP TABLE IF EXISTS admission_applicants CASCADE;
DROP TABLE IF EXISTS admission_batches CASCADE;

-- Mosque / Islamic Events
DROP TABLE IF EXISTS zakat_reports CASCADE;
DROP TABLE IF EXISTS ramadhan_programs CASCADE;
DROP TABLE IF EXISTS islamic_events CASCADE;
DROP TABLE IF EXISTS mosque_activities CASCADE;

-- Audit / Logs
DROP TABLE IF EXISTS login_logs CASCADE;
DROP TABLE IF EXISTS activity_logs CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;

-- Settings
DROP TABLE IF EXISTS system_settings CASCADE;
DROP TABLE IF EXISTS settings CASCADE;

-- Announcements
DROP TABLE IF EXISTS announcements CASCADE;

-- Meetings / Calendar
DROP TABLE IF EXISTS task_comments CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS task_boards CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS meeting_participants CASCADE;
DROP TABLE IF EXISTS meetings CASCADE;

-- Documents / Letters
DROP TABLE IF EXISTS approval_actions CASCADE;
DROP TABLE IF EXISTS approval_requests CASCADE;
DROP TABLE IF EXISTS approval_workflows CASCADE;
DROP TABLE IF EXISTS letter_templates CASCADE;
DROP TABLE IF EXISTS letters CASCADE;
DROP TABLE IF EXISTS document_versions CASCADE;
DROP TABLE IF EXISTS documents CASCADE;

-- Notifications
DROP TABLE IF EXISTS notification_preferences CASCADE;
DROP TABLE IF EXISTS notification_templates CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;

-- Canteen
DROP TABLE IF EXISTS canteen_order_items CASCADE;
DROP TABLE IF EXISTS canteen_orders CASCADE;
DROP TABLE IF EXISTS canteen_products CASCADE;

-- Dormitory
DROP TABLE IF EXISTS dormitory_attendance CASCADE;
DROP TABLE IF EXISTS dormitory_residents CASCADE;
DROP TABLE IF EXISTS dormitory_rooms CASCADE;
DROP TABLE IF EXISTS dormitories CASCADE;

-- Transportation
DROP TABLE IF EXISTS transport_attendance CASCADE;
DROP TABLE IF EXISTS transport_students CASCADE;
DROP TABLE IF EXISTS transport_routes CASCADE;

-- Counseling
DROP TABLE IF EXISTS counseling_sessions CASCADE;

-- Medical
DROP TABLE IF EXISTS student_health_info CASCADE;
DROP TABLE IF EXISTS immunization_records CASCADE;
DROP TABLE IF EXISTS medical_records CASCADE;

-- Library
DROP TABLE IF EXISTS book_reservations CASCADE;
DROP TABLE IF EXISTS book_borrowings CASCADE;
DROP TABLE IF EXISTS library_members CASCADE;
DROP TABLE IF EXISTS library_books CASCADE;

-- Inventory / Assets
DROP TABLE IF EXISTS procurement_items CASCADE;
DROP TABLE IF EXISTS procurements CASCADE;
DROP TABLE IF EXISTS inventory_transactions CASCADE;
DROP TABLE IF EXISTS inventory_items CASCADE;
DROP TABLE IF EXISTS asset_maintenance CASCADE;
DROP TABLE IF EXISTS assets CASCADE;

-- Finance
DROP TABLE IF EXISTS infaq_records CASCADE;
DROP TABLE IF EXISTS waqf_records CASCADE;
DROP TABLE IF EXISTS tax_reports CASCADE;
DROP TABLE IF EXISTS payroll_details CASCADE;
DROP TABLE IF EXISTS payroll_periods CASCADE;
DROP TABLE IF EXISTS journal_entries CASCADE;
DROP TABLE IF EXISTS journals CASCADE;
DROP TABLE IF EXISTS general_ledger CASCADE;
DROP TABLE IF EXISTS budget_items CASCADE;
DROP TABLE IF EXISTS budget_plans CASCADE;
DROP TABLE IF EXISTS cash_transactions CASCADE;
DROP TABLE IF EXISTS bank_accounts CASCADE;
DROP TABLE IF EXISTS payment_gateway_requests CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS invoice_items CASCADE;
DROP TABLE IF EXISTS invoices CASCADE;
DROP TABLE IF EXISTS fee_assignments CASCADE;
DROP TABLE IF EXISTS fee_types CASCADE;

-- Islamic Modules
DROP TABLE IF EXISTS quranic_competencies CASCADE;
DROP TABLE IF EXISTS halaqah_members CASCADE;
DROP TABLE IF EXISTS halaqah_groups CASCADE;
DROP TABLE IF EXISTS islamic_character_notes CASCADE;
DROP TABLE IF EXISTS prayer_attendance CASCADE;
DROP TABLE IF EXISTS mutabaah_yaumiyah CASCADE;
DROP TABLE IF EXISTS tasmi_records CASCADE;
DROP TABLE IF EXISTS tahfidz_targets CASCADE;
DROP TABLE IF EXISTS tahfidz_progress CASCADE;
DROP TABLE IF EXISTS tahfidz_members CASCADE;
DROP TABLE IF EXISTS tahfidz_groups CASCADE;
DROP TABLE IF EXISTS tahfidz_programs CASCADE;

-- Academic
DROP TABLE IF EXISTS student_behavior CASCADE;
DROP TABLE IF EXISTS achievements CASCADE;
DROP TABLE IF EXISTS extracurricular_members CASCADE;
DROP TABLE IF EXISTS extracurriculars CASCADE;
DROP TABLE IF EXISTS teaching_journals CASCADE;
DROP TABLE IF EXISTS lesson_plans CASCADE;
DROP TABLE IF EXISTS assignment_submissions CASCADE;
DROP TABLE IF EXISTS assignments CASCADE;
DROP TABLE IF EXISTS report_cards CASCADE;
DROP TABLE IF EXISTS gradebooks CASCADE;
DROP TABLE IF EXISTS exam_results CASCADE;
DROP TABLE IF EXISTS exams CASCADE;
DROP TABLE IF EXISTS exam_types CASCADE;
DROP TABLE IF EXISTS attendances CASCADE;
DROP TABLE IF EXISTS schedules CASCADE;

-- Teachers / Employees
DROP TABLE IF EXISTS employee_performance CASCADE;
DROP TABLE IF EXISTS employee_trainings CASCADE;
DROP TABLE IF EXISTS leave_requests CASCADE;
DROP TABLE IF EXISTS employee_attendances CASCADE;
DROP TABLE IF EXISTS employees CASCADE;
DROP TABLE IF EXISTS teacher_subjects CASCADE;
DROP TABLE IF EXISTS teachers CASCADE;

-- Students
DROP TABLE IF EXISTS student_documents CASCADE;
DROP TABLE IF EXISTS student_parents CASCADE;
DROP TABLE IF EXISTS students CASCADE;

-- Users / RBAC
DROP TABLE IF EXISTS password_resets CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS user_profiles CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- Core / Master
DROP TABLE IF EXISTS curriculum CASCADE;
DROP TABLE IF EXISTS subjects CASCADE;
DROP TABLE IF EXISTS classes CASCADE;
DROP TABLE IF EXISTS grades CASCADE;
DROP TABLE IF EXISTS semesters CASCADE;
DROP TABLE IF EXISTS academic_years CASCADE;
DROP TABLE IF EXISTS schools CASCADE;

-- Extensions (optional - keep for safety)
-- DROP EXTENSION IF EXISTS pg_trgm;
-- DROP EXTENSION IF EXISTS pgcrypto;
-- DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;
