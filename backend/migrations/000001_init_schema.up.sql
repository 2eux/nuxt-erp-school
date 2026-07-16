-- ============================================================================
-- Islamic School ERP - Initial Database Schema
-- Migration: 000001_init_schema
-- Description: Complete initial schema with all modules
-- ============================================================================

BEGIN;

-- ============================================================================
-- EXTENSIONS
-- ============================================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================================
-- 1. CORE / MASTER DATA
-- ============================================================================

-- 1.1 schools
CREATE TABLE schools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    address TEXT,
    phone VARCHAR(30),
    email VARCHAR(255),
    logo_url TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('islamic', 'integrated')),
    accreditation VARCHAR(10),
    founded_year INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    settings JSONB DEFAULT '{}',
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    locale VARCHAR(10) DEFAULT 'id',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE schools IS 'School entities - top-level tenant';
COMMENT ON COLUMN schools.code IS 'Unique school code identifier';
COMMENT ON COLUMN schools.type IS 'islamic = full islamic school, integrated = integrated islamic curriculum';
COMMENT ON COLUMN schools.settings IS 'School-level configuration: grading rules, attendance rules, branding, etc.';
CREATE INDEX idx_schools_status ON schools (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_schools_code ON schools (code);
CREATE INDEX idx_schools_created_at ON schools (created_at);

-- 1.2 academic_years
CREATE TABLE academic_years (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_academic_years_dates CHECK (end_date > start_date)
);
COMMENT ON TABLE academic_years IS 'Academic year periods (e.g. 2024/2025)';
CREATE INDEX idx_academic_years_school ON academic_years (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_academic_years_active ON academic_years (school_id, is_active);
CREATE UNIQUE INDEX idx_academic_years_active_unique ON academic_years (school_id) WHERE is_active = true AND deleted_at IS NULL;

-- 1.3 semesters
CREATE TABLE semesters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    type VARCHAR(10) NOT NULL CHECK (type IN ('odd', 'even')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_semesters_dates CHECK (end_date > start_date)
);
COMMENT ON TABLE semesters IS 'Semesters within an academic year';
CREATE INDEX idx_semesters_ay ON semesters (academic_year_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_semesters_active ON semesters (academic_year_id, is_active);

-- 1.4 grades
CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    level INTEGER NOT NULL CHECK (level >= 0),
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_grades_code UNIQUE (school_id, code)
);
COMMENT ON TABLE grades IS 'Grade levels (e.g. 1st Grade, 2nd Grade)';
COMMENT ON COLUMN grades.level IS 'Numerical ordering level';
CREATE INDEX idx_grades_school ON grades (school_id) WHERE deleted_at IS NULL;

-- 1.5 classes
CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    grade_id UUID NOT NULL REFERENCES grades(id) ON DELETE RESTRICT,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20),
    homeroom_teacher_id UUID,
    capacity INTEGER DEFAULT 30,
    room_number VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'archived')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_classes_code UNIQUE (school_id, code)
);
COMMENT ON TABLE classes IS 'Classes combining grade + academic year with a homeroom teacher';
CREATE INDEX idx_classes_school ON classes (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_classes_grade ON classes (grade_id);
CREATE INDEX idx_classes_ay ON classes (academic_year_id);
CREATE INDEX idx_classes_teacher ON classes (homeroom_teacher_id);

-- 1.6 subjects
CREATE TABLE subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(20) NOT NULL,
    category VARCHAR(30) NOT NULL CHECK (category IN ('general', 'islamic', 'quran', 'arabic')),
    description TEXT,
    grade_level INTEGER,
    credits DECIMAL(3,1),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_subjects_code UNIQUE (school_id, code)
);
COMMENT ON TABLE subjects IS 'Subjects/courses offered by school';
COMMENT ON COLUMN subjects.category IS 'general, islamic, quran, arabic';
COMMENT ON COLUMN subjects.credits IS 'Credit hours/weight';
CREATE INDEX idx_subjects_school ON subjects (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_subjects_category ON subjects (school_id, category);
CREATE INDEX idx_subjects_fts ON subjects USING GIN (to_tsvector('simple', coalesce(name, '') || ' ' || coalesce(description, '')));

-- 1.7 curriculum
CREATE TABLE curriculum (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE curriculum IS 'Curriculum plans linked to academic years';
CREATE INDEX idx_curriculum_school ON curriculum (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_curriculum_ay ON curriculum (academic_year_id);

-- ============================================================================
-- 2. USERS / RBAC
-- ============================================================================

-- 2.1 roles
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    permissions JSONB DEFAULT '[]',
    is_system BOOLEAN DEFAULT false,
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE roles IS 'User roles with JSONB permission list';
COMMENT ON COLUMN roles.is_system IS 'System roles cannot be modified by school admins';
COMMENT ON COLUMN roles.school_id IS 'NULL for global roles, set for school-specific roles';
CREATE INDEX idx_roles_school ON roles (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_permissions ON roles USING GIN (permissions);

-- 2.2 permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    module VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE permissions IS 'Granular permission definitions (module.action)';
CREATE INDEX idx_permissions_module ON permissions (module);
CREATE INDEX idx_permissions_slug ON permissions (slug);

-- 2.3 role_permissions
CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, permission_id)
);
COMMENT ON TABLE role_permissions IS 'Many-to-many role-permission assignment';
CREATE INDEX idx_role_permissions_role ON role_permissions (role_id);
CREATE INDEX idx_role_permissions_permission ON role_permissions (permission_id);

-- 2.4 users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
    role_id UUID REFERENCES roles(id) ON DELETE SET NULL,
    email VARCHAR(255),
    phone VARCHAR(30),
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login TIMESTAMPTZ,
    settings JSONB DEFAULT '{}',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT chk_users_contact CHECK (email IS NOT NULL OR phone IS NOT NULL)
);
COMMENT ON TABLE users IS 'Core user accounts with multi-school support';
CREATE INDEX idx_users_school ON users (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users (role_id);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_phone ON users (phone);
CREATE INDEX idx_users_active ON users (school_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_fts ON users USING GIN (to_tsvector('simple', coalesce(full_name, '') || ' ' || coalesce(email, '')));

-- 2.5 user_profiles
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    nik VARCHAR(30),
    nisn VARCHAR(20),
    birth_date DATE,
    birth_place VARCHAR(100),
    gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
    address TEXT,
    religion VARCHAR(20) DEFAULT 'Islam',
    blood_type VARCHAR(5),
    photo_url TEXT,
    joining_date DATE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE user_profiles IS 'Extended profile information for users';
COMMENT ON COLUMN user_profiles.nik IS 'National ID number (KTP)';
COMMENT ON COLUMN user_profiles.nisn IS 'National Student ID Number';
CREATE INDEX idx_user_profiles_school ON user_profiles (school_id);

-- 2.6 user_sessions
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    ip_address INET,
    user_agent TEXT,
    is_revoked BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE user_sessions IS 'Active user refresh token sessions';
CREATE INDEX idx_user_sessions_user ON user_sessions (user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions (refresh_token);
CREATE INDEX idx_user_sessions_expires ON user_sessions (expires_at) WHERE is_revoked = false;

-- 2.7 password_resets
CREATE TABLE password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE password_resets IS 'Password reset tokens';
CREATE INDEX idx_password_resets_user ON password_resets (user_id);
CREATE INDEX idx_password_resets_token ON password_resets (token);

-- ============================================================================
-- 3. STUDENTS
-- ============================================================================

-- 3.1 students
CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    nis VARCHAR(30) UNIQUE NOT NULL,
    nisn VARCHAR(20),
    class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    parent_id UUID REFERENCES users(id) ON DELETE SET NULL,
    enrollment_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'alumni', 'transferred', 'dropped')),
    admission_batch VARCHAR(50),
    nik VARCHAR(30),
    family_card_no VARCHAR(50),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE students IS 'Student enrollment records';
COMMENT ON COLUMN students.nis IS 'School Student ID Number';
COMMENT ON COLUMN students.nisn IS 'National Student ID Number';
COMMENT ON COLUMN students.family_card_no IS 'Kartu Keluarga number';
CREATE INDEX idx_students_school ON students (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_class ON students (class_id);
CREATE INDEX idx_students_status ON students (status);
CREATE INDEX idx_students_fts ON students USING GIN (to_tsvector('simple', coalesce(nis, '') || ' ' || coalesce(nisn, '') || ' ' || coalesce(nik, '')));

-- 3.2 student_parents
CREATE TABLE student_parents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    relationship VARCHAR(20) NOT NULL CHECK (relationship IN ('father', 'mother', 'guardian')),
    is_primary BOOLEAN DEFAULT false,
    occupation VARCHAR(100),
    income_level VARCHAR(20),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE student_parents IS 'Parent/guardian relationships for students';
CREATE INDEX idx_student_parents_student ON student_parents (student_id);
CREATE INDEX idx_student_parents_user ON student_parents (user_id);
CREATE UNIQUE INDEX idx_student_parents_primary ON student_parents (student_id, relationship) WHERE deleted_at IS NULL;

-- 3.3 student_documents
CREATE TABLE student_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL CHECK (document_type IN ('birth_certificate', 'family_card', 'ijazah', 'skhun', 'transfer_letter', 'photo', 'ktp_parent', 'health_certificate', 'other')),
    file_url TEXT NOT NULL,
    verified_at TIMESTAMPTZ,
    verified_by UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE student_documents IS 'Student supporting documents';
CREATE INDEX idx_student_documents_student ON student_documents (student_id) WHERE deleted_at IS NULL;

-- ============================================================================
-- 4. TEACHERS / EMPLOYEES
-- ============================================================================

-- 4.1 teachers
CREATE TABLE teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    nip VARCHAR(30) UNIQUE,
    nuptk VARCHAR(30),
    qualification VARCHAR(100),
    certification TEXT,
    specialization VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'retired', 'transferred')),
    join_date DATE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE teachers IS 'Teacher employment records';
COMMENT ON COLUMN teachers.nip IS 'Civil servant/teacher ID';
COMMENT ON COLUMN teachers.nuptk IS 'Unique Teacher and Education Personnel Number';
CREATE INDEX idx_teachers_school ON teachers (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_teachers_nip ON teachers (nip);
CREATE INDEX idx_teachers_nuptk ON teachers (nuptk);
CREATE INDEX idx_teachers_fts ON teachers USING GIN (to_tsvector('simple', coalesce(qualification, '') || ' ' || coalesce(certification, '') || ' ' || coalesce(specialization, '')));

-- 4.2 teacher_subjects
CREATE TABLE teacher_subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE teacher_subjects IS 'Subject assignment mapping for teachers';
CREATE INDEX idx_teacher_subjects_teacher ON teacher_subjects (teacher_id);
CREATE INDEX idx_teacher_subjects_subject ON teacher_subjects (subject_id);
CREATE INDEX idx_teacher_subjects_class ON teacher_subjects (class_id);
CREATE UNIQUE INDEX idx_teacher_subjects_unique ON teacher_subjects (teacher_id, subject_id, class_id) WHERE deleted_at IS NULL;

-- 4.3 employees
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    nik VARCHAR(30),
    nip VARCHAR(30),
    position VARCHAR(100) NOT NULL,
    department VARCHAR(100),
    employment_type VARCHAR(20) NOT NULL CHECK (employment_type IN ('permanent', 'contract', 'honorary')),
    join_date DATE,
    salary_grade VARCHAR(20),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE employees IS 'Non-teaching employee records';
CREATE INDEX idx_employees_school ON employees (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_department ON employees (school_id, department);

-- 4.4 employee_attendances
CREATE TABLE employee_attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    check_in TIME,
    check_out TIME,
    status VARCHAR(20) NOT NULL CHECK (status IN ('present', 'late', 'absent', 'sick', 'leave')),
    overtime_hours DECIMAL(4,2) DEFAULT 0,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_employee_attendance UNIQUE (employee_id, date)
);
COMMENT ON TABLE employee_attendances IS 'Daily attendance records for employees';
CREATE INDEX idx_employee_att_employee ON employee_attendances (employee_id);
CREATE INDEX idx_employee_att_date ON employee_attendances (date);
CREATE INDEX idx_employee_att_status ON employee_attendances (status);

-- 4.5 leave_requests
CREATE TABLE leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type VARCHAR(30) NOT NULL CHECK (leave_type IN ('annual', 'sick', 'maternity', 'family', 'study', 'personal', 'other')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    approved_by UUID REFERENCES users(id),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_leave_dates CHECK (end_date >= start_date)
);
COMMENT ON TABLE leave_requests IS 'Employee leave/absence requests';
CREATE INDEX idx_leave_requests_employee ON leave_requests (employee_id);
CREATE INDEX idx_leave_requests_status ON leave_requests (status);
CREATE INDEX idx_leave_requests_dates ON leave_requests (start_date, end_date);

-- 4.6 employee_trainings
CREATE TABLE employee_trainings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    training_name VARCHAR(255) NOT NULL,
    provider VARCHAR(200),
    start_date DATE NOT NULL,
    end_date DATE,
    certificate_url TEXT,
    status VARCHAR(20) DEFAULT 'completed' CHECK (status IN ('planned', 'in_progress', 'completed', 'cancelled')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_training_dates CHECK (end_date IS NULL OR end_date >= start_date)
);
COMMENT ON TABLE employee_trainings IS 'Training and development records';
CREATE INDEX idx_employee_trainings_employee ON employee_trainings (employee_id);

-- 4.7 employee_performance
CREATE TABLE employee_performance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    semester_id UUID NOT NULL REFERENCES semesters(id) ON DELETE RESTRICT,
    score DECIMAL(5,2) NOT NULL CHECK (score >= 0 AND score <= 100),
    evaluator_id UUID REFERENCES users(id),
    notes TEXT,
    evaluation_date DATE NOT NULL,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_employee_perf UNIQUE (employee_id, semester_id)
);
COMMENT ON TABLE employee_performance IS 'Employee performance evaluations per semester';
CREATE INDEX idx_employee_perf_employee ON employee_performance (employee_id);

-- ============================================================================
-- 5. ACADEMIC
-- ============================================================================

-- 5.1 schedules
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID REFERENCES teachers(id) ON DELETE SET NULL,
    day_of_week SMALLINT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    room VARCHAR(50),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_schedule_time CHECK (end_time > start_time)
);
COMMENT ON TABLE schedules IS 'Class schedules (time tables)';
COMMENT ON COLUMN schedules.day_of_week IS '0=Sunday, 1=Monday, ..., 6=Saturday';
CREATE INDEX idx_schedules_class ON schedules (class_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_schedules_teacher ON schedules (teacher_id);
CREATE INDEX idx_schedules_subject ON schedules (subject_id);
CREATE INDEX idx_schedules_day ON schedules (class_id, day_of_week);

-- 5.2 attendances
CREATE TABLE attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    schedule_id UUID REFERENCES schedules(id) ON DELETE SET NULL,
    date DATE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('present', 'late', 'absent', 'sick', 'leave', 'permission')),
    check_in_time TIME,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_attendance UNIQUE (student_id, schedule_id, date)
);
COMMENT ON TABLE attendances IS 'Student daily attendance records';
CREATE INDEX idx_attendances_student ON attendances (student_id);
CREATE INDEX idx_attendances_schedule ON attendances (schedule_id);
CREATE INDEX idx_attendances_date ON attendances (date);
CREATE INDEX idx_attendances_status ON attendances (status);

