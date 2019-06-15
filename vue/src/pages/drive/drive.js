import Vue from 'vue';
import VueRouter from 'vue-router';
import store from './drive-store.js';
import File from './File.vue';
import Folder from './Folder.vue';
import PhotoAlbum from './PhotoAlbum.vue';
import Image from './Image.vue';
import { mapState } from 'vuex';

import './drive-styles.scss';

import 'typeface-clear-sans/index.css';




Vue.use(VueRouter);

store.dispatch('loadInitialData');

const Drive = {
    computed: {
        ...mapState(['account', 'file', 'content', 'folder', 'type']),
    },
    mounted() {
        console.log('Drive mounted');
        this.removeBodyTextNodes()
    },
    methods: {
        removeBodyTextNodes() {
            Array.from(document.body.childNodes).filter(n => n.nodeType==3).map(n => n.remove());
        }
    },
    render(h) {
        let c = null;
        // if (this.type=='album') {
        //     c = Album;
        // }
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
            path: '/alben*',
            component: PhotoAlbum,
        },
        {
            path: '*',
            component: Drive,
        },
    ],
});

new Vue({
    el: '#app',
    router,
    store,
    mounted() {
        this.$router.beforeEach((to, from, next) => {
            store.dispatch('loadData', to);
            next();
        });
    },
    template: '<router-view></router-view>',
});

//render: h => h(Image)
