// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import { ApolloClient } from 'apollo-client'
import { HttpLink } from 'apollo-link-http'
import { InMemoryCache } from 'apollo-cache-inmemory'
import VueApollo from 'vue-apollo'
import BootstrapVue from 'bootstrap-vue'
import './cyborg.min.css'
import './global.css'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faPlay, faPause, faStop, faTrash, faTimes, faCheck } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faPlay, faPause, faStop, faTrash, faTimes, faCheck)

Vue.component('font-awesome-icon', FontAwesomeIcon)

const httpLink = new HttpLink({
  uri: 'http://' + window.location.host + '/graphql'
  // uri: '/graphql'
})

const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
  connectToDevTools: true
})

const apolloProvider = new VueApollo({
  defaultClient: apolloClient,
  defaultOptions: {
    $loadingKey: 'loading'
  }
})

// store the client global for special use cases
// Vue.prototype.$apolloClient = apolloClient

Vue.use(VueApollo)
Vue.use(BootstrapVue)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  provide: apolloProvider.provide(),
  template: '<App/>'
})
