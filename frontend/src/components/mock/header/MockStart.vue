<template>
  <div class="flex items-center" v-if="activeMock">
    <PlayIcon @click="()=>startMockServer(activeMock.data.id)" class="cursor-pointer text-green-500 h-8 w-8"
              v-if="activeMock.state.status === 'stopped'"/>
    <StopIcon @click="()=>stopMockServer(activeMock.data.id)" class="cursor-pointer text-red-500 h-8 w-8"
              v-if="activeMock.state.status === 'running'"/>
  </div>
</template>

<script setup lang="ts">
import type {Mock} from "@/stores";
import {useMockStore} from "@/stores";
import {PlayIcon, StopIcon} from '@heroicons/vue/solid';

const {stopMockServer, startMockServer} = useMockStore()

defineProps({
  activeMock: {
    type: Object as () => Mock,
    required: false
  }
})

</script>
