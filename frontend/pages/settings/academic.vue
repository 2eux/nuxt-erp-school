<template>
  <div class="space-y-6">
    <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('settings.academic_settings') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('settings.academic_settings_subtitle') }}</p></div>

    <div class="space-y-6">
      <div class="card">
        <div class="card-header">
          <div><h3 class="card-title">{{ $t('settings.academic_years') }}</h3><p class="text-sm text-gray-500 mt-1">{{ $t('settings.academic_years_description') }}</p></div>
          <UButton color="primary" size="sm" icon="i-heroicons-plus" @click="openYearForm">{{ $t('common.add') }}</UButton>
        </div>
        <DataTable :columns="yearColumns" :rows="academicYears" :loading="loading" :empty-title="$t('settings.no_academic_years')" :show-filters="false">
          <template #cell-isActive="{ row }"><StatusBadge :status="row.isActive ? 'active' : 'inactive'" /></template>
          <template #item-actions="{ row }">
            <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editYear(row as Record<string, unknown>)" />
            <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-check" @click="setActiveYear(row)" />
          </template>
        </DataTable>
      </div>

      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('settings.semester_config') }}</h3></div>
        <DataTable :columns="semesterColumns" :rows="semesters" :loading="loading" :empty-title="$t('settings.no_semesters')" :show-filters="false">
          <template #cell-isActive="{ row }"><StatusBadge :status="row.isActive ? 'active' : 'inactive'" /></template>
          <template #item-actions="{ row }">
            <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editSemester(row as Record<string, unknown>)" />
            <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-check" @click="setActiveSemester(row)" />
          </template>
        </DataTable>
      </div>

      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('settings.grading_system') }}</h3></div>
        <div class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <UFormGroup :label="$t('settings.passing_score')"><UInput v-model.number="grading.passingScore" type="number" color="gray" /></UFormGroup>
            <UFormGroup :label="$t('settings.max_score')"><UInput v-model.number="grading.maxScore" type="number" color="gray" /></UFormGroup>
            <UFormGroup :label="$t('settings.grading_type')"><USelect v-model="grading.type" :options="[{label:'0-100',value:'numeric'},{label:'A-E',value:'letter'},{label:'4.0 Scale',value:'gpa'}]" color="gray" /></UFormGroup>
          </div>
          <div v-if="grading.type === 'letter'" class="space-y-2">
            <p class="text-sm font-medium text-gray-900 dark:text-white">{{ $t('settings.grade_ranges') }}</p>
            <div v-for="g in gradeRanges" :key="g.letter" class="flex items-center gap-3">
              <span class="text-sm font-bold w-8">{{ g.letter }}</span>
              <input v-model.number="g.min" type="number" class="w-20 px-2 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white" />
              <span class="text-sm text-gray-500">-</span>
              <input v-model.number="g.max" type="number" class="w-20 px-2 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white" />
            </div>
          </div>
          <div class="flex justify-end"><UButton color="primary" :loading="saving" @click="saveGrading">{{ $t('common.save') }}</UButton></div>
        </div>
      </div>
    </div>

    <FormDialog v-model="showYearForm" :title="editingYear ? $t('settings.edit_year') : $t('settings.add_year')" :loading="savingYear" @submit="saveYear" @cancel="showYearForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('academic.year_name')" required><UInput v-model="yearForm.name" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('common.start_date')"><UInput v-model="yearForm.startDate" type="date" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('common.end_date')"><UInput v-model="yearForm.endDate" type="date" color="gray" /></UFormGroup>
        </div>
        <UToggle v-model="yearForm.isActive" :label="$t('settings.set_active')" />
      </div>
    </FormDialog>

    <FormDialog v-model="showSemesterForm" :title="$t('settings.edit_semester')" :loading="savingSemester" @submit="saveSemester" @cancel="showSemesterForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('academic.year')"><USelect v-model="semesterForm.academicYearId" :options="yearOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('settings.name')" required><UInput v-model="semesterForm.name" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('common.start_date')"><UInput v-model="semesterForm.startDate" type="date" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('common.end_date')"><UInput v-model="semesterForm.endDate" type="date" color="gray" /></UFormGroup>
        </div>
        <UToggle v-model="semesterForm.isActive" :label="$t('settings.set_active')" />
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const toast = useToast()
const loading = ref(false); const saving = ref(false); const savingYear = ref(false); const savingSemester = ref(false)
const showYearForm = ref(false); const showSemesterForm = ref(false)
const editingYear = ref(false); const editYearId = ref<string | null>(null)
const editingSemesterId = ref<string | null>(null)
const academicYears = ref<Record<string, unknown>[]>([]); const semesters = ref<Record<string, unknown>[]>([])

