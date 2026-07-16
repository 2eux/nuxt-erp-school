package database

const (
	GetUserByEmail = `SELECT * FROM users WHERE email=$1 AND deleted_at IS NULL LIMIT 1`

	GetUserWithRoles = `
		SELECT u.*, COALESCE(json_agg(json_build_object('id', r.id, 'name', r.name, 'slug', r.slug))
			FILTER (WHERE r.id IS NOT NULL), '[]') as roles
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.id = $1 AND u.deleted_at IS NULL
		GROUP BY u.id
	`

	GetUserPermissions = `
		SELECT DISTINCT p.slug
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
	`

	InsertSession    = `INSERT INTO user_sessions (id, user_id, refresh_token, user_agent, ip_address, expires_at) VALUES ($1, $2, $3, $4, $5, $6)`
	DeleteSession    = `DELETE FROM user_sessions WHERE id=$1`
	GetSessionByToken = `SELECT * FROM user_sessions WHERE refresh_token=$1 AND expires_at > NOW() LIMIT 1`
	DeleteUserSessions = `DELETE FROM user_sessions WHERE user_id=$1`

	ListSchools = `SELECT * FROM schools WHERE deleted_at IS NULL ORDER BY created_at DESC`
	GetSchoolByID = `SELECT * FROM schools WHERE id=$1`

	ListAcademicYears = `
		SELECT * FROM academic_years WHERE school_id=$1 ORDER BY start_date DESC
		LIMIT $2 OFFSET $3
	`
	CountAcademicYears = `SELECT COUNT(*) FROM academic_years WHERE school_id=$1`

	GetActiveAcademicYear = `SELECT * FROM academic_years WHERE school_id=$1 AND is_active=true LIMIT 1`

	ListClasses = `
		SELECT c.*, g.name as grade_name,
			COALESCE((SELECT COUNT(*) FROM students s WHERE s.class_id = c.id AND s.status='active'), 0) as student_count
		FROM classes c
		JOIN grades g ON c.grade_id = g.id
		WHERE c.school_id=$1
		ORDER BY g.level, c.name
		LIMIT $2 OFFSET $3
	`
	CountClasses = `SELECT COUNT(*) FROM classes WHERE school_id=$1`
	GetClassByID  = `
		SELECT c.*, g.name as grade_name
		FROM classes c JOIN grades g ON c.grade_id = g.id
		WHERE c.id=$1
	`

	ListStudents = `
		SELECT s.*, u.full_name, u.email, up.gender, up.place_of_birth, up.date_of_birth, up.address, u.phone,
			c.name as class_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		JOIN classes c ON s.class_id = c.id
		WHERE s.school_id=$1 AND u.deleted_at IS NULL
	`

	ListSchedules = `
		SELECT sc.*, c.name as class_name, sub.name as subject_name, u.full_name as teacher_name
		FROM schedules sc
		JOIN classes c ON sc.class_id = c.id
		JOIN subjects sub ON sc.subject_id = sub.id
		JOIN teachers t ON sc.teacher_id = t.id
		JOIN users u ON t.user_id = u.id
		WHERE sc.semester_id=$1
	`

	GetStudentAttendanceSummary = `
		SELECT status, COUNT(*) as count
		FROM attendances
		WHERE student_id=$1 AND date >= $2 AND date <= $3
		GROUP BY status
	`

	ListTahfidzProgress = `
		SELECT tp.*, u.full_name as student_name, ut.full_name as teacher_name
		FROM tahfidz_progress tp
		JOIN students s ON tp.student_id = s.id
		JOIN users u ON s.user_id = u.id
		JOIN teachers t ON tp.teacher_id = t.id
		JOIN users ut ON t.user_id = ut.id
		WHERE tp.student_id=$1
		ORDER BY tp.date DESC
	`

	GetTahfidzSummary = `
		SELECT
			COUNT(DISTINCT surah) as total_surahs,
			COALESCE(SUM(CASE WHEN status='memorized' THEN end_ayah - start_ayah + 1 ELSE 0 END), 0) as total_ayahs,
			COUNT(DISTINCT juz) as total_juz,
			COUNT(DISTINCT CASE WHEN status='memorizing' THEN id END) as memorizing_count,
			COUNT(DISTINCT CASE WHEN status='memorized' THEN id END) as memorized_count
		FROM tahfidz_progress
		WHERE student_id=$1
	`

	ListInvoices = `
		SELECT inv.*, u.full_name as student_name
		FROM invoices inv
		JOIN students s ON inv.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE inv.school_id=$1
		ORDER BY inv.created_at DESC
		LIMIT $2 OFFSET $3
	`
	CountInvoices = `SELECT COUNT(*) FROM invoices WHERE school_id=$1`

	GetInvoiceWithItems = `
		SELECT inv.*, u.full_name as student_name
		FROM invoices inv
		JOIN students s ON inv.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE inv.id=$1
	`

	GetInvoiceItems = `SELECT * FROM invoice_items WHERE invoice_id=$1`

	ListPayments = `
		SELECT p.*, inv.invoice_no
		FROM payments p
		JOIN invoices inv ON p.invoice_id = inv.id
		WHERE inv.school_id=$1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`
	CountPayments = `
		SELECT COUNT(*) FROM payments p
		JOIN invoices inv ON p.invoice_id = inv.id
		WHERE inv.school_id=$1
	`

	ListGradebooks = `
		SELECT gb.*, sub.name as subject_name, u.full_name as student_name
		FROM gradebooks gb
		JOIN subjects sub ON gb.subject_id = sub.id
		JOIN students s ON gb.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE gb.class_id=$1 AND gb.subject_id=$2 AND gb.semester_id=$3
	`

	ListExams = `
		SELECT e.*, sub.name as subject_name, c.name as class_name
		FROM exams e
		JOIN subjects sub ON e.subject_id = sub.id
		JOIN classes c ON e.class_id = c.id
		WHERE e.class_id=$1 AND e.semester_id=$2
		ORDER BY e.date DESC
	`

	FinancialSummaryQuery = `
		SELECT
			COALESCE(SUM(CASE WHEN p.status='verified' THEN p.amount ELSE 0 END), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN je.credit > 0 THEN je.credit ELSE 0 END), 0) as total_expense,
			COALESCE(SUM(CASE WHEN inv.status != 'paid' THEN inv.total_amount - inv.paid_amount ELSE 0 END), 0) as outstanding
		FROM invoices inv
		LEFT JOIN payments p ON inv.id = p.invoice_id
		LEFT JOIN journal_entries je ON je.account_code LIKE '5%'
		WHERE inv.school_id=$1
	`

	DashboardStatsQuery = `
		SELECT
			(SELECT COUNT(*) FROM students WHERE school_id=$1 AND status='active') as total_students,
			(SELECT COUNT(*) FROM teachers WHERE school_id=$1) as total_teachers,
			(SELECT COUNT(*) FROM employees WHERE school_id=$1 AND status='active') as total_employees,
			(SELECT COUNT(*) FROM classes WHERE school_id=$1) as total_classes
	`

	ListNotifications = `
		SELECT * FROM notifications
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	CountNotifications      = `SELECT COUNT(*) FROM notifications WHERE user_id=$1`
	CountUnreadNotifications = `SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND is_read=false`

	ListMutabaah = `
		SELECT * FROM mutabaah_yaumiyah
		WHERE student_id=$1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	ListPrayerAttendance = `
		SELECT * FROM prayer_attendances
		WHERE student_id=$1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	ListTeachers = `
		SELECT t.*, u.full_name, u.email, u.phone
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		WHERE t.school_id=$1 AND u.deleted_at IS NULL
		ORDER BY u.full_name
		LIMIT $2 OFFSET $3
	`
	CountTeachers = `SELECT COUNT(*) FROM teachers t JOIN users u ON t.user_id = u.id WHERE t.school_id=$1 AND u.deleted_at IS NULL`

	ListEmployees = `
		SELECT e.*, u.full_name, u.email, u.phone
		FROM employees e
		JOIN users u ON e.user_id = u.id
		WHERE e.school_id=$1 AND u.deleted_at IS NULL
		ORDER BY u.full_name
		LIMIT $2 OFFSET $3
	`
	CountEmployees = `SELECT COUNT(*) FROM employees e JOIN users u ON e.user_id = u.id WHERE e.school_id=$1 AND u.deleted_at IS NULL`

	ListAnnouncements = `
		SELECT * FROM announcements
		WHERE school_id=$1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	ListAIConversations = `
		SELECT * FROM ai_conversations
		WHERE user_id=$1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	ListAIMessages = `
		SELECT * FROM ai_messages
		WHERE conversation_id=$1
		ORDER BY created_at ASC
	`

	ListHalaqahGroups = `
		SELECT hg.*, u.full_name as teacher_name,
			(SELECT COUNT(*) FROM halaqah_members hm WHERE hm.halaqah_id = hg.id) as member_count
		FROM halaqah_groups hg
		JOIN teachers t ON hg.teacher_id = t.id
		JOIN users u ON t.user_id = u.id
		WHERE hg.school_id=$1
		ORDER BY hg.created_at DESC
		LIMIT $2 OFFSET $3
	`

	ListHalaqahMembers = `
		SELECT hm.*, u.full_name
		FROM halaqah_members hm
		JOIN students s ON hm.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE hm.halaqah_id=$1
	`

	ListLeaveRequests = `
		SELECT lr.*, u.full_name
		FROM leave_requests lr
		JOIN employees e ON lr.employee_id = e.id
		JOIN users u ON e.user_id = u.id
		WHERE e.school_id=$1
		ORDER BY lr.created_at DESC
		LIMIT $2 OFFSET $3
	`

	ListStudentParents = `
		SELECT sp.*, u.full_name, u.email, u.phone
		FROM student_parents sp
		JOIN users u ON sp.user_id = u.id
		WHERE sp.student_id=$1
	`

	ListDocuments = `
		SELECT * FROM documents
		WHERE school_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
)
