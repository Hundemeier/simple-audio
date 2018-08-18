<template>
  <div>
    {{time}}
  </div>
</template>

<script>
export default {
  name: 'Clock',
  data () {
    setInterval(this.setTime, 1000)
    return {
      time: new Date().toLocaleTimeString()
    }
  },
  methods: {
    setTime () {
      this.time = new Date().toLocaleTimeString()
    }
  }
}
</script>
