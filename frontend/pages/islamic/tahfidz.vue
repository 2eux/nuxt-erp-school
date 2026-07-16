<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('tahfidz.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('tahfidz.subtitle') }}</p></div>
      <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAddGroup">{{ $t('tahfidz.add_group') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('tahfidz.total_groups')" :value="stats.totalGroups" icon="i-heroicons-user-group" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('tahfidz.total_students')" :value="stats.totalStudents" icon="i-heroicons-users" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('tahfidz.completed_juz')" :value="stats.completedJuz" icon="i-heroicons-check-badge" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('tahfidz.teachers')" :value="stats.totalTeachers" icon="i-heroicons-academic-cap" color="purple" :loading="statsLoading" />
    </div>

    <DataTable :columns="columns" :rows="groups" :loading="loading" :empty-title="$t('tahfidz.no_groups')" :show-export="false">
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewGroup(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editGroup(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('tahfidz.edit_group') : $t('tahfidz.add_group')" :loading="saving" @submit="saveGroup" @cancel="closeDialog">
      <div class="space-y-4">
        <UFormGroup :label="$t('tahfidz.group_name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.teacher')"><USelect v-model="form.teacherId" :options="teacherOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.description')"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.schedule')"><UInput v-model="form.schedule" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.room')"><UInput v-model="form.room" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.members')">
          <USelect v-model="form.studentIds" :options="studentOptions" color="gray" :multiple="true" />
        </UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('tahfidz.delete_group')" :loading="deleting" @confirm="deleteGroup" @cancel="showDelete=false" />
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
const showDialog = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const groups = ref<Record<string, unknown>[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)
const teacherOptions = ref<{ label: string; value: string }[]>([])
const studentOptions = ref<{ label: string; value: string }[]>([])
const stats = reactive({ totalGroups: 0, totalStudents: 0, completedJuz: 0, totalTeachers: 0 })

const columns: TableColumn[] = [
  { key: 'name', label: 'tahfidz.group_name', sortable: true },
  { key: 'teacherName', label: 'tahfidz.teacher' },
  { key: 'schedule', label: 'tahfidz.schedule' },
  { key: 'memberCount', label: 'tahfidz.members', type: 'number' },
  { key: 'room', label: 'tahfidz.room' },
]
const form = reactive({ name: '', teacherId: '', description: '', schedule: '', room: '', studentIds: [] as string[] })

const fetchGroups = async () => { loading.value = true; try { groups.value = await api.paginate('/tahfidz/groups').then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/tahfidz/stats')) } catch {} finally { statsLoading.value = false } }

const openAddGroup = () => { editing.value = false; editId.value = null; Object.assign(form, { name: '', teacherId: '', description: '', schedule: '', room: '', studentIds: [] }); showDialog.value = true }
const editGroup = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false }
const viewGroup = (row: Record<string, unknown>) => { navigateTo(`/islamic/tahfidz/${row.id}`) }

const saveGroup = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/tahfidz/groups/${editId.value}`, form); toast.add({ title: t('tahfidz.group_updated'), color: 'success' }) }
    else { await api.post('/tahfidz/groups', form); toast.add({ title: t('tahfidz.group_created'), color: 'success' }) }
    closeDialog(); fetchGroups(); fetchStats()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteGroup = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/tahfidz/groups/${deleteTarget.value.id}`); toast.add({ title: t('tahfidz.group_deleted'), color: 'success' }); showDelete.value = false; fetchGroups(); fetchStats() } catch {} finally { deleting.value = false } }

const fetchOptions = async () => {
  try { teacherOptions.value = (await api.get<{id:string;fullName:string}[]>('/teachers')).map(t => ({ label: t.fullName, value: t.id })) } catch {}
  try { studentOptions.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })) } catch {}
}

onMounted(() => { fetchGroups(); fetchStats(); fetchOptions() })
</script>
