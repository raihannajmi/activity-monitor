<script setup>
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useUIStore } from './stores/ui'
import { useTimerStore } from './stores/timer'
import { useTasksStore } from './stores/tasks'
import { 
  LayoutDashboard, 
  CheckSquare, 
  Clock, 
  PieChart, 
  Edit3, 
  Calendar,
  X 
} from 'lucide-vue-next'
import PomodoroWidget from './components/PomodoroWidget.vue'

const uiStore = useUIStore()
const { toasts, sidebarCollapsed } = storeToRefs(uiStore)

const timerStore = useTimerStore()
const tasksStore = useTasksStore()

onMounted(() => {
  timerStore.checkActiveSession()
  tasksStore.fetchLabels()
})
</script>

<template>
  <div class="layout-wrapper">
    <!-- Desktop/Tablet Sidebar -->
    <aside class="sidebar" :class="{ 'collapsed': sidebarCollapsed }">
      <div class="sidebar-logo">
        <div class="logo-icon">A</div>
        Activity Monitor
      </div>
      <nav class="sidebar-nav">
        <div class="nav-label">MENU</div>
        <router-link to="/" class="nav-item" active-class="active">
          <LayoutDashboard class="nav-icon" :size="18" />
          <span>Dashboard</span>
        </router-link>
        <router-link to="/tasks" class="nav-item" active-class="active">
          <CheckSquare class="nav-icon" :size="18" />
          <span>Tasks</span>
        </router-link>
        <router-link to="/timeline" class="nav-item" active-class="active">
          <Clock class="nav-icon" :size="18" />
          <span>Timeline</span>
        </router-link>
        <router-link to="/reports" class="nav-item" active-class="active">
          <PieChart class="nav-icon" :size="18" />
          <span>Waktu</span>
        </router-link>
        <router-link to="/notes" class="nav-item" active-class="active">
          <Edit3 class="nav-icon" :size="18" />
          <span>Brain Dump</span>
        </router-link>
      </nav>
      
      <!-- Personal Dropdown at bottom of sidebar -->
      <div style="margin-top: auto; display: flex; align-items: center; justify-content: space-between; padding-top: 16px; border-top: 1px solid var(--color-border); cursor: pointer;">
        <div style="display: flex; align-items: center; gap: 8px;">
          <div style="width: 24px; height: 24px; background: #374151; color: white; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 600;">P</div>
          <span style="font-size: 13px; font-weight: 500; color: var(--color-text-primary);">Personal &middot; Private</span>
        </div>
      </div>
    </aside>

    <!-- Main Content Panel -->
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is={Component} />
        </transition>
      </router-view>
    </main>

    <!-- Mobile Bottom Navigation -->
    <nav class="bottom-nav">
      <router-link to="/" class="bottom-nav-item" active-class="active">
        <LayoutDashboard class="bottom-nav-icon" :size="18" />
        <span class="bottom-nav-label">Dashboard</span>
      </router-link>
      <router-link to="/tasks" class="bottom-nav-item" active-class="active">
        <CheckSquare class="bottom-nav-icon" :size="18" />
        <span class="bottom-nav-label">Tasks</span>
      </router-link>
      <router-link to="/notes" class="bottom-nav-item" active-class="active">
        <Edit3 class="bottom-nav-icon" :size="18" />
        <span class="bottom-nav-label">Catatan</span>
      </router-link>
      <router-link to="/timeline" class="bottom-nav-item" active-class="active">
        <Clock class="bottom-nav-icon" :size="18" />
        <span class="bottom-nav-label">Timeline</span>
      </router-link>
      <router-link to="/reports" class="bottom-nav-item" active-class="active">
        <PieChart class="bottom-nav-icon" :size="18" />
        <span class="bottom-nav-label">Waktu</span>
      </router-link>
    </nav>

    <!-- Floating Pomodoro Widget -->
    <PomodoroWidget />

    <!-- Toast Notification Overlay -->
    <div class="toast-container" id="toast-container">
      <transition-group name="toast">
        <div 
          v-for="toast in toasts" 
          :key="toast.id"
          class="toast"
          :class="'toast-' + toast.type"
        >
          <div class="toast-icon">
            <span v-if="toast.type === 'success'">✓</span>
            <span v-else>⚠️</span>
          </div>
          <div class="toast-content">
            <div class="toast-title">{{ toast.type === 'success' ? 'Berhasil' : 'Kesalahan' }}</div>
            <div class="toast-message">{{ toast.message }}</div>
          </div>
          <button class="toast-close" @click="uiStore.removeToast(toast.id)">
            <X :size="14" />
          </button>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Sidebar collapsed helper rules */
.sidebar.collapsed {
  width: 72px;
}
.sidebar.collapsed .sidebar-logo {
  font-size: 0;
}
.sidebar.collapsed .sidebar-nav .nav-label,
.sidebar.collapsed .sidebar-nav span,
.sidebar.collapsed div:last-child {
  display: none !important;
}
</style>
