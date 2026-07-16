<template>
  <div class="space-y-4">
    <div v-if="$slots.header" class="card">
      <slot name="header" />
    </div>

    <div v-if="showFilters || $slots.filters" class="card flex flex-wrap items-center gap-3">
      <slot name="filters" />
      <slot name="actions" />
    </div>

    <div class="card overflow-hidden !p-0">
      <div v-if="$slots.toolbar" class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-gray-700">
        <slot name="toolbar" />
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-800/50 border-b border-gray-200 dark:border-gray-700">
              <th v-if="selectable" class="w-10 px-3 py-3 text-left">
                <UCheckbox
                  :model-value="allSelected"
                  :indeterminate="someSelected && !allSelected"
                  @change="toggleSelectAll"
                />
              </th>
              <th
                v-for="col in visibleColumns"
                :key="col.key"
                class="px-3 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider whitespace-nowrap"
                :class="{
                  'cursor-pointer select-none hover:text-gray-700 dark:hover:text-gray-300': col.sortable,
                  'text-right': col.align === 'right',
                  'text-center': col.align === 'center',
                }"
                :style="col.width ? { width: `${col.width}px` } : {}"
                @click="col.sortable && toggleSort(col.key)"
              >
                <div class="flex items-center gap-1">
                  <span>{{ $t(col.label) || col.label }}</span>
                  <UIcon
                    v-if="col.sortable && sortKey === col.key"
                    :name="sortOrder === 'asc' ? 'i-heroicons-chevron-up' : 'i-heroicons-chevron-down'"
                    class="w-3.5 h-3.5 text-brand-600"
                  />
                </div>
              </th>
              <th v-if="$slots['item-actions']" class="w-20 px-3 py-3 text-right text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">
                {{ $t('common.actions') }}
              </th>
            </tr>
          </thead>

          <tbody>
            <template v-if="loading">
              <tr v-for="i in skeletonRows" :key="'skeleton-' + i">
                <td v-if="selectable" class="px-3 py-3">
                  <div class="h-4 w-4 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
                </td>
                <td v-for="col in visibleColumns" :key="col.key" class="px-3 py-3">
                  <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-[80%]" />
                </td>
              </tr>
            </template>

            <template v-else-if="rows.length === 0">
              <tr>
                <td
                  :colspan="columns.length + (selectable ? 1 : 0) + ($slots['item-actions'] ? 1 : 0)"
                  class="px-6 py-16"
                >
                  <EmptyState
                    :title="emptyTitle || $t('common.no_data')"
                    :description="emptyDescription || $t('common.no_data_description')"
                    :icon="emptyIcon || 'i-heroicons-inbox'"
                    :action-label="emptyActionLabel"
                    @action="$emit('empty-action')"
                  />
                </td>
              </tr>
            </template>

            <template v-else>
              <tr
                v-for="(row, index) in rows"
                :key="rowKey ? row[rowKey] : index"
                class="border-b border-gray-100 dark:border-gray-750 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                :class="{ 'bg-brand-50/30 dark:bg-brand-900/10': selectedIds?.includes(row[rowKey || 'id']) }"
              >
                <td v-if="selectable" class="px-3 py-3">
                  <UCheckbox
                    :model-value="selectedIds?.includes(row[rowKey || 'id'])"
                    @change="toggleRow(row[rowKey || 'id'])"
                  />
                </td>
                <td
                  v-for="col in visibleColumns"
                  :key="col.key"
                  class="px-3 py-3 whitespace-nowrap"
                  :class="{
                    'text-right': col.align === 'right',
                    'text-center': col.align === 'center',
                  }"
                >
                  <slot
                    :name="`cell-${col.key}`"
                    :row="row"
                    :value="getValue(row, col.key)"
                    :column="col"
                  >
                    <template v-if="col.type === 'status'">
                      <StatusBadge :status="getValue(row, col.key)" />
                    </template>
                    <template v-else-if="col.type === 'currency'">
                      <span class="text-sm text-gray-900 dark:text-white font-mono">
                        {{ formatCurrency(getValue(row, col.key)) }}
                      </span>
                    </template>
                    <template v-else-if="col.type === 'date'">
                      <span class="text-sm text-gray-700 dark:text-gray-300">
                        {{ formatDate(getValue(row, col.key)) }}
                      </span>
                    </template>
                    <template v-else-if="col.type === 'image'">
                      <img
                        :src="getValue(row, col.key)"
                        class="w-8 h-8 rounded-full object-cover"
                        alt=""
                      />
                    </template>
                    <template v-else-if="col.formatter">
                      <span class="text-sm text-gray-700 dark:text-gray-300">
                        {{ col.formatter(getValue(row, col.key), row) }}
                      </span>
                    </template>
                    <template v-else>
                      <span class="text-sm text-gray-700 dark:text-gray-300">
                        {{ getValue(row, col.key) }}
                      </span>
                    </template>
                  </slot>
                </td>
                <td v-if="$slots['item-actions']" class="px-3 py-3 text-right">
                  <slot name="item-actions" :row="row" />
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>

      <div
        v-if="showPagination && pagination && pagination.total > 0"
        class="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-800/50"
      >
        <div class="flex items-center gap-2">
          <span class="text-sm text-gray-600 dark:text-gray-400">
            {{ $t('table.showing') }}
          </span>
          <select
            v-model="pageSize"
            class="text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 px-2 py-1"
            @change="$emit('page-size-change', pageSize)"
          >
            <option v-for="size in pageSizes" :key="size" :value="size">{{ size }}</option>
          </select>
          <span class="text-sm text-gray-600 dark:text-gray-400">
            {{ $t('table.of') }} {{ pagination.total }}
          </span>
        </div>

        <div class="flex items-center gap-1">
          <UButton
            color="gray"
            variant="ghost"
            size="xs"
            icon="i-heroicons-chevron-double-left"
            :disabled="!pagination.hasPreviousPage"
            @click="$emit('page-change', 1)"
          />
          <UButton
            color="gray"
            variant="ghost"
            size="xs"
            icon="i-heroicons-chevron-left"
            :disabled="!pagination.hasPreviousPage"
            @click="$emit('page-change', pagination.page - 1)"
          />
          <span class="text-sm text-gray-600 dark:text-gray-400 px-2">
            {{ $t('table.page') }} {{ pagination.page }} {{ $t('table.of') }} {{ pagination.totalPages }}
          </span>
          <UButton
            color="gray"
            variant="ghost"
            size="xs"
            icon="i-heroicons-chevron-right"
            :disabled="!pagination.hasNextPage"
            @click="$emit('page-change', pagination.page + 1)"
          />
          <UButton
            color="gray"
            variant="ghost"
            size="xs"
            icon="i-heroicons-chevron-double-right"
            :disabled="!pagination.hasNextPage"
            @click="$emit('page-change', pagination.totalPages)"
          />
        </div>
      </div>

      <div v-if="showExport && rows.length > 0" class="flex items-center justify-end gap-2 px-4 py-2 border-t border-gray-100 dark:border-gray-700">
        <UButton
          color="gray"
          variant="ghost"
          size="xs"
          icon="i-heroicons-arrow-down-tray"
          @click="$emit('export', 'csv')"
        >
          CSV
        </UButton>
        <UButton
          color="gray"
          variant="ghost"
          size="xs"
          icon="i-heroicons-document-arrow-down"
          @click="$emit('export', 'pdf')"
        >
          PDF
        </UButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts" generic="T extends Record<string, unknown>">
