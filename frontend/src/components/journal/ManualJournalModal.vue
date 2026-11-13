<script setup lang="ts">
import { ref, computed } from 'vue'
import { useJournalsStore } from '@/stores/journals'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import type { BusinessUnit } from '@/types/accounting'
import { formatIDR } from '@/utils/currency'
import { Plus, Trash2, X, AlertCircle, CheckCircle2 } from 'lucide-vue-next'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const journalsStore = useJournalsStore()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

const journalDate = ref(new Date().toISOString().split('T')[0])
const selectedUnit = ref<BusinessUnit>('KIDINEST')
const branchId = ref('BRANCH-MAIN')
const description = ref('')

interface LineInput {
  accountId: string
  accountCode: string
  accountName: string
  debit: number
  credit: number
  description: string
}

const lines = ref<LineInput[]>([
  {
    accountId: '1-1100',
    accountCode: '1-1100',
    accountName: 'Kas Operasional Kidinest',
    debit: 0,
    credit: 0,
    description: ''
  },
  {
    accountId: '4-1100',
    accountCode: '4-1100',
    accountName: 'Pendapatan Jasa Medis Kidinest',
    debit: 0,
    credit: 0,
    description: ''
  }
])

const availableAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader)
})

function onAccountSelect(index: number, code: string) {
  const acc = availableAccounts.value.find(a => a.code === code)
  if (acc) {
    lines.value[index].accountId = acc.id
    lines.value[index].accountCode = acc.code
    lines.value[index].accountName = acc.name
  }
}

function addLine() {
  lines.value.push({
    accountId: '',
    accountCode: '',
    accountName: '',
    debit: 0,
    credit: 0,
    description: ''
  })
}

function removeLine(index: number) {
  if (lines.value.length > 2) {
    lines.value.splice(index, 1)
  }
}

const totalDebit = computed(() => {
  return lines.value.reduce((s, l) => s + (Number(l.debit) || 0), 0)
})

const totalCredit = computed(() => {
  return lines.value.reduce((s, l) => s + (Number(l.credit) || 0), 0)
})

const isBalanced = computed(() => {
  return totalDebit.value > 0 && totalDebit.value === totalCredit.value
})

const balanceDifference = computed(() => {
  return Math.abs(totalDebit.value - totalCredit.value)
})

function handleSubmit() {
  if (!isBalanced.value) {
    accountingStore.showToast('Jurnal belum seimbang! Total Debit harus sama dengan Total Kredit.', 'error')
    return
  }

  if (!description.value.trim()) {
    accountingStore.showToast('Harap isi keterangan memo jurnal.', 'error')
    return
  }

  const branchName = branchId.value === 'BRANCH-MAIN' ? 'Kebayoran Baru (Pusat)' : 'Pondok Indah'

  try {
    journalsStore.createManualJournal({
      journalDate: journalDate.value,
      businessUnit: selectedUnit.value,
      branchId: branchId.value,
      branchName,
      description: description.value.trim(),
      lines: lines.value.map(l => ({
        accountId: l.accountId || l.accountCode,
        accountCode: l.accountCode,
        accountName: l.accountName,
        debit: Number(l.debit) || 0,
        credit: Number(l.credit) || 0,
        description: l.description || description.value.trim()
      }))
    })

    emit('close')
  } catch (err: any) {
    accountingStore.showToast(err.message, 'error')
  }
}
</script>

