<template>
  <div class="m-5">
    <div class="mt-1 flex rounded-md">
      <div class="w-32">
        <DropdownListFilterable :selected="method" @change="change" :items="items"/>
      </div>
      <input :value="route.path" type="text" name="path" id="path"
             class="flex-1 bg-transparent block w-full min-w-0 rounded-none rounded-r-md sm:text-sm border-gray-300 dark:border-slate-800"/>
    </div>
  </div>

  <div class="m-5">
    <Responses :mock="mock" :route="route"/>
  </div>
</template>

<script setup lang="ts">
import type {Mock, Route} from "@/stores";
import {useMockStore} from "@/stores";
import Responses from "@/components/mock/response/Responses.vue";
import DropdownListFilterable from '@/components/ui/DropdownListFilterable.vue'
import type {Item} from "@/components/ui/PopoverMenu.vue";
import {computed} from "vue";

const {patchRoute} = useMockStore()

const props = defineProps({
  mock: {type: Object as () => Mock, required: true},
  route: {type: Object as () => Route, required: true}
})

let items = [
  {name: 'GET'},
  {name: 'POST'},
  {name: 'PUT'},
  {name: 'PATCH'},
  {name: 'DELETE'},
  {name: 'OPTIONS'},
]
items = items.map(i => ({...i, id: i.name}))
const method = computed(() => {
  return items.find(i => i.name === props.route.method)
})

const change = (value: Item) => {
  patchRoute(props.mock.data.id, props.route.id, {method: value.name})
}
</script>
