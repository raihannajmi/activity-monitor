<script setup>
import { ref } from 'vue'
import { useTasksStore } from '../stores/tasks'
import { useTimerStore } from '../stores/timer'
import { storeToRefs } from 'pinia'
import { 
  X, 
  Trash2, 
  Plus, 
  Play, 
  Clock, 
  PlusCircle, 
  CheckCircle2, 
  User, 
  Smile, 
  Send 
} from 'lucide-vue-next'
import dayjs from 'dayjs'

const tasksStore = useTasksStore()
const { selectedTask, labels } = storeToRefs(tasksStore)
const timerStore = useTimerStore()

const newSubtaskTitle = ref('')
const newSubtaskDeadline = ref('')

const newChecklistTitle = ref('')
const newChecklistItemTitles = ref({}) // map of checklistId -> inputTitle

const newCommentContent = ref('')

// Auto-save updates logic
function saveField(field, val) {
  if (!selectedTask.value) return
  tasksStore.updateTask(selectedTask.value.ID, { [field]: val })
}

async function handleAddSubtask() {
  if (!newSubtaskTitle.value.trim()) return
  await tasksStore.addSubtask(newSubtaskTitle.value, newSubtaskDeadline.value)
  newSubtaskTitle.value = ''
  newSubtaskDeadline.value = ''
}

async function handleAddChecklist() {
  const title = newChecklistTitle.value.trim() || 'Checklist'
  await tasksStore.addChecklist(title)
  newChecklistTitle.value = ''
}

async function handleAddChecklistItem(checklistId) {
  const title = newChecklistItemTitles.value[checklistId]?.trim()
  if (!title) return
  await tasksStore.addChecklistItem(checklistId, title)
  newChecklistItemTitles.value[checklistId] = ''
}

async function handleAddComment() {
  if (!newCommentContent.value.trim()) return
  await tasksStore.addComment(newCommentContent.value)
  newCommentContent.value = ''
}
</script>

