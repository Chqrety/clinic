/**
 * Audio & Voice Synthesizer Announcement Engine
 * Combines harmonic chime (Web Audio API) with Indonesian Text-to-Speech (Web Speech API)
 */

function spellQueueNumber(qNum: string): string {
  // Converts "A-001" or "A001" to "A, nol, nol, satu" for natural Indonesian pronunciation
  const clean = qNum.replace(/[^A-Za-z0-9]/g, '').toUpperCase()
  const parts: string[] = []

  const digitMap: Record<string, string> = {
    '0': 'nol',
    '1': 'satu',
    '2': 'dua',
    '3': 'tiga',
    '4': 'empat',
    '5': 'lima',
    '6': 'enam',
    '7': 'tujuh',
    '8': 'delapan',
    '9': 'sembilan',
  }

  for (const char of clean) {
    if (digitMap[char]) {
      parts.push(digitMap[char])
    } else {
      parts.push(char)
    }
  }

  return parts.join(', ')
}

export function playChimeAudio(): Promise<void> {
  return new Promise((resolve) => {
    try {
      const AudioCtx = window.AudioContext || (window as any).webkitAudioContext
      if (!AudioCtx) {
        resolve()
        return
      }
      const ctx = new AudioCtx()
      
      // Tone 1: D5 (587.33 Hz)
      const osc1 = ctx.createOscillator()
      const gain1 = ctx.createGain()
      osc1.type = 'sine'
      osc1.frequency.setValueAtTime(587.33, ctx.currentTime)
      gain1.gain.setValueAtTime(0.18, ctx.currentTime)
      gain1.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.35)
      osc1.connect(gain1)
      gain1.connect(ctx.destination)
      osc1.start(ctx.currentTime)
      osc1.stop(ctx.currentTime + 0.35)

      // Tone 2: A5 (880 Hz) harmonic
      const osc2 = ctx.createOscillator()
      const gain2 = ctx.createGain()
      osc2.type = 'sine'
      osc2.frequency.setValueAtTime(880, ctx.currentTime + 0.18)
      gain2.gain.setValueAtTime(0.22, ctx.currentTime + 0.18)
      gain2.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.6)
      osc2.connect(gain2)
      gain2.connect(ctx.destination)
      osc2.start(ctx.currentTime + 0.18)
      osc2.stop(ctx.currentTime + 0.6)

      setTimeout(() => {
        resolve()
      }, 500)
    } catch {
      resolve()
    }
  })
}

export async function announcePatientCall(
  queueNumber: string,
  roomName: string,
  _doctorName?: string,
  soundEnabled = true
): Promise<void> {
  if (!soundEnabled) return

  // 1. Play Harmonic Chime
  await playChimeAudio()

  // 2. Text-to-Speech (Indonesian)
  if (!('speechSynthesis' in window)) return

  try {
    window.speechSynthesis.cancel() // Cancel previous utterance if any

    const spelled = spellQueueNumber(queueNumber)
    const roomClean = roomName || 'Ruang Periksa Dokter'
    const speechText = `Nomor antrian, ${spelled}. Silakan menuju, ${roomClean}.`

    const utterance = new SpeechSynthesisUtterance(speechText)
    utterance.lang = 'id-ID'
    utterance.rate = 0.92 // slightly calm and clear pace
    utterance.pitch = 1.05

    // Detect Indonesian Voice if installed in OS
    const voices = window.speechSynthesis.getVoices()
    const idVoice = voices.find((v) => v.lang.startsWith('id') || v.lang.startsWith('ID'))
    if (idVoice) {
      utterance.voice = idVoice
    }

    window.speechSynthesis.speak(utterance)
  } catch (err) {
    console.warn('Voice announcement error:', err)
  }
}
