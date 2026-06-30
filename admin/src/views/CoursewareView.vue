<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Courseware } from '@/api/types'
import {
  createCourseware,
  deleteCourseware,
  fetchCourseware,
  getErrorMessage,
  updateCourseware,
} from '@/api/client'

const loading = ref(true)
const list = ref<Courseware[]>([])
const dialogVisible = ref(false)
const editing = ref<Courseware | null>(null)

const form = reactive({
  title: '',
  cover_url: '',
  content_url: '',
  sort_order: 0,
  is_active: true,
})

function resetForm() {
  form.title = ''
  form.cover_url = ''
  form.content_url = ''
  form.sort_order = 0
  form.is_active = true
}

async function load() {
  loading.value = true
  try {
    list.value = await fetchCourseware()
  } catch (e) {
    ElMessage.error(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  resetForm()
  dialogVisible.value = true
}

function openEdit(row: Courseware) {
  editing.value = row
  form.title = row.title
  form.cover_url = row.cover_url || ''
  form.content_url = row.content_url
  form.sort_order = row.sort_order
  form.is_active = row.is_active
  dialogVisible.value = true
}

async function onSubmit() {
  if (!form.title || !form.content_url) {
    ElMessage.warning('请填写标题和内容链接')
    return
  }
  try {
    if (editing.value) {
      await updateCourseware(editing.value.id, { ...form })
      ElMessage.success('已更新')
    } else {
      await createCourseware({ ...form })
      ElMessage.success('已创建')
    }
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(getErrorMessage(e))
  }
}

async function onDelete(row: Courseware) {
  try {
    await ElMessageBox.confirm(`确认删除课件「${row.title}」？`, '删除确认', { type: 'warning' })
    await deleteCourseware(row.id)
    ElMessage.success('已删除')
    await load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(getErrorMessage(e))
  }
}

async function toggleActive(row: Courseware) {
  try {
    await updateCourseware(row.id, { is_active: !row.is_active })
    ElMessage.success(row.is_active ? '已下架' : '已上架')
    await load()
  } catch (e) {
    ElMessage.error(getErrorMessage(e))
  }
}

onMounted(load)
</script>

<template>
  <div class="page-card">
    <div class="head">
      <h2>课件管理</h2>
      <el-button type="primary" @click="openCreate">新增课件</el-button>
    </div>

    <el-table v-loading="loading" :data="list" stripe>
      <el-table-column prop="title" label="标题" min-width="160" />
      <el-table-column prop="content_url" label="内容链接" min-width="220" show-overflow-tooltip />
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">{{ row.is_active ? '上架' : '下架' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link @click="toggleActive(row)">{{ row.is_active ? '下架' : '上架' }}</el-button>
          <el-button link type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editing ? '编辑课件' : '新增课件'" width="520px">
      <el-form label-width="90px">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="封面链接">
          <el-input v-model="form.cover_url" placeholder="https://..." />
        </el-form-item>
        <el-form-item label="内容链接" required>
          <el-input v-model="form.content_url" placeholder="PDF / H5 绘本 URL" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="上架">
          <el-switch v-model="form.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="onSubmit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

h2 {
  margin: 0;
}
</style>
