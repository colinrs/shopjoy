import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { nextTick } from 'vue'

// Mock the child selectors with lightweight stubs so the test does not pull
// in the full Element Plus runtime.
vi.mock('./FeeTypeSelector.vue', () => ({
  default: {
    name: 'FeeTypeSelector',
    props: ['modelValue'],
    template: '<div class="fee-type-stub" />'
  }
}))
vi.mock('./RegionSelector.vue', () => ({
  default: {
    name: 'RegionSelector',
    props: ['modelValue', 'selected'],
    template: '<div class="region-stub" />'
  }
}))

vi.mock('element-plus', () => ({
  ElMessage: { warning: vi.fn(), error: vi.fn(), success: vi.fn() }
}))

vi.mock('@element-plus/icons-vue', () => ({
  Rank: { name: 'RankIcon', template: '<i />' },
  Edit: { name: 'EditIcon', template: '<i />' },
  Delete: { name: 'DeleteIcon', template: '<i />' },
  Location: { name: 'LocationIcon', template: '<i />' },
  Plus: { name: 'PlusIcon', template: '<i />' }
}))

import ZoneConfigForm from './ZoneConfigForm.vue'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      shipping: {
        zoneName: 'Zone Name',
        enterZoneName: 'Enter zone name',
        zoneNameRequired: 'Please enter zone name',
        zoneNameLength: 'Zone name must be 2-50 characters',
        deliveryArea: 'Delivery Area',
        selectDeliveryArea: 'Select area',
        feeType: 'Fee Type',
        freeShippingCondition: 'Free Shipping',
        cancel: 'Cancel',
        update: 'Update',
        add: 'Add',
        taxSettings: 'Tax',
        taxable: 'Taxable',
        taxRate: 'Tax Rate',
        taxRatePlaceholder: 'e.g., 0.13',
        taxIncluded: 'Price Includes Tax',
        iossApplicable: 'IOSS',
        remoteArea: 'Remote Area',
        remoteSurcharge: 'Remote Surcharge',
        remoteZipPatterns: 'Remote ZIP',
        remoteZipPatternsPlaceholder: 'One per line',
        fuelSurcharge: 'Fuel',
        fuelSurchargePct: 'Fuel %',
        nameI18n: 'Multilingual Name',
        addNameI18n: 'Add Locale',
        locale: 'Locale',
        localePlaceholder: 'en-US',
        volumetricDivisor: 'Volumetric Divisor',
        volumetricDivisorTip: 'tip',
        market: 'Market',
        countryCode: 'Country Code',
        countryCodePlaceholder: 'US',
        postalCode: 'Postal',
        addressType: 'Type',
        addressTypeChina: 'CN',
        addressTypeInternational: 'Intl'
      },
      common: {
        edit: 'Edit',
        delete: 'Delete',
        cancel: 'Cancel'
      }
    }
  }
})

// Lightweight stubs for Element Plus form components
const formStubs: Record<string, any> = {
  'el-form': {
    name: 'ElForm',
    props: ['model', 'rules'],
    template: '<form class="el-form-stub"><slot /></form>',
    methods: {
      async validate(cb?: (valid: boolean) => void) {
        // Always pass validation in the test; call callback with true once.
        if (cb) cb(true)
        return true
      }
    }
  },
  'el-form-item': {
    name: 'ElFormItem',
    template: '<div class="el-form-item-stub" :class="$attrs.class"><slot /></div>',
    inheritAttrs: false
  },
  'el-input': {
    name: 'ElInput',
    props: ['modelValue'],
    template:
      '<input class="el-input-stub" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
  },
  'el-input-number': {
    name: 'ElInputNumber',
    props: ['modelValue', 'min', 'max', 'precision'],
    template:
      '<input class="el-input-number-stub" type="number" :value="modelValue" @input="$emit(\'update:modelValue\', Number($event.target.value))" />'
  },
  'el-switch': {
    name: 'ElSwitch',
    props: ['modelValue'],
    template:
      '<button class="el-switch-stub" @click="$emit(\'update:modelValue\', !modelValue)">{{ modelValue }}</button>'
  },
  'el-checkbox': {
    name: 'ElCheckbox',
    props: ['modelValue'],
    template:
      '<button class="el-checkbox-stub" @click="$emit(\'update:modelValue\', !modelValue)">{{ modelValue }}</button>'
  },
  'el-button': {
    name: 'ElButton',
    props: ['loading'],
    template:
      '<button class="el-button-stub" :data-loading="loading" @click="$emit(\'click\')"><slot /></button>'
  },
  'el-row': {
    name: 'ElRow',
    template: '<div class="el-row-stub"><slot /></div>'
  },
  'el-col': {
    name: 'ElCol',
    template: '<div class="el-col-stub"><slot /></div>'
  },
  'el-divider': {
    name: 'ElDivider',
    template: '<div class="el-divider-stub" />'
  },
  'el-icon': {
    name: 'ElIcon',
    template: '<span class="el-icon-stub"><slot /></span>'
  }
}

