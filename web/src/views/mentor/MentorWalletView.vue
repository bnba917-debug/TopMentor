<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showConfirmDialog, showFailToast, showSuccessToast } from 'vant'
import { NavBar, Button, CellGroup, Cell, Empty, Loading } from 'vant'
import type { WalletSummary } from '@/api/types'
import { fetchMentorWallet, mentorWithdraw, getErrorMessage } from '@/api/client'

const router = useRouter()
const loading = ref(true)
const withdrawing = ref(false)
const wallet = ref<WalletSummary | null>(null)

function txLabel(type: string) {
  return type === 'EARN' ? '课时收入' : type === 'WITHDRAW' ? '提现' : type
}

async function load() {
  loading.value = true
  try {
    wallet.value = await fetchMentorWallet()
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

async function onWithdraw() {
  if (!wallet.value || wallet.value.balance < 1) {
    showFailToast('余额不足')
    return
  }
  const cents = Math.floor(wallet.value.balance * 100)
  try {
    await showConfirmDialog({
      title: '确认提现',
      message: `将提现 ¥${(cents / 100).toFixed(2)} 至微信零钱（开发模式为模拟到账）`,
    })
  } catch {
    return
  }

  withdrawing.value = true
  try {
    const result = await mentorWithdraw(cents)
    showSuccessToast(result.mock_paid ? '模拟提现成功' : '提现申请已提交')
    await load()
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    withdrawing.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="我的钱包" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" class="center-loading" vertical>加载中...</Loading>

    <template v-else-if="wallet">
      <div class="balance-card">
        <p class="label">可提现余额（元）</p>
        <p class="amount">{{ wallet.balance.toFixed(2) }}</p>
        <Button round type="primary" size="small" :loading="withdrawing" @click="onWithdraw">
          全部提现
        </Button>
      </div>

      <CellGroup inset title="最近流水">
        <Empty v-if="wallet.transactions.length === 0" description="暂无流水" />
        <Cell
          v-for="tx in wallet.transactions"
          :key="tx.id"
          :title="txLabel(tx.type)"
          :label="tx.remark || tx.created_at"
          :value="`${tx.amount >= 0 ? '+' : ''}${tx.amount.toFixed(2)}`"
        />
      </CellGroup>
    </template>
  </div>
</template>

<style scoped>
.center-loading {
  margin-top: 80px;
  text-align: center;
}

.balance-card {
  margin: 16px;
  padding: 24px;
  border-radius: 12px;
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  color: #fff;
  text-align: center;
}

.label {
  margin: 0;
  opacity: 0.85;
  font-size: 14px;
}

.amount {
  margin: 8px 0 16px;
  font-size: 36px;
  font-weight: 700;
}
</style>
