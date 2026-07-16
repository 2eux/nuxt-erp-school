<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('meetings.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('meetings.subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <UButton v-if="permissions.can('meetings.manage')" color="outline" size="sm" icon="i-heroicons-calendar" @click="viewCalendar">{{ $t('meetings.calendar_view') }}</UButton>
        <UButton v-if="permissions.can('meetings.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('meetings.schedule_meeting') }}</UButton>
      </div>
    </div>
    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />
    <DataTable :columns="columns" :rows="meetings" :loading="loading" :empty-title="$t('meetings.no_meetings')" :show-export="false">
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'scheduled'" /></template>
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewMeeting(row as Record<string, unknown>)" />
        <UButton v-if="(row.status as string) === 'scheduled'" color="gray" variant="ghost" size="xs" icon="i-heroicons-clipboard-document" @click="addMinutes(row as Record<string, unknown>)" />
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('meetings.edit') : $t('meetings.schedule_meeting')" :loading="saving" @submit="saveMeeting" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('meetings.title')" required><UInput v-model="form.title" color="gray" /></UFormGroup>
        <div class="grid grid-cols-3 gap-4">
          <UFormGroup :label="$t('common.date')" required><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('meetings.start_time')"><UInput v-model="form.startTime" type="time" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('meetings.end_time')"><UInput v-model="form.endTime" type="time" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('meetings.location')"><UInput v-model="form.location" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('meetings.type')"><USelect v-model="form.type" :options="[{label:t('meetings.staff'),value:'staff'},{label:t('meetings.parent'),value:'parent'},{label:t('meetings.committee'),value:'committee'}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('meetings.participants')"><USelect v-model="form.participantIds" :options="userOptions" color="gray" :multiple="true" /></UFormGroup>
        <UFormGroup :label="$t('common.description')"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>

    <FormDialog v-model="showMinutes" :title="$t('meetings.minutes')" :loading="savingMinutes" @submit="saveMinutes" @cancel="showMinutes=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('meetings.minutes')" required><RichEditor v-model="minutesForm.content" /></UFormGroup>
        <UFormGroup :label="$t('common.status')"><USelect v-model="minutesForm.status" :options="[{label:t('status.scheduled'),value:'scheduled'},{label:t('status.completed'),value:'completed'}]" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('meetings.delete')" :loading="deleting" @confirm="deleteMeeting" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const saving = ref(false); const savingMinutes = ref(false); const deleting = ref(false)
const showForm = ref(false); const showMinutes = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const meetings = ref<Record<string, unknown>[]>([]); const minutesTarget = ref<Record<string, unknown> | null>(null)
const deleteTarget = ref<Record<string, unknown> | null>(null)
const userOptions = ref<{ label: string; value: string }[]>([])
const columns: TableColumn[] = [
  { key: 'title', label: 'meetings.title', sortable: true },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'startTime', label: 'meetings.start_time' },
  { key: 'endTime', label: 'meetings.end_time' },
  { key: 'type', label: 'meetings.type' },
  { key: 'location', label: 'meetings.location' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [{ key: 'type', label: 'meetings.type', type: 'select' as const, options: [{ label: t('meetings.staff'), value: 'staff' }, { label: t('meetings.parent'), value: 'parent' }] }]
const form = reactive({ title: '', description: '', date: '', startTime: '08:00', endTime: '09:00', location: '', type: 'staff', participantIds: [] as string[] })
const minutesForm = reactive({ content: '', status: 'completed' })

const fetchMeetings = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { meetings.value = await api.paginate('/meetings', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchMeetings(filters)
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { title: '', description: '', date: $dayjs().format('YYYY-MM-DD'), startTime: '08:00', endTime: '09:00', location: '', type: 'staff', participantIds: [] }); showForm.value = true }
const saveMeeting = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/meetings/${editId.value}`, form); toast.add({ title: t('meetings.updated'), color: 'success' }) } else { await api.post('/meetings', form); toast.add({ title: t('meetings.created'), color: 'success' }) } showForm.value = false; fetchMeetings() } catch {} finally { saving.value = false } }
const viewMeeting = (row: Record<string, unknown>) => { navigateTo(`/meetings/${row.id}`) }
const viewCalendar = () => { /* calendar view */ }
const addMinutes = (row: Record<string, unknown>) => { minutesTarget.value = row; minutesForm.content = (row.minutes as string) || ''; minutesForm.status = 'completed'; showMinutes.value = true }
const saveMinutes = async () => { savingMinutes.value = true; try { await api.patch(`/meetings/${minutesTarget.value?.id}/minutes`, minutesForm); toast.add({ title: t('meetings.minutes_saved'), color: 'success' }); showMinutes.value = false; fetchMeetings() } catch {} finally { savingMinutes.value = false } }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteMeeting = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/meetings/${deleteTarget.value.id}`); toast.add({ title: t('meetings.deleted'), color: 'success' }); showDelete.value = false; fetchMeetings() } catch {} finally { deleting.value = false } }
const fetchOptions = async () => { try { userOptions.value = (await api.get<{id:string;fullName:string}[]>('/users')).map(u => ({ label: u.fullName, value: u.id })) } catch {} }
onMounted(() => { fetchMeetings(); fetchOptions() })
</script>
