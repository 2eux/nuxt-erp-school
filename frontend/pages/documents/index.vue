<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('documents.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('documents.subtitle') }}</p></div>
      <UButton v-if="permissions.can('documents.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('documents.upload') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="folder in folders" :key="folder" class="card cursor-pointer hover:shadow-md transition-shadow" @click="selectedFolder=folder">
        <div class="flex items-center gap-3">
          <UIcon name="i-heroicons-folder" class="w-8 h-8 text-amber-500" />
          <div><p class="text-sm font-medium text-gray-900 dark:text-white">{{ folder }}</p><p class="text-xs text-gray-500">{{ getFileCount(folder) }} {{ $t('documents.files') }}</p></div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <h3 class="card-title">{{ selectedFolder || $t('documents.all_documents') }}</h3>
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-folder-plus" @click="createFolder">{{ $t('documents.new_folder') }}</UButton>
      </div>
      <DataTable :columns="columns" :rows="filteredDocs" :loading="loading" :empty-title="$t('documents.no_documents')" :show-export="false">
        <template #item-actions="{ row }">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-arrow-down-tray" @click="downloadDoc(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </template>
      </DataTable>
    </div>

    <FormDialog v-model="showUpload" :title="$t('documents.upload')" :loading="uploading" @submit="doUpload" @cancel="showUpload=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('documents.title')"><UInput v-model="uploadForm.title" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('documents.folder')"><USelect v-model="uploadForm.folder" :options="folderOptions" color="gray" /></UFormGroup>
        <FileUpload :multiple="true" :accept-hint="$t('documents.upload_hint')" @files-selected="handleFilesSelected" />
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const uploading = ref(false); const showUpload = ref(false)
const documents = ref<Record<string, unknown>[]>([]); const deleteTarget = ref<Record<string, unknown> | null>(null)
const selectedFolder = ref(''); const selectedFiles = ref<File[]>([])
const folders = computed(() => [...new Set(documents.value.map(d => d.folder as string))].filter(Boolean))
const folderOptions = computed(() => [{ label: t('documents.all'), value: '' }, ...folders.value.map(f => ({ label: f, value: f }))])
const filteredDocs = computed(() => selectedFolder.value ? documents.value.filter(d => d.folder === selectedFolder.value) : documents.value)
const getFileCount = (folder: string) => documents.value.filter(d => d.folder === folder).length
const columns: TableColumn[] = [
  { key: 'title', label: 'documents.title' }, { key: 'folder', label: 'documents.folder' },
  { key: 'fileName', label: 'documents.file_name' }, { key: 'fileSize', label: 'documents.file_size' },
  { key: 'createdAt', label: 'common.created_at', type: 'date' }, { key: 'version', label: 'documents.version' },
]
const uploadForm = reactive({ title: '', folder: '' })

const fetchDocuments = async () => { loading.value = true; try { documents.value = await api.paginate('/documents').then(r => r.data) } catch {} finally { loading.value = false } }
const openAdd = () => { uploadForm.title = ''; uploadForm.folder = ''; showUpload.value = true }
const handleFilesSelected = (files: File[]) => { selectedFiles.value = files }
const doUpload = async () => { uploading.value = true; try { for (const file of selectedFiles.value) { const fd = new FormData(); fd.append('file', file); fd.append('title', uploadForm.title || file.name); fd.append('folder', uploadForm.folder); await api.upload('/documents', fd) } toast.add({ title: t('documents.uploaded'), color: 'success' }); showUpload.value = false; fetchDocuments() } catch {} finally { uploading.value = false } }
const createFolder = () => { const name = prompt(t('documents.enter_folder_name')); if (name) { api.post('/documents/folders', { name }); fetchDocuments() } }
const downloadDoc = (row: Record<string, unknown>) => { window.open(row.url as string, '_blank') }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; if (confirm(t('documents.confirm_delete'))) { api.delete(`/documents/${row.id}`); fetchDocuments() } }
onMounted(() => fetchDocuments())
</script>
