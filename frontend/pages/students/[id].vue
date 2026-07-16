<template>
  <div class="space-y-6">
    <div class="flex items-center gap-4">
      <UButton color="gray" variant="ghost" size="sm" icon="i-heroicons-arrow-left" :to="'/students'" />
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ student?.fullName || $t('common.loading') }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400">NIS: {{ student?.nis }} | NISN: {{ student?.nisn }}</p>
      </div>
    </div>

    <template v-if="loading">
      <LoadingSkeleton type="detail" />
    </template>

    <template v-else-if="!student">
      <EmptyState
        :title="$t('students.not_found')"
        :description="$t('students.not_found_description')"
        icon="i-heroicons-user-circle"
        :action-label="$t('common.back_to_list')"
        @action="navigateTo('/students')"
      />
    </template>

    <template v-else>
      <UTabs :items="tabs" :ui="{ list: { tab: { size: 'sm' } } }">
        <template #item="{ item }">
          <div class="pt-4 space-y-6">
            <template v-if="item.key === 'profile'">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('students.personal_info') }}</h3>
                  <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="openEditProfile" />
                </div>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div class="flex flex-col items-center gap-3">
                    <UAvatar
                      :src="(student.photo as string) || undefined"
                      :alt="student.fullName"
                      size="3xl"
                    />
                    <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ student.fullName }}</h2>
                    <StatusBadge :status="student.status" />
                  </div>
                  <div class="md:col-span-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div><span class="text-xs text-gray-500">{{ $t('students.nis') }}</span><p class="text-sm font-medium text-gray-900 dark:text-white">{{ student.nis }}</p></div>
                    <div><span class="text-xs text-gray-500">{{ $t('students.nisn') }}</span><p class="text-sm font-medium text-gray-900 dark:text-white">{{ student.nisn }}</p></div>
                    <div><span class="text-xs text-gray-500">{{ $t('common.gender') }}</span><p class="text-sm font-medium text-gray-900 dark:text-white">{{ $t(`common.${student.gender}`) }}</p></div>
                    <div><span class="text-xs text-gray-500">{{ $t('students.birth_place') }}, {{ $t('students.birth_date') }}</span><p class="text-sm font-medium text-gray-900 dark:text-white">{{ student.birthPlace }}, {{ $dayjs(student.birthDate).format('DD MMMM YYYY') }}</p></div>
                    <div><span class="text-xs text-gray-500">{{ $t('students.religion') }}</span><p class="text-sm font-medium text-gray-900 dark:text-white">{{ student.religion }}</p></div>
                    <div><span class="text-xs text-gray-500">{{ $t('students.class') }}</span><p class="text-sm font-medium text-gray-900 dark:text-white">{{ student.className || '-' }}</p></div>
                    <div class="sm:col-span-2"><span class="text-xs text-gray-500">{{ $t('students.address') }}</span><p class="text-sm font-medium text-gray-900 dark:text-white">{{ student.address }}, {{ student.city }}</p></div>
                  </div>
                </div>
              </div>

              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('students.parent_info') }}</h3>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <h4 class="text-sm font-semibold text-gray-900 dark:text-white mb-2">{{ $t('students.father') }}</h4>
                    <div class="space-y-2">
                      <div><span class="text-xs text-gray-500">{{ $t('common.name') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ student.fatherName }}</p></div>
                      <div><span class="text-xs text-gray-500">{{ $t('common.phone') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ student.fatherPhone }}</p></div>
                      <div><span class="text-xs text-gray-500">{{ $t('common.occupation') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ student.fatherOccupation || '-' }}</p></div>
                    </div>
                  </div>
                  <div>
                    <h4 class="text-sm font-semibold text-gray-900 dark:text-white mb-2">{{ $t('students.mother') }}</h4>
                    <div class="space-y-2">
                      <div><span class="text-xs text-gray-500">{{ $t('common.name') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ student.motherName }}</p></div>
                      <div><span class="text-xs text-gray-500">{{ $t('common.phone') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ student.motherPhone }}</p></div>
                      <div><span class="text-xs text-gray-500">{{ $t('common.occupation') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ student.motherOccupation || '-' }}</p></div>
                    </div>
                  </div>
                </div>
              </div>
            </template>

            <template v-if="item.key === 'academic'">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('students.class_history') }}</h3>
                </div>
                <DataTable
                  :columns="classHistoryColumns"
                  :rows="classHistory"
                  :loading="subLoading"
                  :empty-title="$t('students.no_class_history')"
                  :show-filters="false"
                />
              </div>
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('students.grades') }}</h3>
                </div>
                <DataTable
                  :columns="gradeColumns"
                  :rows="grades"
                  :loading="subLoading"
                  :empty-title="$t('students.no_grades')"
                  :show-filters="false"
                />
              </div>
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('students.report_cards') }}</h3>
                </div>
                <DataTable
                  :columns="reportCardColumns"
                  :rows="reportCards"
                  :loading="subLoading"
                  :empty-title="$t('students.no_report_cards')"
                  :show-filters="false"
                >
                  <template #item-actions="{ row }">
                    <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-printer" @click="viewReportCard(row as Record<string, unknown>)" />
                  </template>
                </DataTable>
              </div>
            </template>

            <template v-if="item.key === 'attendance'">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('attendance.title') }}</h3>
                  <div class="flex items-center gap-2">
                    <UInput v-model="attendanceMonth" type="month" color="gray" size="sm" @change="fetchAttendance" />
                  </div>
                </div>
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-4">
                  <StatCard :label="$t('attendance.present')" :value="attendanceStats.present" icon="i-heroicons-check-circle" color="emerald" :loading="subLoading" />
                  <StatCard :label="$t('attendance.absent')" :value="attendanceStats.absent" icon="i-heroicons-x-circle" color="red" :loading="subLoading" />
                  <StatCard :label="$t('attendance.late')" :value="attendanceStats.late" icon="i-heroicons-clock" color="amber" :loading="subLoading" />
                  <StatCard :label="$t('attendance.sick')" :value="attendanceStats.sick" icon="i-heroicons-heart" color="purple" :loading="subLoading" />
                </div>
                <DataTable
                  :columns="attendanceColumns"
                  :rows="attendanceRecords"
                  :loading="subLoading"
                  :empty-title="$t('attendance.no_records')"
                  :show-filters="false"
                />
              </div>
            </template>

            <template v-if="item.key === 'tahfidz'">
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div class="card">
                  <div class="card-header">
                    <h3 class="card-title">{{ $t('tahfidz.progress') }}</h3>
                  </div>
                  <div class="h-80">
                    <ApexChart type="bar" height="100%" :options="tahfidzChart.options" :series="tahfidzChart.series" />
                  </div>
                </div>
                <div class="card">
                  <div class="card-header">
                    <h3 class="card-title">{{ $t('tahfidz.juz_completion') }}</h3>
                  </div>
                  <div class="space-y-3">
                    <div v-for="juz in juzProgress" :key="juz.number" class="space-y-1">
                      <div class="flex items-center justify-between">
                        <span class="text-sm text-gray-700 dark:text-gray-300">{{ $t('tahfidz.juz') }} {{ juz.number }}</span>
                        <span class="text-sm font-semibold text-brand-600 dark:text-brand-400">{{ juz.progress }}%</span>
                      </div>
                      <div class="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                        <div class="h-full bg-brand-500 rounded-full transition-all" :style="{ width: `${juz.progress}%` }" />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <DataTable
                :columns="tahfidzColumns"
                :rows="tahfidzRecords"
                :loading="subLoading"
                :empty-title="$t('tahfidz.no_records')"
                :show-filters="false"
              />
            </template>

            <template v-if="item.key === 'mutabaah'">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('mutabaah.daily_log') }}</h3>
                  <UInput v-model="mutabaahMonth" type="month" color="gray" size="sm" @change="fetchMutabaah" />
                </div>
                <DataTable
                  :columns="mutabaahColumns"
                  :rows="mutabaahRecords"
                  :loading="subLoading"
                  :empty-title="$t('mutabaah.no_records')"
                  :show-filters="false"
                />
              </div>
            </template>

            <template v-if="item.key === 'finance'">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('finance.invoices') }}</h3>
                </div>
                <DataTable
                  :columns="financeInvoiceColumns"
                  :rows="studentInvoices"
                  :loading="subLoading"
                  :empty-title="$t('finance.no_invoices')"
                  :show-filters="false"
                >
                  <template #cell-status="{ row }">
                    <StatusBadge :status="row.status as string" />
                  </template>
                  <template #cell-amount="{ row }">
                    <span class="text-sm font-mono text-gray-900 dark:text-white">
                      {{ formatCurrency(row.amount as number) }}
                    </span>
                  </template>
                </DataTable>
              </div>
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('finance.payments') }}</h3>
                </div>
                <DataTable
                  :columns="financePaymentColumns"
                  :rows="studentPayments"
                  :loading="subLoading"
                  :empty-title="$t('finance.no_payments')"
                  :show-filters="false"
                >
                  <template #cell-amount="{ row }">
                    <span class="text-sm font-mono text-gray-900 dark:text-white">
                      {{ formatCurrency(row.amount as number) }}
                    </span>
                  </template>
                </DataTable>
              </div>
            </template>

            <template v-if="item.key === 'medical'">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('medical.records') }}</h3>
                  <UButton v-if="permissions.can('medical.create')" color="primary" size="xs" icon="i-heroicons-plus" @click="openAddMedical" />
                </div>
                <DataTable
                  :columns="medicalColumns"
                  :rows="medicalRecords"
                  :loading="subLoading"
                  :empty-title="$t('medical.no_records')"
                  :show-filters="false"
                >
                  <template #cell-date="{ row }">
                    <span class="text-sm text-gray-700 dark:text-gray-300">{{ $dayjs(row.date as string).format('DD MMM YYYY') }}</span>
                  </template>
                </DataTable>
              </div>
            </template>

            <template v-if="item.key === 'documents'">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('documents.title') }}</h3>
                  <UButton color="primary" size="xs" icon="i-heroicons-plus" @click="triggerUpload" />
                </div>
                <FileUpload
                  ref="docUploadRef"
                  :accept="'.pdf,.jpg,.jpeg,.png,.doc,.docx'"
                  :multiple="true"
                  :accept-hint="$t('documents.upload_hint')"
                  @files-selected="handleFileUpload"
                />
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 mt-4">
                  <div
                    v-for="doc in studentDocuments"
                    :key="doc.id"
                    class="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 cursor-pointer"
                  >
                    <UIcon name="i-heroicons-document-text" class="w-8 h-8 text-gray-400" />
                    <div class="flex-1 min-w-0">
                      <p class="text-sm font-medium text-gray-900 dark:text-white truncate">{{ doc.title }}</p>
                      <p class="text-xs text-gray-500">{{ $dayjs(doc.createdAt).format('DD MMM YYYY') }}</p>
                    </div>
                    <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-arrow-down-tray" @click="downloadDocument(doc)" />
                  </div>
                </div>
                <EmptyState
                  v-if="!subLoading && studentDocuments.length === 0"
                  :title="$t('documents.no_documents')"
                  icon="i-heroicons-document"
                />
              </div>
            </template>

            <template v-if="item.key === 'behavior'">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('students.behavior_notes') }}</h3>
                  <UButton color="primary" size="xs" icon="i-heroicons-plus" @click="openAddBehavior" />
                </div>
                <div class="space-y-3">
                  <div
                    v-for="note in behaviorNotes"
                    :key="note.id"
                    class="flex items-start gap-3 p-3 rounded-lg border"
                    :class="note.type === 'positive' ? 'border-emerald-200 dark:border-emerald-800 bg-emerald-50 dark:bg-emerald-900/10' : 'border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/10'"
                  >
                    <UIcon
                      :name="note.type === 'positive' ? 'i-heroicons-hand-thumb-up' : 'i-heroicons-exclamation-triangle'"
                      class="w-5 h-5 mt-0.5"
                      :class="note.type === 'positive' ? 'text-emerald-600' : 'text-red-600'"
                    />
                    <div class="flex-1">
                      <p class="text-sm text-gray-900 dark:text-white">{{ note.description }}</p>
                      <p class="text-xs text-gray-500 mt-1">{{ $dayjs(note.date).format('DD MMM YYYY') }} - {{ note.recordedBy }}</p>
                    </div>
                    <span class="text-xs font-medium px-2 py-0.5 rounded-full" :class="note.type === 'positive' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'">
                      {{ $t(`students.${note.type}`) }}
                    </span>
                  </div>
                </div>
                <EmptyState
                  v-if="!subLoading && behaviorNotes.length === 0"
                  :title="$t('students.no_behavior_notes')"
                  icon="i-heroicons-clipboard-document"
                />
              </div>
            </template>
          </div>
        </template>
      </UTabs>
    </template>

    <FormDialog
      v-model="showProfileForm"
      :title="$t('students.edit_profile')"
      :loading="saving"
      @submit="saveProfile"
      @cancel="showProfileForm = false"
    >
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('students.full_name')" class="sm:col-span-2">
          <UInput v-model="editForm.fullName" color="gray" />
        </UFormGroup>
        <UFormGroup :label="$t('students.birth_place')">
          <UInput v-model="editForm.birthPlace" color="gray" />
        </UFormGroup>
        <UFormGroup :label="$t('students.birth_date')">
          <UInput v-model="editForm.birthDate" type="date" color="gray" />
        </UFormGroup>
        <UFormGroup :label="$t('students.address')" class="sm:col-span-2">
          <UTextarea v-model="editForm.address" color="gray" :rows="2" />
        </UFormGroup>
        <UFormGroup :label="$t('students.phone')">
          <UInput v-model="editForm.phone" color="gray" />
        </UFormGroup>
        <UFormGroup :label="$t('students.status')">
          <USelect v-model="editForm.status" :options="statusOptions" color="gray" />
        </UFormGroup>
      </div>
    </FormDialog>

    <FormDialog
      v-model="showMedicalForm"
      :title="$t('medical.add_record')"
      :loading="savingMedical"
      @submit="saveMedical"
      @cancel="showMedicalForm = false"
    >
      <div class="space-y-4">
        <UFormGroup :label="$t('medical.date')" required>
          <UInput v-model="medicalForm.date" type="date" color="gray" />
        </UFormGroup>
        <UFormGroup :label="$t('medical.complaint')" required>
          <UTextarea v-model="medicalForm.complaint" color="gray" :rows="2" />
        </UFormGroup>
        <UFormGroup :label="$t('medical.diagnosis')">
          <UTextarea v-model="medicalForm.diagnosis" color="gray" :rows="2" />
        </UFormGroup>
        <UFormGroup :label="$t('medical.treatment')">
          <UTextarea v-model="medicalForm.treatment" color="gray" :rows="2" />
        </UFormGroup>
        <UFormGroup :label="$t('medical.medication')">
          <UInput v-model="medicalForm.medication" color="gray" />
        </UFormGroup>
      </div>
    </FormDialog>

    <FormDialog
      v-model="showBehaviorForm"
      :title="$t('students.add_behavior_note')"
      :loading="savingBehavior"
      @submit="saveBehavior"
      @cancel="showBehaviorForm = false"
    >
      <div class="space-y-4">
        <UFormGroup :label="$t('students.type')" required>
          <USelect v-model="behaviorForm.type" :options="[{ label: $t('students.positive'), value: 'positive' }, { label: $t('students.negative'), value: 'negative' }]" color="gray" />
        </UFormGroup>
        <UFormGroup :label="$t('common.description')" required>
          <UTextarea v-model="behaviorForm.description" color="gray" :rows="3" />
        </UFormGroup>
        <UFormGroup :label="$t('common.date')">
          <UInput v-model="behaviorForm.date" type="date" color="gray" />
        </UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn, Student } from '~/types'

