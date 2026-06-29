<script setup>
import { ref, onMounted, computed } from 'vue'
import { useTasksStore } from '../stores/tasks'
import { useTimerStore } from '../stores/timer'
import { storeToRefs } from 'pinia'
import { VueDraggable } from 'vue-draggable-plus'
import { 
  Plus, 
  MessageSquare, 
  Paperclip, 
  Clock, 
  Calendar, 
  CheckCircle2, 
  AlertCircle,
  Play,
  X
} from 'lucide-vue-next'
import TaskDetailDrawer from '../components/TaskDetailDrawer.vue'
import dayjs from 'dayjs'

const tasksStore = useTasksStore()
const { tasks, labels, isDrawerOpen, selectedTask } = storeToRefs(tasksStore)
const timerStore = useTimerStore()

const showCreateModal = ref(false)
const newTaskTitle = ref('')
const newTaskDesc = ref('')
const newTaskPriority = ref('medium')
const newTaskDeadline = ref('')

const todoTasks = computed({
  get: () => tasks.value.filter(t => t.Status === 'todo').sort((a, b) => a.Position - b.Position),
  set: (val) => updateColumnPositions(val, 'todo')
})

const inProgressTasks = computed({
  get: () => tasks.value.filter(t => t.Status === 'in_progress' || t.Status === 'inprogress').sort((a, b) => a.Position - b.Position),
  set: (val) => updateColumnPositions(val, 'in_progress')
})

const doneTasks = computed({
  get: () => tasks.value.filter(t => t.Status === 'done').sort((a, b) => a.Position - b.Position),
  set: (val) => updateColumnPositions(val, 'done')
})

async function updateColumnPositions(newTasks, status) {
  // Find which task moved
  for (let i = 0; i < newTasks.length; i++) {
    const task = newTasks[i]
    // If status or position changed, trigger update
    const currentTask = tasks.value.find(t => t.ID === task.ID)
    if (currentTask && (currentTask.Status !== status || i > 0 && currentTask.Position <= newTasks[i-1].Position)) {
      const prevId = i > 0 ? newTasks[i-1].ID : ''
      const nextId = i < newTasks.length - 1 ? newTasks[i+1].ID : ''
      await tasksStore.moveTask(task.ID, status, prevId, nextId)
      break
    }
  }
  await tasksStore.fetchTasks()
}

onMounted(() => {
  tasksStore.fetchTasks()
})

async function handleCreateTask() {
  if (!newTaskTitle.value.trim()) return
  await tasksStore.createTask(
    newTaskTitle.value, 
    newTaskDesc.value, 
    newTaskPriority.value, 
    newTaskDeadline.value
  )
  showCreateModal.value = false
  newTaskTitle.value = ''
  newTaskDesc.value = ''
  newTaskPriority.value = 'medium'
  newTaskDeadline.value = ''
  tasksStore.fetchTasks()
}

function getChecklistPercent(task) {
  if (!task.Checklists || task.Checklists.length === 0) return 0
  let total = 0
  let completed = 0
  task.Checklists.forEach(c => {
    c.Items?.forEach(item => {
      total++
      if (item.Completed) completed++
    })
  })
  if (total === 0) return 0
  return Math.round((completed / total) * 100)
}

function getChecklistCount(task) {
  if (!task.Checklists || task.Checklists.length === 0) return ''
  let total = 0
  let completed = 0
  task.Checklists.forEach(c => {
    c.Items?.forEach(item => {
      total++
      if (item.Completed) completed++
    })
  })
  if (total === 0) return ''
  return `${completed}/${total}`
}
</script>

