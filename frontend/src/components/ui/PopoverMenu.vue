<template>
  <Popover class="relative flex">
    <PopoverButton class="p-1 text-white items-center text-base font-medium focus:outline-none">
      <DotsVerticalIcon class="text-gray-900 dark:text-slate-500 w-4 h-4"/>
    </PopoverButton>
    <PopoverPanel v-slot="{close}"
                  class="absolute z-10 right-[10px] top-[15px] transform mt-3 px-2 w-32 max-w-md sm:px-0">
      <div
          class="shadow rounded dark:ring-slate-800 overflow-hidden ring-1 ring-black dark:ring-slate-800 ring-opacity-5">
        <div class="relative grid bg-white dark:bg-slate-900">
          <div class="max-h-96 overflow-auto">
            <div @click="onItemClick(item, close)" v-for="item in items" :key="item.name"
                 class="hover:text-green-500 border-transparent border-l-2 flex items-center p-2 cursor-pointer">
              {{ item.name }}
            </div>
          </div>
        </div>
      </div>
    </PopoverPanel>
  </Popover>
</template>

<script setup lang="ts">
import type {PropType} from "vue";
import {Popover, PopoverButton, PopoverPanel} from '@headlessui/vue';
import {DotsVerticalIcon} from '@heroicons/vue/solid';

export interface Item {
  name: string
  path?: string
  click?: () => void
}

const onItemClick = (item: Item, close: () => void) => {
  close()
  if (!item.click) {
    return;
  }

  item.click()
}

const props = defineProps({
  items: {type: Array as PropType<Array<Item>>, required: true},
})

</script>