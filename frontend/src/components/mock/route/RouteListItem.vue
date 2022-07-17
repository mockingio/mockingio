<template>
  <div
      class="flex items-center group block text-sm">
    <router-link v-slot="{ isActive }"
                 class="flex-1"
                 :to="{name: 'route-view', params: {routeId: route.id, id: mock.data.id}}"
    >
      <div
          :class="[isActive ? 'border-green-500' : 'border-transparent', 'pl-3 pr-1 py-1 hover:border-green-500 border-l-2']">
      <span>
          <span :class="`method method-${route.method.toLowerCase()} mr-1 text-xs`">{{ route.method }}</span>
          <span>{{ route.path }}</span>
      </span>
        <p class="text-xs my-1 text-gray-500 dark:text-slate-600">
          {{ route.description }}
        </p>
      </div>

    </router-link>
    <div>
      <PopoverMenu :items="items"/>
    </div>
  </div>
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
  {
    name: "View",
    isExternal: true,
    click: () => {
      console.log(props.route)
    }
  }
]
</script>


