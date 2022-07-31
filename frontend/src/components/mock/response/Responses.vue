<template>
  <div class="overflow-hidden sm:rounded-md">
    <div role="list" class="text-sm">
      <draggable handle=".handle"
                 class="dragArea block list-group w-full divide-y-4 divide-solid dark:divide-slate-900">
        <Response v-for="item in list" :key="item.id" :response="item"/>
      </draggable>
    </div>
  </div>
  <div>
    <a href="#" class="flex justify-start my-5 items-center hover:text-green-500">
      <PlusIcon class="w-4 h-4 mr-2 text-green-500"/>
      <span class="text-sm">Add response</span>
    </a>
  </div>
</template>

<script setup lang="ts">
import {PlusIcon} from '@heroicons/vue/outline';
import Response from './Response.vue';</script>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import {VueDraggableNext} from 'vue-draggable-next';
import type {Mock, Route} from "@/stores";

export default defineComponent({
  components: {
    draggable: VueDraggableNext
  },
  props: {
    mock: {type: Object as () => Mock, required: true},
    route: {type: Object as () => Route, required: true}
  },
  data() {
    return {
      enabled: true,
      list: computed(() => this.route.responses),
      dragging: false
    };
  }
});
</script>
