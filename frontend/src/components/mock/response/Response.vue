<template>
  <div>
    <div class="flex dark:bg-slate-800 bg-white">
      <div class="m-5 block flex-1 flex justify-between">
        <p>{{ response.status }}</p>
        <ViewListIcon class="handle w-5 h-5 ml-5"/>
      </div>
    </div>

    <div class="hidden">
      Detail here
    </div>
  </div>
</template>

<script setup lang="ts">
import type {Response} from "@/stores";
import {ViewListIcon} from '@heroicons/vue/outline';

defineProps({
  response: {type: Object as () => Response, required: true}
})
</script>
