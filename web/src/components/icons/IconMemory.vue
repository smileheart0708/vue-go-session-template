<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 120 90"
    fill="none"
    stroke="currentColor"
    stroke-width="8"
    stroke-linecap="round"
    stroke-linejoin="miter"
  >
    <!-- 主体外框：顶部使用圆角，底部使用直角 -->
    <path
      d="M 10 60 V 20 Q 10 10 20 10 H 100 Q 110 10 110 20 V 60 H 10 Z"
      stroke-linejoin="round"
    />
    <!-- 内部两个芯片：保持微小圆角 -->
    <rect x="23" y="25" width="30" height="20" rx="2" stroke-linejoin="round" />
    <rect x="67" y="25" width="30" height="20" rx="2" stroke-linejoin="round" />
    <!-- 底部引脚：6个引脚，间距相等，最左和最右与外框侧边对齐 -->
    <path d="M 10 60 V 78" />
    <path d="M 30 60 V 78" />
    <path d="M 50 60 V 78" />
    <path d="M 70 60 V 78" />
    <path d="M 90 60 V 78" />
    <path d="M 110 60 V 78" />
  </svg>
</template>
