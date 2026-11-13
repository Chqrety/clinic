<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import StatusBadge from '@/components/common/StatusBadge.vue'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
import type { AccountType, NormalBalance, BusinessUnit } from '@/types/accounting'
import {
  FileSpreadsheet,
  Plus,
  Search,
  Filter,
  X,
  ChevronRight,
  Folder,
  FileText,
  BookOpen,
  Receipt,
  Layers
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

const searchQuery = ref('')
const selectedType = ref<string>('ALL')
const selectedLevel = ref<number | 'ALL'>('ALL')

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
  { value: 'ALL', label: 'Semua Kategori (1 - 7)' },
  { value: 'ASSET', label: '1 - ASET - BCS' },
  { value: 'LIABILITY', label: '2 - KEWAJIBAN - BCS' },
  { value: 'EQUITY', label: '3 - EKUITAS - BCS' },
  { value: 'REVENUE', label: '4 - PENDAPATAN & 6 - NON USAHA' },
  { value: 'COGS', label: '5 - HPP - BCS' },
  { value: 'EXPENSE', label: '52 - BEBAN OPERASIONAL & 7 - NON USAHA' }
]

const filteredAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => {
    if (selectedType.value !== 'ALL') {
      if (selectedType.value === 'REVENUE' && (a.type !== 'REVENUE')) return false
      else if (selectedType.value === 'EXPENSE' && (a.type !== 'EXPENSE')) return false
      else if (selectedType.value !== 'REVENUE' && selectedType.value !== 'EXPENSE' && a.type !== selectedType.value) return false
    }
    if (selectedLevel.value !== 'ALL' && a.level !== selectedLevel.value) {
      return false
    }
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchCode = a.code.toLowerCase().includes(q)
      const matchName = a.name.toLowerCase().includes(q)
      const matchTax = a.taxCategory ? a.taxCategory.toLowerCase().includes(q) : false
      if (!matchCode && !matchName && !matchTax) return false
    }
    return true
  })
})

function openLedgerForAccount(accountCode: string) {
  router.push({ path: '/ledger', query: { account: accountCode } })
}

