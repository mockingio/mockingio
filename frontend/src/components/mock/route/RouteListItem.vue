<template>
  <router-link v-slot="{ isActive }"
               :to="{name: 'route-view', params: {routeId: route.id, id: mock.data.id}}"
  >
    <div
        :class="[isActive ? 'bg-slate-800' : 'border-transparent', 'hover:bg-slate-800 flex items-center group block text-sm py-1']">
      <div class="pl-3 pr-1 py-1 flex-1">
      <span class="flex items-center">
          <span :class="`rounded p-1 text-white w-[40px] method-${route.method.toLowerCase()} text-xxs`">{{
              shortMethod
            }}</span>
          <span>{{ route.path }}</span>
      </span>
        <p class="text-xs my-1 text-gray-500 dark:text-slate-600">
          {{ route.description }}
        </p>
      </div>

      <div>
        <PopoverMenu :items="items"/>
      </div>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import type {Mock, Route} from "@/stores";
import PopoverMenu from "@/components/ui/PopoverMenu.vue";
import {computed} from "vue";

const props = defineProps({
  route: {type: Object as () => Route, required: true},
  mock: {type: Object as () => Mock, required: true},
})

const shortMethodMap: { [key: string]: string } = {
  POST: "POST",
  GET: "GET",
  PUT: "PUT",
  PATCH: "PATCH",
  DELETE: "DEL",
  OPTIONS: "OPT"
}

const shortMethod = computed(() => {
  return shortMethodMap[props.route.method]
})

const items = [
  {
    name: "Duplicate",
  },
  {
    name: "Delete",
  },
]
</script>


