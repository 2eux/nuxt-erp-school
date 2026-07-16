package dto

import (
	"time"

	"github.com/opencode/erp-school-backend/internal/domain"
)

type PaginationRequest struct {
	Page     int    `json:"page" form:"page" validate:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" validate:"min=1,max=100"`
	Search   string `json:"search" form:"search"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortDir  string `json:"sort_dir" form:"sort_dir" validate:"omitempty,oneof=asc desc"`
}

func (p *PaginationRequest) Defaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	if p.SortDir == "" {
		p.SortDir = "asc"
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
}

func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required" validate:"required,min=8"`
}

type CreateSchoolRequest struct {
	Name           string `json:"name" binding:"required" validate:"required,max=255"`
	NPSN           string `json:"npsn" validate:"max=20"`
	Address        string `json:"address" validate:"max=500"`
	City           string `json:"city" validate:"max=100"`
	Province       string `json:"province" validate:"max=100"`
	PostalCode     string `json:"postal_code" validate:"max=10"`
	Phone          string `json:"phone" validate:"max=20"`
	Email          string `json:"email" validate:"email,max=255"`
	Website        string `json:"website" validate:"max=255"`
	Type           string `json:"type" validate:"max=50"`
	Accreditation  string `json:"accreditation" validate:"max=20"`
	EstablishedDate string `json:"established_date"`
}

type UpdateSchoolRequest struct {
	Name           string `json:"name" validate:"max=255"`
	NPSN           string `json:"npsn" validate:"max=20"`
	Address        string `json:"address" validate:"max=500"`
	City           string `json:"city" validate:"max=100"`
	Province       string `json:"province" validate:"max=100"`
	PostalCode     string `json:"postal_code" validate:"max=10"`
	Phone          string `json:"phone" validate:"max=20"`
	Email          string `json:"email" validate:"email,max=255"`
	Website        string `json:"website" validate:"max=255"`
	Type           string `json:"type" validate:"max=50"`
	Accreditation  string `json:"accreditation" validate:"max=20"`
	IsActive       *bool  `json:"is_active"`
}

type CreateAcademicYearRequest struct {
	Name      string    `json:"name" binding:"required" validate:"required,max=100"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type UpdateAcademicYearRequest struct {
	Name      string     `json:"name" validate:"max=100"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsActive  *bool      `json:"is_active"`
}