-- 5.3 exam_types
CREATE TABLE exam_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    weight_percentage DECIMAL(5,2) NOT NULL CHECK (weight_percentage >= 0 AND weight_percentage <= 100),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_exam_types_code UNIQUE (school_id, code)
);
COMMENT ON TABLE exam_types IS 'Types of exams (UTS, UAS, Daily Test, etc.)';
CREATE INDEX idx_exam_types_school ON exam_types (school_id) WHERE deleted_at IS NULL;

-- 5.4 exams
CREATE TABLE exams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    exam_type_id UUID NOT NULL REFERENCES exam_types(id) ON DELETE RESTRICT,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    exam_date DATE,
    duration_minutes INTEGER,
    total_score DECIMAL(6,2) DEFAULT 100,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE exams IS 'Examination schedules and details';
CREATE INDEX idx_exams_subject ON exams (subject_id);
CREATE INDEX idx_exams_type ON exams (exam_type_id);
CREATE INDEX idx_exams_class ON exams (class_id);
CREATE INDEX idx_exams_date ON exams (exam_date);

-- 5.5 exam_results
CREATE TABLE exam_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id UUID NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    score DECIMAL(6,2) CHECK (score >= 0),
    grade VARCHAR(5),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_exam_result UNIQUE (exam_id, student_id)
);
COMMENT ON TABLE exam_results IS 'Individual exam scores per student';
CREATE INDEX idx_exam_results_exam ON exam_results (exam_id);
CREATE INDEX idx_exam_results_student ON exam_results (student_id);

-- 5.6 gradebooks
CREATE TABLE gradebooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    semester_id UUID NOT NULL REFERENCES semesters(id) ON DELETE RESTRICT,
    daily_score DECIMAL(5,2),
    uts_score DECIMAL(5,2),
    uas_score DECIMAL(5,2),
    final_score DECIMAL(5,2),
    grade VARCHAR(5),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_gradebook UNIQUE (class_id, subject_id, student_id, semester_id)
);
COMMENT ON TABLE gradebooks IS 'Cumulative grade records per student per subject per semester';
CREATE INDEX idx_gradebooks_class ON gradebooks (class_id);
CREATE INDEX idx_gradebooks_student ON gradebooks (student_id);
CREATE INDEX idx_gradebooks_semester ON gradebooks (semester_id);

-- 5.7 report_cards
CREATE TABLE report_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT,
    semester_id UUID REFERENCES semesters(id) ON DELETE RESTRICT,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    file_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_report_card UNIQUE (student_id, academic_year_id, semester_id)
);
COMMENT ON TABLE report_cards IS 'Student report cards/rapor';
CREATE INDEX idx_report_cards_student ON report_cards (student_id);
CREATE INDEX idx_report_cards_ay ON report_cards (academic_year_id);

-- 5.8 assignments
CREATE TABLE assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    due_date TIMESTAMPTZ NOT NULL,
    max_score DECIMAL(5,2) DEFAULT 100,
    file_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('draft', 'active', 'closed')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE assignments IS 'Assignments/tasks given by teachers';
CREATE INDEX idx_assignments_subject ON assignments (subject_id);
CREATE INDEX idx_assignments_teacher ON assignments (teacher_id);
CREATE INDEX idx_assignments_due_date ON assignments (due_date);
CREATE INDEX idx_assignments_fts ON assignments USING GIN (to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(description, '')));

-- 5.9 assignment_submissions
CREATE TABLE assignment_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    content TEXT,
    file_url TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    score DECIMAL(5,2),
    feedback TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_submission UNIQUE (assignment_id, student_id)
);
COMMENT ON TABLE assignment_submissions IS 'Student submissions for assignments';
CREATE INDEX idx_submissions_assignment ON assignment_submissions (assignment_id);
CREATE INDEX idx_submissions_student ON assignment_submissions (student_id);

