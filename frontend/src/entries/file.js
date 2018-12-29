import Vue from 'vue'
//import App from './App.vue'
//import router from './router'
//import store from './store'

Vue.config.productionTip = false

import UIkit from 'uikit'
import Icons from 'uikit/dist/js/uikit-icons'
import 'uikit/dist/css/uikit.css'

// loads the Icon plugin
UIkit.use(Icons)

// components can be called from the imported UIkit reference
UIkit.notification('Hello world.')

new Vue({
  data: function() {
    return JSON.parse(document.getElementById('data').innerHTML)
  },
  mounted() {
    console.log(this.org)
  }
  //router,
  //store,
  //render: h => h(App)
}).$mount('.main')