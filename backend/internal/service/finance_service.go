package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
	"go.uber.org/zap"
)

type FinanceService interface {
	ListFeeTypes(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.FeeTypeResponse, int64, error)
	CreateFeeType(ctx context.Context, schoolID string, req dto.CreateFeeTypeRequest) (*dto.FeeTypeResponse, error)

	ListInvoices(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.InvoiceResponse, int64, error)
	GetInvoice(ctx context.Context, id string) (*dto.InvoiceResponse, error)
	CreateInvoice(ctx context.Context, schoolID, createdBy string, req dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error)
	SendInvoice(ctx context.Context, id string) error

	ListPayments(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.PaymentResponse, int64, error)
	CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	VerifyPayment(ctx context.Context, id, verifiedBy string, req dto.VerifyPaymentRequest) error

	ListJournals(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.JournalResponse, int64, error)
	CreateJournal(ctx context.Context, schoolID, createdBy string, req dto.CreateJournalRequest) (*dto.JournalResponse, error)

	ListLedger(ctx context.Context, schoolID string, accountCode string) ([]dto.GeneralLedgerResponse, error)

	ListCashTransactions(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]map[string]interface{}, int64, error)
	CreateCashTransaction(ctx context.Context, schoolID, createdBy string, transType string, amount float64, notes string) error

	ListPayrollPeriods(ctx context.Context, schoolID string) ([]string, error)
	ListPayrollDetails(ctx context.Context, schoolID, period string) ([]dto.PayrollResponse, error)
	ProcessPayroll(ctx context.Context, schoolID string, req dto.CreatePayrollRequest) (*dto.PayrollResponse, error)
}

type financeService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewFinanceService(db *sqlx.DB, logger *zap.Logger) FinanceService {
	return &financeService{db: db, logger: logger}
}

