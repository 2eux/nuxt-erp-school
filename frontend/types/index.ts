export interface School {
  id: string
  name: string
  npsn: string
  logo: string | null
  address: string
  city: string
  province: string
  postalCode: string
  phone: string
  email: string
  website: string | null
  principal: string
  established: string
  accreditation: string
  curriculum: string
  createdAt: string
  updatedAt: string
}

export interface AcademicYear {
  id: string
  schoolId: string
  name: string
  startDate: string
  endDate: string
  isActive: boolean
  terms: Term[]
  createdAt: string
  updatedAt: string
}

export interface Term {
  id: string
  academicYearId: string
  name: string
  startDate: string
  endDate: string
  type: 'semester' | 'trimester' | 'quarter'
  isActive: boolean
}

export interface Class {
  id: string
  schoolId: string
  academicYearId: string
  name: string
  grade: number
  program: string
  capacity: number
  homeRoomTeacherId: string | null
  homeRoomTeacherName: string | null
  studentCount: number
  room: string | null
  createdAt: string
  updatedAt: string
}

export interface Subject {
  id: string
  schoolId: string
  name: string
  code: string
  description: string | null
  category: 'academic' | 'islamic' | 'language' | 'extracurricular' | 'other'
  creditHours: number
  createdAt: string
  updatedAt: string
}

export interface Curriculum {
  id: string
  schoolId: string
  academicYearId: string
  subjectId: string
  classId: string
  teacherId: string | null
  teacherName: string | null
  scheduleDay: string
  scheduleStart: string
  scheduleEnd: string
  room: string | null
  createdAt: string
  updatedAt: string
}

export interface Schedule {
  id: string
  schoolId: string
  academicYearId: string
  classId: string
  className: string
  subjectId: string
  subjectName: string
  teacherId: string | null
  teacherName: string | null
  day: 'monday' | 'tuesday' | 'wednesday' | 'thursday' | 'friday' | 'saturday' | 'sunday'
  startTime: string
  endTime: string
  room: string | null
  semester: number
  createdAt: string
  updatedAt: string
}

export interface Attendance {
  id: string
  schoolId: string
  studentId: string
  studentName: string
  classId: string
  className: string
  date: string
  status: 'present' | 'absent' | 'late' | 'sick' | 'permission' | 'holiday'
  checkInTime: string | null
  checkOutTime: string | null
  note: string | null
  createdAt: string
  updatedAt: string
}

export interface Grade {
  id: string
  schoolId: string
  studentId: string
  studentName: string
  subjectId: string
  subjectName: string
  classId: string
  className: string
  academicYearId: string
  termId: string
  score: number
  type: 'daily' | 'midterm' | 'final' | 'assignment' | 'project'
  description: string | null
  gradeLetter: string
  createdAt: string
  updatedAt: string
}

export interface Exam {
  id: string
  schoolId: string
  academicYearId: string
  termId: string
  name: string
  type: 'midterm' | 'final' | 'daily' | 'semester'
  subjectId: string
  subjectName: string
  classId: string
  className: string
  date: string
  startTime: string
  endTime: string
  maxScore: number
  passingScore: number
  createdAt: string
  updatedAt: string
}

export interface ReportCard {
  id: string
  schoolId: string
  studentId: string
  studentName: string
  classId: string
  className: string
  academicYearId: string
  termId: string
  grades: {
    subjectId: string
    subjectName: string
    score: number
    gradeLetter: string
    description: string
  }[]
  attendanceSummary: {
    present: number
    absent: number
    late: number
    sick: number
    permission: number
  }
  rank: number | null
  gpa: number
  teacherNote: string | null
  principalNote: string | null
  createdAt: string
  updatedAt: string
}

export interface QuranProgress {
  id: string
  schoolId: string
  studentId: string
  studentName: string
  surahStart: string
  ayahStart: number
  surahEnd: string
  ayahEnd: number
  memorizationType: 'tahfidz' | 'murojaah' | 'tilawah'
  status: 'memorized' | 'in_progress' | 'not_memorized'
  teacherId: string | null
  teacherName: string | null
  date: string
  note: string | null
  createdAt: string
  updatedAt: string
}

