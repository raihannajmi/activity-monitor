<script setup>
import { ref, onMounted } from 'vue'
import client from '../api/client'
import { Calendar, Clock, Activity } from 'lucide-vue-next'
import dayjs from 'dayjs'

const activities = ref([])

async function fetchActivities() {
  try {
    const res = await client.get('/timeline')
    activities.value = res.data || []
  } catch (err) {
    console.error('Gagal memuat timeline', err)
  }
}

onMounted(() => {
  fetchActivities()
})
</script>

<template>
  <div class="page-container">
    <header class="page-header" style="border-bottom: 1px solid var(--color-border); padding: 24px 32px;">
      <h1 class="page-title" style="font-size: 24px; font-weight: 700; color: var(--color-text-primary);">Timeline</h1>
      <p class="page-subtitle" style="font-size: 14px; color: var(--color-text-secondary); margin-top: 4px;">Riwayat aktivitas dan catatan kronologis pekerjaan Anda.</p>
    </header>

    <div class="page-content-scroll" style="padding: 32px; max-width: 800px; margin: 0 auto; width: 100%;">
      <div v-if="activities.length === 0" style="padding: 48px; text-align: center; color: var(--color-text-secondary);">
        Belum ada aktivitas tercatat.
      </div>

      <div v-else style="display: flex; flex-direction: column; gap: 24px; position: relative;">
        <!-- Timeline vertical line visual guide -->
        <div style="position: absolute; left: 19px; top: 8px; bottom: 8px; width: 2px; background: var(--color-border); z-index: 1;"></div>

        <div 
          v-for="act in activities" 
          :key="act.ID"
          style="display: flex; gap: 16px; position: relative; z-index: 2;"
        >
          <!-- Point indicator -->
          <div style="width: 40px; height: 40px; background: white; border: 2px solid var(--color-primary); border-radius: 50%; display: flex; align-items: center; justify-content: center; box-shadow: var(--shadow-sm);">
            <Activity :size="14" style="color: var(--color-primary);" />
          </div>

          <!-- Activity details card -->
          <div class="card" style="flex: 1; padding: 16px; display: flex; flex-direction: column; gap: 6px;">
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span style="font-size: 12px; font-weight: 700; color: var(--color-text-primary); text-transform: uppercase;">
                {{ act.Type }}
              </span>
              <span style="font-size: 11px; color: var(--color-text-secondary); display: flex; align-items: center; gap: 4px;">
                <Clock :size="11" />
                {{ dayjs(act.CreatedAt).format('DD MMM YYYY, HH:mm') }}
              </span>
            </div>
            <p style="font-size: 13px; color: var(--color-text-primary); margin: 0;">{{ act.Description }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