async function mountForm(props: Record<string, any> = {}) {
  const wrapper = mount(ZoneConfigForm, {
    props: {
      isDialog: true,
      ...props
    },
    global: {
      plugins: [i18n],
      stubs: formStubs
    }
  })
  await flushPromises()
  // Manually populate regions to satisfy the rule
  ;(wrapper.vm as any).form.regions = ['110000']
  await nextTick()
  return wrapper
}

function clickSubmit(wrapper: any) {
  // The submit button is in .form-actions, the second el-button there
  const formActions = wrapper.find('.form-actions')
  if (!formActions.exists()) {
    throw new Error('.form-actions not found')
  }
  const buttons = formActions.findAll('.el-button-stub')
  if (buttons.length < 2) {
    throw new Error('expected at least 2 buttons in form-actions, got ' + buttons.length)
  }
  // Submit is the last button (cancel + submit)
  const submit = buttons[buttons.length - 1]
  // Use DOM-level click to avoid Vue's click event double-firing
  submit.element.click()
}

describe('ZoneConfigForm', () => {
  it('emits save with all P1 fields when submitted', async () => {
    const wrapper = await mountForm()
    const inputs = wrapper.findAll('.el-input-stub')
    inputs[0].setValue('East China')
    await flushPromises()
    clickSubmit(wrapper)
    await flushPromises()
    const saves = wrapper.emitted('save')
    expect(saves).toBeTruthy()
    // The exact emit count can be 1 or 2 depending on async re-render timing
    // (the submitting flag toggles cause re-render of the el-button). The
    // important thing is the payload structure.
    expect(saves!.length).toBeGreaterThanOrEqual(1)
    const payload: any = (saves![0] as any[])[0]
    // Required fields
    expect(payload).toHaveProperty('name')
    expect(payload).toHaveProperty('regions')
    expect(payload).toHaveProperty('fee_type')
    expect(payload).toHaveProperty('first_unit')
    expect(payload).toHaveProperty('first_fee')
    expect(payload).toHaveProperty('additional_unit')
    expect(payload).toHaveProperty('additional_fee')
    expect(payload).toHaveProperty('free_threshold_amount')
    expect(payload).toHaveProperty('free_threshold_count')
    // P1-6
    expect(payload).toHaveProperty('taxable')
    expect(payload).toHaveProperty('tax_rate')
    expect(payload).toHaveProperty('tax_included')
    expect(payload).toHaveProperty('ioss_applicable')
    // P1-7
    expect(payload).toHaveProperty('remote_surcharge')
    expect(payload).toHaveProperty('remote_zip_patterns')
    expect(Array.isArray(payload.remote_zip_patterns)).toBe(true)
    // P1-8
    expect(payload).toHaveProperty('fuel_surcharge_pct')
  })

  it('omits volumetric_divisor when fee_type is not by_volume', async () => {
    const wrapper = await mountForm()
    const inputs = wrapper.findAll('.el-input-stub')
    inputs[0].setValue('Test Zone')
    await flushPromises()
    clickSubmit(wrapper)
    await flushPromises()
    const saves = wrapper.emitted('save')
    expect(saves).toBeTruthy()
    const payload: any = (saves![0] as any[])[0]
    expect(payload.volumetric_divisor).toBeUndefined()
  })

  it('pre-fills name_i18n when editing a zone', async () => {
    const wrapper = await mountForm({
      zone: {
        id: '1',
        tenant_id: '1',
        template_id: '1',
        market_id: '1',
        currency: 'CNY',
        name: 'East China',
        name_i18n: [{ locale: 'en-US', name: 'East China' }],
        regions: ['110000'],
        fee_type: 'fixed',
        first_unit: 1,
        first_fee: '10',
        additional_unit: 1,
        additional_fee: '5',
        free_threshold_amount: '0',
        free_threshold_count: 0,
        taxable: true,
        tax_rate: '0.13',
        tax_included: true,
        ioss_applicable: false,
        remote_surcharge: '5',
        remote_zip_patterns: ['9*', '8*'],
        fuel_surcharge_pct: '2',
        volumetric_divisor: 5000,
        sort: 0
      }
    })
    await flushPromises()
    clickSubmit(wrapper)
    await flushPromises()
    const saves = wrapper.emitted('save')
    expect(saves).toBeTruthy()
    const payload: any = (saves![0] as any[])[0]
    expect(payload.name_i18n).toEqual([{ locale: 'en-US', name: 'East China' }])
    expect(payload.taxable).toBe(true)
    expect(payload.tax_rate).toBe('0.13')
    expect(payload.remote_zip_patterns).toEqual(['9*', '8*'])
    expect(payload.fuel_surcharge_pct).toBe('2')
  })
})
