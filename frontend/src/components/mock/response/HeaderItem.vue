<template>
  <div class="grid grid-cols-2 gap-4">
    <div class="col-span-1 flex items-center">
      <TrashIcon class="w-5 h-5 mr-5"/>
      <input type="text" name="name" :value="name"
             class="dark:bg-transparent dark:border-slate-700 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"/>
    </div>

    <div class="col-span-1">
      <input type="text" name="value" :value="value"
             class=" dark:bg-transparent dark:border-slate-700 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"/>
    </div>
  </div>
</template>

<script setup lang="ts">
import {TrashIcon} from '@heroicons/vue/outline'
import {defineProps} from "vue";

defineProps({
  name: {type: String, required: true},
  value: {type: String, required: true}
})</script>
