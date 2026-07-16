<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.budget') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.budget_subtitle') }}</p></div>
      <UButton v-if="permissions.can('finance.budget.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('finance.add_budget') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('finance.total_budget')" :value="formatCurrency(stats.totalBudget)" icon="i-heroicons-calculator" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('finance.total_actual')" :value="formatCurrency(stats.totalActual)" icon="i-heroicons-banknotes" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('finance.remaining')" :value="formatCurrency(stats.remaining)" icon="i-heroicons-wallet" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('finance.utilization')" :value="`${stats.utilization}%`" icon="i-heroicons-chart-bar" color="purple" :loading="statsLoading" />
    </div>

    <DataTable :columns="columns" :rows="budgets" :loading="loading" :empty-title="$t('finance.no_budgets')" :show-export="false">
      <template #cell-plannedAmount="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.plannedAmount as number) }}</span></template>
      <template #cell-actualAmount="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.actualAmount as number) }}</span></template>
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'draft'" /></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editBudget(row as Record<string, unknown>)" />
          <UButton v-if="(row.status as string) === 'draft'" color="emerald" variant="ghost" size="xs" icon="i-heroicons-check" @click="approveBudget(row)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('finance.edit_budget') : $t('finance.add_budget')" :loading="saving" @submit="saveBudget" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('finance.budget_name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('finance.category')"><UInput v-model="form.category" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('finance.planned_amount')" required><UInput v-model.number="form.plannedAmount" type="number" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('finance.actual_amount')"><UInput v-model.number="form.actualAmount" type="number" color="gray" /></UFormGroup>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('common.start_date')"><UInput v-model="form.startDate" type="date" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('common.end_date')"><UInput v-model="form.endDate" type="date" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('academic.academic_year')"><USelect v-model="form.academicYearId" :options="yearOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.description')"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('finance.delete_budget')" :loading="deleting" @confirm="deleteBudget" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const statsLoading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const budgets = ref<Record<string, unknown>[]>([])
const yearOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)

const stats = reactive({ totalBudget: 0, totalActual: 0, remaining: 0, utilization: 0 })
const columns: TableColumn[] = [
  { key: 'name', label: 'finance.budget_name', sortable: true },
  { key: 'category', label: 'finance.category' },
  { key: 'plannedAmount', label: 'finance.planned_amount', type: 'currency' },
  { key: 'actualAmount', label: 'finance.actual_amount', type: 'currency' },
  { key: 'startDate', label: 'common.start_date', type: 'date' },
  { key: 'endDate', label: 'common.end_date', type: 'date' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const form = reactive({ name: '', category: '', plannedAmount: 0, actualAmount: 0, startDate: '', endDate: '', academicYearId: '', description: '' })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

const fetchBudgets = async () => { loading.value = true; try { budgets.value = await api.paginate('/finance/budgets').then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/finance/budgets/stats')) } catch {} finally { statsLoading.value = false } }
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { name: '', category: '', plannedAmount: 0, actualAmount: 0, startDate: '', endDate: '', academicYearId: '', description: '' }); showForm.value = true }
const editBudget = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveBudget = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/finance/budgets/${editId.value}`, form); toast.add({ title: t('finance.budget_updated'), color: 'success' }) } else { await api.post('/finance/budgets', form); toast.add({ title: t('finance.budget_created'), color: 'success' }) } showForm.value = false; fetchBudgets(); fetchStats() } catch {} finally { saving.value = false } }
const approveBudget = async (row: Record<string, unknown>) => { try { await api.patch(`/finance/budgets/${row.id}/approve`); toast.add({ title: t('finance.budget_approved'), color: 'success' }); fetchBudgets() } catch {} }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteBudget = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/finance/budgets/${deleteTarget.value.id}`); toast.add({ title: t('finance.budget_deleted'), color: 'success' }); showDelete.value = false; fetchBudgets(); fetchStats() } catch {} finally { deleting.value = false } }
const fetchOptions = async () => { try { yearOptions.value = (await api.get<{id:string;name:string}[]>('/academic-years')).map(y => ({ label: y.name, value: y.id })) } catch {} }

onMounted(() => { fetchBudgets(); fetchStats(); fetchOptions() })
</script>