<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-4xl w-full overflow-hidden flex flex-col max-h-[90vh] animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Header -->
      <div class="px-6 py-4 bg-slate-900 text-white flex items-center justify-between">
        <div>
          <h3 class="text-base font-bold text-white">Buat Jurnal Penyesuaian Manual</h3>
          <p class="text-xs text-slate-400">Pencatatan akuntansi manual dengan validasi keseimbangan debit & kredit</p>
        </div>
        <button
          @click="emit('close')"
          class="p-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition-colors"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 overflow-y-auto space-y-6 text-xs text-slate-700">
        <!-- Form Controls -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block font-bold text-slate-700 mb-1 text-[11px] uppercase">Tanggal Jurnal</label>
            <input
              v-model="journalDate"
              type="date"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
            />
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1 text-[11px] uppercase">Unit Bisnis</label>
            <select
              v-model="selectedUnit"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
            >
              <option value="KIDINEST">Kidinest (Klinik Anak)</option>
              <option value="SSC">SSC (Skin & Slimming Clinic)</option>
              <option value="ALL">ALL (Biaya Bersama / Kantor Pusat)</option>
            </select>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1 text-[11px] uppercase">Cabang</label>
            <select
              v-model="branchId"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
            >
              <option value="BRANCH-MAIN">Kebayoran Baru (Pusat)</option>
              <option value="BRANCH-JAKSEL">Pondok Indah</option>
              <option value="BRANCH-SBY">Surabaya Barat</option>
            </select>
          </div>
        </div>

        <div>
          <label class="block font-bold text-slate-700 mb-1 text-[11px] uppercase">Uraian / Memo Jurnal</label>
          <input
            v-model="description"
            type="text"
            placeholder="Contoh: Penyesuaian biaya administrasi bank akhir bulan / Kas kecil operasional"
            class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <!-- Dynamic Lines Table -->
        <div class="border border-slate-200 rounded-xl overflow-hidden shadow-2xs">
          <div class="bg-slate-100 px-4 py-2 flex items-center justify-between">
            <span class="font-bold text-slate-800">Daftar Akun Jurnal (Double-Entry)</span>
            <button
              @click="addLine"
              type="button"
              class="flex items-center gap-1 px-2.5 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-xs font-semibold shadow-xs"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Tambah Baris</span>
            </button>
          </div>

          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-50 text-[11px] font-bold text-slate-600 uppercase border-b border-slate-200">
                <th class="px-3 py-2.5 w-72">Pilih Akun (COA)</th>
                <th class="px-3 py-2.5">Keterangan Baris</th>
                <th class="px-3 py-2.5 text-right w-36">Debit (Rp)</th>
                <th class="px-3 py-2.5 text-right w-36">Kredit (Rp)</th>
                <th class="px-3 py-2.5 text-center w-12">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="(line, idx) in lines" :key="idx" class="hover:bg-slate-50/50">
                <td class="px-3 py-2">
                  <select
                    :value="line.accountCode"
                    @change="(e: any) => onAccountSelect(idx, e.target.value)"
                    class="w-full px-2 py-1.5 bg-white border border-slate-300 rounded text-xs focus:ring-1 focus:ring-blue-500 font-medium"
                  >
                    <option value="" disabled>-- Pilih Akun --</option>
                    <option
                      v-for="acc in availableAccounts"
                      :key="acc.code"
                      :value="acc.code"
                    >
                      {{ acc.code }} - {{ acc.name }}
                    </option>
                  </select>
                </td>
                <td class="px-3 py-2">
                  <input
                    v-model="line.description"
                    type="text"
                    placeholder="Uraian akun..."
                    class="w-full px-2 py-1.5 bg-white border border-slate-300 rounded text-xs"
                  />
                </td>
                <td class="px-3 py-2 text-right">
                  <input
                    v-model.number="line.debit"
                    type="number"
                    min="0"
                    step="1000"
                    class="w-full px-2 py-1.5 text-right bg-white border border-slate-300 rounded font-mono text-xs"
                  />
                </td>
                <td class="px-3 py-2 text-right">
                  <input
                    v-model.number="line.credit"
                    type="number"
                    min="0"
                    step="1000"
                    class="w-full px-2 py-1.5 text-right bg-white border border-slate-300 rounded font-mono text-xs"
                  />
                </td>
                <td class="px-3 py-2 text-center">
                  <button
                    @click="removeLine(idx)"
                    :disabled="lines.length <= 2"
                    type="button"
                    class="p-1 text-slate-400 hover:text-rose-600 disabled:opacity-30 disabled:pointer-events-none"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="bg-slate-100 font-bold border-t-2 border-slate-300">
                <td colspan="2" class="px-4 py-2.5 text-right text-slate-700">
                  TOTAL:
                </td>
                <td class="px-4 py-2.5 text-right font-mono text-slate-900">
                  {{ formatIDR(totalDebit) }}
                </td>
                <td class="px-4 py-2.5 text-right font-mono text-slate-900">
                  {{ formatIDR(totalCredit) }}
                </td>
                <td></td>
              </tr>
            </tfoot>
          </table>
        </div>

        <!-- Balance Indicator Banner -->
        <div
          v-if="isBalanced"
          class="flex items-center gap-2 p-3 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-800 text-xs font-semibold"
        >
          <CheckCircle2 class="w-4 h-4 text-emerald-600" />
          <span>Posisi Jurnal Seimbang (Total Debit = Total Kredit: {{ formatIDR(totalDebit) }})</span>
        </div>
        <div
          v-else
          class="flex items-center gap-2 p-3 bg-amber-50 border border-amber-200 rounded-xl text-amber-800 text-xs font-medium"
        >
          <AlertCircle class="w-4 h-4 text-amber-600" />
          <span>
            Jurnal Belum Seimbang! Selisih: <strong class="font-mono">{{ formatIDR(balanceDifference) }}</strong>.
            Lengkapi nilai Debit dan Kredit hingga sama.
          </span>
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="px-6 py-4 bg-slate-50 border-t border-slate-200 flex items-center justify-end gap-3">
        <button
          @click="emit('close')"
          type="button"
          class="px-4 py-2 bg-white border border-slate-300 hover:bg-slate-100 text-slate-700 rounded-lg text-xs font-semibold"
        >
          Batal
        </button>
        <button
          @click="handleSubmit"
          :disabled="!isBalanced"
          type="button"
          class="px-5 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          Simpan Draft Jurnal Manual
        </button>
      </div>
    </div>
  </div>
</template>