export interface Mutabaah {
  id: string
  schoolId: string
  studentId: string
  studentName: string
  date: string
  fajr: boolean
  dhuhr: boolean
  asr: boolean
  maghrib: boolean
  isha: boolean
  tahajjud: boolean
  dhuha: boolean
  quranPages: number
  note: string | null
  createdAt: string
  updatedAt: string
}

export interface Halaqah {
  id: string
  schoolId: string
  name: string
  teacherId: string
  teacherName: string
  description: string | null
  schedule: string
  room: string | null
  studentIds: string[]
  createdAt: string
  updatedAt: string
}

export interface IslamicEvent {
  id: string
  schoolId: string
  name: string
  date: string
  type: 'hijri' | 'islamic_holiday' | 'school_islamic'
  description: string | null
  createdAt: string
  updatedAt: string
}

export interface User {
  id: string
  schoolId: string
  email: string
  phone: string | null
  fullName: string
  avatar: string | null
  roleId: string
  roleName: string
  isActive: boolean
  lastLoginAt: string | null
  createdAt: string
  updatedAt: string
}

export interface Role {
  id: string
  schoolId: string
  name: string
  description: string | null
  permissions: string[]
  isSystem: boolean
  createdAt: string
  updatedAt: string
}

export interface Student {
  id: string
  schoolId: string
  nis: string
  nisn: string
  fullName: string
  nickName: string
  gender: 'male' | 'female'
  birthPlace: string
  birthDate: string
  religion: string
  address: string
  city: string
  postalCode: string
  phone: string | null
  email: string | null
  fatherName: string
  fatherPhone: string
  fatherOccupation: string | null
  motherName: string
  motherPhone: string
  motherOccupation: string | null
  guardianName: string | null
  guardianPhone: string | null
  classId: string | null
  className: string | null
  enrollmentDate: string
  status: 'active' | 'inactive' | 'graduated' | 'transferred' | 'dropped'
  photo: string | null
  createdAt: string
  updatedAt: string
}

export interface Teacher {
  id: string
  schoolId: string
  nip: string
  userId: string | null
  fullName: string
  gender: 'male' | 'female'
  birthPlace: string
  birthDate: string
  address: string
  phone: string
  email: string
  education: string
  major: string
  joinDate: string
  status: 'active' | 'inactive' | 'resigned'
  subjects: string[]
  photo: string | null
  createdAt: string
  updatedAt: string
}

export interface Employee {
  id: string
  schoolId: string
  nip: string
  userId: string | null
  fullName: string
  gender: 'male' | 'female'
  position: string
  department: string
  phone: string
  email: string
  joinDate: string
  status: 'active' | 'inactive' | 'resigned'
  photo: string | null
  createdAt: string
  updatedAt: string
}

export interface Invoice {
  id: string
  schoolId: string
  invoiceNumber: string
  studentId: string
  studentName: string
  classId: string
  className: string
  type: 'spp' | 'registration' | 'exam' | 'book' | 'uniform' | 'activity' | 'other'
  amount: number
  dueDate: string
  paidAmount: number
  status: 'unpaid' | 'partial' | 'paid' | 'overdue' | 'cancelled'
  paymentDate: string | null
  academicYearId: string
  month: number
  year: number
  description: string | null
  createdAt: string
  updatedAt: string
}

export interface Payment {
  id: string
  schoolId: string
  invoiceId: string
  invoiceNumber: string
  studentId: string
  studentName: string
  amount: number
  paymentMethod: 'cash' | 'bank_transfer' | 'digital_wallet' | 'other'
  paymentDate: string
  referenceNumber: string | null
  note: string | null
  recordedById: string
  recordedByName: string
  createdAt: string
  updatedAt: string
}

export interface Journal {
  id: string
  schoolId: string
  journalNumber: string
  date: string
  type: 'income' | 'expense' | 'transfer'
  category: string
  amount: number
  description: string
  reference: string | null
  status: 'draft' | 'posted'
  createdAt: string
  updatedAt: string
}

export interface Budget {
  id: string
  schoolId: string
  academicYearId: string
  name: string
  category: string
  plannedAmount: number
  actualAmount: number
  startDate: string
  endDate: string
  status: 'draft' | 'approved' | 'active' | 'closed'
  description: string | null
  createdAt: string
  updatedAt: string
}

