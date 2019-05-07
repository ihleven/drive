import Vue from 'vue';
import store from './store';
import Album from './Album.vue';
import PhotoSwipePlugin from '@/plugins/PhotoSwipePlugin';

import '@/assets/bulma-customize.scss';

const data = JSON.parse(document.getElementById('data').innerHTML);

Vue.use(PhotoSwipePlugin);

store.commit('setAlbum', data.Album);

new Vue({
    store,

    render: h =>
        h(
            Album //, {props: {album: data.Album,},}
        ),
}).$mount('#app');
