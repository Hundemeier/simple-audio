<template>
  <div class="bg-dark">
    <b-navbar toggleable="sm" type="dark">
      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>

      <b-navbar-brand href="/#/">Main</b-navbar-brand>

      <b-collapse is-nav id="nav_collapse">
        <b-navbar-nav>
          <b-nav-item href="/#/pool">Pool</b-nav-item>
        </b-navbar-nav>
      </b-collapse>
      <Clock/>
    </b-navbar>
  </div>
</template>

<script>
import Clock from '@/components/Clock'

export default {
  name: 'Navbar',
  components: {
    Clock
  }
}
</script>
