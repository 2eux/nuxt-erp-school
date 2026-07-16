package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"go.uber.org/zap"
)

type LibraryService interface {
	ListBooks(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.LibraryBookResponse, int64, error)
	GetBook(ctx context.Context, id string) (*dto.LibraryBookResponse, error)
	CreateBook(ctx context.Context, schoolID string, req dto.CreateLibraryBookRequest) (*dto.LibraryBookResponse, error)
	BorrowBook(ctx context.Context, bookID string, req dto.BorrowBookRequest) (*dto.BookBorrowingResponse, error)
	ReturnBook(ctx context.Context, borrowingID string, req dto.ReturnBookRequest) (*dto.BookBorrowingResponse, error)
	ListBorrowings(ctx context.Context, bookID string) ([]dto.BookBorrowingResponse, error)
}

type libraryService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewLibraryService(db *sqlx.DB, logger *zap.Logger) LibraryService {
	return &libraryService{db: db, logger: logger}
}

func (s *libraryService) ListBooks(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.LibraryBookResponse, int64, error) {
	filter.Defaults()

	var items []domain.LibraryBook
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM library_books WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM library_books WHERE school_id=$1 ORDER BY title LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list books", err)
	}

	result := make([]dto.LibraryBookResponse, len(items))
	for i, b := range items {
		result[i] = dto.LibraryBookResponse{
			ID:          b.ID,
			SchoolID:    b.SchoolID,
			ISBN:        b.ISBN,
			Title:       b.Title,
			Author:      b.Author,
			Publisher:   b.Publisher,
			PublishYear: b.PublishYear,
			Category:    b.Category,
			Language:    b.Language,
			TotalCopies: b.TotalCopies,
			Available:   b.Available,
			Location:    b.Location,
			CoverURL:    b.CoverURL,
		}
	}
	return result, total, nil
}

