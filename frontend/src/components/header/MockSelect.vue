<template>
    <div class="flex">
        <div class="flex-1">
            <h2 class="font-semibold">Twitter API</h2>
            <p class="text-xs text-gray-500 dark:text-slate-400">127.0.0.1:3000</p>
        </div>
        <Popover class="relative flex">
            <PopoverButton class="text-white items-center text-base font-medium focus:outline-none">
                <DotsVerticalIcon class="text-gray-900 dark:text-slate-200 w-6 h-6" />
            </PopoverButton>

            <PopoverPanel class="absolute z-10 left-[30px] transform mt-3 px-2 w-screen max-w-md sm:px-0">
                <div class="rounded-lg shadow-lg ring-1 ring-black dark:ring-slate-800 ring-opacity-5 overflow-hidden">
                    <div class="relative grid bg-white dark:bg-slate-900 p-3">
                        <h2 class="ml-4 text-lg font-bold">Mocks</h2>
                        <div class="max-h-96 overflow-auto">
                            <MockSelectItem v-for="item in resources" :key="item.name" :name="item.name"
                                :description="item.description"
                                class="p-4 flex items-start transition ease-in-out duration-150 hover:border-green-500 border-transparent border-l-2">
                            </MockSelectItem>
                        </div>
                        <div class="mx-5 mt-5 flex justify-end">
                            <button type="button"
                                class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-white bg-green-500 focus:outline-none">
                                New mock
                            </button>
                        </div>
                    </div>
                </div>
            </PopoverPanel>
        </Popover>
    </div>
</template>

<script setup lang="ts">
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { DotsVerticalIcon } from '@heroicons/vue/solid';
import MockSelectItem from './MockSelectItem.vue';

const resources = [
    {
        name: 'Google API',
        description:
            'Get all of your questions answered in our forums or contact support.',
        href: '#'
    },
    {
        name: 'Twitter API',
        description:
            'Learn how to maximize our platform to get the most out of it.',
        href: '#'
    },
    {
        name: 'AWS API',
        description:
            'See what meet-ups and other events we might be planning near you.',
        href: '#'
    },
    {
        name: 'Salesforce API',
        description: 'Understand how we take your privacy seriously.',
        href: '#'
    }
];
const recentPosts = [
    { id: 1, name: 'Boost your conversion rate', href: '#' },
    {
        id: 2,
        name: 'How to use search engine optimization to drive traffic to your site',
        href: '#'
    },
    { id: 3, name: 'Improve your customer experience', href: '#' }
];
</script>
