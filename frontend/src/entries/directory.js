import Vue from 'vue'
import axios from 'axios'
import UIkit from 'uikit'
import Icons from 'uikit/dist/js/uikit-icons'

UIkit.use(Icons)

Vue.config.productionTip = false

new Vue({
  data: function () {

    return {
      directory: JSON.parse(document.getElementById('data-directory').innerHTML)
    }
  },
  mounted() {
    this.createThumbnails();
    console.log("asdfasdfasdf")
  },
  methods: {
    createThumbnails: function (dir) {
      console.log("createThumbnails:", dir);
      axios.post(dir, {
          CreateThumbnails: true,
        })
        .then(function (response) {
          console.log(response);
        })
        .catch(function (error) {
          console.log(error);
        });
    }
  }
}).$mount('#vueapp')