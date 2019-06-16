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
        error: null,
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
        error(state, error) {
            state.error = error;
        }
    },

    actions: {
        loadInitialData({ commit }) {
            const d = document.getElementById('data');
            if (d) {
                commit('setData', JSON.parse(d.innerHTML));
            } else {
                axios
                    .get(location.pathname, {
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
            console.log('loaddata', payload);
            commit('error', null);
            axios
                .get(payload.to.path, {
                    //location.hash.substring(1), {
                    headers: {
                        Accept: 'application/json',
                    },
                })
                .then(function(response) {
                    console.log('loaded', response.status);
                    commit('setData', response.data);
                    payload.next();
                })
                .catch(error => {
                    // Error 😨
                    if (error.response) {
                        /*
                         * The request was made and the server responded with a
                         * status code that falls out of the range of 2xx
                         */
                        //console.log(error.response.data);
                        //console.log(error.response.status);
                        //console.log(error.response.headers);
                        commit('error', error.response);
                    } else if (error.request) {
                        /*
                         * The request was made but no response was received, `error.request`
                         * is an instance of XMLHttpRequest in the browser and an instance
                         * of http.ClientRequest in Node.js
                         */
                        console.log(error.request);
                    } else {
                        // Something happened in setting up the request and triggered an Error
                        console.log('Error', error.message);
                    }
                    console.log(error.config);
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
