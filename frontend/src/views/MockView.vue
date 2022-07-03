<template>
  <div class="min-h-ful flex">
    <MainNav />
    <main class="flex-1">
      <Header />
      <div class="flex min-h-screen">
        <div class="flex-none w-72 border-r border-gray-200 dark:border-slate-800">
          <Routes />
        </div>
        <div class="flex-1 bg-gray-100 dark:bg-slate-900 overflow-auto">
          <div class="m-5">
            <div class="mt-1 flex rounded-md shadow-sm">
              <span
                class="inline-flex dark:bg-slate-800 dark:border-slate-700 items-center px-3 rounded-l-md border border-r-0 border-gray-300 bg-gray-50 text-gray-500 dark:text-slate-400 sm:text-sm">
                GET /
              </span>
              <input value="hello/world" type="text" name="username" id="username" autocomplete="username"
                class="flex-1 dark:bg-slate-800 block dark:border-slate-700 w-full min-w-0 rounded-none rounded-r-md sm:text-sm border-gray-300" />
            </div>
          </div>

          <h3 class="text-md m-5 text-gray-900">Responses</h3>

          <div class="m-5">
            <Responses />
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import Header from '../components/header/Header.vue';
import MainNav from '../components/MainNav.vue';
import Responses from '../components/response/Responses.vue';
import Routes from '../components/routes/Routes.vue';
</script>
