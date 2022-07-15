<template>
  <router-link active-class="border-green-500" :to="{name: 'route-view', params: {routeId: route.id, id: mock.data.id}}"
               class="hover:border-green-500 border-transparent border-l-2 flex items-center group block pl-3 pr-1 py-2 border-l-2 text-sm">
    <div class="flex-1">
      <span>
          <span :class="`method method-${method.toLowerCase()} mr-2`">{{ method }}</span>
          <span>{{ path }}</span>
      </span>
      <p class="text-xs text-gray-500 dark:text-slate-400">
        {{ route.description }}
      </p>
    </div>
    <div>
      <Popover class="relative flex">
        <PopoverButton class="text-white items-center text-base font-medium focus:outline-none">
          <DotsVerticalIcon class="text-gray-900 dark:text-slate-200 w-4 h-4"/>
        </PopoverButton>
        <PopoverPanel class="absolute z-10 left-[10px] transform mt-3 px-2 w-32 max-w-md sm:px-0">
          <div
              class="shadow rounded-lg dark:ring-slate-800 overflow-hiddenring-1 ring-black dark:ring-slate-800 ring-opacity-5">
            <div class="relative grid bg-white dark:bg-slate-900 p-3">
              <div class="max-h-96 overflow-auto">
                <div
                    class="hover:border-green-500 border-transparent border-l-2 flex p-2 mb-3 items-start transition ease-in-out duration-150">
                  Duplicate
                </div>
                <div
                    class="hover:border-green-500 border-transparent border-l-2 flex p-2 flex items-start transition ease-in-out duration-150">
                  Delete
                </div>
                <div
                    class="hover:border-green-500 border-transparent border-l-2 flex p-2 flex items-start transition ease-in-out duration-150">
                  Open [->]
                </div>
              </div>
            </div>
          </div>
        </PopoverPanel>
      </Popover>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import {Popover, PopoverButton, PopoverPanel} from '@headlessui/vue';
import {DotsVerticalIcon} from '@heroicons/vue/solid';
import {computed} from "vue";
import type {Mock, Route} from "@/stores";

const props = defineProps({
  route: {type: Object as () => Route, required: true},
  mock: {type: Object as () => Mock, required: true},
})

const method = computed(() => {
  const [method] = props.route.request.split(' ')
  return method
})

const path = computed(() => {
  const [_, path] = props.route.request.split(' ')
  return path
})
</script>


