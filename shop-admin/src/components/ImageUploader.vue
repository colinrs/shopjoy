<template>
  <div class="image-uploader">
    <!-- Single mode: avatar / logo / single image -->
    <template v-if="!multiple">
      <div class="iu-single">
        <div
          v-if="currentSingle"
          class="iu-preview"
          :class="{ 'iu-preview--error': currentSingle.status === 'error' }"
        >
          <img
            :src="currentSingle.preview || currentSingle.url"
            :alt="category"
          >
          <div
            v-if="currentSingle.status !== 'done' && currentSingle.status !== 'pending'"
            class="iu-overlay"
          >
            <el-progress
              v-if="currentSingle.status === 'uploading' || currentSingle.status === 'signing' || currentSingle.status === 'confirming'"
              :percentage="currentSingle.progress"
              :stroke-width="6"
              type="circle"
            />
            <el-icon
              v-else-if="currentSingle.status === 'error'"
              size="24"
              color="#f56c6c"
            >
              <CircleClose />
            </el-icon>
          </div>
          <div class="iu-actions">
            <el-button
              size="small"
              type="danger"
              plain
              @click="onRemoveSingle"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
            <el-button
              v-if="currentSingle.status === 'error'"
              size="small"
              type="primary"
              plain
              @click="onRetry(currentSingle)"
            >
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </div>
        <el-upload
          v-else
          class="iu-dropzone"
          action="#"
          :show-file-list="false"
          :auto-upload="false"
          :accept="accept"
          :on-change="onSingleFileChange"
        >
          <el-icon class="iu-icon">
            <Plus />
          </el-icon>
          <div class="iu-hint">
            {{ placeholder || $t('upload.clickOrDrag') }}
          </div>
        </el-upload>
      </div>
    </template>

    <!-- Multi mode: gallery of previews with drag reorder + add -->
    <template v-else>
      <draggable
        v-model="multiItems"
        item-key="uid"
        :animation="180"
        class="iu-grid"
        ghost-class="iu-grid__ghost"
        @end="onReorder"
      >
        <template #item="{ element }">
          <div
            :key="element.uid"
            class="iu-card"
            :class="{ 'iu-card--error': element.status === 'error' }"
          >
            <img
              :src="element.preview || element.url"
              :alt="element.uid"
            >
            <div
              v-if="element.status !== 'done' && element.status !== 'pending'"
              class="iu-overlay"
            >
              <el-progress
                v-if="element.status === 'uploading' || element.status === 'signing' || element.status === 'confirming'"
                :percentage="element.progress"
                :stroke-width="6"
                type="circle"
              />
              <el-icon
                v-else-if="element.status === 'error'"
                size="22"
                color="#f56c6c"
              >
                <CircleClose />
              </el-icon>
            </div>
            <div class="iu-card__actions">
              <el-button
                size="small"
                type="danger"
                plain
                @click="onRemoveItem(element)"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
              <el-button
                v-if="element.status === 'error'"
                size="small"
                type="primary"
                plain
                @click="onRetry(element)"
              >
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
            <div
              v-if="element.status === 'error'"
              class="iu-card__error-tip"
            >
              {{ element.error || $t('upload.uploadFailed') }}
            </div>
          </div>
        </template>
        <template #footer>
          <el-upload
            v-if="multiItems.length < (max ?? Infinity)"
            class="iu-card iu-card--add"
            action="#"
            :show-file-list="false"
            :multiple="true"
            :auto-upload="false"
            :accept="accept"
            :on-change="onMultiFileChange"
          >
            <el-icon class="iu-icon">
              <Plus />
            </el-icon>
            <div class="iu-hint">
              {{ placeholder || $t('upload.addImage') }}
            </div>
          </el-upload>
        </template>
      </draggable>
    </template>

    <ImageCropperDialog
      v-if="crop"
      v-model:visible="cropperVisible"
      :src="pendingCropSrc"
      :aspect-ratio="cropRatio"
      @confirm="onCropConfirm"
      @cancel="cropperVisible = false"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete, Refresh, CircleClose } from '@element-plus/icons-vue'
