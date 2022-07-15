<template>
  <div class="w-16 border-r dark:border-slate-800">
    <div class="py-6 flex flex-col items-center">
      <div class="font-bold flex-shrink-0 flex items-center">
        mock
      </div>
      <div class="flex-1 mt-6">
        <router-link v-for="item in items" :to="item.href"
                     active-class="border-green-500" class="p-5 flex flex-col items-center">
          <component :is="item.icon"
                     class="text-gray-900 dark:text-slate-200 h-6 w-6 hover:text-green-500 dark:hover:text-green-500"/>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {CogIcon, CollectionIcon} from '@heroicons/vue/outline';

const items = [
  {href: '/', icon: CollectionIcon, current: true},
  {href: '/setting', icon: CogIcon, current: false}
];
</script>