function openTaxModule(taxCategory?: string) {
  router.push({ path: '/tax', query: { taxCategory } })
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Master Data Bagan Akun
          </span>
          <span class="text-xs text-slate-400">&bull; Standar Akuntansi Indonesia (PSAK & Tiered Indexing)</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Bagan Akun Standar (Chart of Accounts - BCS)
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Struktur berjenjang Index (Level 1), Sub-Index (Level 2), Sub-Index ke-2 (Level 3), dan Leaf Detail (Level 4) dengan integrasi otomatis ke modul Perpajakan.
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

    <!-- Toolbar Filters & Search & Level Selectors -->
    <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto">
        <!-- Search -->
        <div class="relative w-full md:w-72">
          <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari kode akun, nama, atau jenis pajak..."
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

        <!-- Tiered Index Level Selector Buttons matching user request -->
        <div class="flex items-center bg-slate-100 p-0.5 rounded-lg border border-slate-200">
          <button
            @click="selectedLevel = 'ALL'"
            :class="[
              'px-2.5 py-1 rounded text-xs font-bold transition-all',
              selectedLevel === 'ALL' ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-600 hover:text-slate-900'
            ]"
          >
            Semua
          </button>
          <button
            @click="selectedLevel = 1"
            :class="[
              'px-2.5 py-1 rounded text-xs font-bold transition-all',
              selectedLevel === 1 ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-600 hover:text-slate-900'
            ]"
            title="Index Utama (Aset, Kewajiban, Ekuitas, Pendapatan, Beban)"
          >
            Lvl 1 (Index)
          </button>
          <button
            @click="selectedLevel = 2"
            :class="[
              'px-2.5 py-1 rounded text-xs font-bold transition-all',
              selectedLevel === 2 ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-600 hover:text-slate-900'
            ]"
            title="Sub-Index Kelompok (Kas & Bank, Piutang, Hutang Pajak, dll.)"
          >
            Lvl 2 (Sub-Index)
          </button>
          <button
            @click="selectedLevel = 3"
            :class="[
              'px-2.5 py-1 rounded text-xs font-bold transition-all',
              selectedLevel === 3 ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-600 hover:text-slate-900'
            ]"
            title="Sub-Index dari Sub-Index (Kasir, Bank BCA, PPh 21, PPh 23, dll.)"
          >
            Lvl 3
          </button>
          <button
            @click="selectedLevel = 4"
            :class="[
              'px-2.5 py-1 rounded text-xs font-bold transition-all',
              selectedLevel === 4 ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-600 hover:text-slate-900'
            ]"
            title="Rincian Detail Rekening Buku Besar"
          >
            Lvl 4 (Detail)
          </button>
        </div>
      </div>

      <div class="text-xs text-slate-500 font-medium">
        Menampilkan <strong>{{ filteredAccounts.length }}</strong> akun
      </div>
    </div>

    <!-- Tree-like Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3 w-36">Kode Akun</th>
              <th class="px-4 py-3">Nama Akun (Hierarki COA)</th>
              <th class="px-4 py-3 text-center w-20">Level</th>
              <th class="px-4 py-3">Klasifikasi Tipe</th>
              <th class="px-4 py-3 text-center">Saldo Normal</th>
              <th class="px-4 py-3 text-center">Unit Bisnis</th>
              <th class="px-4 py-3 text-right">Saldo Saat Ini (Rp)</th>
              <th class="px-4 py-3 text-center w-28">Integrasi Pajak</th>
              <th class="px-4 py-3 text-center w-24">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="acc in filteredAccounts"
              :key="acc.id"
              :class="[
                'transition-colors',
                acc.level === 1
                  ? 'bg-slate-900 text-white font-black hover:bg-slate-800'
                  : acc.level === 2
                  ? 'bg-slate-100/90 font-bold text-slate-900 hover:bg-slate-200/70'
                  : acc.level === 3
                  ? 'bg-slate-50/50 font-semibold text-slate-800 hover:bg-slate-100/80'
                  : 'hover:bg-blue-50/40 text-slate-700'
              ]"
            >
              <!-- Kode Akun -->
              <td class="px-4 py-2.5 font-mono font-bold whitespace-nowrap">
                <span :class="acc.level === 1 ? 'text-amber-400 font-extrabold text-sm' : acc.level === 2 ? 'text-blue-700 font-bold' : 'text-slate-700'">
                  {{ acc.code }}
                </span>
              </td>

              <!-- Nama Akun dengan Indentasi Hierarki -->
              <td class="px-4 py-2.5">
                <div
                  :class="[
                    'flex items-center gap-2',
                    acc.level === 2 ? 'pl-4' : acc.level === 3 ? 'pl-8' : acc.level === 4 ? 'pl-12' : ''
                  ]"
                >
                  <Folder v-if="acc.isHeader" :class="['w-4 h-4 shrink-0', acc.level === 1 ? 'text-amber-400' : 'text-amber-500']" />
                  <FileText v-else class="w-3.5 h-3.5 text-slate-400 shrink-0" />
                  <span :class="[acc.level === 1 ? 'font-extrabold text-white text-xs tracking-wide' : acc.level === 2 ? 'font-bold text-slate-900' : 'font-medium']">
                    {{ acc.name }}
                  </span>
                </div>
              </td>

              <!-- Level Badge -->
              <td class="px-4 py-2.5 text-center font-mono">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold',
                    acc.level === 1 ? 'bg-amber-500/20 text-amber-300 border border-amber-400/40' :
                    acc.level === 2 ? 'bg-blue-100 text-blue-800 border border-blue-200' :
                    acc.level === 3 ? 'bg-purple-100 text-purple-800 border border-purple-200' :
                    'bg-slate-100 text-slate-600'
                  ]"
                >
                  L{{ acc.level }}
                </span>
              </td>

              <!-- Tipe Klasifikasi -->
              <td class="px-4 py-2.5">
                <span :class="['text-[11px] font-semibold', acc.level === 1 ? 'text-slate-300' : 'text-slate-600']">
                  {{ acc.type }}
                </span>
              </td>

              <!-- Saldo Normal -->
              <td class="px-4 py-2.5 text-center">
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

              <!-- Unit Bisnis -->
              <td class="px-4 py-2.5 text-center">
                <StatusBadge :status="acc.businessUnit" size="sm" />
              </td>

              <!-- Saldo Saat Ini -->
              <td :class="['px-4 py-2.5 text-right font-mono font-semibold whitespace-nowrap', acc.level === 1 ? 'text-white font-bold' : 'text-slate-800']">
                {{ acc.currentBalance !== undefined ? formatIDR(acc.currentBalance) : '-' }}
              </td>

              <!-- Integrasi Pajak Badge -->
              <td class="px-4 py-2.5 text-center">
                <button
                  v-if="acc.taxCategory"
                  @click="openTaxModule(acc.taxCategory)"
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-800 hover:bg-emerald-200 text-[10px] font-bold transition-colors cursor-pointer"
                  title="Klik untuk membuka kalkulasi pajak otomatis di Modul Pajak"
                >
                  <Receipt class="w-2.5 h-2.5" />
                  <span>{{ acc.taxCategory }}</span>
                </button>
                <span v-else class="text-slate-300 text-[11px]">-</span>
              </td>

              <!-- Aksi (Buku Besar / Ledger) -->
              <td class="px-4 py-2.5 text-center">
                <button
                  @click="openLedgerForAccount(acc.code)"
                  class="inline-flex items-center gap-1 px-2.5 py-1 bg-slate-100 hover:bg-blue-50 text-slate-600 hover:text-blue-600 rounded-md text-[11px] font-semibold transition-colors border border-slate-200"
                  title="Lihat Buku Besar (General Ledger)"
                >
                  <BookOpen class="w-3 h-3" />
                  <span>Buku Besar</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Approval Signature Section matching Image 1 -->
    <ReportSignatures />

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
              placeholder="Contoh: 2501.003"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-mono font-bold focus:ring-2 focus:ring-blue-500 focus:outline-none"
              required
            />
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Nama Akun</label>
            <input
              v-model="newName"
              type="text"
              placeholder="Contoh: Hutang PPh 21 Tenaga Medis Konsultan"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
              required
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Level Hierarki</label>
              <select
                v-model.number="newLevel"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium focus:ring-2 focus:ring-blue-500 focus:outline-none"
              >
                <option :value="1">Level 1 (Index Utama)</option>
                <option :value="2">Level 2 (Sub-Index)</option>
                <option :value="3">Level 3 (Sub-Index ke-2)</option>
                <option :value="4">Level 4 (Detail Leaf)</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Unit Bisnis</label>
              <select
                v-model="newUnit"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium focus:ring-2 focus:ring-blue-500 focus:outline-none"
              >
                <option value="ALL">Semua Unit (Konsolidasi)</option>
                <option value="KIDINEST">Kidinest</option>
                <option value="SSC">Semarang Skin Center</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Tipe Klasifikasi</label>
              <select
                v-model="newType"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium focus:ring-2 focus:ring-blue-500 focus:outline-none"
              >
                <option value="ASSET">1 - ASET</option>
                <option value="LIABILITY">2 - KEWAJIBAN</option>
                <option value="EQUITY">3 - EKUITAS</option>
                <option value="REVENUE">4 / 6 - PENDAPATAN</option>
                <option value="COGS">5 - HPP</option>
                <option value="EXPENSE">52 / 7 - BEBAN</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Saldo Normal</label>
              <select
                v-model="newNormalBalance"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium focus:ring-2 focus:ring-blue-500 focus:outline-none"
              >
                <option value="DEBIT">DEBIT (Dr)</option>
                <option value="CREDIT">CREDIT (Cr)</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Keterangan / Catatan</label>
            <textarea
              v-model="newDescription"
              rows="2"
              placeholder="Deskripsi peruntukan akun buku besar..."
              class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
            ></textarea>
          </div>

          <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
            <button
              type="button"
              @click="showAddModal = false"
              class="px-4 py-2 border border-slate-200 hover:bg-slate-100 text-slate-700 rounded-lg font-bold"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold shadow-xs"
            >
              Simpan Akun
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

