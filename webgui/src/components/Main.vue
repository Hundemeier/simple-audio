<template>
  <b-container fluid>
    <b-alert class="mt-3" :show="error" variant="danger">Can not reach the server!</b-alert>
    <table class="table table-hover mt-3" id="slotTable">
      <thead>
        <tr>
          <th scope="col" style="width:5rem">Slot</th>
          <th scope="col">Name</th>
          <th scope="col" style="width:5rem">Looping</th>
          <th scope="col" style="width:7rem">Volume</th>
          <th scope="col">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="slot in slots" :key="slot.slot" :class="{'table-success': slot.player.playing}">
          <td>{{slot.slot}}</td>
          <td>{{slot.item.name + '.' + slot.item.extension}}</td>
          <td>
            <b-button v-if="slot.player.loop" variant="outline-info" @click="setLoop(slot.slot,false)"><font-awesome-icon icon="check"></font-awesome-icon></b-button>
            <b-button v-else variant="outline-secondary" @click="setLoop(slot.slot,true)"><font-awesome-icon icon="times"></font-awesome-icon></b-button>
          </td>
          <td><volume-setter :slotNumber="slot.slot" :volume="slot.player.volume"/></td>
          <td><slot-actions :slotNumber="slot.slot" :websocket="ws" /></td>
        </tr>
      </tbody>
    </table>
    <SetSlot/>
  </b-container>
</template>

<script>
import { GET_SLOTS } from '@/constants/queries'
import { SET_PLAYING } from '@/constants/mutations'
import SetSlot from '@/components/SetSlot'
import SlotActions from '@/components/SlotActions'
import VolumeSetter from '@/components/VolumeSetter'

export default {
  name: 'Main',
  components: {
    SetSlot,
    SlotActions,
    VolumeSetter
  },
  data () {
    return {
      slots: [],
      error: false,
      ws: null
    }
  },
  apollo: {
    slots: {
      query: GET_SLOTS,
      pollInterval: 2000, // higher interval, because reloads are triggered via websocket
      manual: true,
      result ({ data, loading }) {
        if (!loading) {
          if (data.slots === undefined) {
            this.error = true
            return
          }
          // check if error was true to trigger a reload
          if (this.error) {
            location.reload()
          }
          this.error = false
          // manually sort the array
          var arr = data.slots.slice()
          this.slots = arr.sort((a, b) => (a.slot - b.slot))
        }
      }
    }
  },
  methods: {
    setLoop (slot, looping) {
      this.$apollo.mutate({
        mutation: SET_PLAYING,
        variables: {
          'slot': slot,
          'player': {
            'loop': looping
          }
        }
      }).then(() => {
        this.$apollo.queries.slots.refetch()
      })
    }
  },
  mounted: function () {
    const addr = 'ws://' + window.location.host + '/websocket'
    console.log('start websocket')
    this.ws = new WebSocket(addr)

    this.ws.onerror = function () {
      console.log('Error on websocket')
    }// .bind(this)

    this.ws.onopen = function () {
      console.log('websocket was opened!')
    }// .bind(this)

    this.ws.onclose = function () {
      console.log('websocket was closed! Try to reconnect...')
      this.ws = new WebSocket(addr)
    }// .bind(this)

    this.ws.onmessage = function (e) {
      // console.log(JSON.parse(e.data))
      // force refetch. Not ideal solution but currently the only one
      this.$apollo.queries.slots.refetch()
    }.bind(this)
  }
}
</script>
