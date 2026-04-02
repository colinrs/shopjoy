<template>
  <el-dialog
    :model-value="visible"
    :title="$t('roles.assignPermissions') + ' - ' + (role?.name || '')"
    width="700px"
    destroy-on-close
    @update:model-value="$emit('update:visible', $event)"
  >
    <div v-loading="loading" class="permission-selector">
      <el-input
        v-model="searchQuery"
        :placeholder="$t('roles.searchPermission')"
        class="search-input"
        clearable
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <div class="permission-tree-container">
        <el-tree
          ref="treeRef"
          :data="permissionTree"
          :props="treeProps"
          node-key="id"
          :default-expand-all="true"
          :expand-on-click-node="false"
          :filter-node-method="filterNode"
          show-checkbox
          check-strictly
          @check="handleCheck"
        >
          <template #default="{ data }">
            <span class="tree-node">
              <span class="node-label">{{ data.name }}</span>
              <el-tag v-if="data.type === 0" type="primary" size="small" class="type-tag">
                {{ $t('roles.menu') }}
              </el-tag>
              <el-tag v-else-if="data.type === 1" type="success" size="small" class="type-tag">
                {{ $t('roles.button') }}
              </el-tag>
              <el-tag v-else type="warning" size="small" class="type-tag">
                {{ $t('roles.api') }}
              </el-tag>
            </span>
          </template>
        </el-tree>
      </div>

      <div class="selected-info">
        <span>{{ $t('roles.selectedCount') }}: {{ selectedCount }}</span>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleCancel">{{ $t('roles.cancel') }}</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
        {{ $t('roles.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import {
  getPermissionList,
  getRoleDetail,
  updateRolePermissions,
  type Role,
  type Permission
} from '@/api/role'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<{
  visible: boolean
  role: Role | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', val: boolean): void
  (e: 'success'): void
}>()

const { t } = useI18n()
const { handleError } = useErrorHandler()

const treeRef = ref()
const loading = ref(false)
const submitLoading = ref(false)
const searchQuery = ref('')
const allPermissions = ref<Permission[]>([])
const selectedPermissionIds = ref<number[]>([])

const treeProps = {
  children: 'children',
  label: 'name'
}

// Build permission tree from flat list
const permissionTree = computed(() => {
  const map = new Map<number, Permission & { children?: Permission[] }>()
  const roots: (Permission & { children?: Permission[] })[] = []

  // First pass: create map
  allPermissions.value.forEach(p => {
    map.set(p.id, { ...p, children: [] })
  })

  // Second pass: build tree
  allPermissions.value.forEach(p => {
    const node = map.get(p.id)!
    if (p.parent_id === 0 || !map.has(p.parent_id)) {
      roots.push(node)
    } else {
      const parent = map.get(p.parent_id)
      if (parent) {
        parent.children = parent.children || []
        parent.children.push(node)
      }
    }
  })

  return roots
})

const selectedCount = computed(() => selectedPermissionIds.value.length)

const filterNode = (query: string, data: Permission) => {
  if (!query) return true
  return data.name.toLowerCase().includes(query.toLowerCase()) ||
         data.code.toLowerCase().includes(query.toLowerCase())
}

const handleCheck = (_data: Permission, { checkedKeys }: { checkedKeys: number[] }) => {
  selectedPermissionIds.value = checkedKeys
}

watch(() => props.visible, async (val) => {
  if (val) {
    await loadPermissions()
    await loadRolePermissions()
  }
})

watch(searchQuery, (val) => {
  treeRef.value?.filter(val)
})

const loadPermissions = async () => {
  loading.value = true
  try {
    const res = await getPermissionList()
    allPermissions.value = res.list || []
  } catch (e) {
    handleError(e, t('roles.loadPermissionsFailed'))
  } finally {
    loading.value = false
  }
}

const loadRolePermissions = async () => {
  if (!props.role) return

  try {
    const res = await getRoleDetail(props.role.id)
    selectedPermissionIds.value = (res.permissions || []).map(p => p.id)

    // Set checked keys in tree
    nextTick(() => {
      if (treeRef.value) {
        treeRef.value.setCheckedKeys(selectedPermissionIds.value)
      }
    })
  } catch (e) {
    handleError(e)
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}

const handleSubmit = async () => {
  if (!props.role) return

  submitLoading.value = true
  try {
    // Get all checked nodes including half-checked parents
    const checkedNodes = treeRef.value?.getCheckedNodes(false, true) || []
    const permissionIds = checkedNodes.map((n: Permission) => n.id)

    await updateRolePermissions(props.role.id, { permission_ids: permissionIds })
    ElMessage.success(t('roles.permissionsUpdatedSuccess'))
    emit('update:visible', false)
    emit('success')
  } catch (e) {
    handleError(e, t('roles.updatePermissionsFailed'))
  } finally {
    submitLoading.value = false
  }
}
</script>

<style scoped>
.permission-selector {
  max-height: 60vh;
  overflow-y: auto;
}

.search-input {
  margin-bottom: 16px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.permission-tree-container {
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 12px;
  max-height: 400px;
  overflow-y: auto;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.node-label {
  flex: 1;
}

.type-tag {
  margin-left: 8px;
}

.selected-info {
  margin-top: 12px;
  padding: 8px 12px;
  background: #F5F3FF;
  border-radius: 8px;
  font-size: 13px;
  color: #6366F1;
}

:deep(.el-tree) {
  background: transparent;
}

:deep(.el-tree-node__content) {
  height: 36px;
}

:deep(.el-tree-node__content:hover) {
  background-color: #F5F3FF;
}

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
