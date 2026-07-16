<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('students.title') }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('students.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <UButton
          v-if="permissions.can('students.export')"
          color="gray"
          variant="outline"
          size="sm"
          icon="i-heroicons-arrow-down-tray"
          @click="exportData"
        >
          {{ $t('common.export') }}
        </UButton>
        <UButton
          v-if="selectedStudents.length > 1 && permissions.can('students.promote')"
          color="gray"
          variant="outline"
          size="sm"
          icon="i-heroicons-arrows-up-down"
          @click="openMassPromotion"
        >
          {{ $t('students.mass_promotion') }}
        </UButton>
        <UButton
          v-if="permissions.can('students.create')"
          color="primary"
          size="sm"
          icon="i-heroicons-plus"
          @click="openAddDialog"
        >
          {{ $t('students.add_student') }}
        </UButton>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard
        :label="$t('students.total_students')"
        :value="stats.totalStudents"
        icon="i-heroicons-users"
        color="blue"
        :loading="statsLoading"
      />
      <StatCard
        :label="$t('students.active_students')"
        :value="stats.activeStudents"
        icon="i-heroicons-check-circle"
        color="emerald"
        :loading="statsLoading"
      />
      <StatCard
        :label="$t('students.new_this_year')"
        :value="stats.newStudents"
        icon="i-heroicons-user-plus"
        color="amber"
        :loading="statsLoading"
      />
      <StatCard
        :label="$t('students.graduated')"
        :value="stats.graduated"
        icon="i-heroicons-academic-cap"
        color="purple"
        :loading="statsLoading"
      />
    </div>

    <div class="card">
      <DataFilter
        :filter-fields="filterFields"
        :searchable="true"
        @apply="handleFilter"
      />
    </div>

    <DataTable
      :columns="columns"
      :rows="students"
      :loading="loading"
      :empty-title="$t('students.no_students')"
      :empty-description="$t('students.no_students_description')"
      :empty-icon="'i-heroicons-users'"
      :show-export="false"
      :selectable="true"
      :selected-ids="selectedStudents"
      @update:selected-ids="selectedStudents = $event"
      @sort="handleSort"
      @page-change="handlePageChange"
      @page-size-change="handlePageSizeChange"
    >
      <template #cell-fullName="{ row }">
        <div class="flex items-center gap-3">
          <UAvatar
            :src="(row.photo as string) || undefined"
            :alt="(row.fullName as string)"
            size="sm"
          />
          <div>
            <NuxtLink
              :to="`/students/${row.id}`"
              class="text-sm font-medium text-gray-900 dark:text-white hover:text-brand-600 dark:hover:text-brand-400"
            >
              {{ row.fullName }}
            </NuxtLink>
            <p class="text-xs text-gray-500">{{ row.nis }}</p>
          </div>
        </div>
      </template>

      <template #cell-status="{ row }">
        <StatusBadge :status="row.status as string" />
      </template>

      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton
            color="gray"
            variant="ghost"
            size="xs"
            icon="i-heroicons-eye"
            :to="`/students/${row.id}`"
          />
          <UButton
            color="gray"
            variant="ghost"
            size="xs"
            icon="i-heroicons-pencil-square"
            @click="openEditDialog(row as StudentRecord)"
          />
          <UPopover v-if="permissions.can('students.delete')" :ui="{ width: 'w-48' }">
            <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-ellipsis-vertical" />
            <template #panel>
              <div class="p-1">
                <UButton
                  color="gray"
                  variant="ghost"
                  block
                  size="xs"
                  icon="i-heroicons-arrow-up-on-square"
                  @click="promoteStudent(row as StudentRecord)"
                >
                  {{ $t('students.promote') }}
                </UButton>
                <UButton
                  color="red"
                  variant="ghost"
                  block
                  size="xs"
                  icon="i-heroicons-trash"
                  @click="confirmDelete(row as StudentRecord)"
                >
                  {{ $t('common.delete') }}
                </UButton>
              </div>
            </template>
          </UPopover>
        </div>
      </template>
    </DataTable>

    <FormDialog
      v-model="showAddDialog"
      :title="editingStudent ? $t('students.edit_student') : $t('students.add_student')"
      :loading="saving"
      @submit="saveStudent"
      @cancel="closeForm"
    >
      <div class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <UFormGroup :label="$t('students.nis')" required>
            <UInput v-model="form.nis" color="gray" :disabled="!!editingStudent" />
          </UFormGroup>
          <UFormGroup :label="$t('students.nisn')">
            <UInput v-model="form.nisn" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.full_name')" required class="sm:col-span-2">
            <UInput v-model="form.fullName" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('common.gender')" required>
            <USelect
              v-model="form.gender"
              :options="genderOptions"
              color="gray"
            />
          </UFormGroup>
          <UFormGroup :label="$t('students.birth_place')">
            <UInput v-model="form.birthPlace" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.birth_date')">
            <UInput v-model="form.birthDate" type="date" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.religion')">
            <USelect
              v-model="form.religion"
              :options="religionOptions"
              color="gray"
            />
          </UFormGroup>
          <UFormGroup :label="$t('students.class')" class="sm:col-span-2">
            <USelect v-model="form.classId" :options="classOptions" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.address')" class="sm:col-span-2">
            <UTextarea v-model="form.address" color="gray" :rows="2" />
          </UFormGroup>
          <UFormGroup :label="$t('students.city')">
            <UInput v-model="form.city" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.phone')">
            <UInput v-model="form.phone" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.status')">
            <USelect
              v-model="form.status"
              :options="statusOptions"
              color="gray"
            />
          </UFormGroup>
          <UFormGroup :label="$t('students.father_name')">
            <UInput v-model="form.fatherName" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.father_phone')">
            <UInput v-model="form.fatherPhone" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.mother_name')">
            <UInput v-model="form.motherName" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.mother_phone')">
            <UInput v-model="form.motherPhone" color="gray" />
          </UFormGroup>
        </div>
      </div>
    </FormDialog>

    <FormDialog
      v-model="showMassPromotionDialog"
      :title="$t('students.mass_promotion')"
      :loading="massPromoting"
      @submit="submitMassPromotion"
      @cancel="showMassPromotionDialog = false"
    >
      <div class="space-y-4">
        <p class="text-sm text-gray-500 dark:text-gray-400">
          {{ $t('students.mass_promotion_description', { count: selectedStudents.length }) }}
        </p>
        <UFormGroup :label="$t('students.target_class')" required>
          <USelect v-model="promotionTarget" :options="classOptions" color="gray" />
        </UFormGroup>
        <UFormGroup :label="$t('students.promotion_date')">
          <UInput v-model="promotionDate" type="date" color="gray" />
        </UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog
      v-model="showDeleteConfirm"
      :title="$t('students.delete_student')"
      :description="$t('students.delete_student_confirm', { name: deleteTarget?.fullName || '' })"
      :loading="deleting"
      @confirm="deleteStudent"
      @cancel="showDeleteConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
import type { Student, TableColumn, PaginatedResponse } from '~/types'

interface StudentRecord extends Record<string, unknown> {
  id: string
  nis: string
  nisn: string
  fullName: string
  gender: string
  birthPlace: string
  birthDate: string
  religion: string
  address: string
  city: string
  phone: string | null
  classId: string | null
  className: string | null
  status: string
  fatherName: string
  fatherPhone: string
  motherName: string
  motherPhone: string
  photo: string | null
}

definePageMeta({
  middleware: ['auth'],
})

const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const { $dayjs } = useNuxtApp()
const toast = useToast()

const loading = ref(false)
const statsLoading = ref(false)
const saving = ref(false)
const massPromoting = ref(false)
const deleting = ref(false)

const students = ref<StudentRecord[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const selectedStudents = ref<string[]>([])

const currentPage = ref(1)
const pageSize = ref(15)
const totalItems = ref(0)
const sortKey = ref('fullName')
const sortOrder = ref('asc')
const filterParams = ref<Record<string, unknown>>({})

const showAddDialog = ref(false)
const showMassPromotionDialog = ref(false)
const showDeleteConfirm = ref(false)
const editingStudent = ref<Student | null>(null)
const deleteTarget = ref<StudentRecord | null>(null)
const promotionTarget = ref('')
const promotionDate = ref('')

const stats = reactive({
  totalStudents: 0,
  activeStudents: 0,
  newStudents: 0,
  graduated: 0,
})

const columns: TableColumn[] = [
  { key: 'fullName', label: 'students.name', sortable: true },
  { key: 'nis', label: 'students.nis', sortable: true },
  { key: 'nisn', label: 'students.nisn', sortable: true },
  { key: 'className', label: 'students.class', sortable: true },
  { key: 'fatherName', label: 'students.parent' },
  { key: 'status', label: 'common.status', type: 'status' },
]

const filterFields = [
  { key: 'classId', label: 'students.class', type: 'select' as const, options: [], placeholder: t('common.select') },
  { key: 'status', label: 'common.status', type: 'select' as const, options: [
    { label: t('status.active'), value: 'active' },
    { label: t('status.inactive'), value: 'inactive' },
    { label: t('status.graduated'), value: 'graduated' },
    { label: t('status.transferred'), value: 'transferred' },
    { label: t('status.dropped'), value: 'dropped' },
  ]},
  { key: 'academicYearId', label: 'academic_year.title', type: 'select' as const, options: [] },
]

const genderOptions = [
  { label: t('common.male'), value: 'male' },
  { label: t('common.female'), value: 'female' },
]

const religionOptions = [
  { label: 'Islam', value: 'Islam' },
  { label: 'Kristen', value: 'Kristen' },
  { label: 'Katolik', value: 'Katolik' },
  { label: 'Hindu', value: 'Hindu' },
  { label: 'Buddha', value: 'Buddha' },
  { label: 'Konghucu', value: 'Konghucu' },
]

const statusOptions = [
  { label: t('status.active'), value: 'active' },
  { label: t('status.inactive'), value: 'inactive' },
  { label: t('status.graduated'), value: 'graduated' },
  { label: t('status.transferred'), value: 'transferred' },
  { label: t('status.dropped'), value: 'dropped' },
]

const form = reactive({
  nis: '',
  nisn: '',
  fullName: '',
  gender: 'male' as string,
  birthPlace: '',
  birthDate: '',
  religion: 'Islam',
  address: '',
  city: '',
  phone: '',
  classId: '',
  status: 'active' as string,
  fatherName: '',
  fatherPhone: '',
  motherName: '',
  motherPhone: '',
})

const fetchClasses = async () => {
  try {
    const res = await api.get<{ id: string; name: string }[]>('/classes', { limit: 100 })
    classOptions.value = res.map((c: { id: string; name: string }) => ({ label: c.name, value: c.id }))
    filterFields[0].options = classOptions.value
  } catch { /* ignore */ }
}

const fetchStats = async () => {
  statsLoading.value = true
  try {
    const res = await api.get<{
      totalStudents: number; activeStudents: number; newStudents: number; graduated: number
    }>('/students/stats')
    Object.assign(stats, res)
  } catch { /* ignore */ }
  finally { statsLoading.value = false }
}

const fetchStudents = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      sortBy: sortKey.value,
      sortOrder: sortOrder.value,
      ...filterParams.value,
    }
    const res = await api.paginate<StudentRecord>('/students', params)
    students.value = res.data as unknown as StudentRecord[]
    totalItems.value = res.pagination.total
  } catch { /* ignore */ }
  finally { loading.value = false }
}

