<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.invoices') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.invoices_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <UButton v-if="permissions.can('finance.invoices.create')" color="gray" variant="outline" size="sm" icon="i-heroicons-arrow-path" @click="openBulkGenerate">{{ $t('finance.bulk_generate') }}</UButton>
        <UButton v-if="permissions.can('finance.invoices.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('finance.create_invoice') }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <StatCard :label="$t('finance.unpaid')" :value="stats.unpaid" icon="i-heroicons-exclamation-circle" color="red" :loading="statsLoading" />
      <StatCard :label="$t('finance.paid')" :value="stats.paid" icon="i-heroicons-check-circle" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('finance.overdue')" :value="stats.overdue" icon="i-heroicons-x-circle" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('finance.total')" :value="formatCurrency(stats.total)" icon="i-heroicons-currency-dollar" color="blue" :loading="statsLoading" />
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="invoices" :loading="loading" :empty-title="$t('finance.no_invoices')" :show-export="true" @export="handleExport">
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'unpaid'" /></template>
      <template #cell-amount="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.amount as number) }}</span></template>
      <template #cell-paidAmount="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.paidAmount as number) }}</span></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewInvoice(row as Record<string, unknown>)" />
          <UButton v-if="(row.status as string) !== 'paid'" color="gray" variant="ghost" size="xs" icon="i-heroicons-envelope" @click="sendReminder(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="$t('finance.create_invoice')" :loading="saving" @submit="saveInvoice" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('tahfidz.student')" required><USelect v-model="form.studentId" :options="studentOptions" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('finance.type')" required><USelect v-model="form.type" :options="typeOptions" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('finance.amount')" required><UInput v-model.number="form.amount" type="number" color="gray" /></UFormGroup>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('finance.month')"><USelect v-model="form.month" :options="monthOptions" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('finance.year')"><UInput v-model.number="form.year" type="number" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('finance.due_date')" required><UInput v-model="form.dueDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.description')"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>

    <FormDialog v-model="showBulkForm" :title="$t('finance.bulk_generate')" :loading="bulkSaving" @submit="submitBulk" @cancel="showBulkForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('finance.type')"><USelect v-model="bulkForm.type" :options="typeOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.class')"><USelect v-model="bulkForm.classId" :options="classOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('finance.amount')"><UInput v-model.number="bulkForm.amount" type="number" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('finance.month')"><USelect v-model="bulkForm.month" :options="monthOptions" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('finance.year')"><UInput v-model.number="bulkForm.year" type="number" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('finance.due_date')"><UInput v-model="bulkForm.dueDate" type="date" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('finance.delete_invoice')" :loading="deleting" @confirm="deleteInvoice" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const statsLoading = ref(false); const saving = ref(false); const bulkSaving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showBulkForm = ref(false); const showDelete = ref(false)
const invoices = ref<Record<string, unknown>[]>([])
const studentOptions = ref<{ label: string; value: string }[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)
const stats = reactive({ unpaid: 0, paid: 0, overdue: 0, total: 0 })

const columns: TableColumn[] = [
  { key: 'invoiceNumber', label: 'finance.invoice_number', sortable: true },
  { key: 'studentName', label: 'students.name' },
  { key: 'type', label: 'finance.type' },
  { key: 'amount', label: 'finance.amount', type: 'currency' },
  { key: 'paidAmount', label: 'finance.paid', type: 'currency' },
  { key: 'dueDate', label: 'finance.due_date', type: 'date' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.unpaid'), value: 'unpaid' }, { label: t('status.paid'), value: 'paid' }, { label: t('status.overdue'), value: 'overdue' }] },
  { key: 'classId', label: 'academic.class', type: 'select' as const, options: [] },
]
const typeOptions = [{ label: 'SPP', value: 'spp' }, { label: t('finance.registration'), value: 'registration' }, { label: t('finance.exam'), value: 'exam' }, { label: t('finance.book'), value: 'book' }, { label: t('finance.uniform'), value: 'uniform' }, { label: t('finance.activity'), value: 'activity' }]
const monthOptions = Array.from({ length: 12 }, (_, i) => ({ label: $dayjs().month(i).format('MMMM'), value: String(i + 1) }))
const form = reactive({ studentId: '', type: 'spp', amount: 0, month: String($dayjs().month() + 1), year: $dayjs().year(), dueDate: '', description: '' })
const bulkForm = reactive({ type: 'spp', classId: '', amount: 0, month: String($dayjs().month() + 1), year: $dayjs().year(), dueDate: '' })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

const fetchInvoices = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { invoices.value = await api.paginate('/finance/invoices', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/finance/invoices/stats')) } catch {} finally { statsLoading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchInvoices(filters)

const openAdd = () => { Object.assign(form, { studentId: '', type: 'spp', amount: 0, month: String($dayjs().month() + 1), year: $dayjs().year(), dueDate: '', description: '' }); showForm.value = true }
const openBulkGenerate = () => { Object.assign(bulkForm, { type: 'spp', classId: '', amount: 0, month: String($dayjs().month() + 1), year: $dayjs().year(), dueDate: '' }); showBulkForm.value = true }

const saveInvoice = async () => { saving.value = true; try { await api.post('/finance/invoices', form); toast.add({ title: t('finance.invoice_created'), color: 'success' }); showForm.value = false; fetchInvoices(); fetchStats() } catch {} finally { saving.value = false } }
const submitBulk = async () => { bulkSaving.value = true; try { await api.post('/finance/invoices/bulk', bulkForm); toast.add({ title: t('finance.invoices_generated'), color: 'success' }); showBulkForm.value = false; fetchInvoices(); fetchStats() } catch {} finally { bulkSaving.value = false } }

const viewInvoice = (row: Record<string, unknown>) => { window.open(`/api/v1/finance/invoices/${row.id}/view`, '_blank') }
const sendReminder = async (row: Record<string, unknown>) => { try { await api.post(`/finance/invoices/${row.id}/remind`); toast.add({ title: t('finance.reminder_sent'), color: 'success' }) } catch {} }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteInvoice = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/finance/invoices/${deleteTarget.value.id}`); toast.add({ title: t('finance.invoice_deleted'), color: 'success' }); showDelete.value = false; fetchInvoices(); fetchStats() } catch {} finally { deleting.value = false } }
const handleExport = (format: string) => { window.open(`/api/v1/finance/invoices/export?format=${format}`, '_blank') }

const fetchOptions = async () => { try { studentOptions.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })) } catch {}; try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); filterFields[1].options = classOptions.value } catch {} }

onMounted(() => { fetchInvoices(); fetchStats(); fetchOptions() })
</script>
