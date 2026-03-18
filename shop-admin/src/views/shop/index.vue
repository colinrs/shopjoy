<template>
  <div class="shop">
    <el-card>
      <template #header>
        <span>店铺设置</span>
      </template>
      
      <el-form :model="shopForm" label-width="120px">
        <el-form-item label="店铺名称">
          <el-input v-model="shopForm.name" />
        </el-form-item>
        <el-form-item label="店铺描述">
          <el-input v-model="shopForm.description" type="textarea" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="shopForm.contact_phone" />
        </el-form-item>
        <el-form-item label="联系邮箱">
          <el-input v-model="shopForm.contact_email" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave">保存</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const shopForm = ref({
  name: '',
  description: '',
  contact_phone: '',
  contact_email: ''
})

const handleSave = () => {
  ElMessage.success('保存成功')
}
</script>
