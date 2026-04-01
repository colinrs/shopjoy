<template>
  <el-table
    ref="tableRef"
    :data="data"
    v-loading="loading"
    stripe
    @selection-change="handleSelectionChange"
    @select="handleSelect"
    @select-all="handleSelectAll"
    @cell-click="handleCellClick"
  >
    <el-table-column type="selection" width="50" :selectable="checkSelectable" />
    <slot />
  </el-table>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { ElTable } from 'element-plus'

export interface TableRow {
  [key: string]: any
}

const props = withDefaults(defineProps<{
  data: TableRow[]
  loading?: boolean
  rowKey?: string
}>(), {
  loading: false,
  rowKey: 'id'
})

const emit = defineEmits<{
  'selection-change': [selection: any[]]
}>()

const tableRef = ref<InstanceType<typeof ElTable>>()
const lastSelectedIndex = ref<number | null>(null)
const lastSelectedId = ref<number | string | null>(null)

// Check if row is selectable (always true for selection column)
const checkSelectable = () => true

// Get index by row id
const getIndexById = (id: number | string): number => {
  return props.data.findIndex(row => row[props.rowKey] === id)
}

// Check if a row is selected by trying to toggle and checking selection
const selectedRows = ref<Set<number | string>>(new Set())

// Check if a row is selected
const isRowSelected = (row: TableRow): boolean => {
  return selectedRows.value.has(row[props.rowKey])
}

// Handle cell click - detect shift/ctrl/cmd clicks on selection column
const handleCellClick = (row: TableRow, column: any, event: MouseEvent) => {
  // Only handle clicks on the selection column (type === 'selection')
  if (column.type !== 'selection') {
    return
  }

  const clickedIndex = getIndexById(row[props.rowKey])

  if (event.shiftKey && lastSelectedIndex.value !== null) {
    // Shift + click: select range
    handleShiftSelection(clickedIndex)
  } else if (event.ctrlKey || event.metaKey) {
    // Ctrl/Cmd + click: toggle single row
    handleCtrlSelection(row, clickedIndex)
  }
  // else: Normal click on checkbox - Element Plus handles it naturally

  // Update last selected after interaction
  lastSelectedIndex.value = clickedIndex
  lastSelectedId.value = row[props.rowKey]
}

// Handle shift selection - select range from last selected to current
const handleShiftSelection = (clickedIndex: number) => {
  if (lastSelectedIndex.value === null) {
    return
  }

  const start = Math.min(lastSelectedIndex.value, clickedIndex)
  const end = Math.max(lastSelectedIndex.value, clickedIndex)
  const rangeRows = props.data.slice(start, end + 1)

  // Select all rows in range (Element Plus will handle preserving other selections)
  rangeRows.forEach(row => {
    tableRef.value?.toggleRowSelection(row, true)
    selectedRows.value.add(row[props.rowKey])
  })
}

// Handle ctrl/cmd selection - toggle single row without affecting others
const handleCtrlSelection = (row: TableRow, clickedIndex: number) => {
  const isCurrentlySelected = isRowSelected(row)

  // Toggle the clicked row
  tableRef.value?.toggleRowSelection(row, !isCurrentlySelected)

  // Update selectedRows set
  if (!isCurrentlySelected) {
    selectedRows.value.add(row[props.rowKey])
  } else {
    selectedRows.value.delete(row[props.rowKey])
  }

  // If we're selecting (not deselecting), update last selected
  if (!isCurrentlySelected) {
    lastSelectedIndex.value = clickedIndex
    lastSelectedId.value = row[props.rowKey]
  }
}

// Handle selection change from Element Plus
const handleSelectionChange = (selection: TableRow[]) => {
  selectedRows.value = new Set(selection.map(row => row[props.rowKey]))
  emit('selection-change', selection)
}

// Handle select event
const handleSelect = (_selection: TableRow[], _row: TableRow) => {
  // Could be used for additional tracking if needed
}

// Handle select-all event
const handleSelectAll = (selection: TableRow[]) => {
  // Reset tracking on select all
  lastSelectedIndex.value = null
  lastSelectedId.value = null
  selectedRows.value = new Set(selection.map(r => r[props.rowKey]))
}

// Clear all selections
const clearSelection = () => {
  tableRef.value?.clearSelection()
  selectedRows.value.clear()
  lastSelectedIndex.value = null
  lastSelectedId.value = null
}

// Toggle selection for a single row
const toggleRowSelection = (row: TableRow, selected?: boolean) => {
  tableRef.value?.toggleRowSelection(row, selected)
}

// Toggle selection for multiple rows
const toggleRowSelectionForRows = (rows: TableRow[], selected: boolean) => {
  rows.forEach(_row => {
    tableRef.value?.toggleRowSelection(_row, selected)
  })
}

// Reset last selected when data changes significantly
watch(() => props.data, (newData, oldData) => {
  // If data was replaced entirely, reset tracking
  if (oldData && newData && oldData.length !== newData.length) {
    lastSelectedIndex.value = null
    lastSelectedId.value = null
  }
}, { deep: true })

// Expose methods for parent components
defineExpose({
  clearSelection,
  toggleRowSelection,
  toggleRowSelectionForRows,
  lastSelectedIndex,
  lastSelectedId
})
</script>

<style scoped>
/* Table styles can be customized here */
</style>
