package mcp

import (
	"context"
	"fmt"
	"time"
)

func (s *MCPServer) RegisterResources() {
	s.registerSchoolResources()
	s.registerStudentResources()
	s.registerFinancialResources()
	s.registerIslamicResources()
	s.registerGeneralResources()
}

func (s *MCPServer) registerSchoolResources() {
	resources := []*Resource{
		{
			URI:         "erp://school/profile",
			Name:        "School Profile",
			Description: "School profile information including vision, mission, accreditation, and basic data",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "name": "Islamic School Name",
  "npsn": "12345678",
  "type": "Islamic Integrated School",
  "accreditation": "A",
  "vision": "To develop righteous, knowledgeable, and globally competitive Muslim generations",
  "mission": "Providing holistic Islamic education integrating Quran, Sunnah, and modern sciences",
  "established": 2010,
  "levels": ["TK", "SD", "SMP", "SMA"],
  "curriculum": "Kurikulum Merdeka with Islamic Integration"
}`, "application/json", nil
			},
		},
		{
			URI:         "erp://school/academic-calendar",
			Name:        "Academic Calendar",
			Description: "Current academic year calendar with important dates",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				year := time.Now().Year()
				return fmt.Sprintf(`{
  "academic_year": "%d/%d",
  "semester": "Ganjil",
  "start_date": "%d-07-15",
  "end_date": "%d-06-20",
  "events": [
    {"name": "MPLS", "date": "%d-07-15", "type": "orientation"},
    {"name": "PTS Ganjil", "date": "%d-09-20", "type": "exam"},
    {"name": "PAS Ganjil", "date": "%d-12-01", "type": "exam"},
    {"name": "Islamic New Year", "date": "varies", "type": "holiday"},
    {"name": "Maulid Nabi", "date": "varies", "type": "holiday"},
    {"name": "Isra Mi'raj", "date": "varies", "type": "holiday"},
    {"name": "Ramadan Activities", "date": "varies", "type": "islamic"},
    {"name": "PTS Genap", "date": "%d-03-15", "type": "exam"},
    {"name": "PAS Genap", "date": "%d-06-01", "type": "exam"}
  ],
  "holidays": [
    "Islamic holidays based on Hijri calendar",
    "National holidays",
    "Semester breaks"
  ]
}`, year, year+1, year, year+1, year, year, year, year+1, year+1), "application/json", nil
			},
		},
		{
			URI:         "erp://school/teacher-profile",
			Name:        "Teacher Profile",
			Description: "Teacher profile and qualification information",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "total_teachers": 45,
  "certified_teachers": 30,
  "education_levels": {
    "S2": 5,
    "S1": 38,
    "D4": 2
  },
  "departments": ["Islamic Studies", "Arabic", "Mathematics", "Science", "Social Studies", "Languages", "Physical Education", "IT"],
  "gender_ratio": {"male": 20, "female": 25},
  "average_experience_years": 7.5
}`, "application/json", nil
			},
		},
	}

	for _, r := range resources {
		s.RegisterResource(r)
	}
}

