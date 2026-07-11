<template>
  <el-dialog
    :model-value="visible"
    title="Crop Image"
    width="640px"
    :close-on-click-modal="false"
    @update:model-value="onVisibleChange"
  >
    <div class="ic-wrapper">
      <img
        v-if="src"
        ref="imgRef"
        :src="src"
        alt="crop"
      >
    </div>
    <template #footer>
      <el-button @click="onCancel">
        Cancel
      </el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="onConfirm"
      >
        Confirm
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, watch, nextTick } from 'vue'
import Cropper from 'cropperjs'

// cropperjs 2.x bundles styles inline (no separate CSS file). Older 1.x had
// `cropperjs/dist/cropper.css` but we use 2.x so no import is needed.

interface Props {
  visible: boolean
  src: string | null
  aspectRatio?: number
}

const props = withDefaults(defineProps<Props>(), {
  aspectRatio: 1
})

const emit = defineEmits<{
  'update:visible': [v: boolean]
  confirm: [blob: Blob]
  cancel: []
}>()

const imgRef = ref<HTMLImageElement | null>(null)
const loading = ref(false)
let cropper: Cropper | null = null

function destroyCropper() {
  if (cropper) {
    cropper.destroy()
    cropper = null
  }
}

async function initCropper() {
  await nextTick()
  if (!imgRef.value) return
  destroyCropper()
  // cropperjs's bundled types are extremely narrow; cast to any so we can pass
  // the full option set the library actually supports at runtime.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  cropper = new (Cropper as any)(imgRef.value, {
    aspectRatio: props.aspectRatio,
    viewMode: 1,
    autoCropArea: 0.9,
    background: false,
    movable: true,
    zoomable: true,
    rotatable: true,
    scalable: false
  })
}

watch(
  () => props.visible,
  (v) => {
    if (v) {
      initCropper()
    } else {
      destroyCropper()
    }
  }
)

watch(
  () => props.src,
  () => {
    if (props.visible) initCropper()
  }
)

onMounted(() => {
  if (props.visible) initCropper()
})

onBeforeUnmount(() => {
  destroyCropper()
})

function onVisibleChange(v: boolean) {
  emit('update:visible', v)
}

function onCancel() {
  emit('update:visible', false)
  emit('cancel')
}

async function onConfirm() {
  if (!cropper) return
  loading.value = true
  try {
    // cropperjs .getCroppedCanvas() exists at runtime; types are incomplete.
    const canvas: HTMLCanvasElement = (cropper as any).getCroppedCanvas({
      maxWidth: 2048,
      maxHeight: 2048,
      fillColor: '#fff',
      imageSmoothingEnabled: true,
      imageSmoothingQuality: 'high'
    })
    if (!canvas) throw new Error('Cropper returned no canvas')
    const blob: Blob = await new Promise((resolve, reject) => {
      canvas.toBlob(
        (b: Blob | null) => {
          if (b) resolve(b)
          else reject(new Error('toBlob returned null'))
        },
        'image/png',
        0.95
      )
    })
    emit('confirm', blob)
  } catch (err: any) {
    // eslint-disable-next-line no-console
    console.error('Crop failed', err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.ic-wrapper {
  width: 100%;
  height: 360px;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;

  img {
    max-width: 100%;
    max-height: 100%;
  }
}
</style>