definePageMeta({
  middleware: ['auth'],
})

const route = useRoute()
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const { $dayjs } = useNuxtApp()
const toast = useToast()

const studentId = computed(() => route.params.id as string)

const loading = ref(true)
const subLoading = ref(false)
const saving = ref(false)
const savingMedical = ref(false)
const savingBehavior = ref(false)

const student = ref<Student | null>(null)

const classHistory = ref<Record<string, unknown>[]>([])
const grades = ref<Record<string, unknown>[]>([])
const reportCards = ref<Record<string, unknown>[]>([])
const attendanceRecords = ref<Record<string, unknown>[]>([])
const tahfidzRecords = ref<Record<string, unknown>[]>([])
const mutabaahRecords = ref<Record<string, unknown>[]>([])
const studentInvoices = ref<Record<string, unknown>[]>([])
const studentPayments = ref<Record<string, unknown>[]>([])
const medicalRecords = ref<Record<string, unknown>[]>([])
const studentDocuments = ref<{ id: string; title: string; url: string; createdAt: string }[]>([])
const behaviorNotes = ref<{ id: string; type: string; description: string; date: string; recordedBy: string }[]>([])

const attendanceMonth = ref($dayjs().format('YYYY-MM'))
const mutabaahMonth = ref($dayjs().format('YYYY-MM'))
const docUploadRef = ref()

