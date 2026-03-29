import request from '@/utils/request'

// ========== Theme Types ==========

export interface ThemeItem {
  id: number
  code: string
  name: string
  description: string
  preview_image: string
  thumbnail: string
  is_preset: boolean
  is_current: boolean
}

export interface ListThemesResponse {
  themes: ThemeItem[]
}

export interface ThemeConfigDTO {
  primary_color: string
  secondary_color: string
  font_family: string
  button_style: string
}

export interface CurrentThemeResponse {
  theme: ThemeItem
  config: ThemeConfigDTO
}

export interface SwitchThemeRequest {
  theme_id: number
}

export interface UpdateThemeConfigRequest {
  config: ThemeConfigDTO
}

// ========== Page Types ==========

export interface PageItem {
  id: number
  page_type: string
  name: string
  slug: string
  is_published: boolean
  version: number
}

export interface ListPagesResponse {
  pages: PageItem[]
}

// API returns block_config as JSON string, but we use object internally
export interface DecorationDTOAPI {
  id: number
  block_type: string
  block_config: string // JSON string from API
  sort_order: number
}

export interface DecorationDTO {
  id: number
  block_type: string
  block_config: Record<string, any> // Parsed object for internal use
  sort_order: number
}

export interface PageDetailResponseAPI {
  page: PageItem
  decorations: DecorationDTOAPI[]
}

export interface PageDetailResponse {
  page: PageItem
  decorations: DecorationDTO[]
}

export interface SaveDraftRequest {
  blocks: {
    block_type: string
    block_config: string // JSON string for API
    sort_order: number
  }[]
}

export interface AddDecorationRequest {
  block_type: string
  block_config: string // JSON string for API
  sort_order: number
}

export interface UpdateDecorationRequest {
  block_config: string // JSON string for API
}

export interface ReorderBlocksRequest {
  block_orders: { id: number; sort_order: number }[]
}

// ========== Version Types ==========

export interface VersionItem {
  id: number
  version: number
  created_by: number
  created_at: string
}

export interface ListVersionsResponse {
  versions: VersionItem[]
}

export interface VersionDetailResponse {
  version: VersionItem
  blocks: DecorationDTO[]
}

export interface RestoreVersionRequest {
  version: number
}

// ========== SEO Types ==========

export interface SEOConfigDTO {
  title: string
  description: string
  keywords: string
}

export interface PageSEOConfigDTO {
  page_type: string
  page_id?: number
  config: SEOConfigDTO
}

export interface ListPageSEOConfigsResponse {
  configs: PageSEOConfigDTO[]
}

export interface UpdateSEOConfigRequest {
  title: string
  description: string
  keywords: string
}

// ========== Theme API ==========

export function listThemes() {
  return request<ListThemesResponse>({
    url: '/api/v1/themes',
    method: 'get'
  })
}

export function getCurrentTheme() {
  return request<CurrentThemeResponse>({
    url: '/api/v1/themes/current',
    method: 'get'
  })
}

export function switchTheme(data: SwitchThemeRequest) {
  return request({
    url: '/api/v1/themes/switch',
    method: 'put',
    data
  })
}

export function updateThemeConfig(data: UpdateThemeConfigRequest) {
  return request({
    url: '/api/v1/themes/config',
    method: 'put',
    data
  })
}

// ========== Page API ==========

export function listPages() {
  return request<ListPagesResponse>({
    url: '/api/v1/pages',
    method: 'get'
  })
}

// Helper to parse decoration from API
function parseDecoration(d: DecorationDTOAPI): DecorationDTO {
  let blockConfig: Record<string, any> = {}
  try {
    blockConfig = JSON.parse(d.block_config || '{}')
  } catch (e) {
    console.warn('Failed to parse block_config:', e)
  }
  return {
    id: d.id,
    block_type: d.block_type,
    block_config: blockConfig,
    sort_order: d.sort_order
  }
}

// Helper to stringify decoration for API
function stringifyDecoration(d: DecorationDTO): { block_type: string; block_config: string; sort_order: number } {
  return {
    block_type: d.block_type,
    block_config: JSON.stringify(d.block_config || {}),
    sort_order: d.sort_order
  }
}

export async function getPage(id: number): Promise<PageDetailResponse> {
  const res = await request<PageDetailResponseAPI>({
    url: `/api/v1/pages/${id}`,
    method: 'get'
  })
  return {
    page: res.page,
    decorations: res.decorations?.map(parseDecoration) || []
  }
}

