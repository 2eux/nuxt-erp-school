package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type AcademicHandler struct {
	academicService service.AcademicService
}

func NewAcademicHandler(academicService service.AcademicService) *AcademicHandler {
	return &AcademicHandler{academicService: academicService}
}

// ListSchedules godoc
// @Summary List schedules
// @Tags Academic
// @Produce json
// @Param semester_id query string true "Semester ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schedules [get]
func (h *AcademicHandler) ListSchedules(c *gin.Context) {
	semesterID := c.Query("semester_id")
	items, err := h.academicService.ListSchedules(c.Request.Context(), semesterID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateSchedule godoc
// @Summary Create schedule
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.CreateScheduleRequest true "Schedule data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/schedules [post]
func (h *AcademicHandler) CreateSchedule(c *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "schedule created", resp))
}

// GetSchedule godoc
// @Summary Get schedule
// @Tags Academic
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schedules/{id} [get]
func (h *AcademicHandler) GetSchedule(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.academicService.GetSchedule(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// ListAttendances godoc
// @Summary List attendances
// @Tags Academic
// @Produce json
// @Param student_id query string false "Student ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/attendances [get]
func (h *AcademicHandler) ListAttendances(c *gin.Context) {
	studentID := c.Query("student_id")
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()
	items, err := h.academicService.ListAttendances(c.Request.Context(), studentID, startDate, endDate)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateAttendance godoc
// @Summary Create attendance
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.CreateAttendanceRequest true "Attendance data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/attendances [post]
func (h *AcademicHandler) CreateAttendance(c *gin.Context) {
	createdBy := middleware.GetUserID(c)
	var req dto.CreateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.CreateAttendance(c.Request.Context(), createdBy, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "attendance created", resp))
}

// BulkCreateAttendance godoc
// @Summary Bulk create attendance
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.BatchAttendanceRequest true "Batch attendance data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/attendances/bulk [post]
func (h *AcademicHandler) BulkCreateAttendance(c *gin.Context) {
	createdBy := middleware.GetUserID(c)
	var req dto.BatchAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	if err := h.academicService.BulkCreateAttendance(c.Request.Context(), createdBy, req.ScheduleID, req); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "attendances created", nil))
}

