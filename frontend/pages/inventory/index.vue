<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('inventory.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('inventory.subtitle') }}</p></div>
      <UButton v-if="permissions.can('inventory.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('inventory.add_item') }}</UButton>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('inventory.total_items')" :value="stats.total" icon="i-heroicons-cube" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('inventory.excellent')" :value="stats.excellent" icon="i-heroicons-check-badge" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('inventory.damaged')" :value="stats.damaged" icon="i-heroicons-exclamation-triangle" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('inventory.lost')" :value="stats.lost" icon="i-heroicons-x-circle" color="red" :loading="statsLoading" />
    </div>
    <DataFilter :filter-fields="filterFields" :searchable="true" @apply="handleFilter" />
    <DataTable :columns="columns" :rows="items" :loading="loading" :empty-title="$t('inventory.no_items')" :show-export="true" @export="handleExport">
      <template #cell-condition="{ row }"><StatusBadge :status="(row.condition as string) || 'good'" /></template>
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editItem(row as Record<string, unknown>)" />
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
      </template>
    </DataTable>
    <FormDialog v-model="showForm" :title="editing ? $t('inventory.edit_item') : $t('inventory.add_item')" :loading="saving" @submit="saveItem" @cancel="showForm=false">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('inventory.code')" required><UInput v-model="form.code" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('inventory.name')" required class="sm:col-span-2"><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('inventory.category')"><UInput v-model="form.category" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('inventory.quantity')"><UInput v-model.number="form.quantity" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('inventory.unit')"><UInput v-model="form.unit" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('inventory.location')"><UInput v-model="form.location" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('inventory.condition')"><USelect v-model="form.condition" :options="conditionOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('inventory.purchase_price')"><UInput v-model.number="form.purchasePrice" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('inventory.purchase_date')"><UInput v-model="form.purchaseDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.description')" class="sm:col-span-2"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>
    <ConfirmDialog v-model="showDelete" :title="$t('inventory.delete_item')" :loading="deleting" @confirm="deleteItem" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const statsLoading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const items = ref<Record<string, unknown>[]>([]); const deleteTarget = ref<Record<string, unknown> | null>(null)
const stats = reactive({ total: 0, excellent: 0, damaged: 0, lost: 0 })
const columns: TableColumn[] = [
  { key: 'code', label: 'inventory.code', sortable: true },
  { key: 'name', label: 'inventory.name' },
  { key: 'category', label: 'inventory.category' },
  { key: 'quantity', label: 'inventory.quantity', type: 'number' },
  { key: 'location', label: 'inventory.location' },
  { key: 'condition', label: 'inventory.condition', type: 'status' },
]
const filterFields = [
  { key: 'category', label: 'inventory.category', type: 'text' as const },
  { key: 'condition', label: 'inventory.condition', type: 'select' as const, options: [{ label: t('inventory.excellent'), value: 'excellent' }, { label: t('inventory.good'), value: 'good' }, { label: t('inventory.fair'), value: 'fair' }, { label: t('inventory.damaged'), value: 'damaged' }, { label: t('inventory.lost'), value: 'lost' }] },
]
const conditionOpts = [{ label: t('inventory.excellent'), value: 'excellent' }, { label: t('inventory.good'), value: 'good' }, { label: t('inventory.fair'), value: 'fair' }, { label: t('inventory.damaged'), value: 'damaged' }, { label: t('inventory.lost'), value: 'lost' }]
const form = reactive({ code: '', name: '', category: '', description: '', quantity: 0, unit: 'pcs', location: '', condition: 'good', purchaseDate: '', purchasePrice: 0 })

const fetchItems = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { items.value = await api.paginate('/inventory', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/inventory/stats')) } catch {} finally { statsLoading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchItems(filters)
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { code: '', name: '', category: '', description: '', quantity: 0, unit: 'pcs', location: '', condition: 'good', purchaseDate: '', purchasePrice: 0 }); showForm.value = true }
const editItem = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveItem = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/inventory/${editId.value}`, form); toast.add({ title: t('inventory.item_updated'), color: 'success' }) } else { await api.post('/inventory', form); toast.add({ title: t('inventory.item_created'), color: 'success' }) } showForm.value = false; fetchItems(); fetchStats() } catch {} finally { saving.value = false } }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteItem = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/inventory/${deleteTarget.value.id}`); toast.add({ title: t('inventory.item_deleted'), color: 'success' }); showDelete.value = false; fetchItems(); fetchStats() } catch {} finally { deleting.value = false } }
const handleExport = (format: string) => { window.open(`/api/v1/inventory/export?format=${format}`, '_blank') }
onMounted(() => { fetchItems(); fetchStats() })
</script>
