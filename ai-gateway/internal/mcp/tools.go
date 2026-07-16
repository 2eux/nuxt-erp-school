package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/opencode/erp-ai-gateway/internal/providers"
	"go.uber.org/zap"
)

func (s *MCPServer) RegisterTools() {
	s.registerAcademicTools()
	s.registerIslamicTools()
	s.registerFinancialTools()
	s.registerAdministrativeTools()
	s.registerHRTools()
	s.registerGeneralTools()
}

func (s *MCPServer) registerAcademicTools() {
	tools := []*Tool{
		{
			Name:        "lesson_plan_generator",
			Description: "Generate comprehensive lesson plans for any subject and grade level, aligned with curriculum standards",
			Category:    "academic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"subject":      map[string]any{"type": "string", "description": "Subject name (e.g., Mathematics, Science, Islamic Studies)"},
					"grade":        map[string]any{"type": "string", "description": "Grade level"},
					"topic":        map[string]any{"type": "string", "description": "Specific topic for the lesson"},
					"duration":     map[string]any{"type": "string", "description": "Lesson duration in minutes"},
					"objectives":   map[string]any{"type": "string", "description": "Learning objectives, comma separated"},
					"curriculum":   map[string]any{"type": "string", "description": "Curriculum standard (e.g., K13, Kurikulum Merdeka, International)"},
					"language":     map[string]any{"type": "string", "description": "Language (en/id/ar)"},
				},
				"required": []string{"subject", "grade", "topic"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "lesson_plan_generator", args, `You are an expert educator. Generate a detailed lesson plan with:
1. Lesson title and overview
2. Learning objectives (SMART)
3. Materials and resources needed
4. Lesson procedure (opening, main activity, closing)
5. Assessment methods
6. Differentiation strategies
7. Homework/additional practice
8. Reflection notes for teacher

Format the output in clear sections with markdown.`)
			},
		},
		{
			Name:        "quiz_generator",
			Description: "Generate quizzes with multiple question types for any subject",
			Category:    "academic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"subject":      map[string]any{"type": "string", "description": "Subject"},
					"topic":        map[string]any{"type": "string", "description": "Topic"},
					"grade":        map[string]any{"type": "string", "description": "Grade level"},
					"num_questions": map[string]any{"type": "number", "description": "Number of questions"},
					"types":         map[string]any{"type": "string", "description": "Question types: multiple_choice,true_false,essay,short_answer"},
					"difficulty":    map[string]any{"type": "string", "description": "easy/medium/hard/mixed"},
					"language":      map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"subject", "topic", "grade"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "quiz_generator", args, `You are an expert quiz creator. Generate a quiz with:
- Clear instructions
- Questions with varying difficulty
- Answer key at the end
- Include the question type requested
- Format properly in markdown`)
			},
		},
		{
			Name:        "exam_generator",
			Description: "Generate comprehensive exams for mid-term, final, or unit tests",
			Category:    "academic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"subject":      map[string]any{"type": "string", "description": "Subject"},
					"grade":        map[string]any{"type": "string", "description": "Grade"},
					"topics":       map[string]any{"type": "string", "description": "Topics covered (comma separated)"},
					"exam_type":    map[string]any{"type": "string", "description": "mid_term/final/unit_test"},
					"total_marks":   map[string]any{"type": "number", "description": "Total marks"},
					"language":      map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"subject", "grade", "topics"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "exam_generator", args, `You are an expert exam designer. Generate a complete exam with:
1. Exam header (school name, subject, grade, date, duration, total marks)
2. General instructions
3. Sections with different question types
4. Mark allocation per section/question
5. Answer key and marking scheme
6. Bloom's taxonomy alignment notes
Format in markdown.`)
			},
		},
		{
			Name:        "worksheet_generator",
			Description: "Generate student worksheets and practice materials",
			Category:    "academic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"subject":    map[string]any{"type": "string", "description": "Subject"},
					"topic":      map[string]any{"type": "string", "description": "Topic"},
					"grade":      map[string]any{"type": "string", "description": "Grade"},
					"pages":      map[string]any{"type": "number", "description": "Number of pages"},
					"activities": map[string]any{"type": "string", "description": "Activity types (fill_blanks,matching,crossword,word_search,etc)"},
					"language":   map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"subject", "topic", "grade"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "worksheet_generator", args, `Create an engaging student worksheet with varied activities. Include:
- Title and learning objective
- Clear instructions for each section
- Mix of activity types
- Space for student name and date
- Answer key
Format in markdown.`)
			},
		},
		{
			Name:        "rubric_generator",
			Description: "Generate assessment rubrics for assignments and projects",
			Category:    "academic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"assignment_type": map[string]any{"type": "string", "description": "Type of assignment"},
					"criteria":        map[string]any{"type": "string", "description": "Assessment criteria (comma separated)"},
					"levels":          map[string]any{"type": "number", "description": "Number of performance levels (default 4)"},
					"max_score":       map[string]any{"type": "number", "description": "Maximum score"},
				},
				"required": []string{"assignment_type", "criteria"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "rubric_generator", args, `Create a detailed assessment rubric as a markdown table with:
- Performance levels as columns
- Criteria as rows
- Clear descriptors for each cell
- Score ranges
Format as a table.`)
			},
		},
		{
			Name:        "homework_generator",
			Description: "Generate homework assignments with practice problems",
			Category:    "academic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"subject":   map[string]any{"type": "string", "description": "Subject"},
					"topic":     map[string]any{"type": "string", "description": "Topic"},
					"grade":     map[string]any{"type": "string", "description": "Grade"},
					"problems":  map[string]any{"type": "number", "description": "Number of problems"},
					"due_date":  map[string]any{"type": "string", "description": "Due date description"},
					"language":  map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"subject", "topic", "grade"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "homework_generator", args, `Generate homework with:
- Assignment title and instructions
- Practice problems
- Bonus/extension questions
- Space for student work
- Answer key for teacher
Format in markdown.`)
			},
		},
		{
			Name:        "student_feedback_generator",
			Description: "Generate personalized student feedback based on performance data",
			Category:    "academic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_name":  map[string]any{"type": "string", "description": "Student name"},
					"subject":       map[string]any{"type": "string", "description": "Subject"},
					"score":         map[string]any{"type": "number", "description": "Score achieved"},
					"max_score":     map[string]any{"type": "number", "description": "Maximum score"},
					"strengths":     map[string]any{"type": "string", "description": "Student strengths"},
					"weaknesses":    map[string]any{"type": "string", "description": "Areas for improvement"},
					"behavior_notes": map[string]any{"type": "string", "description": "Behavior observations"},
					"language":      map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"student_name", "subject", "score"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "student_feedback_generator", args, `Write constructive, personalized student feedback that is:
- Specific and actionable
- Encouraging and supportive
- Balanced (strengths and areas to improve)
- Professional tone suitable for parents
- Include specific suggestions for improvement`)
			},
		},
		{
			Name:        "report_card_comment_generator",
			Description: "Generate report card comments for students",
			Category:    "academic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_name":    map[string]any{"type": "string", "description": "Student name"},
					"grade":           map[string]any{"type": "string", "description": "Grade"},
					"semester":        map[string]any{"type": "string", "description": "Semester (1 or 2)"},
					"academic_performance": map[string]any{"type": "string", "description": "Academic performance summary"},
					"attendance":      map[string]any{"type": "string", "description": "Attendance record"},
					"extracurricular": map[string]any{"type": "string", "description": "Extracurricular activities"},
					"character_notes": map[string]any{"type": "string", "description": "Character/manner notes"},
					"language":        map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"student_name", "grade", "academic_performance"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "report_card_comment_generator", args, `Write professional report card comments covering:
- Academic achievement summary
- Effort and participation
- Social and behavioral development
- Areas for growth
- Recommendations for next semester
Keep it professional, constructive, and parent-friendly.`)
			},
		},
	}

	for _, tool := range tools {
		s.RegisterTool(tool)
	}
}

