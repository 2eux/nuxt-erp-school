package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
	"go.uber.org/zap"
)

type IslamicService interface {
	CreateTahfidzProgress(ctx context.Context, teacherID string, req dto.CreateTahfidzProgressRequest) (*dto.TahfidzProgressResponse, error)
	ListTahfidzProgress(ctx context.Context, studentID string) ([]dto.TahfidzProgressResponse, error)
	GetTahfidzSummary(ctx context.Context, studentID string) (*dto.TahfidzSummary, error)

	CreateMutabaah(ctx context.Context, req dto.CreateMutabaahRequest) (*dto.MutabaahResponse, error)
	ListMutabaah(ctx context.Context, studentID string, startDate, endDate time.Time) ([]dto.MutabaahResponse, error)

	CreatePrayerAttendance(ctx context.Context, req dto.CreatePrayerAttendanceRequest) (*dto.PrayerAttendanceResponse, error)
	ListPrayerAttendance(ctx context.Context, studentID string, startDate, endDate time.Time) ([]dto.PrayerAttendanceResponse, error)

	CreateHalaqahGroup(ctx context.Context, schoolID string, req dto.CreateHalaqahGroupRequest) (*dto.HalaqahGroupResponse, error)
	ListHalaqahGroups(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.HalaqahGroupResponse, int64, error)
	GetHalaqahGroup(ctx context.Context, id string) (*dto.HalaqahGroupResponse, error)
	AddHalaqahMember(ctx context.Context, halaqahID string, req dto.AddHalaqahMemberRequest) error
}

type islamicService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewIslamicService(db *sqlx.DB, logger *zap.Logger) IslamicService {
	return &islamicService{db: db, logger: logger}
}

func (s *islamicService) CreateTahfidzProgress(ctx context.Context, teacherID string, req dto.CreateTahfidzProgressRequest) (*dto.TahfidzProgressResponse, error) {
	tp := &domain.TahfidzProgress{
		ID:        uuid.New().String(),
		StudentID: req.StudentID,
		TeacherID: teacherID,
		Surah:     req.Surah,
		StartAyah: req.StartAyah,
		EndAyah:   req.EndAyah,
		Juz:       req.Juz,
		Page:      req.Page,
		Status:    req.Status,
		Quality:   req.Quality,
		Notes:     req.Notes,
		Date:      req.Date,
	}

	query := `INSERT INTO tahfidz_progress (id, student_id, teacher_id, surah, start_ayah, end_ayah, juz, page, status, quality, notes, date) VALUES (:id, :student_id, :teacher_id, :surah, :start_ayah, :end_ayah, :juz, :page, :status, :quality, :notes, :date)`
	if _, err := s.db.NamedExecContext(ctx, query, tp); err != nil {
		return nil, domain.NewInternalError("failed to create tahfidz progress", err)
	}

	return &dto.TahfidzProgressResponse{
		ID:        tp.ID,
		StudentID: tp.StudentID,
		TeacherID: tp.TeacherID,
		Surah:     tp.Surah,
		StartAyah: tp.StartAyah,
		EndAyah:   tp.EndAyah,
		Juz:       tp.Juz,
		Page:      tp.Page,
		Status:    tp.Status,
		Quality:   tp.Quality,
		Notes:     tp.Notes,
		Date:      tp.Date,
	}, nil
}

func (s *islamicService) ListTahfidzProgress(ctx context.Context, studentID string) ([]dto.TahfidzProgressResponse, error) {
	var items []struct {
		domain.TahfidzProgress
		StudentName string `db:"student_name"`
		TeacherName string `db:"teacher_name"`
	}

	if err := s.db.SelectContext(ctx, &items, database.ListTahfidzProgress, studentID); err != nil {
		return nil, domain.NewInternalError("failed to list tahfidz progress", err)
	}

	result := make([]dto.TahfidzProgressResponse, len(items))
	for i, tp := range items {
		result[i] = dto.TahfidzProgressResponse{
			ID:          tp.ID,
			StudentID:   tp.StudentID,
			StudentName: tp.StudentName,
			TeacherID:   tp.TeacherID,
			TeacherName: tp.TeacherName,
			Surah:       tp.Surah,
			StartAyah:   tp.StartAyah,
			EndAyah:     tp.EndAyah,
			Juz:         tp.Juz,
			Page:        tp.Page,
			Status:      tp.Status,
			Quality:     tp.Quality,
			Notes:       tp.Notes,
			Date:        tp.Date,
		}
	}
	return result, nil
}

