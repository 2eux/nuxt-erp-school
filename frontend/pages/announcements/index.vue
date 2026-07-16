<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('announcements.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('announcements.subtitle') }}</p></div>
      <UButton v-if="permissions.can('announcements.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('announcements.create') }}</UButton>
    </div>
    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />
    <div class="space-y-3">
      <div v-for="a in announcements" :key="a.id" class="card hover:shadow-md transition-shadow cursor-pointer" @click="viewAnnouncement(a)">
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-full flex items-center justify-center shrink-0" :class="typeClasses(a.type as string)">
            <UIcon :name="typeIcons[a.type as string] || 'i-heroicons-megaphone'" class="w-5 h-5" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ a.title }}</h3>
              <StatusBadge :status="a.type as string" />
              <span v-if="a.isPinned" class="text-xs text-amber-600 bg-amber-100 dark:bg-amber-900/30 px-1.5 py-0.5 rounded">{{ $t('announcements.pinned') }}</span>
            </div>
            <p class="text-xs text-gray-500 mt-1 line-clamp-2" v-html="a.content"></p>
            <div class="flex items-center gap-3 mt-2 text-xs text-gray-400">
              <span>{{ a.createdByName }}</span>
              <span>{{ $dayjs(a.publishDate).format('DD MMM YYYY HH:mm') }}</span>
              <span>{{ a.targetAudience?.join(', ') }}</span>
            </div>
          </div>
          <div class="flex items-center gap-1">
            <UButton v-if="permissions.can('announcements.manage')" color="gray" variant="ghost" size="xs" icon="i-heroicons-pin" @click.stop="togglePin(a)" />
            <UButton v-if="permissions.can('announcements.manage')" color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click.stop="editAnnouncement(a as Record<string, unknown>)" />
            <UButton v-if="permissions.can('announcements.manage')" color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click.stop="confirmDelete(a)" />
          </div>
        </div>
      </div>
    </div>
    <EmptyState v-if="!loading && announcements.length === 0" :title="$t('announcements.no_announcements')" icon="i-heroicons-megaphone" :action-label="$t('announcements.create')" @action="openAdd" />

    <FormDialog v-model="showForm" :title="editing ? $t('announcements.edit') : $t('announcements.create')" :loading="saving" @submit="saveAnnouncement" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('announcements.title')" required><UInput v-model="form.title" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('announcements.content')" required><RichEditor v-model="form.content" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('announcements.type')"><USelect v-model="form.type" :options="[{label:t('announcements.general'),value:'general'},{label:t('announcements.academic'),value:'academic'},{label:t('announcements.event'),value:'event'},{label:t('announcements.urgent'),value:'urgent'},{label:t('announcements.finance'),value:'finance'}]" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('announcements.publish_date')"><UInput v-model="form.publishDate" type="datetime-local" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('announcements.target_audience')">
          <div class="flex flex-wrap gap-2">
            <UButton v-for="tgt in audienceOptions" :key="tgt.value" :color="form.targetAudience.includes(tgt.value) ? 'primary' : 'gray'" variant="outline" size="xs" @click="toggleAudience(tgt.value)">{{ tgt.label }}</UButton>
          </div>
        </UFormGroup>
        <UFormGroup :label="$t('announcements.expiry_date')"><UInput v-model="form.expiryDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('announcements.pin')"><UToggle v-model="form.isPinned" /></UFormGroup>
      </div>
    </FormDialog>
    <ConfirmDialog v-model="showDelete" :title="$t('announcements.delete')" :loading="deleting" @confirm="deleteAnnouncement" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const announcements = ref<Record<string, unknown>[]>([]); const deleteTarget = ref<Record<string, unknown> | null>(null)

const filterFields = [{ key: 'type', label: 'announcements.type', type: 'select' as const, options: [{ label: t('announcements.general'), value: 'general' }, { label: t('announcements.urgent'), value: 'urgent' }, { label: t('announcements.academic'), value: 'academic' }] }]
const audienceOptions = [{ label: t('common.all'), value: 'all' }, { label: t('teachers.title'), value: 'teachers' }, { label: t('students.title'), value: 'students' }, { label: t('parents.title'), value: 'parents' }, { label: t('employees.title'), value: 'employees' }]
const typeIcons: Record<string, string> = { urgent: 'i-heroicons-exclamation-triangle', academic: 'i-heroicons-academic-cap', event: 'i-heroicons-calendar-days', finance: 'i-heroicons-currency-dollar', general: 'i-heroicons-megaphone' }
const typeClasses = (t: string) => t === 'urgent' ? 'bg-red-100 dark:bg-red-900/30 text-red-600' : t === 'academic' ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600' : t === 'event' ? 'bg-amber-100 dark:bg-amber-900/30 text-amber-600' : 'bg-gray-100 dark:bg-gray-700 text-gray-600'
const form = reactive({ title: '', content: '', type: 'general', publishDate: '', expiryDate: '', targetAudience: ['all'] as string[], isPinned: false })

const fetchAnnouncements = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { announcements.value = await api.paginate('/announcements', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchAnnouncements(filters)
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { title: '', content: '', type: 'general', publishDate: $dayjs().format('YYYY-MM-DDTHH:mm'), expiryDate: '', targetAudience: ['all'], isPinned: false }); showForm.value = true }
const editAnnouncement = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); form.targetAudience = (row.targetAudience as string[]) || ['all']; showForm.value = true }
const toggleAudience = (val: string) => { const idx = form.targetAudience.indexOf(val); if (idx === -1) form.targetAudience.push(val); else form.targetAudience.splice(idx, 1) }
const saveAnnouncement = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/announcements/${editId.value}`, form); toast.add({ title: t('announcements.updated'), color: 'success' }) } else { await api.post('/announcements', form); toast.add({ title: t('announcements.created'), color: 'success' }) } showForm.value = false; fetchAnnouncements() } catch {} finally { saving.value = false } }
const togglePin = async (row: Record<string, unknown>) => { try { await api.patch(`/announcements/${row.id}/pin`, { isPinned: !row.isPinned }); fetchAnnouncements() } catch {} }
const viewAnnouncement = (row: Record<string, unknown>) => { navigateTo(`/announcements/${row.id}`) }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteAnnouncement = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/announcements/${deleteTarget.value.id}`); toast.add({ title: t('announcements.deleted'), color: 'success' }); showDelete.value = false; fetchAnnouncements() } catch {} finally { deleting.value = false } }
onMounted(() => fetchAnnouncements())
</script>