func (s *MCPServer) registerStudentResources() {
	resources := []*Resource{
		{
			URI:         "erp://students/enrollment",
			Name:        "Student Enrollment Data",
			Description: "Current student enrollment statistics by level and program",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "total_students": 850,
  "levels": {
    "TK": {"classes": 2, "students": 40},
    "SD": {"classes": 12, "students": 360},
    "SMP": {"classes": 9, "students": 270},
    "SMA": {"classes": 6, "students": 180}
  },
  "gender_ratio": {"male": 420, "female": 430},
  "tahfidz_program": {"students": 200, "levels": ["juz_amma", "juz_1_5", "juz_6_10", "juz_11_20", "juz_21_30"]},
  "language_program": {"arabic": 500, "english": 850},
  "scholarship_recipients": 85
}`, "application/json", nil
			},
		},
		{
			URI:         "erp://students/attendance",
			Name:        "Student Attendance Records",
			Description: "Student attendance summary and patterns",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "average_attendance_rate": 95.5,
  "monthly_trend": {
    "January": 96.2, "February": 95.8, "March": 94.5,
    "April": 96.0, "May": 95.3, "June": 93.2,
    "July": 97.1, "August": 96.5, "September": 95.9,
    "October": 95.2, "November": 94.8, "December": 93.5
  },
  "common_absence_reasons": [
    "Illness", "Family events", "Islamic activities/competitions",
    "Weather conditions", "Transportation issues"
  ]
}`, "application/json", nil
			},
		},
		{
			URI:         "erp://students/achievements",
			Name:        "Student Achievements",
			Description: "Student academic and non-academic achievement records",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "academic": {
    "national_olympiad_winners": 12,
    "regional_competitions": 25,
    "average_exam_score": 82.5
  },
  "islamic": {
    "musabaqah_tilawatil_quran_winners": 8,
    "musabaqah_hifdzil_quran_winners": 5,
    "arabic_debate_winners": 3
  },
  "sports": {
    "national_level": 5,
    "regional_level": 18
  },
  "arts": {
    "calligraphy_winners": 7,
    "nasheed_winners": 4
  }
}`, "application/json", nil
			},
		},
	}

	for _, r := range resources {
		s.RegisterResource(r)
	}
}

func (s *MCPServer) registerFinancialResources() {
	resources := []*Resource{
		{
			URI:         "erp://finance/budget-overview",
			Name:        "Budget Overview",
			Description: "Current year budget overview and allocation",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "total_budget": 5000000000,
  "revenue_sources": {
    "SPP": {"amount": 3000000000, "percentage": 60},
    "BOS": {"amount": 800000000, "percentage": 16},
    "Donations": {"amount": 500000000, "percentage": 10},
    "Other": {"amount": 700000000, "percentage": 14}
  },
  "expenditure_categories": {
    "Personnel": {"amount": 2500000000, "percentage": 50},
    "Operations": {"amount": 1000000000, "percentage": 20},
    "Facilities": {"amount": 750000000, "percentage": 15},
    "Programs": {"amount": 500000000, "percentage": 10},
    "Reserve": {"amount": 250000000, "percentage": 5}
  }
}`, "application/json", nil
			},
		},
		{
			URI:         "erp://finance/fee-structure",
			Name:        "Fee Structure",
			Description: "School fee structure and payment policies",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "registration_fee": 5000000,
  "monthly_spp": {
    "TK": 500000,
    "SD": 750000,
    "SMP": 1000000,
    "SMA": 1250000
  },
  "additional_fees": {
    "books": {"annual": 1500000},
    "uniform": {"annual": 1000000},
    "extracurricular": {"per_activity": 250000},
    "field_trip": {"per_event": 500000}
  },
  "payment_methods": ["Bank Transfer", "Virtual Account", "QRIS", "Cash"],
  "late_payment_policy": "2% per month after due date",
  "scholarship_discounts": ["Prestasi (25-100%)", "Hafidz (50-100%)", "Dhuafa (50-100%)", "Sibling (10-25%)"]
}`, "application/json", nil
			},
		},
	}

	for _, r := range resources {
		s.RegisterResource(r)
	}
}

func (s *MCPServer) registerIslamicResources() {
	resources := []*Resource{
		{
			URI:         "erp://islamic/tahfidz-program",
			Name:        "Tahfidz Program Structure",
			Description: "Quran memorization program structure and curriculum",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "program_name": "Tahfidz Al-Quran Program",
  "levels": [
    {"name": "Tahsin", "description": "Quran reading improvement", "duration": "6 months"},
    {"name": "Juz Amma", "description": "Juz 30 memorization", "duration": "1 year"},
    {"name": "Juz 1-5", "description": "Early juz memorization", "duration": "1.5 years"},
    {"name": "Juz 6-15", "description": "Intermediate memorization", "duration": "2 years"},
    {"name": "Juz 16-25", "description": "Advanced memorization", "duration": "2 years"},
    {"name": "Juz 26-30", "description": "Completion", "duration": "1.5 years"}
  ],
  "daily_schedule": {
    "morning_murajaah": "05:30-06:30",
    "new_memorization": "06:30-07:30",
    "afternoon_murajaah": "13:30-14:30",
    "evening_tasmi": "After Maghrib"
  },
  "evaluation_methods": ["Monthly tasmi'", "Semester exam", "Annual khataman"],
  "teaching_methods": ["Talaqqi", "Sima'i", "Kitabah", "Jama'"]
}`, "application/json", nil
			},
		},
		{
			URI:         "erp://islamic/islamic-calendar",
			Name:        "Islamic Calendar Events",
			Description: "Hijri calendar with important Islamic dates and school activities",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				return `{
  "hijri_months": [
    "Muharram", "Safar", "Rabiul Awal", "Rabiul Akhir",
    "Jumadil Awal", "Jumadil Akhir", "Rajab", "Sya'ban",
    "Ramadan", "Syawal", "Dzulqa'dah", "Dzulhijjah"
  ],
  "important_dates": [
    {"event": "Islamic New Year (1 Muharram)", "type": "holiday"},
    {"event": "Ashura (10 Muharram)", "type": "sunnah_fasting"},
    {"event": "Maulid Nabi Muhammad SAW (12 Rabiul Awal)", "type": "celebration"},
    {"event": "Isra Mi'raj (27 Rajab)", "type": "celebration"},
    {"event": "Nisfu Sya'ban (15 Sya'ban)", "type": "worship"},
    {"event": "Ramadan (1-30 Ramadan)", "type": "fasting_month"},
    {"event": "Nuzulul Quran (17 Ramadan)", "type": "celebration"},
    {"event": "Lailatul Qadr (odd nights last 10 days)", "type": "worship"},
    {"event": "Idul Fitri (1 Syawal)", "type": "holiday"},
    {"event": "Idul Adha (10 Dzulhijjah)", "type": "holiday"}
  ],
  "school_islamic_activities": [
    "Daily Dhuha prayer",
    "Weekly Islamic studies (Keputrian/Keputraan)",
    "Monthly khataman",
    "Ramadan pesantren kilat",
    "Qurban training",
    "Islamic competitions"
  ]
}`, "application/json", nil
			},
		},
	}

	for _, r := range resources {
		s.RegisterResource(r)
	}
}

func (s *MCPServer) registerGeneralResources() {
	resources := []*Resource{
		{
			URI:         "erp://system/health",
			Name:        "System Health",
			Description: "Current system health status and provider availability",
			MimeType:    "application/json",
			Handler: func(ctx context.Context, uri string) (string, string, error) {
				providers := s.router.GetEnabledProviders()
				status := make(map[string]string)
				for _, p := range providers {
					if p.IsAvailable(ctx) {
						status[p.Name()] = "healthy"
					} else {
						status[p.Name()] = "unavailable"
					}
				}
				return fmt.Sprintf(`{
  "status": "operational",
  "version": "1.0.0",
  "uptime": "running",
  "providers": %s
}`, formatProviderStatus(status)), "application/json", nil
			},
		},
	}

	for _, r := range resources {
		s.RegisterResource(r)
	}
}

func formatProviderStatus(status map[string]string) string {
	if len(status) == 0 {
		return "{}"
	}
	result := "{\n"
	for name, s := range status {
		result += fmt.Sprintf(`    "%s": "%s",`, name, s) + "\n"
	}
	result = result[:len(result)-2] + "\n  }"
	return result
}
