<script setup lang="ts">
import { ref } from 'vue'
import { formatDate, formatDateTime } from '@/utils/date'
import {
  ShieldAlert,
  Search,
  CheckCircle2,
  Lock,
  UserCheck,
  FileCheck2,
  RotateCcw
} from 'lucide-vue-next'

interface MockAuditLog {
  id: string
  timestamp: string
  userName: string
  userRole: string
  action: string
  actionType: 'create' | 'review' | 'approve' | 'post' | 'reverse' | 'lock'
  entityReference: string
  description: string
  ipAddress: string
}

const auditLogs = ref<MockAuditLog[]>([
  {
    id: 'log-101',
    timestamp: '2026-03-04 10:10:00',
    userName: 'Rian Pratama, S.Ak.',
    userRole: 'Accounting Staff',
    action: 'CREATE_MANUAL_JOURNAL',
    actionType: 'create',
    entityReference: 'JV-2026-03-007',
    description: 'Membuat draft jurnal penyesuaian biaya administrasi bank Rp 75.000',
    ipAddress: '10.20.30.141'
  },
  {
    id: 'log-102',
    timestamp: '2026-03-03 10:01:00',
    userName: 'Budi Santoso, SE, Ak., CA',
    userRole: 'Finance Manager',
    action: 'POST_JOURNAL_TO_GL',
    actionType: 'post',
    entityReference: 'JV-2026-03-005',
    description: 'Posting jurnal penyusutan bulanan USG 4D Kidinest Rp 4.166.667 ke Buku Besar',
    ipAddress: '10.20.30.100'
  },
  {
    id: 'log-103',
    timestamp: '2026-03-03 10:00:00',
    userName: 'Budi Santoso, SE, Ak., CA',
    userRole: 'Finance Manager',
    action: 'APPROVE_JOURNAL',
    actionType: 'approve',
    entityReference: 'JV-2026-03-005',
    description: 'Menyetujui (Approve) voucher jurnal penyusutan aset tetap',
    ipAddress: '10.20.30.100'
  },
  {
    id: 'log-104',
    timestamp: '2026-03-03 09:30:00',
    userName: 'Dewi Lestari, SE, Ak.',
    userRole: 'Accounting Supervisor',
    action: 'REVIEW_JOURNAL',
    actionType: 'review',
    entityReference: 'JV-2026-03-005',
    description: 'Verifikasi kesesuaian tabel depresiasi garis lurus USG 4D 5 tahun',
    ipAddress: '10.20.30.115'
  },
  {
    id: 'log-105',
    timestamp: '2026-03-02 09:18:30',
    userName: 'Budi Santoso, SE, Ak., CA',
    userRole: 'Finance Manager',
    action: 'POST_JOURNAL_TO_GL',
    actionType: 'post',
    entityReference: 'JV-2026-03-002',
    description: 'Posting resmi closing kasir shift harian SSC Rp 27.800.000',
    ipAddress: '10.20.30.100'
  },
  {
    id: 'log-106',
    timestamp: '2026-02-05 18:00:00',
    userName: 'Budi Santoso, SE, Ak., CA',
    userRole: 'Finance Manager',
    action: 'LOCK_ACCOUNTING_PERIOD',
    actionType: 'lock',
    entityReference: 'PERIODE-2026-01',
    description: 'Tutup buku dan penguncian anti-tamper periode Januari 2026',
    ipAddress: '10.20.30.100'
  }
])

const searchQuery = ref('')
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Compliance & Security
          </span>
          <span class="text-xs text-slate-400">&bull; Immutable Event Logging</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Jejak Audit & Riwayat Aktivitas Sistem
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Log permanen mencatat setiap pembuatan draft, review supervisor, approval manager, posting buku besar, hingga penguncian periode.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <span class="px-3 py-1.5 rounded-lg bg-slate-900 text-white text-xs font-mono font-bold">
          SHA-256 Audit Trail
        </span>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3">Waktu Kejadian</th>
              <th class="px-4 py-3">Pengguna & Peran</th>
              <th class="px-4 py-3">Tindakan Sistem (Action)</th>
              <th class="px-4 py-3">Referensi Entitas</th>
              <th class="px-4 py-3">Keterangan Aktivitas</th>
              <th class="px-4 py-3 font-mono">IP Address</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="log in auditLogs"
              :key="log.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5 text-slate-600 whitespace-nowrap font-mono text-[11px]">
                {{ log.timestamp }}
              </td>
              <td class="px-4 py-3.5">
                <span class="font-bold text-slate-900 block">{{ log.userName }}</span>
                <span class="text-[10px] text-slate-500">{{ log.userRole }}</span>
              </td>
              <td class="px-4 py-3.5">
                <span
                  :class="[
                    'px-2 py-0.5 rounded font-mono font-bold text-[10px] inline-block',
                    log.actionType === 'post'
                      ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                      : log.actionType === 'approve'
                      ? 'bg-teal-50 text-teal-700 border border-teal-200'
                      : log.actionType === 'review'
                      ? 'bg-sky-50 text-sky-700 border border-sky-200'
                      : log.actionType === 'lock'
                      ? 'bg-slate-900 text-white'
                      : 'bg-blue-50 text-blue-700 border border-blue-200'
                  ]"
                >
                  {{ log.action }}
                </span>
              </td>
              <td class="px-4 py-3.5 font-mono font-bold text-blue-700">
                {{ log.entityReference }}
              </td>
              <td class="px-4 py-3.5 text-slate-700 max-w-sm">
                {{ log.description }}
              </td>
              <td class="px-4 py-3.5 font-mono text-slate-500 text-[11px]">
                {{ log.ipAddress }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

