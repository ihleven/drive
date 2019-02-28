import Vue from 'vue'

import Autosize from './directives/Autosize.js'


import './assets/bulma-customize.scss';

Vue.config.productionTip = false

new Vue({
  directives: {
    Autosize
  },
  data: function () {
    return {
      edit: false
      // directory: JSON.parse(document.getElementById('data-directory').innerHTML)
    }
  },
  mounted() {
    console.log("asdfasdfasdf")
  },
  methods: {
    createThumbnails: function (dir) {}
  }
}).$mount('#app')