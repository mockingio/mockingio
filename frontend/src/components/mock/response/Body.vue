<template>
  <div class="w-32">
    <DropdownListFilterable :selected="httpStatus" v-on="{change: change('status')}" :items="statuses"/>
  </div>
  <div>
    <textarea :value="response.body" v-on="{input: change('body')}"
              class="bg-transparent border border-slate-800 w-full h-[400px]"/>
  </div>
</template>

<script setup lang="ts">
import {computed} from "vue";
import type {Response} from "@/stores";
import DropdownListFilterable from "@/components/ui/DropdownListFilterable.vue";
import {getStatuses} from "@/helpers";

const props = defineProps({
  response: {type: Object as () => Response, required: true}
})

const emits = defineEmits(['change'])

const httpStatus = computed(() => {
  return props.response.status.toString()
})

const change = (field: string) => (evt: any) => {
  let val = 0
  if (typeof evt === 'string') {
    val = parseInt(evt)
  } else if (typeof evt === 'object') {
    val = parseInt(evt.target.value)
  }
  emits("change", {[field]: val})
}

const statuses = getStatuses()
</script>