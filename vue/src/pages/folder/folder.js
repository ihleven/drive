import Vue from "vue";
import App from "./App.vue";

import '@/assets/bulma-customize.scss';

//Vue.config.productionTip = false;

//new Vue({
//  render: h => h(App)
//}).$mount("#app");




new Vue({
  el: "#app",
  components: {},
  data() {
    return JSON.parse(document.getElementById('data').innerHTML);
  },
  mounted() {
    console.log("mounted", this)
  }
});