-- 5.10 lesson_plans
CREATE TABLE lesson_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    week INTEGER,
    topic VARCHAR(255) NOT NULL,
    objectives TEXT,
    materials TEXT,
    methods TEXT,
    media TEXT,
    activities TEXT,
    evaluation TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'reviewed', 'approved', 'executed')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE lesson_plans IS 'RPP - Lesson Plan documents';
CREATE INDEX idx_lesson_plans_teacher ON lesson_plans (teacher_id);
CREATE INDEX idx_lesson_plans_subject ON lesson_plans (subject_id);
CREATE INDEX idx_lesson_plans_class ON lesson_plans (class_id);

-- 5.11 teaching_journals
CREATE TABLE teaching_journals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    schedule_id UUID REFERENCES schedules(id) ON DELETE SET NULL,
    date DATE NOT NULL,
    topic VARCHAR(255),
    reflection TEXT,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE teaching_journals IS 'Teacher daily reflection journals';
CREATE INDEX idx_teaching_journals_teacher ON teaching_journals (teacher_id);
CREATE INDEX idx_teaching_journals_date ON teaching_journals (date);

-- 5.12 extracurriculars
CREATE TABLE extracurriculars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    coach VARCHAR(200),
    schedule TEXT,
    max_participants INTEGER,
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE extracurriculars IS 'Extracurricular activities offered by school';
CREATE INDEX idx_extracurriculars_school ON extracurriculars (school_id) WHERE deleted_at IS NULL;

-- 5.13 extracurricular_members
CREATE TABLE extracurricular_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    extracurricular_id UUID NOT NULL REFERENCES extracurriculars(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    role VARCHAR(30) DEFAULT 'member',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_ec_member UNIQUE (extracurricular_id, student_id)
);
COMMENT ON TABLE extracurricular_members IS 'Student membership in extracurricular activities';
CREATE INDEX idx_ec_members_ex ON extracurricular_members (extracurricular_id);
CREATE INDEX idx_ec_members_student ON extracurricular_members (student_id);

-- 5.14 achievements
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('academic', 'sport', 'art', 'tahfidz', 'tahfiz', 'quran', 'arabic', 'science', 'technology', 'debate', 'olympiad', 'other')),
    level VARCHAR(30) NOT NULL CHECK (level IN ('school', 'district', 'city', 'province', 'national', 'international')),
    organizer VARCHAR(255),
    date DATE NOT NULL,
    certificate_url TEXT,
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE achievements IS 'Student achievements and awards';
CREATE INDEX idx_achievements_student ON achievements (student_id);
CREATE INDEX idx_achievements_category ON achievements (category);
CREATE INDEX idx_achievements_level ON achievements (level);

-- 5.15 student_behavior
CREATE TABLE student_behavior (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    incident_date DATE NOT NULL,
    category VARCHAR(20) NOT NULL CHECK (category IN ('positive', 'negative')),
    description TEXT NOT NULL,
    action_taken TEXT,
    recorded_by UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE student_behavior IS 'Student behavior/character incident records';
CREATE INDEX idx_student_behavior_student ON student_behavior (student_id);
CREATE INDEX idx_student_behavior_date ON student_behavior (incident_date);

-- ============================================================================
-- 6. ISLAMIC MODULES
-- ============================================================================

-- 6.1 tahfidz_programs
CREATE TABLE tahfidz_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE SET NULL,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tahfidz_programs IS 'Tahfidz (Quran memorization) programs';
CREATE INDEX idx_tahfidz_programs_school ON tahfidz_programs (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tahfidz_programs_ay ON tahfidz_programs (academic_year_id);

-- 6.2 tahfidz_groups
CREATE TABLE tahfidz_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tahfidz_program_id UUID NOT NULL REFERENCES tahfidz_programs(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    target_juz_start INTEGER CHECK (target_juz_start >= 1 AND target_juz_start <= 30),
    target_juz_end INTEGER CHECK (target_juz_end >= 1 AND target_juz_end <= 30),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_juz_range CHECK (target_juz_end >= target_juz_start)
);
COMMENT ON TABLE tahfidz_groups IS 'Tahfidz study groups within a program';
CREATE INDEX idx_tahfidz_groups_program ON tahfidz_groups (tahfidz_program_id);
CREATE INDEX idx_tahfidz_groups_teacher ON tahfidz_groups (teacher_id);

-- 6.3 tahfidz_members
CREATE TABLE tahfidz_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tahfidz_group_id UUID NOT NULL REFERENCES tahfidz_groups(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_tahfidz_member UNIQUE (tahfidz_group_id, student_id)
);
COMMENT ON TABLE tahfidz_members IS 'Student membership in tahfidz groups';
CREATE INDEX idx_tahfidz_members_group ON tahfidz_members (tahfidz_group_id);
CREATE INDEX idx_tahfidz_members_student ON tahfidz_members (student_id);

-- 6.4 tahfidz_progress
CREATE TABLE tahfidz_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    surah_name VARCHAR(100) NOT NULL,
    ayat_start INTEGER NOT NULL CHECK (ayat_start > 0),
    ayat_end INTEGER NOT NULL CHECK (ayat_end >= ayat_start),
    juz INTEGER CHECK (juz >= 1 AND juz <= 30),
    page INTEGER,
    score SMALLINT CHECK (score >= 1 AND score <= 10),
    type VARCHAR(20) NOT NULL CHECK (type IN ('new_memorization', 'murojaah')),
    teacher_id UUID REFERENCES teachers(id) ON DELETE SET NULL,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tahfidz_progress IS 'Daily tahfidz progress tracking';
COMMENT ON COLUMN tahfidz_progress.type IS 'new_memorization = hafalan baru, murojaah = review';
CREATE INDEX idx_tahfidz_progress_student ON tahfidz_progress (student_id);
CREATE INDEX idx_tahfidz_progress_date ON tahfidz_progress (date);
CREATE INDEX idx_tahfidz_progress_teacher ON tahfidz_progress (teacher_id);
CREATE INDEX idx_tahfidz_progress_type ON tahfidz_progress (type);

-- 6.5 tahfidz_targets
CREATE TABLE tahfidz_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    target_juz INTEGER CHECK (target_juz >= 1 AND target_juz <= 30),
    target_surah VARCHAR(100),
    target_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'not_started' CHECK (status IN ('achieved', 'in_progress', 'not_started')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tahfidz_targets IS 'Tahfidz target milestones per student';
CREATE INDEX idx_tahfidz_targets_student ON tahfidz_targets (student_id);
CREATE INDEX idx_tahfidz_targets_status ON tahfidz_targets (status);

-- 6.6 tasmi_records
CREATE TABLE tasmi_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    juz_start INTEGER CHECK (juz_start >= 1 AND juz_start <= 30),
    juz_end INTEGER CHECK (juz_end >= 1 AND juz_end <= 30),
    total_pages INTEGER,
    examiner_id UUID REFERENCES users(id),
    score SMALLINT CHECK (score >= 1 AND score <= 10),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_tasmi_juz CHECK (juz_end >= juz_start)
);
COMMENT ON TABLE tasmi_records IS 'Tasmi (recitation exam) records';
CREATE INDEX idx_tasmi_records_student ON tasmi_records (student_id);
CREATE INDEX idx_tasmi_records_date ON tasmi_records (date);

-- 6.7 mutabaah_yaumiyah
CREATE TABLE mutabaah_yaumiyah (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    fajr BOOLEAN DEFAULT false,
    dhuhr BOOLEAN DEFAULT false,
    asr BOOLEAN DEFAULT false,
    maghrib BOOLEAN DEFAULT false,
    isha BOOLEAN DEFAULT false,
    sunnah_rawatib BOOLEAN DEFAULT false,
    dhuha BOOLEAN DEFAULT false,
    tahajjud BOOLEAN DEFAULT false,
    quran_pages DECIMAL(4,1) DEFAULT 0,
    dhikr BOOLEAN DEFAULT false,
    charity BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'verified', 'revision')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_mutabaah UNIQUE (student_id, date)
);
COMMENT ON TABLE mutabaah_yaumiyah IS 'Daily ibadah (worship) tracking - Mutabaah Yaumiyah';
CREATE INDEX idx_mutabaah_student ON mutabaah_yaumiyah (student_id);
CREATE INDEX idx_mutabaah_date ON mutabaah_yaumiyah (date);

-- 6.8 prayer_attendance
CREATE TABLE prayer_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    prayer_type VARCHAR(20) NOT NULL CHECK (prayer_type IN ('fajr', 'dhuhr', 'asr', 'maghrib', 'isha', 'dhuha', 'tahajjud')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('prayed', 'late', 'missed', 'excused')),
    location VARCHAR(100),
    verified_by UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE prayer_attendance IS 'Prayer attendance tracking per student';
CREATE INDEX idx_prayer_att_student ON prayer_attendance (student_id);
CREATE INDEX idx_prayer_att_date ON prayer_attendance (date);
CREATE INDEX idx_prayer_att_type ON prayer_attendance (prayer_type);
CREATE UNIQUE INDEX idx_prayer_att_unique ON prayer_attendance (student_id, date, prayer_type);

-- 6.9 islamic_character_notes
CREATE TABLE islamic_character_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    category VARCHAR(20) NOT NULL CHECK (category IN ('adab', 'akhlaq', 'ibadah')),
    observation TEXT NOT NULL,
    rating SMALLINT CHECK (rating >= 1 AND rating <= 5),
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE islamic_character_notes IS 'Islamic character/behavior observations';
COMMENT ON COLUMN islamic_character_notes.category IS 'adab = manners, akhlaq = morals, ibadah = worship';
CREATE INDEX idx_char_notes_student ON islamic_character_notes (student_id);
CREATE INDEX idx_char_notes_date ON islamic_character_notes (date);
CREATE INDEX idx_char_notes_teacher ON islamic_character_notes (teacher_id);

-- 6.10 halaqah_groups
CREATE TABLE halaqah_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    day VARCHAR(20) CHECK (day IN ('sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday')),
    start_time TIME,
    end_time TIME,
    room VARCHAR(50),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE halaqah_groups IS 'Halaqah (Islamic study circle) groups';
CREATE INDEX idx_halaqah_groups_school ON halaqah_groups (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_halaqah_groups_teacher ON halaqah_groups (teacher_id);

-- 6.11 halaqah_members
CREATE TABLE halaqah_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    halaqah_group_id UUID NOT NULL REFERENCES halaqah_groups(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_halaqah_member UNIQUE (halaqah_group_id, student_id)
);
COMMENT ON TABLE halaqah_members IS 'Student membership in halaqah groups';
CREATE INDEX idx_halaqah_members_group ON halaqah_members (halaqah_group_id);
CREATE INDEX idx_halaqah_members_student ON halaqah_members (student_id);

-- 6.12 quranic_competencies
CREATE TABLE quranic_competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    semester_id UUID NOT NULL REFERENCES semesters(id) ON DELETE RESTRICT,
    tajwid_score DECIMAL(5,2) CHECK (tajwid_score >= 0 AND tajwid_score <= 100),
    tilawah_score DECIMAL(5,2) CHECK (tilawah_score >= 0 AND tilawah_score <= 100),
    tahfidz_score DECIMAL(5,2) CHECK (tahfidz_score >= 0 AND tahfidz_score <= 100),
    tafsir_score DECIMAL(5,2) CHECK (tafsir_score >= 0 AND tafsir_score <= 100),
    overall_score DECIMAL(5,2) CHECK (overall_score >= 0 AND overall_score <= 100),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_quranic_comp UNIQUE (student_id, semester_id)
);
COMMENT ON TABLE quranic_competencies IS 'Quranic competency assessments per semester';
CREATE INDEX idx_quranic_comp_student ON quranic_competencies (student_id);
CREATE INDEX idx_quranic_comp_semester ON quranic_competencies (semester_id);

-- ============================================================================
-- 7. FINANCE
-- ============================================================================

-- 7.1 fee_types
CREATE TABLE fee_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(50) NOT NULL,
    category VARCHAR(30) NOT NULL CHECK (category IN ('spp', 'registration', 'development', 'uniform', 'book', 'exam', 'event', 'donation')),
    amount DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    is_recurring BOOLEAN DEFAULT false,
    recurrence_period VARCHAR(20) CHECK (recurrence_period IN ('monthly', 'semester', 'annual')),
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_fee_types_code UNIQUE (school_id, code)
);
COMMENT ON TABLE fee_types IS 'Fee type definitions';
CREATE INDEX idx_fee_types_school ON fee_types (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_fee_types_category ON fee_types (school_id, category);

-- 7.2 fee_assignments
CREATE TABLE fee_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fee_type_id UUID NOT NULL REFERENCES fee_types(id) ON DELETE RESTRICT,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    amount DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    discount DECIMAL(12,2) DEFAULT 0 CHECK (discount >= 0),
    total DECIMAL(12,2) GENERATED ALWAYS AS (amount - discount) STORED,
    due_date DATE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE SET NULL,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE fee_assignments IS 'Fee assignments linking fee types to specific students';
CREATE INDEX idx_fee_assignments_fee ON fee_assignments (fee_type_id);
CREATE INDEX idx_fee_assignments_student ON fee_assignments (student_id);
CREATE INDEX idx_fee_assignments_due ON fee_assignments (due_date);

-- 7.3 invoices
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID REFERENCES students(id) ON DELETE SET NULL,
    invoice_number VARCHAR(50) NOT NULL,
    invoice_date DATE NOT NULL,
    due_date DATE NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL CHECK (total_amount >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'paid', 'overdue', 'cancelled')),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_invoice_number UNIQUE (school_id, invoice_number)
);
COMMENT ON TABLE invoices IS 'Student/school invoices';
CREATE INDEX idx_invoices_school ON invoices (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoices_student ON invoices (student_id);
CREATE INDEX idx_invoices_status ON invoices (status);
CREATE INDEX idx_invoices_due ON invoices (due_date) WHERE status NOT IN ('paid', 'cancelled');

-- 7.4 invoice_items
CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    fee_type_id UUID REFERENCES fee_types(id) ON DELETE SET NULL,
    description VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL CHECK (unit_price >= 0),
    amount DECIMAL(12,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    discount DECIMAL(12,2) DEFAULT 0 CHECK (discount >= 0),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE invoice_items IS 'Line items within an invoice';
CREATE INDEX idx_invoice_items_invoice ON invoice_items (invoice_id);

-- 7.5 payments
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    student_id UUID REFERENCES students(id) ON DELETE SET NULL,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    payment_date DATE NOT NULL,
    payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN ('cash', 'transfer', 'va', 'ewallet')),
    reference_number VARCHAR(100),
    proof_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('verified', 'pending', 'rejected')),
    verified_by UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE payments IS 'Payment records against invoices';
