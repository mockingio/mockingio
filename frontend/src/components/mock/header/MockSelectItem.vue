<template>
  <router-link :to="{name: 'routes-view', params: {id}}"
               class="p-2 flex items-start hover:border-green-500 border-transparent border-l-2">
    <div>
      <p class="text-base flex items-center font-medium text-gray-900 dark:text-slate-400">
        <PlayIcon class="text-green-500 h-4 w-4 mr-2"/>
        <span>{{ name }}</span>
      </p>
      <p class="mt-1 text-sm text-gray-500 dark:text-slate-400 dark:text-slate-600">
        {{ description }}
      </p>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import {PlayIcon} from '@heroicons/vue/solid';</script>

<script lang="ts">
export default {
  props: {
    id: {type: String, required: true},
    name: {type: String, required: true},
    description: {type: String, required: false}
  }
};
</script>