type CreateSemesterRequest struct {
	AcademicYearID string    `json:"academic_year_id" binding:"required" validate:"required,uuid"`
	Name           string    `json:"name" binding:"required" validate:"required,max=50"`
	SemesterNumber int       `json:"semester_number" binding:"required" validate:"required,min=1,max=2"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
}

type CreateClassRequest struct {
	GradeID           string  `json:"grade_id" binding:"required" validate:"required,uuid"`
	Name              string  `json:"name" binding:"required" validate:"required,max=100"`
	Capacity          int     `json:"capacity" validate:"min=1,max=100"`
	HomeroomTeacherID *string `json:"homeroom_teacher_id" validate:"omitempty,uuid"`
	AcademicYearID    string  `json:"academic_year_id" binding:"required" validate:"required,uuid"`
}

type UpdateClassRequest struct {
	Name              string  `json:"name" validate:"max=100"`
	Capacity          *int    `json:"capacity" validate:"min=1,max=100"`
	HomeroomTeacherID *string `json:"homeroom_teacher_id" validate:"omitempty,uuid"`
}

type CreateGradeRequest struct {
	Name  string `json:"name" binding:"required" validate:"required,max=100"`
	Level int    `json:"level" binding:"required" validate:"required,min=1,max=12"`
}

type CreateSubjectRequest struct {
	Code        string  `json:"code" binding:"required" validate:"required,max=20"`
	Name        string  `json:"name" binding:"required" validate:"required,max=255"`
	Category    string  `json:"category" binding:"required" validate:"required,max=50"`
	Description string  `json:"description"`
	KKM         float64 `json:"kkm" validate:"min=0,max=100"`
}

type UpdateSubjectRequest struct {
	Code        string  `json:"code" validate:"max=20"`
	Name        string  `json:"name" validate:"max=255"`
	Category    string  `json:"category" validate:"max=50"`
	Description string  `json:"description"`
	KKM         *float64 `json:"kkm" validate:"min=0,max=100"`
	IsActive    *bool   `json:"is_active"`
}

type CreateUserRequest struct {
	Email    string  `json:"email" binding:"required" validate:"required,email,max=255"`
	Username string  `json:"username" binding:"required" validate:"required,min=3,max=50"`
	Password string  `json:"password" binding:"required" validate:"required,min=8"`
	FullName string  `json:"full_name" binding:"required" validate:"required,max=255"`
	Phone    string  `json:"phone" validate:"max=20"`
	RoleIDs  []string `json:"role_ids" validate:"required,min=1,dive,uuid"`
}

type UpdateUserRequest struct {
	Email    string   `json:"email" validate:"email,max=255"`
	Username string   `json:"username" validate:"min=3,max=50"`
	FullName string   `json:"full_name" validate:"max=255"`
	Phone    string   `json:"phone" validate:"max=20"`
	IsActive *bool    `json:"is_active"`
	RoleIDs  []string `json:"role_ids" validate:"dive,uuid"`
}

type CreateStudentRequest struct {
	UserID         string    `json:"user_id" validate:"uuid"`
	Email          string    `json:"email" binding:"required" validate:"required,email"`
	Password       string    `json:"password" binding:"required" validate:"required,min=8"`
	FullName       string    `json:"full_name" binding:"required" validate:"required,max=255"`
	NIS            string    `json:"nis" validate:"max=20"`
	NISN           string    `json:"nisn" validate:"max=20"`
	NIK            string    `json:"nik" validate:"max=20"`
	ClassID        string    `json:"class_id" binding:"required" validate:"required,uuid"`
	AcademicYearID string    `json:"academic_year_id" binding:"required" validate:"required,uuid"`
	EnrollmentDate time.Time `json:"enrollment_date" binding:"required"`
	Gender         domain.Gender `json:"gender" binding:"required" validate:"required,oneof=male female"`
	PlaceOfBirth   string    `json:"place_of_birth" validate:"max=100"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Address        string    `json:"address" validate:"max=500"`
}

type UpdateStudentRequest struct {
	NIS     string `json:"nis" validate:"max=20"`
	NISN    string `json:"nisn" validate:"max=20"`
	NIK     string `json:"nik" validate:"max=20"`
	ClassID string `json:"class_id" validate:"uuid"`
	Status  string `json:"status" validate:"oneof=active inactive graduated transferred"`
}

type CreateStudentParentRequest struct {
	UserID      string `json:"user_id" validate:"uuid"`
	Email       string `json:"email" binding:"required" validate:"required,email"`
	Password    string `json:"password" binding:"required" validate:"required,min=8"`
	FullName    string `json:"full_name" binding:"required" validate:"required,max=255"`
	Relation    string `json:"relation" binding:"required" validate:"required,oneof=father mother guardian"`
	IsPrimary   bool   `json:"is_primary"`
	Occupation  string `json:"occupation" validate:"max=100"`
	Institution string `json:"institution" validate:"max=255"`
	Income      float64 `json:"income"`
	Phone       string `json:"phone" validate:"max=20"`
}

type CreateExamRequest struct {
	SubjectID  string    `json:"subject_id" binding:"required" validate:"required,uuid"`
	ClassID    string    `json:"class_id" binding:"required" validate:"required,uuid"`
	SemesterID string    `json:"semester_id" binding:"required" validate:"required,uuid"`
	Name       string    `json:"name" binding:"required" validate:"required,max=255"`
	ExamType   string    `json:"exam_type" binding:"required" validate:"required,oneof=daily mid final practice"`
	Date       time.Time `json:"date" binding:"required"`
	Duration   int       `json:"duration" validate:"min=1"`
	TotalScore float64   `json:"total_score" validate:"min=0"`
}

type CreateExamResultRequest struct {
	ExamID    string  `json:"exam_id" binding:"required" validate:"required,uuid"`
	StudentID string  `json:"student_id" binding:"required" validate:"required,uuid"`
	Score     float64 `json:"score" binding:"required" validate:"min=0"`
	Grade     string  `json:"grade"`
	Notes     string  `json:"notes"`
}

