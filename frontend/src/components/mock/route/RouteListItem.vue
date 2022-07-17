<template>
  <div
      class="flex items-center group block text-sm">
    <router-link v-slot="{ isActive }"
                 class="flex-1"
                 :to="{name: 'route-view', params: {routeId: route.id, id: mock.data.id}}"
    >
      <div
          :class="[isActive ? 'border-green-500' : 'border-transparent', 'pl-3 pr-1 py-1 hover:border-green-500 border-l-2']">
      <span>
          <span :class="`method method-${method.toLowerCase()} mr-1 text-xs`">{{ method }}</span>
          <span>{{ path }}</span>
      </span>
        <p class="text-xs my-1 text-gray-500 dark:text-slate-600">
          {{ route.description }}
        </p>
      </div>

    </router-link>
    <div>
      <Popover class="relative flex">
        <PopoverButton class="p-1 text-white items-center text-base font-medium focus:outline-none">
          <DotsVerticalIcon class="text-gray-900 dark:text-slate-200 w-4 h-4"/>
        </PopoverButton>
        <PopoverPanel class="absolute z-10 right-[10px] top-[15px] transform mt-3 px-2 w-32 max-w-md sm:px-0">
          <div
              class="shadow rounded dark:ring-slate-800 overflow-hidden ring-1 ring-black dark:ring-slate-800 ring-opacity-5">
            <div class="relative grid bg-white dark:bg-slate-900">
              <div class="max-h-96 overflow-auto">
                <div
                    class="hover:text-green-500 border-transparent border-l-2 flex p-2 cursor-pointer">
                  Duplicate
                </div>
                <div
                    class="hover:text-green-500 border-transparent border-l-2 flex p-2 flex cursor-pointer">
                  Delete
                </div>
                <div
                    class="hover:text-green-500 border-transparent border-l-2 items-center flex p-2 flex cursor-pointer">
                  <span>Open</span>
                  <ExternalLinkIcon class="w-4 h-4 ml-2"/>
                </div>
              </div>
            </div>
          </div>
        </PopoverPanel>
      </Popover>
    </div>
  </div>
</template>

<script setup lang="ts">
import {Popover, PopoverButton, PopoverPanel} from '@headlessui/vue';
import {DotsVerticalIcon, ExternalLinkIcon} from '@heroicons/vue/solid';
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


