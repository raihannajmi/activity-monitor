<script setup>
import { ref, onMounted, computed } from 'vue'
import client from '../api/client'
import { Calendar, Clock, BarChart3, ChevronRight } from 'lucide-vue-next'
import dayjs from 'dayjs'

const logs = ref([])

async function fetchLogs() {
  try {
    const res = await client.get('/reports')
    logs.value = res.data || []
  } catch (err) {
    console.error('Gagal memuat log laporan waktu', err)
  }
}

const totalFocusHours = computed(() => {
  const totalSecs = logs.value
    .filter(l => l.SessionType === 'pomodoro')
    .reduce((acc, curr) => acc + curr.DurationSeconds, 0)
  return (totalSecs / 3600).toFixed(1)
})

const totalBreakHours = computed(() => {
  const totalSecs = logs.value
    .filter(l => l.SessionType !== 'pomodoro')
    .reduce((acc, curr) => acc + curr.DurationSeconds, 0)
  return (totalSecs / 3600).toFixed(1)
})

const sessionTypeCounts = computed(() => {
  const counts = { pomodoro: 0, break: 0 }
  logs.value.forEach(l => {
    if (l.SessionType === 'pomodoro') counts.pomodoro++
    else counts.break++
  })
  return counts
})

onMounted(() => {
  fetchLogs()
})
</script>

<template>
  <div class="page-container">
    <header class="page-header" style="border-bottom: 1px solid var(--color-border); padding: 24px 32px;">
      <h1 class="page-title" style="font-size: 24px; font-weight: 700; color: var(--color-text-primary);">Laporan Waktu</h1>
      <p class="page-subtitle" style="font-size: 14px; color: var(--color-text-secondary); margin-top: 4px;">Analisis detail alokasi waktu kerja dan istirahat Anda.</p>
    </header>

    <div class="page-content-scroll" style="padding: 32px; display: flex; flex-direction: column; gap: 32px;">
      
      <!-- Summaries Cards -->
      <section style="display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 16px;">
        <div class="card" style="padding: 24px; display: flex; align-items: center; gap: 16px; border-left: 4px solid var(--color-danger);">
          <div style="background: #FFF5F5; padding: 12px; border-radius: 12px; font-size: 24px;">⏱️</div>
          <div>
            <div style="font-size: 28px; font-weight: 700; color: var(--color-text-primary);">{{ totalFocusHours }}h</div>
            <div style="font-size: 12px; color: var(--color-text-secondary); margin-top: 2px;">Total Waktu Fokus ({{ sessionTypeCounts.pomodoro }} Sesi)</div>
          </div>
        </div>

        <div class="card" style="padding: 24px; display: flex; align-items: center; gap: 16px; border-left: 4px solid var(--color-success);">
          <div style="background: #EBFDF5; padding: 12px; border-radius: 12px; font-size: 24px;">🏝️</div>
          <div>
            <div style="font-size: 28px; font-weight: 700; color: var(--color-text-primary);">{{ totalBreakHours }}h</div>
            <div style="font-size: 12px; color: var(--color-text-secondary); margin-top: 2px;">Total Waktu Istirahat ({{ sessionTypeCounts.break }} Sesi)</div>
          </div>
        </div>
      </section>

      <!-- Logs list table -->
      <section class="card" style="padding: 24px;">
        <h2 style="font-size: 16px; font-weight: 700; margin-bottom: 16px; color: var(--color-text-primary);">Semua Sesi Waktu</h2>

        <div v-if="logs.length === 0" style="padding: 48px; text-align: center; color: var(--color-text-secondary);">
          Belum ada rekaman sesi waktu kerja.
        </div>

        <table v-else style="width: 100%; border-collapse: collapse; text-align: left;">
          <thead>
            <tr style="border-bottom: 1px solid var(--color-border); font-size: 12px; color: var(--color-text-secondary);">
              <th style="padding: 12px 8px;">Tugas</th>
              <th style="padding: 12px 8px;">Subtask</th>
              <th style="padding: 12px 8px;">Mulai Sesi</th>
              <th style="padding: 12px 8px;">Selesai Sesi</th>
              <th style="padding: 12px 8px;">Durasi Sesi</th>
              <th style="padding: 12px 8px;">Tipe</th>
            </tr>
          </thead>
          <tbody>
            <tr 
              v-for="log in logs" 
              :key="log.ID"
              style="border-bottom: 1px solid var(--color-border); font-size: 13px; color: var(--color-text-primary);"
            >
              <td style="padding: 12px 8px; font-weight: 600;">{{ log.TaskTitle || 'Umum / Tidak Berkaitan' }}</td>
              <td style="padding: 12px 8px; color: var(--color-text-secondary);">{{ log.SubtaskTitle || '-' }}</td>
              <td style="padding: 12px 8px;">{{ dayjs(log.StartTime).format('DD MMM YYYY, HH:mm') }}</td>
              <td style="padding: 12px 8px;">{{ log.EndTime ? dayjs(log.EndTime).format('DD MMM YYYY, HH:mm') : '-' }}</td>
              <td style="padding: 12px 8px; font-weight: 700;">{{ Math.floor(log.DurationSeconds / 60) }}m {{ log.DurationSeconds % 60 }}s</td>
              <td style="padding: 12px 8px;">
                <span 
                  style="font-size: 11px; padding: 2px 6px; border-radius: 4px; font-weight: bold;"
                  :style="{ 
                    background: log.SessionType === 'pomodoro' ? '#FFF5F5' : '#EBFDF5', 
                    color: log.SessionType === 'pomodoro' ? 'var(--color-danger)' : 'var(--color-success)' 
                  }"
                >
                  {{ log.SessionType.toUpperCase() }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>
  </div>
</template>
