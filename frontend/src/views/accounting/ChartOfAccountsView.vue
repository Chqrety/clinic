<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { useLocaleStore } from '@/stores/locale'
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
const localeStore = useLocaleStore()

onMounted(() => {
  ledgerStore.fetchAccounts()
})

const searchQuery = ref('')
const selectedType = ref<string>('ALL')
const selectedSubCategory = ref<string>('ALL')
const selectedLevel = ref<number | 'ALL'>('ALL')
const selectedStatus = ref<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL')

// Kategori Utama (Level 1) diambil dinamis 100% dari struktur COA
const categoryOptions = computed(() => {
  const lvl1 = ledgerStore.accounts
    .filter(a => a.level === 1)
    .sort((a, b) => a.code.localeCompare(b.code))

  const list = lvl1.map(a => {
    const translatedName = localeStore.translateAccountName(a.name)
    return {
      value: a.type,
      code: a.code,
      name: translatedName,
      label: `${a.code} - ${translatedName}`
    }
  })

  return [{ value: 'ALL', code: 'ALL', name: localeStore.t('coa.allCategories'), label: localeStore.t('coa.allCategories') }, ...list]
})

// Sub-Kategori / Kelompok Akun (Level 2 Header) diambil dinamis dari COA sesuai Kategori yang dipilih
const subCategoryOptions = computed(() => {
  let lvl2 = ledgerStore.accounts
    .filter(a => a.level === 2 && a.isHeader)
    .sort((a, b) => a.code.localeCompare(b.code))

  if (selectedType.value !== 'ALL') {
    lvl2 = lvl2.filter(a => a.type === selectedType.value)
  }

  return [
    { value: 'ALL', code: 'ALL', name: localeStore.t('coa.allSubCategories'), label: localeStore.t('coa.allSubCategories') },
    ...lvl2.map(a => {
      const translatedName = localeStore.translateAccountName(a.name)
      return {
        value: a.code,
        code: a.code,
        name: translatedName,
        label: `${a.code} - ${translatedName}`
      }
    })
  ]
})

// Add Account / Sub-Category Modal States
const showAddModal = ref(false)
const isHeaderMode = ref(false) // false = Akun Transaksi / Detail, true = Sub-Kategori / Header Kelompok Baru
const selectedAddCategoryCode = ref('1-0000')
const selectedAddParentCode = ref('')
const newCode = ref('')
const newName = ref('')
const newType = ref<AccountType>('ASSET')
const newNormalBalance = ref<NormalBalance>('DEBIT')
const newLevel = ref(3)
const newUnit = ref<BusinessUnit>('ALL')
const newTaxCategory = ref<string>('NON_TAX')
const newDescription = ref('')

// Daftar Kategori Level 1 untuk Modal
const addCategoryList = computed(() => {
  return ledgerStore.accounts
    .filter(a => a.level === 1)
    .sort((a, b) => a.code.localeCompare(b.code))
})

// Daftar Sub-Kategori (Level 2/3 Header) sesuai kategori yang dipilih untuk modal
const addSubCategoryList = computed(() => {
  return ledgerStore.accounts
    .filter(a => (a.level === 2 || a.level === 3) && a.isHeader && (a.type === newType.value || a.parentCode === selectedAddCategoryCode.value))
    .sort((a, b) => a.code.localeCompare(b.code))
})

function suggestNextCode(parentCode: string): string {
  const children = ledgerStore.accounts.filter(a => a.parentCode === parentCode)
  if (children.length > 0) {
    const codes = children.map(c => c.code)
    const lastCode = codes[codes.length - 1]
    const match = lastCode.match(/^(.*?)(\d+)$/)
    if (match) {
      const prefix = match[1]
      const num = parseInt(match[2], 10) + 1
      const paddedNum = String(num).padStart(match[2].length, '0')
      return `${prefix}${paddedNum}`
    }
  }
  const match = parentCode.match(/^(.*?)(\d+)$/)
  if (match) {
    const prefix = match[1]
    const num = parseInt(match[2], 10) + 10
    return `${prefix}${num}`
  }
  return `${parentCode}.01`
}

function onCategoryChanged() {
  const cat = ledgerStore.accounts.find(a => a.code === selectedAddCategoryCode.value)
  if (cat) {
    newType.value = cat.type
    newNormalBalance.value = cat.normalBalance
  }
  const availableSub = addSubCategoryList.value
  if (!isHeaderMode.value && availableSub.length > 0) {
    selectedAddParentCode.value = availableSub[0].code
    onParentChanged()
  } else {
    selectedAddParentCode.value = selectedAddCategoryCode.value
    newLevel.value = 2
    newCode.value = suggestNextCode(selectedAddCategoryCode.value)
  }
}

