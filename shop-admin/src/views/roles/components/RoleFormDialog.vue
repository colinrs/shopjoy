<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? $t('roles.editRole') : $t('roles.createRole')"
    width="500px"
    destroy-on-close
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
      <el-form-item :label="$t('roles.roleName')" prop="name">
        <el-input
          v-model="formData.name"
          :placeholder="$t('roles.enterRoleName')"
          :disabled="isEdit && !!props.role?.is_system"
        />
      </el-form-item>
      <el-form-item :label="$t('roles.roleCode')" prop="code" v-if="!isEdit">
        <el-input
          v-model="formData.code"
          :placeholder="$t('roles.enterRoleCode')"
        />
      </el-form-item>
      <el-form-item :label="$t('roles.roleCode')" prop="code" v-else>
        <el-input
          v-model="formData.code"
          disabled
        />
      </el-form-item>
      <el-form-item :label="$t('roles.description')" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          :placeholder="$t('roles.enterDescription')"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleCancel">{{ $t('roles.cancel') }}</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
        {{ $t('roles.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { createRole, updateRole, type Role, type CreateRoleParams } from '@/api/role'

const props = defineProps<{
  visible: boolean
  role: Role | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', val: boolean): void
  (e: 'success'): void
}>()

const { t } = useI18n()

const formRef = ref()
const submitLoading = ref(false)

const isEdit = computed(() => !!props.role)

const formData = reactive<CreateRoleParams>({
  name: '',
  code: '',
  description: '',
  permission_ids: []
})

const formRules = computed(() => ({
  name: [
    { required: true, message: t('roles.pleaseEnterRoleName'), trigger: 'blur' },
    { min: 2, max: 50, message: t('roles.roleNameLength'), trigger: 'blur' }
  ],
  code: [
    { required: true, message: t('roles.pleaseEnterRoleCode'), trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: t('roles.roleCodePattern'), trigger: 'blur' },
    { min: 2, max: 50, message: t('roles.roleCodeLength'), trigger: 'blur' }
  ]
}))

watch(() => props.visible, (val) => {
  if (val) {
    if (props.role) {
      formData.name = props.role.name
      formData.code = props.role.code
      formData.description = props.role.description || ''
    } else {
      formData.name = ''
      formData.code = ''
      formData.description = ''
    }
  }
})

const handleCancel = () => {
  emit('update:visible', false)
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (isEdit.value && props.role) {
          await updateRole(props.role.id, {
            name: formData.name,
            description: formData.description
          })
          ElMessage.success(t('roles.updateSuccess'))
        } else {
          await createRole(formData)
          ElMessage.success(t('roles.createSuccess'))
        }
        emit('update:visible', false)
        emit('success')
      } catch (e) {
        console.error('Failed to save role:', e)
        ElMessage.error(t('roles.saveFailed'))
      } finally {
        submitLoading.value = false
      }
    }
  })
}
</script>

<style scoped>
:deep(.el-dialog) {
  border-radius: 16px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: #1E1B4B;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}
</style>
