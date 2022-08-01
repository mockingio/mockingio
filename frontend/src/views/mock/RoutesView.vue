<template>
  <div class="flex-none w-72 border-r border-gray-200 dark:border-slate-800">
    <RouteList :activeMock="activeMock" v-if="activeMock"/>
  </div>
  <div class="flex-1 bg-gray-100 dark:bg-slate-900 overflow-auto">
    <RouteDetail :mock="activeMock" :route="activeRoute" v-if="activeRoute"/>
  </div>
</template>

<script setup lang="ts">
import {useMockStore} from "@/stores";
import {storeToRefs} from 'pinia'
import RouteList from "@/components/mock/route/RouteList.vue";
import RouteDetail from "@/components/mock/route/RouteDetail.vue";
import {useRoute} from "vue-router";
import {provide, watch} from "vue";

const {activeMock, activeRoute} = storeToRefs(useMockStore())
const {setActiveRoute, setDefaultActiveRoute} = useMockStore()

const route = useRoute()
watch(() => route.params.routeId, (newId) => {
  setRouteId(newId as string)
})

const setRouteId = (id?: string) => {
  if (id) {
    setActiveRoute(id as string)
  } else {
    setDefaultActiveRoute()
  }
}
setRouteId(route.params.routeId as string)

provide('mock', activeMock)
provide('route', activeRoute)
</script>