function onParentChanged() {
  const parent = ledgerStore.accounts.find(a => a.code === selectedAddParentCode.value)
  if (parent) {
    newLevel.value = isHeaderMode.value ? parent.level : Math.min(parent.level + 1, 4)
    newCode.value = suggestNextCode(parent.code)
  }
}

function onHeaderModeChanged() {
  if (isHeaderMode.value) {
    selectedAddParentCode.value = selectedAddCategoryCode.value
    newLevel.value = 2
    newCode.value = suggestNextCode(selectedAddCategoryCode.value)
  } else {
    const availableSub = addSubCategoryList.value
    if (availableSub.length > 0) {
      selectedAddParentCode.value = availableSub[0].code
      onParentChanged()
    }
  }
}

function openAddModal(parentAcc?: any) {
  if (parentAcc) {
    // Quick Add di bawah kelompok ini
    let topLvl: any = parentAcc
    while (topLvl.level > 1 && topLvl.parentCode) {
      const p = ledgerStore.accounts.find(a => a.code === topLvl.parentCode)
      if (p) topLvl = p
      else break
    }
    selectedAddCategoryCode.value = topLvl.code
    newType.value = parentAcc.type
    newNormalBalance.value = parentAcc.normalBalance
    selectedAddParentCode.value = parentAcc.code
    isHeaderMode.value = false
    newLevel.value = Math.min(parentAcc.level + 1, 4)
    newCode.value = suggestNextCode(parentAcc.code)
  } else {
    selectedAddCategoryCode.value = addCategoryList.value[0]?.code || '1-0000'
    isHeaderMode.value = false
    onCategoryChanged()
  }
  newName.value = ''
  newDescription.value = ''
  newTaxCategory.value = 'NON_TAX'
  newUnit.value = accountingStore.businessUnit !== 'ALL' ? (accountingStore.businessUnit as BusinessUnit) : 'ALL'
  showAddModal.value = true
}

