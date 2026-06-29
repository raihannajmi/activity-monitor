import { defineStore } from 'pinia'
import { ref } from 'vue'
import client from '../api/client'
import { useUIStore } from './ui'

export const useTasksStore = defineStore('tasks', () => {
  const ui = useUIStore()
  const tasks = ref([])
  const selectedTask = ref(null)
  const isDrawerOpen = ref(false)
  const labels = ref([])

  async function fetchTasks(filter = 'all') {
    try {
      const res = await client.get(`/tasks?filter=${filter}`)
      tasks.value = res.data || []
    } catch (err) {
      ui.showToast('Gagal memuat task', 'error')
    }
  }

  async function fetchLabels() {
    try {
      const res = await client.get('/labels')
      labels.value = res.data || []
    } catch (err) {
      console.error('Gagal memuat label', err)
    }
  }

  async function selectTask(taskId) {
    try {
      const res = await client.get(`/tasks/${taskId}`)
      selectedTask.value = res.data
      isDrawerOpen.value = true
    } catch (err) {
      ui.showToast('Gagal memuat detail task', 'error')
    }
  }

  function closeDrawer() {
    selectedTask.value = null
    isDrawerOpen.value = false
  }

  async function createTask(title, description = '', priority = 'medium', deadline = '') {
    try {
      const res = await client.post('/tasks', { title, description, priority, deadline })
      tasks.value.push(res.data)
      ui.showToast('Task baru berhasil ditambahkan', 'success')
      return res.data
    } catch (err) {
      ui.showToast('Gagal membuat task baru', 'error')
    }
  }

  async function updateTask(taskId, updates) {
    try {
      const res = await client.put(`/tasks/${taskId}`, updates)
      const index = tasks.value.findIndex(t => t.ID === taskId)
      if (index !== -1) {
        tasks.value[index] = { ...tasks.value[index], ...res.data }
      }
      if (selectedTask.value && selectedTask.value.ID === taskId) {
        selectedTask.value = { ...selectedTask.value, ...res.data }
      }
      ui.showToast('Task berhasil diperbarui', 'success')
    } catch (err) {
      ui.showToast('Gagal memperbarui task', 'error')
    }
  }

  async function deleteTask(taskId) {
    try {
      await client.delete(`/tasks/${taskId}`)
      tasks.value = tasks.value.filter(t => t.ID !== taskId)
      if (selectedTask.value && selectedTask.value.ID === taskId) {
        closeDrawer()
      }
      ui.showToast('Task berhasil dihapus', 'success')
    } catch (err) {
      ui.showToast('Gagal menghapus task', 'error')
    }
  }

  async function moveTask(taskId, columnId, prevId = '', nextId = '') {
    try {
      await client.put(`/tasks/${taskId}/move`, { column_id: columnId, prev_id: prevId, next_id: nextId })
      // Local swap to minimize server requests
      const t = tasks.value.find(task => task.ID === taskId)
      if (t) {
        t.ColumnID = columnId
      }
    } catch (err) {
      ui.showToast('Gagal memindahkan task', 'error')
      fetchTasks() // revert
    }
  }

  // Subtasks
  async function addSubtask(title, deadline = '') {
    if (!selectedTask.value) return
    try {
      const res = await client.post(`/tasks/${selectedTask.value.ID}/subtasks`, { title, deadline })
      if (!selectedTask.value.Subtasks) selectedTask.value.Subtasks = []
      selectedTask.value.Subtasks.push(res.data)
      ui.showToast('Subtask berhasil ditambahkan', 'success')
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal menambah subtask', 'error')
    }
  }

  async function toggleSubtask(subtaskId) {
    if (!selectedTask.value) return
    try {
      await client.put(`/tasks/${selectedTask.value.ID}/subtasks/${subtaskId}/toggle`)
      const sub = selectedTask.value.Subtasks.find(s => s.ID === subtaskId)
      if (sub) {
        sub.IsCompleted = !sub.IsCompleted
      }
      ui.showToast('Status subtask diperbarui', 'success')
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal memperbarui subtask', 'error')
    }
  }

  async function deleteSubtask(subtaskId) {
    if (!selectedTask.value) return
    try {
      await client.delete(`/tasks/${selectedTask.value.ID}/subtasks/${subtaskId}`)
      selectedTask.value.Subtasks = selectedTask.value.Subtasks.filter(s => s.ID !== subtaskId)
      ui.showToast('Subtask berhasil dihapus', 'success')
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal menghapus subtask', 'error')
    }
  }

  // Checklists
  async function addChecklist(title) {
    if (!selectedTask.value) return
    try {
      const res = await client.post(`/tasks/${selectedTask.value.ID}/checklists`, { title })
      if (!selectedTask.value.Checklists) selectedTask.value.Checklists = []
      selectedTask.value.Checklists.push(res.data)
      ui.showToast('Checklist ditambahkan', 'success')
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal menambah checklist', 'error')
    }
  }

  async function deleteChecklist(checklistId) {
    if (!selectedTask.value) return
    try {
      await client.delete(`/tasks/${selectedTask.value.ID}/checklists/${checklistId}`)
      selectedTask.value.Checklists = selectedTask.value.Checklists.filter(c => c.ID !== checklistId)
      ui.showToast('Checklist dihapus', 'success')
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal menghapus checklist', 'error')
    }
  }

  async function addChecklistItem(checklistId, title) {
    if (!selectedTask.value) return
    try {
      const res = await client.post(`/tasks/${selectedTask.value.ID}/checklists/${checklistId}/items`, { title })
      const checklist = selectedTask.value.Checklists.find(c => c.ID === checklistId)
      if (checklist) {
        if (!checklist.Items) checklist.Items = []
        checklist.Items.push(res.data)
      }
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal menambah item checklist', 'error')
    }
  }

  async function toggleChecklistItem(checklistId, item) {
    if (!selectedTask.value) return
    try {
      await client.put(`/tasks/${selectedTask.value.ID}/checklists/${checklistId}/items/${item.ID}`, {
        title: item.Title,
        completed: !item.Completed
      })
      item.Completed = !item.Completed
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal memperbarui checklist item', 'error')
    }
  }

  async function deleteChecklistItem(checklistId, itemId) {
    if (!selectedTask.value) return
    try {
      await client.delete(`/tasks/${selectedTask.value.ID}/checklists/${checklistId}/items/${itemId}`)
      const checklist = selectedTask.value.Checklists.find(c => c.ID === checklistId)
      if (checklist) {
        checklist.Items = checklist.Items.filter(i => i.ID !== itemId)
      }
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal menghapus item', 'error')
    }
  }

  // Comments
  async function addComment(content) {
    if (!selectedTask.value) return
    try {
      const res = await client.post(`/tasks/${selectedTask.value.ID}/comments`, { content })
      if (!selectedTask.value.Comments) selectedTask.value.Comments = []
      selectedTask.value.Comments.unshift(res.data)
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal menambah komentar', 'error')
    }
  }

  // Labels
  async function toggleLabel(labelId) {
    if (!selectedTask.value) return
    try {
      await client.post(`/tasks/${selectedTask.value.ID}/labels`, { label_id: labelId })
      const hasLabel = selectedTask.value.Labels?.some(l => l.ID === labelId)
      if (hasLabel) {
        selectedTask.value.Labels = selectedTask.value.Labels.filter(l => l.ID !== labelId)
      } else {
        const labelObj = labels.value.find(l => l.ID === labelId)
        if (labelObj) {
          if (!selectedTask.value.Labels) selectedTask.value.Labels = []
          selectedTask.value.Labels.push(labelObj)
        }
      }
      fetchTasks()
    } catch (err) {
      ui.showToast('Gagal memperbarui label', 'error')
    }
  }

  return {
    tasks,
    selectedTask,
    isDrawerOpen,
    labels,
    fetchTasks,
    fetchLabels,
    selectTask,
    closeDrawer,
    createTask,
    updateTask,
    deleteTask,
    moveTask,
    // Subtasks
    addSubtask,
    toggleSubtask,
    deleteSubtask,
    // Checklists
    addChecklist,
    deleteChecklist,
    addChecklistItem,
    toggleChecklistItem,
    deleteChecklistItem,
    // Comments
    addComment,
    // Labels
    toggleLabel,
  }
})
