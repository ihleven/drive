import Vue from 'vue';
import VueRouter from 'vue-router';
import store from './drive-store.js';
import File from './File.vue';
import Image from './Image.vue';
import Folder from './Folder.vue';
import { mapState } from 'vuex';

import './drive-styles.scss';

import 'typeface-clear-sans/index.css';

Vue.use(VueRouter);

store.dispatch('loadInitialData');

const Drive = {
    computed: {
        ...mapState(['account', 'file', 'content', 'folder']),
    },
    mounted() {
        console.log('Drive mounted');
    },
    render(h) {
        let c = null;
        if (this.folder) {
            c = Folder;
        }
        if (this.file && this.file.mime) {
            if (this.file.mime.Type == 'image') {
                c = Image;
            } else {
                c = File;
            }
        }
        return h(c);
    },
};

const router = new VueRouter({
    mode: 'history',
    routes: [
        {
            path: '/bar',
            component: { template: '<div>bar</div>' },
        },
        {
            path: '*',
            component: Drive,
        },
    ],
});

router.beforeEach((to, from, next) => {
    store.dispatch('loadData', to);
    next();
});

new Vue({
    el: '#app',
    router,
    store,
    template: '<router-view></router-view>',
});


    //render: h => h(Image)