func (s *islamicService) GetTahfidzSummary(ctx context.Context, studentID string) (*dto.TahfidzSummary, error) {
	var summary dto.TahfidzSummary
	query := database.GetTahfidzSummary
	if err := s.db.GetContext(ctx, &summary, query, studentID); err != nil {
		return nil, domain.NewInternalError("failed to get tahfidz summary", err)
	}
	return &summary, nil
}

func (s *islamicService) CreateMutabaah(ctx context.Context, req dto.CreateMutabaahRequest) (*dto.MutabaahResponse, error) {
	mu := &domain.MutabaahYaumiyah{
		ID:        uuid.New().String(),
		StudentID: req.StudentID,
		Date:      req.Date,
	}

	if req.Fajr != nil {
		mu.Fajr = *req.Fajr
	}
	if req.Dhuhr != nil {
		mu.Dhuhr = *req.Dhuhr
	}
	if req.Asr != nil {
		mu.Asr = *req.Asr
	}
	if req.Maghrib != nil {
		mu.Maghrib = *req.Maghrib
	}
	if req.Isha != nil {
		mu.Isha = *req.Isha
	}
	if req.Tahajjud != nil {
		mu.Tahajjud = *req.Tahajjud
	}
	if req.Dhuha != nil {
		mu.Dhuha = *req.Dhuha
	}
	if req.Sunnah != nil {
		mu.Sunnah = *req.Sunnah
	}
	if req.QuranTilawah != nil {
		mu.QuranTilawah = *req.QuranTilawah
	}
	if req.QuranHifdz != nil {
		mu.QuranHifdz = *req.QuranHifdz
	}
	if req.DzikrPagi != nil {
		mu.DzikrPagi = *req.DzikrPagi
	}
	if req.DzikrPetang != nil {
		mu.DzikrPetang = *req.DzikrPetang
	}
	if req.Shadaqah != nil {
		mu.Shadaqah = *req.Shadaqah
	}
	if req.PuasaSunnah != nil {
		mu.PuasaSunnah = *req.PuasaSunnah
	}
	if req.WudhuSebelumTidur != nil {
		mu.WudhuSebelumTidur = *req.WudhuSebelumTidur
	}
	if req.BacaDoaTidur != nil {
		mu.BacaDoaTidur = *req.BacaDoaTidur
	}
	mu.Notes = req.Notes

	query := `INSERT INTO mutabaah_yaumiyah (id, student_id, date, fajr, dhuhr, asr, maghrib, isha, tahajjud, dhuha, sunnah, quran_tilawah, quran_hifdz, dzikr_pagi, dzikr_petang, shadaqah, puasa_sunnah, wudhu_sebelum_tidur, baca_doa_tidur, notes) VALUES (:id, :student_id, :date, :fajr, :dhuhr, :asr, :maghrib, :isha, :tahajjud, :dhuha, :sunnah, :quran_tilawah, :quran_hifdz, :dzikr_pagi, :dzikr_petang, :shadaqah, :puasa_sunnah, :wudhu_sebelum_tidur, :baca_doa_tidur, :notes)`
	if _, err := s.db.NamedExecContext(ctx, query, mu); err != nil {
		return nil, domain.NewInternalError("failed to create mutabaah", err)
	}

	return &dto.MutabaahResponse{
		ID:                mu.ID,
		StudentID:         mu.StudentID,
		Date:              mu.Date,
		Fajr:              mu.Fajr,
		Dhuhr:             mu.Dhuhr,
		Asr:               mu.Asr,
		Maghrib:           mu.Maghrib,
		Isha:              mu.Isha,
		Tahajjud:          mu.Tahajjud,
		Dhuha:             mu.Dhuha,
		Sunnah:            mu.Sunnah,
		QuranTilawah:      mu.QuranTilawah,
		QuranHifdz:        mu.QuranHifdz,
		DzikrPagi:         mu.DzikrPagi,
		DzikrPetang:       mu.DzikrPetang,
		Shadaqah:          mu.Shadaqah,
		PuasaSunnah:       mu.PuasaSunnah,
		WudhuSebelumTidur: mu.WudhuSebelumTidur,
		BacaDoaTidur:      mu.BacaDoaTidur,
		Notes:             mu.Notes,
	}, nil
}

