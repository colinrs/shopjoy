import type { Product, ProductMarket } from '@/api/product'
import type { Warehouse } from '@/api/inventory'
import type { Market } from '@/api/market'

import type { CategoryTree } from '@/api/category'
import type { Brand } from '@/api/brand'

// Product form data for basic info tab
export interface ProductFormData {
  name: string
  description: string
  price: string
  currency: string
  cost_price: string
  stock: number
  status: 'draft' | 'on_sale' | 'off_sale' | 'deleted'
  category_id: string
  sku: string
  brand: string
  tags: string[]
  images: string[]
  is_matrix_product: boolean
  hs_code: string
  coo: string
  weight: string
  weight_unit: string
  length: string
  width: string
  height: string
  dangerous_goods: string[]
}

// Push to market form
export interface PushToMarketFormData {
  selectedMarkets: string[]
  price: number
}

// Variant form
export interface VariantFormData {
  id: string
  code: string
  price: number
  currency: string
  stock: number
  safety_stock: number
  pre_sale_enabled: boolean
  attributes: Record<string, string>
}

// Adjust stock form
export interface AdjustStockFormData {
  warehouse_id: string
  quantity: number
  remark: string
}

// Pricing table row data
export interface PricingRowData extends ProductMarket {
  price_value: number
  compare_at_price_value: number
}

// Supported language
export interface SupportedLanguage {
  code: string
  name: string
}

// Product Detail Props for parent component
export interface ProductDetailProps {
  productId: string
}

// Product Detail Emits
export interface ProductDetailEmits {
  (e: 'save-success'): void
  (e: 'status-change'): void
}

// Tab component props and emits
export interface ProductInfoTabProps {
  product: Product | null
  productForm: ProductFormData
  loading: boolean
  isDirty: boolean
  saveLoading: boolean
  categories: CategoryTree[]
  brands: Brand[]
  categoriesLoading: boolean
  brandsLoading: boolean
}

export interface ProductInfoTabEmits {
  (e: 'update:productForm', value: ProductFormData): void
  (e: 'show-add-image'): void
  (e: 'save'): void
}

// Product Markets Tab
export interface ProductMarketsTabProps {
  productId: string
  productMarkets: ProductMarket[]
  markets: Market[]
  loading: boolean
}

export interface ProductMarketsTabEmits {
  (e: 'update:productMarkets', value: ProductMarket[]): void
  (e: 'refresh'): void
  (e: 'show-push-to-market'): void
}

// Product Variants Tab
export interface ProductVariantsTabProps {
  productId: string
  defaultPrice: string
  defaultCurrency: string
  loading: boolean
}

export interface ProductVariantsTabEmits {
  (e: 'variants-change'): void
  (e: 'edit-variant', value: VariantFormData): void
}

// Product Pricing Tab
export interface ProductPricingTabProps {
  productId: string
  productMarkets: ProductMarket[]
  loading: boolean
}

export interface ProductPricingTabEmits {
  (e: 'refresh'): void
}

// Product Localization Tab
export interface ProductLocalizationTabProps {
  productId: string
  productName: string
  productDescription: string
  loading: boolean
}

export interface ProductLocalizationTabEmits {
  (e: 'localizations-change'): void
}

// Product Inventory Tab
export interface ProductInventoryTabProps {
  productId: string
  sku: string
  loading: boolean
}

export interface ProductInventoryTabEmits {
  (e: 'inventory-change'): void
  (e: 'go-to-variants'): void
}

// Product Reviews Tab
export interface ProductReviewsTabProps {
  productId: string
  loading: boolean
}

// Dialog props and emits
export interface PushToMarketDialogProps {
  visible: boolean
  productId: string
  productMarkets: ProductMarket[]
  markets: Market[]
  productPrice: string
  loading: boolean
}

export interface PushToMarketDialogEmits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

export interface AddImageDialogProps {
  visible: boolean
}

export interface AddImageDialogEmits {
  (e: 'update:visible', value: boolean): void
  (e: 'add', url: string): void
}

export interface AdjustStockDialogProps {
  visible: boolean
  sku: string
  warehouses: Warehouse[]
  loading: boolean
}

export interface AdjustStockDialogEmits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

export interface VariantDialogProps {
  visible: boolean
  productId: string
  isEdit: boolean
  variant: VariantFormData
  loading: boolean
}

export interface VariantDialogEmits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}
