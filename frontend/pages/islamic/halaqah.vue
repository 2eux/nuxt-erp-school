<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('halaqah.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('halaqah.subtitle') }}</p></div>
      <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('halaqah.add_group') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <StatCard :label="$t('halaqah.total_groups')" :value="stats.total" icon="i-heroicons-user-group" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('halaqah.total_members')" :value="stats.totalMembers" icon="i-heroicons-users" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('halaqah.sessions_week')" :value="stats.sessionsPerWeek" icon="i-heroicons-calendar" color="amber" :loading="statsLoading" />
    </div>

    <DataTable :columns="columns" :rows="groups" :loading="loading" :empty-title="$t('halaqah.no_groups')" :show-export="false">
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-users" @click="viewMembers(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editGroup(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('halaqah.edit_group') : $t('halaqah.add_group')" :loading="saving" @submit="saveGroup" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('halaqah.group_name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.teacher')"><USelect v-model="form.teacherId" :options="teacherOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('halaqah.description')"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.schedule')"><UInput v-model="form.schedule" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.room')"><UInput v-model="form.room" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('halaqah.members')"><USelect v-model="form.studentIds" :options="studentOptions" color="gray" :multiple="true" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('halaqah.delete_group')" :loading="deleting" @confirm="deleteGroup" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const loading = ref(false)
const statsLoading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const showForm = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const groups = ref<Record<string, unknown>[]>([])
const teacherOptions = ref<{ label: string; value: string }[]>([])
const studentOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)
const stats = reactive({ total: 0, totalMembers: 0, sessionsPerWeek: 0 })

const columns: TableColumn[] = [
  { key: 'name', label: 'halaqah.group_name', sortable: true },
  { key: 'teacherName', label: 'tahfidz.teacher' },
  { key: 'schedule', label: 'tahfidz.schedule' },
  { key: 'memberCount', label: 'halaqah.members', type: 'number' },
  { key: 'room', label: 'tahfidz.room' },
]
const form = reactive({ name: '', teacherId: '', description: '', schedule: '', room: '', studentIds: [] as string[] })

const fetchGroups = async () => { loading.value = true; try { groups.value = await api.paginate('/halaqah').then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/halaqah/stats')) } catch {} finally { statsLoading.value = false } }
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { name: '', teacherId: '', description: '', schedule: '', room: '', studentIds: [] }); showForm.value = true }
const editGroup = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const viewMembers = (row: Record<string, unknown>) => { navigateTo(`/islamic/halaqah/${row.id}`) }
const saveGroup = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/halaqah/${editId.value}`, form); toast.add({ title: t('halaqah.group_updated'), color: 'success' }) } else { await api.post('/halaqah', form); toast.add({ title: t('halaqah.group_created'), color: 'success' }) } showForm.value = false; fetchGroups(); fetchStats() } catch {} finally { saving.value = false } }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteGroup = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/halaqah/${deleteTarget.value.id}`); toast.add({ title: t('halaqah.group_deleted'), color: 'success' }); showDelete.value = false; fetchGroups(); fetchStats() } catch {} finally { deleting.value = false } }
const fetchOptions = async () => { try { teacherOptions.value = (await api.get<{id:string;fullName:string}[]>('/teachers')).map(t => ({ label: t.fullName, value: t.id })) } catch {} try { studentOptions.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })) } catch {} }

onMounted(() => { fetchGroups(); fetchStats(); fetchOptions() })
</script>
