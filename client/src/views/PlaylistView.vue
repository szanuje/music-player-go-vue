<script setup lang="ts">
import { useFetch } from '@vueuse/core'
import { useCurrentlyPlaying } from '../stores/currentlyPlaying'

const currentlyPlaying = useCurrentlyPlaying()

const { isFetching, error, data } = useFetch<{ title: string }[]>(
  `${import.meta.env.VITE_API_URL}/songs`,
)
  .get()
  .json()
</script>

<template>
  <main>
    <h1>Playlist</h1>
    <div v-if="isFetching" class="loading">Loading...</div>

    <div v-if="error" class="error">{{ error }}</div>

    <ol v-if="data" class="content">
      <li v-for="{ title } in data" :key="title">
        <button @click="currentlyPlaying.set(title)">
          {{ title }}
        </button>
      </li>
    </ol>
  </main>
</template>