<template>
  <div class="page-container kanban-page">
    <header class="page-header" style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--color-border); padding: 24px 32px;">
      <div>
        <h1 class="page-title" style="font-size: 24px; font-weight: 700; color: var(--color-text-primary);">Tasks Board</h1>
        <p class="page-subtitle" style="font-size: 14px; color: var(--color-text-secondary); margin-top: 4px;">Kelola semua pekerjaan kamu</p>
      </div>
      <button class="btn btn-primary" @click="showCreateModal = true" style="display: flex; align-items: center; gap: 8px;">
        <Plus :size="16" /> Tambah Task
      </button>
    </header>

    <div class="page-content-scroll">
      <div class="kanban-board-container" style="padding: 24px;">
        <div class="kanban-board" style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; align-items: start;">
          
          <!-- Todo Column -->
          <div class="kanban-col">
            <div class="kanban-col-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
              <span class="kanban-col-title" style="font-weight: 600; font-size: 14px; color: var(--color-text-primary);">Todo</span>
              <span class="kanban-col-count" style="background: #F3F4F6; color: var(--color-text-secondary); padding: 2px 8px; border-radius: 20px; font-size: 11px;">{{ todoTasks.length }}</span>
            </div>
            
            <VueDraggable
              v-model="todoTasks"
              group="kanban"
              ghost-class="kanban-ghost"
              class="kanban-list"
              style="min-height: 500px; display: flex; flex-direction: column; gap: 12px;"
            >
              <div 
                v-for="task in todoTasks" 
                :key="task.ID" 
                class="card task-card"
                style="padding: 16px; cursor: grab; position: relative;"
                @click="tasksStore.selectTask(task.ID)"
              >
                <!-- Badges -->
                <div style="display: flex; flex-wrap: wrap; gap: 4px; margin-bottom: 8px;">
                  <span :class="'task-badge priority-' + task.Priority" style="font-size: 10px; font-weight: bold; text-transform: uppercase;">
                    {{ task.Priority }}
                  </span>
                  <span 
                    v-for="lbl in task.Labels" 
                    :key="lbl.ID"
                    style="font-size: 10px; padding: 2px 6px; border-radius: 4px; font-weight: bold;"
                    :style="{ background: lbl.Color + '15', color: lbl.Color }"
                  >
                    {{ lbl.Name.toUpperCase() }}
                  </span>
                </div>

                <!-- Title -->
                <h3 style="font-weight: 600; font-size: 14px; color: var(--color-text-primary); margin-bottom: 6px;">{{ task.Title }}</h3>
                <p v-if="task.Description" style="font-size: 12px; color: var(--color-text-secondary); margin-bottom: 12px; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;">{{ task.Description }}</p>

                <!-- Pomodoro & Checklist stats -->
                <div style="display: flex; gap: 12px; align-items: center; margin-bottom: 12px; flex-wrap: wrap;">
                  <div v-if="task.EstimatePomodoro > 0" style="font-size: 11px; color: var(--color-danger); display: flex; align-items: center; gap: 2px;">
                    <span>🍅</span>
                    <strong>{{ task.CompletedPomodoro }}/{{ task.EstimatePomodoro }}</strong>
                  </div>
                  <div v-if="getChecklistCount(task)" style="font-size: 11px; color: var(--color-text-secondary); display: flex; align-items: center; gap: 2px;">
                    <span>☑</span>
                    <strong>{{ getChecklistPercent(task) }}% ({{ getChecklistCount(task) }})</strong>
                  </div>
                </div>

                <!-- Footer -->
                <div style="display: flex; align-items: center; justify-content: space-between; border-top: 1px solid var(--color-border); padding-top: 8px; font-size: 11px;">
                  <div style="display: flex; align-items: center; gap: 4px; color: var(--color-text-secondary);">
                    <Calendar :size="12" />
                    <span>{{ task.Deadline ? dayjs(task.Deadline).format('DD MMM') : 'No due date' }}</span>
                  </div>
                  <div style="display: flex; align-items: center; gap: 8px; color: var(--color-text-muted);">
                    <span v-if="task.Comments?.length > 0" style="display: flex; align-items: center; gap: 2px;">
                      <MessageSquare :size="12" />{{ task.Comments.length }}
                    </span>
                    <span v-if="task.Attachments?.length > 0" style="display: flex; align-items: center; gap: 2px;">
                      <Paperclip :size="12" />{{ task.Attachments.length }}
                    </span>
                  </div>
                </div>
              </div>
            </VueDraggable>
          </div>

          <!-- In Progress Column -->
          <div class="kanban-col">
            <div class="kanban-col-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
              <span class="kanban-col-title" style="font-weight: 600; font-size: 14px; color: var(--color-text-primary);">Doing</span>
              <span class="kanban-col-count" style="background: #F3F4F6; color: var(--color-text-secondary); padding: 2px 8px; border-radius: 20px; font-size: 11px;">{{ inProgressTasks.length }}</span>
            </div>
            
            <VueDraggable
              v-model="inProgressTasks"
              group="kanban"
              ghost-class="kanban-ghost"
              class="kanban-list"
              style="min-height: 500px; display: flex; flex-direction: column; gap: 12px;"
            >
              <div 
                v-for="task in inProgressTasks" 
                :key="task.ID" 
                class="card task-card"
                style="padding: 16px; cursor: grab; position: relative;"
                @click="tasksStore.selectTask(task.ID)"
              >
                <!-- Badges -->
                <div style="display: flex; flex-wrap: wrap; gap: 4px; margin-bottom: 8px;">
                  <span :class="'task-badge priority-' + task.Priority" style="font-size: 10px; font-weight: bold; text-transform: uppercase;">
                    {{ task.Priority }}
                  </span>
                  <span 
                    v-for="lbl in task.Labels" 
                    :key="lbl.ID"
                    style="font-size: 10px; padding: 2px 6px; border-radius: 4px; font-weight: bold;"
                    :style="{ background: lbl.Color + '15', color: lbl.Color }"
                  >
                    {{ lbl.Name.toUpperCase() }}
                  </span>
                </div>

                <!-- Title -->
                <h3 style="font-weight: 600; font-size: 14px; color: var(--color-text-primary); margin-bottom: 6px;">{{ task.Title }}</h3>
                <p v-if="task.Description" style="font-size: 12px; color: var(--color-text-secondary); margin-bottom: 12px; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;">{{ task.Description }}</p>

                <!-- Pomodoro & Checklist stats -->
                <div style="display: flex; gap: 12px; align-items: center; margin-bottom: 12px; flex-wrap: wrap;">
                  <div v-if="task.EstimatePomodoro > 0" style="font-size: 11px; color: var(--color-danger); display: flex; align-items: center; gap: 2px;">
                    <span>🍅</span>
                    <strong>{{ task.CompletedPomodoro }}/{{ task.EstimatePomodoro }}</strong>
                  </div>
                  <div v-if="getChecklistCount(task)" style="font-size: 11px; color: var(--color-text-secondary); display: flex; align-items: center; gap: 2px;">
                    <span>☑</span>
                    <strong>{{ getChecklistPercent(task) }}% ({{ getChecklistCount(task) }})</strong>
                  </div>
                </div>

                <!-- Footer -->
                <div style="display: flex; align-items: center; justify-content: space-between; border-top: 1px solid var(--color-border); padding-top: 8px; font-size: 11px;">
                  <div style="display: flex; align-items: center; gap: 4px; color: var(--color-text-secondary);">
                    <Calendar :size="12" />
                    <span>{{ task.Deadline ? dayjs(task.Deadline).format('DD MMM') : 'No due date' }}</span>
                  </div>
                  <div style="display: flex; align-items: center; gap: 8px; color: var(--color-text-muted);">
                    <span v-if="task.Comments?.length > 0" style="display: flex; align-items: center; gap: 2px;">
                      <MessageSquare :size="12" />{{ task.Comments.length }}
                    </span>
                    <span v-if="task.Attachments?.length > 0" style="display: flex; align-items: center; gap: 2px;">
                      <Paperclip :size="12" />{{ task.Attachments.length }}
                    </span>
                  </div>
                </div>
              </div>
            </VueDraggable>
          </div>

          <!-- Done Column -->
          <div class="kanban-col">
            <div class="kanban-col-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
              <span class="kanban-col-title" style="font-weight: 600; font-size: 14px; color: var(--color-text-primary);">Done</span>
              <span class="kanban-col-count" style="background: #F3F4F6; color: var(--color-text-secondary); padding: 2px 8px; border-radius: 20px; font-size: 11px;">{{ doneTasks.length }}</span>
            </div>
            
            <VueDraggable
              v-model="doneTasks"
              group="kanban"
              ghost-class="kanban-ghost"
              class="kanban-list"
              style="min-height: 500px; display: flex; flex-direction: column; gap: 12px;"
            >
              <div 
                v-for="task in doneTasks" 
                :key="task.ID" 
                class="card task-card"
                style="padding: 16px; cursor: grab; position: relative;"
                @click="tasksStore.selectTask(task.ID)"
              >
                <!-- Badges -->
                <div style="display: flex; flex-wrap: wrap; gap: 4px; margin-bottom: 8px;">
                  <span :class="'task-badge priority-' + task.Priority" style="font-size: 10px; font-weight: bold; text-transform: uppercase;">
                    {{ task.Priority }}
                  </span>
                  <span 
                    v-for="lbl in task.Labels" 
                    :key="lbl.ID"
                    style="font-size: 10px; padding: 2px 6px; border-radius: 4px; font-weight: bold;"
                    :style="{ background: lbl.Color + '15', color: lbl.Color }"
                  >
                    {{ lbl.Name.toUpperCase() }}
                  </span>
                </div>

                <!-- Title -->
                <h3 style="font-weight: 600; font-size: 14px; color: var(--color-text-primary); margin-bottom: 6px; text-decoration: line-through;">{{ task.Title }}</h3>
                <p v-if="task.Description" style="font-size: 12px; color: var(--color-text-secondary); margin-bottom: 12px; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;">{{ task.Description }}</p>

                <!-- Footer -->
                <div style="display: flex; align-items: center; justify-content: space-between; border-top: 1px solid var(--color-border); padding-top: 8px; font-size: 11px;">
                  <div style="display: flex; align-items: center; gap: 4px; color: var(--color-success);">
                    <CheckCircle2 :size="12" />
                    <span>Selesai</span>
                  </div>
                  <div style="display: flex; align-items: center; gap: 8px; color: var(--color-text-muted);">
                    <span v-if="task.Comments?.length > 0" style="display: flex; align-items: center; gap: 2px;">
                      <MessageSquare :size="12" />{{ task.Comments.length }}
                    </span>
                  </div>
                </div>
              </div>
            </VueDraggable>
          </div>

        </div>
      </div>
    </div>

    <!-- Task Detail Drawer component -->
    <TaskDetailDrawer v-if="isDrawerOpen" />

    <!-- Create Task Modal -->
    <div v-if="showCreateModal" class="modal-backdrop" style="z-index: 1001;">
      <div class="modal" style="padding: 24px; max-width: 500px;">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
          <h2 style="font-size: 18px; font-weight: 700; color: var(--color-text-primary);">Tambah Task Baru</h2>
          <button @click="showCreateModal = false" style="background: none; border: none; cursor: pointer;">
            <X :size="18" />
          </button>
        </div>

        <div style="display: flex; flex-direction: column; gap: 16px;">
          <div style="display: flex; flex-direction: column; gap: 6px;">
            <label style="font-size: 12px; font-weight: 600; color: var(--color-text-secondary);">Judul Task</label>
            <input 
              v-model="newTaskTitle"
              type="text" 
              placeholder="Tulis judul tugas..."
              style="width: 100%; padding: 10px 12px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); outline: none; font-size: 14px;"
            />
          </div>

          <div style="display: flex; flex-direction: column; gap: 6px;">
            <label style="font-size: 12px; font-weight: 600; color: var(--color-text-secondary);">Deskripsi</label>
            <textarea 
              v-model="newTaskDesc"
              placeholder="Tambahkan detail deskripsi..."
              rows="3"
              style="width: 100%; padding: 10px 12px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); outline: none; font-size: 14px; resize: none;"
            ></textarea>
          </div>

          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
            <div style="display: flex; flex-direction: column; gap: 6px;">
              <label style="font-size: 12px; font-weight: 600; color: var(--color-text-secondary);">Prioritas</label>
              <select 
                v-model="newTaskPriority"
                style="padding: 10px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); outline: none; font-size: 14px; background: white;"
              >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </div>

            <div style="display: flex; flex-direction: column; gap: 6px;">
              <label style="font-size: 12px; font-weight: 600; color: var(--color-text-secondary);">Deadline</label>
              <input 
                v-model="newTaskDeadline"
                type="date"
                style="padding: 10px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); outline: none; font-size: 14px;"
              />
            </div>
          </div>
        </div>

        <div style="display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px;">
          <button class="btn btn-secondary" @click="showCreateModal = false">Batal</button>
          <button class="btn btn-primary" @click="handleCreateTask">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>
