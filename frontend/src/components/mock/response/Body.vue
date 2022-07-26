<template>
  <div class="w-32">
    <DropdownListFilterable :selected="httpStatus" v-on="{change: change('status')}" :items="statuses"/>
  </div>
  <div>
    <textarea v-on="{input: change('body')}" class="bg-transparent border border-slate-800 w-full h-[400px]"/>
  </div>
</template>

<script setup lang="ts">
import debounce from "lodash/debounce"
import type {Response} from "@/stores";
import {computed} from "vue";
import DropdownListFilterable from "@/components/ui/DropdownListFilterable.vue";

const props = defineProps({
  response: {type: Object as () => Response, required: true}
})

const httpStatus = computed(() => {
  return props.response.status + ""
})

const change = (field: string) => debounce((evt: any) => {
  let val = ""
  if (typeof evt === 'string') {
    val = evt
  } else if (typeof evt === 'object') {
    val = evt.target.value
  }
  console.log({val})
})

const statuses = [
  {"name": "200 OK", "id": "200"},
  {"name": "201 Created", "id": "201"},
  {"name": "204 No Content", "id": "204"},
  {"name": "400 Bad Request", "id": "400"},
  {"name": "401 Unauthorized", "id": "401"},
  {"name": "403 Forbidden", "id": "403"},
  {"name": "404 Not Found", "id": "404"},
  {"name": "405 Method Not Allowed", "id": "405"},
  {"name": "406 Not Acceptable", "id": "406"},
  {"name": "409 Conflict", "id": "409"},
  {"name": "410 Gone", "id": "410"},
  {"name": "411 Length Required", "id": "411"},
  {"name": "412 Precondition Failed", "id": "412"},
  {"name": "413 Request Entity Too Large", "id": "413"},
  {"name": "414 Request-URI Too Long", "id": "414"},
  {"name": "415 Unsupported Media Type", "id": "415"},
  {"name": "416 Requested Range Not Satisfiable", "id": "416"},
  {"name": "417 Expectation Failed", "id": "417"},
  {"name": "418 I'm a teapot", "id": "418"},
  {"name": "422 Unprocessable Entity", "id": "422"},
  {"name": "423 Locked", "id": "423"},
  {"name": "424 Failed Dependency", "id": "424"},
  {"name": "425 Unordered Collection", "id": "425"},
  {"name": "426 Upgrade Required", "id": "426"},
  {"name": "449 Retry With", "id": "449"},
  {"name": "500 Internal Server Error", "id": "500"},
  {"name": "501 Not Implemented", "id": "501"},
  {"name": "502 Bad Gateway", "id": "502"},
  {"name": "503 Service Unavailable", "id": "503"},
  {"name": "504 Gateway Timeout", "id": "504"},
  {"name": "505 HTTP Version Not Supported", "id": "505"},
]
</script>