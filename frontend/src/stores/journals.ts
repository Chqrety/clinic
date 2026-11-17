import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { JournalEntry, JournalStatus, BusinessUnit, JournalSource } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { accountingApi, type BackendJournal } from '@/api/accounting'

export const useJournalsStore = defineStore('journals', () => {
  const accountingStore = useAccountingStore()

  // Real journal entries fetched dynamically from PostgreSQL database via Go API
  const journals = ref<JournalEntry[]>([])

  // Filtered journals based on global filters (Business Unit & Fiscal Month)
  const filteredJournals = computed(() => {
    return journals.value.filter(j => {
      if (accountingStore.businessUnit !== 'ALL' && j.businessUnit !== 'ALL' && j.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      if (accountingStore.fiscalMonth && j.journalDate) {
        const str = String(j.journalDate).trim()
        const isoMatch = str.match(/^(\d{4})[-/](\d{1,2})[-/]/)
        const dmyMatch = str.match(/^\d{1,2}[-/](\d{1,2})[-/]\d{4}/)
        const m = isoMatch ? parseInt(isoMatch[2], 10) : (dmyMatch ? parseInt(dmyMatch[1], 10) : null)
        if (m !== null && m !== accountingStore.fiscalMonth) return false
      }
      return true
    })
  })

  // Statistics
  const stats = computed(() => {
    const list = filteredJournals.value
    return {
      totalCount: list.length,
      draftCount: list.filter(j => j.status === 'DRAFT').length,
      submittedCount: list.filter(j => j.status === 'SUBMITTED').length,
      reviewedCount: list.filter(j => j.status === 'REVIEWED').length,
      approvedCount: list.filter(j => j.status === 'APPROVED').length,
      postedCount: list.filter(j => j.status === 'POSTED').length,
      reversedCount: list.filter(j => j.status === 'REVERSED').length,
      totalPostedDebit: list
        .filter(j => j.status === 'POSTED')
        .reduce((sum, j) => sum + j.totalDebit, 0)
    }
  })

  // Actions
  function getJournalById(id: string): JournalEntry | undefined {
    return journals.value.find(j => j.id === id)
  }

  function submitDraft(id: string) {
    const j = getJournalById(id)
    if (j && j.status === 'DRAFT') {
      j.status = 'SUBMITTED'
      accountingStore.showToast(`Jurnal ${j.journalNumber} berhasil diajukan untuk review`, 'success')
    }
  }

  function reviewJournal(id: string) {
    const j = getJournalById(id)
    if (j && (j.status === 'SUBMITTED' || j.status === 'DRAFT')) {
      j.status = 'REVIEWED'
      j.reviewedBy = accountingStore.currentUserEmail
      j.reviewedByName = accountingStore.currentUserName
      j.reviewedAt = new Date().toISOString()
      accountingStore.showToast(`Jurnal ${j.journalNumber} telah direview`, 'success')
    }
  }

  function approveJournal(id: string) {
    const j = getJournalById(id)
    if (j && j.status === 'REVIEWED') {
      j.status = 'APPROVED'
      j.approvedBy = accountingStore.currentUserEmail
      j.approvedByName = accountingStore.currentUserName
      j.approvedAt = new Date().toISOString()
      accountingStore.showToast(`Jurnal ${j.journalNumber} telah disetujui (Approved)`, 'success')
    }
  }

  function postJournal(id: string) {
    const j = getJournalById(id)
    if (j && j.status === 'APPROVED') {
      j.status = 'POSTED'
      j.postedBy = accountingStore.currentUserEmail
      j.postedByName = accountingStore.currentUserName
      j.postedAt = new Date().toISOString()
      accountingStore.showToast(`Jurnal ${j.journalNumber} resmi diposting ke Buku Besar!`, 'success')
    }
  }

  function rejectJournal(id: string, reason: string) {
    const j = getJournalById(id)
    if (j) {
      j.status = 'REJECTED'
      j.reversalReason = reason
      accountingStore.showToast(`Jurnal ${j.journalNumber} ditolak. Alasan: ${reason}`, 'info')
    }
  }

  function reverseJournal(id: string, reason: string) {
    const original = getJournalById(id)
    if (!original || original.status !== 'POSTED') {
      accountingStore.showToast('Hanya jurnal yang sudah berstatus POSTED yang dapat dibalik', 'error')
      return
    }

    const revNumber = `REV-${original.journalNumber}`
    const revJournal: JournalEntry = {
      id: 'jv-rev-' + Date.now(),
      journalNumber: revNumber,
      journalDate: new Date().toISOString().split('T')[0],
      journalSource: original.journalSource,
      sourceReference: `Reversal of ${original.journalNumber}`,
      businessUnit: original.businessUnit,
      branchId: original.branchId,
      branchName: original.branchName,
      description: `[JURNAL PEMBALIK] ${original.description} - Alasan: ${reason}`,
      totalDebit: original.totalCredit,
      totalCredit: original.totalDebit,
      status: 'POSTED',
      lines: original.lines.map((l, idx) => ({
        id: `jl-rev-${idx}-${Date.now()}`,
        accountId: l.accountId,
        accountCode: l.accountCode,
        accountName: l.accountName,
        // Swap debit and credit for reversal
        debit: l.credit,
        credit: l.debit,
        description: `Pembalikan: ${l.description}`
      })),
      createdBy: accountingStore.currentUserEmail,
      createdByName: accountingStore.currentUserName,
      createdAt: new Date().toISOString(),
      postedBy: accountingStore.currentUserEmail,
      postedByName: accountingStore.currentUserName,
      postedAt: new Date().toISOString(),
      reversalReason: reason,
      isManual: true
    }

    original.status = 'REVERSED'
    original.reversalJournalId = revJournal.id
    original.reversalReason = reason

    journals.value.unshift(revJournal)
    accountingStore.showToast(`Jurnal pembalik ${revNumber} berhasil diposting`, 'success')
  }

  function createManualJournal(data: {
    journalDate: string
    businessUnit: BusinessUnit
    branchId?: string
    branchName?: string
    description: string
    journalSource?: JournalSource
    sourceReference?: string
    status?: JournalStatus
    lines: {
      accountId?: string
      accountCode: string
      accountName: string
      debit: number
      credit: number
      description?: string
    }[]
  }): JournalEntry {
    const totalDebit = data.lines.reduce((s, l) => s + (l.debit || 0), 0)
    const totalCredit = data.lines.reduce((s, l) => s + (l.credit || 0), 0)

    if (Math.abs(totalDebit - totalCredit) > 0.01) {
      throw new Error(`Total Debit (Rp ${totalDebit}) tidak seimbang dengan Total Kredit (Rp ${totalCredit})`)
    }

    const nextSeq = journals.value.length + 1
    const seqStr = String(nextSeq).padStart(3, '0')
    const jNumber = `JV-2026-03-${seqStr}`
    const st = data.status || 'DRAFT'

    const newJournal: JournalEntry = {
      id: 'jv-' + Date.now(),
      journalNumber: jNumber,
      journalDate: data.journalDate,
      journalSource: data.journalSource || 'MANUAL',
      sourceReference: data.sourceReference || `MAN-${Date.now().toString().slice(-6)}`,
      businessUnit: data.businessUnit,
      branchId: data.branchId || 'BRANCH-MAIN',
      branchName: data.branchName || 'Kantor Pusat Semarang',
      description: data.description,
      totalDebit,
      totalCredit,
      status: st,
      lines: data.lines.map((l, idx) => ({
        id: `jl-${idx}-${Date.now()}`,
        accountId: l.accountId || `acc-${l.accountCode}`,
        accountCode: l.accountCode,
        accountName: l.accountName,
        debit: l.debit || 0,
        credit: l.credit || 0,
        description: l.description || data.description
      })),
      createdBy: accountingStore.currentUserEmail,
      createdByName: accountingStore.currentUserName,
      createdAt: new Date().toISOString(),
      postedBy: st === 'POSTED' ? accountingStore.currentUserEmail : undefined,
      postedByName: st === 'POSTED' ? accountingStore.currentUserName : undefined,
      postedAt: st === 'POSTED' ? new Date().toISOString() : undefined,
      isManual: true
    }

    journals.value.unshift(newJournal)
    accountingStore.showToast(`Jurnal ${jNumber} berhasil dibuat (${st})!`, 'success')

    // Persist to PostgreSQL database asynchronously
    accountingApi.createJournal({
      journal_type: 'GENERAL',
      transaction_date: data.journalDate,
      reference_type: data.journalSource || 'MANUAL',
      reference_id: data.sourceReference || jNumber,
      description: data.description,
      auto_post: st === 'POSTED',
      lines: data.lines.map(l => ({
        account_code: l.accountCode,
        debit: l.debit || 0,
        credit: l.credit || 0,
        description: l.description || data.description
      }))
    }).catch(err => {
      console.warn('Note: Journal created in memory, DB sync response:', err.message)
    })

    return newJournal
  }

  const isLoading = ref(false)

  // Fetch journals directly from PostgreSQL database via Go API
  async function fetchJournals() {
    try {
      isLoading.value = true
      const backendJournals = await accountingApi.getJournals({ limit: 100 })
      if (backendJournals && backendJournals.length > 0) {
        journals.value = backendJournals.map(j => {
          let bu: BusinessUnit = 'ALL'
          const descUpper = (j.description || '').toUpperCase()
          if (descUpper.includes('KIDINEST')) bu = 'KIDINEST'
          else if (descUpper.includes('SSC') || descUpper.includes('SEMARANG SKIN CENTER')) bu = 'SSC'

          return {
            id: j.id,
            journalNumber: j.journal_number,
            journalDate: j.transaction_date,
            journalSource: (j.reference_type as JournalSource) || 'MANUAL',
            sourceReference: j.reference_id || j.journal_number,
            businessUnit: bu,
            branchId: j.branch_id || 'BRANCH-MAIN',
            branchName: 'Kantor Pusat Semarang',
            description: j.description,
            totalDebit: Number(j.total_debit),
            totalCredit: Number(j.total_credit),
            status: (j.status as JournalStatus) || 'POSTED',
            lines: (j.lines || []).map((l, idx) => ({
              id: l.id || `jl-${idx}`,
              accountId: l.account_id,
              accountCode: l.account_code,
              accountName: l.account_name,
              debit: Number(l.debit),
              credit: Number(l.credit),
              description: l.description || j.description
            })),
            createdBy: accountingStore.currentUserEmail,
            createdByName: accountingStore.currentUserName,
            createdAt: j.created_at || new Date().toISOString(),
            isManual: true
          }
        })
      }
    } catch (err: any) {
      console.error('Failed to fetch journals from database:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Initial load
  fetchJournals()

  return {
    journals,
    isLoading,
    filteredJournals,
    stats,
    fetchJournals,
    getJournalById,
    submitDraft,
    reviewJournal,
    approveJournal,
    postJournal,
    rejectJournal,
    reverseJournal,
    createManualJournal
  }
})