<template>
  <div v-if="selectedTask" class="modal-backdrop" style="justify-content: flex-end; align-items: stretch; padding: 0; z-index: 1000;" @click.self="tasksStore.closeDrawer">
    <div 
      class="task-detail-drawer" 
      style="width: 550px; max-width: 95vw; background: white; height: 100vh; display: flex; flex-direction: column; box-shadow: var(--shadow-lg); animation: slideInRight 0.2s ease;"
    >
      <!-- Header -->
      <header style="padding: 20px 24px; border-bottom: 1px solid var(--color-border); display: flex; justify-content: space-between; align-items: center;">
        <div style="display: flex; align-items: center; gap: 12px;">
          <select 
            :value="selectedTask.Status"
            @change="saveField('status', $event.target.value)"
            style="padding: 6px 12px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: 13px; font-weight: 600; outline: none; background: white;"
          >
            <option value="todo">Todo</option>
            <option value="in_progress">Doing</option>
            <option value="done">Done</option>
          </select>
          <span :class="'priority-badge priority-' + selectedTask.Priority">{{ selectedTask.Priority }}</span>
        </div>
        <div style="display: flex; align-items: center; gap: 8px;">
          <button class="btn btn-ghost btn-danger-ghost btn-sm" @click="tasksStore.deleteTask(selectedTask.ID)">
            <Trash2 :size="16" />
          </button>
          <button @click="tasksStore.closeDrawer" style="background: none; border: none; cursor: pointer; padding: 4px;">
            <X :size="20" />
          </button>
        </div>
      </header>

      <!-- Scrollable content -->
      <div style="flex: 1; overflow-y: auto; padding: 24px; display: flex; flex-direction: column; gap: 24px;">
        
        <!-- Title Input -->
        <input 
          type="text" 
          :value="selectedTask.Title"
          @blur="saveField('title', $event.target.value)"
          style="font-size: 20px; font-weight: 700; border: none; width: 100%; outline: none; color: var(--color-text-primary);"
        />

        <!-- Description Textarea -->
        <div style="display: flex; flex-direction: column; gap: 8px;">
          <label style="font-size: 11px; font-weight: bold; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.5px;">Deskripsi</label>
          <textarea 
            :value="selectedTask.Description"
            @blur="saveField('description', $event.target.value)"
            placeholder="Tambahkan detail deskripsi untuk tugas ini..."
            rows="3"
            style="width: 100%; border: 1px solid var(--color-border); border-radius: var(--radius-sm); padding: 10px 12px; font-size: 13px; resize: none; outline: none;"
          ></textarea>
        </div>

        <!-- Pomodoro sync configuration -->
        <div style="display: flex; flex-direction: column; gap: 8px; background: #FFF5F5; border: 1px solid #FFE3E3; border-radius: var(--radius-md); padding: 16px;">
          <h3 style="font-size: 12px; font-weight: 700; color: #E53E3E; margin-bottom: 8px; display: flex; align-items: center; gap: 4px;">
            <span>🍅</span> Target Waktu Pomodoro
          </h3>
          <div style="display: flex; gap: 16px; align-items: center;">
            <div style="display: flex; align-items: center; gap: 8px;">
              <span style="font-size: 13px; color: var(--color-text-secondary);">Selesai:</span>
              <input 
                type="number" 
                :value="selectedTask.CompletedPomodoro"
                @blur="saveField('completed_pomodoro', parseInt($event.target.value) || 0)"
                style="width: 60px; padding: 6px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); text-align: center;"
              />
            </div>
            <div style="display: flex; align-items: center; gap: 8px;">
              <span style="font-size: 13px; color: var(--color-text-secondary);">Target:</span>
              <input 
                type="number" 
                :value="selectedTask.EstimatePomodoro"
                @blur="saveField('estimate_pomodoro', parseInt($event.target.value) || 0)"
                style="width: 60px; padding: 6px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); text-align: center;"
              />
            </div>
            <button 
              class="btn btn-danger btn-sm"
              style="margin-left: auto; display: flex; align-items: center; gap: 4px;"
              @click="timerStore.startSession(selectedTask.ID, '', 'pomodoro')"
            >
              <Play :size="12" /> Start Focus
            </button>
          </div>
        </div>

        <!-- Labels selection Checklist -->
        <div style="display: flex; flex-direction: column; gap: 8px;">
          <label style="font-size: 11px; font-weight: bold; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.5px;">Labels</label>
          <div style="display: flex; flex-wrap: wrap; gap: 8px;">
            <label 
              v-for="lbl in labels" 
              :key="lbl.ID"
              style="display: flex; align-items: center; gap: 6px; padding: 4px 10px; border-radius: 20px; font-size: 12px; cursor: pointer; border: 1px solid var(--color-border);"
              :style="{ 
                borderColor: selectedTask.Labels?.some(l => l.ID === lbl.ID) ? lbl.Color : 'var(--color-border)',
                background: selectedTask.Labels?.some(l => l.ID === lbl.ID) ? lbl.Color + '15' : 'white',
                color: selectedTask.Labels?.some(l => l.ID === lbl.ID) ? lbl.Color : 'var(--color-text-secondary)'
              }"
            >
              <input 
                type="checkbox"
                :checked="selectedTask.Labels?.some(l => l.ID === lbl.ID)"
                @change="tasksStore.toggleLabel(lbl.ID)"
                style="display: none;"
              />
              <span>{{ lbl.Name }}</span>
            </label>
          </div>
        </div>

        <!-- Checklist Section -->
        <div style="display: flex; flex-direction: column; gap: 16px;">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <label style="font-size: 11px; font-weight: bold; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.5px;">Checklists</label>
            <div style="display: flex; gap: 8px;">
              <input 
                v-model="newChecklistTitle"
                type="text" 
                placeholder="Nama checklist..."
                style="padding: 4px 8px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: 12px;"
                @keyup.enter="handleAddChecklist"
              />
              <button class="btn btn-secondary btn-sm" @click="handleAddChecklist">Tambah</button>
            </div>
          </div>

          <div v-for="c in selectedTask.Checklists" :key="c.ID" class="card" style="padding: 16px; display: flex; flex-direction: column; gap: 12px;">
            <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--color-border); padding-bottom: 8px;">
              <h4 style="font-weight: 700; font-size: 13px; color: var(--color-text-primary);">{{ c.Title }}</h4>
              <button style="background: none; border: none; cursor: pointer; color: var(--color-danger);" @click="tasksStore.deleteChecklist(c.ID)">
                <Trash2 :size="12" />
              </button>
            </div>

            <!-- List items -->
            <div style="display: flex; flex-direction: column; gap: 8px;">
              <div 
                v-for="item in c.Items" 
                :key="item.ID"
                style="display: flex; align-items: center; gap: 8px; font-size: 13px;"
              >
                <input 
                  type="checkbox"
                  :checked="item.Completed"
                  @change="tasksStore.toggleChecklistItem(c.ID, item)"
                />
                <span :style="{ textDecoration: item.Completed ? 'line-through' : 'none', color: item.Completed ? 'var(--color-text-muted)' : 'var(--color-text-primary)' }">
                  {{ item.Title }}
                </span>
                <button 
                  style="margin-left: auto; background: none; border: none; cursor: pointer; color: var(--color-text-muted);"
                  @click="tasksStore.deleteChecklistItem(c.ID, item.ID)"
                >
                  ✕
                </button>
              </div>
            </div>

            <!-- Add checklist item -->
            <div style="display: flex; gap: 8px; margin-top: 4px;">
              <input 
                v-model="newChecklistItemTitles[c.ID]"
                type="text" 
                placeholder="+ Tambah item checklist..."
                style="flex: 1; padding: 6px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: 12px;"
                @keyup.enter="handleAddChecklistItem(c.ID)"
              />
            </div>
          </div>
        </div>

        <!-- Subtasks Section -->
        <div style="display: flex; flex-direction: column; gap: 12px;">
          <label style="font-size: 11px; font-weight: bold; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.5px;">Subtasks</label>
          
          <div style="display: flex; flex-direction: column; gap: 8px;">
            <div 
              v-for="sub in selectedTask.Subtasks" 
              :key="sub.ID"
              style="display: flex; align-items: center; gap: 8px; padding: 8px 12px; background: #F9FAFB; border-radius: var(--radius-sm); font-size: 13px;"
            >
              <input 
                type="checkbox" 
                :checked="sub.IsCompleted"
                @change="tasksStore.toggleSubtask(sub.ID)"
              />
              <span :style="{ textDecoration: sub.IsCompleted ? 'line-through' : 'none', color: sub.IsCompleted ? 'var(--color-text-muted)' : 'var(--color-text-primary)' }">
                {{ sub.Title }}
              </span>
              <button 
                v-if="!sub.IsCompleted"
                style="margin-left: auto; background: none; border: none; cursor: pointer;"
                @click="timerStore.startSession(selectedTask.ID, sub.ID, 'pomodoro')"
              >
                <Play :size="12" style="color: var(--color-success);" />
              </button>
              <button 
                style="background: none; border: none; cursor: pointer; color: var(--color-danger);"
                @click="tasksStore.deleteSubtask(sub.ID)"
              >
                <Trash2 :size="12" />
              </button>
            </div>
          </div>

          <!-- Add subtask form inline -->
          <div style="display: flex; gap: 8px; margin-top: 4px;">
            <input 
              v-model="newSubtaskTitle"
              type="text" 
              placeholder="+ Tambah subtask baru..."
              style="flex: 1; padding: 8px; border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: 12px;"
              @keyup.enter="handleAddSubtask"
            />
            <button class="btn btn-secondary btn-sm" @click="handleAddSubtask">Tambah</button>
          </div>
        </div>

        <!-- Comments Section -->
        <div style="display: flex; flex-direction: column; gap: 16px; border-top: 1px solid var(--color-border); padding-top: 24px;">
          <label style="font-size: 11px; font-weight: bold; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.5px;">Diskusi</label>
          
          <!-- Comment entry form -->
          <div style="display: flex; gap: 8px;">
            <input 
              v-model="newCommentContent"
              type="text" 
              placeholder="Tulis tanggapan atau komentar..."
              style="flex: 1; padding: 8px 12px; border: 1px solid var(--color-border); border-radius: 20px; font-size: 13px;"
              @keyup.enter="handleAddComment"
            />
            <button class="btn btn-primary btn-icon" style="border-radius: 50%;" @click="handleAddComment">
              <Send :size="14" />
            </button>
          </div>

          <!-- Comment list -->
          <div style="display: flex; flex-direction: column; gap: 12px;">
            <div 
              v-for="comm in selectedTask.Comments" 
              :key="comm.ID"
              style="display: flex; gap: 8px; font-size: 13px;"
            >
              <div style="width: 24px; height: 24px; background: #E5E7EB; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 10px; font-weight: bold; color: var(--color-text-secondary);">
                U
              </div>
              <div style="flex: 1; background: #F3F4F6; padding: 8px 12px; border-radius: 12px;">
                <div style="font-weight: 600; font-size: 11px; color: var(--color-text-primary); margin-bottom: 2px;">
                  User &middot; <span style="font-weight: normal; color: var(--color-text-secondary);">{{ dayjs(comm.CreatedAt).format('DD MMM HH:mm') }}</span>
                </div>
                <div style="color: var(--color-text-primary); line-height: 1.4;">{{ comm.Content }}</div>
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes slideInRight {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}
</style>
