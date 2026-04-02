<template>
  <div
    class="skeleton-card"
    :class="{ 'skeleton-card--interactive': interactive }"
  >
    <!-- Header -->
    <div
      v-if="showHeader"
      class="skeleton-card-header"
    >
      <div class="skeleton-title" />
      <div
        v-if="showAction"
        class="skeleton-action"
      />
    </div>

    <!-- Content -->
    <div class="skeleton-card-content">
      <slot>
        <!-- Default content: stats layout -->
        <div
          v-if="variant === 'stats'"
          class="skeleton-stats"
        >
          <div
            v-for="i in statsCount"
            :key="i"
            class="skeleton-stat-item"
          >
            <div class="skeleton-stat-value" />
            <div class="skeleton-stat-label" />
          </div>
        </div>

        <!-- List layout -->
        <div
          v-else-if="variant === 'list'"
          class="skeleton-list"
        >
          <div
            v-for="i in listCount"
            :key="i"
            class="skeleton-list-item"
          >
            <div class="skeleton-avatar" />
            <div class="skeleton-list-content">
              <div class="skeleton-text skeleton-text--lg" />
              <div class="skeleton-text skeleton-text--sm" />
            </div>
          </div>
        </div>

        <!-- Form layout -->
        <div
          v-else-if="variant === 'form'"
          class="skeleton-form"
        >
          <div
            v-for="i in formCount"
            :key="i"
            class="skeleton-form-item"
          >
            <div class="skeleton-label" />
            <div class="skeleton-input" />
          </div>
        </div>

        <!-- Default: simple content -->
        <div
          v-else
          class="skeleton-content-default"
        >
          <div class="skeleton-text" />
          <div class="skeleton-text skeleton-text--md" />
          <div class="skeleton-text skeleton-text--sm" />
        </div>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  showHeader: {
    type: Boolean,
    default: true
  },
  showAction: {
    type: Boolean,
    default: false
  },
  interactive: {
    type: Boolean,
    default: false
  },
  variant: {
    type: String,
    default: 'default', // 'default' | 'stats' | 'list' | 'form'
    validator: (v: string) => ['default', 'stats', 'list', 'form'].includes(v)
  },
  statsCount: {
    type: Number,
    default: 4
  },
  listCount: {
    type: Number,
    default: 5
  },
  formCount: {
    type: Number,
    default: 4
  }
})
</script>

<style scoped>
.skeleton-card {
  background: var(--color-bg-card, #FFFFFF);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid var(--color-border, rgba(99, 102, 241, 0.06));
}

.skeleton-card--interactive {
  transition: box-shadow 0.2s ease;
  animation: pulse-subtle 2s infinite;
}

@keyframes pulse-subtle {
  0%, 100% {
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  }
  50% {
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.1);
  }
}

.skeleton-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--color-border-light, #F3F4F6);
}

.skeleton-title {
  width: 140px;
  height: 20px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 6px;
}

[data-theme="dark"] .skeleton-title {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-action {
  width: 80px;
  height: 32px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 16px;
}

[data-theme="dark"] .skeleton-action {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

/* Stats variant */
.skeleton-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 24px;
}

.skeleton-stat-item {
  text-align: center;
}

.skeleton-stat-value {
  width: 80px;
  height: 32px;
  margin: 0 auto 8px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 6px;
}

[data-theme="dark"] .skeleton-stat-value {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-stat-label {
  width: 60px;
  height: 14px;
  margin: 0 auto;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

[data-theme="dark"] .skeleton-stat-label {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

/* List variant */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.skeleton-list-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.skeleton-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  flex-shrink: 0;
}

[data-theme="dark"] .skeleton-avatar {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-list-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

/* Form variant */
.skeleton-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.skeleton-form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skeleton-label {
  width: 80px;
  height: 14px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

[data-theme="dark"] .skeleton-label {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-input {
  height: 40px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 8px;
}

[data-theme="dark"] .skeleton-input {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

/* Default content */
.skeleton-content-default {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-text {
  height: 14px;
  width: 100%;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

[data-theme="dark"] .skeleton-text {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-text--lg {
  height: 18px;
  width: 60%;
}

.skeleton-text--md {
  width: 80%;
}

.skeleton-text--sm {
  width: 40%;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>