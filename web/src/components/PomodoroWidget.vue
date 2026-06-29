<script setup>
import { useTimerStore } from '../stores/timer'
import { storeToRefs } from 'pinia'
import { Play, Square, Timer } from 'lucide-vue-next'

const timerStore = useTimerStore()
const { activeSession, isRunning, formattedTime, progressPercent } = storeToRefs(timerStore)
</script>

<template>
  <div v-if="activeSession" class="pomodoro-widget" style="z-index: 999; display: flex; align-items: center; gap: 16px; box-shadow: var(--shadow-lg); padding: 12px 20px;">
    <div style="display: flex; flex-direction: column; gap: 4px;">
      <div style="font-size: 11px; font-weight: bold; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.5px; display: flex; align-items: center; gap: 4px;">
        <span style="color: var(--color-danger);">🍅</span> 
        <span>Fokus Session</span>
      </div>
      <div style="font-size: 24px; font-weight: 700; font-family: monospace; color: var(--color-text-primary); line-height: 1;">
        {{ formattedTime }}
      </div>
    </div>
    
    <div style="display: flex; gap: 8px;">
      <button 
        v-if="!isRunning" 
        class="btn btn-primary btn-icon btn-sm"
        @click="timerStore.startSession(activeSession.TaskID, activeSession.SubtaskID, activeSession.SessionType)"
      >
        <Play :size="14" />
      </button>
      <button 
        class="btn btn-danger btn-icon btn-sm"
        @click="timerStore.stopSession"
      >
        <Square :size="14" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.pomodoro-widget {
  position: fixed;
  bottom: 24px;
  right: 24px;
  background: white;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
}

@media (max-width: 768px) {
  .pomodoro-widget {
    bottom: 96px;
    right: 16px;
    left: 16px;
    justify-content: space-between;
  }
}
</style>
