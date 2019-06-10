import Vue from 'vue';
import store from './album-store';
import Album from './Album.vue';
import PhotoSwipePlugin from '@/plugins/PhotoSwipePlugin';
import axios from 'axios';

import './album-bulma.scss';

Vue.use(PhotoSwipePlugin);

const d = document.getElementById('data');
let data = null;
if (d) {
    data = JSON.parse(d.innerHTML);
    //store.commit('SET_FOLDER', data.folder);
    //store.commit('SET_ACCOUNT', data.account);
    store.commit('setAlbum', data.album);

} else {
    axios.get('http://localhost:3000/alben/' + location.hash.substring(1), {
            headers: {
                'Accept': 'application/json',
            }
        })
        .then(function (response) {
            // handle success
            console.log(response.data.album, location.hash);
            store.commit('setAlbum', response.data.album);
        })
    data = {}
}



new Vue({
    store,
    created() {
        console.log("album.js")
    },
    render: h =>
        h(
            Album //, {props: {album: data.Album,},}
        ),
}).$mount('#app');