const showProfileForm = ref(false)
const showMedicalForm = ref(false)
const showBehaviorForm = ref(false)

const attendanceStats = reactive({ present: 0, absent: 0, late: 0, sick: 0 })

const juzProgress = Array.from({ length: 30 }, (_, i) => ({
  number: 30 - i,
  progress: Math.floor(Math.random() * 100),
}))

const editForm = reactive({
  fullName: '',
  birthPlace: '',
  birthDate: '',
  address: '',
  phone: '',
  status: '',
})

const medicalForm = reactive({
  date: $dayjs().format('YYYY-MM-DD'),
  complaint: '',
  diagnosis: '',
  treatment: '',
  medication: '',
})

const behaviorForm = reactive({
  type: 'positive',
  description: '',
  date: $dayjs().format('YYYY-MM-DD'),
})

const statusOptions = [
  { label: t('status.active'), value: 'active' },
  { label: t('status.inactive'), value: 'inactive' },
  { label: t('status.graduated'), value: 'graduated' },
  { label: t('status.transferred'), value: 'transferred' },
  { label: t('status.dropped'), value: 'dropped' },
]

const tabs = computed(() => [
  { key: 'profile', label: t('students.profile') },
  { key: 'academic', label: t('students.academic') },
  { key: 'attendance', label: t('attendance.title') },
  { key: 'tahfidz', label: t('tahfidz.title') },
  { key: 'mutabaah', label: t('mutabaah.title') },
  { key: 'finance', label: t('finance.title') },
  { key: 'medical', label: t('medical.title') },
  { key: 'documents', label: t('documents.title') },
  { key: 'behavior', label: t('students.behavior') },
])

