import Vue from 'vue'

Vue.config.productionTip = false

new Vue({
  data: function () {
    return {
      // directory: JSON.parse(document.getElementById('data-directory').innerHTML)
    }
  },
  mounted() {
    console.log("asdfasdfasdf")
  },
  methods: {
    createThumbnails: function (dir) {}
  }
}).$mount('#vue')