func (s *islamicService) ListMutabaah(ctx context.Context, studentID string, startDate, endDate time.Time) ([]dto.MutabaahResponse, error) {
	var items []domain.MutabaahYaumiyah
	if err := s.db.SelectContext(ctx, &items, database.ListMutabaah, studentID, startDate, endDate); err != nil {
		return nil, domain.NewInternalError("failed to list mutabaah", err)
	}

	result := make([]dto.MutabaahResponse, len(items))
	for i, m := range items {
		result[i] = dto.MutabaahResponse{
			ID:                m.ID,
			StudentID:         m.StudentID,
			Date:              m.Date,
			Fajr:              m.Fajr,
			Dhuhr:             m.Dhuhr,
			Asr:               m.Asr,
			Maghrib:           m.Maghrib,
			Isha:              m.Isha,
			Tahajjud:          m.Tahajjud,
			Dhuha:             m.Dhuha,
			Sunnah:            m.Sunnah,
			QuranTilawah:      m.QuranTilawah,
			QuranHifdz:        m.QuranHifdz,
			DzikrPagi:         m.DzikrPagi,
			DzikrPetang:       m.DzikrPetang,
			Shadaqah:          m.Shadaqah,
			PuasaSunnah:       m.PuasaSunnah,
			WudhuSebelumTidur: m.WudhuSebelumTidur,
			BacaDoaTidur:      m.BacaDoaTidur,
			Notes:             m.Notes,
		}
	}
	return result, nil
}

func (s *islamicService) CreatePrayerAttendance(ctx context.Context, req dto.CreatePrayerAttendanceRequest) (*dto.PrayerAttendanceResponse, error) {
	pa := &domain.PrayerAttendance{
		ID:        uuid.New().String(),
		StudentID: req.StudentID,
		Date:      req.Date,
		Notes:     req.Notes,
	}
	if req.Fajr != nil {
		pa.Fajr = *req.Fajr
	}
	if req.Dhuhr != nil {
		pa.Dhuhr = *req.Dhuhr
	}
	if req.Asr != nil {
		pa.Asr = *req.Asr
	}
	if req.Maghrib != nil {
		pa.Maghrib = *req.Maghrib
	}
	if req.Isha != nil {
		pa.Isha = *req.Isha
	}

	query := `INSERT INTO prayer_attendances (id, student_id, date, fajr, dhuhr, asr, maghrib, isha, notes) VALUES (:id, :student_id, :date, :fajr, :dhuhr, :asr, :maghrib, :isha, :notes)`
	if _, err := s.db.NamedExecContext(ctx, query, pa); err != nil {
		return nil, domain.NewInternalError("failed to create prayer attendance", err)
	}

	return &dto.PrayerAttendanceResponse{
		ID:        pa.ID,
		StudentID: pa.StudentID,
		Date:      pa.Date,
		Fajr:      pa.Fajr,
		Dhuhr:     pa.Dhuhr,
		Asr:       pa.Asr,
		Maghrib:   pa.Maghrib,
		Isha:      pa.Isha,
		Notes:     pa.Notes,
	}, nil
}

func (s *islamicService) ListPrayerAttendance(ctx context.Context, studentID string, startDate, endDate time.Time) ([]dto.PrayerAttendanceResponse, error) {
	var items []domain.PrayerAttendance
	if err := s.db.SelectContext(ctx, &items, database.ListPrayerAttendance, studentID, startDate, endDate); err != nil {
		return nil, domain.NewInternalError("failed to list prayer attendance", err)
	}

	result := make([]dto.PrayerAttendanceResponse, len(items))
	for i, pa := range items {
		result[i] = dto.PrayerAttendanceResponse{
			ID:        pa.ID,
			StudentID: pa.StudentID,
			Date:      pa.Date,
			Fajr:      pa.Fajr,
			Dhuhr:     pa.Dhuhr,
			Asr:       pa.Asr,
			Maghrib:   pa.Maghrib,
			Isha:      pa.Isha,
			Notes:     pa.Notes,
		}
	}
	return result, nil
}

