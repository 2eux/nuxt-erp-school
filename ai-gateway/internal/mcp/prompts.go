package mcp

import (
	"context"
	"fmt"
)

func (s *MCPServer) RegisterPrompts() {
	s.registerAcademicPrompts()
	s.registerIslamicPrompts()
	s.registerFinancialPrompts()
	s.registerAdministrativePrompts()
	s.registerHRPrompts()
}

func (s *MCPServer) registerAcademicPrompts() {
	prompts := []*Prompt{
		{
			Name:        "lesson_plan_prompt",
			Description: "Generate a lesson plan for a specific subject and topic",
			Arguments: []PromptArg{
				{Name: "subject", Description: "Subject name", Required: true},
				{Name: "grade", Description: "Grade level", Required: true},
				{Name: "topic", Description: "Lesson topic", Required: true},
				{Name: "duration", Description: "Duration in minutes", Required: false},
				{Name: "curriculum", Description: "Curriculum standard", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				subject := args["subject"]
				grade := args["grade"]
				topic := args["topic"]
				duration := "60"
				if d, ok := args["duration"]; ok && d != "" {
					duration = d
				}
				curriculum := "Kurikulum Merdeka"
				if c, ok := args["curriculum"]; ok && c != "" {
					curriculum = c
				}
				return fmt.Sprintf(`Create a comprehensive lesson plan following this structure:

SUBJECT: %s
GRADE: %s
TOPIC: %s
DURATION: %s minutes
CURRICULUM: %s

Please include:
1. **Learning Objectives** (SMART - Specific, Measurable, Achievable, Relevant, Time-bound)
2. **Materials and Resources** needed
3. **Opening Activity** (5-10 minutes) - Hook/attention grabber
4. **Main Activity** (35-40 minutes) - Step-by-step instruction
5. **Closing Activity** (5-10 minutes) - Wrap-up and reflection
6. **Assessment Method** - How to measure learning
7. **Differentiation** - For different learning levels
8. **Homework/Follow-up** - Practice assignment
9. **Teacher Reflection Notes**

For Islamic schools, include relevant Islamic integration where appropriate.
Format in clear, structured markdown.`, subject, grade, topic, duration, curriculum), nil
			},
		},
		{
			Name:        "quiz_prompt",
			Description: "Generate a quiz for a subject topic",
			Arguments: []PromptArg{
				{Name: "subject", Description: "Subject", Required: true},
				{Name: "topic", Description: "Topic", Required: true},
				{Name: "grade", Description: "Grade level", Required: true},
				{Name: "num_questions", Description: "Number of questions", Required: false},
				{Name: "question_types", Description: "Question types (mc, tf, essay, short)", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				numQuestions := "10"
				if n, ok := args["num_questions"]; ok && n != "" {
					numQuestions = n
				}
				types := "multiple_choice,true_false"
				if t, ok := args["question_types"]; ok && t != "" {
					types = t
				}
				return fmt.Sprintf(`Create a quiz with the following:

Subject: %s
Topic: %s
Grade: %s
Number of questions: %s
Question types: %s

Requirements:
- Mix of easy, medium, and challenging questions
- Clear instructions for each section
- Answer key at the end
- Mark allocation per question
- Include a mix of knowledge, comprehension, and application questions
- Align with Bloom's taxonomy

Format in markdown with proper question numbering.`, args["subject"], args["topic"], args["grade"], numQuestions, types), nil
			},
		},
		{
			Name:        "exam_prompt",
			Description: "Generate a comprehensive exam paper",
			Arguments: []PromptArg{
				{Name: "subject", Description: "Subject", Required: true},
				{Name: "grade", Description: "Grade level", Required: true},
				{Name: "topics", Description: "Topics to cover (comma-separated)", Required: true},
				{Name: "exam_type", Description: "mid_term/final/unit_test", Required: false},
				{Name: "total_marks", Description: "Total marks", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				examType := "mid_term"
				if e, ok := args["exam_type"]; ok && e != "" {
					examType = e
				}
				totalMarks := "100"
				if m, ok := args["total_marks"]; ok && m != "" {
					totalMarks = m
				}
				return fmt.Sprintf(`Create a formal exam paper:

Subject: %s
Grade: %s
Topics Covered: %s
Exam Type: %s
Total Marks: %s

Include:
1. **Header**: School name placeholder, subject, grade, date, time allocation
2. **General Instructions**: Read carefully before starting
3. **Section A - Multiple Choice** (25 marks): 25 questions
4. **Section B - Short Answer** (35 marks): 5-7 questions
5. **Section C - Essay/Extended Response** (40 marks): 2-3 questions
6. **Answer Key and Marking Scheme**: For each section
7. **Bloom's Taxonomy Alignment**: Map each question to cognitive level

Format professionally in markdown.`, args["subject"], args["grade"], args["topics"], examType, totalMarks), nil
			},
		},
	}

	for _, p := range prompts {
		s.RegisterPrompt(p)
	}
}

func (s *MCPServer) registerIslamicPrompts() {
	prompts := []*Prompt{
		{
			Name:        "tahfidz_plan_prompt",
			Description: "Create a personalized Quran memorization plan",
			Arguments: []PromptArg{
				{Name: "student_name", Description: "Student name", Required: true},
				{Name: "current_memorization", Description: "Current juz/surah memorized", Required: true},
				{Name: "target", Description: "Target memorization", Required: true},
				{Name: "daily_capacity", Description: "Daily pages/verses capacity", Required: false},
				{Name: "language", Description: "Output language", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				language := "en"
				if l, ok := args["language"]; ok && l != "" {
					language = l
				}
				dailyCapacity := "1 page"
				if d, ok := args["daily_capacity"]; ok && d != "" {
					dailyCapacity = d
				}
				return fmt.Sprintf(`You are an experienced Quran teacher (Ustadz/Ustadzah). Create a detailed Tahfidz plan:

Student: %s
Current Memorization: %s
Target: %s
Daily Capacity: %s
Language: %s

Create a comprehensive tahfidz plan including:
1. Current status assessment
2. Weekly targets with daily breakdown for new memorization (ziyadah)
3. Daily muraja'ah (review) schedule with specific portions
4. Best times for memorization (after Fajr, etc.)
5. Effective memorization techniques (repetition, listening, writing)
6. Milestone markers and celebration suggestions
7. Du'a recommendations for memorization
8. Motivation with relevant hadith about Quran memorization

Format as a structured, practical daily/weekly schedule.`, args["student_name"], args["current_memorization"], args["target"], dailyCapacity, language), nil
			},
		},
		{
			Name:        "islamic_story_prompt",
			Description: "Generate an Islamic story with moral lessons",
			Arguments: []PromptArg{
				{Name: "story_type", Description: "prophets/sahabah/islamic_history/moral_story", Required: true},
				{Name: "age_group", Description: "Target age group (4-6, 7-10, 11-14, 15+)", Required: true},
				{Name: "character", Description: "Specific character if applicable", Required: false},
				{Name: "moral_lesson", Description: "Main moral to convey", Required: false},
				{Name: "language", Description: "Output language", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				moral := "honesty and integrity"
				if m, ok := args["moral_lesson"]; ok && m != "" {
					moral = m
				}
				return fmt.Sprintf(`Create an engaging Islamic story:

Story Type: %s
Age Group: %s
Character: %s
Moral Lesson: %s
Language: %s

Requirements:
- Accurate Islamic content with authentic sources
- Age-appropriate language and complexity
- Engaging narrative with dialogue
- Quran/hadith references where relevant
- Clear moral lesson at the end
- Reflection questions for children
- Related activity suggestions
- Arabic vocabulary used with translations

The story should inspire and educate while being enjoyable to read.`, args["story_type"], args["age_group"], args["character"], moral, args["language"]), nil
			},
		},
		{
			Name:        "mutabaah_analysis_prompt",
			Description: "Analyze mutaba'ah records and provide growth recommendations",
			Arguments: []PromptArg{
				{Name: "student_name", Description: "Student name", Required: true},
				{Name: "period", Description: "weekly/monthly/semester", Required: true},
				{Name: "prayer_data", Description: "Prayer consistency record", Required: false},
				{Name: "quran_data", Description: "Quran reading record", Required: false},
				{Name: "language", Description: "Output language", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				return fmt.Sprintf(`Analyze the following mutaba'ah record:

Student: %s
Period: %s

Prayer Record: %s
Quran Reading Record: %s

Provide:
1. Overall ibadah consistency summary
2. Strengths and consistent practices to maintain
3. Areas needing gentle encouragement and improvement
4. Personalized spiritual development recommendations
5. Encouraging Islamic reminders (Quran verses, hadith)
6. SMART goals for next period
7. Suggested du'a practices
8. Parent communication guidance

Keep the tone encouraging and supportive. Focus on growth, not criticism.`, args["student_name"], args["period"], args["prayer_data"], args["quran_data"]), nil
			},
		},
	}

	for _, p := range prompts {
		s.RegisterPrompt(p)
	}
}

func (s *MCPServer) registerFinancialPrompts() {
	prompts := []*Prompt{
		{
			Name:        "rkas_prompt",
			Description: "Generate RKAS (School Budget Plan) analysis and recommendations",
			Arguments: []PromptArg{
				{Name: "school_name", Description: "School name", Required: true},
				{Name: "budget_period", Description: "Budget period (e.g., 2024/2025)", Required: true},
				{Name: "total_students", Description: "Total students", Required: false},
				{Name: "bop_fund", Description: "BOS fund amount", Required: false},
				{Name: "language", Description: "Output language (id/en)", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				lang := "id"
				if l, ok := args["language"]; ok && l != "" {
					lang = l
				}
				return fmt.Sprintf(`You are an Islamic school finance expert. Create an RKAS analysis:

School: %s
Budget Period: %s
Students: %s
BOS Fund: %s
Language: %s

Generate a comprehensive RKAS with:
1. Revenue projection (BOS, SPP, donations, other sources)
2. Expenditure breakdown by 8 National Education Standards (SNP)
3. Program-based budgeting
4. Islamic education-specific allocations
5. Budget efficiency analysis
6. Risk assessment and mitigation
7. Compliance checklist for BOS reporting
8. Monitoring and evaluation framework

If Indonesian, use proper RKAS terminology (K1-K8, etc.).`, args["school_name"], args["budget_period"], args["total_students"], args["bop_fund"], lang), nil
			},
		},
	}

	for _, p := range prompts {
		s.RegisterPrompt(p)
	}
}

func (s *MCPServer) registerAdministrativePrompts() {
	prompts := []*Prompt{
		{
			Name:        "meeting_minutes_prompt",
			Description: "Generate professional meeting minutes from raw discussion notes",
			Arguments: []PromptArg{
				{Name: "meeting_title", Description: "Meeting title", Required: true},
				{Name: "date", Description: "Meeting date", Required: true},
				{Name: "attendees", Description: "Attendees list", Required: true},
				{Name: "notes", Description: "Raw discussion notes", Required: true},
				{Name: "language", Description: "Output language", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				return fmt.Sprintf(`Generate professional meeting minutes:

Meeting: %s
Date: %s
Attendees: %s
Raw Notes: %s

Create structured minutes with:
1. Header (title, date, time, venue, attendees, absentees)
2. Opening remarks
3. Agenda items reviewed with discussion summary
4. Decisions and resolutions made
5. Action items table (Task, PIC, Deadline, Status)
6. Next meeting schedule
7. Closing remarks
8. Signature section for approval

Use formal business language.`, args["meeting_title"], args["date"], args["attendees"], args["notes"]), nil
			},
		},
		{
			Name:        "official_letter_prompt",
			Description: "Generate official school letter",
			Arguments: []PromptArg{
				{Name: "letter_type", Description: "Type of official letter", Required: true},
				{Name: "subject", Description: "Letter subject", Required: true},
				{Name: "recipient", Description: "Letter recipient", Required: true},
				{Name: "content_points", Description: "Key content points", Required: true},
				{Name: "language", Description: "id/en", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				lang := args["language"]
				if lang == "id" {
					return fmt.Sprintf(`Buat surat resmi sekolah dengan format berikut:

Jenis Surat: %s
Perihal: %s
Penerima: %s
Poin Konten: %s

Format surat resmi Indonesia:
1. Kop surat (placeholder)
2. Nomor surat [placeholder]
3. Lampiran dan Perihal
4. Tanggal surat
5. Alamat tujuan
6. Salam pembuka
7. Isi surat (paragraf pembuka, isi, penutup)
8. Salam penutup
9. Tanda tangan dan nama terang
10. Tembusan (jika ada)`, args["letter_type"], args["subject"], args["recipient"], args["content_points"])
				}
				return fmt.Sprintf(`Generate official school letter in formal format:

Letter Type: %s
Subject: %s
Recipient: %s
Content Points: %s

Include:
1. Letterhead
2. Reference number [placeholder]
3. Date
4. Recipient address
5. Salutation
6. Body paragraphs (opening, content, closing)
7. Signature block
8. Enclosures notation if applicable`, args["letter_type"], args["subject"], args["recipient"], args["content_points"])
			},
		},
	}

	for _, p := range prompts {
		s.RegisterPrompt(p)
	}
}

func (s *MCPServer) registerHRPrompts() {
	prompts := []*Prompt{
		{
			Name:        "recruitment_prompt",
			Description: "Generate recruitment documents for a teaching position",
			Arguments: []PromptArg{
				{Name: "position", Description: "Position title", Required: true},
				{Name: "department", Description: "Department", Required: true},
				{Name: "requirements", Description: "Specific requirements", Required: false},
				{Name: "islamic_required", Description: "Whether Islamic qualifications are needed", Required: false},
			},
			Handler: func(ctx context.Context, args map[string]string) (string, error) {
				islamicSection := ""
				if args["islamic_required"] == "true" || args["islamic_required"] == "yes" {
					islamicSection = `
**Islamic Qualifications:**
- Hafidz Quran or minimum juz 30 memorized
- Understanding of basic fiqh and aqidah
- Can read Al-Quran with proper tajwid
- Demonstrates Islamic character and akhlaq
- Able to lead prayers and Islamic activities`
				}
				return fmt.Sprintf(`Create recruitment documents for an Islamic school teaching position:

Position: %s
Department: %s
Requirements: %s

Generate:
1. **Job Description**
   - Position summary
   - Key responsibilities
   - Working hours and schedule
   - Reporting structure
   %s

2. **Job Advertisement**
   - Engaging headline
   - About the school
   - Position overview
   - Requirements and qualifications
   - Benefits package
   - Application instructions

3. **Interview Questions**
   - Teaching philosophy (3-5 questions)
   - Technical/subject matter (3-5 questions)
   - Classroom management (3-5 questions)
   - Islamic education philosophy (2-3 questions)
   - Scenario-based questions (2-3 questions)

4. **Assessment Criteria**
   - Scoring rubric for interviews
   - Demo teaching evaluation form
   - Reference check template

Format professionally for Islamic school context.`, args["position"], args["department"], args["requirements"], islamicSection), nil
			},
		},
	}

	for _, p := range prompts {
		s.RegisterPrompt(p)
	}
}