func (s *financeService) ListFeeTypes(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.FeeTypeResponse, int64, error) {
	filter.Defaults()
	var items []domain.FeeType
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM fee_types WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM fee_types WHERE school_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list fee types", err)
	}

	result := make([]dto.FeeTypeResponse, len(items))
	for i, ft := range items {
		result[i] = dto.FeeTypeResponse{
			ID:          ft.ID,
			SchoolID:    ft.SchoolID,
			Name:        ft.Name,
			Amount:      ft.Amount,
			Category:    ft.Category,
			Frequency:   ft.Frequency,
			Description: ft.Description,
			IsActive:    ft.IsActive,
			CreatedAt:   ft.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *financeService) CreateFeeType(ctx context.Context, schoolID string, req dto.CreateFeeTypeRequest) (*dto.FeeTypeResponse, error) {
	ft := &domain.FeeType{
		ID:          uuid.New().String(),
		SchoolID:    schoolID,
		Name:        req.Name,
		Amount:      req.Amount,
		Category:    req.Category,
		Frequency:   req.Frequency,
		Description: req.Description,
		IsActive:    true,
	}

	query := `INSERT INTO fee_types (id, school_id, name, amount, category, frequency, description, is_active) VALUES (:id, :school_id, :name, :amount, :category, :frequency, :description, :is_active)`
	if _, err := s.db.NamedExecContext(ctx, query, ft); err != nil {
		return nil, domain.NewInternalError("failed to create fee type", err)
	}

	return &dto.FeeTypeResponse{
		ID:          ft.ID,
		SchoolID:    ft.SchoolID,
		Name:        ft.Name,
		Amount:      ft.Amount,
		Category:    ft.Category,
		Frequency:   ft.Frequency,
		Description: ft.Description,
		IsActive:    ft.IsActive,
		CreatedAt:   ft.CreatedAt,
	}, nil
}

func (s *financeService) ListInvoices(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.InvoiceResponse, int64, error) {
	filter.Defaults()

	type invoiceRow struct {
		domain.Invoice
		StudentName string `db:"student_name"`
	}

	var rows []invoiceRow
	if err := s.db.SelectContext(ctx, &rows, database.ListInvoices, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list invoices", err)
	}

	var total int64
	s.db.GetContext(ctx, &total, database.CountInvoices, schoolID)

	result := make([]dto.InvoiceResponse, len(rows))
	for i, inv := range rows {
		var items []domain.InvoiceItem
		s.db.SelectContext(ctx, &items, database.GetInvoiceItems, inv.ID)

		itemResponses := make([]dto.InvoiceItemResponse, len(items))
		for j, item := range items {
			itemResponses[j] = dto.InvoiceItemResponse{
				ID:        item.ID,
				FeeTypeID: item.FeeTypeID,
				Name:      item.Name,
				Amount:    item.Amount,
			}
		}

		result[i] = dto.InvoiceResponse{
			ID:          inv.ID,
			SchoolID:    inv.SchoolID,
			StudentID:   inv.StudentID,
			StudentName: inv.StudentName,
			InvoiceNo:   inv.InvoiceNo,
			TotalAmount: inv.TotalAmount,
			PaidAmount:  inv.PaidAmount,
			Status:      inv.Status,
			DueDate:     inv.DueDate,
			PaidAt:      inv.PaidAt,
			SemesterID:  inv.SemesterID,
			Notes:       inv.Notes,
			Items:       itemResponses,
			CreatedAt:   inv.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *financeService) GetInvoice(ctx context.Context, id string) (*dto.InvoiceResponse, error) {
	type invoiceRow struct {
		domain.Invoice
		StudentName string `db:"student_name"`
	}
	var row invoiceRow
	if err := s.db.GetContext(ctx, &row, database.GetInvoiceWithItems, id); err != nil {
		return nil, domain.NewNotFoundError("invoice", id)
	}

	var items []domain.InvoiceItem
	s.db.SelectContext(ctx, &items, database.GetInvoiceItems, id)

	itemResponses := make([]dto.InvoiceItemResponse, len(items))
	for j, item := range items {
		itemResponses[j] = dto.InvoiceItemResponse{
			ID:        item.ID,
			FeeTypeID: item.FeeTypeID,
			Name:      item.Name,
			Amount:    item.Amount,
		}
	}

	return &dto.InvoiceResponse{
		ID:          row.ID,
		SchoolID:    row.SchoolID,
		StudentID:   row.StudentID,
		StudentName: row.StudentName,
		InvoiceNo:   row.InvoiceNo,
		TotalAmount: row.TotalAmount,
		PaidAmount:  row.PaidAmount,
		Status:      row.Status,
		DueDate:     row.DueDate,
		PaidAt:      row.PaidAt,
		SemesterID:  row.SemesterID,
		Notes:       row.Notes,
		Items:       itemResponses,
		CreatedAt:   row.CreatedAt,
	}, nil
}

func (s *financeService) CreateInvoice(ctx context.Context, schoolID, createdBy string, req dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error) {
	invoiceNo := fmt.Sprintf("INV-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%100000)

	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Amount
	}

	inv := &domain.Invoice{
		ID:          uuid.New().String(),
		SchoolID:    schoolID,
		StudentID:   req.StudentID,
		InvoiceNo:   invoiceNo,
		TotalAmount: totalAmount,
		DueDate:     req.DueDate,
		SemesterID:  req.SemesterID,
		Notes:       req.Notes,
		CreatedBy:   createdBy,
		Status:      "unpaid",
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	invQuery := `INSERT INTO invoices (id, school_id, student_id, invoice_no, total_amount, due_date, semester_id, notes, created_by, status) VALUES (:id, :school_id, :student_id, :invoice_no, :total_amount, :due_date, :semester_id, :notes, :created_by, :status)`
	if _, err := tx.NamedExecContext(ctx, invQuery, inv); err != nil {
		return nil, domain.NewInternalError("failed to create invoice", err)
	}

	for _, item := range req.Items {
		ii := &domain.InvoiceItem{
			ID:        uuid.New().String(),
			InvoiceID: inv.ID,
			FeeTypeID: item.FeeTypeID,
			Name:      item.Name,
			Amount:    item.Amount,
		}
		itemQuery := `INSERT INTO invoice_items (id, invoice_id, fee_type_id, name, amount) VALUES (:id, :invoice_id, :fee_type_id, :name, :amount)`
		if _, err := tx.NamedExecContext(ctx, itemQuery, ii); err != nil {
			return nil, domain.NewInternalError("failed to create invoice item", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, domain.NewInternalError("failed to commit invoice", err)
	}

	return s.GetInvoice(ctx, inv.ID)
}

func (s *financeService) SendInvoice(ctx context.Context, id string) error {
	return nil
}

func (s *financeService) ListPayments(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.PaymentResponse, int64, error) {
	filter.Defaults()

	type paymentRow struct {
		domain.Payment
		InvoiceNo string `db:"invoice_no"`
	}

	var rows []paymentRow
	if err := s.db.SelectContext(ctx, &rows, database.ListPayments, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list payments", err)
	}

	var total int64
	s.db.GetContext(ctx, &total, database.CountPayments, schoolID)

	result := make([]dto.PaymentResponse, len(rows))
	for i, p := range rows {
		result[i] = dto.PaymentResponse{
			ID:         p.ID,
			InvoiceID:  p.InvoiceID,
			InvoiceNo:  p.InvoiceNo,
			PaymentNo:  p.PaymentNo,
			Amount:     p.Amount,
			Method:     p.Method,
			Status:     p.Status,
			ProofURL:   p.ProofURL,
			PaidAt:     p.PaidAt,
			VerifiedBy: p.VerifiedBy,
			VerifiedAt: p.VerifiedAt,
			Notes:      p.Notes,
			CreatedAt:  p.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *financeService) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	paymentNo := fmt.Sprintf("PAY-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%100000)

	paidAt := req.PaidAt
	if paidAt.IsZero() {
		paidAt = time.Now()
	}

	p := &domain.Payment{
		ID:        uuid.New().String(),
		InvoiceID: req.InvoiceID,
		PaymentNo: paymentNo,
		Amount:    req.Amount,
		Method:    req.Method,
		PaidAt:    paidAt,
		Status:    "pending",
		Notes:     req.Notes,
	}

	query := `INSERT INTO payments (id, invoice_id, payment_no, amount, method, paid_at, status, notes) VALUES (:id, :invoice_id, :payment_no, :amount, :method, :paid_at, :status, :notes)`
	if _, err := s.db.NamedExecContext(ctx, query, p); err != nil {
		return nil, domain.NewInternalError("failed to create payment", err)
	}

	return &dto.PaymentResponse{
		ID:         p.ID,
		InvoiceID:  p.InvoiceID,
		PaymentNo:  p.PaymentNo,
		Amount:     p.Amount,
		Method:     p.Method,
		Status:     p.Status,
		PaidAt:     p.PaidAt,
		Notes:      p.Notes,
		CreatedAt:  p.CreatedAt,
	}, nil
}

func (s *financeService) VerifyPayment(ctx context.Context, id, verifiedBy string, req dto.VerifyPaymentRequest) error {
	now := time.Now()
	query := `UPDATE payments SET status=$1, verified_by=$2, verified_at=$3, notes=$4, updated_at=NOW() WHERE id=$5`
	if _, err := s.db.ExecContext(ctx, query, req.Status, verifiedBy, now, req.Notes, id); err != nil {
		return domain.NewInternalError("failed to verify payment", err)
	}

	if req.Status == "verified" {
		s.db.ExecContext(ctx, `UPDATE invoices SET status='paid', paid_amount=total_amount, paid_at=$1, updated_at=NOW() WHERE id=(SELECT invoice_id FROM payments WHERE id=$2)`, now, id)
	}
	return nil
}

func (s *financeService) ListJournals(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.JournalResponse, int64, error) {
	filter.Defaults()

	var items []domain.Journal
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM journals WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM journals WHERE school_id=$1 ORDER BY date DESC LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list journals", err)
	}

	result := make([]dto.JournalResponse, len(items))
	for i, j := range items {
		var entries []domain.JournalEntry
		s.db.SelectContext(ctx, &entries, `SELECT * FROM journal_entries WHERE journal_id=$1`, j.ID)

		var totalDebit, totalCredit float64
		entryResponses := make([]dto.JournalEntryResponse, len(entries))
		for k, e := range entries {
			entryResponses[k] = dto.JournalEntryResponse{
				ID:          e.ID,
				AccountCode: e.AccountCode,
				Description: e.Description,
				Debit:       e.Debit,
				Credit:      e.Credit,
			}
			totalDebit += e.Debit
			totalCredit += e.Credit
		}

		result[i] = dto.JournalResponse{
			ID:          j.ID,
			SchoolID:    j.SchoolID,
			JournalNo:   j.JournalNo,
			Date:        j.Date,
			Description: j.Description,
			Status:      j.Status,
			TotalDebit:  totalDebit,
			TotalCredit: totalCredit,
			Entries:     entryResponses,
			CreatedBy:   j.CreatedBy,
			CreatedAt:   j.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *financeService) CreateJournal(ctx context.Context, schoolID, createdBy string, req dto.CreateJournalRequest) (*dto.JournalResponse, error) {
	journalNo := fmt.Sprintf("JRN-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%100000)

	j := &domain.Journal{
		ID:          uuid.New().String(),
		SchoolID:    schoolID,
		JournalNo:   journalNo,
		Date:        req.Date,
		Description: req.Description,
		Status:      "draft",
		CreatedBy:   createdBy,
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	jQuery := `INSERT INTO journals (id, school_id, journal_no, date, description, status, created_by) VALUES (:id, :school_id, :journal_no, :date, :description, :status, :created_by)`
	if _, err := tx.NamedExecContext(ctx, jQuery, j); err != nil {
		return nil, domain.NewInternalError("failed to create journal", err)
	}

	for _, entry := range req.Entries {
		e := &domain.JournalEntry{
			ID:          uuid.New().String(),
			JournalID:   j.ID,
			AccountCode: entry.AccountCode,
			Description: entry.Description,
			Debit:       entry.Debit,
			Credit:      entry.Credit,
		}
		eQuery := `INSERT INTO journal_entries (id, journal_id, account_code, description, debit, credit) VALUES (:id, :journal_id, :account_code, :description, :debit, :credit)`
		if _, err := tx.NamedExecContext(ctx, eQuery, e); err != nil {
			return nil, domain.NewInternalError("failed to create journal entry", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, domain.NewInternalError("failed to commit journal", err)
	}

	result, _, _ := s.ListJournals(ctx, schoolID, dto.PaginationRequest{Page: 1, PageSize: 1})
	if len(result) > 0 {
		return &result[0], nil
	}

	return &dto.JournalResponse{
		ID:          j.ID,
		SchoolID:    j.SchoolID,
		JournalNo:   j.JournalNo,
		Date:        j.Date,
		Description: j.Description,
		Status:      j.Status,
		CreatedBy:   j.CreatedBy,
		CreatedAt:   j.CreatedAt,
	}, nil
}

func (s *financeService) ListLedger(ctx context.Context, schoolID string, accountCode string) ([]dto.GeneralLedgerResponse, error) {
	var items []domain.GeneralLedger
	query := `SELECT * FROM general_ledgers WHERE school_id=$1`
	args := []interface{}{schoolID}
	if accountCode != "" {
		query += ` AND account_code=$2`
		args = append(args, accountCode)
	}
	query += ` ORDER BY date DESC LIMIT 100`

	if err := s.db.SelectContext(ctx, &items, query, args...); err != nil {
		return nil, domain.NewInternalError("failed to list ledger", err)
	}

	result := make([]dto.GeneralLedgerResponse, len(items))
	for i, l := range items {
		result[i] = dto.GeneralLedgerResponse{
			ID:          l.ID,
			AccountCode: l.AccountCode,
			AccountName: l.AccountName,
			Date:        l.Date,
			Description: l.Description,
			Debit:       l.Debit,
			Credit:      l.Credit,
			Balance:     l.Balance,
		}
	}
	return result, nil
}

func (s *financeService) ListCashTransactions(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}

func (s *financeService) CreateCashTransaction(ctx context.Context, schoolID, createdBy string, transType string, amount float64, notes string) error {
	return nil
}

func (s *financeService) ListPayrollPeriods(ctx context.Context, schoolID string) ([]string, error) {
	var periods []string
	query := `SELECT DISTINCT period FROM payroll_details WHERE school_id=$1 ORDER BY period DESC`
	if err := s.db.SelectContext(ctx, &periods, query, schoolID); err != nil {
		return nil, domain.NewInternalError("failed to list periods", err)
	}
	return periods, nil
}

func (s *financeService) ListPayrollDetails(ctx context.Context, schoolID, period string) ([]dto.PayrollResponse, error) {
	type payrollRow struct {
		domain.PayrollDetail
		FullName string `db:"full_name"`
	}

	query := `SELECT pd.*, u.full_name FROM payroll_details pd JOIN employees e ON pd.employee_id = e.id JOIN users u ON e.user_id = u.id WHERE pd.school_id=$1 AND pd.period=$2`

	var rows []payrollRow
	if err := s.db.SelectContext(ctx, &rows, query, schoolID, period); err != nil {
		return nil, domain.NewInternalError("failed to list payroll details", err)
	}

	result := make([]dto.PayrollResponse, len(rows))
	for i, r := range rows {
		result[i] = dto.PayrollResponse{
			ID:         r.ID,
			EmployeeID: r.EmployeeID,
			FullName:   r.FullName,
			Period:     r.Period,
			BaseSalary: r.BaseSalary,
			Allowance:  r.Allowance,
			Deduction:  r.Deduction,
			NetSalary:  r.NetSalary,
			Status:     r.Status,
			PaidAt:     r.PaidAt,
			CreatedAt:  r.CreatedAt,
		}
	}
	return result, nil
}

func (s *financeService) ProcessPayroll(ctx context.Context, schoolID string, req dto.CreatePayrollRequest) (*dto.PayrollResponse, error) {
	var emp domain.Employee
	if err := s.db.GetContext(ctx, &emp, `SELECT * FROM employees WHERE id=$1`, req.EmployeeID); err != nil {
		return nil, domain.NewNotFoundError("employee", req.EmployeeID)
	}

	netSalary := emp.BaseSalary + req.Allowance - req.Deduction
	pd := &domain.PayrollDetail{
		ID:         uuid.New().String(),
		EmployeeID: req.EmployeeID,
		SchoolID:   schoolID,
		Period:     req.Period,
		BaseSalary: emp.BaseSalary,
		Allowance:  req.Allowance,
		Deduction:  req.Deduction,
		NetSalary:  netSalary,
		Status:     "pending",
		Notes:      req.Notes,
	}

	query := `INSERT INTO payroll_details (id, employee_id, school_id, period, base_salary, allowance, deduction, net_salary, status, notes) VALUES (:id, :employee_id, :school_id, :period, :base_salary, :allowance, :deduction, :net_salary, :status, :notes)`
	if _, err := s.db.NamedExecContext(ctx, query, pd); err != nil {
		return nil, domain.NewInternalError("failed to process payroll", err)
	}

	return &dto.PayrollResponse{
		ID:         pd.ID,
		EmployeeID: pd.EmployeeID,
		Period:     pd.Period,
		BaseSalary: pd.BaseSalary,
		Allowance:  pd.Allowance,
		Deduction:  pd.Deduction,
		NetSalary:  pd.NetSalary,
		Status:     pd.Status,
	}, nil
}
