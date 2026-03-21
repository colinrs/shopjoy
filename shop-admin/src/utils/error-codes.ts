/**
 * Admin Error Code Mapping
 * Maps backend business error codes to user-friendly messages and actions
 */

import type { ErrorCodeHandler, ErrorResponse } from './types'

export interface AdminErrorCodeConfig {
  message: string
  action?: 'logout' | 'redirect' | 'toast' | 'none'
  redirectPath?: string
}

// Admin-specific error code mappings
export const adminErrorCodeMap: Record<number, AdminErrorCodeConfig> = {
  // Auth errors (401xx)
  40100: { message: '未授权，请先登录', action: 'logout', redirectPath: '/login' },
  40101: { message: '登录已过期，请重新登录', action: 'logout', redirectPath: '/login' },
  40102: { message: '无效的身份凭证', action: 'logout', redirectPath: '/login' },

  // Admin User errors (10xxx)
  10001: { message: '邮箱格式不正确', action: 'toast' },
  10002: { message: '手机号格式不正确', action: 'toast' },
  10003: { message: '密码强度不足', action: 'toast' },
  10004: { message: '管理员不存在', action: 'toast' },
  10005: { message: '该账号已存在', action: 'toast' },
  10006: { message: '密码错误', action: 'toast' },
  10007: { message: '不能删除自己的账号', action: 'toast' },
  10008: { message: '该用户已被删除', action: 'toast' },
  10009: { message: '账号已被禁用', action: 'toast' },
  10010: { message: '两次密码不一致', action: 'toast' },
  10011: { message: '没有操作权限', action: 'toast' },

  // Auth Module (130xxx)
  130001: { message: '身份凭证无效', action: 'logout', redirectPath: '/login' },
  130002: { message: '身份凭证已过期', action: 'logout', redirectPath: '/login' },

  // Product errors (30xxx)
  30001: { message: '商品名称不能为空', action: 'toast' },
  30002: { message: '商品价格必须大于0', action: 'toast' },
  30003: { message: '币种不匹配', action: 'toast' },
  30004: { message: '金额不足', action: 'toast' },
  30005: { message: '商品已删除', action: 'toast' },
  30006: { message: '无效的状态转换', action: 'toast' },
  30007: { message: '库存不能为0', action: 'toast' },
  30008: { message: '库存不能为负数', action: 'toast' },
  30009: { message: '商品未上架', action: 'toast' },
  30010: { message: '无效的数量', action: 'toast' },
  30011: { message: '库存不足', action: 'toast' },
  30012: { message: '商品不存在', action: 'toast' },
  30013: { message: '商品ID无效', action: 'toast' },

  // Category errors (301xx)
  30101: { message: '分类不存在', action: 'toast' },
  30102: { message: '分类已存在', action: 'toast' },
  30103: { message: '分类无效', action: 'toast' },
  30104: { message: '该分类下存在子分类', action: 'toast' },
  30105: { message: '该分类下存在商品', action: 'toast' },

  // Order errors (40xxx)
  40001: { message: '订单不存在', action: 'toast' },
  40002: { message: '订单状态无效', action: 'toast' },
  40003: { message: '订单已支付', action: 'toast' },
  40004: { message: '订单未支付', action: 'toast' },
  40005: { message: '订单已过期', action: 'toast' },
  40006: { message: '库存不足', action: 'toast' },
  40007: { message: '金额无效', action: 'toast' },
  40008: { message: '购物车为空', action: 'toast' },

  // Tenant errors (90xxx)
  90001: { message: '租户不存在', action: 'toast' },
  90002: { message: '租户已存在', action: 'toast' },
  90003: { message: '域名格式无效', action: 'toast' },
  90004: { message: '租户已停用', action: 'toast' },
  90005: { message: '无法停用过期租户', action: 'toast' },
  90006: { message: '租户ID无效', action: 'toast' },
  90007: { message: '租户名称不能为空', action: 'toast' },
  90008: { message: '租户编码不能为空', action: 'toast' },

  // Role errors (100xxx)
  100001: { message: '角色不存在', action: 'toast' },
  100002: { message: '角色已存在', action: 'toast' },
  100003: { message: '角色无效', action: 'toast' },

  // Market errors (150xxx)
  150001: { message: '市场不存在', action: 'toast' },
  150002: { message: '市场编码已存在', action: 'toast' },
  150003: { message: '市场编码不能为空', action: 'toast' },
  150004: { message: '市场名称不能为空', action: 'toast' },
  150005: { message: '市场币种不能为空', action: 'toast' },
  150006: { message: '市场已停用', action: 'toast' },
  150007: { message: '该市场已是默认市场', action: 'toast' },
  150008: { message: '不能删除默认市场', action: 'toast' },

  // ProductMarket errors (160xxx)
  160001: { message: '商品市场配置不存在', action: 'toast' },
  160002: { message: '商品已在该市场中', action: 'toast' },
  160003: { message: '市场价格不能为空', action: 'toast' },

  // Inventory errors (170xxx)
  170001: { message: '可用库存不足', action: 'toast' },
  170002: { message: '锁定库存不足', action: 'toast' },
  170003: { message: 'SKU不存在', action: 'toast' },
  170004: { message: '仓库不存在', action: 'toast' },
  170005: { message: '仓库编码已存在', action: 'toast' },

  // Brand errors (180xxx)
  180001: { message: '品牌不存在', action: 'toast' },
  180002: { message: '品牌名称已存在', action: 'toast' },
  180003: { message: '品牌下存在商品，无法删除', action: 'toast' },

  // Common errors
  20001: { message: '数据验证失败', action: 'toast' },
  20002: { message: '数据库错误', action: 'toast' },
  20003: { message: '业务已存在', action: 'toast' },
  20004: { message: '业务不存在', action: 'toast' },
}

// HTTP status code handlers for admin
export const adminHttpStatusHandlers: Record<number, AdminErrorCodeConfig> = {
  401: { message: '未授权，请先登录', action: 'logout', redirectPath: '/login' },
  403: { message: '没有权限访问', action: 'toast' },
  404: { message: '请求的资源不存在', action: 'toast' },
  405: { message: '请求方法不允许', action: 'toast' },
  429: { message: '请求过于频繁，请稍后再试', action: 'toast' },
  500: { message: '服务器内部错误', action: 'toast' },
  502: { message: '网关错误', action: 'toast' },
  503: { message: '服务暂时不可用', action: 'toast' },
  504: { message: '网关超时', action: 'toast' },
}

export const handleAdminError: ErrorCodeHandler = (error: ErrorResponse) => {
  // First check business error code
  if (error.code && adminErrorCodeMap[error.code]) {
    return adminErrorCodeMap[error.code]
  }

  // Then check HTTP status code
  if (error.httpStatus && adminHttpStatusHandlers[error.httpStatus]) {
    return adminHttpStatusHandlers[error.httpStatus]
  }

  // Default error
  return {
    message: error.msg || error.message || '请求失败',
    action: 'toast'
  }
}