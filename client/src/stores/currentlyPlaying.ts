import { defineStore } from 'pinia'

export const useCurrentlyPlaying = defineStore('currentlyPlaying', {
  state: () => {
    return { song: '' }
  },
  actions: {
    set(song: string) {
      console.log('Changing currently playing song to: %s', song)
      this.song = song
    },
  },
})