const yearColumns: TableColumn[] = [
  { key: 'name', label: 'academic.year_name' },
  { key: 'startDate', label: 'common.start_date', type: 'date' },
  { key: 'endDate', label: 'common.end_date', type: 'date' },
  { key: 'isActive', label: 'common.status', type: 'status' },
]
const semesterColumns: TableColumn[] = [
  { key: 'name', label: 'settings.name' },
  { key: 'startDate', label: 'common.start_date', type: 'date' },
  { key: 'endDate', label: 'common.end_date', type: 'date' },
  { key: 'isActive', label: 'common.status', type: 'status' },
]

const yearOpts = computed(() => academicYears.value.map(y => ({ label: y.name as string, value: y.id as string })))
const yearForm = reactive({ name: '', startDate: '', endDate: '', isActive: false })
const semesterForm = reactive({ academicYearId: '', name: '', startDate: '', endDate: '', isActive: false })

const grading = reactive({ passingScore: 75, maxScore: 100, type: 'numeric' })
const gradeRanges = reactive([
  { letter: 'A', min: 90, max: 100 }, { letter: 'B', min: 80, max: 89 },
  { letter: 'C', min: 70, max: 79 }, { letter: 'D', min: 60, max: 69 }, { letter: 'E', min: 0, max: 59 },
])

const fetchData = async () => {
  loading.value = true
  try { const [years, sems] = await Promise.all([api.get('/academic-years'), api.get('/terms')]); academicYears.value = years; semesters.value = sems } catch {} finally { loading.value = false }
}
const openYearForm = () => { editingYear.value = false; editYearId.value = null; Object.assign(yearForm, { name: '', startDate: '', endDate: '', isActive: false }); showYearForm.value = true }
const editYear = (row: Record<string, unknown>) => { editingYear.value = true; editYearId.value = row.id as string; Object.assign(yearForm, row); showYearForm.value = true }
const saveYear = async () => { savingYear.value = true; try { if (editingYear.value && editYearId.value) { await api.put(`/academic-years/${editYearId.value}`, yearForm) } else { await api.post('/academic-years', yearForm) } toast.add({ title: t('settings.year_saved'), color: 'success' }); showYearForm.value = false; fetchData() } catch {} finally { savingYear.value = false } }
const setActiveYear = async (row: Record<string, unknown>) => { try { await api.patch(`/academic-years/${row.id}/activate`); toast.add({ title: t('settings.year_activated'), color: 'success' }); fetchData() } catch {} }
const editSemester = (row: Record<string, unknown>) => { editingSemesterId.value = row.id as string; Object.assign(semesterForm, row); showSemesterForm.value = true }
const saveSemester = async () => { savingSemester.value = true; try { await api.put(`/terms/${editingSemesterId.value}`, semesterForm); toast.add({ title: t('settings.semester_saved'), color: 'success' }); showSemesterForm.value = false; fetchData() } catch {} finally { savingSemester.value = false } }
const setActiveSemester = async (row: Record<string, unknown>) => { try { await api.patch(`/terms/${row.id}/activate`); toast.add({ title: t('settings.semester_activated'), color: 'success' }); fetchData() } catch {} }
const saveGrading = async () => { saving.value = true; try { await api.put('/settings/grading', grading); toast.add({ title: t('settings.grading_saved'), color: 'success' }) } catch {} finally { saving.value = false } }

onMounted(() => fetchData())
</script>
