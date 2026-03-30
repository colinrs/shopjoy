<template>
  <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
    <el-form-item :label="$t('adminUsers.emailLabel')" prop="email">
      <el-input v-model="formData.email" :placeholder="$t('adminUsers.pleaseEnterEmail')" />
    </el-form-item>
    <el-form-item :label="$t('adminUsers.mobileLabel')" prop="mobile">
      <el-input v-model="formData.mobile" :placeholder="$t('adminUsers.pleaseEnterEmail')" />
    </el-form-item>
    <el-form-item :label="$t('adminUsers.realNameLabel')" prop="real_name">
      <el-input v-model="formData.real_name" :placeholder="$t('adminUsers.pleaseEnterRealName')" />
    </el-form-item>
    <el-form-item :label="$t('adminUsers.initialPassword')" prop="password" v-if="!isEdit">
      <el-input v-model="formData.password" type="password" :placeholder="$t('adminUsers.initialPasswordPlaceholder')" show-password />
    </el-form-item>
    <el-form-item :label="$t('adminUsers.userType')" prop="type" v-if="showTypeSelect">
      <el-select v-model="formData.type" :placeholder="$t('adminUsers.pleaseSelectUserType')" style="width: 100%">
        <el-option :label="$t('adminUsers.merchantAdmin')" :value="2" />
        <el-option :label="$t('adminUsers.merchantSubAccount')" :value="3" />
      </el-select>
    </el-form-item>
    <el-form-item :label="$t('adminUsers.tenantId')" prop="tenant_id" v-if="showTenantId">
      <el-input-number v-model="formData.tenant_id" :min="0" :placeholder="$t('adminUsers.tenantIdPlaceholder')" style="width: 100%" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import type { CreateAdminUserParams, UpdateAdminUserParams, AdminUser } from '@/api/admin-user'

const { t } = useI18n()

const props = defineProps<{
  modelValue: CreateAdminUserParams | UpdateAdminUserParams
  isEdit?: boolean
  showTypeSelect?: boolean
  showTenantId?: boolean
  adminUser?: AdminUser | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: CreateAdminUserParams | UpdateAdminUserParams]
}>()

const formRef = ref<FormInstance>()

const formData = reactive({
  email: '',
  mobile: '',
  real_name: '',
  password: '',
  type: 2 as number,
  tenant_id: undefined as number | undefined
})

const rules = computed<FormRules>(() => ({
  email: [
    { required: true, message: t('adminUsers.pleaseEnterEmail'), trigger: 'blur' },
    { type: 'email', message: t('adminUsers.pleaseEnterValidEmail'), trigger: 'blur' }
  ],
  real_name: [{ required: true, message: t('adminUsers.pleaseEnterRealName'), trigger: 'blur' }],
  password: props.isEdit ? [] : [
    { required: true, message: t('adminUsers.pleaseEnterPassword'), trigger: 'blur' },
    { min: 6, message: t('adminUsers.passwordMinLength'), trigger: 'blur' }
  ],
  type: [{ required: true, message: t('adminUsers.pleaseSelectUserType'), trigger: 'change' }]
}))

// Sync from parent
watch(() => props.adminUser, (user) => {
  if (user && props.isEdit) {
    formData.email = user.email || ''
    formData.mobile = user.mobile || ''
    formData.real_name = user.real_name || ''
    formData.type = user.type || 2
  }
}, { immediate: true })

// Sync to parent
watch(formData, (val) => {
  if (props.isEdit) {
    emit('update:modelValue', {
      email: val.email,
      mobile: val.mobile,
      real_name: val.real_name
    })
  } else {
    emit('update:modelValue', {
      email: val.email,
      mobile: val.mobile || undefined,
      real_name: val.real_name,
      password: val.password,
      type: val.type,
      tenant_id: val.tenant_id
    })
  }
}, { deep: true })

const validate = async () => {
  if (!formRef.value) return false
  return formRef.value.validate()
}

defineExpose({
  validate
})
</script>
