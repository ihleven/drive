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
        content: '',
        account: {},
        breadcrumbs: [],
    },

    mutations: {
        setData(state, data) {
            state.file = data.File;
            state.image = data.Image;
            state.content = data.Content;
            state.account = data.User;
            //state.breadcrumbs = data.Breadcrumbs;
            console.log("data:", data)
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
                    .get('http://localhost:3000/' + location.hash.substring(1), {
                        headers: {
                            Accept: 'application/json',
                        },
                    })
                    .then(function(response) {
                        commit('setData', response.data);
                    });
            }
        },
    },
});
