<template>
  <b-container fluid>
    <b-alert class="mt-3" :show="error" variant="danger">Can not reach the server!</b-alert>
    <table class="table table-hover mt-3" id="poolTable">
      <thead>
        <tr>
          <th scope="col">Name</th>
          <th scope="col">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in pool" :key="calcFilename(item.name, item.extension)">
          <td>{{calcFilename(item.name, item.extension)}}</td>
          <td>
            <b-button v-if="item.used" variant="outline-danger" size="sm" disabled><font-awesome-icon icon="times"></font-awesome-icon> In use</b-button>
            <b-button v-else variant="outline-danger" size="sm" @click="deleteFile(calcFilename(item.name, item.extension))"><font-awesome-icon icon="trash"></font-awesome-icon></b-button>
          </td>
        </tr>
      </tbody>
    </table>
    <Upload/>
  </b-container>
</template>

<script>
import { GET_POOL } from '@/constants/queries'
import { REMOVE_FILE } from '@/constants/mutations'
import Upload from '@/components/Upload'

export default {
  name: 'Pool',
  data () {
    return {
      error: false
    }
  },
  methods: {
    calcFilename: (name, ext) => {
      return name + '.' + ext
    },
    deleteFile (filename) {
      this.$apollo.mutate({
        mutation: REMOVE_FILE,
        variables: {
          'file': filename
        }
      })
    }
  },
  apollo: {
    pool: {
      query: GET_POOL,
      pollInterval: 1000
    }
  },
  components: {
    Upload
  }
}
</script>