const handleFilter = (filters: Record<string, unknown>) => {
  filterParams.value = filters
  currentPage.value = 1
  fetchStudents()
}

const handleSort = (key: string, order: string) => {
  sortKey.value = key
  sortOrder.value = order as 'asc' | 'desc'
  fetchStudents()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchStudents()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchStudents()
}

const openAddDialog = () => {
  editingStudent.value = null
  Object.assign(form, {
    nis: '', nisn: '', fullName: '', gender: 'male', birthPlace: '', birthDate: '',
    religion: 'Islam', address: '', city: '', phone: '', classId: '', status: 'active',
    fatherName: '', fatherPhone: '', motherName: '', motherPhone: '',
  })
  showAddDialog.value = true
}

const openEditDialog = (student: StudentRecord) => {
  editingStudent.value = student as unknown as Student
  Object.assign(form, {
    nis: student.nis, nisn: student.nisn, fullName: student.fullName,
    gender: student.gender, birthPlace: student.birthPlace, birthDate: student.birthDate,
    religion: student.religion, address: student.address, city: student.city,
    phone: student.phone, classId: student.classId, status: student.status,
    fatherName: student.fatherName, fatherPhone: student.fatherPhone,
    motherName: student.motherName, motherPhone: student.motherPhone,
  })
  showAddDialog.value = true
}

