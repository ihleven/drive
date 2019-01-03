import Vue from 'vue'
import HelloWorld from '../components/HelloWorld.vue'
import Breadcrumbs from '../components/Breadcrumbs.vue'


import autosize from '../directives/Autosize.js'

Vue.config.productionTip = false

import UIkit from 'uikit'
import Icons from 'uikit/dist/js/uikit-icons'
//import 'uikit/dist/css/uikit.css'

// loads the Icon plugin
UIkit.use(Icons)

// components can be called from the imported UIkit reference
UIkit.notification('Hello world.')

new Vue({
  data: function () {
    return JSON.parse(document.getElementById('data').innerHTML)
  },
  mounted() {
    console.log('mounted', this.org)
  },
  components: {
    HelloWorld,
    Breadcrumbs
  },

  directives: {
    autosize: autosize
  }
  //router,
  //store,
  //render: h => h(App)
}).$mount('#vueapp')