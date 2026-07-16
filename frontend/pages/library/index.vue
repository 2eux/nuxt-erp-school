<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('library.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('library.subtitle') }}</p></div>
      <UButton v-if="permissions.can('library.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAddBook">{{ $t('library.add_book') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('library.total_books')" :value="stats.totalBooks" icon="i-heroicons-book-open" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('library.available')" :value="stats.available" icon="i-heroicons-check-circle" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('library.borrowed')" :value="stats.borrowed" icon="i-heroicons-arrow-path" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('library.overdue')" :value="stats.overdue" icon="i-heroicons-exclamation-triangle" color="red" :loading="statsLoading" />
    </div>

    <UTabs :items="tabs">
      <template #item="{ item }">
        <div class="pt-4">
          <template v-if="item.key === 'books'">
            <DataFilter :filter-fields="bookFilterFields" :searchable="true" @apply="fetchBooks" />
            <DataTable :columns="bookColumns" :rows="books" :loading="loading" :empty-title="$t('library.no_books')" :show-export="false">
              <template #cell-available="{ row }"><span class="text-sm font-semibold" :class="(row.available as number) > 0 ? 'text-emerald-600' : 'text-red-600'">{{ row.available }}/{{ row.quantity }}</span></template>
              <template #item-actions="{ row }">
                <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editBook(row as Record<string, unknown>)" />
                <UButton v-if="(row.available as number) > 0" color="emerald" variant="ghost" size="xs" icon="i-heroicons-arrow-right-on-rectangle" @click="borrowBook(row as Record<string, unknown>)" />
              </template>
            </DataTable>
          </template>
          <template v-else>
            <DataTable :columns="borrowingColumns" :rows="borrowings" :loading="loading" :empty-title="$t('library.no_borrowings')" :show-export="false">
              <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'borrowed'" /></template>
              <template #item-actions="{ row }">
                <UButton v-if="(row.status as string) === 'borrowed' || (row.status as string) === 'overdue'" color="emerald" variant="ghost" size="xs" icon="i-heroicons-arrow-uturn-left" @click="returnBook(row)" />
              </template>
            </DataTable>
          </template>
        </div>
      </template>
    </UTabs>

    <FormDialog v-model="showBookForm" :title="editing ? $t('library.edit_book') : $t('library.add_book')" :loading="saving" @submit="saveBook" @cancel="showBookForm=false">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('library.isbn')"><UInput v-model="bookForm.isbn" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.title')" required><UInput v-model="bookForm.title" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.author')" required><UInput v-model="bookForm.author" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.publisher')"><UInput v-model="bookForm.publisher" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.category')"><UInput v-model="bookForm.category" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.publish_year')"><UInput v-model.number="bookForm.publishYear" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.quantity')"><UInput v-model.number="bookForm.quantity" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.rack')"><UInput v-model="bookForm.rack" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <FormDialog v-model="showBorrowForm" :title="$t('library.borrow_book')" :loading="saving" @submit="saveBorrow" @cancel="showBorrowForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('library.borrower')" required><USelect v-model="borrowForm.borrowerId" :options="borrowerOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.borrower_type')" required><USelect v-model="borrowForm.borrowerType" :options="[{label:t('students.title'),value:'student'},{label:t('teachers.title'),value:'teacher'},{label:t('employees.title'),value:'employee'}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('library.due_date')" required><UInput v-model="borrowForm.dueDate" type="date" color="gray" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const statsLoading = ref(false); const saving = ref(false)
const showBookForm = ref(false); const showBorrowForm = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const books = ref<Record<string, unknown>[]>([]); const borrowings = ref<Record<string, unknown>[]>([])
const borrowerOptions = ref<{ label: string; value: string }[]>([])
const stats = reactive({ totalBooks: 0, available: 0, borrowed: 0, overdue: 0 })
const borrowTarget = ref<Record<string, unknown> | null>(null)

const tabs = [{ key: 'books', label: t('library.books') }, { key: 'borrowings', label: t('library.borrowings') }]
const bookColumns: TableColumn[] = [
  { key: 'isbn', label: 'library.isbn' }, { key: 'title', label: 'library.title', sortable: true },
  { key: 'author', label: 'library.author' }, { key: 'category', label: 'library.category' },
  { key: 'available', label: 'library.available' }, { key: 'rack', label: 'library.rack' },
]
const borrowingColumns: TableColumn[] = [
  { key: 'bookTitle', label: 'library.title' }, { key: 'borrowerName', label: 'library.borrower' },
  { key: 'borrowDate', label: 'library.borrow_date', type: 'date' }, { key: 'dueDate', label: 'library.due_date', type: 'date' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const bookFilterFields = [{ key: 'category', label: 'library.category', type: 'text' as const }]
const bookForm = reactive({ isbn: '', title: '', author: '', publisher: '', category: '', publishYear: 2024, quantity: 1, rack: '' })
const borrowForm = reactive({ borrowerId: '', borrowerType: 'student', dueDate: '' })

const fetchBooks = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { books.value = await api.paginate('/library/books', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const fetchBorrowings = async () => { try { borrowings.value = await api.paginate('/library/borrowings').then(r => r.data) } catch {} }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/library/stats')) } catch {} finally { statsLoading.value = false } }
const openAddBook = () => { editing.value = false; editId.value = null; Object.assign(bookForm, { isbn: '', title: '', author: '', publisher: '', category: '', publishYear: 2024, quantity: 1, rack: '' }); showBookForm.value = true }
const editBook = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(bookForm, row); showBookForm.value = true }
const saveBook = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/library/books/${editId.value}`, bookForm); toast.add({ title: t('library.book_updated'), color: 'success' }) } else { await api.post('/library/books', bookForm); toast.add({ title: t('library.book_created'), color: 'success' }) } showBookForm.value = false; fetchBooks(); fetchStats() } catch {} finally { saving.value = false } }
const borrowBook = (row: Record<string, unknown>) => { borrowTarget.value = row; borrowForm.dueDate = $dayjs().add(7, 'day').format('YYYY-MM-DD'); showBorrowForm.value = true }
const saveBorrow = async () => { saving.value = true; try { await api.post('/library/borrowings', { ...borrowForm, bookId: borrowTarget.value?.id }); toast.add({ title: t('library.book_borrowed'), color: 'success' }); showBorrowForm.value = false; fetchBooks(); fetchBorrowings(); fetchStats() } catch {} finally { saving.value = false } }
const returnBook = async (row: Record<string, unknown>) => { try { await api.patch(`/library/borrowings/${row.id}/return`); toast.add({ title: t('library.book_returned'), color: 'success' }); fetchBooks(); fetchBorrowings(); fetchStats() } catch {} }
const fetchBorrowers = async () => { try { const students = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })); borrowerOptions.value = students } catch {} }
onMounted(() => { fetchBooks(); fetchBorrowings(); fetchStats(); fetchBorrowers() })
</script>
