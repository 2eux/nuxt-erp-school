<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('islamic.islamic_events') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('islamic.islamic_events_subtitle') }}</p></div>
      <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('islamic.add_event') }}</UButton>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="card lg:col-span-2">
        <div class="card-header"><h3 class="card-title">{{ $t('islamic.upcoming_events') }}</h3></div>
        <div class="space-y-3">
          <div v-for="event in upcomingEvents" :key="event.id" class="flex items-center gap-4 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50">
            <div class="w-14 h-14 rounded-xl bg-emerald-100 dark:bg-emerald-900/30 flex flex-col items-center justify-center shrink-0">
              <span class="text-sm font-bold text-emerald-600 dark:text-emerald-400">{{ $dayjs(event.date).format('DD') }}</span>
              <span class="text-xs text-emerald-500">{{ $dayjs(event.date).format('MMM') }}</span>
            </div>
            <div class="flex-1"><p class="text-sm font-medium text-gray-900 dark:text-white">{{ event.name }}</p><p class="text-xs text-gray-500">{{ event.description }}</p></div>
            <StatusBadge :status="event.type" />
          </div>
        </div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('islamic.hijri_calendar') }}</h3></div>
        <div class="text-center py-4">
          <p class="text-3xl font-bold text-brand-600 dark:text-brand-400">{{ hijriDate }}</p>
          <p class="text-sm text-gray-500 mt-1">{{ $t('islamic.hijri_date') }}</p>
        </div>
      </div>
    </div>

    <DataTable :columns="columns" :rows="events" :loading="loading" :empty-title="$t('islamic.no_events')" :show-export="false">
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editEvent(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('islamic.edit_event') : $t('islamic.add_event')" :loading="saving" @submit="saveEvent" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('islamic.event_name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.date')" required><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('islamic.event_type')"><USelect v-model="form.type" :options="[{label:t('islamic.hijri'),value:'hijri'},{label:t('islamic.islamic_holiday'),value:'islamic_holiday'},{label:t('islamic.school_islamic'),value:'school_islamic'}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.description')"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('islamic.delete_event')" :loading="deleting" @confirm="deleteEvent" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const events = ref<Record<string, unknown>[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)

const columns: TableColumn[] = [
  { key: 'name', label: 'islamic.event_name', sortable: true },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'type', label: 'islamic.event_type' },
  { key: 'description', label: 'common.description' },
]
const form = reactive({ name: '', date: '', type: 'islamic_holiday', description: '' })

const upcomingEvents = computed(() => events.value.filter(e => $dayjs(e.date as string).isAfter($dayjs())).slice(0, 5))
const hijriDate = computed(() => { return 'Muharram 1447 H' })

const fetchEvents = async () => { loading.value = true; try { events.value = await api.paginate('/islamic-events').then(r => r.data) } catch {} finally { loading.value = false } }
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { name: '', date: '', type: 'islamic_holiday', description: '' }); showForm.value = true }
const editEvent = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveEvent = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/islamic-events/${editId.value}`, form); toast.add({ title: t('islamic.event_updated'), color: 'success' }) } else { await api.post('/islamic-events', form); toast.add({ title: t('islamic.event_created'), color: 'success' }) } showForm.value = false; fetchEvents() } catch {} finally { saving.value = false } }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteEvent = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/islamic-events/${deleteTarget.value.id}`); toast.add({ title: t('islamic.event_deleted'), color: 'success' }); showDelete.value = false; fetchEvents() } catch {} finally { deleting.value = false } }
onMounted(() => fetchEvents())
</script>