import type { TableColumn } from '~/types'

const props = withDefaults(
  defineProps<{
    columns: TableColumn[]
    rows: T[]
    rowKey?: string
    loading?: boolean
    emptyTitle?: string
    emptyDescription?: string
    emptyIcon?: string
    emptyActionLabel?: string
    showPagination?: boolean
    showFilters?: boolean
    showExport?: boolean
    selectable?: boolean
    selectedIds?: string[]
    skeletonRows?: number
    pageSizes?: number[]
  }>(),
  {
    rowKey: 'id',
    loading: false,
    showPagination: true,
    showFilters: true,
    showExport: false,
    selectable: false,
    skeletonRows: 8,
    pageSizes: () => [10, 25, 50, 100],
  }
)

const emit = defineEmits<{
  'update:selectedIds': [ids: string[]]
  'sort': [key: string, order: string]
  'page-change': [page: number]
  'page-size-change': [size: number]
  'export': [format: 'csv' | 'pdf']
  'empty-action': []
  'row-click': [row: T]
}>()

const sortKey = ref('')
const sortOrder = ref<'asc' | 'desc'>('asc')
const pageSize = ref(props.pageSizes[0])

const visibleColumns = computed(() => props.columns.filter(c => c.type !== 'action'))

const pagination = computed(() => {
  if (!props.showPagination) return null
  const total = props.rows.length
  const totalPages = Math.ceil(total / pageSize.value)
  return {
    page: 1,
    limit: pageSize.value,
    total,
    totalPages,
    hasNextPage: totalPages > 1,
    hasPreviousPage: false,
  }
})

const allSelected = computed(() => {
  if (!props.selectedIds) return false
  return props.rows.length > 0 && props.rows.every(r => props.selectedIds?.includes(r[props.rowKey] as string))
})

const someSelected = computed(() => {
  if (!props.selectedIds) return false
  return props.rows.some(r => props.selectedIds?.includes(r[props.rowKey] as string)) && !allSelected.value
})

const { $dayjs } = useNuxtApp()
const { t, locale } = useI18n()

const getValue = (row: T, key: string): unknown => {
  return row[key]
}

const formatCurrency = (value: unknown): string => {
  const num = Number(value)
  if (isNaN(num)) return String(value)
  return new Intl.NumberFormat(locale.value === 'id' ? 'id-ID' : 'en-US', {
    style: 'currency',
    currency: locale.value === 'id' ? 'IDR' : 'USD',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(num)
}

const formatDate = (value: unknown): string => {
  if (!value) return '-'
  return $dayjs(String(value)).format('DD MMM YYYY')
}

const toggleSort = (key: string) => {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortOrder.value = 'asc'
  }
  emit('sort', sortKey.value, sortOrder.value)
}

const toggleSelectAll = (checked: boolean) => {
  const ids = checked ? props.rows.map(r => r[props.rowKey] as string) : []
  emit('update:selectedIds', ids)
}

const toggleRow = (id: string) => {
  const current = [...(props.selectedIds || [])]
  const idx = current.indexOf(id)
  if (idx === -1) {
    current.push(id)
  } else {
    current.splice(idx, 1)
  }
  emit('update:selectedIds', current)
}
</script>
