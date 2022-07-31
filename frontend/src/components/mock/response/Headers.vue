<template>
  <div class="divide-y-8 divide-transparent">
    <HeaderItem @change="change" v-for="item in list" :key="item.id" :item="item"/>
  </div>
  <div>
    <a href="#" class="flex justify-start my-5 items-center hover:text-green-500">
      <PlusIcon class="w-4 h-4 mr-2 text-green-500"/>
      <span class="text-sm">Add header</span>
    </a>
  </div>
</template>

<script setup lang="ts">
import HeaderItem from "./HeaderItem.vue"
import {PlusIcon} from "@heroicons/vue/outline"
import {computed} from "vue";
import type {Response} from "@/stores";

interface Item {
  id: string;
  name: string;
  value: string;
}

const emits = defineEmits(['change'])

const props = defineProps({
  response: {type: Object as () => Response, required: true},
})

const list = computed((): Item[] => {
  const headers = props.response.headers || []
  return Object.keys(headers).map((key, idx) => ({
    id: idx.toString(),
    name: key,
    value: headers[key]
  }))
})

const change = (item: any) => {
  const data = list.value.reduce((obj, cur) => {
    if (cur.id === item.id) {
      obj[item.name] = item.value
    } else {
      obj[cur.name] = cur.value
    }
    return obj
  }, {} as { [key: string]: string })

  emits("change", {headers: data})
}
</script>