type CreateGradebookRequest struct {
	ClassID       string  `json:"class_id" binding:"required" validate:"required,uuid"`
	SubjectID     string  `json:"subject_id" binding:"required" validate:"required,uuid"`
	SemesterID    string  `json:"semester_id" binding:"required" validate:"required,uuid"`
	StudentID     string  `json:"student_id" binding:"required" validate:"required,uuid"`
	DailyScore    float64 `json:"daily_score"`
	MidScore      float64 `json:"mid_score"`
	FinalScore    float64 `json:"final_score"`
	PracticeScore float64 `json:"practice_score"`
	Attitude      string  `json:"attitude"`
	Notes         string  `json:"notes"`
}

type UpdateGradebookRequest struct {
	DailyScore    *float64 `json:"daily_score"`
	MidScore      *float64 `json:"mid_score"`
	FinalScore    *float64 `json:"final_score"`
	PracticeScore *float64 `json:"practice_score"`
	Attitude      string   `json:"attitude"`
	Notes         string   `json:"notes"`
}

type CreateScheduleRequest struct {
	ClassID    string          `json:"class_id" binding:"required" validate:"required,uuid"`
	SubjectID  string          `json:"subject_id" binding:"required" validate:"required,uuid"`
	TeacherID  string          `json:"teacher_id" binding:"required" validate:"required,uuid"`
	Day        domain.DayOfWeek `json:"day" binding:"required" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime  string          `json:"start_time" binding:"required" validate:"required"`
	EndTime    string          `json:"end_time" binding:"required" validate:"required"`
	Room       string          `json:"room" validate:"max=50"`
	SemesterID string          `json:"semester_id" binding:"required" validate:"required,uuid"`
}

type CreateAttendanceRequest struct {
	StudentID  string    `json:"student_id" binding:"required" validate:"required,uuid"`
	ScheduleID string    `json:"schedule_id" binding:"required" validate:"required,uuid"`
	Date       time.Time `json:"date" binding:"required"`
	Status     string    `json:"status" binding:"required" validate:"required,oneof=present absent sick permit late"`
	Notes      string    `json:"notes"`
}

type CreateAssignmentRequest struct {
	SubjectID   string    `json:"subject_id" binding:"required" validate:"required,uuid"`
	ClassID     string    `json:"class_id" binding:"required" validate:"required,uuid"`
	Title       string    `json:"title" binding:"required" validate:"required,max=255"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	MaxScore    float64   `json:"max_score" validate:"min=0"`
}

type SubmitAssignmentRequest struct {
	Content string `json:"content"`
}

type CreateTahfidzProgressRequest struct {
	StudentID string    `json:"student_id" binding:"required" validate:"required,uuid"`
	Surah     string    `json:"surah" binding:"required" validate:"required"`
	StartAyah int       `json:"start_ayah" binding:"required" validate:"required,min=1"`
	EndAyah   int       `json:"end_ayah" binding:"required" validate:"required,min=1"`
	Juz       int       `json:"juz" validate:"min=1,max=30"`
	Page      int       `json:"page"`
	Status    string    `json:"status" binding:"required" validate:"required,oneof=memorizing memorized reviewing"`
	Quality   string    `json:"quality" validate:"oneof=excellent good fair poor"`
	Notes     string    `json:"notes"`
	Date      time.Time `json:"date" binding:"required"`
}

type CreateMutabaahRequest struct {
	StudentID         string    `json:"student_id" binding:"required" validate:"required,uuid"`
	Date              time.Time `json:"date" binding:"required"`
	Fajr              *bool     `json:"fajr"`
	Dhuhr             *bool     `json:"dhuhr"`
	Asr               *bool     `json:"asr"`
	Maghrib           *bool     `json:"maghrib"`
	Isha              *bool     `json:"isha"`
	Tahajjud          *bool     `json:"tahajjud"`
	Dhuha             *bool     `json:"dhuha"`
	Sunnah            *bool     `json:"sunnah"`
	QuranTilawah      *int      `json:"quran_tilawah"`
	QuranHifdz        *int      `json:"quran_hifdz"`
	DzikrPagi         *bool     `json:"dzikr_pagi"`
	DzikrPetang       *bool     `json:"dzikr_petang"`
	Shadaqah          *bool     `json:"shadaqah"`
	PuasaSunnah       *bool     `json:"puasa_sunnah"`
	WudhuSebelumTidur *bool     `json:"wudhu_sebelum_tidur"`
	BacaDoaTidur      *bool     `json:"baca_doa_tidur"`
	Notes             string    `json:"notes"`
}

