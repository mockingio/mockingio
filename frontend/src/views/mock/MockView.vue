<template>
  <Header :mocks="mocks" :activeMock="activeMock" v-if="activeMock"/>
  <div class="flex min-h-screen">
    <RouterView/>
  </div>
</template>

<script setup lang="ts">
import {RouterView, useRoute} from 'vue-router';
import {useMockStore} from "@/stores";
import {watch} from "vue";
import {storeToRefs} from 'pinia'
import Header from '@/components/mock/header/Header.vue';

const route = useRoute()
const {fetchMocks, setActiveMock, setDefaultActiveMock} = useMockStore()
const {mocks, error, activeMock} = storeToRefs(useMockStore())

watch(() => route.params.id, (newId) => {
  if (newId) {
    setActiveMock(newId as string)
  } else {
    setDefaultActiveMock()
  }
})

fetchMocks().then(() => {
  if (route.params.id) {
    setActiveMock(route.params.id as string)
  } else {
    setDefaultActiveMock()
  }
})
</script>
