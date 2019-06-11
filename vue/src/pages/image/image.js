import Vue from 'vue';
import store from '../drive/drive-store.js';

import Image from './Image.vue';

import '@/assets/bulma.scss';


//import(/* webpackPreload: true */ 'typeface-clear-sans/index.css');
import 'typeface-clear-sans/index.css';

store.dispatch('loadInitialData');


new Vue({
    store,
    components: {},
    data() {
        return {};
    },
    mounted() {
      console.log('mounted image.js', this);
    },
    render: h => h(Image)
  }).$mount('#app');

  
    

  