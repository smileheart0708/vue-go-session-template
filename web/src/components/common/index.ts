import type { App, Plugin } from 'vue'
import AppPagination from './AppPagination.vue'
import AppSwitch from './AppSwitch.vue'
import AppTable from './AppTable.vue'
import BaseButton from './BaseButton.vue'
import DropdownDrawer from './DropdownDrawer.vue'
import IconButton from './IconButton.vue'
import StatsCard from './StatsCard.vue'
import ThemeToggle from './ThemeToggle.vue'
import ToastMessage from './ToastMessage.vue'

const commonComponents = [
  ['AppPagination', AppPagination],
  ['AppSwitch', AppSwitch],
  ['AppTable', AppTable],
  ['BaseButton', BaseButton],
  ['DropdownDrawer', DropdownDrawer],
  ['IconButton', IconButton],
  ['StatsCard', StatsCard],
  ['ThemeToggle', ThemeToggle],
  ['ToastMessage', ToastMessage],
] as const

export const CommonComponentsPlugin: Plugin = {
  install(app: App): void {
    for (const [name, component] of commonComponents) {
      app.component(name, component)
    }
  },
}

export {
  AppPagination,
  AppSwitch,
  AppTable,
  BaseButton,
  DropdownDrawer,
  IconButton,
  StatsCard,
  ThemeToggle,
  ToastMessage,
}

export type {
  AppTableColumn,
  ColumnAlign,
  ColumnKey,
  CssSize,
  GlassTableColumn,
} from './AppTable.vue'

export default CommonComponentsPlugin