const classHistoryColumns: TableColumn[] = [
  { key: 'academicYear', label: 'academic_year.title' },
  { key: 'className', label: 'students.class' },
  { key: 'grade', label: 'students.grade' },
  { key: 'status', label: 'common.status' },
]

const gradeColumns: TableColumn[] = [
  { key: 'subjectName', label: 'subjects.title' },
  { key: 'type', label: 'grades.type' },
  { key: 'score', label: 'grades.score', type: 'number' },
  { key: 'gradeLetter', label: 'grades.grade' },
  { key: 'term', label: 'academic_year.term' },
]

const reportCardColumns: TableColumn[] = [
  { key: 'term', label: 'academic_year.term' },
  { key: 'academicYear', label: 'academic_year.title' },
  { key: 'gpa', label: 'students.gpa', type: 'number' },
  { key: 'rank', label: 'students.rank', type: 'number' },
]

const attendanceColumns: TableColumn[] = [
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'status', label: 'common.status', type: 'status' },
  { key: 'checkInTime', label: 'attendance.check_in' },
  { key: 'checkOutTime', label: 'attendance.check_out' },
  { key: 'note', label: 'common.note' },
]

const tahfidzColumns: TableColumn[] = [
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'surah', label: 'quran.surah' },
  { key: 'ayahRange', label: 'quran.ayah' },
  { key: 'type', label: 'tahfidz.type' },
  { key: 'status', label: 'common.status', type: 'status' },
]

