<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
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
  ChevronDown,
  ChevronsUpDown,
  Folder,
  FolderOpen,
  FileText,
  BookOpen,
  Receipt,
  Layers,
  Database,
  BookmarkCheck,
  Edit3,
  Power
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

onMounted(() => {
  ledgerStore.fetchAccounts()
})

const searchQuery = ref('')
const selectedType = ref<string>('ALL')
const selectedLevel = ref<number | 'ALL'>('ALL')
const selectedStatus = ref<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL')

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

// Edit account modal & handlers
const showEditModal = ref(false)
const editingAccount = ref<any>(null)
const editName = ref('')
const editDescription = ref('')

function openEditModal(acc: any) {
  editingAccount.value = acc
  editName.value = acc.name
  editDescription.value = acc.description || ''
  showEditModal.value = true
}

async function handleUpdateAccount() {
  if (!editName.value.trim() || !editingAccount.value) {
    accountingStore.showToast('Nama Akun wajib diisi!', 'error')
    return
  }

  await ledgerStore.updateAccount(editingAccount.value.id, {
    name: editName.value.trim(),
    description: editDescription.value.trim()
  })

  showEditModal.value = false
  editingAccount.value = null
}

async function handleToggleActive(acc: any) {
  const newStatus = acc.isActive === false ? true : false
  const promptText = newStatus
    ? `Aktifkan kembali akun ${acc.code} - ${acc.name}?`
    : `Nonaktifkan akun ${acc.code} - ${acc.name}? Akun tidak akan muncul di pilihan transaksi baru.`

  if (confirm(promptText)) {
    await ledgerStore.updateAccount(acc.id, {
      isActive: newStatus
    })
  }
}

const typeOptions = [
  { value: 'ALL', label: 'Semua Kategori (1 - 7)' },
  { value: 'ASSET', label: '1 - ASET' },
  { value: 'LIABILITY', label: '2 - KEWAJIBAN' },
  { value: 'EQUITY', label: '3 - EKUITAS' },
  { value: 'REVENUE', label: '4 - PENDAPATAN & 6 - NON USAHA' },
  { value: 'COGS', label: '5 - HPP' },
  { value: 'EXPENSE', label: '52 - BEBAN OPERASIONAL & 7 - NON USAHA' }
]

// Expand & Collapse Tree Nodes
const collapsedCodes = ref<Set<string>>(new Set())

function toggleCollapse(code: string) {
  const next = new Set(collapsedCodes.value)
  if (next.has(code)) {
    next.delete(code)
  } else {
    next.add(code)
  }
  collapsedCodes.value = next
}

function expandAll() {
  collapsedCodes.value = new Set()
}

function collapseAll() {
  const headers = ledgerStore.accounts
    .filter(a => a.isHeader || a.level <= 2)
    .map(a => a.code)
  collapsedCodes.value = new Set(headers)
}

function isAncestorCollapsed(acc: any): boolean {
  if (searchQuery.value.trim() || selectedLevel.value !== 'ALL') return false
  let curr = acc.parentCode
  while (curr) {
    if (collapsedCodes.value.has(curr)) return true
    const parentAcc = ledgerStore.accounts.find(a => a.code === curr)
    curr = parentAcc ? parentAcc.parentCode : undefined
  }
  return false
}

const filteredAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => {
    if (accountingStore.businessUnit !== 'ALL' && a.businessUnit && a.businessUnit !== 'ALL' && a.businessUnit !== accountingStore.businessUnit) {
      return false
    }
    if (selectedType.value !== 'ALL') {
      if (selectedType.value === 'REVENUE' && (a.type !== 'REVENUE')) return false
      else if (selectedType.value === 'EXPENSE' && (a.type !== 'EXPENSE')) return false
      else if (selectedType.value !== 'REVENUE' && selectedType.value !== 'EXPENSE' && a.type !== selectedType.value) return false
    }
    if (selectedLevel.value !== 'ALL' && a.level !== selectedLevel.value) {
      return false
    }
    if (selectedStatus.value === 'ACTIVE' && a.isActive === false) {
      return false
    }
    if (selectedStatus.value === 'INACTIVE' && a.isActive !== false) {
      return false
    }
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchCode = a.code.toLowerCase().includes(q)
      const matchName = a.name.toLowerCase().includes(q)
      const matchTax = a.taxCategory ? a.taxCategory.toLowerCase().includes(q) : false
      if (!matchCode && !matchName && !matchTax) return false
    }
    if (isAncestorCollapsed(a)) {
      return false
    }
    return true
  })
})

