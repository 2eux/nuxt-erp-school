package domain

import (
	"time"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "monday"
	Tuesday   DayOfWeek = "tuesday"
	Wednesday DayOfWeek = "wednesday"
	Thursday  DayOfWeek = "thursday"
	Friday    DayOfWeek = "friday"
	Saturday  DayOfWeek = "saturday"
	Sunday    DayOfWeek = "sunday"
)

type School struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	NPSN        string    `json:"npsn" db:"npsn"`
	Address     string    `json:"address" db:"address"`
	City        string    `json:"city" db:"city"`
	Province    string    `json:"province" db:"province"`
	PostalCode  string    `json:"postal_code" db:"postal_code"`
	Phone       string    `json:"phone" db:"phone"`
	Email       string    `json:"email" db:"email"`
	Website     string    `json:"website" db:"website"`
	LogoURL     string    `json:"logo_url" db:"logo_url"`
	Type        string    `json:"type" db:"type"`
	Accreditation string  `json:"accreditation" db:"accreditation"`
	EstablishedDate string `json:"established_date" db:"established_date"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type AcademicYear struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Name      string    `json:"name" db:"name"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Semester struct {
	ID            string    `json:"id" db:"id"`
	AcademicYearID string   `json:"academic_year_id" db:"academic_year_id"`
	Name          string    `json:"name" db:"name"`
	SemesterNumber int     `json:"semester_number" db:"semester_number"`
	StartDate     time.Time `json:"start_date" db:"start_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Grade struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Name      string    `json:"name" db:"name"`
	Level     int       `json:"level" db:"level"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Class struct {
	ID         string    `json:"id" db:"id"`
	SchoolID   string    `json:"school_id" db:"school_id"`
	GradeID    string    `json:"grade_id" db:"grade_id"`
	Name       string    `json:"name" db:"name"`
	Capacity   int       `json:"capacity" db:"capacity"`
	HomeroomTeacherID *string `json:"homeroom_teacher_id,omitempty" db:"homeroom_teacher_id"`
	AcademicYearID string  `json:"academic_year_id" db:"academic_year_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Subject struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	Code        string    `json:"code" db:"code"`
	Name        string    `json:"name" db:"name"`
	Category    string    `json:"category" db:"category"`
	Description string    `json:"description" db:"description"`
	KKM         float64   `json:"kkm" db:"kkm"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Curriculum struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	GradeID   string    `json:"grade_id" db:"grade_id"`
	SubjectID string    `json:"subject_id" db:"subject_id"`
	SemesterID string   `json:"semester_id" db:"semester_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Role struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	IsSystem    bool      `json:"is_system" db:"is_system"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Permission struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Module      string    `json:"module" db:"module"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type RolePermission struct {
	RoleID       string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

type User struct {
	ID             string    `json:"id" db:"id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	Email          string    `json:"email" db:"email"`
	Username       string    `json:"username" db:"username"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	FullName       string    `json:"full_name" db:"full_name"`
	AvatarURL      string    `json:"avatar_url" db:"avatar_url"`
	Phone          string    `json:"phone" db:"phone"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	PasswordChangedAt time.Time `json:"password_changed_at" db:"password_changed_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UserRole struct {
	UserID string `json:"user_id" db:"user_id"`
	RoleID string `json:"role_id" db:"role_id"`
}

type UserProfile struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	PlaceOfBirth string    `json:"place_of_birth" db:"place_of_birth"`
	DateOfBirth  time.Time `json:"date_of_birth" db:"date_of_birth"`
	Gender       Gender    `json:"gender" db:"gender"`
	Religion     string    `json:"religion" db:"religion"`
	Address      string    `json:"address" db:"address"`
	Bio          string    `json:"bio" db:"bio"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserSession struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"-" db:"refresh_token"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Student struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	NIS            string    `json:"nis" db:"nis"`
	NISN           string    `json:"nisn" db:"nisn"`
	NIK            string    `json:"nik" db:"nik"`
	ClassID        string    `json:"class_id" db:"class_id"`
	AcademicYearID string    `json:"academic_year_id" db:"academic_year_id"`
	EnrollmentDate time.Time `json:"enrollment_date" db:"enrollment_date"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type StudentParent struct {
	ID          string    `json:"id" db:"id"`
	StudentID   string    `json:"student_id" db:"student_id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Relation    string    `json:"relation" db:"relation"`
	IsPrimary   bool      `json:"is_primary" db:"is_primary"`
	Occupation  string    `json:"occupation" db:"occupation"`
	Institution string    `json:"institution" db:"institution"`
	Income      float64   `json:"income" db:"income"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type StudentDocument struct {
	ID         string    `json:"id" db:"id"`
	StudentID  string    `json:"student_id" db:"student_id"`
	Name       string    `json:"name" db:"name"`
	DocType    string    `json:"doc_type" db:"doc_type"`
	FileURL    string    `json:"file_url" db:"file_url"`
	Status     string    `json:"status" db:"status"`
	Notes      string    `json:"notes" db:"notes"`
	VerifiedAt *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	VerifiedBy *string   `json:"verified_by,omitempty" db:"verified_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Teacher struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	NIP            string    `json:"nip" db:"nip"`
	NIK            string    `json:"nik" db:"nik"`
	NUPTK          string    `json:"nupk" db:"nupk"`
	Status         string    `json:"status" db:"status"`
	JoinDate       time.Time `json:"join_date" db:"join_date"`
	EducationLevel string    `json:"education_level" db:"education_level"`
	Major          string    `json:"major" db:"major"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type TeacherSubject struct {
	ID        string    `json:"id" db:"id"`
	TeacherID string    `json:"teacher_id" db:"teacher_id"`
	SubjectID string    `json:"subject_id" db:"subject_id"`
	ClassID   string    `json:"class_id" db:"class_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Employee struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	SchoolID   string    `json:"school_id" db:"school_id"`
	NIP        string    `json:"nip" db:"nip"`
	NIK        string    `json:"nik" db:"nik"`
	Position   string    `json:"position" db:"position"`
	Department string    `json:"department" db:"department"`
	JoinDate   time.Time `json:"join_date" db:"join_date"`
	Status     string    `json:"status" db:"status"`
	BaseSalary float64   `json:"base_salary" db:"base_salary"`
	BankAccount string   `json:"bank_account" db:"bank_account"`
	BankName   string    `json:"bank_name" db:"bank_name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type EmployeeAttendance struct {
	ID         string    `json:"id" db:"id"`
	EmployeeID string    `json:"employee_id" db:"employee_id"`
	Date       time.Time `json:"date" db:"date"`
	CheckIn    time.Time `json:"check_in" db:"check_in"`
	CheckOut   *time.Time `json:"check_out,omitempty" db:"check_out"`
	Status     string    `json:"status" db:"status"`
	Notes      string    `json:"notes" db:"notes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type LeaveRequest struct {
	ID         string    `json:"id" db:"id"`
	EmployeeID string    `json:"employee_id" db:"employee_id"`
	LeaveType  string    `json:"leave_type" db:"leave_type"`
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    time.Time `json:"end_date" db:"end_date"`
	Reason     string    `json:"reason" db:"reason"`
	Status     string    `json:"status" db:"status"`
	ApprovedBy *string   `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	AttachmentURL *string `json:"attachment_url,omitempty" db:"attachment_url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Schedule struct {
	ID        string    `json:"id" db:"id"`
	ClassID   string    `json:"class_id" db:"class_id"`
	SubjectID string    `json:"subject_id" db:"subject_id"`
	TeacherID string    `json:"teacher_id" db:"teacher_id"`
	Day       DayOfWeek `json:"day" db:"day"`
	StartTime string    `json:"start_time" db:"start_time"`
	EndTime   string    `json:"end_time" db:"end_time"`
	Room      string    `json:"room" db:"room"`
	SemesterID string   `json:"semester_id" db:"semester_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Attendance struct {
	ID         string    `json:"id" db:"id"`
	StudentID  string    `json:"student_id" db:"student_id"`
	ScheduleID string    `json:"schedule_id" db:"schedule_id"`
	Date       time.Time `json:"date" db:"date"`
	Status     string    `json:"status" db:"status"`
	Notes      string    `json:"notes" db:"notes"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Exam struct {
	ID         string    `json:"id" db:"id"`
	SubjectID  string    `json:"subject_id" db:"subject_id"`
	ClassID    string    `json:"class_id" db:"class_id"`
	SemesterID string    `json:"semester_id" db:"semester_id"`
	Name       string    `json:"name" db:"name"`
	ExamType   string    `json:"exam_type" db:"exam_type"`
	Date       time.Time `json:"date" db:"date"`
	Duration   int       `json:"duration" db:"duration"`
	TotalScore float64   `json:"total_score" db:"total_score"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ExamResult struct {
	ID        string    `json:"id" db:"id"`
	ExamID    string    `json:"exam_id" db:"exam_id"`
	StudentID string    `json:"student_id" db:"student_id"`
	Score     float64   `json:"score" db:"score"`
	Grade     string    `json:"grade" db:"grade"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Gradebook struct {
	ID          string    `json:"id" db:"id"`
	ClassID     string    `json:"class_id" db:"class_id"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	SemesterID  string    `json:"semester_id" db:"semester_id"`
	StudentID   string    `json:"student_id" db:"student_id"`
	DailyScore  float64   `json:"daily_score" db:"daily_score"`
	MidScore    float64   `json:"mid_score" db:"mid_score"`
	FinalScore  float64   `json:"final_score" db:"final_score"`
	PracticeScore float64 `json:"practice_score" db:"practice_score"`
	Attitude    string    `json:"attitude" db:"attitude"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ReportCard struct {
	ID             string    `json:"id" db:"id"`
	StudentID      string    `json:"student_id" db:"student_id"`
	SemesterID     string    `json:"semester_id" db:"semester_id"`
	ClassID        string    `json:"class_id" db:"class_id"`
	AverageScore   float64   `json:"average_score" db:"average_score"`
	Rank           int       `json:"rank" db:"rank"`
	AbsentCount    int       `json:"absent_count" db:"absent_count"`
	SickCount      int       `json:"sick_count" db:"sick_count"`
	PermitCount    int       `json:"permit_count" db:"permit_count"`
	HomeroomComment string   `json:"homeroom_comment" db:"homeroom_comment"`
	ParentSignature bool     `json:"parent_signature" db:"parent_signature"`
	ApprovedBy     string    `json:"approved_by" db:"approved_by"`
	ApprovedAt     time.Time `json:"approved_at" db:"approved_at"`
	PublishedAt    *time.Time `json:"published_at,omitempty" db:"published_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type Assignment struct {
	ID          string    `json:"id" db:"id"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	ClassID     string    `json:"class_id" db:"class_id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	MaxScore    float64   `json:"max_score" db:"max_score"`
	Attachments string    `json:"attachments" db:"attachments"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type AssignmentSubmission struct {
	ID           string     `json:"id" db:"id"`
	AssignmentID string     `json:"assignment_id" db:"assignment_id"`
	StudentID    string     `json:"student_id" db:"student_id"`
	Content      string     `json:"content" db:"content"`
	FileURL      string     `json:"file_url" db:"file_url"`
	Score        *float64   `json:"score,omitempty" db:"score"`
	Feedback     string     `json:"feedback" db:"feedback"`
	SubmittedAt  time.Time  `json:"submitted_at" db:"submitted_at"`
	GradedAt     *time.Time `json:"graded_at,omitempty" db:"graded_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type LessonPlan struct {
	ID          string    `json:"id" db:"id"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	ClassID     string    `json:"class_id" db:"class_id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	Title       string    `json:"title" db:"title"`
	Date        time.Time `json:"date" db:"date"`
	Objectives  string    `json:"objectives" db:"objectives"`
	Materials   string    `json:"materials" db:"materials"`
	Activities  string    `json:"activities" db:"activities"`
	Assessment  string    `json:"assessment" db:"assessment"`
	Reflection  string    `json:"reflection" db:"reflection"`
	Status      string    `json:"status" db:"status"`
	ApprovedBy  *string   `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TeachingJournal struct {
	ID          string    `json:"id" db:"id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	ScheduleID  string    `json:"schedule_id" db:"schedule_id"`
	Date        time.Time `json:"date" db:"date"`
	Material    string    `json:"material" db:"material"`
	Method      string    `json:"method" db:"method"`
	AttendCount int       `json:"attend_count" db:"attend_count"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TahfidzProgress struct {
	ID        string    `json:"id" db:"id"`
	StudentID string    `json:"student_id" db:"student_id"`
	TeacherID string    `json:"teacher_id" db:"teacher_id"`
	Surah     string    `json:"surah" db:"surah"`
	StartAyah int       `json:"start_ayah" db:"start_ayah"`
	EndAyah   int       `json:"end_ayah" db:"end_ayah"`
	Juz       int       `json:"juz" db:"juz"`
	Page      int       `json:"page" db:"page"`
	Status    string    `json:"status" db:"status"`
	Quality   string    `json:"quality" db:"quality"`
	Notes     string    `json:"notes" db:"notes"`
	Date      time.Time `json:"date" db:"date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type MutabaahYaumiyah struct {
	ID        string    `json:"id" db:"id"`
	StudentID string    `json:"student_id" db:"student_id"`
	Date      time.Time `json:"date" db:"date"`
	Fajr      bool      `json:"fajr" db:"fajr"`
	Dhuhr     bool      `json:"dhuhr" db:"dhuhr"`
	Asr       bool      `json:"asr" db:"asr"`
	Maghrib   bool      `json:"maghrib" db:"maghrib"`
	Isha      bool      `json:"isha" db:"isha"`
	Tahajjud  bool      `json:"tahajjud" db:"tahajjud"`
	Dhuha     bool      `json:"dhuha" db:"dhuha"`
	Sunnah    bool      `json:"sunnah" db:"sunnah"`
	QuranTilawah int    `json:"quran_tilawah" db:"quran_tilawah"`
	QuranHifdz   int    `json:"quran_hifdz" db:"quran_hifdz"`
	DzikrPagi    bool   `json:"dzikr_pagi" db:"dzikr_pagi"`
	DzikrPetang  bool   `json:"dzikr_petang" db:"dzikr_petang"`
	Shadaqah     bool   `json:"shadaqah" db:"shadaqah"`
	PuasaSunnah  bool   `json:"puasa_sunnah" db:"puasa_sunnah"`
	WudhuSebelumTidur bool `json:"wudhu_sebelum_tidur" db:"wudhu_sebelum_tidur"`
	BacaDoaTidur bool   `json:"baca_doa_tidur" db:"baca_doa_tidur"`
	Notes       string  `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PrayerAttendance struct {
	ID        string    `json:"id" db:"id"`
	StudentID string    `json:"student_id" db:"student_id"`
	Date      time.Time `json:"date" db:"date"`
	Fajr      bool      `json:"fajr" db:"fajr"`
	Dhuhr     bool      `json:"dhuhr" db:"dhuhr"`
	Asr       bool      `json:"asr" db:"asr"`
	Maghrib   bool      `json:"maghrib" db:"maghrib"`
	Isha      bool      `json:"isha" db:"isha"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type HalaqahGroup struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Name      string    `json:"name" db:"name"`
	TeacherID string    `json:"teacher_id" db:"teacher_id"`
	Room      string    `json:"room" db:"room"`
	Day       DayOfWeek `json:"day" db:"day"`
	StartTime string    `json:"start_time" db:"start_time"`
	EndTime   string    `json:"end_time" db:"end_time"`
	MaxMember int       `json:"max_member" db:"max_member"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type HalaqahMember struct {
	ID          string    `json:"id" db:"id"`
	HalaqahID   string    `json:"halaqah_id" db:"halaqah_id"`
	StudentID   string    `json:"student_id" db:"student_id"`
	JoinedAt    time.Time `json:"joined_at" db:"joined_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type FeeType struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	Name        string    `json:"name" db:"name"`
	Amount      float64   `json:"amount" db:"amount"`
	Category    string    `json:"category" db:"category"`
	Frequency   string    `json:"frequency" db:"frequency"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Invoice struct {
	ID          string     `json:"id" db:"id"`
	SchoolID    string     `json:"school_id" db:"school_id"`
	StudentID   string     `json:"student_id" db:"student_id"`
	InvoiceNo   string     `json:"invoice_no" db:"invoice_no"`
	TotalAmount float64    `json:"total_amount" db:"total_amount"`
	PaidAmount  float64    `json:"paid_amount" db:"paid_amount"`
	Status      string     `json:"status" db:"status"`
	DueDate     time.Time  `json:"due_date" db:"due_date"`
	PaidAt      *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	SemesterID  string     `json:"semester_id" db:"semester_id"`
	Notes       string     `json:"notes" db:"notes"`
	CreatedBy   string     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type InvoiceItem struct {
	ID        string    `json:"id" db:"id"`
	InvoiceID string    `json:"invoice_id" db:"invoice_id"`
	FeeTypeID string    `json:"fee_type_id" db:"fee_type_id"`
	Name      string    `json:"name" db:"name"`
	Amount    float64   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Payment struct {
	ID          string    `json:"id" db:"id"`
	InvoiceID   string    `json:"invoice_id" db:"invoice_id"`
	PaymentNo   string    `json:"payment_no" db:"payment_no"`
	Amount      float64   `json:"amount" db:"amount"`
	Method      string    `json:"method" db:"method"`
	Status      string    `json:"status" db:"status"`
	ProofURL    string    `json:"proof_url" db:"proof_url"`
	PaidAt      time.Time `json:"paid_at" db:"paid_at"`
	VerifiedBy  *string   `json:"verified_by,omitempty" db:"verified_by"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Journal struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	JournalNo   string    `json:"journal_no" db:"journal_no"`
	Date        time.Time `json:"date" db:"date"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	ApprovedBy  *string   `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type JournalEntry struct {
	ID          string    `json:"id" db:"id"`
	JournalID   string    `json:"journal_id" db:"journal_id"`
	AccountCode string    `json:"account_code" db:"account_code"`
	Description string    `json:"description" db:"description"`
	Debit       float64   `json:"debit" db:"debit"`
	Credit      float64   `json:"credit" db:"credit"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type GeneralLedger struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	AccountCode string    `json:"account_code" db:"account_code"`
	AccountName string    `json:"account_name" db:"account_name"`
	Date        time.Time `json:"date" db:"date"`
	Description string    `json:"description" db:"description"`
	Debit       float64   `json:"debit" db:"debit"`
	Credit      float64   `json:"credit" db:"credit"`
	Balance     float64   `json:"balance" db:"balance"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type PayrollDetail struct {
	ID         string    `json:"id" db:"id"`
	EmployeeID string    `json:"employee_id" db:"employee_id"`
	SchoolID   string    `json:"school_id" db:"school_id"`
	Period     string    `json:"period" db:"period"`
	BaseSalary float64   `json:"base_salary" db:"base_salary"`
	Allowance  float64   `json:"allowance" db:"allowance"`
	Deduction  float64   `json:"deduction" db:"deduction"`
	NetSalary  float64   `json:"net_salary" db:"net_salary"`
	Status     string    `json:"status" db:"status"`
	PaidAt     *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	Notes      string    `json:"notes" db:"notes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Asset struct {
	ID             string    `json:"id" db:"id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	Code           string    `json:"code" db:"code"`
	Name           string    `json:"name" db:"name"`
	Category       string    `json:"category" db:"category"`
	PurchaseDate   time.Time `json:"purchase_date" db:"purchase_date"`
	PurchasePrice  float64   `json:"purchase_price" db:"purchase_price"`
	CurrentValue   float64   `json:"current_value" db:"current_value"`
	Location       string    `json:"location" db:"location"`
	Condition      string    `json:"condition" db:"condition"`
	Status         string    `json:"status" db:"status"`
	DepreciationRate float64 `json:"depreciation_rate" db:"depreciation_rate"`
	ResponsibleID  string    `json:"responsible_id" db:"responsible_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type InventoryItem struct {
	ID         string    `json:"id" db:"id"`
	SchoolID   string    `json:"school_id" db:"school_id"`
	Code       string    `json:"code" db:"code"`
	Name       string    `json:"name" db:"name"`
	Category   string    `json:"category" db:"category"`
	Unit       string    `json:"unit" db:"unit"`
	StockIn    int       `json:"stock_in" db:"stock_in"`
	StockOut   int       `json:"stock_out" db:"stock_out"`
	StockMin   int       `json:"stock_min" db:"stock_min"`
	Location   string    `json:"location" db:"location"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type StockMovement struct {
	ID             string    `json:"id" db:"id"`
	ItemID         string    `json:"item_id" db:"item_id"`
	Type           string    `json:"type" db:"type"`
	Quantity       int       `json:"quantity" db:"quantity"`
	ReferenceType  string    `json:"reference_type" db:"reference_type"`
	ReferenceID    string    `json:"reference_id" db:"reference_id"`
	Notes          string    `json:"notes" db:"notes"`
	CreatedBy      string    `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type LibraryBook struct {
	ID           string    `json:"id" db:"id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	ISBN         string    `json:"isbn" db:"isbn"`
	Title        string    `json:"title" db:"title"`
	Author       string    `json:"author" db:"author"`
	Publisher    string    `json:"publisher" db:"publisher"`
	PublishYear  int       `json:"publish_year" db:"publish_year"`
	Category     string    `json:"category" db:"category"`
	Language     string    `json:"language" db:"language"`
	TotalCopies  int       `json:"total_copies" db:"total_copies"`
	Available    int       `json:"available" db:"available"`
	Location     string    `json:"location" db:"location"`
	CoverURL     string    `json:"cover_url" db:"cover_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type BookBorrowing struct {
	ID           string     `json:"id" db:"id"`
	BookID       string     `json:"book_id" db:"book_id"`
	BorrowerID   string     `json:"borrower_id" db:"borrower_id"`
	BorrowerType string     `json:"borrower_type" db:"borrower_type"`
	BorrowDate   time.Time  `json:"borrow_date" db:"borrow_date"`
	DueDate      time.Time  `json:"due_date" db:"due_date"`
	ReturnDate   *time.Time `json:"return_date,omitempty" db:"return_date"`
	Status       string     `json:"status" db:"status"`
	Fine         float64    `json:"fine" db:"fine"`
	FinePaid     bool       `json:"fine_paid" db:"fine_paid"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type MedicalRecord struct {
	ID             string    `json:"id" db:"id"`
	StudentID      string    `json:"student_id" db:"student_id"`
	Date           time.Time `json:"date" db:"date"`
	Diagnosis      string    `json:"diagnosis" db:"diagnosis"`
	Treatment      string    `json:"treatment" db:"treatment"`
	Medication     string    `json:"medication" db:"medication"`
	DocType        string    `json:"doc_type" db:"doc_type"`
	Notes          string    `json:"notes" db:"notes"`
	CreatedBy      string    `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CounselingSession struct {
	ID          string    `json:"id" db:"id"`
	StudentID   string    `json:"student_id" db:"student_id"`
	CounselorID string    `json:"counselor_id" db:"counselor_id"`
	Date        time.Time `json:"date" db:"date"`
	Cause       string    `json:"cause" db:"cause"`
	Action      string    `json:"action" db:"action"`
	Mediation   string    `json:"mediation" db:"mediation"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Announcement struct {
	ID         string     `json:"id" db:"id"`
	SchoolID   string     `json:"school_id" db:"school_id"`
	Title      string     `json:"title" db:"title"`
	Content    string     `json:"content" db:"content"`
	TargetRole string     `json:"target_role" db:"target_role"`
	Priority   string     `json:"priority" db:"priority"`
	StartDate  time.Time  `json:"start_date" db:"start_date"`
	EndDate    time.Time  `json:"end_date" db:"end_date"`
	CreatedBy  string     `json:"created_by" db:"created_by"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Notification struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Message   string    `json:"message" db:"message"`
	Type      string    `json:"type" db:"type"`
	RefType   string    `json:"ref_type" db:"ref_type"`
	RefID     string    `json:"ref_id" db:"ref_id"`
	IsRead    bool      `json:"is_read" db:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty" db:"read_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Document struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Title     string    `json:"title" db:"title"`
	DocType   string    `json:"doc_type" db:"doc_type"`
	FileURL   string    `json:"file_url" db:"file_url"`
	FileSize  int64     `json:"file_size" db:"file_size"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Letter struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	LetterNo    string    `json:"letter_no" db:"letter_no"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	LetterType  string    `json:"letter_type" db:"letter_type"`
	Status      string    `json:"status" db:"status"`
	From        string    `json:"from" db:"from"`
	To          string    `json:"to" db:"to"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	ApprovedBy  *string   `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Meeting struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Title     string    `json:"title" db:"title"`
	Agenda    string    `json:"agenda" db:"agenda"`
	Date      time.Time `json:"date" db:"date"`
	StartTime string    `json:"start_time" db:"start_time"`
	EndTime   string    `json:"end_time" db:"end_time"`
	Location  string    `json:"location" db:"location"`
	Minutes   string    `json:"minutes" db:"minutes"`
	Status    string    `json:"status" db:"status"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type MeetingAttendee struct {
	ID        string    `json:"id" db:"id"`
	MeetingID string    `json:"meeting_id" db:"meeting_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Event struct {
	ID          string    `json:"id" db:"id"`
	SchoolID    string    `json:"school_id" db:"school_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	EventType   string    `json:"event_type" db:"event_type"`
	Date        time.Time `json:"date" db:"date"`
	StartTime   string    `json:"start_time" db:"start_time"`
	EndTime     string    `json:"end_time" db:"end_time"`
	Location    string    `json:"location" db:"location"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Task struct {
	ID          string     `json:"id" db:"id"`
	SchoolID    string     `json:"school_id" db:"school_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	AssignedTo  string     `json:"assigned_to" db:"assigned_to"`
	CreatedBy   string     `json:"created_by" db:"created_by"`
	Status      string     `json:"status" db:"status"`
	Priority    string     `json:"priority" db:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type Setting struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	Type      string    `json:"type" db:"type"`
	Module    string    `json:"module" db:"module"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AuditLog struct {
	ID         string    `json:"id" db:"id"`
	SchoolID   string    `json:"school_id" db:"school_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Action     string    `json:"action" db:"action"`
	Entity     string    `json:"entity" db:"entity"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	OldValues  string    `json:"old_values" db:"old_values"`
	NewValues  string    `json:"new_values" db:"new_values"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type ActivityLog struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Action    string    `json:"action" db:"action"`
	Module    string    `json:"module" db:"module"`
	Detail    string    `json:"detail" db:"detail"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type AdmissionApplicant struct {
	ID               string    `json:"id" db:"id"`
	SchoolID         string    `json:"school_id" db:"school_id"`
	FullName         string    `json:"full_name" db:"full_name"`
	Gender           Gender    `json:"gender" db:"gender"`
	PlaceOfBirth     string    `json:"place_of_birth" db:"place_of_birth"`
	DateOfBirth      time.Time `json:"date_of_birth" db:"date_of_birth"`
	PreviousSchool   string    `json:"previous_school" db:"previous_school"`
	GradeID          string    `json:"grade_id" db:"grade_id"`
	RegistrationNo   string    `json:"registration_no" db:"registration_no"`
	ParentName       string    `json:"parent_name" db:"parent_name"`
	ParentPhone      string    `json:"parent_phone" db:"parent_phone"`
	ParentEmail      string    `json:"parent_email" db:"parent_email"`
	Status           string    `json:"status" db:"status"`
	TestScore        *float64  `json:"test_score,omitempty" db:"test_score"`
	InterviewScore   *float64  `json:"interview_score,omitempty" db:"interview_score"`
	AcceptedAt       *time.Time `json:"accepted_at,omitempty" db:"accepted_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type GraduationCandidate struct {
	ID          string    `json:"id" db:"id"`
	StudentID   string    `json:"student_id" db:"student_id"`
	ClassID     string    `json:"class_id" db:"class_id"`
	SemesterID  string    `json:"semester_id" db:"semester_id"`
	Status      string    `json:"status" db:"status"`
	FinalGrade  float64   `json:"final_grade" db:"final_grade"`
	CertificateNo string  `json:"certificate_no" db:"certificate_no"`
	CertificateURL string  `json:"certificate_url" db:"certificate_url"`
	GraduatedAt *time.Time `json:"graduated_at,omitempty" db:"graduated_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type KnowledgeDocument struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	DocType   string    `json:"doc_type" db:"doc_type"`
	Module    string    `json:"module" db:"module"`
	ChunkIDs  string    `json:"chunk_ids" db:"chunk_ids"`
	EmbeddingStatus string `json:"embedding_status" db:"embedding_status"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AIConversation struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Title     string    `json:"title" db:"title"`
	Model     string    `json:"model" db:"model"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AIMessage struct {
	ID             string    `json:"id" db:"id"`
	ConversationID string    `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"`
	Content        string    `json:"content" db:"content"`
	TokenCount     int       `json:"token_count" db:"token_count"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type UserPermission struct {
	UserID         string `json:"user_id" db:"user_id"`
	PermissionSlug string `json:"permission_slug" db:"permission_slug"`
}
