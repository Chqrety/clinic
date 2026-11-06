import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { QueueItem } from '../api/client'

export type RoleType = 'patient' | 'cashier' | 'doctor' | 'pharmacy' | 'inventory' | 'billing' | 'satusehat'

export const useClinicStore = defineStore('clinic', () => {
  const currentRole = ref<RoleType>('patient')
  const activePatientQueue = ref<QueueItem | null>(null)
  const selectedDoctorQueueId = ref<number | null>(null)
  const soundEnabled = ref(true)

  function setRole(role: RoleType) {
    currentRole.value = role
  }

  function setActivePatientQueue(q: QueueItem | null) {
    activePatientQueue.value = q
  }

  function setSelectedDoctorQueue(queueId: number | null) {
    selectedDoctorQueueId.value = queueId
  }

  function playChime() {
    if (!soundEnabled.value) return
    try {
      const ctx = new (window.AudioContext || (window as any).webkitAudioContext)()
      const osc = ctx.createOscillator()
      const gain = ctx.createGain()
      osc.type = 'sine'
      osc.frequency.setValueAtTime(587.33, ctx.currentTime) // D5
      osc.frequency.exponentialRampToValueAtTime(880, ctx.currentTime + 0.15) // A5
      gain.gain.setValueAtTime(0.15, ctx.currentTime)
      gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.5)
      osc.connect(gain)
      gain.connect(ctx.destination)
      osc.start()
      osc.stop(ctx.currentTime + 0.5)
    } catch (e) {
      // Audio autoplay policy fallback
    }
  }

  return {
    currentRole,
    activePatientQueue,
    selectedDoctorQueueId,
    soundEnabled,
    setRole,
    setActivePatientQueue,
    setSelectedDoctorQueue,
    playChime,
  }
})