function openLedgerForAccount(accountCode: string) {
  router.push({ path: '/accounting/ledger', query: { account: accountCode } })
}

function openTaxModule(taxCategory?: string) {
  router.push({ path: '/accounting/tax/spt-masa', query: { taxCategory } })
}
</script>

<template>
  <div class="space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Minimalist Clinical Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Bagan Akun Standar
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">PSAK & Hierarki Akuntansi Klinik</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Chart of Accounts (COA)
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="router.push('/accounting/opening-balance')"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <BookmarkCheck class="w-3.5 h-3.5 text-zinc-600" />
          <span>Atur Saldo Awal</span>
        </button>

        <button
          @click="showAddModal = true"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Tambah Akun Baru</span>
        </button>
      </div>
    </header>

    <!-- Toolbar Filters & Search & Level Selectors -->
    <div class="bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto">
        <!-- Search -->
        <div class="relative w-full md:w-72">
          <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari kode akun, nama, atau jenis pajak..."
            class="w-full pl-9 pr-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs focus:ring-1 focus:ring-zinc-950 focus:outline-none"
          />
        </div>

        <!-- Type Selector -->
        <div class="relative">
          <select
            v-model="selectedType"
            class="px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium text-zinc-800 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
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

        <!-- Tiered Index Level Selector Buttons -->
        <div class="flex items-center bg-zinc-100 p-0.5 rounded-xl border border-zinc-200/90">
          <button
            @click="selectedLevel = 'ALL'"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 'ALL' ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            Semua
          </button>
          <button
            @click="selectedLevel = 1"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 1 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
            title="Index Utama (Aset, Kewajiban, Ekuitas, Pendapatan, Beban)"
          >
            Lvl 1 (Index)
          </button>
          <button
            @click="selectedLevel = 2"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 2 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
            title="Sub-Index Kelompok (Kas & Bank, Piutang, Hutang Pajak, dll.)"
          >
            Lvl 2 (Sub-Index)
          </button>
          <button
            @click="selectedLevel = 3"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 3 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
            title="Sub-Index ke-2 (Kasir, Bank BCA, PPh 21, PPh 23, dll.)"
          >
            Lvl 3
          </button>
          <button
            @click="selectedLevel = 4"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 4 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
            title="Rincian Detail Rekening Buku Besar"
          >
            Lvl 4 (Detail)
          </button>
        </div>

        <!-- Status Filter: Semua / Aktif / Nonaktif -->
        <div class="flex items-center bg-zinc-100 p-0.5 rounded-xl border border-zinc-200/90 text-xs">
          <button
            @click="selectedStatus = 'ALL'"
            :class="[
              'px-2.5 py-1 rounded-lg font-bold transition-all',
              selectedStatus === 'ALL' ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            Semua
          </button>
          <button
            @click="selectedStatus = 'ACTIVE'"
            :class="[
              'px-2.5 py-1 rounded-lg font-bold transition-all',
              selectedStatus === 'ACTIVE' ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            Aktif
          </button>
          <button
            @click="selectedStatus = 'INACTIVE'"
            :class="[
              'px-2.5 py-1 rounded-lg font-bold transition-all',
              selectedStatus === 'INACTIVE' ? 'bg-rose-600 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            Nonaktif
          </button>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <div class="flex items-center gap-1.5">
          <button
            type="button"
            @click="expandAll"
            class="px-2.5 py-1 text-xs font-semibold text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100 rounded-lg border border-zinc-200 transition-colors cursor-pointer"
            title="Buka seluruh sub-akun"
          >
            Buka Semua
          </button>
          <button
            type="button"
            @click="collapseAll"
            class="px-2.5 py-1 text-xs font-semibold text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100 rounded-lg border border-zinc-200 transition-colors cursor-pointer"
            title="Ciutkan seluruh sub-akun"
          >
            Ciutkan Semua
          </button>
        </div>
        <div class="text-xs text-zinc-500 font-medium">
          Menampilkan <strong>{{ filteredAccounts.length }}</strong> akun
        </div>
      </div>
    </div>

    <!-- Tree-like Table with Clear Visual Hierarchy -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200/90 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-4 py-3 w-36">Kode Akun</th>
              <th class="px-4 py-3">Nama Akun (Hierarki COA)</th>
              <th class="px-4 py-3 text-center w-20">Level</th>
              <th class="px-4 py-3">Klasifikasi Tipe</th>
              <th class="px-4 py-3 text-center">Saldo Normal</th>
              <th class="px-4 py-3 text-center">Unit Bisnis</th>
              <th class="px-4 py-3 text-right">Saldo Saat Ini (Rp)</th>
              <th class="px-4 py-3 text-center w-28">Integrasi Pajak</th>
              <th class="px-4 py-3 text-center w-32">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr
              v-for="acc in filteredAccounts"
              :key="acc.id"
              :class="[
                'transition-colors',
                acc.isActive === false ? 'opacity-50 bg-zinc-50/50' : '',
                acc.level === 1
                  ? 'bg-zinc-100/90 text-zinc-950 font-extrabold border-t-2 border-b border-zinc-200/90 hover:bg-zinc-200/60'
                  : acc.level === 2
                  ? 'bg-zinc-50/80 font-bold text-zinc-900 border-b border-zinc-100 hover:bg-zinc-100/70'
                  : acc.level === 3
                  ? 'bg-white font-medium text-zinc-800 hover:bg-zinc-50/70'
                  : 'bg-white/80 text-zinc-600 hover:bg-zinc-50/70'
              ]"
            >
              <!-- Kode Akun -->
              <td class="px-4 py-2.5 font-mono font-bold whitespace-nowrap">
                <span :class="acc.level === 1 ? 'text-zinc-950 font-black text-xs tracking-tight' : acc.level === 2 ? 'text-zinc-900 font-bold text-xs' : 'text-zinc-600 text-xs'">
                  {{ acc.code }}
                </span>
              </td>

              <!-- Nama Akun dengan Indentasi Hierarki & Collapsible Toggle -->
              <td class="px-4 py-2.5">
                <div
                  :class="[
                    'flex items-center gap-2',
                    acc.level === 2 ? 'pl-4 border-l-2 border-zinc-200/80 ml-1' :
                    acc.level === 3 ? 'pl-8 border-l-2 border-zinc-200/50 ml-1' :
                    acc.level === 4 ? 'pl-12 border-l-2 border-zinc-200/30 ml-1' : ''
                  ]"
                >
                  <!-- Expand / Collapse chevron for headers -->
                  <button
                    v-if="acc.isHeader"
                    type="button"
                    @click.stop="toggleCollapse(acc.code)"
                    class="p-0.5 -ml-1 text-zinc-500 hover:text-zinc-900 hover:bg-zinc-200/80 rounded transition-colors cursor-pointer"
                    :title="collapsedCodes.has(acc.code) ? 'Buka Sub-Akun' : 'Ciutkan Sub-Akun'"
                  >
                    <ChevronRight v-if="collapsedCodes.has(acc.code)" class="w-3.5 h-3.5" />
                    <ChevronDown v-else class="w-3.5 h-3.5" />
                  </button>

                  <Folder v-if="acc.isHeader" :class="['w-4 h-4 shrink-0', acc.level === 1 ? 'text-amber-600' : 'text-amber-500']" />
                  <FileText v-else class="w-3.5 h-3.5 text-zinc-400 shrink-0" />

                  <span
                    :class="[
                      acc.level === 1 ? 'font-black text-zinc-950 text-xs tracking-wide uppercase' : acc.level === 2 ? 'font-bold text-zinc-900' : acc.level === 3 ? 'font-semibold text-zinc-800' : 'font-normal text-zinc-600',
                      acc.isActive === false ? 'line-through text-zinc-400' : ''
                    ]"
                  >
                    {{ acc.name }}
                  </span>
                  <span
                    v-if="acc.isActive === false"
                    class="ml-2 px-1.5 py-0.5 rounded text-[9px] font-bold bg-rose-50 text-rose-700 border border-rose-200 uppercase tracking-wider"
                  >
                    Nonaktif
                  </span>
                </div>
              </td>

              <!-- Level Badge -->
              <td class="px-4 py-2.5 text-center font-mono">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold',
                    acc.level === 1 ? 'bg-zinc-200 text-zinc-900 border border-zinc-300' :
                    acc.level === 2 ? 'bg-zinc-100 text-zinc-800 border border-zinc-200' :
                    acc.level === 3 ? 'bg-zinc-50 text-zinc-700 border border-zinc-200/70' :
                    'bg-transparent text-zinc-400 border border-zinc-200/50'
                  ]"
                >
                  L{{ acc.level }}
                </span>
              </td>

              <!-- Tipe Klasifikasi -->
              <td class="px-4 py-2.5">
                <span :class="['text-[11px] font-semibold', acc.level === 1 ? 'text-zinc-900 font-bold' : 'text-zinc-600']">
                  {{ acc.type }}
                </span>
              </td>

              <!-- Saldo Normal -->
              <td class="px-4 py-2.5 text-center">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold',
                    acc.normalBalance === 'DEBIT'
                      ? 'bg-emerald-50 text-emerald-800 border border-emerald-200'
                      : 'bg-rose-50 text-rose-800 border border-rose-200'
                  ]"
                >
                  {{ acc.normalBalance }}
                </span>
              </td>

              <!-- Unit Bisnis -->
              <td class="px-4 py-2.5 text-center">
                <StatusBadge :status="acc.businessUnit" size="sm" />
              </td>

              <!-- Saldo Saat Ini (Strict Tabular Nums) -->
              <td :class="['px-4 py-2.5 text-right tabular-nums font-bold whitespace-nowrap', acc.level === 1 ? 'text-zinc-950 font-black' : 'text-zinc-900']">
                {{ acc.currentBalance !== undefined ? formatIDR(acc.currentBalance) : '-' }}
              </td>

              <!-- Integrasi Pajak Badge -->
              <td class="px-4 py-2.5 text-center">
                <button
                  v-if="acc.taxCategory"
                  @click="openTaxModule(acc.taxCategory)"
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-800 hover:bg-emerald-100 text-[10px] font-bold transition-colors cursor-pointer border border-emerald-200"
                  title="Klik untuk membuka kalkulasi pajak otomatis di Modul Pajak"
                >
                  <Receipt class="w-2.5 h-2.5" />
                  <span>{{ acc.taxCategory }}</span>
                </button>
                <span v-else class="text-zinc-300 text-[11px]">-</span>
              </td>

              <!-- Aksi (Buku Besar, Edit, Nonaktifkan) -->
              <td class="px-4 py-2.5 text-center">
                <div class="flex items-center justify-center gap-1.5">
                  <button
                    @click="openLedgerForAccount(acc.code)"
                    class="p-1.5 bg-white hover:bg-zinc-100 text-zinc-600 hover:text-zinc-950 rounded-lg transition-colors border border-zinc-200 shadow-2xs cursor-pointer"
                    title="Buku Besar (General Ledger)"
                  >
                    <BookOpen class="w-3.5 h-3.5" />
                  </button>

                  <button
                    v-if="!acc.isHeader"
                    @click="openEditModal(acc)"
                    class="p-1.5 bg-white hover:bg-zinc-100 text-zinc-600 hover:text-zinc-950 rounded-lg transition-colors border border-zinc-200 shadow-2xs cursor-pointer"
                    title="Edit Nama / Keterangan Akun"
                  >
                    <Edit3 class="w-3.5 h-3.5" />
                  </button>

                  <button
                    v-if="!acc.isHeader"
                    @click="handleToggleActive(acc)"
                    :class="[
                      'p-1.5 rounded-lg transition-colors border shadow-2xs cursor-pointer',
                      acc.isActive !== false
                        ? 'bg-white hover:bg-rose-50 text-zinc-500 hover:text-rose-600 border-zinc-200'
                        : 'bg-emerald-600 hover:bg-emerald-700 text-white border-emerald-600'
                    ]"
                    :title="acc.isActive !== false ? 'Nonaktifkan Akun' : 'Aktifkan Kembali Akun'"
                  >
                    <Power class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Approval Signature Section -->
    <ReportSignatures />

    <!-- Add Account Modal (Official Clinical Modal Specification) -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-black/40 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="p-5 border-b border-zinc-100 flex items-center justify-between">
          <h3 class="text-sm font-bold text-zinc-950">Tambah Akun Bagan Baru (COA)</h3>
          <button @click="showAddModal = false" class="p-1.5 text-zinc-400 hover:text-zinc-700 rounded-lg hover:bg-zinc-100 transition-colors cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleSaveAccount" class="p-6 space-y-4 text-xs">
          <div>
            <label class="block font-bold text-zinc-800 mb-1">Kode Akun</label>
            <input
              v-model="newCode"
              type="text"
              placeholder="Contoh: 2-1130 atau 2501.003"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-mono font-bold focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              required
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-800 mb-1">Nama Akun</label>
            <input
              v-model="newName"
              type="text"
              placeholder="Contoh: Hutang PPh 21 Tenaga Medis Konsultan"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              required
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-zinc-800 mb-1">Level Hierarki</label>
              <select
                v-model.number="newLevel"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              >
                <option :value="1">Level 1 (Index Utama)</option>
                <option :value="2">Level 2 (Sub-Index)</option>
                <option :value="3">Level 3 (Sub-Index ke-2)</option>
                <option :value="4">Level 4 (Detail Leaf)</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-zinc-800 mb-1">Unit Bisnis</label>
              <select
                v-model="newUnit"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              >
                <option value="ALL">Semua Unit (Konsolidasi)</option>
                <option value="KIDINEST">Kidinest</option>
                <option value="SSC">Semarang Skin Center</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-zinc-800 mb-1">Tipe Klasifikasi</label>
              <select
                v-model="newType"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none"
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
              <label class="block font-bold text-zinc-800 mb-1">Saldo Normal</label>
              <select
                v-model="newNormalBalance"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              >
                <option value="DEBIT">DEBIT (Dr)</option>
                <option value="CREDIT">CREDIT (Cr)</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block font-bold text-zinc-800 mb-1">Keterangan / Catatan</label>
            <textarea
              v-model="newDescription"
              rows="2"
              placeholder="Deskripsi peruntukan akun buku besar..."
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            ></textarea>
          </div>

          <div class="p-4 border-t border-zinc-100 bg-zinc-50/70 -mx-6 -mb-6 mt-4 flex items-center justify-end gap-2.5">
            <button
              type="button"
              @click="showAddModal = false"
              class="px-4 py-2 border border-zinc-200 hover:bg-zinc-100 text-zinc-700 rounded-xl font-bold cursor-pointer transition-colors shadow-2xs"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold shadow-2xs cursor-pointer active:scale-95 transition-all"
            >
              Simpan Akun
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Account Modal (Official Clinical Modal Specification) -->
    <div
      v-if="showEditModal && editingAccount"
      class="fixed inset-0 z-50 overflow-y-auto bg-black/40 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-md w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="p-5 border-b border-zinc-100 flex items-center justify-between">
          <div>
            <h3 class="text-sm font-bold text-zinc-950">Edit Akun COA</h3>
            <p class="text-[11px] font-mono text-zinc-500 mt-0.5">{{ editingAccount.code }} &bull; Level {{ editingAccount.level }}</p>
          </div>
          <button
            @click="showEditModal = false"
            class="p-1.5 text-zinc-400 hover:text-zinc-700 rounded-lg hover:bg-zinc-100 transition-colors cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleUpdateAccount" class="p-6 space-y-4 text-xs">
          <div>
            <label class="block font-bold text-zinc-800 mb-1">Kode Akun</label>
            <input
              :value="editingAccount.code"
              disabled
              class="w-full px-3 py-2 bg-zinc-100 border border-zinc-200 rounded-xl text-xs font-mono font-bold text-zinc-500 cursor-not-allowed"
            />
            <p class="text-[10px] text-zinc-400 mt-1">*Kode akun dikunci untuk menjaga integritas riwayat jurnal</p>
          </div>

          <div>
            <label class="block font-bold text-zinc-800 mb-1">Nama Akun <span class="text-rose-500">*</span></label>
            <input
              v-model="editName"
              required
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-bold focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-800 mb-1">Keterangan / Catatan</label>
            <textarea
              v-model="editDescription"
              rows="3"
              placeholder="Peruntukan akun buku besar..."
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            ></textarea>
          </div>

          <div class="p-4 border-t border-zinc-100 bg-zinc-50/70 -mx-6 -mb-6 mt-4 flex items-center justify-end gap-2.5">
            <button
              type="button"
              @click="showEditModal = false"
              class="px-4 py-2 border border-zinc-200 hover:bg-zinc-100 text-zinc-700 rounded-xl font-bold cursor-pointer transition-colors shadow-2xs"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold shadow-2xs cursor-pointer active:scale-95 transition-all"
            >
              Simpan Perubahan
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

