<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Left Side - Image -->
      <div class="login-image">
        <div class="image-overlay">
          <h2>欢迎回来</h2>
          <p>发现精彩商品，享受购物乐趣</p>
        </div>
      </div>

      <!-- Right Side - Form -->
      <div class="login-form-wrapper">
        <div class="form-header">
          <h1>登录</h1>
          <p>还没有账号？ <router-link to="/register">立即注册</router-link></p>
        </div>

        <form class="login-form" @submit.prevent="handleLogin">
          <div class="form-group">
            <label>邮箱 / 手机号</label>
            <input 
              v-model="form.email"
              type="text"
              placeholder="请输入邮箱或手机号"
              required
            />
          </div>

          <div class="form-group">
            <label>密码</label>
            <div class="password-input">
              <input 
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="请输入密码"
                required
              />
              <button type="button" class="toggle-btn" @click="showPassword = !showPassword">
                <EyeIcon v-if="!showPassword" class="icon" />
                <EyeSlashIcon v-else class="icon" />
              </button>
            </div>
          </div>

          <div class="form-options">
            <label class="remember-me">
              <input v-model="form.remember" type="checkbox" />
              <span>记住我</span>
            </label>
            <router-link to="/forgot-password" class="forgot-link">忘记密码？</router-link>
          </div>

          <button type="submit" class="btn-login" :disabled="loading">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </form>

        <div class="divider">
          <span>或使用以下方式登录</span>
        </div>

        <div class="social-login">
          <button class="social-btn wechat">
            <svg viewBox="0 0 24 24" fill="currentColor"><path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 01.213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 00.167-.054l1.903-1.114a.864.864 0 01.717-.098 10.16 10.16 0 002.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 01-1.162 1.178A1.17 1.17 0 014.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 01-1.162 1.178 1.17 1.17 0 01-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 01.598.082l1.584.926a.272.272 0 00.14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 01-.023-.156.49.49 0 01.201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-7.062-6.122zM14.51 13.88a.94.94 0 01.939-.944c.519 0 .94.423.94.944a.94.94 0 01-.94.943.94.94 0 01-.94-.943zm4.963 0a.94.94 0 01.94-.944c.519 0 .939.423.939.944a.94.94 0 01-.94.943.94.94 0 01-.939-.943z"/></svg>
            微信登录
          </button>
          <button class="social-btn phone">
            <DevicePhoneMobileIcon class="icon" />
            手机号登录
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { EyeIcon, EyeSlashIcon, DevicePhoneMobileIcon } from '@heroicons/vue/24/outline'

const router = useRouter()
const loading = ref(false)
const showPassword = ref(false)

const form = reactive({
  email: '',
  password: '',
  remember: false
})

const handleLogin = async () => {
  loading.value = true
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 1500))
  loading.value = false
  router.push('/')
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #ECFDF5 0%, #D1FAE5 100%);
  padding: 24px;
}

.login-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  width: 100%;
  max-width: 1000px;
  min-height: 600px;
  background: white;
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

/* Login Image */
.login-image {
  position: relative;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-overlay {
  text-align: center;
  color: white;
  padding: 40px;
}

.image-overlay h2 {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 16px 0;
}

.image-overlay p {
  font-size: 18px;
  opacity: 0.9;
  margin: 0;
}

/* Login Form */
.login-form-wrapper {
  padding: 48px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-header {
  text-align: center;
  margin-bottom: 32px;
}

.form-header h1 {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 8px 0;
}

.form-header p {
  font-size: 15px;
  color: #6B7280;
  margin: 0;
}

.form-header a {
  color: #059669;
  font-weight: 600;
  text-decoration: none;
}

.form-header a:hover {
  text-decoration: underline;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input {
  padding: 14px 16px;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  font-size: 15px;
  outline: none;
  transition: all 0.2s;
}

.form-group input:focus {
  border-color: #059669;
  box-shadow: 0 0 0 3px rgba(5, 150, 105, 0.1);
}

.password-input {
  position: relative;
}

.password-input input {
  width: 100%;
  padding-right: 48px;
}

.toggle-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
}

.toggle-btn .icon {
  width: 20px;
  height: 20px;
  color: #9CA3AF;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.remember-me {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #4B5563;
  cursor: pointer;
}

.remember-me input {
  width: 18px;
  height: 18px;
  accent-color: #059669;
}

.forgot-link {
  font-size: 14px;
  color: #059669;
  text-decoration: none;
}

.forgot-link:hover {
  text-decoration: underline;
}

.btn-login {
  padding: 16px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-login:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px -10px rgba(5, 150, 105, 0.5);
}

.btn-login:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.divider {
  display: flex;
  align-items: center;
  margin: 24px 0;
  color: #9CA3AF;
  font-size: 14px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #E5E7EB;
}

.divider span {
  padding: 0 16px;
}

.social-login {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.social-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  background: #F3F4F6;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.social-btn:hover {
  background: #E5E7EB;
}

.social-btn svg,
.social-btn .icon {
  width: 20px;
  height: 20px;
}

.social-btn.wechat {
  color: #22C55E;
}

.social-btn.phone {
  color: #3B82F6;
}

/* Responsive */
@media (max-width: 768px) {
  .login-container {
    grid-template-columns: 1fr;
  }

  .login-image {
    display: none;
  }

  .login-form-wrapper {
    padding: 32px 24px;
  }
}
</style>
