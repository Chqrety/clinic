<script setup lang="ts">
import { ref, computed } from 'vue'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import StatusBadge from '@/components/common/StatusBadge.vue'
import type { AccountType, NormalBalance, BusinessUnit } from '@/types/accounting'
import {
  FileSpreadsheet,
  Plus,
  Search,
  Filter,
  X,
  ChevronRight,
  Folder,
  FileText
} from 'lucide-vue-next'

const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

const searchQuery = ref('')
const selectedType = ref<string>('ALL')

// Add account modal
const showAddModal = ref(false)
const newCode = ref('')
const newName = ref('')
const newType = ref<AccountType>('EXPENSE')
const newNormalBalance = ref<NormalBalance>('DEBIT')
const newLevel = ref(3)
const newUnit = ref<BusinessUnit>('ALL')
const newDescription = ref('')

function handleSaveAccount() {
  if (!newCode.value.trim() || !newName.value.trim()) {
    accountingStore.showToast('Kode Akun dan Nama Akun wajib diisi!', 'error')
    return
  }

  ledgerStore.addAccount({
    code: newCode.value.trim(),
    name: newName.value.trim(),
    type: newType.value,
    normalBalance: newNormalBalance.value,
    level: newLevel.value,
    businessUnit: newUnit.value,
    description: newDescription.value.trim()
  })

  showAddModal.value = false
  newCode.value = ''
  newName.value = ''
  newDescription.value = ''
}

const typeOptions = [
  { value: 'ALL', label: 'Semua Tipe Akun' },
  { value: 'ASSET', label: '1 - ASET (Aktiva)' },
  { value: 'LIABILITY', label: '2 - KEWAJIBAN (Hutang)' },
  { value: 'EQUITY', label: '3 - EKUITAS (Modal)' },
  { value: 'REVENUE', label: '4 - PENDAPATAN' },
  { value: 'COGS', label: '5 - BEBAN POKOK (HPP)' },
  { value: 'EXPENSE', label: '6 - BEBAN OPERASIONAL' }
]

const filteredAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => {
    if (selectedType.value !== 'ALL' && a.type !== selectedType.value) {
      return false
    }
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchCode = a.code.toLowerCase().includes(q)
      const matchName = a.name.toLowerCase().includes(q)
      if (!matchCode && !matchName) return false
    }
    return true
  })
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Master Data Keuangan
          </span>
          <span class="text-xs text-slate-400">&bull; Standar Akuntansi Indonesia (PSAK)</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Bagan Akun Standar (Chart of Accounts)
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Struktur hierarki akun buku besar klinik untuk mengelompokkan aset, kewajiban, ekuitas, pendapatan, HPP dan beban.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="showAddModal = true"
          class="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah Akun Baru</span>
        </button>
      </div>
    </div>

    <!-- Toolbar Filters & Search -->
    <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto">
        <!-- Search -->
        <div class="relative w-full md:w-80">
          <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari kode akun atau nama..."
            class="w-full pl-9 pr-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <!-- Type Selector -->
        <div class="relative">
          <select
            v-model="selectedType"
            class="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium text-slate-800 focus:ring-2 focus:ring-blue-500 focus:outline-none cursor-pointer"
          >
            <option
              v-for="opt in typeOptions"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </option>
          </select>
        </div>
      </div>

      <div class="text-xs text-slate-500 font-medium">
        Total <strong>{{ filteredAccounts.length }}</strong> akun terdaftar
      </div>
    </div>

    <!-- Tree-like Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3 w-40">Kode Akun</th>
              <th class="px-4 py-3">Nama Akun</th>
              <th class="px-4 py-3">Klasifikasi Tipe</th>
              <th class="px-4 py-3 text-center">Saldo Normal</th>
              <th class="px-4 py-3 text-center">Unit Bisnis</th>
              <th class="px-4 py-3 text-right">Saldo Saat Ini (Rp)</th>
              <th class="px-4 py-3 text-center w-24">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="acc in filteredAccounts"
              :key="acc.id"
              :class="[
                'transition-colors',
                acc.isHeader
                  ? 'bg-slate-50/80 font-bold text-slate-900'
                  : 'hover:bg-slate-50/50 text-slate-700'
              ]"
            >
              <td class="px-4 py-3 font-mono font-bold whitespace-nowrap">
                <span :class="acc.isHeader ? 'text-slate-900 font-extrabold' : 'text-blue-700'">
                  {{ acc.code }}
                </span>
              </td>
              <td class="px-4 py-3">
                <div
                  :class="[
                    'flex items-center gap-2',
                    acc.level === 2 ? 'pl-4' : acc.level === 3 ? 'pl-8' : ''
                  ]"
                >
                  <Folder v-if="acc.isHeader" class="w-4 h-4 text-amber-500 shrink-0" />
                  <FileText v-else class="w-3.5 h-3.5 text-slate-400 shrink-0" />
                  <span :class="acc.isHeader ? 'font-bold text-slate-900' : 'font-medium'">
                    {{ acc.name }}
                  </span>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="text-[11px] font-semibold text-slate-600">
                  {{ acc.type }}
                </span>
              </td>
              <td class="px-4 py-3 text-center">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold',
                    acc.normalBalance === 'DEBIT'
                      ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                      : 'bg-rose-50 text-rose-700 border border-rose-200'
                  ]"
                >
                  {{ acc.normalBalance }}
                </span>
              </td>
              <td class="px-4 py-3 text-center">
                <StatusBadge :status="acc.businessUnit" size="sm" />
              </td>
              <td class="px-4 py-3 text-right font-mono font-semibold whitespace-nowrap">
                {{ acc.currentBalance ? formatIDR(acc.currentBalance) : '-' }}
              </td>
              <td class="px-4 py-3 text-center">
                <span class="inline-flex w-2 h-2 rounded-full bg-emerald-500" title="Aktif"></span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Account Modal -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="px-6 py-4 bg-slate-900 text-white flex items-center justify-between">
          <h3 class="text-sm font-bold text-white">Tambah Akun Bagan Baru (COA)</h3>
          <button @click="showAddModal = false" class="text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleSaveAccount" class="p-6 space-y-4 text-xs">
          <div>
            <label class="block font-bold text-slate-700 mb-1">Kode Akun</label>
            <input
              v-model="newCode"
              type="text"
              placeholder="Contoh: 6-1600"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-mono"
            />
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Nama Akun</label>
            <input
              v-model="newName"
              type="text"
              placeholder="Contoh: Beban Kebersihan & Sanitasi Medis"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Klasifikasi Tipe</label>
              <select
                v-model="newType"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              >
                <option value="ASSET">ASET</option>
                <option value="LIABILITY">KEWAJIBAN</option>
                <option value="EQUITY">EKUITAS</option>
                <option value="REVENUE">PENDAPATAN</option>
                <option value="COGS">HPP</option>
                <option value="EXPENSE">BEBAN</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Saldo Normal</label>
              <select
                v-model="newNormalBalance"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              >
                <option value="DEBIT">DEBIT</option>
                <option value="CREDIT">CREDIT</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Unit Bisnis</label>
              <select
                v-model="newUnit"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              >
                <option value="ALL">Semua Unit (Konsolidasi)</option>
                <option value="KIDINEST">Kidinest</option>
                <option value="SSC">SSC</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Level Hierarki</label>
              <select
                v-model.number="newLevel"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              >
                <option :value="2">Level 2 (Header Sub-Grup)</option>
                <option :value="3">Level 3 (Akun Transaksional)</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Deskripsi Tambahan</label>
            <input
              v-model="newDescription"
              type="text"
              placeholder="Catatan penggunaan akun..."
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
            <button
              @click="showAddModal = false"
              type="button"
              class="px-4 py-2 bg-white border border-slate-300 rounded-lg font-semibold text-slate-700 hover:bg-slate-100"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold shadow-xs"
            >
              Simpan Akun
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

