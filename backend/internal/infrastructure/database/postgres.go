package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/config"
	"go.uber.org/zap"
)

func NewPostgresDB(cfg config.DatabaseConfig, logger *zap.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	logger.Info("connected to postgres",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("db", cfg.Name),
	)

	return db, nil
}

func RunMigrations(db *sqlx.DB, logger *zap.Logger) error {
	schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";

	CREATE TABLE IF NOT EXISTS schools (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		npsn VARCHAR(20),
		address TEXT,
		city VARCHAR(100),
		province VARCHAR(100),
		postal_code VARCHAR(10),
		phone VARCHAR(20),
		email VARCHAR(255),
		website VARCHAR(255),
		logo_url TEXT,
		type VARCHAR(50) DEFAULT 'formal',
		accreditation VARCHAR(20),
		established_date VARCHAR(20),
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS academic_years (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		name VARCHAR(100) NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE NOT NULL,
		is_active BOOLEAN DEFAULT false,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS semesters (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		academic_year_id UUID NOT NULL REFERENCES academic_years(id),
		name VARCHAR(50) NOT NULL,
		semester_number INT NOT NULL CHECK (semester_number IN (1,2)),
		start_date DATE NOT NULL,
		end_date DATE NOT NULL,
		is_active BOOLEAN DEFAULT false,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS grades (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		name VARCHAR(100) NOT NULL,
		level INT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(school_id, level)
	);

	CREATE TABLE IF NOT EXISTS classes (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		grade_id UUID NOT NULL REFERENCES grades(id),
		name VARCHAR(100) NOT NULL,
		capacity INT DEFAULT 40,
		homeroom_teacher_id UUID,
		academic_year_id UUID NOT NULL REFERENCES academic_years(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS subjects (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		code VARCHAR(20) NOT NULL,
		name VARCHAR(255) NOT NULL,
		category VARCHAR(50) NOT NULL DEFAULT 'general',
		description TEXT,
		kkm DECIMAL(5,2) DEFAULT 75.00,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(school_id, code)
	);

	CREATE TABLE IF NOT EXISTS curriculums (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		grade_id UUID NOT NULL REFERENCES grades(id),
		subject_id UUID NOT NULL REFERENCES subjects(id),
		semester_id UUID NOT NULL REFERENCES semesters(id),
		content TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS permissions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(100) NOT NULL,
		slug VARCHAR(100) NOT NULL UNIQUE,
		module VARCHAR(50) NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS roles (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		name VARCHAR(100) NOT NULL,
		slug VARCHAR(100) NOT NULL,
		description TEXT,
		is_system BOOLEAN DEFAULT false,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(school_id, slug)
	);

	CREATE TABLE IF NOT EXISTS role_permissions (
		role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
		permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
		PRIMARY KEY (role_id, permission_id)
	);

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		email VARCHAR(255) NOT NULL,
		username VARCHAR(50) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		full_name VARCHAR(255) NOT NULL,
		avatar_url TEXT,
		phone VARCHAR(20),
		is_active BOOLEAN DEFAULT true,
		email_verified_at TIMESTAMPTZ,
		last_login_at TIMESTAMPTZ,
		password_changed_at TIMESTAMPTZ DEFAULT NOW(),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		deleted_at TIMESTAMPTZ,
		UNIQUE(school_id, email),
		UNIQUE(school_id, username)
	);

	CREATE TABLE IF NOT EXISTS user_roles (
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
		PRIMARY KEY (user_id, role_id)
	);

	CREATE TABLE IF NOT EXISTS user_profiles (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		place_of_birth VARCHAR(100),
		date_of_birth DATE,
		gender VARCHAR(10) CHECK (gender IN ('male','female')),
		religion VARCHAR(50) DEFAULT 'islam',
		address TEXT,
		bio TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(user_id)
	);

	CREATE TABLE IF NOT EXISTS user_sessions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		refresh_token TEXT NOT NULL,
		user_agent TEXT,
		ip_address VARCHAR(50),
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS students (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id),
		school_id UUID NOT NULL REFERENCES schools(id),
		nis VARCHAR(20),
		nisn VARCHAR(20),
		nik VARCHAR(20),
		class_id UUID NOT NULL REFERENCES classes(id),
		academic_year_id UUID NOT NULL REFERENCES academic_years(id),
		enrollment_date DATE NOT NULL,
		status VARCHAR(20) DEFAULT 'active',
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(school_id, nis)
	);

	CREATE TABLE IF NOT EXISTS student_parents (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id),
		relation VARCHAR(20) NOT NULL CHECK (relation IN ('father','mother','guardian')),
		is_primary BOOLEAN DEFAULT false,
		occupation VARCHAR(100),
		institution VARCHAR(255),
		income DECIMAL(15,2) DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(student_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS student_documents (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		doc_type VARCHAR(50) NOT NULL,
		file_url TEXT NOT NULL,
		status VARCHAR(20) DEFAULT 'pending',
		notes TEXT,
		verified_at TIMESTAMPTZ,
		verified_by UUID REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS teachers (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id),
		school_id UUID NOT NULL REFERENCES schools(id),
		nip VARCHAR(30),
		nik VARCHAR(20),
		nupk VARCHAR(20),
		status VARCHAR(20) DEFAULT 'permanent',
		join_date DATE NOT NULL,
		education_level VARCHAR(10),
		major VARCHAR(100),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(school_id, nip)
	);

	CREATE TABLE IF NOT EXISTS teacher_subjects (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
		subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
		class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(teacher_id, subject_id, class_id)
	);

	CREATE TABLE IF NOT EXISTS employees (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id),
		school_id UUID NOT NULL REFERENCES schools(id),
		nip VARCHAR(30),
		nik VARCHAR(20),
		position VARCHAR(100) NOT NULL,
		department VARCHAR(100) NOT NULL,
		join_date DATE NOT NULL,
		status VARCHAR(20) DEFAULT 'active',
		base_salary DECIMAL(15,2) DEFAULT 0,
		bank_account VARCHAR(50),
		bank_name VARCHAR(100),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS employee_attendances (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		employee_id UUID NOT NULL REFERENCES employees(id),
		date DATE NOT NULL,
		check_in TIMESTAMPTZ NOT NULL,
		check_out TIMESTAMPTZ,
		status VARCHAR(20) DEFAULT 'present',
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS leave_requests (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		employee_id UUID NOT NULL REFERENCES employees(id),
		leave_type VARCHAR(20) NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE NOT NULL,
		reason TEXT NOT NULL,
		status VARCHAR(20) DEFAULT 'pending',
		approved_by UUID REFERENCES users(id),
		approved_at TIMESTAMPTZ,
		attachment_url TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS schedules (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		class_id UUID NOT NULL REFERENCES classes(id),
		subject_id UUID NOT NULL REFERENCES subjects(id),
		teacher_id UUID NOT NULL REFERENCES teachers(id),
		day VARCHAR(10) NOT NULL,
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		room VARCHAR(50),
		semester_id UUID NOT NULL REFERENCES semesters(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS attendances (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id),
		schedule_id UUID NOT NULL REFERENCES schedules(id),
		date DATE NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'present',
		notes TEXT,
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS exams (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		subject_id UUID NOT NULL REFERENCES subjects(id),
		class_id UUID NOT NULL REFERENCES classes(id),
		semester_id UUID NOT NULL REFERENCES semesters(id),
		name VARCHAR(255) NOT NULL,
		exam_type VARCHAR(20) NOT NULL,
		date DATE NOT NULL,
		duration INT DEFAULT 60,
		total_score DECIMAL(5,2) DEFAULT 100.00,
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS exam_results (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		exam_id UUID NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
		student_id UUID NOT NULL REFERENCES students(id),
		score DECIMAL(5,2) NOT NULL,
		grade VARCHAR(5),
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(exam_id, student_id)
	);

	CREATE TABLE IF NOT EXISTS gradebooks (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		class_id UUID NOT NULL REFERENCES classes(id),
		subject_id UUID NOT NULL REFERENCES subjects(id),
		teacher_id UUID NOT NULL REFERENCES teachers(id),
		semester_id UUID NOT NULL REFERENCES semesters(id),
		student_id UUID NOT NULL REFERENCES students(id),
		daily_score DECIMAL(5,2) DEFAULT 0,
		mid_score DECIMAL(5,2) DEFAULT 0,
		final_score DECIMAL(5,2) DEFAULT 0,
		practice_score DECIMAL(5,2) DEFAULT 0,
		attitude VARCHAR(20),
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(class_id, subject_id, student_id, semester_id)
	);

	CREATE TABLE IF NOT EXISTS report_cards (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id),
		semester_id UUID NOT NULL REFERENCES semesters(id),
		class_id UUID NOT NULL REFERENCES classes(id),
		average_score DECIMAL(5,2) DEFAULT 0,
		rank INT,
		absent_count INT DEFAULT 0,
		sick_count INT DEFAULT 0,
		permit_count INT DEFAULT 0,
		homeroom_comment TEXT,
		parent_signature BOOLEAN DEFAULT false,
		approved_by UUID NOT NULL REFERENCES users(id),
		approved_at TIMESTAMPTZ DEFAULT NOW(),
		published_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(student_id, semester_id)
	);

	CREATE TABLE IF NOT EXISTS assignments (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		subject_id UUID NOT NULL REFERENCES subjects(id),
		class_id UUID NOT NULL REFERENCES classes(id),
		teacher_id UUID NOT NULL REFERENCES teachers(id),
		title VARCHAR(255) NOT NULL,
		description TEXT,
		due_date DATE NOT NULL,
		max_score DECIMAL(5,2) DEFAULT 100.00,
		attachments TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS assignment_submissions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
		student_id UUID NOT NULL REFERENCES students(id),
		content TEXT,
		file_url TEXT,
		score DECIMAL(5,2),
		feedback TEXT,
		submitted_at TIMESTAMPTZ DEFAULT NOW(),
		graded_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(assignment_id, student_id)
	);

	CREATE TABLE IF NOT EXISTS lesson_plans (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		subject_id UUID NOT NULL REFERENCES subjects(id),
		class_id UUID NOT NULL REFERENCES classes(id),
		teacher_id UUID NOT NULL REFERENCES teachers(id),
		title VARCHAR(255) NOT NULL,
		date DATE NOT NULL,
		objectives TEXT NOT NULL,
		materials TEXT NOT NULL,
		activities TEXT NOT NULL,
		assessment TEXT,
		reflection TEXT,
		status VARCHAR(20) DEFAULT 'draft',
		approved_by UUID REFERENCES users(id),
		approved_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS teaching_journals (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		teacher_id UUID NOT NULL REFERENCES teachers(id),
		schedule_id UUID NOT NULL REFERENCES schedules(id),
		date DATE NOT NULL,
		material TEXT NOT NULL,
		method VARCHAR(50) NOT NULL,
		attend_count INT DEFAULT 0,
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(teacher_id, schedule_id, date)
	);

	CREATE TABLE IF NOT EXISTS tahfidz_progress (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id),
		teacher_id UUID NOT NULL REFERENCES teachers(id),
		surah VARCHAR(100) NOT NULL,
		start_ayah INT NOT NULL,
		end_ayah INT NOT NULL,
		juz INT CHECK (juz BETWEEN 1 AND 30),
		page INT,
		status VARCHAR(20) DEFAULT 'memorizing',
		quality VARCHAR(20) CHECK (quality IN ('excellent','good','fair','poor')),
		notes TEXT,
		date DATE NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS mutabaah_yaumiyah (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id),
		date DATE NOT NULL,
		fajr BOOLEAN DEFAULT false,
		dhuhr BOOLEAN DEFAULT false,
		asr BOOLEAN DEFAULT false,
		maghrib BOOLEAN DEFAULT false,
		isha BOOLEAN DEFAULT false,
		tahajjud BOOLEAN DEFAULT false,
		dhuha BOOLEAN DEFAULT false,
		sunnah BOOLEAN DEFAULT false,
		quran_tilawah INT DEFAULT 0,
		quran_hifdz INT DEFAULT 0,
		dzikr_pagi BOOLEAN DEFAULT false,
		dzikr_petang BOOLEAN DEFAULT false,
		shadaqah BOOLEAN DEFAULT false,
		puasa_sunnah BOOLEAN DEFAULT false,
		wudhu_sebelum_tidur BOOLEAN DEFAULT false,
		baca_doa_tidur BOOLEAN DEFAULT false,
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(student_id, date)
	);

	CREATE TABLE IF NOT EXISTS prayer_attendances (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id),
		date DATE NOT NULL,
		fajr BOOLEAN DEFAULT false,
		dhuhr BOOLEAN DEFAULT false,
		asr BOOLEAN DEFAULT false,
		maghrib BOOLEAN DEFAULT false,
		isha BOOLEAN DEFAULT false,
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(student_id, date)
	);

	CREATE TABLE IF NOT EXISTS halaqah_groups (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		name VARCHAR(255) NOT NULL,
		teacher_id UUID NOT NULL REFERENCES teachers(id),
		room VARCHAR(50),
		day VARCHAR(10) NOT NULL,
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		max_member INT DEFAULT 30,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS halaqah_members (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		halaqah_id UUID NOT NULL REFERENCES halaqah_groups(id) ON DELETE CASCADE,
		student_id UUID NOT NULL REFERENCES students(id),
		joined_at TIMESTAMPTZ DEFAULT NOW(),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(halaqah_id, student_id)
	);

	CREATE TABLE IF NOT EXISTS fee_types (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		name VARCHAR(255) NOT NULL,
		amount DECIMAL(15,2) NOT NULL,
		category VARCHAR(50) NOT NULL,
		frequency VARCHAR(20) NOT NULL,
		description TEXT,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS invoices (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		student_id UUID NOT NULL REFERENCES students(id),
		invoice_no VARCHAR(50) NOT NULL UNIQUE,
		total_amount DECIMAL(15,2) NOT NULL,
		paid_amount DECIMAL(15,2) DEFAULT 0,
		status VARCHAR(20) DEFAULT 'unpaid',
		due_date DATE NOT NULL,
		paid_at TIMESTAMPTZ,
		semester_id UUID NOT NULL REFERENCES semesters(id),
		notes TEXT,
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS invoice_items (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
		fee_type_id UUID NOT NULL REFERENCES fee_types(id),
		name VARCHAR(255) NOT NULL,
		amount DECIMAL(15,2) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS payments (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		invoice_id UUID NOT NULL REFERENCES invoices(id),
		payment_no VARCHAR(50) NOT NULL UNIQUE,
		amount DECIMAL(15,2) NOT NULL,
		method VARCHAR(20) NOT NULL,
		status VARCHAR(20) DEFAULT 'pending',
		proof_url TEXT,
		paid_at TIMESTAMPTZ NOT NULL,
		verified_by UUID REFERENCES users(id),
		verified_at TIMESTAMPTZ,
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS journals (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		journal_no VARCHAR(50) NOT NULL UNIQUE,
		date DATE NOT NULL,
		description TEXT NOT NULL,
		status VARCHAR(20) DEFAULT 'draft',
		created_by UUID NOT NULL REFERENCES users(id),
		approved_by UUID REFERENCES users(id),
		approved_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS journal_entries (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		journal_id UUID NOT NULL REFERENCES journals(id) ON DELETE CASCADE,
		account_code VARCHAR(20) NOT NULL,
		description TEXT,
		debit DECIMAL(15,2) DEFAULT 0,
		credit DECIMAL(15,2) DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS general_ledgers (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		account_code VARCHAR(20) NOT NULL,
		account_name VARCHAR(255) NOT NULL,
		date DATE NOT NULL,
		description TEXT,
		debit DECIMAL(15,2) DEFAULT 0,
		credit DECIMAL(15,2) DEFAULT 0,
		balance DECIMAL(15,2) DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS payroll_details (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		employee_id UUID NOT NULL REFERENCES employees(id),
		school_id UUID NOT NULL REFERENCES schools(id),
		period VARCHAR(10) NOT NULL,
		base_salary DECIMAL(15,2) NOT NULL,
		allowance DECIMAL(15,2) DEFAULT 0,
		deduction DECIMAL(15,2) DEFAULT 0,
		net_salary DECIMAL(15,2) NOT NULL,
		status VARCHAR(20) DEFAULT 'pending',
		paid_at TIMESTAMPTZ,
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS assets (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		code VARCHAR(50) NOT NULL,
		name VARCHAR(255) NOT NULL,
		category VARCHAR(100) NOT NULL,
		purchase_date DATE NOT NULL,
		purchase_price DECIMAL(15,2) NOT NULL,
		current_value DECIMAL(15,2),
		location VARCHAR(255),
		condition VARCHAR(20) DEFAULT 'good',
		status VARCHAR(20) DEFAULT 'active',
		depreciation_rate DECIMAL(5,2) DEFAULT 0,
		responsible_id UUID REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(school_id, code)
	);

	CREATE TABLE IF NOT EXISTS inventory_items (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		code VARCHAR(50) NOT NULL,
		name VARCHAR(255) NOT NULL,
		category VARCHAR(100),
		unit VARCHAR(20) NOT NULL,
		stock_in INT DEFAULT 0,
		stock_out INT DEFAULT 0,
		stock_min INT DEFAULT 0,
		location VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(school_id, code)
	);

	CREATE TABLE IF NOT EXISTS stock_movements (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		item_id UUID NOT NULL REFERENCES inventory_items(id),
		type VARCHAR(10) NOT NULL CHECK (type IN ('in','out')),
		quantity INT NOT NULL,
		reference_type VARCHAR(50),
		reference_id VARCHAR(50),
		notes TEXT,
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS library_books (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		isbn VARCHAR(20),
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255),
		publisher VARCHAR(255),
		publish_year INT,
		category VARCHAR(100),
		language VARCHAR(50) DEFAULT 'indonesia',
		total_copies INT NOT NULL DEFAULT 1,
		available INT NOT NULL DEFAULT 1,
		location VARCHAR(100),
		cover_url TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS book_borrowings (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		book_id UUID NOT NULL REFERENCES library_books(id),
		borrower_id UUID NOT NULL,
		borrower_type VARCHAR(20) NOT NULL CHECK (borrower_type IN ('student','teacher','employee')),
		borrow_date DATE NOT NULL,
		due_date DATE NOT NULL,
		return_date DATE,
		status VARCHAR(20) DEFAULT 'borrowed',
		fine DECIMAL(10,2) DEFAULT 0,
		fine_paid BOOLEAN DEFAULT false,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS medical_records (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id),
		date DATE NOT NULL,
		diagnosis TEXT NOT NULL,
		treatment TEXT,
		medication TEXT,
		doc_type VARCHAR(50) DEFAULT 'checkup',
		notes TEXT,
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS counseling_sessions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id),
		counselor_id UUID NOT NULL REFERENCES users(id),
		date DATE NOT NULL,
		cause TEXT NOT NULL,
		action TEXT NOT NULL,
		mediation TEXT,
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS announcements (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		target_role VARCHAR(50),
		priority VARCHAR(20) DEFAULT 'medium',
		start_date DATE,
		end_date DATE,
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		deleted_at TIMESTAMPTZ
	);

	CREATE TABLE IF NOT EXISTS notifications (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id),
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		type VARCHAR(20) DEFAULT 'info',
		ref_type VARCHAR(50),
		ref_id VARCHAR(50),
		is_read BOOLEAN DEFAULT false,
		read_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS documents (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		title VARCHAR(255) NOT NULL,
		doc_type VARCHAR(50) NOT NULL,
		file_url TEXT NOT NULL,
		file_size BIGINT DEFAULT 0,
		mime_type VARCHAR(100),
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS letters (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		letter_no VARCHAR(50) NOT NULL UNIQUE,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		letter_type VARCHAR(20) NOT NULL,
		status VARCHAR(20) DEFAULT 'draft',
		"from" VARCHAR(255),
		"to" VARCHAR(255) NOT NULL,
		created_by UUID NOT NULL REFERENCES users(id),
		approved_by UUID REFERENCES users(id),
		approved_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS meetings (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		title VARCHAR(255) NOT NULL,
		agenda TEXT NOT NULL,
		date DATE NOT NULL,
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		location VARCHAR(255),
		minutes TEXT,
		status VARCHAR(20) DEFAULT 'scheduled',
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS meeting_attendees (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id),
		status VARCHAR(20) DEFAULT 'pending',
		created_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(meeting_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		title VARCHAR(255) NOT NULL,
		description TEXT,
		event_type VARCHAR(20) NOT NULL,
		date DATE NOT NULL,
		start_time TIME,
		end_time TIME,
		location VARCHAR(255),
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		title VARCHAR(255) NOT NULL,
		description TEXT,
		assigned_to UUID NOT NULL REFERENCES users(id),
		created_by UUID NOT NULL REFERENCES users(id),
		status VARCHAR(20) DEFAULT 'pending',
		priority VARCHAR(20) DEFAULT 'medium',
		due_date TIMESTAMPTZ,
		completed_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS settings (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		key VARCHAR(100) NOT NULL,
		value TEXT,
		type VARCHAR(20) DEFAULT 'string',
		module VARCHAR(50) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(school_id, key)
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		user_id UUID REFERENCES users(id),
		action VARCHAR(50) NOT NULL,
		entity VARCHAR(50) NOT NULL,
		entity_id VARCHAR(50),
		old_values JSONB,
		new_values JSONB,
		ip_address VARCHAR(50),
		user_agent TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS activity_logs (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id),
		school_id UUID NOT NULL REFERENCES schools(id),
		action VARCHAR(100) NOT NULL,
		module VARCHAR(50) NOT NULL,
		detail TEXT,
		ip_address VARCHAR(50),
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS admission_applicants (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		full_name VARCHAR(255) NOT NULL,
		gender VARCHAR(10) NOT NULL,
		place_of_birth VARCHAR(100),
		date_of_birth DATE,
		previous_school VARCHAR(255),
		grade_id UUID NOT NULL REFERENCES grades(id),
		registration_no VARCHAR(50) UNIQUE,
		parent_name VARCHAR(255),
		parent_phone VARCHAR(20),
		parent_email VARCHAR(255),
		status VARCHAR(20) DEFAULT 'pending',
		test_score DECIMAL(5,2),
		interview_score DECIMAL(5,2),
		accepted_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS graduation_candidates (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		student_id UUID NOT NULL REFERENCES students(id),
		class_id UUID NOT NULL REFERENCES classes(id),
		semester_id UUID NOT NULL REFERENCES semesters(id),
		status VARCHAR(20) DEFAULT 'candidate',
		final_grade DECIMAL(5,2),
		certificate_no VARCHAR(50) UNIQUE,
		certificate_url TEXT,
		graduated_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS knowledge_documents (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		school_id UUID NOT NULL REFERENCES schools(id),
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		doc_type VARCHAR(50) NOT NULL,
		module VARCHAR(50) NOT NULL,
		chunk_ids TEXT,
		embedding_status VARCHAR(20) DEFAULT 'pending',
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS ai_conversations (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id),
		school_id UUID NOT NULL REFERENCES schools(id),
		title VARCHAR(255),
		model VARCHAR(100) DEFAULT 'gpt-4o',
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS ai_messages (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		conversation_id UUID NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
		role VARCHAR(20) NOT NULL CHECK (role IN ('user','assistant','system')),
		content TEXT NOT NULL,
		token_count INT DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("database migrations completed")
	return nil
}

type TxFunc func(tx *sqlx.Tx) error

func WithTransaction(db *sqlx.DB, fn TxFunc) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v (original: %w)", rbErr, err)
		}
		return err
	}
	return tx.Commit()
}
