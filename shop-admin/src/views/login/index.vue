<template>
  <div class="login-container">
    <div class="login-wrapper">
      <!-- Left Side - Branding -->
      <div class="login-branding">
        <div class="branding-content">
          <div class="logo">
            <el-icon size="48" color="#fff"><ShoppingBag /></el-icon>
            <h1>ShopJoy</h1>
          </div>
          <p class="tagline">专业的电商管理平台</p>
          <div class="features">
            <div class="feature-item">
              <el-icon><Check /></el-icon>
              <span>多店铺统一管理</span>
            </div>
            <div class="feature-item">
              <el-icon><Check /></el-icon>
              <span>智能数据分析</span>
            </div>
            <div class="feature-item">
              <el-icon><Check /></el-icon>
              <span>全渠道订单同步</span>
            </div>
          </div>
        </div>
        <div class="branding-footer">
          <p>© 2024 ShopJoy. All rights reserved.</p>
        </div>
      </div>

      <!-- Right Side - Login Form -->
      <div class="login-form-section">
        <div class="form-container">
          <h2 class="form-title">欢迎回来</h2>
          <p class="form-subtitle">请登录您的管理员账户</p>

          <el-form 
            :model="loginForm" 
            :rules="loginRules" 
            ref="loginFormRef"
            class="login-form"
          >
            <el-form-item prop="account">
              <el-input 
                v-model="loginForm.account" 
                placeholder="请输入邮箱/用户名/手机号"
                size="large"
                class="custom-input"
              >
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item prop="password">
              <el-input 
                v-model="loginForm.password" 
                type="password" 
                placeholder="请输入密码"
                size="large"
                class="custom-input"
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <div class="form-options">
              <el-checkbox v-model="rememberMe">记住我</el-checkbox>
              <a href="#" class="forgot-link">忘记密码?</a>
            </div>
            <el-form-item>
              <el-button 
                type="primary" 
                @click="handleLogin" 
                :loading="loading" 
                size="large"
                class="login-btn"
              >
                登 录
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { adminLogin } from '@/api/admin-user'
import { ShoppingBag, Message, Lock, Check } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const rememberMe = ref(false)
const loginFormRef = ref()

const loginForm = reactive({
  account: '',
  password: ''
})

const loginRules = {
  account: [
    { required: true, message: '请输入邮箱/用户名/手机号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        const res = await adminLogin({
          account: loginForm.account,
          password: loginForm.password
        })
        userStore.setToken(res.access_token)
        userStore.userInfo = res.user
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error: any) {
        ElMessage.error(error.message || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  background: linear-gradient(135deg, #059669 0%, #10B981 50%, #059669 100%);
}

.login-wrapper {
  display: flex;
  width: 100%;
  max-width: 1200px;
  margin: auto;
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  min-height: 600px;
}

/* Left Side - Branding */
.login-branding {
  flex: 1;
  background: linear-gradient(135deg, #059669 0%, #064E3B 100%);
  padding: 60px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  color: white;
}

.branding-content {
  flex: 1;
}

.logo {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.logo h1 {
  font-size: 36px;
  font-weight: 700;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
}

.tagline {
  font-size: 18px;
  opacity: 0.9;
  margin-bottom: 60px;
  font-weight: 300;
}

.features {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
}

.feature-item .el-icon {
  background: rgba(255, 255, 255, 0.2);
  padding: 8px;
  border-radius: 50%;
}

.branding-footer {
  font-size: 14px;
  opacity: 0.7;
}

/* Right Side - Login Form */
.login-form-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  background: #fafafa;
}

.form-container {
  width: 100%;
  max-width: 400px;
}

.form-title {
  font-size: 32px;
  font-weight: 700;
  color: #064E3B;
  margin: 0 0 8px 0;
  font-family: 'Fira Sans', sans-serif;
}

.form-subtitle {
  font-size: 16px;
  color: #6B7280;
  margin: 0 0 32px 0;
}

.login-form :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.login-form :deep(.el-input__inner) {
  height: 48px;
  font-size: 15px;
}

.login-form :deep(.el-input__prefix) {
  color: #9CA3AF;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 16px 0 24px;
}

.forgot-link {
  color: #059669;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.forgot-link:hover {
  color: #047857;
}

.terms-link {
  color: #059669;
  text-decoration: none;
}

.terms-link:hover {
  text-decoration: underline;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px -10px rgba(5, 150, 105, 0.5);
}

/* Responsive */
@media (max-width: 768px) {
  .login-container {
    padding: 20px;
    background: #fff;
  }
  
  .login-wrapper {
    flex-direction: column;
    border-radius: 0;
    box-shadow: none;
  }
  
  .login-branding {
    padding: 40px 30px;
    min-height: auto;
  }
  
  .features {
    display: none;
  }
  
  .login-form-section {
    padding: 40px 30px;
  }
}
</style>