import draggable from 'vuedraggable'
import { useImageUpload, type UploadItem } from '@/composables/useImageUpload'
import ImageCropperDialog from './ImageCropperDialog.vue'

interface Props {
  /**
   * Bound value. `string` for single mode, `string[]` for multi mode.
   * Consumed as `v-model:value="..."`.
   */
  value: string | string[]
  category: 'product' | 'banner' | 'avatar' | string
  multiple?: boolean
  max?: number
  maxSize?: number
  accept?: string
  compress?: boolean
  crop?: boolean
  cropRatio?: number
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  multiple: false,
  max: 9,
  maxSize: 5 * 1024 * 1024,
  accept: 'image/*',
  compress: true,
  crop: false,
  cropRatio: 1,
  placeholder: ''
})

const emit = defineEmits<{
  'update:value': [value: string | string[]]
  change: [value: string | string[]]
}>()

const { items, uploadOne, retry, remove: removeItem } = useImageUpload({
  category: props.category,
  maxSize: props.maxSize,
  compress: props.compress
})

// --- Single mode bookkeeping ---
// We keep the items[] array as the source of truth and project it onto
// value (string in single mode).
const currentSingle = computed<UploadItem | undefined>(() => {
  if (props.multiple) return undefined
  // First item, or undefined if none.
  return items.value[0]
})

const multiItems = computed({
  get: () => items.value,
  set: () => {
    /* draggable writes via v-model; we sync via onReorder */
  }
})

function syncSingle() {
  if (props.multiple) return
  // In single mode the most recently completed upload should win over
  // any stale "existing" item seeded by the external-value watcher.
  // Otherwise a user picking a new file in the brand/avatar edit dialog
  // would never have their selection reach the parent form.
  const latestDone = [...items.value].reverse().find((it) => it.status === 'done')
  const done = latestDone?.url ?? ''
  if (done !== props.value) {
    emit('update:value', done)
  }
}

function syncMulti() {
  if (!props.multiple) return
  const urls = items.value.filter((it) => it.status === 'done' && it.url).map((it) => it.url!) as string[]
  if (JSON.stringify(urls) !== JSON.stringify(props.value)) {
    emit('update:value', urls)
  }
}

watch(() => items.value.map((it) => `${it.uid}:${it.status}:${it.url ?? ''}`).join('|'), () => {
  if (props.multiple) syncMulti()
  else syncSingle()
})

// External value → items (e.g. when form loads or after save re-fetch).
// Rebuilds items from props.value so pre-existing images render in the
// gallery (multi mode) or preview (single mode). Preserves any items that
// are still uploading so a concurrent upload is not lost when the parent
// re-pushes the value (e.g. after Save → loadProduct).
watch(
  () => props.value,
  (val) => {
    const stamp = Date.now()
    if (props.multiple) {
      const arr = Array.isArray(val) ? val : []
      const existing = items.value
        .filter((it) => it.status === 'done')
        .map((it) => it.url)
      if (JSON.stringify(existing) !== JSON.stringify(arr)) {
        const externalItems: UploadItem[] = arr.map((url, idx) => ({
          uid: `existing-${idx}-${stamp}`,
          status: 'done',
          progress: 100,
          url,
          preview: url
        }))
        const inFlight = items.value.filter((it) => it.status !== 'done')
        items.value = [...inFlight, ...externalItems]
      }
    } else {
      // Single mode: val is a single URL. Seed one synthetic done item
      // from it so the existing logo / avatar renders. We do nothing
      // when val is empty — clearing items on empty input would erase
      // an item the user has already uploaded and is just about to save.
      const url = typeof val === 'string' ? val : ''
      if (url === '') return
      const currentDoneUrl = items.value.find((it) => it.status === 'done')?.url
      if (currentDoneUrl !== url) {
        const externalItem: UploadItem = {
          uid: `existing-0-${stamp}`,
          status: 'done',
          progress: 100,
          url,
          preview: url
        }
        const inFlight = items.value.filter((it) => it.status !== 'done')
        items.value = [externalItem, ...inFlight]
      }
    }
  },
  { immediate: true }
)

