<script setup>
import { ref, onMounted } from 'vue'
import client from '../api/client'
import { useTimerStore } from '../stores/timer'
import { useTasksStore } from '../stores/tasks'
import { Play, Calendar, Clock, CheckCircle2, AlertCircle } from 'lucide-vue-next'
import dayjs from 'dayjs'

const stats = ref({
  active_tasks: 0,
  done_tasks: 0,
  today_tasks: 0,
  today_reminders: 0,
  tasks_due_today: [],
  reminders: [],
  daily_recap: null,
  recent_logs: [],
})

const timerStore = useTimerStore()
const tasksStore = useTasksStore()

async function fetchStats() {
  try {
    const res = await client.get('/dashboard')
    stats.value = res.data
  } catch (err) {
    console.error('Gagal memuat dashboard stats', err)
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <div class="page-container">
    <header class="page-header" style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--color-border); padding: 24px 32px;">
      <div>
        <h1 class="page-title" style="font-size: 24px; font-weight: 700; color: var(--color-text-primary);">Dashboard</h1>
        <p class="page-subtitle" style="font-size: 14px; color: var(--color-text-secondary); margin-top: 4px;">Selamat datang di workspace produktivitas Anda.</p>
      </div>
    </header>

    <div class="page-content-scroll" style="padding: 32px; display: flex; flex-direction: column; gap: 32px;">
      <!-- Stats Summary Banner -->
      <section style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px;">
        <div class="card" style="padding: 20px; display: flex; align-items: center; gap: 16px;">
          <div style="font-size: 24px; background: #EFF6FF; color: var(--color-primary); padding: 12px; border-radius: 12px;">📋</div>
          <div>
            <div style="font-size: 24px; font-weight: 700; color: var(--color-text-primary);">{{ stats.active_tasks }}</div>
            <div style="font-size: 12px; color: var(--color-text-secondary); margin-top: 2px;">Task Aktif</div>
          </div>
        </div>

        <div class="card" style="padding: 20px; display: flex; align-items: center; gap: 16px;">
          <div style="font-size: 24px; background: #ECFDF5; color: var(--color-success); padding: 12px; border-radius: 12px;">✅</div>
          <div>
            <div style="font-size: 24px; font-weight: 700; color: var(--color-text-primary);">{{ stats.done_tasks }}</div>
            <div style="font-size: 12px; color: var(--color-text-secondary); margin-top: 2px;">Selesai</div>
          </div>
        </div>

        <div class="card" style="padding: 20px; display: flex; align-items: center; gap: 16px;">
          <div style="font-size: 24px; background: #FFFBEB; color: var(--color-warning); padding: 12px; border-radius: 12px;">⏰</div>
          <div>
            <div style="font-size: 24px; font-weight: 700; color: var(--color-text-primary);">{{ stats.today_reminders }}</div>
            <div style="font-size: 12px; color: var(--color-text-secondary); margin-top: 2px;">Reminder Hari Ini</div>
          </div>
        </div>
      </section>

      <!-- Grid Dashboard -->
      <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px;">
        <!-- Left Side: Tasks Due Today & Recent logs -->
        <div style="display: flex; flex-direction: column; gap: 24px;">
          <section class="card" style="padding: 24px;">
            <h2 style="font-size: 16px; font-weight: 700; margin-bottom: 16px; color: var(--color-text-primary);">Tugas Hari Ini</h2>
            
            <div v-if="stats.tasks_due_today.length === 0" style="padding: 32px; text-align: center; color: var(--color-text-secondary);">
              Tidak ada tugas yang jatuh tempo hari ini.
            </div>
            
            <div v-else style="display: flex; flex-direction: column; gap: 12px;">
              <div 
                v-for="task in stats.tasks_due_today" 
                :key="task.ID"
                class="card"
                style="padding: 16px; display: flex; justify-content: space-between; align-items: center; cursor: pointer; border-color: var(--color-border);"
                @click="tasksStore.selectTask(task.ID)"
              >
                <div>
                  <h3 style="font-weight: 600; font-size: 14px; color: var(--color-text-primary);">{{ task.Title }}</h3>
                  <div style="font-size: 12px; color: var(--color-text-secondary); display: flex; gap: 8px; margin-top: 4px;">
                    <span :class="'priority-badge priority-' + task.Priority">{{ task.Priority }}</span>
                    <span>Status: {{ task.Status }}</span>
                  </div>
                </div>
                <button 
                  class="btn btn-secondary btn-icon"
                  style="border-radius: 50%;"
                  @click.stop="timerStore.startSession(task.ID, '', 'pomodoro')"
                >
                  <Play :size="14" />
                </button>
              </div>
            </div>
          </section>

          <section class="card" style="padding: 24px;">
            <h2 style="font-size: 16px; font-weight: 700; margin-bottom: 16px; color: var(--color-text-primary);">Waktu Kerja Terakhir</h2>
            
            <div v-if="stats.recent_logs.length === 0" style="padding: 32px; text-align: center; color: var(--color-text-secondary);">
              Belum ada log waktu kerja.
            </div>
            
            <table v-else style="width: 100%; border-collapse: collapse;">
              <thead>
                <tr style="border-bottom: 1px solid var(--color-border); text-align: left; font-size: 12px; color: var(--color-text-secondary);">
                  <th style="padding: 12px 8px;">Tugas</th>
                  <th style="padding: 12px 8px;">Mulai</th>
                  <th style="padding: 12px 8px;">Durasi</th>
                  <th style="padding: 12px 8px;">Tipe</th>
                </tr>
              </thead>
              <tbody>
                <tr 
                  v-for="log in stats.recent_logs" 
                  :key="log.ID"
                  style="border-bottom: 1px solid var(--color-border); font-size: 13px;"
                >
                  <td style="padding: 12px 8px; font-weight: 500;">{{ log.TaskTitle || 'Umum' }}</td>
                  <td style="padding: 12px 8px; color: var(--color-text-secondary);">{{ dayjs(log.StartTime).format('DD MMM, HH:mm') }}</td>
                  <td style="padding: 12px 8px; font-weight: 600;">{{ Math.floor(log.DurationSeconds / 60) }}m</td>
                  <td style="padding: 12px 8px;">
                    <span :style="{ color: log.SessionType === 'pomodoro' ? 'var(--color-danger)' : 'var(--color-success)' }">
                      {{ log.SessionType }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </section>
        </div>

        <!-- Right Side: Focus Recap & Reminders -->
        <div style="display: flex; flex-direction: column; gap: 24px;">
          <!-- Daily Focus Goal -->
          <section class="card" style="padding: 24px;">
            <h2 style="font-size: 14px; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 16px; letter-spacing: 0.5px;">Fokus Hari Ini</h2>
            
            <div style="display: flex; flex-direction: column; gap: 16px;">
              <div style="display: flex; gap: 16px; align-items: center;">
                <div style="font-size: 24px;">🍅</div>
                <div>
                  <div style="font-size: 16px; font-weight: 700; color: var(--color-danger);">
                    {{ stats.daily_recap?.TotalFocusSessions || 0 }}<span style="color: var(--color-text-secondary); font-size: 12px;">/5 Sesi</span>
                  </div>
                  <div style="font-size: 11px; color: var(--color-text-secondary);">Sesi Pomodoro Selesai</div>
                </div>
              </div>

              <div style="display: flex; gap: 16px; align-items: center;">
                <div style="font-size: 24px;">⏱️</div>
                <div>
                  <div style="font-size: 16px; font-weight: 700; color: var(--color-primary);">
                    {{ Math.floor((stats.daily_recap?.TotalFocusSeconds || 0) / 3600) }}h {{ Math.floor(((stats.daily_recap?.TotalFocusSeconds || 0) % 3600) / 60) }}m
                  </div>
                  <div style="font-size: 11px; color: var(--color-text-secondary);">Total Waktu Kerja</div>
                </div>
              </div>

              <div style="margin-top: 8px;">
                <div style="height: 6px; background: #E5E7EB; border-radius: 3px; overflow: hidden;">
                  <div 
                    style="height: 100%; background: var(--color-primary);"
                    :style="{ width: Math.min(((stats.daily_recap?.TotalFocusSeconds || 0) / 7200) * 100, 100) + '%' }"
                  ></div>
                </div>
                <div style="font-size: 10px; color: var(--color-text-muted); text-align: center; margin-top: 8px;">
                  Target harian: 2 jam (7200 detik)
                </div>
              </div>
            </div>
          </section>

          <!-- Reminders -->
          <section class="card" style="padding: 24px;">
            <h2 style="font-size: 16px; font-weight: 700; margin-bottom: 16px; color: var(--color-text-primary);">Reminder Aktif</h2>
            
            <div v-if="stats.reminders.length === 0" style="padding: 24px; text-align: center; color: var(--color-text-secondary);">
              Tidak ada reminder hari ini.
            </div>
            
            <div v-else style="display: flex; flex-direction: column; gap: 8px;">
              <div 
                v-for="rem in stats.reminders" 
                :key="rem.ID"
                class="card"
                style="padding: 12px; display: flex; align-items: flex-start; gap: 8px; border-color: var(--color-border);"
              >
                <div style="color: var(--color-warning);">⚠️</div>
                <div style="flex: 1;">
                  <div style="font-size: 13px; font-weight: 500; color: var(--color-text-primary);">{{ rem.Note }}</div>
                  <div style="font-size: 11px; color: var(--color-text-secondary); margin-top: 2px;">
                    Jam: {{ dayjs(rem.RemindAt).format('HH:mm') }} &middot; {{ rem.TaskTitle || 'Umum' }}
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>
