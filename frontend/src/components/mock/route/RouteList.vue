<template>
  <nav class="space-y-1 mt-2" v-if="activeMock">
    <RouteListItem v-for="route in activeMock.data.routes" :route="route" :mock="activeMock" :key="route.id"/>
  </nav>
</template>

<script setup lang="ts">
import type {Mock} from "@/stores";
import RouteListItem from "./RouteListItem.vue";

defineProps({
  activeMock: {type: Object as () => Mock, required: true}
})
</script>