func (s *MCPServer) registerIslamicTools() {
	tools := []*Tool{
		{
			Name:        "tahfidz_planner",
			Description: "Create personalized Quran memorization (tahfidz) plans for students",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_name":        map[string]any{"type": "string", "description": "Student name"},
					"current_memorization": map[string]any{"type": "string", "description": "Current juz and surah memorized"},
					"target":              map[string]any{"type": "string", "description": "Target (e.g., 30 juz, specific surahs)"},
					"daily_capacity":      map[string]any{"type": "string", "description": "Daily memorization capacity (pages/verses)"},
					"timeframe_months":    map[string]any{"type": "number", "description": "Target timeframe in months"},
					"language":            map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"student_name", "current_memorization", "target"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "tahfidz_planner", args, `You are an expert Quran teacher (ustadz/ustadzah). Create a detailed tahfidz plan with:
1. Current status assessment
2. Weekly targets with daily breakdown
3. Muraja'ah (review) schedule
4. Tips for effective memorization
5. Milestone celebration suggestions
6. Progress tracking template
Include Islamic motivation and relevant hadith about Quran memorization.
Format beautifully in markdown.`)
			},
		},
		{
			Name:        "memorization_recommendation",
			Description: "Provide personalized memorization technique recommendations",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_name":   map[string]any{"type": "string", "description": "Student name"},
					"age":            map[string]any{"type": "number", "description": "Student age"},
					"learning_style": map[string]any{"type": "string", "description": "Learning style (visual/auditory/kinesthetic)"},
					"challenges":     map[string]any{"type": "string", "description": "Specific challenges faced"},
					"current_level":  map[string]any{"type": "string", "description": "Current memorization level"},
				},
				"required": []string{"student_name", "age", "learning_style"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "memorization_recommendation", args, `As a Quran memorization expert, provide:
1. Personalized technique recommendations based on learning style
2. Daily routine suggestions
3. Memory aid tools and methods
4. Common pitfalls and how to avoid them
5. Motivational advice with Islamic references
6. Parent/teacher support guide`)
			},
		},
		{
			Name:        "tilawah_assistant",
			Description: "Provide tilawah (Quran recitation) guidance and analysis",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"surah":         map[string]any{"type": "string", "description": "Surah name or number"},
					"ayah_range":    map[string]any{"type": "string", "description": "Ayah range (e.g., 1-10)"},
					"focus_area":    map[string]any{"type": "string", "description": "tajwid/tartil/makhraj/all"},
					"reciter_style": map[string]any{"type": "string", "description": "Preferred qira'ah style"},
					"language":      map[string]any{"type": "string", "description": "Language for explanation"},
				},
				"required": []string{"surah"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "tilawah_assistant", args, `Provide tilawah guidance with:
1. Surah introduction and context (asbabun nuzul if relevant)
2. Key tajwid rules in the verses
3. Makhraj (articulation) notes
4. Waqf (stopping) guidance
5. Recommended practice approach
6. Common mistakes to avoid
7. Audio resource recommendations`)
			},
		},
		{
			Name:        "tajwid_assistant",
			Description: "Provide detailed tajwid rule explanations and examples",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"rule_name":    map[string]any{"type": "string", "description": "Tajwid rule name (idgham, ikhfa, iqlab, etc.)"},
					"surah_example": map[string]any{"type": "string", "description": "Specific surah for examples"},
					"difficulty":   map[string]any{"type": "string", "description": "Student level (beginner/intermediate/advanced)"},
					"language":     map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"rule_name", "difficulty"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "tajwid_assistant", args, `Explain the tajwid rule with:
1. Definition and Arabic terminology
2. How to identify it in the mushaf
3. Multiple Quranic examples with ayah references
4. Step-by-step pronunciation guide
5. Practice exercises
6. Common errors and corrections
7. Visual representation if applicable`)
			},
		},
		{
			Name:        "arabic_vocabulary_trainer",
			Description: "Generate Arabic vocabulary training materials",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"theme":        map[string]any{"type": "string", "description": "Theme (daily_life/quranic/school/ibadah/conversation)"},
					"num_words":    map[string]any{"type": "number", "description": "Number of words"},
					"level":        map[string]any{"type": "string", "description": "Beginner/intermediate/advanced"},
					"include_harakat": map[string]any{"type": "boolean", "description": "Include harakat (vowel marks)"},
					"language":     map[string]any{"type": "string", "description": "Language for translations"},
				},
				"required": []string{"theme", "level"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "arabic_vocabulary_trainer", args, `Generate Arabic vocabulary training material:
1. Word in Arabic (with harakat if requested)
2. Transliteration
3. Translation
4. Example sentence
5. Root word analysis
6. Related words
7. Memory tips
Format as a table for easy study.`)
			},
		},
		{
			Name:        "islamic_quiz_generator",
			Description: "Generate Islamic knowledge quizzes",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"topic":       map[string]any{"type": "string", "description": "aqidah/fiqh/seerah/tafsir/hadith/akhlaq/tarikh"},
					"num_questions": map[string]any{"type": "number", "description": "Number of questions"},
					"difficulty":  map[string]any{"type": "string", "description": "easy/medium/hard/mixed"},
					"age_group":   map[string]any{"type": "string", "description": "children/teenager/adult"},
					"language":    map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"topic"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "islamic_quiz_generator", args, `Create an Islamic knowledge quiz:
1. Questions with authentic sources
2. Multiple choice with 4 options
3. Correct answer marked
4. Brief explanation/dalil for each answer
5. Source references (Quran ayah, hadith)
Format clearly in markdown.`)
			},
		},
		{
			Name:        "islamic_story_generator",
			Description: "Generate Islamic stories for children with moral lessons",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"story_type":    map[string]any{"type": "string", "description": "prophets/sahabah/islamic_history/moral_story"},
					"character":     map[string]any{"type": "string", "description": "Specific prophet/sahabi/historical figure"},
					"age_group":     map[string]any{"type": "string", "description": "Age group (4-6, 7-10, 11-14, 15+)"},
					"moral_lesson":  map[string]any{"type": "string", "description": "Moral lesson to focus on"},
					"story_length":  map[string]any{"type": "string", "description": "short/medium/long"},
					"language":      map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"story_type", "age_group"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "islamic_story_generator", args, `Create an engaging Islamic story:
1. Title
2. Introduction setting the context
3. The main story with dialogue
4. Quran/hadith references
5. Moral lesson and reflection questions
6. Related activities for children
7. Vocabulary list (Arabic terms used)
Be authentic while maintaining engagement.`)
			},
		},
		{
			Name:        "akhlaq_recommendation",
			Description: "Provide akhlaq (character) development recommendations",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_name":     map[string]any{"type": "string", "description": "Student name"},
					"current_behavior": map[string]any{"type": "string", "description": "Current behavioral observations"},
					"target_akhlaq":    map[string]any{"type": "string", "description": "Target character trait to develop"},
					"age":              map[string]any{"type": "number", "description": "Student age"},
					"language":         map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"student_name", "target_akhlaq"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "akhlaq_recommendation", args, `Provide akhlaq development guidance:
1. Islamic foundation for the character trait (Quran & Sunnah)
2. Prophet Muhammad SAW's example
3. Practical daily activities
4. Du'a recommendations
5. Progress tracking suggestions
6. Parent involvement guide
7. Stories of sahabah demonstrating this trait`)
			},
		},
		{
			Name:        "mutabaah_analyzer",
			Description: "Analyze students' mutaba'ah (daily ibadah records) for patterns and growth",
			Category:    "islamic",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_name":    map[string]any{"type": "string", "description": "Student name"},
					"period":          map[string]any{"type": "string", "description": "Analysis period (weekly/monthly/semester)"},
					"prayer_record":   map[string]any{"type": "string", "description": "Prayer consistency data"},
					"quran_reading":   map[string]any{"type": "string", "description": "Quran reading record"},
					"sunah_prayers":   map[string]any{"type": "string", "description": "Sunnah prayer records"},
					"dhikr_consistency": map[string]any{"type": "string", "description": "Dhikr consistency"},
					"language":        map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"student_name", "period"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeIslamicTool(ctx, "mutabaah_analyzer", args, `Analyze mutaba'ah record and provide:
1. Summary statistics and trends
2. Strengths and consistent practices
3. Areas needing improvement
4. Personalized recommendations
5. Encouraging Islamic reminders
6. Goal setting for next period
7. Parent communication points`)
			},
		},
	}

	for _, tool := range tools {
		s.RegisterTool(tool)
	}
}

