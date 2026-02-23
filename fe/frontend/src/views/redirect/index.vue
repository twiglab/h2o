<template>
  <div></div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

onMounted(() => {
  const { params, query } = route
  const { path } = params
  router.replace({
    path: '/' + (Array.isArray(path) ? path.join('/') : path),
    query
  })
})
</script>