export async function getPageBySlug(slug: string): Promise<PageDetailResponse> {
  const res = await request<PageDetailResponseAPI>({
    url: `/api/v1/pages/slug/${slug}`,
    method: 'get'
  })
  return {
    page: res.page,
    decorations: res.decorations?.map(parseDecoration) || []
  }
}

export function saveDraft(pageId: number, blocks: DecorationDTO[]) {
  return request({
    url: `/api/v1/pages/${pageId}/draft`,
    method: 'put',
    data: {
      blocks: blocks.map(stringifyDecoration)
    }
  })
}

export function publishPage(pageId: number) {
  return request({
    url: `/api/v1/pages/${pageId}/publish`,
    method: 'put'
  })
}

export function unpublishPage(pageId: number) {
  return request({
    url: `/api/v1/pages/${pageId}/unpublish`,
    method: 'put'
  })
}

// ========== Decoration API ==========

export async function getDecorations(pageId: number): Promise<DecorationDTO[]> {
  const res = await request<DecorationDTOAPI[]>({
    url: `/api/v1/pages/${pageId}/decorations`,
    method: 'get'
  })
  return res?.map(parseDecoration) || []
}

export async function addDecoration(pageId: number, data: { block_type: string; block_config: Record<string, any>; sort_order: number }): Promise<DecorationDTO> {
  const res = await request<DecorationDTOAPI>({
    url: `/api/v1/pages/${pageId}/decorations`,
    method: 'post',
    data: {
      block_type: data.block_type,
      block_config: JSON.stringify(data.block_config),
      sort_order: data.sort_order
    }
  })
  return parseDecoration(res)
}

export function updateDecoration(decorationId: number, blockConfig: Record<string, any>) {
  return request({
    url: `/api/v1/decorations/${decorationId}`,
    method: 'put',
    data: {
      block_config: JSON.stringify(blockConfig)
    }
  })
}

export function deleteDecoration(decorationId: number) {
  return request({
    url: `/api/v1/decorations/${decorationId}`,
    method: 'delete'
  })
}

export function reorderBlocks(pageId: number, data: ReorderBlocksRequest) {
  return request({
    url: `/api/v1/pages/${pageId}/blocks/reorder`,
    method: 'put',
    data
  })
}

// ========== Version API ==========

export function listVersions(pageId: number, limit: number = 20) {
  return request<ListVersionsResponse>({
    url: `/api/v1/pages/${pageId}/versions`,
    method: 'get',
    params: { limit }
  })
}

export function getVersion(pageId: number, version: number) {
  return request<VersionDetailResponse>({
    url: `/api/v1/pages/${pageId}/versions/${version}`,
    method: 'get'
  })
}

export function restoreVersion(pageId: number, data: RestoreVersionRequest) {
  return request({
    url: `/api/v1/pages/${pageId}/restore`,
    method: 'put',
    data
  })
}

// ========== SEO API ==========

export function getGlobalSEO() {
  return request<SEOConfigDTO>({
    url: '/api/v1/seo/global',
    method: 'get'
  })
}

export function updateGlobalSEO(data: UpdateSEOConfigRequest) {
  return request({
    url: '/api/v1/seo/global',
    method: 'put',
    data
  })
}

export function listPageSEO() {
  return request<ListPageSEOConfigsResponse>({
    url: '/api/v1/seo/pages',
    method: 'get'
  })
}

export function getPageSEO(pageType: string, pageId?: number) {
  return request<PageSEOConfigDTO>({
    url: `/api/v1/seo/pages/${pageType}`,
    method: 'get',
    params: pageId ? { page_id: pageId } : {}
  })
}

export function updatePageSEO(pageType: string, data: UpdateSEOConfigRequest, pageId?: number) {
  return request({
    url: `/api/v1/seo/pages/${pageType}`,
    method: 'put',
    data,
    params: pageId ? { page_id: pageId } : {}
  })
}

// ========== Block Types Configuration ==========

export const BLOCK_TYPES = [
  { type: 'banner', name: '轮播图', icon: 'Picture' },
  { type: 'product_grid', name: '商品网格', icon: 'Grid' },
  { type: 'rich_text', name: '富文本', icon: 'Document' },
  { type: 'image_carousel', name: '图片轮播', icon: 'PictureFilled' },
  { type: 'featured_products', name: '推荐商品', icon: 'Goods' },
  { type: 'categories', name: '分类展示', icon: 'Menu' },
  { type: 'divider', name: '分割线', icon: 'Minus' },
  { type: 'video', name: '视频', icon: 'VideoPlay' },
  { type: 'spacer', name: '间距', icon: 'Rank' },
  { type: 'custom_html', name: '自定义HTML', icon: 'Document' }
] as const

export type BlockType = typeof BLOCK_TYPES[number]['type']