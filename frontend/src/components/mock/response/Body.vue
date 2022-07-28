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
import type {Ref} from "vue";
import {computed, inject} from "vue";
import type {Mock, Response, Route} from "@/stores";
import {useMockStore} from "@/stores";
import DropdownListFilterable from "@/components/ui/DropdownListFilterable.vue";
import {getStatuses} from "@/helpers";

const {patchResponse} = useMockStore()

const props = defineProps({
  response: {type: Object as () => Response, required: true}
})

const httpStatus = computed(() => {
  return props.response.status.toString()
})

const mock = inject<Ref<Mock>>("mock")
const route = inject<Ref<Route>>("route")

const change = (field: string) => (evt: any) => {
  let val = 0
  if (typeof evt === 'string') {
    val = parseInt(evt)
  } else if (typeof evt === 'object') {
    val = parseInt(evt.target.value)
  }
  patchResponse(mock!.value.data.id, route!.value.id, props.response.id, {[field]: val})
}

const statuses = getStatuses()
</script>