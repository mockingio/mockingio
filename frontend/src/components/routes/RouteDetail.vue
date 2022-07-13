<template>
  <div class="m-5">
    <div class="mt-1 flex rounded-md shadow-sm">
              <span
                  class="inline-flex dark:bg-slate-800 dark:border-slate-700 items-center px-3 rounded-l-md border border-r-0 border-gray-300 bg-gray-50 text-gray-500 dark:text-slate-400 sm:text-sm">
                {{ method }}
              </span>
      <input :value="path" type="text" name="path" id="path"
             class="flex-1 dark:bg-slate-800 block dark:border-slate-700 w-full min-w-0 rounded-none rounded-r-md sm:text-sm border-gray-300"/>
    </div>
  </div>

  <h3 class="text-md m-5 text-gray-900">Responses</h3>
  <div class="m-5">
    <Responses :mock="mock" :route="route"/>
  </div>
</template>

<script setup lang="ts">
import {computed} from "vue";
import Responses from "@/components/mock/response/Responses.vue";
import type {Mock, Route} from "@/stores";

const props = defineProps({
  mock: {type: Object as () => Mock, required: true},
  route: {type: Object as () => Route, required: true}
})

const method = computed(() => {
  const [method] = props.route.request.split(' ')
  return method
})

const path = computed(() => {
  const [_, path] = props.route.request.split(' ')
  return path
})
</script>
