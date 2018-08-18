<template>
  <b-button-group>
    <b-button variant="success" @click="play()"><font-awesome-icon icon="play"></font-awesome-icon></b-button>
    <b-button variant="info" @click="pause()"><font-awesome-icon icon="pause"></font-awesome-icon></b-button>
    <b-button variant="danger" @click="stop()"><font-awesome-icon icon="stop"></font-awesome-icon></b-button>
    <b-button variant="danger" @click="removeSlot()"><font-awesome-icon icon="trash"></font-awesome-icon></b-button>
  </b-button-group>
</template>

<script>
import { REMOVE_SLOT } from '@/constants/mutations'

export default {
  name: 'SlotActions',
  props: {
    slotNumber: {
      type: Number,
      required: true
    },
    websocket: {
      type: WebSocket,
      required: true
    }
  },
  methods: {
    removeSlot () {
      this.$apollo.mutate({
        mutation: REMOVE_SLOT,
        variables: {
          'slot': this.slotNumber
        }
      }).then(() => {
        this.$apollo.queries.slots.refetch()
      })
    },
    play () {
      this.websocket.send(JSON.stringify({
        'slot': this.slotNumber,
        'player': {
          'playing': true
        }
      }))
    },
    pause () {
      this.websocket.send(JSON.stringify({
        'slot': this.slotNumber,
        'player': {
          'playing': false
        }
      }))
    },
    stop () {
      this.websocket.send(JSON.stringify({
        'slot': this.slotNumber,
        'player': {
          'stop': true
        }
      }))
    }
  }
}
</script>