CREATE INDEX idx_payments_invoice ON payments (invoice_id);
CREATE INDEX idx_payments_student ON payments (student_id);
CREATE INDEX idx_payments_date ON payments (payment_date);
CREATE INDEX idx_payments_status ON payments (status);
CREATE INDEX idx_payments_reference ON payments (reference_number);

-- 7.6 payment_gateway_requests
CREATE TABLE payment_gateway_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    gateway VARCHAR(30) NOT NULL CHECK (gateway IN ('va', 'bank_transfer', 'ewallet')),
    request_data JSONB,
    response_data JSONB,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'callback_received')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE payment_gateway_requests IS 'Payment gateway integration logs';
CREATE INDEX idx_pg_requests_payment ON payment_gateway_requests (payment_id);
CREATE INDEX idx_pg_requests_status ON payment_gateway_requests (status);

-- 7.7 bank_accounts
CREATE TABLE bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    bank_name VARCHAR(100) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    account_holder VARCHAR(200) NOT NULL,
    branch VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_bank_account UNIQUE (school_id, account_number)
);
COMMENT ON TABLE bank_accounts IS 'School bank account definitions';
CREATE INDEX idx_bank_accounts_school ON bank_accounts (school_id) WHERE deleted_at IS NULL;

-- 7.8 cash_transactions
CREATE TABLE cash_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    transaction_date DATE NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('cash_in', 'cash_out')),
    category VARCHAR(50) NOT NULL,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    description TEXT,
    reference VARCHAR(100),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE cash_transactions IS 'Cash inflow/outflow records';
CREATE INDEX idx_cash_tx_school ON cash_transactions (school_id);
CREATE INDEX idx_cash_tx_date ON cash_transactions (transaction_date);
CREATE INDEX idx_cash_tx_type ON cash_transactions (type);

-- 7.9 budget_plans
CREATE TABLE budget_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(50),
    planned_amount DECIMAL(14,2) NOT NULL CHECK (planned_amount >= 0),
    approved_amount DECIMAL(14,2),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'submitted', 'approved', 'rejected')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE budget_plans IS 'Budget plans per academic year';
CREATE INDEX idx_budget_plans_school ON budget_plans (school_id);
CREATE INDEX idx_budget_plans_ay ON budget_plans (academic_year_id);

-- 7.10 budget_items
CREATE TABLE budget_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_plan_id UUID NOT NULL REFERENCES budget_plans(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL CHECK (unit_price >= 0),
    amount DECIMAL(14,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE budget_items IS 'Line items within a budget plan';
CREATE INDEX idx_budget_items_plan ON budget_items (budget_plan_id);

-- 7.11 general_ledger
CREATE TABLE general_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    account_code VARCHAR(50) NOT NULL,
    account_name VARCHAR(200) NOT NULL,
    transaction_date DATE NOT NULL,
    debit DECIMAL(14,2) DEFAULT 0 CHECK (debit >= 0),
    credit DECIMAL(14,2) DEFAULT 0 CHECK (credit >= 0),
    balance DECIMAL(14,2) DEFAULT 0,
    reference VARCHAR(100),
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_ledger_debit_credit CHECK (debit = 0 OR credit = 0)
);
COMMENT ON TABLE general_ledger IS 'General ledger entries (COA-based)';
CREATE INDEX idx_gl_school ON general_ledger (school_id);
CREATE INDEX idx_gl_account ON general_ledger (account_code);
CREATE INDEX idx_gl_date ON general_ledger (transaction_date);

-- 7.12 journals
CREATE TABLE journals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    journal_date DATE NOT NULL,
    journal_number VARCHAR(50) NOT NULL,
    description TEXT,
    reference VARCHAR(100),
    total_debit DECIMAL(14,2) DEFAULT 0 CHECK (total_debit >= 0),
    total_credit DECIMAL(14,2) DEFAULT 0 CHECK (total_credit >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'posted', 'reversed')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_journal_number UNIQUE (school_id, journal_number)
);
COMMENT ON TABLE journals IS 'Journal entries (debit/credit pairs)';
CREATE INDEX idx_journals_school ON journals (school_id);
CREATE INDEX idx_journals_date ON journals (journal_date);