func (s *MCPServer) registerFinancialTools() {
	tools := []*Tool{
		{
			Name:        "rkas_assistant",
			Description: "Assist in creating and reviewing RKAS (School Budget Plan)",
			Category:    "financial",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"school_name":    map[string]any{"type": "string", "description": "School name"},
					"budget_period":   map[string]any{"type": "string", "description": "Budget period (e.g., 2024/2025)"},
					"total_students":  map[string]any{"type": "number", "description": "Total students"},
					"total_teachers":  map[string]any{"type": "number", "description": "Total teachers"},
					"priorities":      map[string]any{"type": "string", "description": "Budget priorities"},
					"bop_fund":        map[string]any{"type": "number", "description": "BOS fund amount"},
					"spp_fund":        map[string]any{"type": "number", "description": "SPP fund amount"},
					"language":        map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"school_name", "budget_period"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "rkas_assistant", args, `As an Islamic school finance expert, provide RKAS guidance:
1. Revenue projection breakdown
2. Expenditure categories with allocations
3. 8 National Education Standards alignment
4. BOS fund utilization recommendations
5. Budget efficiency analysis
6. Risk mitigation strategies
7. Compliance checklist`)
			},
		},
		{
			Name:        "budget_recommendation",
			Description: "Provide budget allocation recommendations for school programs",
			Category:    "financial",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"total_budget":  map[string]any{"type": "number", "description": "Total budget amount"},
					"program_areas": map[string]any{"type": "string", "description": "Program areas (comma separated)"},
					"student_count": map[string]any{"type": "number", "description": "Student count"},
					"period":        map[string]any{"type": "string", "description": "Budget period"},
					"language":      map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"total_budget", "program_areas"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "budget_recommendation", args, `Provide budget allocation recommendations:
1. Percentage allocation per program area
2. Per-student cost analysis
3. Cost-saving opportunities
4. Priority-based allocation matrix
5. Islamic education-specific considerations
6. Monitoring and evaluation framework`)
			},
		},
		{
			Name:        "cashflow_predictor",
			Description: "Predict and analyze cashflow trends for school finance",
			Category:    "financial",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"income_sources":    map[string]any{"type": "string", "description": "Income sources and amounts (JSON or description)"},
					"expense_patterns":  map[string]any{"type": "string", "description": "Monthly expense patterns"},
					"forecast_months":   map[string]any{"type": "number", "description": "Months to forecast"},
					"reserve_fund":      map[string]any{"type": "number", "description": "Current reserve fund"},
					"language":          map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"income_sources", "forecast_months"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "cashflow_predictor", args, `Analyze cashflow and provide:
1. Monthly cashflow projection
2. Critical months identification
3. Recommended reserve fund levels
4. Revenue diversification suggestions
5. Expense optimization opportunities
6. Risk factors and mitigation
7. Visual description of trends`)
			},
		},
		{
			Name:        "payroll_assistant",
			Description: "Assist with payroll calculations and teacher compensation analysis",
			Category:    "financial",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"employee_count":    map[string]any{"type": "number", "description": "Total employees"},
					"teacher_count":     map[string]any{"type": "number", "description": "Teacher count"},
					"staff_count":       map[string]any{"type": "number", "description": "Staff count"},
					"salary_ranges":     map[string]any{"type": "string", "description": "Salary ranges description"},
					"benefits_scheme":   map[string]any{"type": "string", "description": "Benefits scheme"},
					"monthly_budget":    map[string]any{"type": "number", "description": "Monthly payroll budget"},
					"language":          map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"employee_count", "monthly_budget"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "payroll_assistant", args, `Provide payroll analysis:
1. Salary structure recommendations
2. Benefits package optimization
3. Tax compliance considerations (PPH 21)
4. BPJS calculations
5. Honor/gaji pokok + tunjangan breakdown
6. Teacher certification allowance (TPG)
7. Budget sustainability analysis`)
			},
		},
		{
			Name:        "financial_health_analyzer",
			Description: "Analyze overall financial health of the school",
			Category:    "financial",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"total_revenue":    map[string]any{"type": "number", "description": "Total annual revenue"},
					"total_expenses":   map[string]any{"type": "number", "description": "Total annual expenses"},
					"total_assets":     map[string]any{"type": "number", "description": "Total assets"},
					"total_liabilities": map[string]any{"type": "number", "description": "Total liabilities"},
					"receivables":      map[string]any{"type": "number", "description": "Accounts receivable (tunggakan)"},
					"student_count":    map[string]any{"type": "number", "description": "Current student count"},
					"language":         map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"total_revenue", "total_expenses"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "financial_health_analyzer", args, `Analyze financial health and provide:
1. Key financial ratios
2. Revenue vs expense trend analysis
3. Per-student cost analysis
4. Collection rate analysis
5. Financial sustainability score
6. Recommendations for improvement
7. Benchmarking against similar schools
8. Risk assessment`)
			},
		},
		{
			Name:        "donation_analytics",
			Description: "Analyze donation patterns and provide fundraising recommendations",
			Category:    "financial",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"donation_data":   map[string]any{"type": "string", "description": "Donation data summary"},
					"period":          map[string]any{"type": "string", "description": "Analysis period"},
					"donor_count":     map[string]any{"type": "number", "description": "Total donor count"},
					"campaign_types":  map[string]any{"type": "string", "description": "Campaign types (zakat/infaq/sedekah/waqaf/donasi)"},
					"language":        map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"donation_data", "period"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "donation_analytics", args, `Analyze donations and provide:
1. Donation trends over the period
2. Campaign type performance comparison
3. Donor retention analysis
4. Peak donation periods
5. Fundraising strategy recommendations
6. Zakat/infaq/sedekah compliance guidance
7. Communication templates for donors
8. Islamic fundraising best practices`)
			},
		},
	}

	for _, tool := range tools {
		s.RegisterTool(tool)
	}
}