export interface Payroll {
  id: string
  schoolId: string
  employeeId: string
  employeeName: string
  month: number
  year: number
  baseSalary: number
  allowances: number
  deductions: number
  netSalary: number
  status: 'draft' | 'approved' | 'paid'
  paymentDate: string | null
  createdAt: string
  updatedAt: string
}

export interface Inventory {
  id: string
  schoolId: string
  code: string
  name: string
  category: string
  description: string | null
  quantity: number
  unit: string
  location: string | null
  condition: 'excellent' | 'good' | 'fair' | 'damaged' | 'lost'
  purchaseDate: string | null
  purchasePrice: number | null
  createdBy: string
  createdAt: string
  updatedAt: string
}

export interface Asset {
  id: string
  schoolId: string
  code: string
  name: string
  category: string
  description: string | null
  acquisitionDate: string
  acquisitionPrice: number
  currentValue: number
  depreciationRate: number
  location: string | null
  condition: 'excellent' | 'good' | 'fair' | 'damaged' | 'lost'
  status: 'active' | 'disposed' | 'maintenance'
  createdAt: string
  updatedAt: string
}

export interface Book {
  id: string
  schoolId: string
  isbn: string | null
  title: string
  author: string
  publisher: string
  category: string
  publishYear: number
  quantity: number
  available: number
  rack: string | null
  cover: string | null
  createdAt: string
  updatedAt: string
}

export interface Borrowing {
  id: string
  schoolId: string
  bookId: string
  bookTitle: string
  borrowerId: string
  borrowerName: string
  borrowerType: 'student' | 'teacher' | 'employee'
  borrowDate: string
  dueDate: string
  returnDate: string | null
  status: 'borrowed' | 'returned' | 'overdue' | 'lost'
  fine: number | null
  createdAt: string
  updatedAt: string
}

export interface MedicalRecord {
  id: string
  schoolId: string
  studentId: string
  studentName: string
  date: string
  complaint: string
  diagnosis: string | null
  treatment: string | null
  medication: string | null
  referredTo: string | null
  recordedById: string
  recordedByName: string
  createdAt: string
  updatedAt: string
}

export interface CounselingRecord {
  id: string
  schoolId: string
  studentId: string
  studentName: string
  counselorId: string
  counselorName: string
  date: string
  type: 'individual' | 'group' | 'parent'
  category: 'academic' | 'behavior' | 'social' | 'emotional' | 'career' | 'other'
  description: string
  followUp: string | null
  isConfidential: boolean
  createdAt: string
  updatedAt: string
}

