<template>
  <div class="flex-1 flex items-center">
    <div class="flex-1" v-if="activeMock">
      <h2 class="font-semibold">{{ activeMock.data.name }}</h2>
      <p class="text-xs text-gray-500 dark:text-slate-400" v-if="activeMock.state.url">
        <span>{{ activeMock.state.url }}</span>
      </p>
    </div>
    <Popover class="relative flex">
      <PopoverButton class="text-white items-center text-base font-medium focus:outline-none">
        <DotsVerticalIcon class="text-gray-900 dark:text-slate-200 w-6 h-6"/>
      </PopoverButton>

      <PopoverPanel class="absolute z-10 left-[10px] top-[15px] transform mt-3 px-2 w-screen max-w-xs sm:px-0">
        <div class="shadow-lg ring-opacity-5 overflow-hidden">
          <div class="relative grid bg-white dark:bg-slate-800">
            <div class="max-h-96 overflow-auto">
              <MockSelectItem v-for="mock in mocks" :key="mock.data.id" :name="mock.data.name"
                              :description="mock.data.description"
                              class="p-4 flex items-start hover:border-green-500 border-transparent border-l-2">
              </MockSelectItem>
            </div>
            <div class="mb-2 mr-2 mt-2 flex justify-end">
              <button type="button"
                      class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-white bg-green-500 focus:outline-none">
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
import MockSelectItem from './MockSelectItem.vue';

import {storeToRefs} from 'pinia'
import {useMockStore} from "@/stores";

const {mocks, error, activeMock} = storeToRefs(useMockStore())
const {fetchMocks} = useMockStore()

fetchMocks()
</script>
