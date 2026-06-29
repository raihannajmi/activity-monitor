import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import client from '../api/client'
import { useUIStore } from './ui'

export const useTimerStore = defineStore('timer', () => {
  const ui = useUIStore()
  const activeSession = ref(null)
  const isRunning = ref(false)
  const timeLeft = ref(25 * 60) // default 25 min
  const initialTime = ref(25 * 60)
  const timerInterval = ref(null)

  const formattedTime = computed(() => {
    const minutes = Math.floor(timeLeft.value / 60)
    const seconds = timeLeft.value % 60
    return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
  })

  const progressPercent = computed(() => {
    if (initialTime.value === 0) return 0
    return ((initialTime.value - timeLeft.value) / initialTime.value) * 100
  })

  async function checkActiveSession() {
    try {
      const res = await client.get('/timer/active')
      if (res.data) {
        activeSession.value = res.data
        const start = new Date(res.data.StartTime).getTime()
        const elapsed = Math.floor((Date.now() - start) / 1000)
        
        let sessionDuration = 25 * 60
        if (res.data.SessionType === 'short_break') sessionDuration = 5 * 60
        else if (res.data.SessionType === 'long_break') sessionDuration = 15 * 60

        initialTime.value = sessionDuration
        if (elapsed < sessionDuration) {
          timeLeft.value = sessionDuration - elapsed
          startLocalTimer()
        } else {
          // auto stopped
          timeLeft.value = 0
          isRunning.value = false
        }
      } else {
        activeSession.value = null
        isRunning.value = false
      }
    } catch (err) {
      console.error(err)
    }
  }

  function startLocalTimer() {
    isRunning.value = true
    if (timerInterval.value) clearInterval(timerInterval.value)
    timerInterval.value = setInterval(() => {
      if (timeLeft.value > 0) {
        timeLeft.value--
      } else {
        stopSession()
      }
    }, 1000)
  }

  async function startSession(taskId = '', subtaskId = '', type = 'pomodoro') {
    try {
      const res = await client.post('/timer/start', {
        task_id: taskId,
        subtask_id: subtaskId,
        session_type: type
      })
      activeSession.value = res.data
      let dur = 25 * 60
      if (type === 'short_break') dur = 5 * 60
      else if (type === 'long_break') dur = 15 * 60
      timeLeft.value = dur
      initialTime.value = dur
      startLocalTimer()
      ui.showToast('Timer fokus dimulai!', 'success')
    } catch (err) {
      ui.showToast('Gagal memulai timer', 'error')
    }
  }

  async function stopSession() {
    if (timerInterval.value) {
      clearInterval(timerInterval.value)
      timerInterval.value = null
    }
    if (!activeSession.value) return

    const duration = initialTime.value - timeLeft.value
    try {
      await client.post(`/timer/${activeSession.value.ID}/stop`, {
        duration_seconds: duration
      })
      ui.showToast('Timer fokus selesai/dihentikan!', 'success')
    } catch (err) {
      ui.showToast('Gagal mencatat sesi timer', 'error')
    } finally {
      activeSession.value = null
      isRunning.value = false
      timeLeft.value = 25 * 60
      initialTime.value = 25 * 60
    }
  }

  return {
    activeSession,
    isRunning,
    timeLeft,
    initialTime,
    formattedTime,
    progressPercent,
    checkActiveSession,
    startSession,
    stopSession,
  }
})
