<template>
  <HeaderItem v-for="(name, value) in response.headers" :key="name" :name="name" :value="value"/>
</template>

<script setup lang="ts">
import {defineProps} from "vue";
import HeaderItem from "./HeaderItem.vue";

const props = defineProps({
  response: {type: Object as () => Response, required: true},
})
</script>