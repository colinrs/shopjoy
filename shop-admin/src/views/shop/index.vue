<template>
  <div class="shop-settings">
    <el-row :gutter="20">
      <!-- Left Column - Basic Info -->
      <el-col :xs="24" :lg="16">
        <el-card class="settings-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Shop /></el-icon>
                店铺基本信息
              </span>
            </div>
          </template>
          
          <el-form :model="shopForm" label-width="120px" :rules="shopRules" ref="shopFormRef">
            <el-form-item label="店铺Logo">
              <el-upload
                class="logo-uploader"
                action="#"
                :show-file-list="false"
                :auto-upload="false"
                :on-change="handleLogoChange"
              >
                <img v-if="shopForm.logo" :src="shopForm.logo" class="logo-preview" />
                <div v-else class="logo-placeholder">
                  <el-icon size="32"><Plus /></el-icon>
                  <span>上传Logo</span>
                </div>
              </el-upload>
              <span class="upload-tip">建议尺寸 200x200px，支持 JPG、PNG 格式</span>
            </el-form-item>
            
            <el-form-item label="店铺名称" prop="name">
              <el-input v-model="shopForm.name" placeholder="请输入店铺名称" maxlength="50" show-word-limit />
            </el-form-item>
            
            <el-form-item label="店铺简介" prop="description">
              <el-input 
                v-model="shopForm.description" 
                type="textarea" 
                :rows="4"
                placeholder="请输入店铺简介，这将展示在店铺首页"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>
            
            <el-form-item label="主营类目">
              <el-select v-model="shopForm.category" placeholder="请选择主营类目" style="width: 100%">
                <el-option label="数码电子" value="electronics" />
                <el-option label="服装配饰" value="clothing" />
                <el-option label="家居生活" value="home" />
                <el-option label="运动户外" value="sports" />
                <el-option label="美妆护肤" value="beauty" />
                <el-option label="食品饮料" value="food" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="客服电话" prop="contact_phone">
              <el-input v-model="shopForm.contact_phone" placeholder="请输入客服电话">
                <template #prefix>
                  <el-icon><Phone /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item label="客服邮箱" prop="contact_email">
              <el-input v-model="shopForm.contact_email" placeholder="请输入客服邮箱">
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item label="店铺地址">
              <el-input v-model="shopForm.address" placeholder="请输入店铺地址">
                <template #prefix>
                  <el-icon><Location /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- SEO Settings -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Search /></el-icon>
                SEO 设置
              </span>
            </div>
          </template>
          
          <el-form :model="seoForm" label-width="120px">
            <el-form-item label="页面标题">
              <el-input v-model="seoForm.title" placeholder="SEO标题，显示在浏览器标签页" maxlength="60" show-word-limit />
            </el-form-item>
            
            <el-form-item label="关键词">
              <el-select
                v-model="seoForm.keywords"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="请输入关键词"
                style="width: 100%"
              />
            </el-form-item>
            
            <el-form-item label="页面描述">
              <el-input 
                v-model="seoForm.description" 
                type="textarea" 
                :rows="3"
                placeholder="SEO描述，搜索引擎结果中显示的摘要"
                maxlength="200"
                show-word-limit
              />
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- Right Column - Status & Settings -->
      <el-col :xs="24" :lg="8">
        <!-- Shop Status -->
        <el-card class="settings-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><InfoFilled /></el-icon>
                店铺状态
              </span>
            </div>
          </template>
          
          <div class="status-item">
            <span class="status-label">店铺状态</span>
            <el-switch
              v-model="shopStatus.isOpen"
              active-text="营业中"
              inactive-text="已打烊"
            />
          </div>
          
          <div class="status-item">
            <span class="status-label">新订单通知</span>
            <el-switch v-model="shopStatus.newOrderNotify" />
          </div>
          
          <div class="status-item">
            <span class="status-label">低库存提醒</span>
            <el-switch v-model="shopStatus.lowStockAlert" />
          </div>
          
          <div class="status-item">
            <span class="status-label">自动接单</span>
            <el-switch v-model="shopStatus.autoAccept" />
          </div>
        </el-card>

        <!-- Notification Settings -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Bell /></el-icon>
                通知设置
              </span>
            </div>
          </template>
          
          <el-form :model="notifyForm" label-position="top">
            <el-form-item label="新订单通知邮箱">
              <el-input v-model="notifyForm.orderEmail" placeholder="多个邮箱用逗号分隔" />
            </el-form-item>
            
            <el-form-item label="库存预警阈值">
              <el-input-number v-model="notifyForm.stockThreshold" :min="1" :max="100" style="width: 100%" />
            </el-form-item>
            
            <el-form-item label="通知方式">
              <el-checkbox-group v-model="notifyForm.methods">
                <el-checkbox label="email">邮件</el-checkbox>
                <el-checkbox label="sms">短信</el-checkbox>
                <el-checkbox label="app">App推送</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- Quick Actions -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Operation /></el-icon>
                快捷操作
              </span>
            </div>
          </template>
          
          <div class="quick-actions">
            <el-button plain style="width: 100%; margin-bottom: 12px">
              <el-icon><Document /></el-icon>
              查看店铺主页
            </el-button>
            <el-button plain style="width: 100%; margin-bottom: 12px">
              <el-icon><Share /></el-icon>
              分享店铺链接
            </el-button>
            <el-button type="danger" plain style="width: 100%">
              <el-icon><Delete /></el-icon>
              清空缓存
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Save Button -->
    <div class="save-bar">
      <el-button type="primary" size="large" @click="handleSave" :loading="saving">
        <el-icon><Check /></el-icon>
        保存设置
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Shop, Plus, Phone, Message, Location, Search,
  InfoFilled, Bell, Operation, Document, Share, Delete, Check
} from '@element-plus/icons-vue'
import { uploadImage } from '@/api/upload'

