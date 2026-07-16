<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('islamic.prayer_attendance') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('islamic.prayer_attendance_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedClass" :options="classOptions" :placeholder="$t('academic.select_class')" color="gray" size="sm" class="w-36" @change="fetchData" />
        <UInput v-model="selectedDate" type="date" color="gray" size="sm" @change="fetchData" />
        <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openBulkEntry">{{ $t('attendance.bulk_entry') }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-5 gap-4">
      <StatCard :label="$t('prayer.fajr')" :value="`${stats.fajr}%`" icon="i-heroicons-sun" color="amber" :loading="loading" />
      <StatCard :label="$t('prayer.dhuhr')" :value="`${stats.dhuhr}%`" icon="i-heroicons-sun" color="sky" :loading="loading" />
      <StatCard :label="$t('prayer.asr')" :value="`${stats.asr}%`" icon="i-heroicons-sun" color="orange" :loading="loading" />
      <StatCard :label="$t('prayer.maghrib')" :value="`${stats.maghrib}%`" icon="i-heroicons-moon" color="purple" :loading="loading" />
      <StatCard :label="$t('prayer.isha')" :value="`${stats.isha}%`" icon="i-heroicons-moon" color="indigo" :loading="loading" />
    </div>

    <div class="card">
      <div class="card-header"><h3 class="card-title">{{ $t('islamic.prayer_compliance') }}</h3></div>
      <div class="h-80">
        <ApexChart type="bar" height="100%" :options="prayerChart.options" :series="prayerChart.series" />
      </div>
    </div>

    <DataTable :columns="columns" :rows="records" :loading="loading" :empty-title="$t('islamic.no_prayer_records')" :show-filters="false">
      <template #cell-fajr="{ row }"><UIcon :name="row.fajr ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.fajr ? 'text-emerald-500' : 'text-gray-300'" /></template>
      <template #cell-dhuhr="{ row }"><UIcon :name="row.dhuhr ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.dhuhr ? 'text-emerald-500' : 'text-gray-300'" /></template>
      <template #cell-asr="{ row }"><UIcon :name="row.asr ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.asr ? 'text-emerald-500' : 'text-gray-300'" /></template>
      <template #cell-maghrib="{ row }"><UIcon :name="row.maghrib ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.maghrib ? 'text-emerald-500' : 'text-gray-300'" /></template>
      <template #cell-isha="{ row }"><UIcon :name="row.isha ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.isha ? 'text-emerald-500' : 'text-gray-300'" /></template>
    </DataTable>

    <FormDialog v-model="showBulkDialog" :title="$t('attendance.bulk_entry')" :loading="bulkSaving" @submit="submitBulk" @cancel="showBulkDialog=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('common.date')"><UInput v-model="bulkDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.class')"><USelect v-model="bulkClass" :options="classOptions" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-3">
          <UFormGroup v-for="p in prayers" :key="p.key" :label="p.label"><UToggle v-model="bulkPrayers[p.key]" /></UFormGroup>
        </div>
      </div>
    </FormDialog>
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
const bulkSaving = ref(false)
const showBulkDialog = ref(false)
const records = ref<Record<string, unknown>[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const selectedClass = ref('')
const selectedDate = ref($dayjs().format('YYYY-MM-DD'))
const bulkDate = ref($dayjs().format('YYYY-MM-DD'))
const bulkClass = ref('')
const stats = reactive({ fajr: 0, dhuhr: 0, asr: 0, maghrib: 0, isha: 0 })
const prayers = [
  { key: 'fajr', label: t('prayer.fajr') }, { key: 'dhuhr', label: t('prayer.dhuhr') },
  { key: 'asr', label: t('prayer.asr') }, { key: 'maghrib', label: t('prayer.maghrib') }, { key: 'isha', label: t('prayer.isha') },
]
const bulkPrayers = reactive({ fajr: true, dhuhr: true, asr: true, maghrib: true, isha: true })

const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name' },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'fajr', label: 'prayer.fajr' }, { key: 'dhuhr', label: 'prayer.dhuhr' },
  { key: 'asr', label: 'prayer.asr' }, { key: 'maghrib', label: 'prayer.maghrib' }, { key: 'isha', label: 'prayer.isha' },
]

const prayerChart = computed(() => ({
  series: [{ name: t('islamic.compliance'), data: [stats.fajr, stats.dhuhr, stats.asr, stats.maghrib, stats.isha] }],
  options: { chart: { type: 'bar' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669'], plotOptions: { bar: { borderRadius: 4 } }, xaxis: { categories: [t('prayer.fajr'), t('prayer.dhuhr'), t('prayer.asr'), t('prayer.maghrib'), t('prayer.isha')] }, yaxis: { max: 100 }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))

const fetchData = async () => { if (!selectedClass.value) return; loading.value = true; try { const [rec, s] = await Promise.all([api.get('/prayer-attendance', { classId: selectedClass.value, date: selectedDate.value }), api.get<{fajr:number;dhuhr:number;asr:number;maghrib:number;isha:number}>('/prayer-attendance/stats', { classId: selectedClass.value })]); records.value = rec; Object.assign(stats, s) } catch {} finally { loading.value = false } }

const openBulkEntry = () => { bulkDate.value = selectedDate.value; bulkClass.value = selectedClass.value; showBulkDialog.value = true }
const submitBulk = async () => { bulkSaving.value = true; try { await api.post('/prayer-attendance/bulk', { date: bulkDate.value, classId: bulkClass.value, prayers: bulkPrayers }); toast.add({ title: t('attendance.bulk_saved'), color: 'success' }); showBulkDialog.value = false; fetchData() } catch {} finally { bulkSaving.value = false } }

const fetchClasses = async () => { try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); if (classOptions.value.length > 0) { selectedClass.value = classOptions.value[0].value; fetchData() } } catch {} }

onMounted(() => fetchClasses())
</script>
