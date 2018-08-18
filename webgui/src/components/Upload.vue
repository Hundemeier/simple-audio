<template>
  <b-container fluid>
    <b-input-group>
      <b-form-file class="mt-3" v-model="file" :state="Boolean(file)" accept=".mp3, .wav" placeholder="Choose an audio-file to upload..."></b-form-file>
      <b-input-group-append>
        <b-button class="mt-3" @click="upload(file)">Upload</b-button>
      </b-input-group-append>
    </b-input-group>
    <b-progress class="mt-3" :value="uploadDone" :max="uploadTotal" show-progress></b-progress>

    <b-alert class="mt-3" dismissible @dismissed="error=false" fade :show="error" variant="danger">Could not upload file! <samp>{{errorText}}</samp></b-alert>
    <b-alert class="mt-3" dismissible @dismissed="success=false" fade :show="success" variant="success">Succesfully uploaded the file!</b-alert>
  </b-container>
</template>

<script>
export default {
  name: 'Upload',
  methods: {
    upload: function (file) {
      var formData = new FormData()
      formData.append('uploadfile', file)

      var xhr = new XMLHttpRequest()
      // xhr.timeout = 1000
      xhr.open('post', '/pool/upload', true)
      xhr.responseType = 'text'
      xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
          if (xhr.status >= 300) {
            this.errorText = 'Bad response from server'
            this.error = true
          } else {
            this.error = false
            this.success = true
          }
        }
      }
      xhr.onerror = () => {
        this.error = true
        this.errorText = 'An error occured!'
      }
      xhr.ontimeout = () => {
        xhr.abort()
        this.error = true
        this.errorText = 'Timeout occured!'
      }
      xhr.upload.onprogress = (e) => {
        this.uploadDone = e.position || e.loaded
        this.uploadTotal = e.totalSize || e.total
      }
      xhr.send(formData)
    }
  },
  data () {
    return {
      file: null,
      error: false,
      errorText: '',
      success: false,
      uploadDone: 0,
      uploadTotal: 10
    }
  }
}
</script>
