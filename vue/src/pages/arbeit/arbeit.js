import Vue from 'vue';
import router from './arbeit-routes';
import store from './arbeit-store';
import ArbeitApp from './ArbeitApp.vue';
import 'typeface-raleway';
import {
    Timepicker
} from 'buefy/dist/components/timepicker';
Vue.component('b-timepicker', Timepicker);
import '@/filters';

import './arbeit.scss';

new Vue({
    router,
    store,
    mounted() {
        console.log('Arbeit mounted');
        this.removeBodyTextNodes()
    },
    methods: {
        removeBodyTextNodes() {
            Array.from(document.body.childNodes).filter(n => n.nodeType==3).map(n => n.remove());
        }
    },
    render: h => h(ArbeitApp),
}).$mount('#app');