package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/middleware"
)

type Router struct {
	authHandler         *AuthHandler
	schoolHandler       *SchoolHandler
	userHandler         *UserHandler
	studentHandler      *StudentHandler
	teacherHandler      *TeacherHandler
	employeeHandler     *EmployeeHandler
	academicHandler     *AcademicHandler
	islamicHandler      *IslamicHandler
	financeHandler      *FinanceHandler
	notificationHandler *NotificationHandler
	documentHandler     *DocumentHandler
	admissionHandler    *AdmissionHandler
	analyticsHandler    *AnalyticsHandler
	aiHandler           *AIHandler
	authMiddleware      *middleware.AuthMiddleware
}

func NewRouter(
	authHandler *AuthHandler,
	schoolHandler *SchoolHandler,
	userHandler *UserHandler,
	studentHandler *StudentHandler,
	teacherHandler *TeacherHandler,
	employeeHandler *EmployeeHandler,
	academicHandler *AcademicHandler,
	islamicHandler *IslamicHandler,
	financeHandler *FinanceHandler,
	notificationHandler *NotificationHandler,
	documentHandler *DocumentHandler,
	admissionHandler *AdmissionHandler,
	analyticsHandler *AnalyticsHandler,
	aiHandler *AIHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		authHandler:         authHandler,
		schoolHandler:       schoolHandler,
		userHandler:         userHandler,
		studentHandler:      studentHandler,
		teacherHandler:      teacherHandler,
		employeeHandler:     employeeHandler,
		academicHandler:     academicHandler,
		islamicHandler:      islamicHandler,
		financeHandler:      financeHandler,
		notificationHandler: notificationHandler,
		documentHandler:     documentHandler,
		admissionHandler:    admissionHandler,
		analyticsHandler:    analyticsHandler,
		aiHandler:           aiHandler,
		authMiddleware:      authMiddleware,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	api := engine.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		auth.POST("/logout", r.authHandler.Logout)
		auth.POST("/forgot-password", r.authHandler.ForgotPassword)
		auth.POST("/reset-password", r.authHandler.ResetPassword)
		auth.GET("/me", r.authMiddleware.Authenticate(), r.authHandler.GetMe)
		auth.PUT("/profile", r.authMiddleware.Authenticate(), r.authHandler.UpdateProfile)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(r.authMiddleware.Authenticate())
	{
		// Schools
		schools := protected.Group("/schools")
		{
			schools.GET("", r.schoolHandler.ListSchools)
			schools.POST("", r.schoolHandler.CreateSchool)
			schools.GET("/:id", r.schoolHandler.GetSchool)
			schools.PUT("/:id", r.schoolHandler.UpdateSchool)
			schools.DELETE("/:id", r.schoolHandler.DeleteSchool)
			schools.GET("/:id/academic-years", r.schoolHandler.ListAcademicYears)
			schools.POST("/:id/academic-years", r.schoolHandler.CreateAcademicYear)
			schools.PUT("/:id/academic-years/:ayid", r.schoolHandler.UpdateAcademicYear)
			schools.GET("/:id/grades", r.schoolHandler.ListGrades)
			schools.POST("/:id/grades", r.schoolHandler.CreateGrade)
			schools.GET("/:id/classes", r.schoolHandler.ListClasses)
			schools.POST("/:id/classes", r.schoolHandler.CreateClass)
			schools.GET("/:id/subjects", r.schoolHandler.ListSubjects)
			schools.POST("/:id/subjects", r.schoolHandler.CreateSubject)
		}

		// Academic Years -> Semesters
		academicYears := protected.Group("/academic-years")
		{
			academicYears.GET("/:id/semesters", r.schoolHandler.ListSemesters)
			academicYears.POST("/:id/semesters", r.schoolHandler.CreateSemester)
		}

		// Users & Roles
		users := protected.Group("/users")
		{
			users.GET("", r.userHandler.ListUsers)
			users.POST("", r.userHandler.CreateUser)
			users.GET("/:id", r.userHandler.GetUser)
			users.PUT("/:id", r.userHandler.UpdateUser)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}

		roles := protected.Group("/roles")
		{
			roles.GET("", r.userHandler.ListRoles)
			roles.POST("", r.userHandler.CreateRole)
			roles.PUT("/:id", r.userHandler.UpdateRole)
		}

		protected.GET("/permissions", r.userHandler.ListPermissions)

		// Students
		students := protected.Group("/students")
		{
			students.GET("", r.studentHandler.ListStudents)
			students.POST("", r.studentHandler.CreateStudent)
			students.GET("/:id", r.studentHandler.GetStudent)
			students.PUT("/:id", r.studentHandler.UpdateStudent)
			students.DELETE("/:id", r.studentHandler.DeleteStudent)
			students.POST("/:id/promote", r.studentHandler.PromoteStudent)
			students.POST("/:id/transfer", r.studentHandler.TransferStudent)
			students.GET("/:id/parents", r.studentHandler.ListParents)
			students.POST("/:id/parents", r.studentHandler.LinkParent)
			students.GET("/:id/documents", r.studentHandler.ListDocuments)
			students.POST("/:id/documents", r.studentHandler.UploadDocument)
		}

		// Teachers
		teachers := protected.Group("/teachers")
		{
			teachers.GET("", r.teacherHandler.ListTeachers)
			teachers.POST("", r.teacherHandler.CreateTeacher)
			teachers.GET("/:id", r.teacherHandler.GetTeacher)
			teachers.PUT("/:id", r.teacherHandler.UpdateTeacher)
			teachers.POST("/:id/subjects", r.teacherHandler.AssignSubjects)
			teachers.GET("/:id/schedule", r.teacherHandler.GetSchedule)
		}

		// Employees
		employees := protected.Group("/employees")
		{
			employees.GET("", r.employeeHandler.ListEmployees)
			employees.POST("", r.employeeHandler.CreateEmployee)
			employees.GET("/:id", r.employeeHandler.GetEmployee)
			employees.PUT("/:id", r.employeeHandler.UpdateEmployee)
			employees.GET("/:id/attendances", r.employeeHandler.ListAttendances)
			employees.POST("/:id/attendances", r.employeeHandler.CreateAttendance)
		}

		leaveRequests := protected.Group("/leave-requests")
		{
			leaveRequests.GET("", r.employeeHandler.ListLeaveRequests)
			leaveRequests.POST("", r.employeeHandler.SubmitLeaveRequest)
			leaveRequests.PUT("/:id/approve", r.employeeHandler.ApproveLeave)
			leaveRequests.PUT("/:id/reject", r.employeeHandler.RejectLeave)
		}

		// Academic
		protected.GET("/schedules", r.academicHandler.ListSchedules)
		protected.POST("/schedules", r.academicHandler.CreateSchedule)
		protected.GET("/schedules/:id", r.academicHandler.GetSchedule)

		attendances := protected.Group("/attendances")
		{
			attendances.GET("", r.academicHandler.ListAttendances)
			attendances.POST("", r.academicHandler.CreateAttendance)
			attendances.POST("/bulk", r.academicHandler.BulkCreateAttendance)
		}

		exams := protected.Group("/exams")
		{
			exams.GET("", r.academicHandler.ListExams)
			exams.POST("", r.academicHandler.CreateExam)
			exams.GET("/:id/results", r.academicHandler.GetExamResults)
			exams.POST("/:id/results", r.academicHandler.CreateExamResult)
		}

		protected.GET("/gradebooks", r.academicHandler.ListGradebooks)
		protected.POST("/gradebooks", r.academicHandler.CreateGradebook)

		protected.GET("/report-cards", r.academicHandler.ListReportCards)
		protected.POST("/report-cards/generate", r.academicHandler.GenerateReportCard)

		assignments := protected.Group("/assignments")
		{
			assignments.GET("", r.academicHandler.ListAssignments)
			assignments.POST("", r.academicHandler.CreateAssignment)
			assignments.GET("/:id/submissions", r.academicHandler.ListSubmissions)
			assignments.POST("/:id/submissions", r.academicHandler.SubmitAssignment)
		}

		protected.GET("/lesson-plans", r.academicHandler.ListLessonPlans)
		protected.POST("/lesson-plans", r.academicHandler.CreateLessonPlan)
		protected.GET("/teaching-journals", r.academicHandler.ListTeachingJournals)
		protected.POST("/teaching-journals", r.academicHandler.CreateTeachingJournal)

		// Islamic
		tahfidz := protected.Group("/tahfidz")
		{
			tahfidz.GET("/progress", r.islamicHandler.ListTahfidzProgress)
			tahfidz.POST("/progress", r.islamicHandler.CreateTahfidzProgress)
		}

		protected.GET("/mutabaah", r.islamicHandler.ListMutabaah)
		protected.POST("/mutabaah", r.islamicHandler.CreateMutabaah)
		protected.GET("/prayer-attendance", r.islamicHandler.ListPrayerAttendance)
		protected.POST("/prayer-attendance", r.islamicHandler.CreatePrayerAttendance)

		halaqah := protected.Group("/halaqah")
		{
			halaqah.GET("/groups", r.islamicHandler.ListHalaqahGroups)
			halaqah.POST("/groups", r.islamicHandler.CreateHalaqahGroup)
			halaqah.GET("/groups/:id", r.islamicHandler.GetHalaqahGroup)
			halaqah.POST("/groups/:id/members", r.islamicHandler.AddHalaqahMember)
		}

		// Finance
		fees := protected.Group("/fees")
		{
			fees.GET("/types", r.financeHandler.ListFeeTypes)
			fees.POST("/types", r.financeHandler.CreateFeeType)
		}

		invoices := protected.Group("/invoices")
		{
			invoices.GET("", r.financeHandler.ListInvoices)
			invoices.POST("", r.financeHandler.CreateInvoice)
			invoices.GET("/:id", r.financeHandler.GetInvoice)
		}

		payments := protected.Group("/payments")
		{
			payments.GET("", r.financeHandler.ListPayments)
			payments.POST("", r.financeHandler.CreatePayment)
			payments.PUT("/:id/verify", r.financeHandler.VerifyPayment)
		}

		protected.GET("/journals", r.financeHandler.ListJournals)
		protected.POST("/journals", r.financeHandler.CreateJournal)
		protected.GET("/ledger", r.financeHandler.ListLedger)

		payroll := protected.Group("/payroll")
		{
			payroll.GET("/periods", r.financeHandler.ListPayrollPeriods)
			payroll.GET("/details", r.financeHandler.ListPayrollDetails)
			payroll.POST("/process", r.financeHandler.ProcessPayroll)
		}

		// Analytics
		analytics := protected.Group("/analytics")
		{
			analytics.GET("/dashboard", r.analyticsHandler.GetDashboard)
			analytics.GET("/academic", r.analyticsHandler.GetAcademicAnalytics)
			analytics.GET("/finance", r.analyticsHandler.GetFinanceAnalytics)
			analytics.GET("/tahfidz", r.analyticsHandler.GetTahfidzAnalytics)
			analytics.GET("/admissions", r.analyticsHandler.GetAdmissionAnalytics)
			analytics.GET("/attendance", r.analyticsHandler.GetAttendanceAnalytics)
		}

		// Notifications
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", r.notificationHandler.ListNotifications)
			notifications.PUT("/:id/read", r.notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", r.notificationHandler.MarkAllAsRead)
		}

		// Documents & Letters
		documents := protected.Group("/documents")
		{
			documents.GET("", r.documentHandler.ListDocuments)
			documents.POST("", r.documentHandler.UploadDocument)
			documents.GET("/:id", r.documentHandler.GetDocument)
			documents.DELETE("/:id", r.documentHandler.DeleteDocument)
		}

		letters := protected.Group("/letters")
		{
			letters.GET("", r.documentHandler.ListLetters)
			letters.POST("", r.documentHandler.CreateLetter)
			letters.GET("/:id", r.documentHandler.GetLetter)
		}

		// Admissions
		admissions := protected.Group("/admissions")
		{
			admissions.GET("/applicants", r.admissionHandler.ListApplicants)
			admissions.POST("/applicants", r.admissionHandler.CreateApplicant)
			admissions.POST("/applicants/:id/accept", r.admissionHandler.AcceptApplicant)
			admissions.POST("/applicants/:id/enroll", r.admissionHandler.EnrollApplicant)
		}

		// AI
		ai := protected.Group("/ai")
		{
			ai.POST("/chat", r.aiHandler.Chat)
			ai.POST("/generate", r.aiHandler.Generate)
			ai.GET("/conversations", r.aiHandler.ListConversations)
			ai.GET("/conversations/:id/messages", r.aiHandler.GetMessages)
			ai.POST("/conversations/:id/messages", r.aiHandler.SendMessage)
			ai.POST("/knowledge/upload", r.aiHandler.UploadKnowledge)
		}
	}
}
