<template>
  <div class="w-16 border-r dark:border-slate-800">
    <div class="py-6 flex flex-col items-center">
      <div class="font-bold flex-shrink-0 flex items-center">
        mock
      </div>
      <div class="flex-1 mt-6">
        <router-link v-slot="{isActive}" v-for="item in items" :to="item.path"
                     class="p-5 flex flex-col items-center">
          <component :is="item.icon"
                     :class="[isActive ? 'text-green-500' : 'text-gray-900 dark:text-slate-200', 'h-6 w-6 hover:text-green-500 dark:hover:text-green-500']"/>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {CogIcon, CollectionIcon} from '@heroicons/vue/outline';

const items = [
  {path: '/mocks', icon: CollectionIcon, current: true},
  {path: '/settings', icon: CogIcon, current: false}
];
</script>

