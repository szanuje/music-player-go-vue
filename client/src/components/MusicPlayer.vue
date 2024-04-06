<script lang="ts">
export default {
  props: ['song'],
  mounted() {
    const currentSong = this.song
    if (!currentSong) {
      console.log('Missing current song')
      return
    }

    let sourceBuffer: SourceBuffer
    let mediaStream = new Uint8Array()
    let bytesRead = 0
    let chunkSize = 1000000

    const waitFor = (ms = 10) => new Promise((res) => setTimeout(() => res(null), ms))

    const mediaSource = new MediaSource()
    const musicPlayer = this.$refs.musicPlayer as HTMLAudioElement
    musicPlayer.src = window.URL.createObjectURL(mediaSource)

    const musicStreamURL = `${import.meta.env.VITE_API_URL}/stream/${currentSong}`

    mediaSource.addEventListener(
      'sourceopen',
      () => {
        sourceBuffer = mediaSource.addSourceBuffer('audio/mpeg')

        const processChunk = () => {
          const chunk = mediaStream.subarray(bytesRead, bytesRead + chunkSize)
          console.log('Processing chunk of size: %d', chunk.length)
          bytesRead = Math.min(bytesRead + chunkSize, mediaStream.length)
          sourceBuffer.appendBuffer(chunk)
          console.log('[%s] Bytes read: %d/%d', currentSong, bytesRead, mediaStream.length)
        }

        fetch(musicStreamURL)
          .then((response) => {
            if (!response.ok) {
              throw new Error('Network response was not ok')
            }
            return response.body
          })
          .then(async (body) => {
            const reader = body?.getReader()

            while (true && reader) {
              const { done, value } = await reader.read()

              if (done) {
                break
              }

              const updatedMediaStream = new Uint8Array(mediaStream.length + value.length)
              updatedMediaStream.set(mediaStream)
              updatedMediaStream.set(value, mediaStream.length)
              mediaStream = updatedMediaStream

              if (!sourceBuffer.updating && mediaStream.length > bytesRead) {
                processChunk()
              }
            }
          })
          .then(async () => {
            while (mediaStream.length > bytesRead) {
              if (!sourceBuffer.updating) {
                processChunk()
              } else {
                await waitFor(10)
              }
            }
          })
          .catch((error) => {
            console.error('Error fetching music stream:', error)
          })
      },
      false,
    )
  },
}
</script>

<template>
  <div class="player">
    <h1 class="green">Now playing</h1>
    <h3 class="green">{{ song }}</h3>
    <audio id="musicPlayer" ref="musicPlayer" controls>
      Your browser does not support the audio element.
    </audio>
  </div>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  position: relative;
  top: -10px;
}

h3 {
  font-size: 1.2rem;
}

.player {
  text-align: center;
}

@media (min-width: 1024px) {
  .player h1,
  .player h3 {
    text-align: center;
  }
}
</style>