-- 7.13 journal_entries
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_id UUID NOT NULL REFERENCES journals(id) ON DELETE CASCADE,
    account_code VARCHAR(50) NOT NULL,
    account_name VARCHAR(200) NOT NULL,
    debit DECIMAL(14,2) DEFAULT 0 CHECK (debit >= 0),
    credit DECIMAL(14,2) DEFAULT 0 CHECK (credit >= 0),
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_journal_entry_dc CHECK (debit = 0 OR credit = 0)
);
COMMENT ON TABLE journal_entries IS 'Individual entries within a journal';
CREATE INDEX idx_journal_entries_journal ON journal_entries (journal_id);

-- 7.14 payroll_periods
CREATE TABLE payroll_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'calculated', 'approved', 'paid')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_payroll_dates CHECK (end_date > start_date)
);
COMMENT ON TABLE payroll_periods IS 'Payroll period definitions';
CREATE INDEX idx_payroll_periods_school ON payroll_periods (school_id);
CREATE INDEX idx_payroll_periods_status ON payroll_periods (status);

-- 7.15 payroll_details
CREATE TABLE payroll_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_period_id UUID NOT NULL REFERENCES payroll_periods(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    base_salary DECIMAL(12,2) NOT NULL CHECK (base_salary >= 0),
    allowances DECIMAL(12,2) DEFAULT 0 CHECK (allowances >= 0),
    deductions DECIMAL(12,2) DEFAULT 0 CHECK (deductions >= 0),
    tax DECIMAL(12,2) DEFAULT 0 CHECK (tax >= 0),
    bpjs_health DECIMAL(12,2) DEFAULT 0 CHECK (bpjs_health >= 0),
    bpjs_employment DECIMAL(12,2) DEFAULT 0 CHECK (bpjs_employment >= 0),
    net_salary DECIMAL(12,2) GENERATED ALWAYS AS (base_salary + allowances - deductions - tax - bpjs_health - bpjs_employment) STORED,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_payroll_employee UNIQUE (payroll_period_id, employee_id)
);
COMMENT ON TABLE payroll_details IS 'Individual employee payroll calculations';
CREATE INDEX idx_payroll_details_period ON payroll_details (payroll_period_id);
CREATE INDEX idx_payroll_details_employee ON payroll_details (employee_id);

-- 7.16 tax_reports
CREATE TABLE tax_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    report_type VARCHAR(50) NOT NULL,
    total_tax DECIMAL(14,2) NOT NULL CHECK (total_tax >= 0),
    file_url TEXT,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'filed', 'verified')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_tax_report_dates CHECK (period_end >= period_start)
);
COMMENT ON TABLE tax_reports IS 'Tax reporting records';
CREATE INDEX idx_tax_reports_school ON tax_reports (school_id);
CREATE INDEX idx_tax_reports_period ON tax_reports (period_start, period_end);

-- 7.17 waqf_records
CREATE TABLE waqf_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    waqif_name VARCHAR(200) NOT NULL,
    asset_type VARCHAR(50) NOT NULL,
    asset_description TEXT,
    estimated_value DECIMAL(14,2),
    certificate_url TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'active', 'released')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE waqf_records IS 'Waqf (endowment) records';

-- 7.18 infaq_records
CREATE TABLE infaq_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    donor_name VARCHAR(200),
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    payment_method VARCHAR(30),
    category VARCHAR(50),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE infaq_records IS 'Infaq (charity/donation) records';
CREATE INDEX idx_infaq_records_school ON infaq_records (school_id);
CREATE INDEX idx_infaq_records_date ON infaq_records (date);

-- ============================================================================
-- 8. INVENTORY / ASSETS
-- ============================================================================

-- 8.1 assets
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    purchase_date DATE,
    purchase_price DECIMAL(14,2),
    current_value DECIMAL(14,2),
    condition VARCHAR(30) CHECK (condition IN ('excellent', 'good', 'fair', 'poor', 'damaged')),
    location VARCHAR(100),
    responsible_person_id UUID REFERENCES users(id),
    depreciation_rate DECIMAL(5,2) CHECK (depreciation_rate >= 0 AND depreciation_rate <= 100),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'in_maintenance', 'disposed', 'lost')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_asset_code UNIQUE (school_id, code)
);
COMMENT ON TABLE assets IS 'School fixed assets';
CREATE INDEX idx_assets_school ON assets (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_assets_category ON assets (school_id, category);
CREATE INDEX idx_assets_fts ON assets USING GIN (to_tsvector('simple', coalesce(name, '') || ' ' || coalesce(code, '')));

-- 8.2 asset_maintenance
CREATE TABLE asset_maintenance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    maintenance_date DATE NOT NULL,
    description TEXT,
    cost DECIMAL(12,2) CHECK (cost >= 0),
    next_maintenance_date DATE,
    vendor VARCHAR(200),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE asset_maintenance IS 'Asset maintenance/service records';
CREATE INDEX idx_asset_maint_asset ON asset_maintenance (asset_id);
CREATE INDEX idx_asset_maint_date ON asset_maintenance (maintenance_date);

-- 8.3 inventory_items
CREATE TABLE inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    unit VARCHAR(30) NOT NULL DEFAULT 'pcs',
    stock DECIMAL(10,2) DEFAULT 0 CHECK (stock >= 0),
    min_stock DECIMAL(10,2) DEFAULT 0 CHECK (min_stock >= 0),
    unit_price DECIMAL(12,2),
    supplier_id UUID,
    location VARCHAR(100),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_inv_item_code UNIQUE (school_id, code)
);
COMMENT ON TABLE inventory_items IS 'Inventory/stock items';
CREATE INDEX idx_inv_items_school ON inventory_items (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_inv_items_category ON inventory_items (school_id, category);
CREATE INDEX idx_inv_items_low_stock ON inventory_items (school_id) WHERE stock <= min_stock AND deleted_at IS NULL;

-- 8.4 inventory_transactions
CREATE TABLE inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    transaction_date DATE NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('in', 'out', 'transfer')),
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2),
    total DECIMAL(12,2),
    reference VARCHAR(100),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE inventory_transactions IS 'Stock in/out/transfer records';
CREATE INDEX idx_inv_tx_item ON inventory_transactions (item_id);
CREATE INDEX idx_inv_tx_date ON inventory_transactions (transaction_date);
CREATE INDEX idx_inv_tx_type ON inventory_transactions (type);

-- 8.5 procurements
CREATE TABLE procurements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    request_date DATE NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id),
    department VARCHAR(100),
    description TEXT,
    estimated_cost DECIMAL(14,2) CHECK (estimated_cost >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'requested' CHECK (status IN ('requested', 'approved', 'ordered', 'received', 'cancelled')),
    approved_by UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE procurements IS 'Procurement/purchase requests';
CREATE INDEX idx_procurements_school ON procurements (school_id);
CREATE INDEX idx_procurements_status ON procurements (status);

-- 8.6 procurement_items
CREATE TABLE procurement_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    procurement_id UUID NOT NULL REFERENCES procurements(id) ON DELETE CASCADE,
    item_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    estimated_unit_price DECIMAL(12,2) CHECK (estimated_unit_price >= 0),
    actual_unit_price DECIMAL(12,2),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE procurement_items IS 'Line items within a procurement request';
CREATE INDEX idx_procurement_items_proc ON procurement_items (procurement_id);

-- ============================================================================
-- 9. LIBRARY
-- ============================================================================

-- 9.1 library_books
CREATE TABLE library_books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    isbn VARCHAR(20),
    title VARCHAR(500) NOT NULL,
    author VARCHAR(300),
    publisher VARCHAR(200),
    publish_year INTEGER,
    category VARCHAR(100),
    classification VARCHAR(50),
    language VARCHAR(50) DEFAULT 'Indonesia',
    total_copies INTEGER NOT NULL DEFAULT 1 CHECK (total_copies >= 0),
    available_copies INTEGER NOT NULL DEFAULT 1 CHECK (available_copies >= 0),
    shelf_location VARCHAR(100),
    cover_url TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE library_books IS 'Library book catalog';
CREATE INDEX idx_library_books_school ON library_books (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_library_books_isbn ON library_books (isbn);
CREATE INDEX idx_library_books_category ON library_books (school_id, category);
CREATE INDEX idx_library_books_fts ON library_books USING GIN (to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(author, '') || ' ' || coalesce(publisher, '') || ' ' || coalesce(isbn, '')));

-- 9.2 library_members
CREATE TABLE library_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    member_id VARCHAR(50) NOT NULL,
    membership_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expiry_date DATE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'expired')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_lib_member_id UNIQUE (school_id, member_id)
);
COMMENT ON TABLE library_members IS 'Library membership records';
CREATE INDEX idx_lib_members_user ON library_members (user_id);
CREATE INDEX idx_lib_members_school ON library_members (school_id) WHERE deleted_at IS NULL;

-- 9.3 book_borrowings
CREATE TABLE book_borrowings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID NOT NULL REFERENCES library_books(id) ON DELETE RESTRICT,
    member_id UUID NOT NULL REFERENCES library_members(id) ON DELETE RESTRICT,
    borrow_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL,
    return_date DATE,
    late_fee DECIMAL(10,2) DEFAULT 0 CHECK (late_fee >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'borrowed' CHECK (status IN ('borrowed', 'returned', 'lost', 'damaged')),
    issued_by UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_borrow_dates CHECK (due_date >= borrow_date)
);
COMMENT ON TABLE book_borrowings IS 'Book borrowing records';
CREATE INDEX idx_borrowings_book ON book_borrowings (book_id);
CREATE INDEX idx_borrowings_member ON book_borrowings (member_id);
CREATE INDEX idx_borrowings_status ON book_borrowings (status);
CREATE INDEX idx_borrowings_due ON book_borrowings (due_date) WHERE status = 'borrowed';

