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
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 50%, #6366F1 100%);
}

.login-wrapper {
  display: flex;
  width: 100%;
  max-width: 1200px;
  margin: auto;
  background: #fff;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(99, 102, 241, 0.35);
  min-height: 600px;
}

/* Left Side - Branding */
.login-branding {
  flex: 1;
  background: linear-gradient(135deg, #4338CA 0%, #6366F1 100%);
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
  letter-spacing: -0.5px;
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
  border-radius: 10px;
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
  background: #F5F3FF;
}

.form-container {
  width: 100%;
  max-width: 400px;
}

.form-title {
  font-size: 32px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 8px 0;
  font-family: 'Fira Sans', sans-serif;
}

.form-subtitle {
  font-size: 16px;
  color: #6B7280;
  margin: 0 0 32px 0;
}

.login-form :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
}

.login-form :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
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
  color: #6366F1;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.forgot-link:hover {
  color: #4338CA;
}

.terms-link {
  color: #6366F1;
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
  border-radius: 12px;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  border: none;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px -10px rgba(99, 102, 241, 0.5);
}

.login-btn:active {
  transform: translateY(0);
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
    background: #fff;
  }
}
</style>
