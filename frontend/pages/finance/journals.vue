<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.journals') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.journals_subtitle') }}</p></div>
      <UButton v-if="permissions.can('finance.journals.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('finance.create_journal') }}</UButton>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="journals" :loading="loading" :empty-title="$t('finance.no_journals')" :show-export="true" @export="handleExport">
      <template #cell-type="{ row }"><StatusBadge :status="(row.type as string) || 'income'" /></template>
      <template #cell-amount="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.amount as number) }}</span></template>
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'draft'" /></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editJournal(row as Record<string, unknown>)" />
          <UButton v-if="(row.status as string) === 'draft'" color="emerald" variant="ghost" size="xs" icon="i-heroicons-check" @click="postJournal(row)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('finance.edit_journal') : $t('finance.create_journal')" :loading="saving" @submit="saveJournal" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('finance.journal_number')"><UInput v-model="form.journalNumber" color="gray" readonly /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('common.date')" required><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('finance.type')" required><USelect v-model="form.type" :options="[{label:t('finance.income'),value:'income'},{label:t('finance.expense'),value:'expense'},{label:t('finance.transfer'),value:'transfer'}]" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('finance.category')"><UInput v-model="form.category" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('finance.amount')" required><UInput v-model.number="form.amount" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('finance.debit_account')"><USelect v-model="form.debitAccount" :options="accountOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('finance.credit_account')"><USelect v-model="form.creditAccount" :options="accountOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.description')" required><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
        <UFormGroup :label="$t('finance.reference')"><UInput v-model="form.reference" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('finance.delete_journal')" :loading="deleting" @confirm="deleteJournal" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const journals = ref<Record<string, unknown>[]>([])
const accountOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)

const columns: TableColumn[] = [
  { key: 'journalNumber', label: 'finance.journal_number', sortable: true },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'type', label: 'finance.type' },
  { key: 'category', label: 'finance.category' },
  { key: 'amount', label: 'finance.amount', type: 'currency' },
  { key: 'description', label: 'common.description' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [
  { key: 'type', label: 'finance.type', type: 'select' as const, options: [{ label: t('finance.income'), value: 'income' }, { label: t('finance.expense'), value: 'expense' }, { label: t('finance.transfer'), value: 'transfer' }] },
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.draft'), value: 'draft' }, { label: t('status.posted'), value: 'posted' }] },
]
const form = reactive({ journalNumber: '', date: $dayjs().format('YYYY-MM-DD'), type: 'income', category: '', amount: 0, debitAccount: '', creditAccount: '', description: '', reference: '' })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

const fetchJournals = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { journals.value = await api.paginate('/finance/journals', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchJournals(filters)
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { journalNumber: '', date: $dayjs().format('YYYY-MM-DD'), type: 'income', category: '', amount: 0, debitAccount: '', creditAccount: '', description: '', reference: '' }); showForm.value = true }
const editJournal = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveJournal = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/finance/journals/${editId.value}`, form); toast.add({ title: t('finance.journal_updated'), color: 'success' }) } else { await api.post('/finance/journals', form); toast.add({ title: t('finance.journal_created'), color: 'success' }) } showForm.value = false; fetchJournals() } catch {} finally { saving.value = false } }
const postJournal = async (row: Record<string, unknown>) => { try { await api.patch(`/finance/journals/${row.id}/post`); toast.add({ title: t('finance.journal_posted'), color: 'success' }); fetchJournals() } catch {} }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteJournal = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/finance/journals/${deleteTarget.value.id}`); toast.add({ title: t('finance.journal_deleted'), color: 'success' }); showDelete.value = false; fetchJournals() } catch {} finally { deleting.value = false } }
const handleExport = (format: string) => { window.open(`/api/v1/finance/journals/export?format=${format}`, '_blank') }
const fetchAccounts = async () => { try { accountOptions.value = (await api.get<{id:string;name:string}[]>('/finance/accounts')).map(a => ({ label: a.name, value: a.id })) } catch {} }

onMounted(() => { fetchJournals(); fetchAccounts() })
</script>
