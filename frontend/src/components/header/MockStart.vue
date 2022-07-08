<template>
  <div class="flex items-center" v-if="activeMock">
    <PlayIcon class="text-green-500 h-8 w-8" v-if="activeMock.state.status === 'stopped'"/>
    <StopIcon class="text-red-500 h-8 w-8" v-if="activeMock.state.status === 'running'"/>
  </div>
</template>

<script setup lang="ts">
import {PlayIcon, StopIcon} from '@heroicons/vue/solid';
import {storeToRefs} from "pinia";
import {useMockStore} from "@/stores";

const {activeMock} = storeToRefs(useMockStore())
</script>
