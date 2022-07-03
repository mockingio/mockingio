<template>
    <div class="hidden w-16 overflow-y-auto md:block border-gray-200 border-r dark:border-slate-800">
        <div class="w-full py-6 flex flex-col items-center">
            <div class="flex-shrink-0 flex items-center">
                <img class="h-8 w-auto" src="https://tailwindui.com/img/logos/workflow-mark.svg?color=white"
                    alt="Workflow" />
            </div>
            <div class="flex-1 mt-6 w-full px-4 space-y-1">
                <a v-for="item in sidebarNavigation" :href="item.href"
                    class="text-gray-900 group w-full p-3 rounded-md flex flex-col items-center text-xs font-medium"
                    :aria-current="item.current ? 'page' : undefined">
                    <component :is="item.icon" class="'text-gray-900 dark:text-slate-200 h-6 w-6' ]"
                        aria-hidden="true" />
                </a>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { CogIcon, CollectionIcon } from '@heroicons/vue/outline';

const sidebarNavigation = [
    { href: '#', icon: CollectionIcon, current: true },
    { href: '#', icon: CogIcon, current: false }
];
</script>