func (s *MCPServer) registerAdministrativeTools() {
	tools := []*Tool{
		{
			Name:        "meeting_minutes_generator",
			Description: "Generate professional meeting minutes from notes",
			Category:    "administrative",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"meeting_title":     map[string]any{"type": "string", "description": "Meeting title"},
					"date":              map[string]any{"type": "string", "description": "Meeting date"},
					"attendees":         map[string]any{"type": "string", "description": "Attendees list"},
					"agenda_items":      map[string]any{"type": "string", "description": "Agenda items"},
					"discussions_notes": map[string]any{"type": "string", "description": "Discussion notes/transcript"},
					"decisions":         map[string]any{"type": "string", "description": "Decisions made"},
					"action_items":      map[string]any{"type": "string", "description": "Action items"},
					"language":          map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"meeting_title", "date", "discussions_notes"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "meeting_minutes_generator", args, `Generate professional meeting minutes:
1. Header (title, date, time, location, attendees)
2. Agenda items reviewed
3. Discussion summary per agenda item
4. Decisions and resolutions
5. Action items with PIC and deadlines
6. Next meeting schedule
7. Approval section
Format as formal document.`)
			},
		},
		{
			Name:        "letter_generator",
			Description: "Generate official school letters and correspondence",
			Category:    "administrative",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"letter_type":   map[string]any{"type": "string", "description": "pemberitahuan/undangan/permohonan/pengumuman/surat_keputusan/edaran"},
					"recipient":     map[string]any{"type": "string", "description": "Recipient (orang tua/instansi/guru/siswa)"},
					"subject":       map[string]any{"type": "string", "description": "Letter subject"},
					"content_points": map[string]any{"type": "string", "description": "Key content points"},
					"sender_info":   map[string]any{"type": "string", "description": "Sender information"},
					"language":      map[string]any{"type": "string", "description": "en/id"},
				},
				"required": []string{"letter_type", "subject", "content_points"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "letter_generator", args, `Generate a formal school letter:
1. Letterhead (kop surat)
2. Reference number placeholder
3. Date and address
4. Salutation
5. Body content
6. Closing
7. Signature block
Follow Indonesian formal letter format (if in Bahasa Indonesia).`)
			},
		},
		{
			Name:        "sop_generator",
			Description: "Generate Standard Operating Procedures for school processes",
			Category:    "administrative",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"process_name":     map[string]any{"type": "string", "description": "Process name"},
					"department":       map[string]any{"type": "string", "description": "Department"},
					"purpose":          map[string]any{"type": "string", "description": "SOP purpose"},
					"scope":            map[string]any{"type": "string", "description": "SOP scope"},
					"steps_description": map[string]any{"type": "string", "description": "Step descriptions"},
					"stakeholders":     map[string]any{"type": "string", "description": "Involved stakeholders"},
					"language":         map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"process_name", "purpose"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "sop_generator", args, `Create a detailed SOP document:
1. Document control (nomor, revisi, tanggal)
2. Purpose and objectives
3. Scope and applicability
4. Definitions and acronyms
5. Step-by-step procedure with flowchart description
6. Roles and responsibilities
7. Quality indicators
8. Related documents and references
9. Approval section`)
			},
		},
		{
			Name:        "policy_assistant",
			Description: "Assist in drafting and reviewing school policies",
			Category:    "administrative",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"policy_area":     map[string]any{"type": "string", "description": "Policy area (academic/discipline/attendance/admission/financial/safety)"},
					"policy_title":    map[string]any{"type": "string", "description": "Policy title"},
					"existing_policy": map[string]any{"type": "string", "description": "Existing policy text (for review)"},
					"requirements":    map[string]any{"type": "string", "description": "Policy requirements"},
					"islamic_school":  map[string]any{"type": "boolean", "description": "Whether this is for Islamic school"},
					"language":        map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"policy_area", "policy_title"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "policy_assistant", args, `Draft a comprehensive school policy:
1. Policy title and number
2. Rationale and legal basis
3. Policy statement
4. Definitions
5. Procedures
6. Rights and responsibilities
7. Consequences/sanctions
8. Appeal process
9. Review and revision schedule
Include Islamic values consideration for Islamic schools.`)
			},
		},
	}

	for _, tool := range tools {
		s.RegisterTool(tool)
	}
}

