<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('employees.title') }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('employees.subtitle') }}</p>
      </div>
      <UButton v-if="permissions.can('employees.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAddDialog">
        {{ $t('employees.add_employee') }}
      </UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('employees.total')" :value="stats.total" icon="i-heroicons-briefcase" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('employees.active')" :value="stats.active" icon="i-heroicons-check-circle" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('employees.by_department')" :value="stats.departments" icon="i-heroicons-building-office" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('employees.on_leave')" :value="stats.onLeave" icon="i-heroicons-clock" color="purple" :loading="statsLoading" />
    </div>

    <DataFilter :filter-fields="filterFields" :searchable="true" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="employees" :loading="loading" :empty-title="$t('employees.no_employees')" :show-export="false">
      <template #cell-fullName="{ row }">
        <div class="flex items-center gap-3">
          <UAvatar :src="(row.photo as string) || undefined" size="sm" />
          <div><span class="text-sm font-medium text-gray-900 dark:text-white">{{ row.fullName }}</span><p class="text-xs text-gray-500">{{ row.nip }}</p></div>
        </div>
      </template>
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'active'" /></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewEmployee(row)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editEmployee(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('employees.edit_employee') : $t('employees.add_employee')" :loading="saving" @submit="saveEmployee" @cancel="closeDialog">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('employees.nip')" required><UInput v-model="form.nip" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.full_name')" required class="sm:col-span-2"><UInput v-model="form.fullName" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.gender')"><USelect v-model="form.gender" :options="[{label:t('common.male'),value:'male'},{label:t('common.female'),value:'female'}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.position')"><UInput v-model="form.position" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.department')"><USelect v-model="form.department" :options="departmentOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.phone')"><UInput v-model="form.phone" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.email')"><UInput v-model="form.email" type="email" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.join_date')"><UInput v-model="form.joinDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.status')"><USelect v-model="form.status" :options="[{label:t('status.active'),value:'active'},{label:t('status.inactive'),value:'inactive'},{label:t('status.resigned'),value:'resigned'}]" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDeleteConfirm" :title="$t('employees.delete_employee')" :loading="deleting" @confirm="deleteEmployee" @cancel="showDeleteConfirm=false" />
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
const employees = ref<Record<string, unknown>[]>([])
const editing = ref(false)
const editId = ref<string | null>(null)
const showDialog = ref(false)
const showDeleteConfirm = ref(false)
const deleteTarget = ref<Record<string, unknown> | null>(null)
const stats = reactive({ total: 0, active: 0, departments: 0, onLeave: 0 })

const columns: TableColumn[] = [
  { key: 'fullName', label: 'employees.name', sortable: true },
  { key: 'nip', label: 'employees.nip' },
  { key: 'position', label: 'employees.position' },
  { key: 'department', label: 'employees.department' },
  { key: 'phone', label: 'common.phone' },
  { key: 'joinDate', label: 'employees.join_date', type: 'date' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.active'), value: 'active' }, { label: t('status.inactive'), value: 'inactive' }] },
  { key: 'department', label: 'employees.department', type: 'text' as const },
]
const departmentOptions = [
  { label: 'Administrasi', value: 'Administrasi' }, { label: 'Keuangan', value: 'Keuangan' },
  { label: 'Kesiswaan', value: 'Kesiswaan' }, { label: 'TU', value: 'TU' },
  { label: 'Kebersihan', value: 'Kebersihan' }, { label: 'Keamanan', value: 'Keamanan' },
  { label: 'Dapur', value: 'Dapur' }, { label: 'Perpustakaan', value: 'Perpustakaan' },
]

const form = reactive({ nip: '', fullName: '', gender: 'male', position: '', department: '', phone: '', email: '', joinDate: '', status: 'active' })

const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get<{total:number;active:number;departments:number;onLeave:number}>('/employees/stats')) } catch { /* ignore */ } finally { statsLoading.value = false } }
const fetchEmployees = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { employees.value = await api.paginate('/employees', { ...filters }).then(r => r.data) } catch { /* ignore */ } finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchEmployees(filters)

const openAddDialog = () => { editing.value = false; editId.value = null; Object.assign(form, { nip: '', fullName: '', gender: 'male', position: '', department: '', phone: '', email: '', joinDate: '', status: 'active' }); showDialog.value = true }
const editEmployee = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false; editId.value = null }
const viewEmployee = (row: Record<string, unknown>) => navigateTo(`/employees/${row.id}`)

const saveEmployee = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/employees/${editId.value}`, form); toast.add({ title: t('employees.employee_updated'), color: 'success' }) }
    else { await api.post('/employees', form); toast.add({ title: t('employees.employee_created'), color: 'success' }) }
    closeDialog(); fetchEmployees(); fetchStats()
  } catch {/* handled */} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDeleteConfirm.value = true }
const deleteEmployee = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/employees/${deleteTarget.value.id}`); toast.add({ title: t('employees.employee_deleted'), color: 'success' }); showDeleteConfirm.value = false; fetchEmployees(); fetchStats() } catch {/* handled */} finally { deleting.value = false } }

onMounted(() => { fetchStats(); fetchEmployees() })
</script>
