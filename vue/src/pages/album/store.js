import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';

const album = {
    state: {
        file: {},
        diaryNames: [],
        diaries: {},
        sources: [],
        pages: [],
    },
    mutations: {
        setAlbum(state, album) {
            state.file = album.file;
            state.sources = album.sources;
            state.pages = album.pages;
            console.log(state.pages);
            album.diaries.forEach(diary => {
                state.diaries[diary.name] = diary;
                state.diaryNames.push(diary.name);
            });
        },
        diaryImage(state, payload) {
            let diary = state.diaries[payload.diaryName];
            if (payload.mode === 'remove') {
                const index = diary.images.findIndex(img => img.name === payload.image.name);
                console.log(index, payload.image.name);
                diary.images.splice(index, 1);
            } else {
                diary.images.push(payload.image);
            }
        },
    },
};

export default new Vuex.Store({
    modules: {
        album: album,
    },

    strict: debug,
});
