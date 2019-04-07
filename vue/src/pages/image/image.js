import Vue from 'vue';

import '@/assets/bulma-customize.scss';

const data = JSON.parse(document.getElementById('data').innerHTML);

//import(/* webpackPreload: true */ 'typeface-clear-sans/index.css');
import 'typeface-clear-sans/index.css';

new Vue({
    el: '#app',
    components: {},
    data() {
        return {
            data,
        };
    },
    mounted() {
        console.log('file =>', this);
    },
});