type CreatePrayerAttendanceRequest struct {
	StudentID string    `json:"student_id" binding:"required" validate:"required,uuid"`
	Date      time.Time `json:"date" binding:"required"`
	Fajr      *bool     `json:"fajr"`
	Dhuhr     *bool     `json:"dhuhr"`
	Asr       *bool     `json:"asr"`
	Maghrib   *bool     `json:"maghrib"`
	Isha      *bool     `json:"isha"`
	Notes     string    `json:"notes"`
}

type CreateInvoiceRequest struct {
	StudentID   string              `json:"student_id" binding:"required" validate:"required,uuid"`
	SemesterID  string              `json:"semester_id" binding:"required" validate:"required,uuid"`
	DueDate     time.Time           `json:"due_date" binding:"required"`
	Notes       string              `json:"notes"`
	Items       []InvoiceItemRequest `json:"items" binding:"required,min=1,dive"`
}

type InvoiceItemRequest struct {
	FeeTypeID string  `json:"fee_type_id" binding:"required" validate:"required,uuid"`
	Name      string  `json:"name" binding:"required" validate:"required,max=255"`
	Amount    float64 `json:"amount" binding:"required" validate:"required,min=0"`
}

type CreatePaymentRequest struct {
	InvoiceID string  `json:"invoice_id" binding:"required" validate:"required,uuid"`
	Amount    float64 `json:"amount" binding:"required" validate:"required,min=1"`
	Method    string  `json:"method" binding:"required" validate:"required,oneof=cash transfer qris"`
	PaidAt    time.Time `json:"paid_at"`
	Notes     string  `json:"notes"`
}

type VerifyPaymentRequest struct {
	Status string `json:"status" binding:"required" validate:"required,oneof=verified rejected"`
	Notes  string `json:"notes"`
}

type CreateFeeTypeRequest struct {
	Name        string  `json:"name" binding:"required" validate:"required,max=255"`
	Amount      float64 `json:"amount" binding:"required" validate:"required,min=0"`
	Category    string  `json:"category" binding:"required" validate:"required,oneof=spp development activity uniform book other"`
	Frequency   string  `json:"frequency" binding:"required" validate:"required,oneof=monthly semester annually once"`
	Description string  `json:"description"`
}

type CreateJournalRequest struct {
	Date        time.Time              `json:"date" binding:"required"`
	Description string                 `json:"description" binding:"required" validate:"required"`
	Entries     []JournalEntryRequest  `json:"entries" binding:"required,min=1,dive"`
}

type JournalEntryRequest struct {
	AccountCode string  `json:"account_code" binding:"required" validate:"required,max=20"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit" validate:"min=0"`
	Credit      float64 `json:"credit" validate:"min=0"`
}

type CreateLessonPlanRequest struct {
	SubjectID  string    `json:"subject_id" binding:"required" validate:"required,uuid"`
	ClassID    string    `json:"class_id" binding:"required" validate:"required,uuid"`
	Title      string    `json:"title" binding:"required" validate:"required,max=255"`
	Date       time.Time `json:"date" binding:"required"`
	Objectives string    `json:"objectives" binding:"required"`
	Materials  string    `json:"materials" binding:"required"`
	Activities string    `json:"activities" binding:"required"`
	Assessment string    `json:"assessment"`
	Reflection string    `json:"reflection"`
}

type CreateTeachingJournalRequest struct {
	ScheduleID string    `json:"schedule_id" binding:"required" validate:"required,uuid"`
	Date       time.Time `json:"date" binding:"required"`
	Material   string    `json:"material" binding:"required"`
	Method     string    `json:"method" binding:"required"`
	AttendCount int      `json:"attend_count" validate:"min=0"`
	Notes      string    `json:"notes"`
}

