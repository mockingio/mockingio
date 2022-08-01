<template>
  <Disclosure as="nav">
    <div class="h-20 flex border-b border-gray-200 dark:border-slate-800 text-sm dark:text-sm">
      <div class="flex flex-1 items-center">
        <div class="flex w-72">
          <MockStart :activeMock="activeMock" class="mx-4"/>
          <MockSelect :mocks="mocks" :activeMock="activeMock" class="flex-1"/>
        </div>
        <div class="flex flex-1 justify-between">
          <div class="ml-5">
            <router-link v-for="item in menuItems" active-class="text-green-500" class="pr-3 py-2 hover:text-green-500"
                         :to="{name: item.route, params: {id: activeMock.data.id}}">
              {{ item.name }}
            </router-link>
          </div>
          <div class="mr-5 flex">
            <a href="#">
              <DuplicateIcon
                  class="text-gray-500 dark:text-slate-400 dark:hover:text-green-500 hover:text-green-500 h-5 w-5 mr-5"/>
            </a> <a href="#">
            <ShareIcon
                class="text-gray-500 dark:text-slate-400 dark:hover:text-green-500 hover:text-green-500 h-5 w-5 mr-5"/>
          </a>
            <a href="#">
              <SaveIcon
                  class="text-gray-500 dark:text-slate-400 hover:text-green-500 dark:hover:text-green-500 h-5 w-5"/>
            </a>
          </div>
        </div>
      </div>
    </div>
  </Disclosure>
</template>

<script setup lang="ts">
import {Disclosure} from '@headlessui/vue';
import {DuplicateIcon, SaveIcon, ShareIcon} from '@heroicons/vue/outline';
import MockSelect from './MockSelect.vue';
import MockStart from './MockStart.vue';

const menuItems: { name: string, route: string }[] = [
  {
    route: 'routes-view',
    name: 'Routes',
  },
  {
    route: 'route-proxy-view',
    name: 'Proxy',
  },
  {
    route: 'route-log-view',
    name: 'Logs',
  },
  {
    route: 'route-settings-view',
    name: 'Settings',
  }
];</script>


<script lang="ts">
import type {Mock} from "@/stores";

export default {
  props: {
    mocks: {type: Object as () => Mock[], required: true},
    activeMock: {type: Object as () => Mock, required: false}
  }
};
</script>