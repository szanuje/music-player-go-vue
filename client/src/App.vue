<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import MusicPlayer from './components/MusicPlayer.vue'
import { useCurrentlyPlaying } from './stores/currentlyPlaying'
import { storeToRefs } from 'pinia'

const currentlyPlaying = useCurrentlyPlaying()
const { song } = storeToRefs(currentlyPlaying)
</script>

<template>
  <header>
    <div class="wrapper">
      <MusicPlayer v-bind:song="song" :key="song" />

      <nav>
        <RouterLink to="/">Playlist</RouterLink>
        <RouterLink to="/upload">Upload song</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

header .wrapper {
  justify-content: center;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-radius: 0.375rem;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
