import Vue from 'vue';

//import '@/assets/bulma.scss';

//const data = JSON.parse(document.getElementById('data').innerHTML);

//import(/* webpackPreload: true */ 'typeface-clear-sans/index.css');
//import 'typeface-clear-sans/index.css';

new Vue({
    el: '#app',
    components: {},
    data() {
        return {};
    },
    mounted() {
        console.log('file =>', document.documentElement.clientWidth);
    },
});