const mutabaahColumns: TableColumn[] = [
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'fajr', label: 'prayer.fajr' },
  { key: 'dhuhr', label: 'prayer.dhuhr' },
  { key: 'asr', label: 'prayer.asr' },
  { key: 'maghrib', label: 'prayer.maghrib' },
  { key: 'isha', label: 'prayer.isha' },
  { key: 'quranPages', label: 'quran.pages' },
]

const financeInvoiceColumns: TableColumn[] = [
  { key: 'invoiceNumber', label: 'finance.invoice_number' },
  { key: 'type', label: 'finance.type' },
  { key: 'amount', label: 'finance.amount', type: 'currency' },
  { key: 'paidAmount', label: 'finance.paid', type: 'currency' },
  { key: 'dueDate', label: 'finance.due_date', type: 'date' },
  { key: 'status', label: 'common.status', type: 'status' },
]

const financePaymentColumns: TableColumn[] = [
  { key: 'invoiceNumber', label: 'finance.invoice_number' },
  { key: 'amount', label: 'finance.amount', type: 'currency' },
  { key: 'paymentMethod', label: 'finance.payment_method' },
  { key: 'paymentDate', label: 'finance.payment_date', type: 'date' },
]

const medicalColumns: TableColumn[] = [
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'complaint', label: 'medical.complaint' },
  { key: 'diagnosis', label: 'medical.diagnosis' },
  { key: 'treatment', label: 'medical.treatment' },
]

