import Vue from 'vue';
import router from './arbeit-routes';
import store from './arbeit-store';
import ArbeitApp from './ArbeitApp.vue';
import 'typeface-raleway';
import {
    Timepicker
} from 'buefy/dist/components/timepicker';
Vue.component('b-timepicker', Timepicker);

import './arbeit.scss';

new Vue({
    router,
    store,
    render: h => h(ArbeitApp),
}).$mount('#app');