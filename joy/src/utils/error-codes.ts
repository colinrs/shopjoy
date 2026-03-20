/**
 * Customer (C-end) Error Code Mapping
 * Maps backend business error codes to user-friendly messages and actions
 */

import type { ErrorCodeHandler, ErrorResponse } from './types'

export interface CustomerErrorCodeConfig {
  message: string
  action?: 'logout' | 'redirect' | 'toast' | 'none'
  redirectPath?: string
}

// Customer-specific error code mappings
export const customerErrorCodeMap: Record<number, CustomerErrorCodeConfig> = {
  // Auth errors (401xx)
  40100: { message: '请先登录', action: 'redirect', redirectPath: '/login' },
  40101: { message: '登录已过期，请重新登录', action: 'logout', redirectPath: '/login' },
  40102: { message: '登录信息无效，请重新登录', action: 'logout', redirectPath: '/login' },

  // User errors (11xxx)
  11001: { message: '邮箱格式不正确', action: 'toast' },
  11002: { message: '手机号格式不正确', action: 'toast' },
  11003: { message: '密码强度不足，请设置更复杂的密码', action: 'toast' },
  11004: { message: '用户不存在', action: 'toast' },
  11005: { message: '该账号已注册，请直接登录', action: 'toast' },
  11006: { message: '密码错误，请重试', action: 'toast' },
  11007: { message: '该账号已被注销', action: 'toast' },
  11008: { message: '两次输入的密码不一致', action: 'toast' },

  // Auth Module (130xxx)
  130001: { message: '登录已失效，请重新登录', action: 'logout', redirectPath: '/login' },
  130002: { message: '登录已过期，请重新登录', action: 'logout', redirectPath: '/login' },

  // Product errors (30xxx)
  30001: { message: '商品名称不能为空', action: 'toast' },
  30002: { message: '商品价格异常', action: 'toast' },
  30003: { message: '币种不匹配', action: 'toast' },
  30004: { message: '余额不足', action: 'toast' },
  30005: { message: '商品已下架', action: 'toast' },
  30006: { message: '商品状态异常', action: 'toast' },
  30007: { message: '商品库存不足', action: 'toast' },
  30008: { message: '库存异常', action: 'toast' },
  30009: { message: '商品未上架', action: 'toast' },
  30010: { message: '数量无效', action: 'toast' },
  30011: { message: '商品库存不足', action: 'toast' },
  30012: { message: '商品不存在', action: 'toast' },
  30013: { message: '商品信息错误', action: 'toast' },

  // Order errors (40xxx)
  40001: { message: '订单不存在', action: 'toast' },
  40002: { message: '订单状态异常', action: 'toast' },
  40003: { message: '订单已支付，请勿重复操作', action: 'toast' },
  40004: { message: '订单尚未支付', action: 'toast' },
  40005: { message: '订单已过期，请重新下单', action: 'toast' },
  40006: { message: '商品库存不足，请调整购买数量', action: 'toast' },
  40007: { message: '订单金额异常', action: 'toast' },
  40008: { message: '购物车为空，请先添加商品', action: 'toast' },

  // Payment errors (50xxx)
  50001: { message: '支付记录不存在', action: 'toast' },
  50002: { message: '支付金额异常', action: 'toast' },
  50003: { message: '支付失败，请重试', action: 'toast' },
  50004: { message: '该订单已支付', action: 'toast' },
  50005: { message: '支付已超时，请重新下单', action: 'toast' },

  // Cart errors (60xxx)
  60001: { message: '购物车商品不存在', action: 'toast' },
  60002: { message: '商品数量无效', action: 'toast' },
  60003: { message: '购物车为空', action: 'toast' },

  // Coupon errors (70xxx)
  70001: { message: '优惠券不存在', action: 'toast' },
  70002: { message: '优惠券已过期', action: 'toast' },
  70003: { message: '优惠券已被领完', action: 'toast' },
  70004: { message: '优惠券活动尚未开始', action: 'toast' },
  70005: { message: '该优惠券已使用', action: 'toast' },
  70006: { message: '优惠券码无效', action: 'toast' },
  70007: { message: '订单金额未达到优惠券使用门槛', action: 'toast' },

  // Promotion errors (80xxx)
  80001: { message: '活动不存在', action: 'toast' },
  80002: { message: '活动信息无效', action: 'toast' },
  80003: { message: '活动已结束', action: 'toast' },
  80004: { message: '活动尚未开始', action: 'toast' },

  // Common errors
  10003: { message: '请求参数有误', action: 'toast' },
  10004: { message: '系统繁忙，请稍后重试', action: 'toast' },
  20001: { message: '数据验证失败', action: 'toast' },
  20002: { message: '系统繁忙，请稍后重试', action: 'toast' },
  20003: { message: '该数据已存在', action: 'toast' },
  20004: { message: '数据不存在', action: 'toast' },
}

// HTTP status code handlers for customer
export const customerHttpStatusHandlers: Record<number, CustomerErrorCodeConfig> = {
  401: { message: '请先登录', action: 'redirect', redirectPath: '/login' },
  403: { message: '没有权限访问', action: 'toast' },
  404: { message: '请求的资源不存在', action: 'toast' },
  405: { message: '请求方法不允许', action: 'toast' },
  429: { message: '操作过于频繁，请稍后再试', action: 'toast' },
  500: { message: '服务器开小差了，请稍后重试', action: 'toast' },
  502: { message: '网络异常，请稍后重试', action: 'toast' },
  503: { message: '服务暂时不可用，请稍后重试', action: 'toast' },
  504: { message: '网络超时，请稍后重试', action: 'toast' },
}

export const handleCustomerError: ErrorCodeHandler = (error: ErrorResponse) => {
  // First check business error code
  if (error.code && customerErrorCodeMap[error.code]) {
    return customerErrorCodeMap[error.code]
  }

  // Then check HTTP status code
  if (error.httpStatus && customerHttpStatusHandlers[error.httpStatus]) {
    return customerHttpStatusHandlers[error.httpStatus]
  }

  // Default error
  return {
    message: error.msg || error.message || '请求失败，请稍后重试',
    action: 'toast'
  }
}