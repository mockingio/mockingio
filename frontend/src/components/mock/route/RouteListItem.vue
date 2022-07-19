<template>
  <router-link v-slot="{ isActive }"
               :to="{name: 'route-view', params: {routeId: route.id, id: mock.data.id}}"
  >
    <div
        :class="[isActive ? 'bg-slate-800' : 'border-transparent', 'hover:bg-slate-800 flex items-center group block text-sm py-1']">
      <div class="pl-3 pr-1 py-1 flex-1">
      <span class="flex items-center justify-between">
          <span>{{ route.path }}</span>
          <span :class="`rounded p-1 text-white method-${route.method.toLowerCase()} text-xs`">{{
              route.method
            }}</span>
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

const props = defineProps({
  route: {type: Object as () => Route, required: true},
  mock: {type: Object as () => Mock, required: true},
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