-- 9.4 book_reservations
CREATE TABLE book_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID NOT NULL REFERENCES library_books(id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES library_members(id) ON DELETE CASCADE,
    reserve_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expire_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'fulfilled', 'cancelled', 'expired')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_reserve_dates CHECK (expire_date >= reserve_date)
);
COMMENT ON TABLE book_reservations IS 'Book reservation records';
CREATE INDEX idx_reservations_book ON book_reservations (book_id);
CREATE INDEX idx_reservations_member ON book_reservations (member_id);

-- ============================================================================
-- 10. MEDICAL
-- ============================================================================

-- 10.1 medical_records
CREATE TABLE medical_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    visit_date DATE NOT NULL,
    complaint TEXT,
    diagnosis TEXT,
    treatment TEXT,
    medication TEXT,
    doctor_name VARCHAR(200),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE medical_records IS 'Student medical visit records';
CREATE INDEX idx_medical_records_student ON medical_records (student_id);
CREATE INDEX idx_medical_records_date ON medical_records (visit_date);

-- 10.2 immunization_records
CREATE TABLE immunization_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    vaccine_name VARCHAR(200) NOT NULL,
    vaccination_date DATE NOT NULL,
    next_dose_date DATE,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE immunization_records IS 'Student immunization/vaccination records';
CREATE INDEX idx_immun_records_student ON immunization_records (student_id);

-- 10.3 student_health_info
CREATE TABLE student_health_info (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL UNIQUE REFERENCES students(id) ON DELETE CASCADE,
    blood_type VARCHAR(5),
    allergies TEXT,
    medical_conditions TEXT,
    medications TEXT,
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(30),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE student_health_info IS 'Student health profile information';
CREATE INDEX idx_student_health_student ON student_health_info (student_id);

-- ============================================================================
-- 11. COUNSELING
-- ============================================================================

-- 11.1 counseling_sessions
CREATE TABLE counseling_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    counselor_id UUID NOT NULL REFERENCES users(id),
    session_date DATE NOT NULL,
    session_type VARCHAR(20) NOT NULL CHECK (session_type IN ('individual', 'group')),
    category VARCHAR(30) NOT NULL CHECK (category IN ('academic', 'behavioral', 'personal', 'career')),
    issue TEXT,
    action_plan TEXT,
    follow_up_date DATE,
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'resolved', 'referred')),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE counseling_sessions IS 'Student counseling/BK sessions';
CREATE INDEX idx_counseling_student ON counseling_sessions (student_id);
CREATE INDEX idx_counseling_counselor ON counseling_sessions (counselor_id);
CREATE INDEX idx_counseling_date ON counseling_sessions (session_date);
CREATE INDEX idx_counseling_status ON counseling_sessions (status);

-- ============================================================================
-- 12. TRANSPORTATION
-- ============================================================================

-- 12.1 transport_routes
CREATE TABLE transport_routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    driver_name VARCHAR(200),
    vehicle_number VARCHAR(30),
    capacity INTEGER CHECK (capacity > 0),
    route_description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE transport_routes IS 'School transportation routes';
CREATE INDEX idx_transport_routes_school ON transport_routes (school_id) WHERE deleted_at IS NULL;

-- 12.2 transport_students
CREATE TABLE transport_students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES transport_routes(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    pickup_point VARCHAR(200),
    pickup_time TIME,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_transport_student UNIQUE (route_id, student_id)
);
COMMENT ON TABLE transport_students IS 'Students assigned to transport routes';
CREATE INDEX idx_transport_students_route ON transport_students (route_id);
CREATE INDEX idx_transport_students_student ON transport_students (student_id);

-- 12.3 transport_attendance
CREATE TABLE transport_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES transport_routes(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    pickup_status VARCHAR(20) CHECK (pickup_status IN ('ontime', 'late', 'absent', 'no_show')),
    dropoff_status VARCHAR(20) CHECK (dropoff_status IN ('completed', 'pending', 'no_show')),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE transport_attendance IS 'Daily transport attendance';
CREATE INDEX idx_transport_att_route ON transport_attendance (route_id);
CREATE INDEX idx_transport_att_student ON transport_attendance (student_id);
CREATE INDEX idx_transport_att_date ON transport_attendance (date);

-- ============================================================================
-- 13. DORMITORY
-- ============================================================================

-- 13.1 dormitories
CREATE TABLE dormitories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    building VARCHAR(100),
    floor INTEGER,
    capacity INTEGER CHECK (capacity > 0),
    supervisor_id UUID REFERENCES users(id),
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE dormitories IS 'Dormitory/asrama buildings';
CREATE INDEX idx_dormitories_school ON dormitories (school_id) WHERE deleted_at IS NULL;

-- 13.2 dormitory_rooms
CREATE TABLE dormitory_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dormitory_id UUID NOT NULL REFERENCES dormitories(id) ON DELETE CASCADE,
    room_number VARCHAR(30) NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    occupied INTEGER DEFAULT 0 CHECK (occupied >= 0),
    type VARCHAR(30) CHECK (type IN ('regular', 'vip', 'special')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_dorm_room UNIQUE (dormitory_id, room_number),
    CONSTRAINT chk_dorm_occupied CHECK (occupied <= capacity)
);
COMMENT ON TABLE dormitory_rooms IS 'Rooms within a dormitory';
CREATE INDEX idx_dorm_rooms_dormitory ON dormitory_rooms (dormitory_id);

-- 13.3 dormitory_residents
CREATE TABLE dormitory_residents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES dormitory_rooms(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    check_in_date DATE NOT NULL,
    check_out_date DATE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'checked_out', 'transferred')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_dorm_resident UNIQUE (room_id, student_id),
    CONSTRAINT chk_resident_dates CHECK (check_out_date IS NULL OR check_out_date >= check_in_date)
);
COMMENT ON TABLE dormitory_residents IS 'Student residency records';
CREATE INDEX idx_dorm_residents_room ON dormitory_residents (room_id);
CREATE INDEX idx_dorm_residents_student ON dormitory_residents (student_id);

-- 13.4 dormitory_attendance
CREATE TABLE dormitory_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resident_id UUID NOT NULL REFERENCES dormitory_residents(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    morning_checkin BOOLEAN DEFAULT false,
    evening_checkin BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'present' CHECK (status IN ('present', 'absent', 'sick', 'leave', 'late')),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_dorm_att UNIQUE (resident_id, date)
);
COMMENT ON TABLE dormitory_attendance IS 'Daily dormitory check-in records';
CREATE INDEX idx_dorm_att_resident ON dormitory_attendance (resident_id);
CREATE INDEX idx_dorm_att_date ON dormitory_attendance (date);

-- ============================================================================
-- 14. CANTEEN
-- ============================================================================

-- 14.1 canteen_products
CREATE TABLE canteen_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(50),
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    nutritional_info TEXT,
    is_halal BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE canteen_products IS 'Canteen product catalog';
CREATE INDEX idx_canteen_products_school ON canteen_products (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_canteen_products_category ON canteen_products (school_id, category);
CREATE INDEX idx_canteen_products_fts ON canteen_products USING GIN (to_tsvector('simple', coalesce(name, '')));

-- 14.2 canteen_orders
CREATE TABLE canteen_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    order_date DATE NOT NULL DEFAULT CURRENT_DATE,
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    payment_method VARCHAR(30) CHECK (payment_method IN ('cash', 'wallet', 'subscription')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'preparing', 'served', 'cancelled')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE canteen_orders IS 'Canteen order transactions';
CREATE INDEX idx_canteen_orders_student ON canteen_orders (student_id);
CREATE INDEX idx_canteen_orders_date ON canteen_orders (order_date);

-- 14.3 canteen_order_items
CREATE TABLE canteen_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES canteen_orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES canteen_products(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    subtotal DECIMAL(10,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE canteen_order_items IS 'Line items in a canteen order';
CREATE INDEX idx_canteen_oi_order ON canteen_order_items (order_id);

-- ============================================================================
-- 15. NOTIFICATIONS
-- ============================================================================

-- 15.1 notifications
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'info' CHECK (type IN ('info', 'success', 'warning', 'error')),
    category VARCHAR(50),
    reference_type VARCHAR(50),
    reference_id UUID,
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE notifications IS 'User notifications';
CREATE INDEX idx_notifications_user ON notifications (user_id, is_read, created_at DESC);
CREATE INDEX idx_notifications_school ON notifications (school_id);
CREATE INDEX idx_notifications_category ON notifications (user_id, category);
CREATE INDEX idx_notifications_type ON notifications (user_id, type);
CREATE INDEX idx_notifications_ref ON notifications (reference_type, reference_id);

-- 15.2 notification_templates
CREATE TABLE notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(50) NOT NULL,
    subject_template TEXT,
    body_template TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    channel VARCHAR(20) NOT NULL CHECK (channel IN ('email', 'sms', 'push', 'in_app')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_notif_template_code UNIQUE (school_id, code)
);
COMMENT ON TABLE notification_templates IS 'Notification message templates';
COMMENT ON COLUMN notification_templates.channel IS 'Delivery channel: email, sms, push, in_app';
CREATE INDEX idx_notif_templates_school ON notification_templates (school_id) WHERE deleted_at IS NULL;

-- 15.3 notification_preferences
CREATE TABLE notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category VARCHAR(50) NOT NULL,
    email_enabled BOOLEAN DEFAULT true,
    sms_enabled BOOLEAN DEFAULT false,
    push_enabled BOOLEAN DEFAULT true,
    in_app_enabled BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_notif_pref UNIQUE (user_id, category)
);
COMMENT ON TABLE notification_preferences IS 'User notification channel preferences';
CREATE INDEX idx_notif_pref_user ON notification_preferences (user_id);

-- ============================================================================
-- 16. DOCUMENTS / LETTERS
-- ============================================================================

-- 16.1 documents
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    description TEXT,
    file_url TEXT NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    version INTEGER DEFAULT 1,
    tags TEXT[],
    is_public BOOLEAN DEFAULT false,
    uploaded_by UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE documents IS 'Document management with versioning';
CREATE INDEX idx_documents_school ON documents (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_documents_type ON documents (school_id, document_type);
CREATE INDEX idx_documents_tags ON documents USING GIN (tags);
CREATE INDEX idx_documents_fts ON documents USING GIN (to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(description, '')));

-- 16.2 document_versions
CREATE TABLE document_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    file_url TEXT NOT NULL,
    changes_description TEXT,
    uploaded_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_doc_version UNIQUE (document_id, version)
);
COMMENT ON TABLE document_versions IS 'Historical versions of documents';
CREATE INDEX idx_doc_versions_doc ON document_versions (document_id);

-- 16.3 letters
CREATE TABLE letters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    letter_number VARCHAR(100),
    letter_type VARCHAR(50) NOT NULL,
    subject VARCHAR(500) NOT NULL,
    recipient VARCHAR(300),
    content TEXT,
    sender_id UUID NOT NULL REFERENCES users(id),
    send_date DATE,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'received', 'archived')),
    file_url TEXT,
    template_id UUID,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE letters IS 'Official school letters/correspondence';