async function handleSaveAccount() {
  if (!newCode.value.trim() || !newName.value.trim()) {
    accountingStore.showToast('Kode Akun dan Nama Akun wajib diisi!', 'error')
    return
  }

  const codeTrimmed = newCode.value.trim()
  const exists = ledgerStore.accounts.some(a => a.code.toLowerCase() === codeTrimmed.toLowerCase())
  if (exists) {
    accountingStore.showToast(`Kode Akun ${codeTrimmed} sudah digunakan! Silakan gunakan kode unik.`, 'error')
    return
  }

  await ledgerStore.addAccount({
    code: codeTrimmed,
    name: newName.value.trim(),
    type: newType.value,
    normalBalance: newNormalBalance.value,
    level: newLevel.value,
    isHeader: isHeaderMode.value,
    businessUnit: newUnit.value,
    parentCode: selectedAddParentCode.value || selectedAddCategoryCode.value,
    taxCategory: newTaxCategory.value !== 'NON_TAX' ? (newTaxCategory.value as any) : undefined,
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
    if (selectedType.value !== 'ALL' && a.type !== selectedType.value) {
      return false
    }
    if (selectedSubCategory.value !== 'ALL') {
      const target = selectedSubCategory.value
      if (a.code !== target && a.parentCode !== target) {
        let p = a.parentCode
        let match = false
        while (p) {
          if (p === target) { match = true; break }
          const pAcc = ledgerStore.accounts.find(x => x.code === p)
          p = pAcc?.parentCode
        }
        if (!match) return false
      }
    }
    if (selectedLevel.value !== 'ALL') {
      if (selectedLevel.value === 4) {
        // Level 4 (Detail): tampilkan semua akun detail / non-header
        if (a.isHeader) return false
        if (a.level < 3) return false
      } else if (selectedLevel.value === 3) {
        // Level 3 (Sub-Kelompok Header): tampilkan header level 3
        if (a.level !== 3) return false
        if (!a.isHeader && ledgerStore.accounts.some(x => x.level === 4)) return false
      } else if (a.level !== selectedLevel.value) {
        return false
      }
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
            {{ localeStore.t('coa.category') }}
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">{{ localeStore.t('coa.subtitle') }}</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          {{ localeStore.t('coa.title') }}
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="router.push('/accounting/opening-balance')"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <BookmarkCheck class="w-3.5 h-3.5 text-zinc-600" />
          <span>{{ localeStore.t('coa.setOpeningBalance') }}</span>
        </button>

        <button
          @click="openAddModal()"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>{{ localeStore.t('coa.addNewAccount') }}</span>
        </button>
      </div>
    </header>

    <!-- Toolbar Filters & Search & Level Selectors -->
    <div class="bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto">
        <!-- Search -->
        <div class="relative w-full md:w-64">
          <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="localeStore.t('coa.searchPlaceholder')"
            class="w-full pl-9 pr-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs focus:ring-1 focus:ring-zinc-950 focus:outline-none"
          />
        </div>

        <!-- Dynamic Category Selector (Level 1 dari COA) -->
        <div class="relative">
          <select
            v-model="selectedType"
            @change="selectedSubCategory = 'ALL'"
            class="px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-bold text-zinc-800 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            :title="localeStore.t('coa.fieldCategory')"
          >
            <option
              v-for="opt in categoryOptions"
              :key="opt.code"
              :value="opt.value"
            >
              {{ opt.label }}
            </option>
          </select>
        </div>

        <!-- Dynamic Sub-Category Selector (Level 2 dari COA) -->
        <div class="relative">
          <select
            v-model="selectedSubCategory"
            class="px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium text-zinc-800 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            :title="localeStore.t('coa.fieldParent')"
          >
            <option
              v-for="sub in subCategoryOptions"
              :key="sub.code"
              :value="sub.value"
            >
              {{ sub.label }}
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
            {{ localeStore.t('coa.levelAll') }}
          </button>
          <button
            @click="selectedLevel = 1"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 1 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
            title="Index Utama (Aset, Kewajiban, Ekuitas, Pendapatan, Beban)"
          >
            {{ localeStore.t('coa.level1') }}
          </button>
          <button
            @click="selectedLevel = 2"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 2 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
            title="Sub-Index Kelompok (Kas & Bank, Piutang, Hutang Pajak, dll.)"
          >
            {{ localeStore.t('coa.level2') }}
          </button>
          <button
            @click="selectedLevel = 3"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 3 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
            title="Sub-Index ke-2 (Kasir, Bank BCA, PPh 21, PPh 23, dll.)"
          >
            {{ localeStore.t('coa.level3') }}
          </button>
          <button
            @click="selectedLevel = 4"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-bold transition-all',
              selectedLevel === 4 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
            title="Rincian Detail Rekening Buku Besar"
          >
            {{ localeStore.t('coa.level4') }}
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
            {{ localeStore.t('coa.statusAll') }}
          </button>
          <button
            @click="selectedStatus = 'ACTIVE'"
            :class="[
              'px-2.5 py-1 rounded-lg font-bold transition-all',
              selectedStatus === 'ACTIVE' ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            {{ localeStore.t('coa.statusActive') }}
          </button>
          <button
            @click="selectedStatus = 'INACTIVE'"
            :class="[
              'px-2.5 py-1 rounded-lg font-bold transition-all',
              selectedStatus === 'INACTIVE' ? 'bg-rose-600 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            {{ localeStore.t('coa.statusInactive') }}
          </button>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <div class="flex items-center gap-1.5">
          <button
            type="button"
            @click="expandAll"
            class="px-2.5 py-1 text-xs font-semibold text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100 rounded-lg border border-zinc-200 transition-colors cursor-pointer"
            :title="localeStore.t('coa.expandAll')"
          >
            {{ localeStore.t('coa.expandAll') }}
          </button>
          <button
            type="button"
            @click="collapseAll"
            class="px-2.5 py-1 text-xs font-semibold text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100 rounded-lg border border-zinc-200 transition-colors cursor-pointer"
            :title="localeStore.t('coa.collapseAll')"
          >
            {{ localeStore.t('coa.collapseAll') }}
          </button>
        </div>
        <div class="text-xs text-zinc-500 font-medium">
          {{ localeStore.t('coa.showingAccounts', 'Menampilkan {count} akun').replace('{count}', String(filteredAccounts.length)) }}
        </div>
      </div>
    </div>

    <!-- Tree-like Table with Clear Visual Hierarchy -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full min-w-[850px] text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200/90 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-4 py-3 w-36">{{ localeStore.t('coa.colCode') }}</th>
              <th class="px-4 py-3">{{ localeStore.t('coa.colName') }}</th>
              <th class="px-4 py-3 text-center w-20">{{ localeStore.t('coa.colLevel') }}</th>
              <th class="px-4 py-3">{{ localeStore.t('coa.colType') }}</th>
              <th class="px-4 py-3 text-center">{{ localeStore.t('coa.colNormalBalance') }}</th>
              <th class="px-4 py-3 text-center">{{ localeStore.t('coa.colUnit') }}</th>
              <th class="px-4 py-3 text-right">{{ localeStore.t('coa.colBalance') }}</th>
              <th class="px-4 py-3 text-center w-28">{{ localeStore.t('coa.colTax') }}</th>
              <th class="px-4 py-3 text-center w-32">{{ localeStore.t('coa.colActions') }}</th>
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
                    :title="collapsedCodes.has(acc.code) ? localeStore.t('coa.expandAll') : localeStore.t('coa.collapseAll')"
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
                    {{ localeStore.translateAccountName(acc.name) }}
                  </span>
                  <span
                    v-if="acc.isActive === false"
                    class="ml-2 px-1.5 py-0.5 rounded text-[9px] font-bold bg-rose-50 text-rose-700 border border-rose-200 uppercase tracking-wider"
                  >
                    {{ localeStore.t('coa.badgeInactive') }}
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
                    'bg-emerald-50 text-emerald-800 border border-emerald-200/70'
                  ]"
                >
                  L{{ acc.level || 4 }}
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

              <!-- Aksi (Buku Besar, Edit, Nonaktifkan, + Tambah Sub-Akun) -->
              <td class="px-4 py-2.5 text-center">
                <div class="flex items-center justify-center gap-1.5">
                  <!-- Tombol Cepat Tambah Sub-Akun pada Header / Kategori -->
                  <button
                    v-if="acc.isHeader"
                    @click.stop="openAddModal(acc)"
                    class="p-1.5 bg-emerald-50 hover:bg-emerald-100 text-emerald-700 hover:text-emerald-900 rounded-lg transition-colors border border-emerald-200 shadow-2xs cursor-pointer"
                    title="Tambah Sub-Akun langsung di bawah kelompok ini"
                  >
                    <Plus class="w-3.5 h-3.5" />
                  </button>

                  <button
                    @click="openLedgerForAccount(acc.code)"
                    class="p-1.5 bg-white hover:bg-zinc-100 text-zinc-600 hover:text-zinc-950 rounded-lg transition-colors border border-zinc-200 shadow-2xs cursor-pointer"
                    :title="localeStore.t('coa.actionLedger')"
                  >
                    <BookOpen class="w-3.5 h-3.5" />
                  </button>

                  <button
                    v-if="!acc.isHeader"
                    @click="openEditModal(acc)"
                    class="p-1.5 bg-white hover:bg-zinc-100 text-zinc-600 hover:text-zinc-950 rounded-lg transition-colors border border-zinc-200 shadow-2xs cursor-pointer"
                    :title="localeStore.t('coa.actionEdit')"
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
                    :title="acc.isActive !== false ? localeStore.t('coa.actionDeactivate') : localeStore.t('coa.actionActivate')"
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



    <!-- Add Account Modal (Dynamic Category & Sub-Category Hierarchy) -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-black/40 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="p-5 border-b border-zinc-100 flex items-center justify-between">
          <div>
            <h3 class="text-sm font-bold text-zinc-950">Tambah Akun Bagan Baru (COA)</h3>
            <p class="text-[11px] text-zinc-500 mt-0.5">Penetapan terstruktur per Kategori & Sub-Kategori Induk</p>
          </div>
          <button @click="showAddModal = false" class="p-1.5 text-zinc-400 hover:text-zinc-700 rounded-lg hover:bg-zinc-100 transition-colors cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleSaveAccount" class="p-6 space-y-4 text-xs">
          <!-- Toggle: Akun Transaksi vs Sub-Kategori Baru -->
          <div class="bg-zinc-100 p-1 rounded-xl flex items-center text-xs font-bold">
            <button
              type="button"
              @click="isHeaderMode = false; onHeaderModeChanged()"
              :class="[
                'flex-1 py-1.5 px-3 rounded-lg text-center transition cursor-pointer',
                !isHeaderMode ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              📄 Akun Transaksi (Detail Buku Besar)
            </button>
            <button
              type="button"
              @click="isHeaderMode = true; onHeaderModeChanged()"
              :class="[
                'flex-1 py-1.5 px-3 rounded-lg text-center transition cursor-pointer',
                isHeaderMode ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              📁 Sub-Kategori / Header Baru
            </button>
          </div>

          <!-- 1. Kategori Utama (Level 1) -->
          <div>
            <label class="block font-bold text-zinc-800 mb-1">
              1. Kategori Utama (Induk Level 1) <span class="text-rose-500">*</span>
            </label>
            <select
              v-model="selectedAddCategoryCode"
              @change="onCategoryChanged"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option
                v-for="cat in addCategoryList"
                :key="cat.code"
                :value="cat.code"
              >
                {{ cat.code }} - {{ cat.name }} ({{ cat.type }})
              </option>
            </select>
          </div>

          <!-- 2. Sub-Kategori Induk (Level 2/3) - Only for Detail Account Mode -->
          <div v-if="!isHeaderMode">
            <label class="block font-bold text-zinc-800 mb-1">
              2. Sub-Kategori Induk (Parent Group) <span class="text-rose-500">*</span>
            </label>
            <select
              v-model="selectedAddParentCode"
              @change="onParentChanged"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-semibold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option
                v-for="sub in addSubCategoryList"
                :key="sub.code"
                :value="sub.code"
              >
                {{ sub.code }} - {{ sub.name }} (Level {{ sub.level }})
              </option>
            </select>
            <p class="text-[10px] text-zinc-500 mt-1">Akun baru akan otomatis berada di bawah sub-kategori ini pada buku besar.</p>
          </div>

          <!-- 3. Kode & Nama Akun -->
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label class="block font-bold text-zinc-800 mb-1">
                Kode Akun <span class="text-rose-500">*</span>
              </label>
              <input
                v-model="newCode"
                type="text"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-mono font-bold focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                required
              />
              <span class="text-[9px] text-emerald-700 font-semibold block mt-0.5">Saran otomatis sistem</span>
            </div>

            <div class="sm:col-span-2">
              <label class="block font-bold text-zinc-800 mb-1">
                Nama Akun <span class="text-rose-500">*</span>
              </label>
              <input
                v-model="newName"
                type="text"
                placeholder="Contoh: Kas Kasir Kidinest Shift Pagi"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-bold focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                required
              />
            </div>
          </div>

          <!-- 4. Otomatis: Klasifikasi, Saldo Normal, Unit Bisnis -->
          <div class="grid grid-cols-3 gap-3 bg-zinc-50 p-3 rounded-2xl border border-zinc-200">
            <div>
              <span class="block text-[10px] font-bold text-zinc-500 uppercase">Tipe Akun</span>
              <span class="font-bold text-xs text-zinc-900">{{ newType }}</span>
            </div>
            <div>
              <span class="block text-[10px] font-bold text-zinc-500 uppercase">Saldo Normal</span>
              <span
                :class="[
                  'font-bold text-xs inline-block px-1.5 py-0.5 rounded text-[10px]',
                  newNormalBalance === 'DEBIT' ? 'bg-emerald-100 text-emerald-900' : 'bg-rose-100 text-rose-900'
                ]"
              >
                {{ newNormalBalance }}
              </span>
            </div>
            <div>
              <span class="block text-[10px] font-bold text-zinc-500 uppercase">Hierarki</span>
              <span class="font-bold text-xs text-zinc-900">Level {{ newLevel }}</span>
            </div>
          </div>

          <!-- 5. Unit Bisnis & Integrasi Pajak -->
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-zinc-800 mb-1">Unit Bisnis</label>
              <select
                v-model="newUnit"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              >
                <option value="ALL">Semua Unit (Konsolidasi)</option>
                <option value="KIDINEST">Kidinest Pediatric</option>
                <option value="SSC">Semarang Skin Center</option>
                <option value="MINIMO">Minimo</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-zinc-800 mb-1">Integrasi Pajak (Opsional)</label>
              <select
                v-model="newTaxCategory"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              >
                <option value="NON_TAX">Non-Pajak / Umum</option>
                <option value="PPH_21">PPh Pasal 21 (Gaji & Fee Dokter)</option>
                <option value="PPH_23">PPh Pasal 23 (Jasa Vendor / IT)</option>
                <option value="PPH_4_2">PPh Final Pasal 4(2) (Sewa Gedung)</option>
                <option value="PPN_OUT">PPN Keluaran</option>
                <option value="PPN_IN">PPN Masukan</option>
                <option value="PPH_25">PPh Pasal 25 (Angsuran Badan)</option>
                <option value="PPH_29">PPh Pasal 29 (Kurang Bayar Tahunan)</option>
              </select>
            </div>
          </div>

          <!-- 6. Keterangan -->
          <div>
            <label class="block font-bold text-zinc-800 mb-1">Keterangan / Catatan</label>
            <textarea
              v-model="newDescription"
              rows="2"
              placeholder="Deskripsi peruntukan akun buku besar ini..."
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

