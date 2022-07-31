<template>
  <div>
    <div @click="toggleOpen" class="select-none flex dark:bg-slate-800 bg-white cursor-pointer handle group">
      <div class="m-3 block flex-1 flex justify-between">
        <div>{{ status }}</div>
        <TrashIcon class="w-5 h-5 ml-5 hover:text-red-500 invisible group-hover:visible"/>
      </div>
    </div>

    <div :class="open ? 'block' : 'hidden'">
      <div class="dark:border-slate-800 border border-t-0 p-3 pt-0">
        <TabGroup>
          <div class="border-b border-gray-200 dark:border-slate-800">
            <TabList class="-mb-px flex space-x-8">
              <Tab v-for="item in tabs" :key="item.name" v-slot="{ selected }" as="template">
                <button
                    :class="[selected ? 'border-green-600': 'border-transparent', 'focus:outline-none whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm hover:border-green-600']">
                  {{ item.name }}
                </button>
              </Tab>
            </TabList>
          </div>
          <TabPanels>
            <TabPanel v-for="item in tabs" :key="item.name">
              <div class="my-5">
                <component :is="item.component" :response="props.response" @change="change"/>
              </div>
            </TabPanel>
          </TabPanels>
        </TabGroup>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, inject, ref, Ref} from "vue";
import {TrashIcon} from '@heroicons/vue/outline';
import {Tab, TabGroup, TabList, TabPanel, TabPanels} from '@headlessui/vue'
import Body from "@/components/mock/response/Body.vue";
import Headers from "@/components/mock/response/Headers.vue";
import Rules from "./Rules.vue";
import {getStatusById} from "@/helpers";
import type {Mock, Response, Route} from "@/stores";
import {useMockStore} from "@/stores";

const tabs = [
  {
    name: 'Status & Body',
    component: Body,
  },
  {
    name: 'Headers',
    component: Headers,
  },
  {
    name: 'Rules',
    component: Rules,
  }
];

const props = defineProps({
  response: {type: Object as () => Response, required: true},
})

const status = computed(() => getStatusById(props.response.status.toString())?.name)

const open = ref(false)

function toggleOpen() {
  open.value = !open.value
}

const mock = inject<Ref<Mock>>("mock")
const route = inject<Ref<Route>>("route")
const {patchResponse} = useMockStore()

const change = (data: any) => {
  patchResponse(mock!.value.data.id, route!.value.id, props.response.id, data)
}

</script>