CREATE INDEX idx_letters_school ON letters (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_letters_type ON letters (school_id, letter_type);
CREATE INDEX idx_letters_status ON letters (status);
CREATE INDEX idx_letters_fts ON letters USING GIN (to_tsvector('simple', coalesce(subject, '') || ' ' || coalesce(letter_number, '')));

-- 16.4 letter_templates
CREATE TABLE letter_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    letter_type VARCHAR(50) NOT NULL,
    content_template TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE letter_templates IS 'Reusable letter templates';
CREATE INDEX idx_letter_templates_school ON letter_templates (school_id) WHERE deleted_at IS NULL;

-- 16.5 approval_workflows
CREATE TABLE approval_workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    module VARCHAR(100) NOT NULL,
    steps JSONB NOT NULL DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE approval_workflows IS 'Approval workflow definitions';
COMMENT ON COLUMN approval_workflows.steps IS 'JSON array of step definitions with role/user assignees';
CREATE INDEX idx_approval_workflows_school ON approval_workflows (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_approval_workflows_module ON approval_workflows (school_id, module);

-- 16.6 approval_requests
CREATE TABLE approval_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES approval_workflows(id) ON DELETE RESTRICT,
    reference_type VARCHAR(50) NOT NULL,
    reference_id UUID NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id),
    current_step INTEGER DEFAULT 1 CHECK (current_step > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE approval_requests IS 'Active approval request instances';
CREATE INDEX idx_approval_requests_workflow ON approval_requests (workflow_id);
CREATE INDEX idx_approval_requests_ref ON approval_requests (reference_type, reference_id);
CREATE INDEX idx_approval_requests_status ON approval_requests (status);
CREATE INDEX idx_approval_requests_requester ON approval_requests (requested_by);

-- 16.7 approval_actions
CREATE TABLE approval_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    approval_request_id UUID NOT NULL REFERENCES approval_requests(id) ON DELETE CASCADE,
    step INTEGER NOT NULL CHECK (step > 0),
    approved_by UUID NOT NULL REFERENCES users(id),
    action VARCHAR(20) NOT NULL CHECK (action IN ('approve', 'reject', 'revise')),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE approval_actions IS 'Individual approval step actions';
CREATE INDEX idx_approval_actions_request ON approval_actions (approval_request_id);
CREATE INDEX idx_approval_actions_user ON approval_actions (approved_by);

-- ============================================================================
-- 17. MEETINGS / CALENDAR
-- ============================================================================

-- 17.1 meetings
CREATE TABLE meetings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    meeting_type VARCHAR(50) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    location VARCHAR(200),
    organizer_id UUID NOT NULL REFERENCES users(id),
    minutes_file_url TEXT,
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'ongoing', 'completed', 'cancelled')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_meeting_time CHECK (end_time > start_time)
);
COMMENT ON TABLE meetings IS 'Scheduled meetings';
CREATE INDEX idx_meetings_school ON meetings (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_meetings_organizer ON meetings (organizer_id);
CREATE INDEX idx_meetings_time ON meetings (start_time, end_time);

-- 17.2 meeting_participants
CREATE TABLE meeting_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    attendance_status VARCHAR(20) DEFAULT 'pending' CHECK (attendance_status IN ('present', 'absent', 'excused', 'pending')),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_meeting_participant UNIQUE (meeting_id, user_id)
);
COMMENT ON TABLE meeting_participants IS 'Meeting participant list';
CREATE INDEX idx_meeting_parts_meeting ON meeting_participants (meeting_id);
CREATE INDEX idx_meeting_parts_user ON meeting_participants (user_id);

-- 17.3 events
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    event_type VARCHAR(50) NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    location VARCHAR(200),
    organizer_id UUID REFERENCES users(id),
    is_public BOOLEAN DEFAULT false,
    color VARCHAR(20) DEFAULT '#3788d8',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_event_dates CHECK (end_date >= start_date)
);
COMMENT ON TABLE events IS 'School events calendar';
CREATE INDEX idx_events_school ON events (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_dates ON events (school_id, start_date, end_date);
CREATE INDEX idx_events_type ON events (school_id, event_type);

-- 17.4 task_boards
CREATE TABLE task_boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    team_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE task_boards IS 'Kanban-style task boards';
CREATE INDEX idx_task_boards_school ON task_boards (school_id) WHERE deleted_at IS NULL;

-- 17.5 tasks
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID NOT NULL REFERENCES task_boards(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    assignee_id UUID REFERENCES users(id),
    priority VARCHAR(10) NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    due_date TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'in_progress', 'done', 'cancelled')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tasks IS 'Task cards on boards';
CREATE INDEX idx_tasks_board ON tasks (board_id);
CREATE INDEX idx_tasks_assignee ON tasks (assignee_id);
CREATE INDEX idx_tasks_status ON tasks (status);
CREATE INDEX idx_tasks_due_date ON tasks (due_date) WHERE status NOT IN ('done', 'cancelled');

-- 17.6 task_comments
CREATE TABLE task_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE task_comments IS 'Comments on task cards';
CREATE INDEX idx_task_comments_task ON task_comments (task_id);

-- ============================================================================
-- 18. ANNOUNCEMENTS
-- ============================================================================

-- 18.1 announcements
CREATE TABLE announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(50),
    priority VARCHAR(10) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    target_roles JSONB,
    target_grades JSONB,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    is_pinned BOOLEAN DEFAULT false,
    published_by UUID REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_announcement_dates CHECK (end_date IS NULL OR end_date >= start_date)
);
COMMENT ON TABLE announcements IS 'School announcements and bulletins';
CREATE INDEX idx_announcements_school ON announcements (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_announcements_status ON announcements (status);
CREATE INDEX idx_announcements_dates ON announcements (school_id, start_date, end_date) WHERE status = 'published';
CREATE INDEX idx_announcements_fts ON announcements USING GIN (to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(content, '')));

-- ============================================================================
-- 19. SETTINGS
-- ============================================================================

-- 19.1 settings
CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    key VARCHAR(200) NOT NULL,
    value JSONB NOT NULL DEFAULT '{}',
    description TEXT,
    is_public BOOLEAN DEFAULT false,
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_settings_key UNIQUE (school_id, key)
);
COMMENT ON TABLE settings IS 'School-level configuration settings';
CREATE INDEX idx_settings_school ON settings (school_id);

-- 19.2 system_settings
CREATE TABLE system_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(200) UNIQUE NOT NULL,
    value JSONB NOT NULL DEFAULT '{}',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE system_settings IS 'System-wide configuration settings (not school-specific)';
CREATE INDEX idx_system_settings_key ON system_settings (key);

-- ============================================================================
-- 20. AUDIT / LOGS
-- ============================================================================

-- 20.1 audit_logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE audit_logs IS 'Data change audit trail';
CREATE INDEX idx_audit_logs_school ON audit_logs (school_id);
CREATE INDEX idx_audit_logs_user ON audit_logs (user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_type, entity_id);
CREATE INDEX idx_audit_logs_action ON audit_logs (action);
CREATE INDEX idx_audit_logs_created ON audit_logs (created_at DESC);

-- 20.2 activity_logs
CREATE TABLE activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    module VARCHAR(100),
    description TEXT,
    metadata JSONB,
    ip_address INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE activity_logs IS 'User activity tracking logs';
CREATE INDEX idx_activity_logs_user ON activity_logs (user_id);
CREATE INDEX idx_activity_logs_school ON activity_logs (school_id);
CREATE INDEX idx_activity_logs_module ON activity_logs (module);
CREATE INDEX idx_activity_logs_created ON activity_logs (created_at DESC);

-- 20.3 login_logs
CREATE TABLE login_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    login_time TIMESTAMPTZ NOT NULL,
    logout_time TIMESTAMPTZ,
    ip_address INET,
    user_agent TEXT,
    status VARCHAR(10) NOT NULL CHECK (status IN ('success', 'failed')),
    failure_reason VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE login_logs IS 'User login/logout audit records';
CREATE INDEX idx_login_logs_user ON login_logs (user_id);
CREATE INDEX idx_login_logs_time ON login_logs (login_time DESC);
CREATE INDEX idx_login_logs_status ON login_logs (status);