func (s *MCPServer) registerHRTools() {
	tools := []*Tool{
		{
			Name:        "hr_assistant",
			Description: "Assist with HR management tasks and teacher development",
			Category:    "hr",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"task_type":    map[string]any{"type": "string", "description": "performance_review/training_plan/conflict_resolution/retention/job_description"},
					"employee_info": map[string]any{"type": "string", "description": "Employee information"},
					"context":      map[string]any{"type": "string", "description": "Additional context"},
					"language":     map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"task_type"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "hr_assistant", args, `Provide HR assistance with:
1. Professional and Islamic approach
2. Legal compliance (UU Ketenagakerjaan)
3. Best practices specific to Islamic education
4. Actionable recommendations
5. Template documents where applicable
6. Follow-up plan suggestions`)
			},
		},
		{
			Name:        "recruitment_assistant",
			Description: "Assist with teacher and staff recruitment",
			Category:    "hr",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"position":         map[string]any{"type": "string", "description": "Position title"},
					"department":       map[string]any{"type": "string", "description": "Department"},
					"requirements":     map[string]any{"type": "string", "description": "Specific requirements"},
					"qualifications":   map[string]any{"type": "string", "description": "Required qualifications"},
					"islamic_required": map[string]any{"type": "boolean", "description": "Islamic qualifications needed"},
					"language":         map[string]any{"type": "string", "description": "Language"},
				},
				"required": []string{"position", "requirements"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "recruitment_assistant", args, `Generate recruitment materials:
1. Job description and specification
2. Job advertisement text
3. Interview questions (general + Islamic + technical)
4. Assessment criteria and scoring rubric
5. Onboarding checklist
6. Required documents list
Include Islamic education-specific requirements.`)
			},
		},
	}

	for _, tool := range tools {
		s.RegisterTool(tool)
	}
}

