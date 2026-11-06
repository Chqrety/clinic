<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Drug } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { Package, AlertTriangle, AlertCircle, Plus, RefreshCw, CheckCircle2 } from 'lucide-vue-next'

const store = useClinicStore()
const drugs = ref<Drug[]>([])
const loading = ref(false)
const selectedDrug = ref<Drug | null>(null)

// Restock modal state
const showRestockModal = ref(false)
const restockBatch = ref('')
const restockED = ref('')
const restockQty = ref(50)
const restockBuyPrice = ref(1000)
const submittingRestock = ref(false)

async function fetchDrugs() {
  loading.value = true
  try {
    drugs.value = await api.getDrugs()
    if (selectedDrug.value) {
      selectedDrug.value = await api.getDrugDetail(selectedDrug.value.id)
    }
  } catch (e) {
    // ignore
  } finally {
    loading.value = false
  }
}

async function viewDetail(drug: Drug) {
  try {
    selectedDrug.value = await api.getDrugDetail(drug.id)
  } catch (e) {
    selectedDrug.value = drug
  }
}

function openRestock(drug: Drug) {
  selectedDrug.value = drug
  const nextYear = new Date()
  nextYear.setFullYear(nextYear.getFullYear() + 2)
  restockED.value = nextYear.toISOString().split('T')[0]
  restockBatch.value = `BCH-${new Date().getFullYear()}-${Math.floor(100 + Math.random() * 900)}`
  restockQty.value = 50
  restockBuyPrice.value = Math.floor(drug.sell_price * 0.6)
  showRestockModal.value = true
}

async function submitRestock() {
  if (!selectedDrug.value) return
  submittingRestock.value = true
  try {
    await api.restockDrug(selectedDrug.value.id, {
      batch_number: restockBatch.value,
      expired_date: restockED.value,
      quantity: restockQty.value,
      buy_price: restockBuyPrice.value,
    })
    store.playChime()
    showRestockModal.value = false
    await fetchDrugs()
  } catch (err: any) {
    alert(err.message || 'Gagal menambah stok obat')
  } finally {
    submittingRestock.value = false
  }
}

