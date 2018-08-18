<template>
  <b-card header="Set Slot-Mapping">
    <b-form inline>
      <b-form-input style="width:7rem" class="m-2" v-model="slot" type="number" max="512" min="1" required></b-form-input>
      <b-form-select class="m-2" v-model="selectedName" :options="options" required/>
      <b-button class="m-2" type="submit" variant="secondary" @click="setMapping()">Set</b-button>
    </b-form>
  </b-card>
</template>

<script>
import { GET_POOL } from '@/constants/queries'
import { SET_MAPPING } from '@/constants/mutations'

export default {
  name: 'SetSlot',
  data () {
    return {
      slot: null,
      selectedName: null,
      options: []
    }
  },
  methods: {
    calcFilename: (name, ext) => {
      return name + '.' + ext
    },
    setMapping () {
      if (this.slot == null || this.selectedName == null) {
        return
      }
      this.$apollo.mutate({
        mutation: SET_MAPPING,
        variables: {
          'slot': parseInt(this.slot),
          'name': this.selectedName
        }
      }).then(() => {
        this.$apollo.queries.slots.refetch()
      })
    }
  },
  apollo: {
    pool: {
      query: GET_POOL,
      pollInterval: 2000,
      manual: true,
      result ({ data, loading }) {
        if (!loading) {
          this.options = []
          var item
          for (item of data.pool) {
            this.options.push({
              value: this.calcFilename(item.name, item.extension),
              text: this.calcFilename(item.name, item.extension)
            })
          }
        }
      }
    }
  }
}
</script>
