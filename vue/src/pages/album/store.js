import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';

const album = {
    state: {
        file: {},
        title: '',
        subtitle: '',
        description: '',
        keywords: null,
        images: null,
        image: null,
        from: null,
        until: null,
        diaryNames: [],
        diaries: {},
        sources: [],
        pages: [],
    },
    mutations: {
        setAlbum(state, album) {
            state.file = album.file;
            state.title = album.title;
            state.subtitle = album.subtitle;
            state.description = album.description;
            state.keywords = album.keywords;
            state.images = album.images;
            state.image = album.image;
            state.from = album.from;
            state.until = album.until;
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