const tahfidzChart = computed(() => ({
  series: [{ name: t('tahfidz.ayahs_memorized'), data: [15, 22, 18, 35, 28, 42, 38, 50, 45, 55, 48, 60] }],
  options: {
    chart: { type: 'bar' as const, toolbar: { show: false }, background: 'transparent' },
    colors: ['#059669'],
    plotOptions: { bar: { borderRadius: 4, columnWidth: '60%' } },
    xaxis: { categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'] },
    yaxis: { title: { text: t('quran.ayahs') } },
    grid: { borderColor: '#e5e7eb', strokeDashArray: 4 },
  },
}))

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value)
}

const fetchStudent = async () => {
  loading.value = true
  try {
    student.value = await api.get<Student>(`/students/${studentId.value}`)
  } catch {
    student.value = null
  } finally {
    loading.value = false
  }
}

const fetchClassHistory = async () => {
  subLoading.value = true
  try {
    classHistory.value = await api.get(`/students/${studentId.value}/class-history`)
  } catch { /* ignore */ }
  finally { subLoading.value = false }
}

const fetchGrades = async () => {
  try {
    grades.value = await api.get(`/students/${studentId.value}/grades`)
  } catch { /* ignore */ }
}

const fetchReportCards = async () => {
  try {
    reportCards.value = await api.get(`/students/${studentId.value}/report-cards`)
  } catch { /* ignore */ }
}

const fetchAttendance = async () => {
  subLoading.value = true
  try {
    const [records, statsRes] = await Promise.all([
      api.get(`/students/${studentId.value}/attendance`, { month: attendanceMonth.value }),
      api.get<{ present: number; absent: number; late: number; sick: number }>(`/students/${studentId.value}/attendance-stats`, { month: attendanceMonth.value }),
    ])
    attendanceRecords.value = records
    Object.assign(attendanceStats, statsRes)
  } catch { /* ignore */ }
  finally { subLoading.value = false }
}

const fetchTahfidz = async () => {
  try {
    tahfidzRecords.value = await api.get(`/students/${studentId.value}/quran-progress`)
  } catch { /* ignore */ }
}

