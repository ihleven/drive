import Vue from 'vue';
import Vuex from 'vuex';

import axios from 'axios';
Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';

export default new Vuex.Store({
    strict: debug,

    state: {
        file: {},
        image: {},
        folder: {},
        content: '',
        account: {},
        siblings: {},
        breadcrumbs: [],
        baseURL: null,
        serveURL: null,

        images: [],
        meta: {},
        sources: [],
        type: null,
        storage: {},
    },

    mutations: {
        
        
        setData(state, data) {
            state.file = data.File;
            state.image = data.Image;
            state.content = data.Content;
            state.account = data.Account;
            state.siblings = data.Siblings;
            //state.breadcrumbs = data.Breadcrumbs;
            state.folder = data.Folder;
            state.breadcrumbs = data.Breadcrumbs;
            if (data.album) {
                state.type = 'album';
                state.images = data.album.images;
                if (data.album.images && data.album.images.length) {
                    let index = Math.random() * data.album.images.length;
                    state.image = data.album.images[Math.floor(index)];
                }
                state.baseURL = data.album.baseURL;
                state.serveURL = data.album.serveURL;
                state.meta.title = data.album.title;
                state.meta.subtitle = data.album.subtitle;
                state.sources = Object.values(data.album.sources);
            }
            state.storage = data.storage;
            console.log('data:', data);

        },
        updateContent(state, content) {
            state.content = content;
        },
    },

    actions: {
        loadInitialData({ commit }) {
            const d = document.getElementById('data');
            if (d) {
                commit('setData', JSON.parse(d.innerHTML));
            } else {
                axios
                    .get('http://localhost:3000' + location.pathname, {
                        //location.hash.substring(1), {
                        headers: {
                            Accept: 'application/json',
                        },
                    })
                    .then(function(response) {
                        
                        
                            commit('setData', response.data);
                        
                    });
            }
        },
        loadData({ commit }, payload) {
            console.log('loaddata', payload.path);
            axios
                .get('http://localhost:3000' + payload.path, {
                    //location.hash.substring(1), {
                    headers: {
                        Accept: 'application/json',
                    },
                })
                .then(function(response) {
                    console.log('loaded', response.data);
                    commit('setData', response.data);
                });
        },

        submitFileForm({ commit }, content) {
            // console.log("submitFileForm", content)
            let formData = new FormData();
            formData.set(
                'file',
                new Blob([content], {
                    type: 'text/plain',
                })
            );
            axios({
                method: 'post',
                url: location.pathname,
                data: formData,
                config: { headers: { 'Content-Type': 'multipart/form-data' } },
            })
                .then(() => {
                    commit('updateContent', content);
                })
                .catch(function(response) {
                    console.log('ERROR submitFileForm =>', response);
                });
        },
    },
});