// ListExams godoc
// @Summary List exams
// @Tags Academic
// @Produce json
// @Param class_id query string true "Class ID"
// @Param semester_id query string true "Semester ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/exams [get]
func (h *AcademicHandler) ListExams(c *gin.Context) {
	classID := c.Query("class_id")
	semesterID := c.Query("semester_id")
	items, err := h.academicService.ListExams(c.Request.Context(), classID, semesterID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateExam godoc
// @Summary Create exam
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.CreateExamRequest true "Exam data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/exams [post]
func (h *AcademicHandler) CreateExam(c *gin.Context) {
	createdBy := middleware.GetUserID(c)
	var req dto.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.CreateExam(c.Request.Context(), createdBy, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "exam created", resp))
}

// GetExamResults godoc
// @Summary Get exam results
// @Tags Academic
// @Produce json
// @Param id path string true "Exam ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/exams/{id}/results [get]
func (h *AcademicHandler) GetExamResults(c *gin.Context) {
	examID := c.Param("id")
	items, err := h.academicService.GetExamResults(c.Request.Context(), examID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateExamResult godoc
// @Summary Create exam result
// @Tags Academic
// @Accept json
// @Produce json
// @Param id path string true "Exam ID"
// @Param request body dto.CreateExamResultRequest true "Exam result data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/exams/{id}/results [post]
func (h *AcademicHandler) CreateExamResult(c *gin.Context) {
	var req dto.CreateExamResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	req.ExamID = c.Param("id")
	resp, err := h.academicService.CreateExamResult(c.Request.Context(), req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "exam result created", resp))
}

// ListGradebooks godoc
// @Summary List gradebooks
// @Tags Academic
// @Produce json
// @Param class_id query string true "Class ID"
// @Param subject_id query string true "Subject ID"
// @Param semester_id query string true "Semester ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/gradebooks [get]
func (h *AcademicHandler) ListGradebooks(c *gin.Context) {
	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	semesterID := c.Query("semester_id")
	items, err := h.academicService.ListGradebooks(c.Request.Context(), classID, subjectID, semesterID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateGradebook godoc
// @Summary Create gradebook
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.CreateGradebookRequest true "Gradebook data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/gradebooks [post]
func (h *AcademicHandler) CreateGradebook(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	var req dto.CreateGradebookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.CreateGradebook(c.Request.Context(), teacherID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "gradebook created", resp))
}

// ListReportCards godoc
// @Summary List report cards
// @Tags Academic
// @Produce json
// @Param class_id query string true "Class ID"
// @Param semester_id query string true "Semester ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/report-cards [get]
func (h *AcademicHandler) ListReportCards(c *gin.Context) {
	classID := c.Query("class_id")
	semesterID := c.Query("semester_id")
	items, err := h.academicService.ListReportCards(c.Request.Context(), classID, semesterID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// GenerateReportCard godoc
// @Summary Generate report card
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.CreateReportCardRequest true "Report card data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/report-cards/generate [post]
func (h *AcademicHandler) GenerateReportCard(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.CreateReportCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.GenerateReportCard(c.Request.Context(), userID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "report card generated", resp))
}

// ListAssignments godoc
// @Summary List assignments
// @Tags Academic
// @Produce json
// @Param class_id query string true "Class ID"
// @Param subject_id query string true "Subject ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/assignments [get]
func (h *AcademicHandler) ListAssignments(c *gin.Context) {
	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	items, err := h.academicService.ListAssignments(c.Request.Context(), classID, subjectID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateAssignment godoc
// @Summary Create assignment
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.CreateAssignmentRequest true "Assignment data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/assignments [post]
func (h *AcademicHandler) CreateAssignment(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	var req dto.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.CreateAssignment(c.Request.Context(), teacherID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "assignment created", resp))
}

// ListSubmissions godoc
// @Summary List submissions
// @Tags Academic
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/assignments/{id}/submissions [get]
func (h *AcademicHandler) ListSubmissions(c *gin.Context) {
	assignmentID := c.Param("id")
	items, err := h.academicService.ListSubmissions(c.Request.Context(), assignmentID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// SubmitAssignment godoc
// @Summary Submit assignment
// @Tags Academic
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Param request body dto.SubmitAssignmentRequest true "Submission data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/assignments/{id}/submissions [post]
func (h *AcademicHandler) SubmitAssignment(c *gin.Context) {
	studentID := middleware.GetUserID(c)
	assignmentID := c.Param("id")
	var req dto.SubmitAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.SubmitAssignment(c.Request.Context(), studentID, assignmentID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "assignment submitted", resp))
}

// ListLessonPlans godoc
// @Summary List lesson plans
// @Tags Academic
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/lesson-plans [get]
func (h *AcademicHandler) ListLessonPlans(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	filter := getPagination(c)
	items, total, err := h.academicService.ListLessonPlans(c.Request.Context(), teacherID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateLessonPlan godoc
// @Summary Create lesson plan
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.CreateLessonPlanRequest true "Lesson plan data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/lesson-plans [post]
func (h *AcademicHandler) CreateLessonPlan(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	var req dto.CreateLessonPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.CreateLessonPlan(c.Request.Context(), teacherID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "lesson plan created", resp))
}

// ListTeachingJournals godoc
// @Summary List teaching journals
// @Tags Academic
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/teaching-journals [get]
func (h *AcademicHandler) ListTeachingJournals(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	items, err := h.academicService.ListTeachingJournals(c.Request.Context(), teacherID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateTeachingJournal godoc
// @Summary Create teaching journal
// @Tags Academic
// @Accept json
// @Produce json
// @Param request body dto.CreateTeachingJournalRequest true "Journal data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/teaching-journals [post]
func (h *AcademicHandler) CreateTeachingJournal(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	var req dto.CreateTeachingJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.academicService.CreateTeachingJournal(c.Request.Context(), teacherID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "journal created", resp))
}