type CreateAnnouncementRequest struct {
	Title      string    `json:"title" binding:"required" validate:"required,max=255"`
	Content    string    `json:"content" binding:"required"`
	TargetRole string    `json:"target_role" validate:"max=50"`
	Priority   string    `json:"priority" validate:"oneof=low medium high urgent"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

type CreateMeetingRequest struct {
	Title     string    `json:"title" binding:"required" validate:"required,max=255"`
	Agenda    string    `json:"agenda" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Location  string    `json:"location" validate:"max=255"`
	Attendees []string  `json:"attendees" validate:"required,min=1,dive,uuid"`
}

type AIChatRequest struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message" binding:"required" validate:"required"`
	Model          string `json:"model"`
}

type AIGenerateRequest struct {
	Type    string                 `json:"type" binding:"required" validate:"required,oneof=lesson_plan quiz exam letter meeting_minutes sop tahfidz_plan"`
	Context map[string]interface{} `json:"context"`
	Model   string                 `json:"model"`
}

type AIAnalyzeRequest struct {
	Type   string                 `json:"type" binding:"required" validate:"required,oneof=student_performance financial_health attendance_trend teacher_effectiveness"`
	Data   map[string]interface{} `json:"data" binding:"required"`
	Model  string                 `json:"model"`
}

type CreateEmployeeRequest struct {
	UserID      string    `json:"user_id" validate:"uuid"`
	Email       string    `json:"email" binding:"required" validate:"required,email"`
	Password    string    `json:"password" binding:"required" validate:"required,min=8"`
	FullName    string    `json:"full_name" binding:"required" validate:"required,max=255"`
	NIP         string    `json:"nip" validate:"max=30"`
	NIK         string    `json:"nik" validate:"max=20"`
	Position    string    `json:"position" binding:"required" validate:"required,max=100"`
	Department  string    `json:"department" binding:"required" validate:"required,max=100"`
	JoinDate    time.Time `json:"join_date" binding:"required"`
	BaseSalary  float64   `json:"base_salary" validate:"min=0"`
	BankAccount string    `json:"bank_account" validate:"max=50"`
	BankName    string    `json:"bank_name" validate:"max=100"`
}

type CreateTeacherRequest struct {
	UserID         string    `json:"user_id" validate:"uuid"`
	Email          string    `json:"email" binding:"required" validate:"required,email"`
	Password       string    `json:"password" binding:"required" validate:"required,min=8"`
	FullName       string    `json:"full_name" binding:"required" validate:"required,max=255"`
	NIP            string    `json:"nip" validate:"max=30"`
	NIK            string    `json:"nik" validate:"max=20"`
	NUPTK          string    `json:"nupk" validate:"max=20"`
	Status         string    `json:"status" binding:"required" validate:"required,oneof=permanent contract honorer"`
	JoinDate       time.Time `json:"join_date" binding:"required"`
	EducationLevel string    `json:"education_level" validate:"max=10"`
	Major          string    `json:"major" validate:"max=100"`
}

type CreateLeaveRequest struct {
	LeaveType string    `json:"leave_type" binding:"required" validate:"required,oneof=annual sick maternity personal other"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
}

type ApproveLeaveRequest struct {
	Status string `json:"status" binding:"required" validate:"required,oneof=approved rejected"`
}

type CreateHalaqahGroupRequest struct {
	Name      string          `json:"name" binding:"required" validate:"required,max=255"`
	TeacherID string          `json:"teacher_id" binding:"required" validate:"required,uuid"`
	Room      string          `json:"room" validate:"max=50"`
	Day       domain.DayOfWeek `json:"day" binding:"required" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime string          `json:"start_time" binding:"required"`
	EndTime   string          `json:"end_time" binding:"required"`
	MaxMember int             `json:"max_member" validate:"min=1,max=50"`
}

type AddHalaqahMemberRequest struct {
	StudentID string `json:"student_id" binding:"required" validate:"required,uuid"`
}

type CreateLibraryBookRequest struct {
	ISBN        string `json:"isbn" validate:"max=20"`
	Title       string `json:"title" binding:"required" validate:"required,max=255"`
	Author      string `json:"author" validate:"max=255"`
	Publisher   string `json:"publisher" validate:"max=255"`
	PublishYear int    `json:"publish_year"`
	Category    string `json:"category" validate:"max=100"`
	Language    string `json:"language" validate:"max=50"`
	TotalCopies int    `json:"total_copies" binding:"required" validate:"required,min=1"`
	Location    string `json:"location" validate:"max=100"`
}

