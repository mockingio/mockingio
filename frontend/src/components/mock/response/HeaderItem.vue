<template>
  <div class="grid grid-cols-2 gap-4">
    <div class="col-span-1 flex items-center">
      <input type="text" name="name" :value="item.name" v-on="{input: change('name')}"
             class="dark:bg-transparent dark:border-slate-800 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"/>
    </div>

    <div class="col-span-1 flex items-center">
      <input type="text" name="value" :value="item.value" v-on="{input: change('value')}"
             class=" dark:bg-transparent dark:border-slate-800 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"/>
      <TrashIcon class="w-5 h-5 ml-5"/>
    </div>
  </div>
</template>

<script setup lang="ts">
import {TrashIcon} from '@heroicons/vue/outline'
import {defineProps} from "vue";

const props = defineProps({
  item: {type: Object as () => { id: string, name: string, value: string }, required: true},
})

const emits = defineEmits(['change'])

const change = (field: string) => (evt: any) => {
  const data = {...props.item, [field]: evt.target.value}
  emits("change", data)
}
</script>
