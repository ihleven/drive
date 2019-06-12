import Vue from 'vue';
import VueRouter from 'vue-router';
import store from './drive-store.js';
import File from './File.vue';
import Image from '../image/Image.vue';
import { mapState } from 'vuex';

import './drive-styles.scss';

import 'typeface-clear-sans/index.css';

Vue.use(VueRouter);

store.dispatch('loadInitialData');

//  const comp = {template: '<div><component :is="$route.params.tab"></component></div>'}
console.log('Vue.config.errorHandler', Vue.config.errorHandler);
const router = new VueRouter({
    mode: 'history',
    routes: [
        {
            path: '/bar',
            component: { template: '<div>bar</div>' },
        },
        {
            path: '*',
            component: File,
        },
    ],
});
router.beforeEach((to, from, next) => {
        store.dispatch("loadData", to)
        next()
  })
new Vue({
    el: '#app',
    router,
    store,
    computed: {
        ...mapState(['account', 'file', 'content']),
        ViewComponent() {
            if (this.file.mime && this.file.mime.Type == 'image') {
                return Image;
            }
            return File;
        },
    },
    
    mounted() {
        console.log('drive mounted');
    },
    render (h) { return h(this.ViewComponent) }
});