type BorrowBookRequest struct {
	BorrowerID   string    `json:"borrower_id" binding:"required" validate:"required,uuid"`
	BorrowerType string    `json:"borrower_type" binding:"required" validate:"required,oneof=student teacher employee"`
	DueDate      time.Time `json:"due_date" binding:"required"`
}

type ReturnBookRequest struct {
	ReturnDate time.Time `json:"return_date" binding:"required"`
}

type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required" validate:"required,max=100"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions" binding:"required,min=1,dive,uuid"`
}

type UpdateRoleRequest struct {
	Name        string   `json:"name" validate:"max=100"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions" validate:"dive,uuid"`
}

type CreateReportCardRequest struct {
	StudentID        string  `json:"student_id" binding:"required" validate:"required,uuid"`
	SemesterID       string  `json:"semester_id" binding:"required" validate:"required,uuid"`
	ClassID          string  `json:"class_id" binding:"required" validate:"required,uuid"`
	AbsentCount      int     `json:"absent_count"`
	SickCount        int     `json:"sick_count"`
	PermitCount      int     `json:"permit_count"`
	HomeroomComment  string  `json:"homeroom_comment"`
}

type CreateAssetRequest struct {
	Code           string    `json:"code" binding:"required" validate:"required,max=50"`
	Name           string    `json:"name" binding:"required" validate:"required,max=255"`
	Category       string    `json:"category" binding:"required" validate:"required,max=100"`
	PurchaseDate   time.Time `json:"purchase_date" binding:"required"`
	PurchasePrice  float64   `json:"purchase_price" binding:"required" validate:"required,min=0"`
	Location       string    `json:"location" validate:"max=255"`
	DepreciationRate float64 `json:"depreciation_rate" validate:"min=0,max=100"`
	ResponsibleID  string    `json:"responsible_id" validate:"uuid"`
}

type CreateInventoryItemRequest struct {
	Code     string `json:"code" binding:"required" validate:"required,max=50"`
	Name     string `json:"name" binding:"required" validate:"required,max=255"`
	Category string `json:"category" validate:"max=100"`
	Unit     string `json:"unit" binding:"required" validate:"required,max=20"`
	StockMin int    `json:"stock_min" validate:"min=0"`
	Location string `json:"location" validate:"max=255"`
}

type StockMovementRequest struct {
	Type          string `json:"type" binding:"required" validate:"required,oneof=in out"`
	Quantity      int    `json:"quantity" binding:"required" validate:"required,min=1"`
	ReferenceType string `json:"reference_type" validate:"max=50"`
	ReferenceID   string `json:"reference_id" validate:"max=50"`
	Notes         string `json:"notes"`
}

type CreateMedicalRecordRequest struct {
	StudentID  string    `json:"student_id" binding:"required" validate:"required,uuid"`
	Date       time.Time `json:"date" binding:"required"`
	Diagnosis  string    `json:"diagnosis" binding:"required"`
	Treatment  string    `json:"treatment"`
	Medication string    `json:"medication"`
	DocType    string    `json:"doc_type"`
	Notes      string    `json:"notes"`
}

type CreateCounselingSessionRequest struct {
	StudentID string    `json:"student_id" binding:"required" validate:"required,uuid"`
	Date      time.Time `json:"date" binding:"required"`
	Cause     string    `json:"cause" binding:"required"`
	Action    string    `json:"action" binding:"required"`
	Mediation string    `json:"mediation"`
	Notes     string    `json:"notes"`
}

type CreateAdmissionRequest struct {
	FullName       string        `json:"full_name" binding:"required" validate:"required,max=255"`
	Gender         domain.Gender `json:"gender" binding:"required" validate:"required,oneof=male female"`
	PlaceOfBirth   string        `json:"place_of_birth" binding:"required" validate:"required,max=100"`
	DateOfBirth    time.Time     `json:"date_of_birth" binding:"required"`
	PreviousSchool string        `json:"previous_school" validate:"max=255"`
	GradeID        string        `json:"grade_id" binding:"required" validate:"required,uuid"`
	ParentName     string        `json:"parent_name" binding:"required" validate:"required,max=255"`
	ParentPhone    string        `json:"parent_phone" binding:"required" validate:"required,max=20"`
	ParentEmail    string        `json:"parent_email" validate:"email,max=255"`
}

