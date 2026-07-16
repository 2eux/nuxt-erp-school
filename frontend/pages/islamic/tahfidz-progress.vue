<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('tahfidz.progress') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('tahfidz.progress_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedGroup" :options="groupOptions" :placeholder="$t('tahfidz.select_group')" color="gray" size="sm" class="w-40" @change="fetchProgress" />
        <USelect v-model="selectedStudent" :options="studentOpts" :placeholder="$t('tahfidz.select_student')" color="gray" size="sm" class="w-48" @change="fetchProgress" />
        <UInput v-model="startDate" type="date" color="gray" size="sm" @change="fetchProgress" />
        <UInput v-model="endDate" type="date" color="gray" size="sm" @change="fetchProgress" />
        <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAddEntry">{{ $t('tahfidz.add_entry') }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="card lg:col-span-2">
        <div class="card-header"><h3 class="card-title">{{ $t('tahfidz.memorization_trend') }}</h3></div>
        <div class="h-80">
          <ApexChart type="line" height="100%" :options="progressChart.options" :series="progressChart.series" />
        </div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('tahfidz.juz_completion') }}</h3></div>
        <div class="h-80">
          <ApexChart type="donut" height="100%" :options="juzChart.options" :series="juzChart.series" />
        </div>
      </div>
    </div>

    <DataTable :columns="columns" :rows="progressRecords" :loading="loading" :empty-title="$t('tahfidz.no_progress')" :show-export="false">
      <template #cell-memorizationType="{ row }">
        <StatusBadge :status="(row.memorizationType as string) || 'tahfidz'" />
      </template>
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'in_progress'" /></template>
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editEntry(row as Record<string, unknown>)" />
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
      </template>
    </DataTable>

    <FormDialog v-model="showEntryForm" :title="editing ? $t('tahfidz.edit_entry') : $t('tahfidz.add_entry')" :loading="saving" @submit="saveEntry" @cancel="showEntryForm=false">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('tahfidz.student')"><USelect v-model="form.studentId" :options="studentOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.date')"><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.surah_start')"><UInput v-model="form.surahStart" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.ayah_start')"><UInput v-model.number="form.ayahStart" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.surah_end')"><UInput v-model="form.surahEnd" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.ayah_end')"><UInput v-model.number="form.ayahEnd" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.type')"><USelect v-model="form.memorizationType" :options="typeOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.status')"><USelect v-model="form.status" :options="statusOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.note')" class="sm:col-span-2"><UTextarea v-model="form.note" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('tahfidz.delete_entry')" :loading="deleting" @confirm="deleteEntry" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const { $dayjs } = useNuxtApp()
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const showEntryForm = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const deleteTarget = ref<Record<string, unknown> | null>(null)
const progressRecords = ref<Record<string, unknown>[]>([])
const groupOptions = ref<{ label: string; value: string }[]>([])
const studentOpts = ref<{ label: string; value: string }[]>([])
const selectedGroup = ref('')
const selectedStudent = ref('')
const startDate = ref($dayjs().startOf('month').format('YYYY-MM-DD'))
const endDate = ref($dayjs().format('YYYY-MM-DD'))

const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name' },
  { key: 'surahStart', label: 'quran.surah' },
  { key: 'ayahStart', label: 'quran.ayah' },
  { key: 'surahEnd', label: 'quran.to_surah' },
  { key: 'ayahEnd', label: 'quran.to_ayah' },
  { key: 'memorizationType', label: 'tahfidz.type' },
  { key: 'status', label: 'common.status', type: 'status' },
  { key: 'date', label: 'common.date', type: 'date' },
]
const typeOptions = [
  { label: t('tahfidz.tahfidz'), value: 'tahfidz' }, { label: t('tahfidz.murojaah'), value: 'murojaah' }, { label: t('tahfidz.tilawah'), value: 'tilawah' },
]
const statusOpts = [
  { label: t('tahfidz.memorized'), value: 'memorized' }, { label: t('tahfidz.in_progress'), value: 'in_progress' }, { label: t('tahfidz.not_memorized'), value: 'not_memorized' },
]
const form = reactive({ studentId: '', date: $dayjs().format('YYYY-MM-DD'), surahStart: 'Al-Fatihah', ayahStart: 1, surahEnd: 'Al-Fatihah', ayahEnd: 1, memorizationType: 'tahfidz', status: 'in_progress', note: '' })

const progressChart = computed(() => ({
  series: [{ name: t('tahfidz.ayahs_memorized'), data: [10, 25, 18, 35, 28, 45, 38, 52, 48, 60] }],
  options: { chart: { type: 'line' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669'], stroke: { curve: 'smooth' as const, width: 2 }, xaxis: { categories: ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct'] }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))
const juzChart = computed(() => ({
  series: [5, 8, 12, 6],
  options: { chart: { type: 'donut' as const, background: 'transparent' }, labels: [t('tahfidz.juz') + ' 30', t('tahfidz.juz') + ' 29', t('tahfidz.juz') + ' 1', t('tahfidz.juz') + ' 2'], colors: ['#059669','#10b981','#34d399','#6ee7b7'], legend: { position: 'bottom' as const } },
}))

const fetchProgress = async () => { loading.value = true; try { progressRecords.value = await api.paginate('/tahfidz/progress', { groupId: selectedGroup.value, studentId: selectedStudent.value, startDate: startDate.value, endDate: endDate.value }).then(r => r.data) } catch {} finally { loading.value = false } }

const openAddEntry = () => { editing.value = false; editId.value = null; form.date = $dayjs().format('YYYY-MM-DD'); form.studentId = selectedStudent.value; form.surahStart = 'Al-Fatihah'; form.ayahStart = 1; form.surahEnd = 'Al-Fatihah'; form.ayahEnd = 1; form.memorizationType = 'tahfidz'; form.status = 'in_progress'; form.note = ''; showEntryForm.value = true }
const editEntry = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showEntryForm.value = true }

const saveEntry = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/tahfidz/progress/${editId.value}`, form); toast.add({ title: t('tahfidz.entry_updated'), color: 'success' }) }
    else { await api.post('/tahfidz/progress', form); toast.add({ title: t('tahfidz.entry_created'), color: 'success' }) }
    showEntryForm.value = false; fetchProgress()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteEntry = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/tahfidz/progress/${deleteTarget.value.id}`); toast.add({ title: t('tahfidz.entry_deleted'), color: 'success' }); showDelete.value = false; fetchProgress() } catch {} finally { deleting.value = false } }

const fetchOptions = async () => {
  try { groupOptions.value = (await api.get<{id:string;name:string}[]>('/tahfidz/groups')).map(g => ({ label: g.name, value: g.id })) } catch {}
  try { studentOpts.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })) } catch {}
}

onMounted(() => { fetchOptions(); fetchProgress() })
</script>