export interface Admission {
  id: string
  schoolId: string
  registrationNumber: string
  fullName: string
  gender: 'male' | 'female'
  birthPlace: string
  birthDate: string
  previousSchool: string | null
  parentName: string
  parentPhone: string
  parentEmail: string | null
  appliedGrade: number
  registrationDate: string
  status: 'pending' | 'documents_submitted' | 'test_scheduled' | 'tested' | 'accepted' | 'rejected' | 'enrolled'
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface Announcement {
  id: string
  schoolId: string
  title: string
  content: string
  type: 'general' | 'academic' | 'event' | 'urgent' | 'finance'
  targetAudience: ('all' | 'teachers' | 'students' | 'parents' | 'employees')[]
  publishDate: string
  expiryDate: string | null
  isActive: boolean
  createdById: string
  createdByName: string
  createdAt: string
  updatedAt: string
}

export interface Message {
  id: string
  schoolId: string
  senderId: string
  senderName: string
  recipientIds: string[]
  recipientType: 'individual' | 'class' | 'role' | 'all'
  subject: string
  content: string
  type: 'message' | 'notification' | 'alert'
  isRead: boolean[]
  attachments: string[]
  createdAt: string
  updatedAt: string
}

export interface Meeting {
  id: string
  schoolId: string
  title: string
  description: string | null
  date: string
  startTime: string
  endTime: string
  location: string | null
  type: 'staff' | 'parent' | 'committee' | 'other'
  organizerId: string
  organizerName: string
  participantIds: string[]
  status: 'scheduled' | 'ongoing' | 'completed' | 'cancelled'
  minutes: string | null
  createdAt: string
  updatedAt: string
}

export interface Document {
  id: string
  schoolId: string
  title: string
  type: 'certificate' | 'letter' | 'form' | 'report' | 'other'
  template: string | null
  content: string
  variables: Record<string, string>
  generatedUrl: string | null
  createdById: string
  createdByName: string
  createdAt: string
  updatedAt: string
}

export interface Notification {
  id: string
  schoolId: string
  userId: string
  title: string
  message: string
  type: 'info' | 'success' | 'warning' | 'error'
  category: string
  referenceId: string | null
  referenceType: string | null
  isRead: boolean
  createdAt: string
  updatedAt: string
}

export interface AuditLog {
  id: string
  schoolId: string
  userId: string
  userName: string
  action: string
  entity: string
  entityId: string
  oldData: Record<string, unknown> | null
  newData: Record<string, unknown> | null
  ipAddress: string | null
  userAgent: string | null
  createdAt: string
}

export type Gender = 'male' | 'female'
export type DayOfWeek = 'monday' | 'tuesday' | 'wednesday' | 'thursday' | 'friday' | 'saturday' | 'sunday'
export type AttendanceStatus = 'present' | 'absent' | 'late' | 'sick' | 'permission' | 'holiday'
export type StudentStatus = 'active' | 'inactive' | 'graduated' | 'transferred' | 'dropped'
export type StaffStatus = 'active' | 'inactive' | 'resigned'
export type InvoiceStatus = 'unpaid' | 'partial' | 'paid' | 'overdue' | 'cancelled'
export type PaymentMethod = 'cash' | 'bank_transfer' | 'digital_wallet' | 'other'
export type AdmissionStatus = 'pending' | 'documents_submitted' | 'test_scheduled' | 'tested' | 'accepted' | 'rejected' | 'enrolled'
export type Condition = 'excellent' | 'good' | 'fair' | 'damaged' | 'lost'
export type Permission = string

export interface AuthTokens {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

export interface AuthUser {
  id: string
  fullName: string
  email: string
  phone: string | null
  avatar: string | null
  schoolId: string
  schoolName: string
  roleId: string
  roleName: string
  permissions: string[]
  isSuperAdmin: boolean
}

export interface LoginRequest {
  email: string
  password: string
  rememberMe?: boolean
}

export interface LoginResponse {
  user: AuthUser
  tokens: AuthTokens
}

export interface ForgotPasswordRequest {
  email: string
}

export interface ResetPasswordRequest {
  token: string
  password: string
  passwordConfirmation: string
}

export interface ChangePasswordRequest {
  currentPassword: string
  newPassword: string
}

export interface PaginationParams {
  page?: number
  limit?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface FilterParams {
  search?: string
  startDate?: string
  endDate?: string
  status?: string
  [key: string]: string | number | boolean | undefined
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
    hasNextPage: boolean
    hasPreviousPage: boolean
  }
}

export interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
  error?: string
}

export interface DashboardStats {
  totalStudents: number
  totalTeachers: number
  totalEmployees: number
  totalClasses: number
  totalRevenue: number
  attendanceRate: number
  studentGrowth: number
  revenueGrowth: number
}

export interface ChartData {
  labels: string[]
  series: {
    name: string
    data: number[]
  }[]
}

export interface RecentActivity {
  id: string
  type: string
  title: string
  description: string
  user: string
  timestamp: string
}

export interface FormSelectOption {
  label: string
  value: string | number
  disabled?: boolean
}

export interface MenuItem {
  label: string
  icon: string
  to?: string
  children?: MenuItem[]
  permission?: string
  badge?: number
  exact?: boolean
}

export interface TableColumn {
  key: string
  label: string
  sortable?: boolean
  filterable?: boolean
  type?: 'text' | 'number' | 'date' | 'status' | 'currency' | 'image' | 'action'
  width?: number
  align?: 'left' | 'center' | 'right'
  formatter?: (value: unknown, row: Record<string, unknown>) => string
}

export interface SidebarState {
  isCollapsed: boolean
  isMobileOpen: boolean
  activeGroup: string | null
}

export interface ThemePreference {
  mode: 'light' | 'dark' | 'system'
  primaryColor: string
  accentColor: string
  borderRadius: 'none' | 'sm' | 'md' | 'lg'
}

export interface NotificationPreference {
  email: boolean
  push: boolean
  sms: boolean
  categories: string[]
}
