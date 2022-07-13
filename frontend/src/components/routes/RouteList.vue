<template>
  <form class="sm:flex sm:items-center m-2">
    <div class="w-full sm:max-w-xs relative">
      <input type="email" name="email" id="email"
             class="shadow-sm block w-full sm:text-sm border-gray-300 dark:border-slate-700 rounded-md dark:bg-slate-800"
             placeholder="find route"/>
      <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
        <FilterIcon class="h-5 w-5 dark:text-slate-500 text-gray-400" aria-hidden="true"/>
      </div>
    </div>
    <PlusIcon class="w-5 h-5 ml-2 text-green-500"/>
  </form>
  <nav class="space-y-1" v-if="activeMock">
    <RouteListItem v-for="route in activeMock.data.routes" :route="route" :mock="activeMock"
                   :key="route.id" :request="route.request"/>
  </nav>
</template>

<script setup lang="ts">
import {FilterIcon, PlusIcon} from '@heroicons/vue/outline';
import type {Mock} from "@/stores";
import RouteListItem from "@/components/routes/RouteListItem.vue";

defineProps({
  activeMock: {type: Object as () => Mock, required: true}
})
</script>
