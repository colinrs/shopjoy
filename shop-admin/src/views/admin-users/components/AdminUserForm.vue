<template>
  <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
    <el-form-item label="邮箱" prop="email">
      <el-input v-model="formData.email" placeholder="请输入邮箱" />
    </el-form-item>
    <el-form-item label="手机号" prop="mobile">
      <el-input v-model="formData.mobile" placeholder="请输入手机号" />
    </el-form-item>
    <el-form-item label="真实姓名" prop="real_name">
      <el-input v-model="formData.real_name" placeholder="请输入真实姓名" />
    </el-form-item>
    <el-form-item label="初始密码" prop="password" v-if="!isEdit">
      <el-input v-model="formData.password" type="password" placeholder="请输入初始密码" show-password />
    </el-form-item>
    <el-form-item label="用户类型" prop="type" v-if="showTypeSelect">
      <el-select v-model="formData.type" placeholder="请选择用户类型" style="width: 100%">
        <el-option label="商家管理员" :value="2" />
        <el-option label="商家子账号" :value="3" />
      </el-select>
    </el-form-item>
    <el-form-item label="租户ID" prop="tenant_id" v-if="showTenantId">
      <el-input-number v-model="formData.tenant_id" :min="0" placeholder="创建商家管理员时使用" style="width: 100%" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { CreateAdminUserParams, UpdateAdminUserParams, AdminUser } from '@/api/admin-user'

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
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  password: props.isEdit ? [] : [
    { required: true, message: '请输入初始密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择用户类型', trigger: 'change' }]
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