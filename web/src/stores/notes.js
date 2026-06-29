import { defineStore } from 'pinia'
import { ref } from 'vue'
import client from '../api/client'
import { useUIStore } from './ui'

export const useNotesStore = defineStore('notes', () => {
  const ui = useUIStore()
  const notes = ref([])
  const selectedNote = ref(null)

  async function fetchNotes(searchQuery = '') {
    try {
      const res = await client.get(`/notes?q=${searchQuery}`)
      notes.value = res.data || []
    } catch (err) {
      ui.showToast('Gagal memuat catatan', 'error')
    }
  }

  async function selectNote(note) {
    selectedNote.value = { ...note }
  }

  async function createNote(title = 'Catatan Baru', content = '') {
    try {
      const res = await client.post('/notes', { title, content })
      notes.value.unshift(res.data)
      selectedNote.value = res.data
      ui.showToast('Catatan dibuat', 'success')
      return res.data
    } catch (err) {
      ui.showToast('Gagal membuat catatan', 'error')
    }
  }

  async function updateNote(noteId, title, content) {
    try {
      const res = await client.put(`/notes/${noteId}`, { title, content })
      const idx = notes.value.findIndex(n => n.ID === noteId)
      if (idx !== -1) {
        notes.value[idx] = res.data
      }
      if (selectedNote.value && selectedNote.value.ID === noteId) {
        selectedNote.value = res.data
      }
    } catch (err) {
      ui.showToast('Gagal menyimpan catatan', 'error')
    }
  }

  async function deleteNote(noteId) {
    try {
      await client.delete(`/notes/${noteId}`)
      notes.value = notes.value.filter(n => n.ID !== noteId)
      if (selectedNote.value && selectedNote.value.ID === noteId) {
        selectedNote.value = null
      }
      ui.showToast('Catatan dihapus', 'success')
    } catch (err) {
      ui.showToast('Gagal menghapus catatan', 'error')
    }
  }

  return {
    notes,
    selectedNote,
    fetchNotes,
    selectNote,
    createNote,
    updateNote,
    deleteNote,
  }
})