func (s *MCPServer) registerGeneralTools() {
	tools := []*Tool{
		{
			Name:        "document_rag_query",
			Description: "Query the school document knowledge base using RAG for contextual answers",
			Category:    "general",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query":            map[string]any{"type": "string", "description": "Natural language question about school documents"},
					"document_types":   map[string]any{"type": "string", "description": "Filter by document types (curriculum,sop,policy,financial,etc)"},
					"max_results":      map[string]any{"type": "number", "description": "Maximum results"},
					"include_sources":  map[string]any{"type": "boolean", "description": "Include source document references"},
				},
				"required": []string{"query"},
			},
			Handler: func(ctx context.Context, args map[string]any) (string, error) {
				return s.executeChatTool(ctx, "document_rag_query", args, `Answer based on the provided school documents context.
1. Answer directly using the provided context
2. Cite specific document sections when relevant
3. If context is insufficient, state clearly
4. Be concise but comprehensive
5. Use professional educational terminology`)
			},
		},
	}

	for _, tool := range tools {
		s.RegisterTool(tool)
	}
}

func (s *MCPServer) executeChatTool(ctx context.Context, toolName string, args map[string]any, systemPrompt string) (string, error) {
	provider, err := s.GetProvider(ctx)
	if err != nil {
		return "", fmt.Errorf("no provider available: %w", err)
	}

	argsJSON, _ := json.Marshal(args)

	messages := []providers.Message{
		{Role: providers.RoleSystem, Content: systemPrompt},
		{Role: providers.RoleUser, Content: fmt.Sprintf("Generate output for tool '%s' with these parameters:\n%s\n\nPlease provide the complete, well-formatted output directly.", toolName, string(argsJSON))},
	}

	cacheKey := fmt.Sprintf("mcp:%s:%x", toolName, hashString(string(argsJSON)))
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
			s.logger.Debug("Using cached MCP response", zap.String("tool", toolName))
			return cached, nil
		}
	}

	temp := 0.3
	maxTokens := 4096

	resp, err := provider.Chat(ctx, &providers.ChatRequest{
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	result := resp.Message.Content

	if s.cache != nil {
		ttl := 3600
		if strings.Contains(toolName, "lesson") || strings.Contains(toolName, "exam") {
			ttl = 86400
		}
		_ = s.cache.Set(ctx, cacheKey, result, timeDuration(ttl))
	}

	return result, nil
}

func (s *MCPServer) executeIslamicTool(ctx context.Context, toolName string, args map[string]any, systemPrompt string) (string, error) {
	provider, err := s.GetIslamicProvider(ctx)
	if err != nil {
		provider, err = s.GetProvider(ctx)
		if err != nil {
			return "", fmt.Errorf("no provider available: %w", err)
		}
	}

	argsJSON, _ := json.Marshal(args)

	messages := []providers.Message{
		{Role: providers.RoleSystem, Content: systemPrompt + "\n\nImportant: All references to Quran must be accurate. All hadith references must include source (Bukhari, Muslim, etc.). Maintain Islamic scholarly integrity."},
		{Role: providers.RoleUser, Content: fmt.Sprintf("Generate output for Islamic education tool '%s' with these parameters:\n%s\n\nPlease provide the complete, well-formatted, Islamically-grounded output directly.", toolName, string(argsJSON))},
	}

	cacheKey := fmt.Sprintf("mcp:islamic:%s:%x", toolName, hashString(string(argsJSON)))
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
			s.logger.Debug("Using cached Islamic MCP response", zap.String("tool", toolName))
			return cached, nil
		}
	}

	temp := 0.3
	maxTokens := 4096

	resp, err := provider.Chat(ctx, &providers.ChatRequest{
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	result := resp.Message.Content

	if s.cache != nil {
		_ = s.cache.Set(ctx, cacheKey, result, timeDuration(86400))
	}

	return result, nil
}

func hashString(s string) string {
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	return fmt.Sprintf("%08x", h)
}

func timeDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}