-- ============================================================================
-- 21. MOSQUE / ISLAMIC EVENTS
-- ============================================================================

-- 21.1 mosque_activities
CREATE TABLE mosque_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(300) NOT NULL,
    activity_type VARCHAR(50) NOT NULL,
    schedule TEXT,
    location VARCHAR(200),
    imam_id UUID REFERENCES users(id),
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE mosque_activities IS 'School mosque activity schedules';
CREATE INDEX idx_mosque_activities_school ON mosque_activities (school_id) WHERE deleted_at IS NULL;

-- 21.2 islamic_events
CREATE TABLE islamic_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(300) NOT NULL,
    event_type VARCHAR(30) NOT NULL CHECK (event_type IN ('isra_miraj', 'maulid', 'ramadhan', 'eid', 'hijri_new_year', 'other')),
    islamic_date VARCHAR(50),
    gregorian_date DATE,
    description TEXT,
    status VARCHAR(20) DEFAULT 'planned' CHECK (status IN ('planned', 'ongoing', 'completed', 'cancelled')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE islamic_events IS 'Islamic holiday and event records';
CREATE INDEX idx_islamic_events_school ON islamic_events (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_islamic_events_type ON islamic_events (event_type);
CREATE INDEX idx_islamic_events_date ON islamic_events (gregorian_date);

-- 21.3 ramadhan_programs
CREATE TABLE ramadhan_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT,
    name VARCHAR(300) NOT NULL,
    program_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    description TEXT,
    target_participants VARCHAR(200),
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'completed')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_ramadhan_dates CHECK (end_date > start_date)
);
COMMENT ON TABLE ramadhan_programs IS 'Special Ramadhan program management';
CREATE INDEX idx_ramadhan_programs_school ON ramadhan_programs (school_id);
CREATE INDEX idx_ramadhan_programs_ay ON ramadhan_programs (academic_year_id);

-- 21.4 zakat_reports
CREATE TABLE zakat_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    report_date DATE NOT NULL,
    zakat_type VARCHAR(20) NOT NULL CHECK (zakat_type IN ('fitrah', 'maal')),
    muzakki_name VARCHAR(200),
    amount DECIMAL(14,2) NOT NULL CHECK (amount >= 0),
    recipient_count INTEGER CHECK (recipient_count >= 0),
    distribution_date DATE,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE zakat_reports IS 'Zakat collection and distribution records';
CREATE INDEX idx_zakat_reports_school ON zakat_reports (school_id);
CREATE INDEX idx_zakat_reports_date ON zakat_reports (report_date);
CREATE INDEX idx_zakat_reports_type ON zakat_reports (zakat_type);

-- ============================================================================
-- 22. ADMISSIONS (PPDB)
-- ============================================================================

-- 22.1 admission_batches
CREATE TABLE admission_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    quota INTEGER CHECK (quota > 0),
    registration_fee DECIMAL(12,2),
    requirements JSONB DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'open', 'closed', 'completed')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_batch_dates CHECK (end_date > start_date)
);
COMMENT ON TABLE admission_batches IS 'PPDB admission batch definitions';
CREATE INDEX idx_admission_batches_school ON admission_batches (school_id);
CREATE INDEX idx_admission_batches_ay ON admission_batches (academic_year_id);
CREATE INDEX idx_admission_batches_status ON admission_batches (status);

-- 22.2 admission_applicants
CREATE TABLE admission_applicants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES admission_batches(id) ON DELETE RESTRICT,
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
    previous_school VARCHAR(200),
    parent_name VARCHAR(255),
    parent_phone VARCHAR(30),
    parent_email VARCHAR(255),
    address TEXT,
    registration_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'registered' CHECK (status IN ('registered', 'tested', 'interviewed', 'accepted', 'rejected', 'enrolled')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE admission_applicants IS 'PPDB applicant data';
CREATE INDEX idx_applicants_batch ON admission_applicants (batch_id);
CREATE INDEX idx_applicants_status ON admission_applicants (status);
CREATE INDEX idx_applicants_reg_number ON admission_applicants (registration_number);
CREATE INDEX idx_applicants_fts ON admission_applicants USING GIN (to_tsvector('simple', coalesce(full_name, '') || ' ' || coalesce(registration_number, '')));

-- 22.3 admission_exams
CREATE TABLE admission_exams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id UUID NOT NULL REFERENCES admission_applicants(id) ON DELETE CASCADE,
    exam_type VARCHAR(30) NOT NULL CHECK (exam_type IN ('written', 'interview', 'quran_test')),
    exam_date DATE,
    score DECIMAL(6,2) CHECK (score >= 0),
    max_score DECIMAL(6,2) DEFAULT 100,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE admission_exams IS 'Admission test/exam results';
CREATE INDEX idx_admission_exams_applicant ON admission_exams (applicant_id);
CREATE INDEX idx_admission_exams_type ON admission_exams (exam_type);

-- 22.4 admission_documents
CREATE TABLE admission_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id UUID NOT NULL REFERENCES admission_applicants(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    file_url TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE admission_documents IS 'Admission applicant documents';
CREATE INDEX idx_admission_docs_applicant ON admission_documents (applicant_id);

-- ============================================================================
-- 23. GRADUATION
-- ============================================================================

-- 23.1 graduation_batches
CREATE TABLE graduation_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    ceremony_date DATE,
    location VARCHAR(200),
    status VARCHAR(20) DEFAULT 'planned' CHECK (status IN ('planned', 'preparing', 'held', 'completed')),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE graduation_batches IS 'Graduation batch/ceremony records';
CREATE INDEX idx_grad_batches_school ON graduation_batches (school_id);
CREATE INDEX idx_grad_batches_ay ON graduation_batches (academic_year_id);

-- 23.2 graduation_candidates
CREATE TABLE graduation_candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES graduation_batches(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'graduated', 'certificate_issued')),
    final_gpa DECIMAL(4,2) CHECK (final_gpa >= 0 AND final_gpa <= 4.00),
    ranking INTEGER,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_grad_candidate UNIQUE (batch_id, student_id)
);
COMMENT ON TABLE graduation_candidates IS 'Students nominated for graduation';
CREATE INDEX idx_grad_candidates_batch ON graduation_candidates (batch_id);
CREATE INDEX idx_grad_candidates_student ON graduation_candidates (student_id);

-- ============================================================================
-- 24. AI / KNOWLEDGE BASE
-- ============================================================================

-- 24.1 knowledge_documents
CREATE TABLE knowledge_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    content TEXT,
    document_type VARCHAR(50),
    file_url TEXT,
    processed_text TEXT,
    embedding_status VARCHAR(20) DEFAULT 'pending' CHECK (embedding_status IN ('pending', 'processing', 'completed', 'failed')),
    qdrant_id VARCHAR(100),
    chunk_count INTEGER DEFAULT 0,
    uploaded_by UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE knowledge_documents IS 'AI knowledge base documents with embedding support';
CREATE INDEX idx_knowledge_docs_school ON knowledge_documents (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_knowledge_docs_status ON knowledge_documents (embedding_status);
CREATE INDEX idx_knowledge_docs_qdrant ON knowledge_documents (qdrant_id);

-- 24.2 knowledge_chunks
CREATE TABLE knowledge_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES knowledge_documents(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL CHECK (chunk_index >= 0),
    content TEXT NOT NULL,
    embedding_qdrant_id VARCHAR(100),
    token_count INTEGER,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_knowledge_chunk UNIQUE (document_id, chunk_index)
);
COMMENT ON TABLE knowledge_chunks IS 'Embedding chunks from knowledge documents';
CREATE INDEX idx_knowledge_chunks_doc ON knowledge_chunks (document_id);

-- 24.3 ai_conversations
CREATE TABLE ai_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(300),
    context_type VARCHAR(50),
    context_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE ai_conversations IS 'AI chat conversation threads';
CREATE INDEX idx_ai_conversations_user ON ai_conversations (user_id);
CREATE INDEX idx_ai_conversations_context ON ai_conversations (context_type, context_id);

-- 24.4 ai_messages
CREATE TABLE ai_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    tokens_used INTEGER,
    model VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE ai_messages IS 'Individual messages within AI conversations';
CREATE INDEX idx_ai_messages_conversation ON ai_messages (conversation_id, created_at);

-- 24.5 ai_prompts
CREATE TABLE ai_prompts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    prompt_type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE ai_prompts IS 'AI prompt templates';
CREATE INDEX idx_ai_prompts_school ON ai_prompts (school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ai_prompts_type ON ai_prompts (school_id, prompt_type);

-- 24.6 ai_generation_history
CREATE TABLE ai_generation_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    prompt_id UUID REFERENCES ai_prompts(id) ON DELETE SET NULL,
    input_data JSONB,
    output_data JSONB,
    model_used VARCHAR(100),
    tokens_used INTEGER,
    rating SMALLINT CHECK (rating >= 1 AND rating <= 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE ai_generation_history IS 'AI generation output history and feedback';
CREATE INDEX idx_ai_gen_history_user ON ai_generation_history (user_id);
CREATE INDEX idx_ai_gen_history_prompt ON ai_generation_history (prompt_id);

-- ============================================================================
-- TRIGGER: Auto-update updated_at
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at trigger to all tables with updated_at column
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN
        SELECT table_name
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND column_name = 'updated_at'
          AND table_name NOT IN ('audit_logs', 'activity_logs', 'login_logs', 'notification_preferences')
    LOOP
        EXECUTE format(
            'CREATE TRIGGER trg_%I_updated_at BEFORE UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()',
            r.table_name,
            r.table_name
        );
    END LOOP;
END
$$;

COMMIT;
