<template>
  <div class="min-h-full flex">
    <MainNav/>
    <main class="flex-1">
      <RouterView/>
    </main>
  </div>
</template>

<style>
@import '@/assets/base.css';
</style>

<script setup lang="ts">
import {RouterView} from 'vue-router';
import MainNav from './components/MainNav.vue';</script>