const fetchMutabaah = async () => {
  subLoading.value = true
  try {
    mutabaahRecords.value = await api.get(`/students/${studentId.value}/mutabaah`, { month: mutabaahMonth.value })
  } catch { /* ignore */ }
  finally { subLoading.value = false }
}

const fetchFinance = async () => {
  try {
    const [invoices, payments] = await Promise.all([
      api.get(`/students/${studentId.value}/invoices`),
      api.get(`/students/${studentId.value}/payments`),
    ])
    studentInvoices.value = invoices
    studentPayments.value = payments
  } catch { /* ignore */ }
}

const fetchMedical = async () => {
  try {
    medicalRecords.value = await api.get(`/students/${studentId.value}/medical-records`)
  } catch { /* ignore */ }
}

const fetchDocuments = async () => {
  try {
    studentDocuments.value = await api.get(`/students/${studentId.value}/documents`)
  } catch { /* ignore */ }
}

const fetchBehavior = async () => {
  try {
    behaviorNotes.value = await api.get(`/students/${studentId.value}/behavior-notes`)
  } catch { /* ignore */ }
}

const openEditProfile = () => {
  if (!student.value) return
  editForm.fullName = student.value.fullName
  editForm.birthPlace = student.value.birthPlace
  editForm.birthDate = student.value.birthDate
  editForm.address = student.value.address
  editForm.phone = student.value.phone || ''
  editForm.status = student.value.status
  showProfileForm.value = true
}

const saveProfile = async () => {
  saving.value = true
  try {
    await api.put(`/students/${studentId.value}`, editForm)
    toast.add({ title: t('students.profile_updated'), color: 'success' })
    showProfileForm.value = false
    fetchStudent()
  } catch { /* handled by api */ }
  finally { saving.value = false }
}

const openAddMedical = () => {
  medicalForm.date = $dayjs().format('YYYY-MM-DD')
  medicalForm.complaint = ''
  medicalForm.diagnosis = ''
  medicalForm.treatment = ''
  medicalForm.medication = ''
  showMedicalForm.value = true
}

const saveMedical = async () => {
  savingMedical.value = true
  try {
    await api.post(`/students/${studentId.value}/medical-records`, medicalForm)
    toast.add({ title: t('medical.record_created'), color: 'success' })
    showMedicalForm.value = false
    fetchMedical()
  } catch { /* handled by api */ }
  finally { savingMedical.value = false }
}

const openAddBehavior = () => {
  behaviorForm.type = 'positive'
  behaviorForm.description = ''
  behaviorForm.date = $dayjs().format('YYYY-MM-DD')
  showBehaviorForm.value = true
}

const saveBehavior = async () => {
  savingBehavior.value = true
  try {
    await api.post(`/students/${studentId.value}/behavior-notes`, behaviorForm)
    toast.add({ title: t('students.behavior_note_created'), color: 'success' })
    showBehaviorForm.value = false
    fetchBehavior()
  } catch { /* handled by api */ }
  finally { savingBehavior.value = false }
}

const triggerUpload = () => {
  // file upload is handled by FileUpload component events
}

const handleFileUpload = async (files: File[]) => {
  for (const file of files) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('studentId', studentId.value)
    try {
      await api.upload(`/students/${studentId.value}/documents`, formData)
    } catch { /* ignore */ }
  }
  toast.add({ title: t('documents.upload_success'), color: 'success' })
  fetchDocuments()
}

const downloadDocument = (doc: { id: string; title: string; url: string; createdAt: string }) => {
  window.open(doc.url, '_blank')
}

const viewReportCard = (row: Record<string, unknown>) => {
  window.open(`/api/v1/students/${studentId.value}/report-cards/${row.id}/pdf`, '_blank')
}

onMounted(async () => {
  await fetchStudent()
  if (student.value) {
    await Promise.all([
      fetchClassHistory(),
      fetchGrades(),
      fetchReportCards(),
      fetchAttendance(),
      fetchTahfidz(),
      fetchMutabaah(),
      fetchFinance(),
      fetchMedical(),
      fetchDocuments(),
      fetchBehavior(),
    ])
  }
})
</script>
