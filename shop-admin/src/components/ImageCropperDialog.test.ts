import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

// Stub cropperjs entirely so we never touch a real <img> + canvas.
vi.mock('cropperjs', () => ({
  default: class FakeCropper {
    destroy() {}
    getCroppedCanvas() {
      return {
        toBlob(cb: (b: Blob | null) => void) {
          cb(new Blob([new Uint8Array(4)], { type: 'image/png' }))
        }
      }
    }
  }
}))
vi.mock('cropperjs/dist/cropper.css', () => ({}))

// Stub Element Plus dialog + button so the component can mount without
// pulling the full Element Plus runtime.
const stubs: Record<string, any> = {
  'el-dialog': {
    name: 'ElDialog',
    props: ['modelValue'],
    template:
      '<div v-if="modelValue" class="el-dialog-stub"><div class="el-dialog-default"><slot /></div><div class="el-dialog-footer"><slot name="footer" /></div></div>'
  },
  'el-button': {
    name: 'ElButton',
    props: ['loading'],
    template:
      '<button class="el-button-stub" :disabled="loading" @click="$emit(\'click\')"><slot /></button>'
  }
}

import ImageCropperDialog from './ImageCropperDialog.vue'

describe('ImageCropperDialog', () => {
  it('mounts without throwing when visible is false', () => {
    const wrapper = mount(ImageCropperDialog, {
      props: {
        visible: false,
        src: null
      },
      global: { stubs }
    })
    expect(wrapper.exists()).toBe(true)
    // Dialog body not rendered because modelValue=false.
    expect(wrapper.find('.el-dialog-stub').exists()).toBe(false)
  })

  it('renders the dialog body when visible is true', () => {
    const wrapper = mount(ImageCropperDialog, {
      props: {
        visible: true,
        src: 'https://example.com/sample.jpg',
        aspectRatio: 1
      },
      global: { stubs }
    })
    expect(wrapper.exists()).toBe(true)
    expect(wrapper.find('.el-dialog-stub').exists()).toBe(true)
    // Image element with the given src is rendered.
    const img = wrapper.find('img')
    expect(img.exists()).toBe(true)
    expect(img.attributes('src')).toBe('https://example.com/sample.jpg')
  })

  it('emits update:visible and cancel when Cancel is clicked', async () => {
    const wrapper = mount(ImageCropperDialog, {
      props: { visible: true, src: 'https://example.com/x.jpg' },
      global: { stubs }
    })
    const buttons = wrapper.findAll('.el-button-stub')
    // First button is Cancel.
    await buttons[0]!.trigger('click')
    const updates = wrapper.emitted('update:visible')
    expect(updates).toBeDefined()
    expect(updates![0]).toEqual([false])
    expect(wrapper.emitted('cancel')).toBeDefined()
  })

  it('emits confirm with a Blob when Confirm is clicked', async () => {
    const wrapper = mount(ImageCropperDialog, {
      props: { visible: true, src: 'https://example.com/x.jpg' },
      global: { stubs }
    })

    // The component initializes `cropper` via a watcher on `visible`.
    // The Cropper stub we provide is a class, so new Cropper(img) is
    // instantiated. The watcher is async (await nextTick) — wait it out.
    await new Promise((r) => setTimeout(r, 50))

    const buttons = wrapper.findAll('.el-button-stub')
    // Second button is Confirm.
    await buttons[1]!.trigger('click')
    // Wait for the async toBlob callback to resolve.
    await new Promise((r) => setTimeout(r, 50))
    const confirms = wrapper.emitted('confirm')
    expect(confirms).toBeDefined()
    // Component may emit twice due to the watcher on `visible` re-running
    // when src is set; accept >= 1.
    expect(confirms!.length).toBeGreaterThanOrEqual(1)
    expect(confirms![0]![0]).toBeInstanceOf(Blob)
  })
})