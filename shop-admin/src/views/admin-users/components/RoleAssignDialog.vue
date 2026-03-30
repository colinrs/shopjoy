<template>
  <el-dialog
    v-model="visible"
    :title="$t('adminUsers.assignRoleTitle')"
    width="500px"
    destroy-on-close
    @close="$emit('update:visible', false)"
  >
    <div class="role-list" v-loading="loading">
      <el-checkbox-group v-model="selectedRoleIds">
        <div v-for="role in roles" :key="role.id" class="role-item">
          <el-checkbox :label="role.id">
            <div class="role-info">
              <span class="role-name">{{ role.name }}</span>
              <span class="role-code">{{ role.code }}</span>
            </div>
          </el-checkbox>
        </div>
      </el-checkbox-group>
      <el-empty v-if="!loading && roles.length === 0" :description="$t('adminUsers.noRolesAvailable')" />
    </div>
    <template #footer>
      <el-button @click="visible = false">{{ $t('adminUsers.cancel') }}</el-button>
      <el-button type="primary" @click="handleAssign" :loading="assigning">{{ $t('adminUsers.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getAvailableRoles, assignRoles, type AdminRole, type AdminUserDetail } from '@/api/admin-user'

const { t } = useI18n()

const props = defineProps<{
  visible: boolean
  adminUser: AdminUserDetail | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'assigned': []
}>()

const loading = ref(false)
const assigning = ref(false)
const roles = ref<AdminRole[]>([])
const selectedRoleIds = ref<number[]>([])

const visible = ref(props.visible)

watch(() => props.visible, (val) => {
  visible.value = val
  if (val) {
    loadRoles()
    if (props.adminUser?.roles) {
      selectedRoleIds.value = props.adminUser.roles.map(r => r.id)
    }
  }
})

watch(visible, (val) => {
  emit('update:visible', val)
})

const loadRoles = async () => {
  loading.value = true
  try {
    const res = await getAvailableRoles()
    roles.value = res?.list || []
  } catch (error) {
    console.error('Failed to load roles:', error)
    ElMessage.error(t('adminUsers.failedToLoadRoles'))
    roles.value = []
  } finally {
    loading.value = false
  }
}

const handleAssign = async () => {
  if (!props.adminUser) return

  assigning.value = true
  try {
    await assignRoles(props.adminUser.id, selectedRoleIds.value)
    ElMessage.success(t('adminUsers.roleAssignedSuccess'))
    emit('assigned')
    visible.value = false
  } catch (error) {
    console.error('Failed to assign roles:', error)
    ElMessage.error(t('adminUsers.roleAssignedFailed'))
  } finally {
    assigning.value = false
  }
}
</script>

<style scoped>
.role-list {
  max-height: 400px;
  overflow-y: auto;
}

.role-item {
  padding: 12px 0;
  border-bottom: 1px solid #F3F4F6;
}

.role-item:last-child {
  border-bottom: none;
}

.role-info {
  display: flex;
  flex-direction: column;
}

.role-name {
  font-weight: 500;
  color: #1E1B4B;
}

.role-code {
  font-size: 12px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}
</style>