type UpdateAdmissionRequest struct {
	Status         string   `json:"status" validate:"oneof=pending test interview accepted rejected"`
	TestScore      *float64 `json:"test_score"`
	InterviewScore *float64 `json:"interview_score"`
}

type CreatePayrollRequest struct {
	EmployeeID string  `json:"employee_id" binding:"required" validate:"required,uuid"`
	Period     string  `json:"period" binding:"required" validate:"required"`
	Allowance  float64 `json:"allowance"`
	Deduction  float64 `json:"deduction"`
	Notes      string  `json:"notes"`
}

type CreateDocumentRequest struct {
	Title    string `json:"title" binding:"required" validate:"required,max=255"`
	DocType  string `json:"doc_type" binding:"required" validate:"required,max=50"`
}

type CreateLetterRequest struct {
	Title      string `json:"title" binding:"required" validate:"required,max=255"`
	Content    string `json:"content" binding:"required"`
	LetterType string `json:"letter_type" binding:"required" validate:"required,oneof=official decision circular memo certificate"`
	To         string `json:"to" binding:"required" validate:"required,max=255"`
}

type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required" validate:"required,max=255"`
	Description string    `json:"description"`
	EventType   string    `json:"event_type" binding:"required" validate:"required,oneof=academic religious extracurricular national"`
	Date        time.Time `json:"date" binding:"required"`
	StartTime   string    `json:"start_time" validate:"max=10"`
	EndTime     string    `json:"end_time" validate:"max=10"`
	Location    string    `json:"location" validate:"max=255"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required" validate:"required,max=255"`
	Description string     `json:"description"`
	AssignedTo  string     `json:"assigned_to" binding:"required" validate:"required,uuid"`
	Priority    string     `json:"priority" validate:"oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Status  string     `json:"status" validate:"oneof=pending in_progress completed cancelled"`
	Priority string    `json:"priority" validate:"oneof=low medium high urgent"`
	DueDate *time.Time `json:"due_date"`
}

type CreateCurriculumRequest struct {
	GradeID   string `json:"grade_id" binding:"required" validate:"required,uuid"`
	SubjectID string `json:"subject_id" binding:"required" validate:"required,uuid"`
	SemesterID string `json:"semester_id" binding:"required" validate:"required,uuid"`
	Content   string `json:"content" binding:"required"`
}

type SettingsRequest struct {
	Settings map[string]string `json:"settings" binding:"required"`
}

type CreateNotificationRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1,dive,uuid"`
	Title   string   `json:"title" binding:"required" validate:"required,max=255"`
	Message string   `json:"message" binding:"required"`
	Type    string   `json:"type" validate:"oneof=info warning success error"`
	RefType string   `json:"ref_type" validate:"max=50"`
	RefID   string   `json:"ref_id" validate:"max=50"`
}

type BatchAttendanceRequest struct {
	Date        string                    `json:"date" binding:"required"`
	ScheduleID  string                    `json:"schedule_id" binding:"required" validate:"required,uuid"`
	Attendances []SingleAttendanceRequest `json:"attendances" binding:"required,min=1,dive"`
}

type SingleAttendanceRequest struct {
	StudentID string `json:"student_id" binding:"required" validate:"required,uuid"`
	Status    string `json:"status" binding:"required" validate:"required,oneof=present absent sick permit late"`
	Notes     string `json:"notes"`
}

type CreateEmployeeAttendanceRequest struct {
	Date   time.Time `json:"date" binding:"required"`
	Status string    `json:"status" binding:"required" validate:"required,oneof=present absent sick permit late"`
	Notes  string    `json:"notes"`
}

type DateRangeFilter struct {
	StartDate time.Time `json:"start_date" form:"start_date"`
	EndDate   time.Time `json:"end_date" form:"end_date"`
}

type GradeStudentFilter struct {
	ClassID    string `json:"class_id" form:"class_id" validate:"uuid"`
	SubjectID  string `json:"subject_id" form:"subject_id" validate:"uuid"`
	SemesterID string `json:"semester_id" form:"semester_id" validate:"uuid"`
}