func (s *libraryService) GetBook(ctx context.Context, id string) (*dto.LibraryBookResponse, error) {
	var b domain.LibraryBook
	if err := s.db.GetContext(ctx, &b, `SELECT * FROM library_books WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("book", id)
	}
	return &dto.LibraryBookResponse{
		ID:          b.ID,
		SchoolID:    b.SchoolID,
		ISBN:        b.ISBN,
		Title:       b.Title,
		Author:      b.Author,
		Publisher:   b.Publisher,
		PublishYear: b.PublishYear,
		Category:    b.Category,
		Language:    b.Language,
		TotalCopies: b.TotalCopies,
		Available:   b.Available,
		Location:    b.Location,
		CoverURL:    b.CoverURL,
	}, nil
}

func (s *libraryService) CreateBook(ctx context.Context, schoolID string, req dto.CreateLibraryBookRequest) (*dto.LibraryBookResponse, error) {
	b := &domain.LibraryBook{
		ID:          uuid.New().String(),
		SchoolID:    schoolID,
		ISBN:        req.ISBN,
		Title:       req.Title,
		Author:      req.Author,
		Publisher:   req.Publisher,
		PublishYear: req.PublishYear,
		Category:    req.Category,
		Language:    req.Language,
		TotalCopies: req.TotalCopies,
		Available:   req.TotalCopies,
		Location:    req.Location,
	}

	query := `INSERT INTO library_books (id, school_id, isbn, title, author, publisher, publish_year, category, language, total_copies, available, location) VALUES (:id, :school_id, :isbn, :title, :author, :publisher, :publish_year, :category, :language, :total_copies, :available, :location)`
	if _, err := s.db.NamedExecContext(ctx, query, b); err != nil {
		return nil, domain.NewInternalError("failed to create book", err)
	}

	return s.GetBook(ctx, b.ID)
}

func (s *libraryService) BorrowBook(ctx context.Context, bookID string, req dto.BorrowBookRequest) (*dto.BookBorrowingResponse, error) {
	var book domain.LibraryBook
	if err := s.db.GetContext(ctx, &book, `SELECT * FROM library_books WHERE id=$1`, bookID); err != nil {
		return nil, domain.NewNotFoundError("book", bookID)
	}

	if book.Available <= 0 {
		return nil, domain.NewInvalidInputError("book not available")
	}

	borrowing := &domain.BookBorrowing{
		ID:           uuid.New().String(),
		BookID:       bookID,
		BorrowerID:   req.BorrowerID,
		BorrowerType: req.BorrowerType,
		BorrowDate:   time.Now(),
		DueDate:      req.DueDate,
		Status:       "borrowed",
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO book_borrowings (id, book_id, borrower_id, borrower_type, borrow_date, due_date, status) VALUES (:id, :book_id, :borrower_id, :borrower_type, :borrow_date, :due_date, :status)`
	if _, err := tx.NamedExecContext(ctx, query, borrowing); err != nil {
		return nil, domain.NewInternalError("failed to borrow book", err)
	}

	tx.ExecContext(ctx, `UPDATE library_books SET available=available-1, updated_at=NOW() WHERE id=$1`, bookID)

	if err := tx.Commit(); err != nil {
		return nil, domain.NewInternalError("failed to commit", err)
	}

	return &dto.BookBorrowingResponse{
		ID:           borrowing.ID,
		BookID:       borrowing.BookID,
		BookTitle:    book.Title,
		BorrowerID:   borrowing.BorrowerID,
		BorrowDate:   borrowing.BorrowDate,
		DueDate:      borrowing.DueDate,
		Status:       borrowing.Status,
	}, nil
}

func (s *libraryService) ReturnBook(ctx context.Context, borrowingID string, req dto.ReturnBookRequest) (*dto.BookBorrowingResponse, error) {
	var borrowing domain.BookBorrowing
	if err := s.db.GetContext(ctx, &borrowing, `SELECT * FROM book_borrowings WHERE id=$1`, borrowingID); err != nil {
		return nil, domain.NewNotFoundError("borrowing", borrowingID)
	}

	returnDate := req.ReturnDate
	if returnDate.IsZero() {
		returnDate = time.Now()
	}

	var fine float64
	if returnDate.After(borrowing.DueDate) {
		daysLate := int(returnDate.Sub(borrowing.DueDate).Hours() / 24)
		fine = float64(daysLate) * 5000
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	tx.ExecContext(ctx, `UPDATE book_borrowings SET return_date=$1, status='returned', fine=$2, updated_at=NOW() WHERE id=$3`, returnDate, fine, borrowingID)
	tx.ExecContext(ctx, `UPDATE library_books SET available=available+1, updated_at=NOW() WHERE id=$1`, borrowing.BookID)

	if err := tx.Commit(); err != nil {
		return nil, domain.NewInternalError("failed to commit", err)
	}

	borrowing.ReturnDate = &returnDate
	borrowing.Status = "returned"
	borrowing.Fine = fine

	return &dto.BookBorrowingResponse{
		ID:           borrowing.ID,
		BookID:       borrowing.BookID,
		BorrowerID:   borrowing.BorrowerID,
		BorrowDate:   borrowing.BorrowDate,
		DueDate:      borrowing.DueDate,
		ReturnDate:   borrowing.ReturnDate,
		Status:       borrowing.Status,
		Fine:         borrowing.Fine,
	}, nil
}

func (s *libraryService) ListBorrowings(ctx context.Context, bookID string) ([]dto.BookBorrowingResponse, error) {
	var items []domain.BookBorrowing
	query := `SELECT * FROM book_borrowings WHERE book_id=$1 ORDER BY borrow_date DESC`
	if err := s.db.SelectContext(ctx, &items, query, bookID); err != nil {
		return nil, domain.NewInternalError("failed to list borrowings", err)
	}

	result := make([]dto.BookBorrowingResponse, len(items))
	for i, b := range items {
		result[i] = dto.BookBorrowingResponse{
			ID:           b.ID,
			BookID:       b.BookID,
			BorrowerID:   b.BorrowerID,
			BorrowDate:   b.BorrowDate,
			DueDate:      b.DueDate,
			ReturnDate:   b.ReturnDate,
			Status:       b.Status,
			Fine:         b.Fine,
		}
	}
	return result, nil
}
