<template>
  <el-dialog
    :model-value="visible"
    :title="$t('products.addImageUrl')"
    width="400px"
    @update:model-value="emit('update:visible', $event)"
  >
    <el-input
      v-model="imageUrl"
      :placeholder="$t('products.enterImageUrl')"
    />
    <template #footer>
      <el-button @click="handleClose">
        {{ $t('common.cancel') }}
      </el-button>
      <el-button
        type="primary"
        @click="handleConfirm"
      >
        {{ $t('common.add') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { AddImageDialogProps, AddImageDialogEmits } from '../types'

const props = defineProps<AddImageDialogProps>()
const emit = defineEmits<AddImageDialogEmits>()

const imageUrl = ref('')

watch(() => props.visible, (newVal) => {
  if (newVal) {
    imageUrl.value = ''
  }
})

const handleClose = () => {
  emit('update:visible', false)
}

const handleConfirm = () => {
  if (imageUrl.value.trim()) {
    emit('add', imageUrl.value.trim())
    handleClose()
  }
}
</script>
