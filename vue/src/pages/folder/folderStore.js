import Vue from 'vue';
import Vuex from 'vuex';

import axios from 'axios';
Vue.use(Vuex);

const debug = true; //process.env.NODE_ENV !== 'production';


export default new Vuex.Store({

    strict: debug,

    state: {
        folder: {},
        account: {},
        breadcrumbs: [],
    },

    mutations: {

        setData(state, data) {
            state.folder = data.Folder;
            state.account = data.Account;
            state.breadcrumbs = data.Breadcrumbs;

        },
    },
    actions: {
        loadInitialData({
            commit
        }) {

            const d = document.getElementById('data');
            if (d) {
                commit('setData', JSON.parse(d.innerHTML));
            } else {
                document.body.style.lineHeight = 0;
                console.log(document.body.style.lineHeight);
                axios.get('http://localhost:3000/' + location.hash.substring(1), {
                        headers: {
                            'Accept': 'application/json',
                        }
                    })
                    .then(function (response) {
                        // handle success
                        console.log(response.data.album, location.hash);
                        commit('setData', response.data);
                    })

            }

        }
    }
})