<template>
    <div class="bg-white dark:bg-slate-800 shadow overflow-hidden sm:rounded-md">
        <div role="list" class="text-sm">
            <draggable handle=".handle"
                class="dragArea block list-group w-full divide-y divide-solid dark:divide-slate-700" :list="list">
                <Response />
                <Response />
                <Response />
                <div>
                    <div class="flex">
                        <div class="m-5 block flex-1 flex justify-between">
                            <p>200 OK - When data is correct</p>
                            <ViewListIcon class="handle w-5 h-5 ml-5" />
                        </div>
                    </div>
                    <div>
                        <div class="hidden sm:block">
                            <div class="border-b border-gray-200 dark:border-slate-700 mx-5">
                                <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                                    <a v-for="tab in tabs" :key="tab.name" :href="tab.href" :class="[
                                        tab.current
                                            ? 'border-green-600'
                                            : 'border-transparent text-gray-500 dark:text-slate-400',
                                        'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm hover:border-green-600'
                                    ]" :aria-current="tab.current ? 'page' : undefined">
                                        {{ tab.name }}
                                    </a>
                                </nav>
                            </div>
                        </div>
                        <div class="m-5">
                            <Rules />
                        </div>
                        <div class="m-5 flex justify-end">
                            <button type="button"
                                class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-red-500 border-red-500 hover:text-white hover:bg-red-500">
                                Delete
                            </button>
                        </div>
                    </div>
                </div>
            </draggable>
        </div>
    </div>
    <div>
        <a href="#" class="flex justify-start my-5 items-center hover:text-green-500">
            <PlusIcon class="w-4 h-4 mr-2 text-green-500" />
            <span class="text-sm">Add new response</span>
        </a>
    </div>
</template>

<script setup lang="ts">
import Rules from '../rule/Rules.vue';
import { ViewListIcon, PlusIcon } from '@heroicons/vue/outline';
import Response from './Response.vue';

const tabs = [
    { name: 'Body', href: '#', current: false },
    { name: 'Headers', href: '#', current: false },
    { name: 'Rules', href: '#', current: true }
];
</script>

<script lang="ts">
import { defineComponent } from 'vue';
import { VueDraggableNext } from 'vue-draggable-next';

export default defineComponent({
    components: {
        draggable: VueDraggableNext
    },
    data() {
        return {
            enabled: true,
            list: [
                { name: 'John', id: 1 },
                { name: 'Joao', id: 2 },
                { name: 'Jean', id: 3 },
                { name: 'Gerard', id: 4 }
            ],
            dragging: false
        };
    }
});
</script>
