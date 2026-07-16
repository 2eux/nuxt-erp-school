<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.schedules') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.schedules_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedClass" :options="classOptions" :placeholder="$t('academic.select_class')" color="gray" size="sm" class="w-40" @change="fetchSchedule" />
        <UButton color="gray" variant="outline" size="sm" icon="i-heroicons-printer" @click="printSchedule">{{ $t('common.print') }}</UButton>
        <UButton v-if="permissions.can('academic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('academic.add_schedule') }}</UButton>
      </div>
    </div>

    <div class="card overflow-hidden !p-0">
      <div class="overflow-x-auto">
        <table class="w-full" id="schedule-table">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-800/50">
              <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">
                {{ $t('academic.time') }}
              </th>
              <th v-for="day in days" :key="day.key" class="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase border-l border-gray-200 dark:border-gray-700">
                {{ day.label }}
              </th>
            </tr>
          </thead>
          <tbody>
            <template v-if="loading">
              <tr v-for="i in 8" :key="'sk-'+i">
                <td class="px-4 py-3"><div class="h-4 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-20" /></td>
                <td v-for="d in 6" :key="d" class="px-4 py-3 border-l border-gray-100 dark:border-gray-750"><div class="h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" /></td>
              </tr>
            </template>
            <template v-else-if="scheduleData.length === 0">
              <tr><td :colspan="7" class="px-6 py-16 text-center">
                <EmptyState :title="$t('academic.no_schedule')" icon="i-heroicons-clock" />
              </td></tr>
            </template>
            <template v-else>
              <tr v-for="timeSlot in timeSlots" :key="timeSlot">
                <td class="px-4 py-2 text-sm font-medium text-gray-900 dark:text-white whitespace-nowrap border-t border-gray-100 dark:border-gray-750">
                  {{ timeSlot }}
                </td>
                <td v-for="day in days" :key="day.key" class="px-3 py-2 border-l border-t border-gray-100 dark:border-gray-750 min-w-[140px]">
                  <div v-for="item in getSlotItems(day.key, timeSlot)" :key="item.id" class="mb-1 p-2 rounded-lg text-xs"
                    :class="item.bgClass">
                    <p class="font-semibold">{{ item.subjectName }}</p>
                    <p class="opacity-75">{{ item.room || '-' }}</p>
                    <p class="opacity-75">{{ item.teacherName || '-' }}</p>
                    <p class="opacity-60">{{ item.startTime }} - {{ item.endTime }}</p>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <FormDialog v-model="showDialog" :title="$t('academic.add_schedule')" :loading="saving" @submit="saveSchedule" @cancel="showDialog=false">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('academic.class')" required><USelect v-model="form.classId" :options="classOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('subjects.title')" required><USelect v-model="form.subjectId" :options="subjectOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.day')" required><USelect v-model="form.day" :options="dayOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.teacher')"><USelect v-model="form.teacherId" :options="teacherOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.start_time')" required><UInput v-model="form.startTime" type="time" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.end_time')" required><UInput v-model="form.endTime" type="time" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.semester')"><USelect v-model="form.semester" :options="[{label:'1',value:1},{label:'2',value:2}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.room')"><UInput v-model="form.room" color="gray" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { Schedule } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const { $dayjs } = useNuxtApp()
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const selectedClass = ref('')
const scheduleData = ref<Schedule[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const subjectOptions = ref<{ label: string; value: string }[]>([])
const teacherOptions = ref<{ label: string; value: string }[]>([])

const days = [
  { key: 'monday', label: t('day.monday') }, { key: 'tuesday', label: t('day.tuesday') },
  { key: 'wednesday', label: t('day.wednesday') }, { key: 'thursday', label: t('day.thursday') },
  { key: 'friday', label: t('day.friday') }, { key: 'saturday', label: t('day.saturday') },
]
const dayOpts = days.map(d => ({ label: d.label, value: d.key }))

const timeSlots = computed(() => {
  const times = ['07:00', '07:45', '08:30', '09:15', '10:00', '10:45', '11:30', '12:15', '13:00', '13:45', '14:30', '15:15']
  return times.slice(0, 8)
})

const getSlotItems = (day: string, timeStr: string) => {
  const items = scheduleData.value.filter(s => s.day === day && s.startTime <= timeStr && s.endTime > timeStr)
  const colorMap = ['bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300',
    'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-800 dark:text-emerald-300',
    'bg-amber-100 dark:bg-amber-900/30 text-amber-800 dark:text-amber-300',
    'bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-300',
    'bg-rose-100 dark:bg-rose-900/30 text-rose-800 dark:text-rose-300',
    'bg-cyan-100 dark:bg-cyan-900/30 text-cyan-800 dark:text-cyan-300']
  return items.map((item, i) => ({ ...(item as unknown as Record<string, unknown>), id: item.id, bgClass: colorMap[i % colorMap.length] }))
}

const form = reactive({ classId: '', subjectId: '', day: 'monday', startTime: '07:00', endTime: '07:45', teacherId: '', room: '', semester: 1 })

const fetchSchedule = async () => {
  if (!selectedClass.value) return
  loading.value = true
  try { scheduleData.value = await api.get<Schedule[]>('/schedules', { classId: selectedClass.value }) }
  catch {} finally { loading.value = false }
}

const openAdd = () => { Object.assign(form, { classId: selectedClass.value, subjectId: '', day: 'monday', startTime: '07:00', endTime: '07:45', teacherId: '', room: '', semester: 1 }); showDialog.value = true }

const saveSchedule = async () => {
  saving.value = true
  try { await api.post('/schedules', form); toast.add({ title: t('academic.schedule_created'), color: 'success' }); showDialog.value = false; fetchSchedule() }
  catch {} finally { saving.value = false }
}

const printSchedule = () => { window.print() }

const fetchOptions = async () => {
  try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); if (classOptions.value.length > 0) { selectedClass.value = classOptions.value[0].value; fetchSchedule() } } catch {}
  try { subjectOptions.value = (await api.get<{id:string;name:string}[]>('/subjects')).map(s => ({ label: s.name, value: s.id })) } catch {}
  try { teacherOptions.value = (await api.get<{id:string;fullName:string}[]>('/teachers', { limit: 100 })).map(t => ({ label: t.fullName, value: t.id })) } catch {}
}

onMounted(() => fetchOptions())
</script>
