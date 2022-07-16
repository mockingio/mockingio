<template>
  <div class="grid grid-cols-4 gap-4">
    <div class="col-span-1 flex items-center">
      <TrashIcon class="w-5 h-5 mr-5" />
      <select id="country" name="country" autocomplete="country-name"
        class="block dark:bg-transparent dark:border-slate-700 w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none sm:text-sm">
        <option>Body</option>
        <option>Request Number</option>
        <option>Cookie</option>
      </select>
    </div>

    <div class="col-span-1">
      <input type="text" name="city" id="city" autocomplete="address-level2"
        class=" dark:bg-transparent dark:border-slate-700 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md" />
    </div>

    <div class="col-span-1">
      <select id="country" name="country" autocomplete="country-name"
        class="block w-full py-2 px-3  dark:bg-transparent dark:border-slate-700 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none sm:text-sm">
        <option>Equals</option>
        <option>Regex match</option>
      </select>
    </div>

    <div class="col-span-1">
      <div class="flex items-center">
        <input type="text" name="postal-code" id="postal-code" autocomplete="postal-code"
          class="dark:bg-transparent dark:border-slate-700 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md" />
        <ViewListIcon class="cursor-move handle w-5 h-5 ml-5" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { TrashIcon, ViewListIcon } from '@heroicons/vue/outline';
</script>