func (s *islamicService) CreateHalaqahGroup(ctx context.Context, schoolID string, req dto.CreateHalaqahGroupRequest) (*dto.HalaqahGroupResponse, error) {
	hg := &domain.HalaqahGroup{
		ID:        uuid.New().String(),
		SchoolID:  schoolID,
		Name:      req.Name,
		TeacherID: req.TeacherID,
		Room:      req.Room,
		Day:       req.Day,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		MaxMember: req.MaxMember,
	}

	query := `INSERT INTO halaqah_groups (id, school_id, name, teacher_id, room, day, start_time, end_time, max_member) VALUES (:id, :school_id, :name, :teacher_id, :room, :day, :start_time, :end_time, :max_member)`
	if _, err := s.db.NamedExecContext(ctx, query, hg); err != nil {
		return nil, domain.NewInternalError("failed to create halaqah group", err)
	}

	return &dto.HalaqahGroupResponse{
		ID:        hg.ID,
		Name:      hg.Name,
		TeacherID: hg.TeacherID,
		Room:      hg.Room,
		Day:       hg.Day,
		StartTime: hg.StartTime,
		EndTime:   hg.EndTime,
		MaxMember: hg.MaxMember,
	}, nil
}

func (s *islamicService) ListHalaqahGroups(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.HalaqahGroupResponse, int64, error) {
	filter.Defaults()

	type halaqahRow struct {
		domain.HalaqahGroup
		TeacherName string `db:"teacher_name"`
		MemberCount int    `db:"member_count"`
	}

	var rows []halaqahRow
	if err := s.db.SelectContext(ctx, &rows, database.ListHalaqahGroups, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list halaqah groups", err)
	}

	var total int64
	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM halaqah_groups WHERE school_id=$1`, schoolID)

	result := make([]dto.HalaqahGroupResponse, len(rows))
	for i, r := range rows {
		result[i] = dto.HalaqahGroupResponse{
			ID:          r.ID,
			Name:        r.Name,
			TeacherID:   r.TeacherID,
			TeacherName: r.TeacherName,
			Room:        r.Room,
			Day:         r.Day,
			StartTime:   r.StartTime,
			EndTime:     r.EndTime,
			MaxMember:   r.MaxMember,
			MemberCount: r.MemberCount,
		}
	}
	return result, total, nil
}

func (s *islamicService) GetHalaqahGroup(ctx context.Context, id string) (*dto.HalaqahGroupResponse, error) {
	var hg domain.HalaqahGroup
	if err := s.db.GetContext(ctx, &hg, `SELECT * FROM halaqah_groups WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("halaqah group", id)
	}

	var members []struct {
		domain.HalaqahMember
		FullName string `db:"full_name"`
	}
	if err := s.db.SelectContext(ctx, &members, database.ListHalaqahMembers, id); err != nil {
		s.logger.Warn("failed to get members", zap.Error(err))
	}

	memberResponses := make([]dto.HalaqahMemberResponse, len(members))
	for i, m := range members {
		memberResponses[i] = dto.HalaqahMemberResponse{
			ID:        m.ID,
			StudentID: m.StudentID,
			FullName:  m.FullName,
			JoinedAt:  m.JoinedAt,
		}
	}

	return &dto.HalaqahGroupResponse{
		ID:          hg.ID,
		Name:        hg.Name,
		TeacherID:   hg.TeacherID,
		Room:        hg.Room,
		Day:         hg.Day,
		StartTime:   hg.StartTime,
		EndTime:     hg.EndTime,
		MaxMember:   hg.MaxMember,
		MemberCount: len(members),
		Members:     memberResponses,
	}, nil
}

func (s *islamicService) AddHalaqahMember(ctx context.Context, halaqahID string, req dto.AddHalaqahMemberRequest) error {
	hm := &domain.HalaqahMember{
		ID:        uuid.New().String(),
		HalaqahID: halaqahID,
		StudentID: req.StudentID,
	}

	query := `INSERT INTO halaqah_members (id, halaqah_id, student_id) VALUES (:id, :halaqah_id, :student_id)`
	if _, err := s.db.NamedExecContext(ctx, query, hm); err != nil {
		return domain.NewInternalError("failed to add member", err)
	}
	return nil
}