const saving = ref(false)
const shopFormRef = ref()

const shopForm = reactive({
  name: 'ShopJoy 官方店铺',
  description: '专注于高品质数码产品，为用户提供优质的购物体验。我们承诺正品保证，七天无理由退换。',
  logo: '',
  category: 'electronics',
  contact_phone: '400-888-8888',
  contact_email: 'support@shopjoy.com',
  address: '深圳市南山区科技园'
})

const seoForm = reactive({
  title: 'ShopJoy - 您的品质生活选择',
  keywords: ['数码', '电子产品', '品质生活'],
  description: 'ShopJoy 官方店铺，提供高品质数码产品，正品保证，售后无忧。'
})

const shopStatus = reactive({
  isOpen: true,
  newOrderNotify: true,
  lowStockAlert: true,
  autoAccept: false
})

const notifyForm = reactive({
  orderEmail: 'admin@shopjoy.com',
  stockThreshold: 10,
  methods: ['email']
})

const shopRules = {
  name: [
    { required: true, message: '请输入店铺名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '最多 500 个字符', trigger: 'blur' }
  ],
  contact_phone: [
    { pattern: /^1[3-9]\d{9}$|^400-\d{3}-\d{4}$/, message: '请输入正确的电话号码', trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const handleLogoChange = async (file: any) => {
  try {
    const response = await uploadImage(file.raw, 'banner')
    shopForm.logo = response.url
  } catch (error) {
    console.error('Upload failed:', error)
    ElMessage.error('Logo上传失败')
  }
}

const handleSave = async () => {
  if (!shopFormRef.value) return
  
  await shopFormRef.value.validate((valid: boolean) => {
    if (valid) {
      saving.value = true
      setTimeout(() => {
        saving.value = false
        ElMessage.success('店铺设置保存成功')
      }, 1000)
    }
  })
}
</script>

<style scoped>
.shop-settings {
  padding: 0;
}

.settings-card {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Logo Uploader */
.logo-uploader {
  border: 2px dashed #D1D5DB;
  border-radius: 12px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s;
  width: 120px;
  height: 120px;
}

.logo-uploader:hover {
  border-color: #059669;
}

.logo-preview {
  width: 120px;
  height: 120px;
  object-fit: cover;
  display: block;
}

.logo-placeholder {
  width: 120px;
  height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #9CA3AF;
  gap: 8px;
}

.logo-placeholder span {
  font-size: 12px;
}

.upload-tip {
  margin-left: 16px;
  font-size: 12px;
  color: #9CA3AF;
}

/* Status Items */
.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #E5E7EB;
}

.status-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.status-label {
  font-size: 14px;
  color: #374151;
}

/* Quick Actions */
.quick-actions {
  display: flex;
  flex-direction: column;
}

.quick-actions .el-button {
  justify-content: flex-start;
  gap: 8px;
}

/* Save Bar */
.save-bar {
  position: fixed;
  bottom: 0;
  left: 220px;
  right: 0;
  padding: 16px 24px;
  background: #fff;
  border-top: 1px solid #E5E7EB;
  display: flex;
  justify-content: center;
  z-index: 100;
}

.save-bar .el-button {
  min-width: 160px;
}

/* Responsive */
@media (max-width: 768px) {
  .save-bar {
    left: 0;
  }
}
</style>
