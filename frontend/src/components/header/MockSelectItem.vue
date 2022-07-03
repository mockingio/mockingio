<template>
    <a href="#">
        <div>
            <p class="text-base flex items-center font-medium text-gray-900 dark:text-slate-400">
                <PlayIcon class="text-green-500 h-4 w-4 mr-2" />
                <span>{{ name }}</span>
            </p>
            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400 dark:text-slate-600">
                {{ description }}
            </p>
        </div>
    </a>
</template>

<script setup lang="ts">
import { PlayIcon } from '@heroicons/vue/solid';
</script>

<script lang="ts">
export default {
    props: {
        name: { type: String, required: true },
        description: { type: String, required: true }
    }
};
</script>
