<template>
  <div class="flex">
    <draggable handle=".handle" class="dragArea block list-group w-full" :list="list">
      <Rule v-for="element in list" :key="element.name" class="list-group-item mb-2"/>
    </draggable>
    <div class="ml-5 border-gray-400 dark:border-green-500 border-solid border-2 border-l-0 w-[10px]"></div>
    <div class="flex items-center ml-5">
      <span class="relative z-0 inline-flex shadow-sm rounded-md">
        <button type="button"
                class="relative hover:text-green-500 text-green-500 dark:text-green-500 inline-flex items-center px-4 py-2 rounded-l-md dark:border-slate-700 border border-gray-300 text-xs font-bold focus:outline-none focus:ring-1">
          AND
        </button>
        <button type="button"
                class="-ml-px hover:text-green-500 dark:hover:text-green-500 dark:text-slate-400 relative inline-flex items-center px-4 py-2 dark:border-slate-700 rounded-r-md border border-gray-300 text-xs font-medium text-gray-700 hover:bg-gray-50 focus:z-10 focus:outline-none focus:ring-1">
          OR
        </button>
      </span>
    </div>
  </div>
  <div>
    <a href="#" class="flex justify-start my-5 items-center hover:text-green-500">
      <PlusIcon class="w-4 h-4 mr-2 text-green-500"/>
      Add new rule
    </a>
  </div>
</template>

<script setup lang="ts">
import Rule from './Rule.vue';
import {PlusIcon} from '@heroicons/vue/outline';</script>

<script lang="ts">
import {defineComponent} from 'vue';
import {VueDraggableNext} from 'vue-draggable-next';

export default defineComponent({
  components: {
    draggable: VueDraggableNext
  },
  data() {
    return {
      enabled: true,
      list: [
        {name: 'John', id: 1},
        {name: 'Joao', id: 2},
        {name: 'Jean', id: 3},
        {name: 'Gerard', id: 4}
      ],
      dragging: false
    };
  }
});
</script>
