<template>
  <div class="space-y-6">
    <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.subtitle') }}</p></div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('academic.academic_years')" :value="stats.academicYears" icon="i-heroicons-calendar" color="blue" :loading="loading" />
      <StatCard :label="$t('academic.active_semester')" :value="activeSemesterName" icon="i-heroicons-clock" color="emerald" :loading="loading" />
      <StatCard :label="$t('academic.total_classes')" :value="stats.totalClasses" icon="i-heroicons-building-office" color="amber" :loading="loading" />
      <StatCard :label="$t('academic.total_subjects')" :value="stats.totalSubjects" icon="i-heroicons-book-open" color="purple" :loading="loading" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('academic.current_academic_year') }}</h3></div>
        <div v-if="loading"><LoadingSkeleton type="detail" /></div>
        <div v-else-if="!activeYear" class="py-4"><EmptyState :title="$t('academic.no_active_year')" icon="i-heroicons-calendar" /></div>
        <div v-else class="space-y-3">
          <div class="flex items-center justify-between"><span class="text-sm text-gray-500">{{ $t('academic.year_name') }}</span><span class="text-sm font-medium text-gray-900 dark:text-white">{{ activeYear.name }}</span></div>
          <div class="flex items-center justify-between"><span class="text-sm text-gray-500">{{ $t('academic.period') }}</span><span class="text-sm font-medium text-gray-900 dark:text-white">{{ $dayjs(activeYear.startDate).format('DD MMM YYYY') }} - {{ $dayjs(activeYear.endDate).format('DD MMM YYYY') }}</span></div>
          <div class="flex items-center justify-between"><span class="text-sm text-gray-500">{{ $t('academic.active_term') }}</span><StatusBadge :status="'active'" /></div>
        </div>
      </div>

      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('academic.quick_links') }}</h3></div>
        <div class="grid grid-cols-2 gap-3">
          <NuxtLink v-for="link in quickLinks" :key="link.to" :to="link.to" class="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
            <UIcon :name="link.icon" class="w-5 h-5 text-brand-600 dark:text-brand-400" />
            <span class="text-sm font-medium text-gray-900 dark:text-white">{{ link.label }}</span>
          </NuxtLink>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header"><h3 class="card-title">{{ $t('academic.academic_years') }}</h3><UButton v-if="permissions.can('academic.manage')" color="primary" size="xs" icon="i-heroicons-plus" @click="openYearForm">{{ $t('common.add') }}</UButton></div>
      <DataTable :columns="yearColumns" :rows="academicYears" :loading="loading" :empty-title="$t('academic.no_years')" :show-filters="false">
        <template #cell-isActive="{ row }"><StatusBadge :status="row.isActive ? 'active' : 'inactive'" /></template>
        <template #item-actions="{ row }">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editYear(row as Record<string, unknown>)" />
        </template>
      </DataTable>
    </div>

    <FormDialog v-model="showYearForm" :title="editingYear ? $t('academic.edit_year') : $t('academic.add_year')" :loading="saving" @submit="saveYear" @cancel="showYearForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('academic.year_name')" required><UInput v-model="yearForm.name" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('common.start_date')" required><UInput v-model="yearForm.startDate" type="date" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('common.end_date')" required><UInput v-model="yearForm.endDate" type="date" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('academic.set_active')"><UToggle v-model="yearForm.isActive" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn, AcademicYear } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const { $dayjs } = useNuxtApp()
const schoolStore = useSchoolStore()
const loading = ref(false)
const saving = ref(false)
const showYearForm = ref(false)
const editingYear = ref(false)
const editYearId = ref<string | null>(null)
const academicYears = ref<Record<string, unknown>[]>([])
const activeYear = ref<AcademicYear | null>(null)

const stats = reactive({ academicYears: 0, totalClasses: 0, totalSubjects: 0 })
const activeSemesterName = computed(() => {
  const year = schoolStore.currentAcademicYear
  if (!year) return '-'
  const activeTerm = year.terms?.find(t => t.isActive)
  return activeTerm ? `${activeTerm.name} - ${year.name}` : year.name
})

const quickLinks = [
  { to: '/academic/classes', label: t('academic.classes'), icon: 'i-heroicons-building-office' },
  { to: '/academic/subjects', label: t('academic.subjects'), icon: 'i-heroicons-book-open' },
  { to: '/academic/curriculum', label: t('academic.curriculum'), icon: 'i-heroicons-clipboard-document-list' },
  { to: '/academic/schedules', label: t('academic.schedules'), icon: 'i-heroicons-clock' },
  { to: '/academic/attendance', label: t('attendance.title'), icon: 'i-heroicons-check-badge' },
  { to: '/academic/exams', label: t('academic.exams'), icon: 'i-heroicons-document-text' },
  { to: '/academic/gradebook', label: t('academic.gradebook'), icon: 'i-heroicons-table-cells' },
  { to: '/academic/report-cards', label: t('academic.report_cards'), icon: 'i-heroicons-document-chart-bar' },
]

const yearColumns: TableColumn[] = [
  { key: 'name', label: 'academic.year_name', sortable: true },
  { key: 'startDate', label: 'common.start_date', type: 'date' },
  { key: 'endDate', label: 'common.end_date', type: 'date' },
  { key: 'isActive', label: 'common.status', type: 'status' },
]

const yearForm = reactive({ name: '', startDate: '', endDate: '', isActive: false })

const fetchData = async () => {
  loading.value = true
  try {
    const [years, activeYearRes, statsRes] = await Promise.all([
      api.get<AcademicYear[]>('/academic-years'),
      api.get<AcademicYear>('/academic-years/active'),
      api.get<{ academicYears: number; totalClasses: number; totalSubjects: number }>('/academic/stats'),
    ])
    academicYears.value = years.map(y => ({ ...y, isActive: y.isActive ? 'active' : 'inactive' }))
    activeYear.value = activeYearRes
    Object.assign(stats, statsRes)
  } catch {/* ignore */}
  finally { loading.value = false }
}

const openYearForm = () => { editingYear.value = false; editYearId.value = null; yearForm.name = ''; yearForm.startDate = ''; yearForm.endDate = ''; yearForm.isActive = false; showYearForm.value = true }
const editYear = (row: Record<string, unknown>) => { editingYear.value = true; editYearId.value = row.id as string; Object.assign(yearForm, row); showYearForm.value = true }

const saveYear = async () => {
  saving.value = true
  try {
    if (editingYear.value && editYearId.value) { await api.put(`/academic-years/${editYearId.value}`, yearForm); toast.add({ title: t('academic.year_updated'), color: 'success' }) }
    else { await api.post('/academic-years', yearForm); toast.add({ title: t('academic.year_created'), color: 'success' }) }
    showYearForm.value = false; fetchData()
  } catch {} finally { saving.value = false }
}

onMounted(() => fetchData())
</script>
