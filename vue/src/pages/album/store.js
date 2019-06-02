import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';


// diaryImage(state, payload) {
//     let diary = state.diaries[payload.diaryName];
//     if (payload.mode === 'remove') {
//         const index = diary.images.findIndex(img => img.name === payload.image.name);
//         console.log(index, payload.image.name);
//         diary.images.splice(index, 1);
//     } else {
//         diary.images.push(payload.image);
//     }
// },

export default new Vuex.Store({
    state: {
        baseURL: null,
        serveURL: null,
        image: {},
        images: [],
        meta: {},
        sources: []
    },
    mutations: {

        setAlbum(state, album) {
            state.images = album.images;
            if (album.images && album.images.length) {
                let index = Math.random() * album.images.length;
                state.image = album.images[Math.floor(index)];
            }
            state.baseURL = album.baseURL;
            state.serveURL = album.serveURL;
            state.meta.title = album.title;
            state.meta.subtitle = album.subtitle;
            state.sources = Object.values(album.sources);
        },
    },
    strict: debug,
});