// --- File intake ---

async function onSingleFileChange(file: { raw?: File }) {
  if (!file.raw) return
  if (props.crop) {
    pendingCropSrc.value = URL.createObjectURL(file.raw)
    pendingCropFile.value = file.raw
    cropperVisible.value = true
    return
  }
  try {
    await uploadOne(file.raw)
    syncSingle()
    emit('change', props.value)
  } catch {
    // Error already set on item.
  }
}

async function onMultiFileChange(file: { raw?: File }) {
  if (!file.raw) return
  if (items.value.length >= (props.max ?? Infinity)) {
    ElMessage.warning(`Max ${props.max} images`)
    return
  }
  try {
    await uploadOne(file.raw)
    syncMulti()
    emit('change', props.value)
  } catch {
    // Error already set on item.
  }
}

// --- Crop flow ---

const cropperVisible = ref(false)
const pendingCropSrc = ref<string | null>(null)
const pendingCropFile = ref<File | null>(null)

async function onCropConfirm(croppedBlob: Blob) {
  if (!pendingCropFile.value) return
  const original = pendingCropFile.value
  const croppedFile = new File([croppedBlob], original.name, { type: original.type })
  cropperVisible.value = false
  pendingCropSrc.value = null
  try {
    await uploadOne(croppedFile)
    syncSingle()
    emit('change', props.value)
  } catch {
    // ignore
  } finally {
    pendingCropFile.value = null
  }
}

// --- Reorder ---

function onReorder() {
  if (!props.multiple) return
  // vuedraggable mutates items.value in place; just re-sync.
  syncMulti()
}

// --- Remove / retry ---

async function onRemoveSingle() {
  const it = currentSingle.value
  if (it) await removeItem(it)
  syncSingle()
  emit('change', props.value)
}

async function onRemoveItem(it: UploadItem) {
  await removeItem(it)
  syncMulti()
  emit('change', props.value)
}

async function onRetry(it: UploadItem) {
  try {
    await retry(it)
    if (props.multiple) syncMulti()
    else syncSingle()
    emit('change', props.value)
  } catch {
    // ignore
  }
}
</script>

<style scoped lang="scss">
.image-uploader {
  width: 100%;
}

/* Single mode */
.iu-single {
  display: inline-block;
}
.iu-preview {
  position: relative;
  width: 120px;
  height: 120px;
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  overflow: hidden;
  background: #fafafa;
  display: inline-flex;
  align-items: center;
  justify-content: center;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  &--error {
    border-color: var(--el-color-danger);
  }
}
.iu-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}
.iu-actions {
  position: absolute;
  right: 4px;
  bottom: 4px;
  display: flex;
  gap: 4px;
}
.iu-dropzone {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 120px;
  height: 120px;
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  flex-direction: column;
  gap: 4px;
  background: #fafafa;
}
.iu-icon {
  font-size: 24px;
  color: var(--el-text-color-secondary);
}
.iu-hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

/* Multi mode */
.iu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 8px;

  &__ghost {
    opacity: 0.4;
  }
}
.iu-card {
  position: relative;
  width: 100%;
  padding-top: 100%; /* 1:1 */
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  overflow: hidden;
  cursor: grab;
  background: #fafafa;

  img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  &--error {
    border-color: var(--el-color-danger);
  }
  &--add {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    flex-direction: column;
    gap: 4px;
    border-style: dashed;
  }
}
.iu-card__actions {
  position: absolute;
  right: 4px;
  bottom: 4px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.15s ease;
}
.iu-card:hover .iu-card__actions {
  opacity: 1;
}
.iu-card__error-tip {
  position: absolute;
  left: 4px;
  right: 4px;
  bottom: 26px;
  font-size: 11px;
  color: #fff;
  background: rgba(245, 108, 108, 0.85);
  padding: 2px 4px;
  border-radius: 2px;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
