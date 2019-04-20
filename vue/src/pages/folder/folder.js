import Vue from 'vue';
import Folder from './Folder.vue';

import '@/assets/bulma-customize.scss';

//Vue.config.productionTip = false;

new Vue({
  data() {
    return {
      //folder: JSON.parse(document.getElementById('data').innerHTML),
    };
  },
  mounted() {
    console.log('mounted', this);
  },
  render: h => h(Folder),
}).$mount('#app');

// new Vue({
//     el: '#app',
//     components: {
//         Folder,
//     },
//     data() {
//         return {
//             folder: JSON.parse(document.getElementById('data').innerHTML),
//         };
//     },
//     mounted() {
//         console.log('mounted', this);
//     },
// });