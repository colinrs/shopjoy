<template>
  <div ref="chartRef" class="echarts-container"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'
import type { TrendDataPoint } from '@/api/points'

const { t } = useI18n()

interface Props {
  data: TrendDataPoint[]
  loading?: boolean
}

const props = defineProps<Props>()

const chartRef = ref<HTMLElement>()
let chartInstance: echarts.ECharts | null = null

const initChart = () => {
  if (!chartRef.value) return

  chartInstance = echarts.init(chartRef.value)
  updateChart()
}

const updateChart = () => {
  if (!chartInstance) return

  const dates = props.data.map(d => formatDateLabel(d.date))
  const earnedData = props.data.map(d => d.earned)
  const redeemedData = props.data.map(d => d.redeemed)
  const expiredData = props.data.map(d => d.expired)

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#E5E7EB',
      borderWidth: 1,
      textStyle: {
        color: '#1E1B4B'
      },
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6B7280'
        }
      },
      formatter: (params: unknown) => {
        if (!Array.isArray(params) || params.length === 0) return ''
        const typedParams = params as Array<{ axisValue?: string; color?: unknown; seriesName?: string; value?: number }>
        const firstParam = typedParams[0]
        const date = firstParam.axisValue
        let html = `<div style="font-weight: 600; margin-bottom: 8px;">${date}</div>`
        typedParams.forEach((param) => {
          const color = param.color instanceof Object ? (param.color as { colorStops?: Array<{ color?: string }> }).colorStops?.[0]?.color : param.color
          html += `<div style="display: flex; align-items: center; gap: 8px; margin: 4px 0;">
            <span style="display: inline-block; width: 10px; height: 10px; border-radius: 2px; background: ${color};"></span>
            <span>${param.seriesName}: <strong>${param.value?.toLocaleString() ?? 0}</strong></span>
          </div>`
        })
        return html
      }
    },
    legend: {
      data: [t('points.earned'), t('points.redeemed2'), t('points.expired')],
      bottom: 0,
      textStyle: {
        color: '#6B7280'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: {
        lineStyle: {
          color: '#E5E7EB'
        }
      },
      axisLabel: {
        color: '#9CA3AF'
      }
    },
    yAxis: {
      type: 'value',
      splitLine: {
        lineStyle: {
          color: '#F3F4F6'
        }
      },
      axisLine: {
        show: false
      },
      axisLabel: {
        color: '#9CA3AF'
      }
    },
    series: [
      {
        name: t('points.earned'),
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        data: earnedData,
        lineStyle: {
          color: '#10B981',
          width: 3
        },
        itemStyle: {
          color: '#10B981'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(16, 185, 129, 0.3)' },
            { offset: 1, color: 'rgba(16, 185, 129, 0.05)' }
          ])
        }
      },
      {
        name: t('points.redeemed2'),
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        data: redeemedData,
        lineStyle: {
          color: '#3B82F6',
          width: 3
        },
        itemStyle: {
          color: '#3B82F6'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
            { offset: 1, color: 'rgba(59, 130, 246, 0.05)' }
          ])
        }
      },
      {
        name: t('points.expired'),
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        data: expiredData,
        lineStyle: {
          color: '#EF4444',
          width: 3,
          type: 'dashed'
        },
        itemStyle: {
          color: '#EF4444'
        }
      }
    ]
  }

  chartInstance.setOption(option)
}

const formatDateLabel = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

const handleResize = () => {
  chartInstance?.resize()
}

onMounted(() => {
  nextTick(() => {
    initChart()
    window.addEventListener('resize', handleResize)
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
})

watch(() => props.data, () => {
  nextTick(() => {
    updateChart()
  })
}, { deep: true })
</script>

<style scoped>
.echarts-container {
  width: 100%;
  height: 320px;
}
</style>
