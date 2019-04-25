import Vue from 'vue';
import Album from './Album.vue';

import '@/assets/bulma-customize.scss';

const data = JSON.parse(document.getElementById('data').innerHTML);
new Vue({
    render: h =>
        h(Album, {
            props: {
                album: data.Album,
            },
        }),
}).$mount('#app');
