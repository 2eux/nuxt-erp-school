<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('assets.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('assets.subtitle') }}</p></div>
      <UButton v-if="permissions.can('assets.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('assets.add_asset') }}</UButton>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('assets.total')" :value="stats.total" icon="i-heroicons-building-office" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('assets.total_value')" :value="formatCurrency(stats.totalValue)" icon="i-heroicons-currency-dollar" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('assets.depreciation')" :value="formatCurrency(stats.depreciation)" icon="i-heroicons-arrow-trending-down" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('assets.maintenance')" :value="stats.maintenance" icon="i-heroicons-wrench-screwdriver" color="purple" :loading="statsLoading" />
    </div>
    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />
    <DataTable :columns="columns" :rows="assets" :loading="loading" :empty-title="$t('assets.no_assets')" :show-export="true" @export="handleExport">
      <template #cell-condition="{ row }"><StatusBadge :status="(row.condition as string) || 'good'" /></template>
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'active'" /></template>
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editAsset(row as Record<string, unknown>)" />
        <UButton v-if="(row.status as string) !== 'maintenance'" color="amber" variant="ghost" size="xs" icon="i-heroicons-wrench" @click="setMaintenance(row)" />
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
      </template>
    </DataTable>
    <FormDialog v-model="showForm" :title="editing ? $t('assets.edit_asset') : $t('assets.add_asset')" :loading="saving" @submit="saveAsset" @cancel="showForm=false">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('assets.code')"><UInput v-model="form.code" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.name')" required class="sm:col-span-2"><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.category')"><UInput v-model="form.category" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.acquisition_price')"><UInput v-model.number="form.acquisitionPrice" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.current_value')"><UInput v-model.number="form.currentValue" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.depreciation_rate')"><UInput v-model.number="form.depreciationRate" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.acquisition_date')"><UInput v-model="form.acquisitionDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.location')"><UInput v-model="form.location" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.condition')"><USelect v-model="form.condition" :options="conditionOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('assets.status')"><USelect v-model="form.status" :options="[{label:t('status.active'),value:'active'},{label:t('assets.maintenance'),value:'maintenance'},{label:t('assets.disposed'),value:'disposed'}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.description')" class="sm:col-span-2"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>
    <ConfirmDialog v-model="showDelete" :title="$t('assets.delete_asset')" :loading="deleting" @confirm="deleteAsset" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast()
const loading = ref(false); const statsLoading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const assets = ref<Record<string, unknown>[]>([]); const deleteTarget = ref<Record<string, unknown> | null>(null)
const stats = reactive({ total: 0, totalValue: 0, depreciation: 0, maintenance: 0 })
const columns: TableColumn[] = [
  { key: 'code', label: 'assets.code', sortable: true }, { key: 'name', label: 'assets.name' },
  { key: 'category', label: 'assets.category' }, { key: 'acquisitionDate', label: 'assets.acquisition_date', type: 'date' },
  { key: 'currentValue', label: 'assets.current_value', type: 'currency' },
  { key: 'condition', label: 'assets.condition', type: 'status' }, { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [
  { key: 'category', label: 'assets.category', type: 'text' as const },
  { key: 'condition', label: 'assets.condition', type: 'select' as const, options: [{ label: t('inventory.excellent'), value: 'excellent' }, { label: t('inventory.good'), value: 'good' }, { label: t('inventory.damaged'), value: 'damaged' }] },
]
const conditionOpts = [{ label: t('inventory.excellent'), value: 'excellent' }, { label: t('inventory.good'), value: 'good' }, { label: t('inventory.fair'), value: 'fair' }, { label: t('inventory.damaged'), value: 'damaged' }]
const form = reactive({ code: '', name: '', category: '', description: '', acquisitionDate: '', acquisitionPrice: 0, currentValue: 0, depreciationRate: 0, location: '', condition: 'good', status: 'active' })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

const fetchAssets = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { assets.value = await api.paginate('/assets', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/assets/stats')) } catch {} finally { statsLoading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchAssets(filters)
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { code: '', name: '', category: '', description: '', acquisitionDate: '', acquisitionPrice: 0, currentValue: 0, depreciationRate: 0, location: '', condition: 'good', status: 'active' }); showForm.value = true }
const editAsset = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveAsset = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/assets/${editId.value}`, form); toast.add({ title: t('assets.asset_updated'), color: 'success' }) } else { await api.post('/assets', form); toast.add({ title: t('assets.asset_created'), color: 'success' }) } showForm.value = false; fetchAssets(); fetchStats() } catch {} finally { saving.value = false } }
const setMaintenance = async (row: Record<string, unknown>) => { try { await api.patch(`/assets/${row.id}/maintenance`); toast.add({ title: t('assets.maintenance_set'), color: 'success' }); fetchAssets(); fetchStats() } catch {} }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteAsset = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/assets/${deleteTarget.value.id}`); toast.add({ title: t('assets.asset_deleted'), color: 'success' }); showDelete.value = false; fetchAssets(); fetchStats() } catch {} finally { deleting.value = false } }
const handleExport = (format: string) => { window.open(`/api/v1/assets/export?format=${format}`, '_blank') }
onMounted(() => { fetchAssets(); fetchStats() })
</script>
