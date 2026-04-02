<template>
  <div class="roles-page">
    <el-card
      class="filter-card"
      shadow="never"
    >
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('roles.searchPlaceholder')"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="filterStatus"
            :placeholder="$t('roles.filterStatus')"
            clearable
            class="filter-select"
            @change="handleSearch"
          >
            <el-option
              :label="$t('roles.allStatus')"
              :value="0"
            />
            <el-option
              :label="$t('roles.enabled')"
              :value="1"
            />
            <el-option
              :label="$t('roles.disabled')"
              :value="2"
            />
          </el-select>
        </div>
        <div class="filter-right">
          <el-button
            type="primary"
            @click="handleCreate"
          >
            <el-icon><Plus /></el-icon>{{ $t('roles.addRole') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card
      class="table-card"
      shadow="never"
    >
      <el-table
        v-loading="loading"
        :data="tableData"
        stripe
      >
        <el-table-column
          :label="$t('roles.roleName')"
          min-width="200"
        >
          <template #default="{ row }">
            <div class="role-cell">
              <div class="role-info">
                <span class="role-name">{{ row.name }}</span>
                <span class="role-code">{{ row.code }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('roles.description')"
          min-width="200"
        >
          <template #default="{ row }">
            <span class="description-text">{{ row.description || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('roles.isSystem')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-tag
              v-if="row.is_system"
              type="warning"
              size="small"
            >
              {{ $t('roles.system') }}
            </el-tag>
            <span
              v-else
              class="normal-text"
            >-</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('roles.status')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="2"
              :disabled="row.is_system"
              @change="(val: number) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column
          prop="created_at"
          :label="$t('roles.createdAt')"
          width="160"
        >
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('roles.actions')"
          width="200"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleEdit(row)"
            >
              {{ $t('roles.edit') }}
            </el-button>
            <el-button
              type="primary"
              link
              size="small"
              @click="handleAssignPermissions(row)"
            >
              {{ $t('roles.assignPermissions') }}
            </el-button>
            <el-dropdown
              trigger="click"
              @command="(cmd: string) => handleCommand(cmd, row)"
            >
              <el-button
                type="primary"
                link
                size="small"
              >
                {{ $t('roles.more') }}<el-icon class="el-icon--right">
                  <ArrowDown />
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="viewPermissions">
                    {{ $t('roles.viewPermissions') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                    command="delete"
                    style="color: #EF4444;"
                    :disabled="row.is_system"
                  >
                    {{ $t('roles.delete') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- Role Form Dialog (Create/Edit) -->
    <RoleFormDialog
      v-model:visible="formDialogVisible"
      :role="currentRole"
      @success="loadData"
    />

    <!-- Permission Selector Dialog -->
    <PermissionSelector
      v-model:visible="permissionDialogVisible"
      :role="currentRole"
      @success="loadData"
    />

    <!-- Role Permissions View Dialog -->
    <el-dialog
      v-model="viewPermissionsVisible"
      :title="$t('roles.viewPermissionsTitle')"
      width="600px"
      destroy-on-close
    >
      <div
        v-if="currentRolePermissions.length > 0"
        class="permissions-list"
      >
        <el-tag
          v-for="perm in currentRolePermissions"
          :key="perm.id"
          :type="getPermissionTagType(perm.type)"
          class="permission-tag"
        >
          {{ perm.name }}
        </el-tag>
      </div>
      <el-empty
        v-else
        :description="$t('roles.noPermissions')"
      />
      <template #footer>
        <el-button @click="viewPermissionsVisible = false">
          {{ $t('roles.close') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, ArrowDown } from '@element-plus/icons-vue'
import RoleFormDialog from './components/RoleFormDialog.vue'
import PermissionSelector from './components/PermissionSelector.vue'
import {
  getRoleList,
  deleteRole,
  updateRoleStatus,
  getRoleDetail,
  type Role,
  type Permission,
  type ListRolesParams
} from '@/api/role'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { t } = useI18n()
const { handleError } = useErrorHandler()

const loading = ref(false)
const searchQuery = ref('')
const filterStatus = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref<Role[]>([])

const formDialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const viewPermissionsVisible = ref(false)
const currentRole = ref<Role | null>(null)
const currentRolePermissions = ref<Permission[]>([])

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleCreate = () => {
  currentRole.value = null
  formDialogVisible.value = true
}

const handleEdit = (row: Role) => {
  currentRole.value = row
  formDialogVisible.value = true
}

const handleAssignPermissions = async (row: Role) => {
  currentRole.value = row
  permissionDialogVisible.value = true
}

const handleCommand = async (command: string, row: Role) => {
  if (command === 'delete') {
    await handleDelete(row)
  } else if (command === 'viewPermissions') {
    await handleViewPermissions(row)
  }
}

const handleViewPermissions = async (row: Role) => {
  try {
    const res = await getRoleDetail(row.id)
    currentRolePermissions.value = res.permissions || []
    currentRole.value = row
    viewPermissionsVisible.value = true
  } catch (error) {
    handleError(error, t('roles.loadPermissionsFailed'))
  }
}

const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm(
      t('roles.confirmDelete', { name: row.name }),
      t('roles.deleteConfirmTitle'),
      {
        confirmButtonText: t('roles.confirm'),
        cancelButtonText: t('roles.cancel'),
        type: 'warning'
      }
    )
    await deleteRole(row.id)
    ElMessage.success(t('roles.deleteSuccess'))
    loadData()
  } catch {
    // User cancelled or error
  }
}

const handleStatusChange = async (row: Role, val: number) => {
  try {
    await ElMessageBox.confirm(
      val === 1 ? t('roles.confirmEnable') : t('roles.confirmDisable'),
      t('roles.prompt'),
      {
        confirmButtonText: t('roles.confirm'),
        cancelButtonText: t('roles.cancel'),
        type: 'warning'
      }
    )
    await updateRoleStatus(row.id, { status: val })
    ElMessage.success(t('roles.statusUpdateSuccess'))
  } catch {
    // User cancelled or error - revert switch
    row.status = val === 1 ? 2 : 1
  }
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadData()
}

const getPermissionTagType = (type: number): string => {
  const types: Record<number, string> = {
    0: 'primary', // menu
    1: 'success', // button
    2: 'warning'  // api
  }
  return types[type] || 'info'
}

const loadData = async () => {
  loading.value = true
  try {
    const params: ListRolesParams = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (searchQuery.value) params.name = searchQuery.value
    if (filterStatus.value) params.status = filterStatus.value === 1 ? 1 : 2

    const res = await getRoleList(params)
    tableData.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    handleError(error, t('roles.loadFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.roles-page {
  padding: 0;
}

.filter-card {
  margin-bottom: 16px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 260px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

.table-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.role-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.role-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.role-name {
  font-weight: 600;
  color: #1E1B4B;
}

.role-code {
  font-size: 12px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

.description-text {
  font-size: 13px;
  color: #6B7280;
}

.normal-text {
  color: #9CA3AF;
}

.system-tag {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

/* Tags */
:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

/* Switch */
:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #10B981;
}

/* Dialog */
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

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

.permissions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.permission-tag {
  margin: 4px;
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }
}
</style>
