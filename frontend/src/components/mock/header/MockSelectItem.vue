<template>
  <router-link @click="close" :to="{name: 'routes-view', params: {id: mock.data.id}}"
               class="p-2 flex items-start border-transparent border-l-2 hover:text-green-500">
    <div>
      <p class="text-base flex items-center text-sm text-gray-900 dark:text-slate-400 dark:hover:text-green-500  cursor-pointer">
        <span class="inline-flex items-center">
          <svg :class="[mock.state.status === 'running' ? ' text-green-500' : 'text-gray-500', 'mr-1.5 h-2 w-2']"
               fill="currentColor" viewBox="0 0 8 8">
            <circle cx="4" cy="4" r="3"/>
          </svg>

          {{ mock.data.name }}
        </span>
      </p>
      <p class="mt-1 text-sm text-gray-500 dark:text-slate-400 dark:text-slate-600">
        {{ mock.data.description }}
      </p>
    </div>
  </router-link>
</template>

<script lang="ts">
import type {Mock} from "@/stores";

export default {
  props: {
    close: {type: Function as () => void, required: true},
    mock: {type: Object as () => Mock, required: true},
  }
};
</script>
