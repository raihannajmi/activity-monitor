<script setup>
import { onMounted, ref, watch } from 'vue'
import { useNotesStore } from '../stores/notes'
import { storeToRefs } from 'pinia'
import { Plus, Search, Trash2, Edit3, Calendar } from 'lucide-vue-next'
import dayjs from 'dayjs'

const notesStore = useNotesStore()
const { notes, selectedNote } = storeToRefs(notesStore)

const searchQ = ref('')
let saveTimeout = null

onMounted(() => {
  notesStore.fetchNotes()
})

// Trigger fetch notes on search query updates
watch(searchQ, (newQ) => {
  notesStore.fetchNotes(newQ)
})

// Auto-save changes in selected note
function triggerAutoSave() {
  if (!selectedNote.value) return
  if (saveTimeout) clearTimeout(saveTimeout)
  saveTimeout = setTimeout(() => {
    notesStore.updateNote(selectedNote.value.ID, selectedNote.value.Title, selectedNote.value.Content)
  }, 1000)
}

async function handleCreateNote() {
  await notesStore.createNote()
}
</script>

<template>
  <div class="page-container" style="height: 100vh; overflow: hidden;">
    <header class="page-header" style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--color-border); padding: 20px 32px;">
      <div>
        <h1 class="page-title" style="font-size: 22px; font-weight: 700; color: var(--color-text-primary);">Brain Dump</h1>
        <p class="page-subtitle" style="font-size: 13px; color: var(--color-text-secondary); margin-top: 2px;">Tangkap dan simpan ide-ide cemerlang Anda.</p>
      </div>
      <button class="btn btn-primary" @click="handleCreateNote" style="display: flex; align-items: center; gap: 8px;">
        <Plus :size="16" /> Catatan Baru
      </button>
    </header>

    <div class="notes-container">
      <!-- Notes Sidebar -->
      <aside class="notes-sidebar">
        <!-- Search bar -->
        <div style="padding: 16px; border-bottom: 1px solid var(--color-border); position: relative;">
          <Search :size="14" style="position: absolute; left: 26px; top: 29px; color: var(--color-text-muted);" />
          <input 
            v-model="searchQ"
            type="text" 
            placeholder="Cari catatan..."
            style="width: 100%; padding: 8px 12px 8px 32px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); outline: none; font-size: 13px;"
          />
        </div>

        <!-- Note List -->
        <div style="flex: 1; overflow-y: auto; display: flex; flex-direction: column;">
          <div 
            v-for="note in notes" 
            :key="note.ID"
            class="note-list-item"
            style="padding: 16px 20px; border-bottom: 1px solid var(--color-border); cursor: pointer; transition: background 0.15s;"
            :style="{ background: selectedNote?.ID === note.ID ? '#F3F4F6' : 'white' }"
            @click="notesStore.selectNote(note)"
          >
            <h3 style="font-weight: 600; font-size: 14px; color: var(--color-text-primary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-bottom: 4px;">
              {{ note.Title || 'Tanpa Judul' }}
            </h3>
            <p style="font-size: 12px; color: var(--color-text-secondary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-bottom: 6px;">
              {{ note.Content || 'Belum ada konten...' }}
            </p>
            <div style="font-size: 10px; color: var(--color-text-muted); display: flex; align-items: center; gap: 4px;">
              <Calendar :size="10" />
              <span>{{ dayjs(note.UpdatedAt).format('DD MMM YYYY') }}</span>
            </div>
          </div>
        </div>
      </aside>

      <!-- Note Editor -->
      <section class="notes-editor">
        <div v-if="selectedNote" style="display: flex; flex-direction: column; height: 100%;">
          <div style="display: flex; justify-content: flex-end; margin-bottom: 12px;">
            <button class="btn btn-ghost btn-danger-ghost btn-sm" @click="notesStore.deleteNote(selectedNote.ID)">
              <Trash2 :size="16" /> Hapus
            </button>
          </div>

          <input 
            v-model="selectedNote.Title"
            type="text" 
            class="note-title-input" 
            placeholder="Judul Catatan..."
            @input="triggerAutoSave"
          />

          <textarea 
            v-model="selectedNote.Content"
            class="note-textarea"
            placeholder="Tulis ide atau catatan kamu disini..."
            @input="triggerAutoSave"
          ></textarea>
        </div>

        <div v-else style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; color: var(--color-text-secondary); gap: 12px;">
          <Edit3 :size="48" style="color: var(--color-text-muted);" />
          <p style="font-size: 14px; font-weight: 500;">Pilih catatan atau buat catatan baru.</p>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.note-list-item:hover {
  background-color: #F9FAFB !important;
}
</style>
