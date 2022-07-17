<template>
  <div class="flex-1 flex items-center">
    <div class="flex-1" v-if="activeMock">
      <h2 class="font-semibold">{{ activeMock.data.name }}</h2>
      <p class="text-xs text-gray-500 dark:text-slate-400">
        <span v-if="activeMock.state.url">{{ activeMock.state.url }}</span>
        <span v-else>offline</span>
      </p>
    </div>
    <Popover class="relative flex">
      <PopoverButton class="text-white items-center text-base font-medium focus:outline-none">
        <DotsVerticalIcon class="text-gray-900 dark:text-slate-200 w-6 h-6"/>
      </PopoverButton>

      <PopoverPanel class="absolute z-10 right-[10px] top-[15px] transform mt-3 px-2  w-64">
        <div
            class="shadow overflow-hidden grid bg-white dark:bg-slate-900 rounded ring-1 ring-black dark:ring-slate-800 ring-opacity-5">
          <div
              class="relative">
            <div class="max-h-96 overflow-auto">
              <MockSelectItem v-for="mock in mocks" :key="mock.data.id" :mock="mock"/>
            </div>
            <div class="mb-2 mr-2 mt-2 flex justify-end">
              <button type="button"
                      class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-white dark:text-green-500 focus:outline-none">
                New mock
              </button>
            </div>
          </div>
        </div>
      </PopoverPanel>
    </Popover>
  </div>
</template>

<script setup lang="ts">
import {Popover, PopoverButton, PopoverPanel} from '@headlessui/vue';
import {DotsVerticalIcon} from '@heroicons/vue/solid';
import MockSelectItem from './MockSelectItem.vue';</script>

<script lang="ts">
import type {Mock} from "@/stores";

export default {
  props: {
    mocks: {type: Object as () => Mock[], required: true},
    activeMock: {type: Object as () => Mock, required: false}
  }
};
</script>