const closeForm = () => {
  showAddDialog.value = false
  editingStudent.value = null
}

const saveStudent = async () => {
  saving.value = true
  try {
    if (editingStudent.value) {
      await api.put(`/students/${editingStudent.value.id}`, { ...form })
      toast.add({ title: t('students.student_updated'), color: 'success' })
    } else {
      await api.post('/students', { ...form })
      toast.add({ title: t('students.student_created'), color: 'success' })
    }
    showAddDialog.value = false
    editingStudent.value = null
    fetchStudents()
    fetchStats()
  } catch { /* handled by api */ }
  finally { saving.value = false }
}

const confirmDelete = (student: StudentRecord) => {
  deleteTarget.value = student
  showDeleteConfirm.value = true
}

const deleteStudent = async () => {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.delete(`/students/${deleteTarget.value.id}`)
    toast.add({ title: t('students.student_deleted'), color: 'success' })
    showDeleteConfirm.value = false
    deleteTarget.value = null
    fetchStudents()
    fetchStats()
  } catch { /* handled by api */ }
  finally { deleting.value = false }
}

const promoteStudent = async (student: StudentRecord) => {
  selectedStudents.value = [student.id]
  promotionTarget.value = ''
  promotionDate.value = $dayjs().format('YYYY-MM-DD')
  showMassPromotionDialog.value = true
}