onMounted(() => {
  fetchDrugs()
})
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 px-4 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
      <div>
        <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold bg-zinc-950 text-white mb-2 shadow-sm">
          Gudang & Inventori Apotek
        </span>
        <h1 class="text-3xl font-extrabold text-zinc-950 tracking-tight">
          Monitoring Stok, Batch & Expired Date
        </h1>
        <p class="text-sm text-zinc-500 mt-1">
          Formula: <code class="px-1.5 py-0.5 rounded bg-stone-200/60 font-mono text-zinc-900">Stok Tersedia = Stok Fisik - Stok Reserved</code>. Prioritas kadaluarsa (FEFO).
        </p>
      </div>
      <button
        @click="fetchDrugs"
        :disabled="loading"
        class="inline-flex items-center gap-2 px-5 py-2.5 rounded-full bg-white border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 shadow-sm transition-all active:scale-95 disabled:opacity-50 self-start sm:self-auto"
      >
        <RefreshCw :class="['w-3.5 h-3.5', loading && 'animate-spin']" />
        Segarkan Stok
      </button>
    </div>

    <!-- Alert Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
      <div class="bg-rose-50/70 border border-rose-200 rounded-3xl p-5 flex items-center justify-between">
        <div>
          <div class="text-xs font-bold uppercase tracking-wider text-rose-700">Mendekati Kadaluarsa (&lt; 3 Bulan)</div>
          <div class="text-3xl font-extrabold text-rose-900 mt-1">
            {{ drugs.filter(d => d.has_expiring).length }} Jenis Obat
          </div>
        </div>
        <div class="w-10 h-10 rounded-2xl bg-rose-200/50 flex items-center justify-center text-rose-800">
          <AlertCircle class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-amber-50/70 border border-amber-200 rounded-3xl p-5 flex items-center justify-between">
        <div>
          <div class="text-xs font-bold uppercase tracking-wider text-amber-700">Peringatan Stok Menipis</div>
          <div class="text-3xl font-extrabold text-amber-900 mt-1">
            {{ drugs.filter(d => d.is_low_stock).length }} Jenis Obat
          </div>
        </div>
        <div class="w-10 h-10 rounded-2xl bg-amber-200/50 flex items-center justify-center text-amber-800">
          <AlertTriangle class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white border border-black/5 rounded-3xl p-5 flex items-center justify-between shadow-[0_2px_12px_rgba(0,0,0,0.03)]">
        <div>
          <div class="text-xs font-bold uppercase tracking-wider text-zinc-400">Total Master Obat</div>
          <div class="text-3xl font-extrabold text-zinc-950 mt-1">{{ drugs.length }} Item</div>
        </div>
        <div class="w-10 h-10 rounded-2xl bg-stone-100 flex items-center justify-center text-zinc-800">
          <Package class="w-5 h-5" />
        </div>
      </div>
    </div>

    <!-- Main Table -->
    <div class="bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-sm">
          <thead>
            <tr class="border-b border-zinc-100 text-[11px] font-bold uppercase tracking-wider text-zinc-400">
              <th class="pb-3 pl-2">Nama Obat</th>
              <th class="pb-3">Kategori & KFA</th>
              <th class="pb-3 text-center">Stok Fisik</th>
              <th class="pb-3 text-center">Stok Dipesan (Reserved)</th>
              <th class="pb-3 text-center">Stok Tersedia</th>
              <th class="pb-3">Status Peringatan</th>
              <th class="pb-3 text-right pr-2">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-50">
            <tr
              v-for="d in drugs"
              :key="d.id"
              class="hover:bg-stone-50/60 transition-colors"
            >
              <td class="py-4 pl-2 font-bold text-zinc-950">
                {{ d.name }}
              </td>
              <td class="py-4">
                <div class="text-xs font-semibold text-zinc-800">{{ d.category }}</div>
                <div class="text-[11px] font-mono text-zinc-400">KFA: {{ d.kfa_code }}</div>
              </td>
              <td class="py-4 text-center font-bold text-zinc-900">
                {{ d.physical_stock }} {{ d.unit }}
              </td>
              <td class="py-4 text-center font-bold text-amber-600">
                {{ d.reserved_stock }} {{ d.unit }}
              </td>
              <td class="py-4 text-center font-extrabold text-base text-zinc-950">
                <span :class="d.is_low_stock ? 'text-rose-600' : 'text-emerald-700'">
                  {{ d.available_stock }}
                </span>
                <span class="text-xs font-normal text-zinc-400 ml-1">{{ d.unit }}</span>
              </td>
              <td class="py-4">
                <div class="flex flex-col gap-1 items-start">
                  <span
                    v-if="d.has_expiring"
                    class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-rose-100 text-rose-800 border border-rose-200"
                  >
                    ⚠ ED &lt; 3 Bulan
                  </span>
                  <span
                    v-if="d.is_low_stock"
                    class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-amber-100 text-amber-800 border border-amber-200"
                  >
                    Stok Menipis
                  </span>
                  <span
                    v-if="!d.has_expiring && !d.is_low_stock"
                    class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-700"
                  >
                    Aman
                  </span>
                </div>
              </td>
              <td class="py-4 text-right pr-2 space-x-2 whitespace-nowrap">
                <button
                  @click="viewDetail(d)"
                  class="px-3 py-1.5 rounded-full border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 transition-all active:scale-95"
                >
                  Detail Batch
                </button>
                <button
                  @click="openRestock(d)"
                  class="px-3.5 py-1.5 rounded-full bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold transition-all active:scale-95"
                >
                  + Restock
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Batch Detail Drawer / Section -->
    <div v-if="selectedDrug && selectedDrug.batches" class="mt-8 bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)]">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h2 class="text-lg font-bold text-zinc-950">
            Rincian Batch Obat: {{ selectedDrug.name }}
          </h2>
          <p class="text-xs text-zinc-500">Urutan FEFO (Batch yang mendekati expired berada di atas).</p>
        </div>
        <button
          @click="selectedDrug = null"
          class="text-xs font-bold text-zinc-400 hover:text-zinc-900"
        >
          Tutup Rincian
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
        <div
          v-for="b in selectedDrug.batches"
          :key="b.id"
          :class="[
            'p-4 rounded-2xl border transition-all',
            b.is_expiring_soon ? 'border-rose-200 bg-rose-50/50' : 'border-black/5 bg-[#F9F9F7]'
          ]"
        >
          <div class="flex items-center justify-between mb-1">
            <span class="font-mono font-bold text-xs text-zinc-900">{{ b.batch_number }}</span>
            <span
              v-if="b.is_expiring_soon"
              class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-rose-200 text-rose-800"
            >
              ED Segera
            </span>
          </div>
          <div class="text-xs text-zinc-500">Expired: <strong class="text-zinc-900">{{ b.expired_date }}</strong></div>
          <div class="flex justify-between items-center mt-3 pt-2 border-t border-black/5 text-xs">
            <span class="text-zinc-500">Fisik: {{ b.physical_stock }}</span>
            <span class="text-amber-700 font-semibold">Reserved: {{ b.reserved_stock }}</span>
            <span class="font-extrabold text-zinc-900">Tersedia: {{ b.available_stock }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Restock -->
    <div
      v-if="showRestockModal && selectedDrug"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/30 backdrop-blur-sm"
    >
      <div class="bg-white rounded-3xl p-6 sm:p-8 max-w-md w-full shadow-2xl border border-black/10 animate-in zoom-in-95 duration-150">
        <h3 class="text-xl font-extrabold text-zinc-950 mb-1">Tambah Stok (Restock)</h3>
        <p class="text-xs text-zinc-500 mb-4">{{ selectedDrug.name }}</p>

        <form @submit.prevent="submitRestock" class="space-y-3">
          <div>
            <label class="block text-[11px] font-bold text-zinc-500 uppercase tracking-wider mb-1">Nomor Batch Faktur</label>
            <input
              v-model="restockBatch"
              type="text"
              required
              class="w-full px-4 py-2.5 rounded-2xl bg-[#F9F9F7] border border-black/10 text-xs font-bold text-zinc-900 outline-none"
            />
          </div>

          <div>
            <label class="block text-[11px] font-bold text-zinc-500 uppercase tracking-wider mb-1">Tanggal Kadaluarsa (ED)</label>
            <input
              v-model="restockED"
              type="date"
              required
              class="w-full px-4 py-2.5 rounded-2xl bg-[#F9F9F7] border border-black/10 text-xs font-bold text-zinc-900 outline-none"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-bold text-zinc-500 uppercase tracking-wider mb-1">Jumlah Masuk</label>
              <input
                v-model.number="restockQty"
                type="number"
                min="1"
                required
                class="w-full px-4 py-2.5 rounded-2xl bg-[#F9F9F7] border border-black/10 text-xs font-bold text-zinc-900 outline-none"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-zinc-500 uppercase tracking-wider mb-1">Harga Beli / Unit</label>
              <input
                v-model.number="restockBuyPrice"
                type="number"
                min="100"
                required
                class="w-full px-4 py-2.5 rounded-2xl bg-[#F9F9F7] border border-black/10 text-xs font-bold text-zinc-900 outline-none"
              />
            </div>
          </div>

          <div class="flex items-center justify-end gap-2 pt-4">
            <button
              type="button"
              @click="showRestockModal = false"
              class="px-4 py-2.5 rounded-full text-xs font-bold text-zinc-500 hover:text-zinc-900"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="submittingRestock"
              class="px-5 py-2.5 rounded-full bg-amber-400 hover:bg-amber-500 text-zinc-950 text-xs font-extrabold transition-all shadow-sm active:scale-95 disabled:opacity-50"
            >
              {{ submittingRestock ? 'Menyimpan...' : 'Simpan Masuk Gudang' }}
            </button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>
