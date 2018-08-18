<template>
  <div>
    <b-form-input v-b-tooltip.hover title="Set the volume. Currently not in dB" v-model="volumeData" type="number" max="3" step="0.1" min="-7" @change="setVolume()"></b-form-input>
  </div>
</template>

<script>
import { SET_PLAYING } from '@/constants/mutations'

export default {
  name: 'VolumeSetter',
  data () {
    return {
      volumeData: this.volume
    }
  },
  props: {
    slotNumber: {
      type: Number,
      required: true
    },
    volume: {
      type: Number,
      required: true
    }
  },
  watch: {
    volume: function (newVal, oldVal) {
      this.volumeData = newVal
    }
  },
  methods: {
    setVolume () {
      this.$apollo.mutate({
        mutation: SET_PLAYING,
        variables: {
          'slot': this.slotNumber,
          'player': {
            'volume': this.volumeData
          }
        }
      }).then(() => {
        this.$apollo.queries.slots.refetch()
      })
    }
  }
}
</script>
