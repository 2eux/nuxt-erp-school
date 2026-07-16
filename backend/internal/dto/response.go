package dto

import (
	"time"

	"github.com/opencode/erp-school-backend/internal/domain"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Items interface{} `json:"items"`
	Meta  Meta        `json:"meta"`
}

func NewAPIResponse(code int, message string, data interface{}) APIResponse {
	return APIResponse{
		Success: code >= 200 && code < 300,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(code int, message string, err string) APIResponse {
	return APIResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error:   err,
	}
}

func NewPaginatedResponse(items interface{}, total int64, page, pageSize int) PaginatedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	return PaginatedResponse{
		Items: items,
		Meta: Meta{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserBrief `json:"user"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type UserBrief struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Username  string   `json:"username"`
	FullName  string   `json:"full_name"`
	AvatarURL string   `json:"avatar_url"`
	Phone     string   `json:"phone"`
	Roles     []string `json:"roles"`
	SchoolID  string   `json:"school_id"`
}

type UserDetail struct {
	ID               string     `json:"id"`
	SchoolID         string     `json:"school_id"`
	Email            string     `json:"email"`
	Username         string     `json:"username"`
	FullName         string     `json:"full_name"`
	AvatarURL        string     `json:"avatar_url"`
	Phone               string     `json:"phone"`
	IsActive            bool       `json:"is_active"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	PasswordChangedAt   time.Time  `json:"password_changed_at"`
	Roles               []RoleBrief `json:"roles"`
	Permissions         []string    `json:"permissions"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type RoleBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type RoleDetail struct {
	ID          string           `json:"id"`
	SchoolID    string           `json:"school_id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	IsSystem    bool             `json:"is_system"`
	Permissions []PermissionBrief `json:"permissions"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type PermissionBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type SchoolResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	NPSN           string    `json:"npsn"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	Province       string    `json:"province"`
	PostalCode     string    `json:"postal_code"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Website        string    `json:"website"`
	LogoURL        string    `json:"logo_url"`
	Type           string    `json:"type"`
	Accreditation  string    `json:"accreditation"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AcademicYearResponse struct {
	ID        string    `json:"id"`
	SchoolID  string    `json:"school_id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SemesterResponse struct {
	ID             string    `json:"id"`
	AcademicYearID string    `json:"academic_year_id"`
	Name           string    `json:"name"`
	SemesterNumber int       `json:"semester_number"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	IsActive       bool      `json:"is_active"`
}

type ClassResponse struct {
	ID                string          `json:"id"`
	SchoolID          string          `json:"school_id"`
	GradeID           string          `json:"grade_id"`
	GradeName         string          `json:"grade_name"`
	Name              string          `json:"name"`
	Capacity          int             `json:"capacity"`
	StudentCount      int             `json:"student_count"`
	HomeroomTeacherID *string         `json:"homeroom_teacher_id,omitempty"`
	HomeroomTeacher   *TeacherBrief   `json:"homeroom_teacher,omitempty"`
	AcademicYearID    string          `json:"academic_year_id"`
	AcademicYear      string          `json:"academic_year"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type TeacherBrief struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	NIP      string `json:"nip"`
}

type SubjectResponse struct {
	ID          string    `json:"id"`
	SchoolID    string    `json:"school_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	KKM         float64   `json:"kkm"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GradeResponse struct {
	ID        string    `json:"id"`
	SchoolID  string    `json:"school_id"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CurriculumResponse struct {
	ID         string    `json:"id"`
	SchoolID   string    `json:"school_id"`
	GradeID    string    `json:"grade_id"`
	SubjectID  string    `json:"subject_id"`
	SemesterID string    `json:"semester_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type StudentResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	SchoolID       string    `json:"school_id"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
	NIS            string    `json:"nis"`
	NISN           string    `json:"nisn"`
	NIK            string    `json:"nik"`
	ClassID        string    `json:"class_id"`
	ClassName      string    `json:"class_name"`
	AcademicYearID string    `json:"academic_year_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	Gender         string    `json:"gender"`
	PlaceOfBirth   string    `json:"place_of_birth"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StudentDetail struct {
	StudentResponse
	Parents []StudentParentResponse `json:"parents"`
}

type StudentParentResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Relation    string    `json:"relation"`
	IsPrimary   bool      `json:"is_primary"`
	Occupation  string    `json:"occupation"`
	Institution string    `json:"institution"`
	CreatedAt   time.Time `json:"created_at"`
}

type TeacherResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	SchoolID       string    `json:"school_id"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
	NIP            string    `json:"nip"`
	NIK            string    `json:"nik"`
	NUPTK          string    `json:"nupk"`
	Status         string    `json:"status"`
	JoinDate       time.Time `json:"join_date"`
	EducationLevel string    `json:"education_level"`
	Major          string    `json:"major"`
	Phone          string    `json:"phone"`
	SubjectsCount  int       `json:"subjects_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type EmployeeResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	SchoolID   string    `json:"school_id"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	NIP        string    `json:"nip"`
	NIK        string    `json:"nik"`
	Position   string    `json:"position"`
	Department string    `json:"department"`
	JoinDate   time.Time `json:"join_date"`
	Status     string    `json:"status"`
	BaseSalary float64   `json:"base_salary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ScheduleResponse struct {
	ID         string          `json:"id"`
	ClassID    string          `json:"class_id"`
	ClassName  string          `json:"class_name"`
	SubjectID  string          `json:"subject_id"`
	SubjectName string         `json:"subject_name"`
	TeacherID  string          `json:"teacher_id"`
	TeacherName string         `json:"teacher_name"`
	Day        domain.DayOfWeek `json:"day"`
	StartTime  string          `json:"start_time"`
	EndTime    string          `json:"end_time"`
	Room       string          `json:"room"`
	SemesterID string          `json:"semester_id"`
}

type AttendanceResponse struct {
	ID           string    `json:"id"`
	StudentID    string    `json:"student_id"`
	StudentName  string    `json:"student_name"`
	ScheduleID   string    `json:"schedule_id"`
	SubjectName  string    `json:"subject_name"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
	Notes        string    `json:"notes"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type ExamResponse struct {
	ID          string    `json:"id"`
	SubjectID   string    `json:"subject_id"`
	SubjectName string    `json:"subject_name"`
	ClassID     string    `json:"class_id"`
	ClassName   string    `json:"class_name"`
	SemesterID  string    `json:"semester_id"`
	Name        string    `json:"name"`
	ExamType    string    `json:"exam_type"`
	Date        time.Time `json:"date"`
	Duration    int       `json:"duration"`
	TotalScore  float64   `json:"total_score"`
	CreatedAt   time.Time `json:"created_at"`
}

type ExamResultResponse struct {
	ID          string    `json:"id"`
	ExamID      string    `json:"exam_id"`
	ExamName    string    `json:"exam_name"`
	StudentID   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	Score       float64   `json:"score"`
	Grade       string    `json:"grade"`
	Notes       string    `json:"notes"`
}

type GradebookResponse struct {
	ID            string    `json:"id"`
	ClassID       string    `json:"class_id"`
	SubjectID     string    `json:"subject_id"`
	SubjectName   string    `json:"subject_name"`
	StudentID     string    `json:"student_id"`
	StudentName   string    `json:"student_name"`
	SemesterID    string    `json:"semester_id"`
	DailyScore    float64   `json:"daily_score"`
	MidScore      float64   `json:"mid_score"`
	FinalScore    float64   `json:"final_score"`
	PracticeScore float64   `json:"practice_score"`
	TotalScore    float64   `json:"total_score"`
	Attitude      string    `json:"attitude"`
	Notes         string    `json:"notes"`
}

type ReportCardResponse struct {
	ID               string     `json:"id"`
	StudentID        string     `json:"student_id"`
	StudentName      string     `json:"student_name"`
	SemesterID       string     `json:"semester_id"`
	ClassName        string     `json:"class_name"`
	AverageScore     float64    `json:"average_score"`
	Rank             int        `json:"rank"`
	AbsentCount      int        `json:"absent_count"`
	SickCount        int        `json:"sick_count"`
	PermitCount      int        `json:"permit_count"`
	HomeroomComment  string     `json:"homeroom_comment"`
	ParentSignature  bool       `json:"parent_signature"`
	PublishedAt      *time.Time `json:"published_at,omitempty"`
	Grades           []GradebookResponse `json:"grades"`
}

type AssignmentResponse struct {
	ID          string    `json:"id"`
	SubjectID   string    `json:"subject_id"`
	SubjectName string    `json:"subject_name"`
	ClassID     string    `json:"class_id"`
	ClassName   string    `json:"class_name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	MaxScore        float64   `json:"max_score"`
	SubmissionCount int      `json:"submission_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type AssignmentSubmissionResponse struct {
	ID           string     `json:"id"`
	AssignmentID string     `json:"assignment_id"`
	StudentID    string     `json:"student_id"`
	StudentName  string     `json:"student_name"`
	Content      string     `json:"content"`
	FileURL      string     `json:"file_url"`
	Score        *float64   `json:"score,omitempty"`
	Feedback     string     `json:"feedback"`
	SubmittedAt  time.Time  `json:"submitted_at"`
	GradedAt     *time.Time `json:"graded_at,omitempty"`
}

type LessonPlanResponse struct {
	ID         string     `json:"id"`
	SubjectID  string     `json:"subject_id"`
	ClassID    string     `json:"class_id"`
	Title      string     `json:"title"`
	Date       time.Time  `json:"date"`
	Objectives string     `json:"objectives"`
	Materials  string     `json:"materials"`
	Activities string     `json:"activities"`
	Assessment string     `json:"assessment"`
	Reflection string     `json:"reflection"`
	Status     string     `json:"status"`
	ApprovedBy *string    `json:"approved_by,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type TeachingJournalResponse struct {
	ID          string    `json:"id"`
	TeacherID   string    `json:"teacher_id"`
	ScheduleID  string    `json:"schedule_id"`
	SubjectName string    `json:"subject_name"`
	ClassName   string    `json:"class_name"`
	Date        time.Time `json:"date"`
	Material    string    `json:"material"`
	Method      string    `json:"method"`
	AttendCount int       `json:"attend_count"`
	Notes       string    `json:"notes"`
}

type TahfidzProgressResponse struct {
	ID          string    `json:"id"`
	StudentID   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	TeacherID   string    `json:"teacher_id"`
	TeacherName string    `json:"teacher_name"`
	Surah       string    `json:"surah"`
	StartAyah   int       `json:"start_ayah"`
	EndAyah     int       `json:"end_ayah"`
	Juz         int       `json:"juz"`
	Page        int       `json:"page"`
	Status      string    `json:"status"`
	Quality     string    `json:"quality"`
	Notes       string    `json:"notes"`
	Date        time.Time `json:"date"`
}

type TahfidzSummary struct {
	TotalSurahs    int `json:"total_surahs"`
	TotalJuz       int `json:"total_juz"`
	TotalAyahs     int `json:"total_ayahs"`
	MemorizingCount int `json:"memorizing_count"`
	MemorizedCount  int `json:"memorized_count"`
}

type MutabaahResponse struct {
	ID        string    `json:"id"`
	StudentID string    `json:"student_id"`
	Date      time.Time `json:"date"`
	Fajr      bool      `json:"fajr"`
	Dhuhr     bool      `json:"dhuhr"`
	Asr       bool      `json:"asr"`
	Maghrib   bool      `json:"maghrib"`
	Isha      bool      `json:"isha"`
	Tahajjud  bool      `json:"tahajjud"`
	Dhuha     bool      `json:"dhuha"`
	Sunnah    bool      `json:"sunnah"`
	QuranTilawah int    `json:"quran_tilawah" db:"quran_tilawah"`
	QuranHifdz   int    `json:"quran_hifdz" db:"quran_hifdz"`
	DzikrPagi    bool   `json:"dzikr_pagi" db:"dzikr_pagi"`
	DzikrPetang  bool   `json:"dzikr_petang" db:"dzikr_petang"`
	Shadaqah     bool   `json:"shadaqah" db:"shadaqah"`
	PuasaSunnah  bool   `json:"puasa_sunnah" db:"puasa_sunnah"`
	WudhuSebelumTidur bool `json:"wudhu_sebelum_tidur" db:"wudhu_sebelum_tidur"`
	BacaDoaTidur bool   `json:"baca_doa_tidur" db:"baca_doa_tidur"`
	Notes       string  `json:"notes"`
}

type PrayerAttendanceResponse struct {
	ID        string    `json:"id"`
	StudentID string    `json:"student_id"`
	Date      time.Time `json:"date"`
	Fajr      bool      `json:"fajr"`
	Dhuhr     bool      `json:"dhuhr"`
	Asr       bool      `json:"asr"`
	Maghrib   bool      `json:"maghrib"`
	Isha      bool      `json:"isha"`
	Notes     string    `json:"notes"`
}

type InvoiceResponse struct {
	ID          string     `json:"id"`
	SchoolID    string     `json:"school_id"`
	StudentID   string     `json:"student_id"`
	StudentName string     `json:"student_name"`
	InvoiceNo   string     `json:"invoice_no"`
	TotalAmount float64    `json:"total_amount"`
	PaidAmount  float64    `json:"paid_amount"`
	Status      string     `json:"status"`
	DueDate     time.Time  `json:"due_date"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	SemesterID  string     `json:"semester_id"`
	Notes       string     `json:"notes"`
	Items       []InvoiceItemResponse `json:"items"`
	CreatedAt   time.Time  `json:"created_at"`
}

type InvoiceItemResponse struct {
	ID        string    `json:"id"`
	FeeTypeID string    `json:"fee_type_id"`
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
}

type PaymentResponse struct {
	ID         string     `json:"id"`
	InvoiceID  string     `json:"invoice_id"`
	InvoiceNo  string     `json:"invoice_no"`
	PaymentNo  string     `json:"payment_no"`
	Amount     float64    `json:"amount"`
	Method     string     `json:"method"`
	Status     string     `json:"status"`
	ProofURL   string     `json:"proof_url"`
	PaidAt     time.Time  `json:"paid_at"`
	VerifiedBy *string    `json:"verified_by,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	Notes      string     `json:"notes"`
	CreatedAt  time.Time  `json:"created_at"`
}

type JournalResponse struct {
	ID          string                `json:"id"`
	SchoolID    string                `json:"school_id"`
	JournalNo   string                `json:"journal_no"`
	Date        time.Time             `json:"date"`
	Description string                `json:"description"`
	Status      string                `json:"status"`
	TotalDebit  float64               `json:"total_debit"`
	TotalCredit float64               `json:"total_credit"`
	Entries     []JournalEntryResponse `json:"entries"`
	CreatedBy   string                `json:"created_by"`
	CreatedAt   time.Time             `json:"created_at"`
}

type JournalEntryResponse struct {
	ID          string  `json:"id"`
	AccountCode string  `json:"account_code"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}

type FeeTypeResponse struct {
	ID          string    `json:"id"`
	SchoolID    string    `json:"school_id"`
	Name        string    `json:"name"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Frequency   string    `json:"frequency"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type GeneralLedgerResponse struct {
	ID          string    `json:"id"`
	AccountCode string    `json:"account_code"`
	AccountName string    `json:"account_name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Debit       float64   `json:"debit"`
	Credit      float64   `json:"credit"`
	Balance     float64   `json:"balance"`
}

type PayrollResponse struct {
	ID         string     `json:"id"`
	EmployeeID string     `json:"employee_id"`
	FullName   string     `json:"full_name"`
	Period     string     `json:"period"`
	BaseSalary float64    `json:"base_salary"`
	Allowance  float64    `json:"allowance"`
	Deduction  float64    `json:"deduction"`
	NetSalary  float64    `json:"net_salary"`
	Status     string     `json:"status"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type AssetResponse struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"school_id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	PurchaseDate     time.Time `json:"purchase_date"`
	PurchasePrice    float64   `json:"purchase_price"`
	CurrentValue     float64   `json:"current_value"`
	Location         string    `json:"location"`
	Condition        string    `json:"condition"`
	Status           string    `json:"status"`
	DepreciationRate float64   `json:"depreciation_rate"`
	CreatedAt        time.Time `json:"created_at"`
}

type InventoryItemResponse struct {
	ID        string    `json:"id"`
	SchoolID  string    `json:"school_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Unit      string    `json:"unit"`
	StockIn   int       `json:"stock_in"`
	StockOut  int       `json:"stock_out"`
	StockCurrent int    `json:"stock_current"`
	StockMin  int       `json:"stock_min"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

type LibraryBookResponse struct {
	ID          string    `json:"id"`
	SchoolID    string    `json:"school_id"`
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	PublishYear int       `json:"publish_year"`
	Category    string    `json:"category"`
	Language    string    `json:"language"`
	TotalCopies int       `json:"total_copies"`
	Available   int       `json:"available"`
	Location    string    `json:"location"`
	CoverURL    string    `json:"cover_url"`
}

type BookBorrowingResponse struct {
	ID           string     `json:"id"`
	BookID       string     `json:"book_id"`
	BookTitle    string     `json:"book_title"`
	BorrowerID   string     `json:"borrower_id"`
	BorrowerName string     `json:"borrower_name"`
	BorrowDate   time.Time  `json:"borrow_date"`
	DueDate      time.Time  `json:"due_date"`
	ReturnDate   *time.Time `json:"return_date,omitempty"`
	Status       string     `json:"status"`
	Fine         float64    `json:"fine"`
}

type MedicalRecordResponse struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`
	Date       time.Time `json:"date"`
	Diagnosis  string    `json:"diagnosis"`
	Treatment  string    `json:"treatment"`
	Medication string    `json:"medication"`
	DocType    string    `json:"doc_type"`
	Notes      string    `json:"notes"`
	CreatedBy  string    `json:"created_by"`
}

type CounselingSessionResponse struct {
	ID           string    `json:"id"`
	StudentID    string    `json:"student_id"`
	StudentName  string    `json:"student_name"`
	CounselorID  string    `json:"counselor_id"`
	CounselorName string   `json:"counselor_name"`
	Date         time.Time `json:"date"`
	Cause        string    `json:"cause"`
	Action       string    `json:"action"`
	Mediation    string    `json:"mediation"`
	Notes        string    `json:"notes"`
}

type AnnouncementResponse struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	TargetRole string     `json:"target_role"`
	Priority   string     `json:"priority"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date"`
	CreatedBy  string     `json:"created_by"`
	CreatedAt  time.Time  `json:"created_at"`
}

type NotificationResponse struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      string     `json:"type"`
	RefType   string     `json:"ref_type"`
	RefID     string     `json:"ref_id"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type DocumentResponse struct {
	ID        string    `json:"id"`
	SchoolID  string    `json:"school_id"`
	Title     string    `json:"title"`
	DocType   string    `json:"doc_type"`
	FileURL   string    `json:"file_url"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type MeetingResponse struct {
	ID        string                  `json:"id"`
	Title     string                  `json:"title"`
	Agenda    string                  `json:"agenda"`
	Date      time.Time               `json:"date"`
	StartTime string                  `json:"start_time"`
	EndTime   string                  `json:"end_time"`
	Location  string                  `json:"location"`
	Minutes   string                  `json:"minutes"`
	Status    string                  `json:"status"`
	CreatedBy string                  `json:"created_by"`
	Attendees []MeetingAttendeeResponse `json:"attendees"`
	CreatedAt time.Time               `json:"created_at"`
}

type MeetingAttendeeResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`
}

type EventResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventType   string    `json:"event_type"`
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Location    string    `json:"location"`
	CreatedBy   string    `json:"created_by"`
}

type TaskResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	AssignedTo  string     `json:"assigned_to"`
	CreatedBy   string     `json:"created_by"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type SettingResponse struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type AdmissionApplicantResponse struct {
	ID             string     `json:"id"`
	FullName       string     `json:"full_name"`
	Gender         string     `json:"gender"`
	PlaceOfBirth   string     `json:"place_of_birth"`
	DateOfBirth    time.Time  `json:"date_of_birth"`
	PreviousSchool string     `json:"previous_school"`
	GradeID        string     `json:"grade_id"`
	GradeName      string     `json:"grade_name"`
	RegistrationNo string     `json:"registration_no"`
	ParentName     string     `json:"parent_name"`
	ParentPhone    string     `json:"parent_phone"`
	ParentEmail    string     `json:"parent_email"`
	Status         string     `json:"status"`
	TestScore      *float64   `json:"test_score,omitempty"`
	InterviewScore *float64   `json:"interview_score,omitempty"`
	AcceptedAt     *time.Time `json:"accepted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type GraduationCandidateResponse struct {
	ID             string     `json:"id"`
	StudentID      string     `json:"student_id"`
	StudentName    string     `json:"student_name"`
	ClassID        string     `json:"class_id"`
	ClassName      string     `json:"class_name"`
	Status         string     `json:"status"`
	FinalGrade     float64    `json:"final_grade"`
	CertificateNo  string     `json:"certificate_no"`
	GraduatedAt    *time.Time `json:"graduated_at,omitempty"`
}

type LeaveRequestResponse struct {
	ID          string     `json:"id"`
	EmployeeID  string     `json:"employee_id"`
	FullName    string     `json:"full_name"`
	LeaveType   string     `json:"leave_type"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	ApprovedBy  *string    `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type HalaqahGroupResponse struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	TeacherID    string              `json:"teacher_id"`
	TeacherName  string              `json:"teacher_name"`
	Room         string              `json:"room"`
	Day          domain.DayOfWeek     `json:"day"`
	StartTime    string              `json:"start_time"`
	EndTime      string              `json:"end_time"`
	MaxMember    int                 `json:"max_member"`
	MemberCount  int                 `json:"member_count"`
	Members      []HalaqahMemberResponse `json:"members,omitempty"`
}

type HalaqahMemberResponse struct {
	ID        string    `json:"id"`
	StudentID string    `json:"student_id"`
	FullName  string    `json:"full_name"`
	JoinedAt  time.Time `json:"joined_at"`
}

type DashboardStats struct {
	TotalStudents    int64 `json:"total_students"`
	TotalTeachers    int64 `json:"total_teachers"`
	TotalEmployees   int64 `json:"total_employees"`
	TotalClasses     int64 `json:"total_classes"`
	TotalRevenue     float64 `json:"total_revenue"`
	OutstandingFees  float64 `json:"outstanding_fees"`
	AttendanceRate   float64 `json:"attendance_rate"`
}

type StudentAttendanceSummary struct {
	Present int `json:"present"`
	Absent  int `json:"absent"`
	Sick    int `json:"sick"`
	Permit  int `json:"permit"`
	Late    int `json:"late"`
}

type FinancialSummary struct {
	TotalRevenue  float64 `json:"total_revenue"`
	TotalExpense  float64 `json:"total_expense"`
	Balance        float64 `json:"balance"`
	Outstanding    float64 `json:"outstanding"`
	CollectionRate float64 `json:"collection_rate"`
}

type AIChatResponse struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
	Role           string `json:"role"`
	TokenCount     int    `json:"token_count"`
}

type AIStreamResponse struct {
	Delta string `json:"delta"`
}

type AIGeneratedResponse struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type AIAnalysisResponse struct {
	Type    string      `json:"type"`
	Summary string      `json:"summary"`
	Details interface{} `json:"details"`
}

type AIConversationResponse struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Model        string              `json:"model"`
	MessageCount int                 `json:"message_count"`
	Messages     []AIMessageResponse `json:"messages,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type AIMessageResponse struct {
	ID         string    `json:"id"`
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	TokenCount int       `json:"token_count"`
	CreatedAt  time.Time `json:"created_at"`
}