const openMassPromotion = () => {
  promotionTarget.value = ''
  promotionDate.value = $dayjs().format('YYYY-MM-DD')
  showMassPromotionDialog.value = true
}

const submitMassPromotion = async () => {
  if (!promotionTarget.value) return
  massPromoting.value = true
  try {
    await api.post('/students/mass-promotion', {
      studentIds: selectedStudents.value,
      targetClassId: promotionTarget.value,
      promotionDate: promotionDate.value,
    })
    toast.add({ title: t('students.promotion_success'), color: 'success' })
    showMassPromotionDialog.value = false
    selectedStudents.value = []
    fetchStudents()
  } catch { /* handled by api */ }
  finally { massPromoting.value = false }
}

const exportData = async () => {
  try {
    const params = { ...filterParams.value, sortBy: sortKey.value, sortOrder: sortOrder.value }
    const res = await api.get<{ data: Student[] }>('/students/export', params)
    const csv = convertToCSV(res.data)
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `students_export_${$dayjs().format('YYYY-MM-DD')}.csv`
    link.click()
    URL.revokeObjectURL(url)
    toast.add({ title: t('common.export_success'), color: 'success' })
  } catch { /* handled by api */ }
}

const convertToCSV = (data: unknown[]): string => {
  if (!data || data.length === 0) return ''
  const headers = Object.keys(data[0] as Record<string, unknown>)
  const csv = [headers.join(',')]
  for (const row of data) {
    csv.push(headers.map(h => `"${(row as Record<string, unknown>)[h] || ''}"`).join(','))
  }
  return csv.join('\n')
}

onMounted(() => {
  fetchClasses()
  fetchStats()
  fetchStudents()
})
</script>
