<template>
  <Disclosure as="nav">
    <div class="h-20 flex border-b border-gray-200 dark:border-slate-800">
      <div class="flex flex-1 items-center">
        <div class="flex w-72">
          <MockStart class="mx-4" />
          <MockSelect class="flex-1" />
        </div>
        <div class="flex flex-1 justify-between">
          <div class="ml-5">
            <a href="#" class="pr-3 py-2 text-green-500">Routes</a>
            <a href="#" class="px-3 py-2 hover:text-green-500">Proxy</a>
            <a href="#" class="px-3 py-2 hover:text-green-500">Logs</a>
            <a href="#" class="px-3 py-2 hover:text-green-500">Settings</a>
          </div>
          <div class="mr-5 flex">
            <a href="#">
              <ShareIcon
                class="text-gray-500 dark:text-slate-400 dark:hover:text-green-500 hover:text-green-500 h-5 w-5 mr-5" />
            </a>
            <a href="#">
              <SaveIcon
                class="text-gray-500 dark:text-slate-400 hover:text-green-500 dark:hover:text-green-500 h-5 w-5" />
            </a>
          </div>
        </div>
      </div>
    </div>
  </Disclosure>
</template>

<script setup lang="ts">
import { Disclosure } from '@headlessui/vue';
import { ShareIcon, SaveIcon } from '@heroicons/vue/outline';
import MockSelect from './MockSelect.vue';
import MockStart from './MockStart.vue';
</script>
