<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.payments') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.payments_subtitle') }}</p></div>
      <UButton v-if="permissions.can('finance.payments.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('finance.record_payment') }}</UButton>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="payments" :loading="loading" :empty-title="$t('finance.no_payments')" :show-export="true" @export="handleExport">
      <template #cell-amount="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.amount as number) }}</span></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewDetail(row as Record<string, unknown>)" />
          <UButton v-if="(row.status as string) === 'pending'" color="emerald" variant="ghost" size="xs" icon="i-heroicons-check-circle" @click="verifyPayment(row)" />
          <UButton v-if="(row.status as string) === 'pending'" color="red" variant="ghost" size="xs" icon="i-heroicons-x-circle" @click="rejectPayment(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="$t('finance.record_payment')" :loading="saving" @submit="savePayment" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('finance.invoice')" required><USelect v-model="form.invoiceId" :options="invoiceOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('finance.amount')" required><UInput v-model.number="form.amount" type="number" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('finance.payment_method')"><USelect v-model="form.paymentMethod" :options="[{label:'Cash',value:'cash'},{label:'Bank Transfer',value:'bank_transfer'},{label:'Digital Wallet',value:'digital_wallet'}]" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('finance.payment_date')"><UInput v-model="form.paymentDate" type="date" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('finance.reference_number')"><UInput v-model="form.referenceNumber" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.note')"><UTextarea v-model="form.note" color="gray" :rows="2" /></UFormGroup>
        <FileUpload accept=".jpg,.jpeg,.png,.pdf" :multiple="false" accept-hint="Upload payment proof" @files-selected="handleProofSelected" />
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const saving = ref(false)
const showForm = ref(false)
const payments = ref<Record<string, unknown>[]>([])
const invoiceOptions = ref<{ label: string; value: string }[]>([])
const proofFile = ref<File | null>(null)

const columns: TableColumn[] = [
  { key: 'invoiceNumber', label: 'finance.invoice_number' },
  { key: 'studentName', label: 'students.name' },
  { key: 'amount', label: 'finance.amount', type: 'currency' },
  { key: 'paymentMethod', label: 'finance.payment_method' },
  { key: 'paymentDate', label: 'finance.payment_date', type: 'date' },
  { key: 'referenceNumber', label: 'finance.reference_number' },
  { key: 'recordedByName', label: 'finance.recorded_by' },
]
const filterFields = [
  { key: 'paymentMethod', label: 'finance.payment_method', type: 'select' as const, options: [{ label: 'Cash', value: 'cash' }, { label: 'Bank Transfer', value: 'bank_transfer' }, { label: 'Digital Wallet', value: 'digital_wallet' }] },
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.pending'), value: 'pending' }, { label: t('status.verified'), value: 'verified' }, { label: t('status.rejected'), value: 'rejected' }] },
]
const form = reactive({ invoiceId: '', amount: 0, paymentMethod: 'cash', paymentDate: $dayjs().format('YYYY-MM-DD'), referenceNumber: '', note: '' })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)
const handleProofSelected = (files: File[]) => { if (files.length > 0) proofFile.value = files[0] }

const fetchPayments = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { payments.value = await api.paginate('/finance/payments', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchPayments(filters)
const openAdd = () => { Object.assign(form, { invoiceId: '', amount: 0, paymentMethod: 'cash', paymentDate: $dayjs().format('YYYY-MM-DD'), referenceNumber: '', note: '' }); showForm.value = true }

const savePayment = async () => {
  saving.value = true
  try {
    const formData = new FormData(); formData.append('data', JSON.stringify(form))
    if (proofFile.value) formData.append('proof', proofFile.value)
    await api.upload('/finance/payments', formData)
    toast.add({ title: t('finance.payment_recorded'), color: 'success' }); showForm.value = false; fetchPayments()
  } catch {} finally { saving.value = false }
}

const viewDetail = (row: Record<string, unknown>) => { window.open(`/api/v1/finance/payments/${row.id}`, '_blank') }
const verifyPayment = async (row: Record<string, unknown>) => { try { await api.patch(`/finance/payments/${row.id}/verify`); toast.add({ title: t('finance.payment_verified'), color: 'success' }); fetchPayments() } catch {} }
const rejectPayment = async (row: Record<string, unknown>) => { try { await api.patch(`/finance/payments/${row.id}/reject`); toast.add({ title: t('finance.payment_rejected'), color: 'info' }); fetchPayments() } catch {} }
const handleExport = (format: string) => { window.open(`/api/v1/finance/payments/export?format=${format}`, '_blank') }

const fetchInvoices = async () => { try { invoiceOptions.value = (await api.get<{id:string;invoiceNumber:string;studentName:string}[]>('/finance/invoices', { status: 'unpaid' })).map(i => ({ label: `${i.invoiceNumber} - ${i.studentName}`, value: i.id })) } catch {} }

onMounted(() => { fetchPayments(); fetchInvoices() })
</script>
