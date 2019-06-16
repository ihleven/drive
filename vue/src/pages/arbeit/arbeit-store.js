import Vue from 'vue';
import Vuex from 'vuex';

import axios from 'axios';
Vue.use(Vuex);

const debug = true; //process.env.NODE_ENV !== 'production';

const SET_ARBEITSJAHR = 'SET_ARBEITSJAHR';
const SET_ARBEITSTAG = 'SET_ARBEITSTAG';
//const SET_DATE = 'SET_DATE';
const SET_LOADING = 'SET_LOADING';
const SET_ACTIVE = 'SET_ACTIVE';
const SET_ACTIVE_DATE = 'SET_ACTIVE_DATE';
const MONATE = {
    0: 'Jänner',
    1: 'Februar',
    2: 'Maerz',
    3: 'April',
    4: 'Mai',
    5: 'Juni',
    6: 'Juli',
    7: 'August',
    8: 'September',
    9: 'Oktober',
    10: 'November',
    11: 'Dezember',
};
const WOCHENTAGE = {
    0: 'Sunndig',
    1: 'Mändig',
    2: 'Zischdig',
    3: 'Mittwoch',
    4: 'Dunschdig',
    5: 'Friddig',
    6: 'Samschdig',
};

export default new Vuex.Store({
    strict: debug,

    state: {
        activeDate: {},
        arbeitsjahr: {},
        arbeitstag: {},
        date: null,
        year: null,
        month: null,
        day: null,
        loading: false,
    },
    getters: {
        monat: state => {
            return state.month > 0 && state.month <= 12 ? MONATE[state.month] : '';
        },
        wochentag: state => {
            let day = state.date ? state.date.getDate() : -1;
            return day >= 0 && day <= 6 ? WOCHENTAGE[day] : '';
        },
    },
    mutations: {
        [SET_ACTIVE](state, payload) {
            console.log('COMMIT_ACTIVE', payload);
            let d = payload.day ? new Date(payload.year, payload.month - 1, payload.day) : null;
            state.activeDate = {
                date: d,
                year: payload.year,
                month: payload.month,
                day: payload.day,
                monat: MONATE[payload.month - 1],
                wochentag: d ? WOCHENTAGE[d.getDay()] : '',
            };
        },
        [SET_ACTIVE_DATE](state, d) {
            console.log('COMMIT_ACTIVE-DATE', d);
            state.activeDate = {
                date: d,
                year: d.getFullYear(),
                month: d.getMonth() + 1,
                day: d.getDate(),
                monat: MONATE[d.getMonth()],
                wochentag: WOCHENTAGE[d.getDay()],
            };
        },
        [SET_LOADING](state, zustand) {
            state.loading = !!zustand;
        },
        [SET_ARBEITSJAHR](state, payload) {
            console.log('COMMIT_ARBEITSJAHR', payload);
            state.arbeitsjahr = payload;
        },
        [SET_ARBEITSTAG](state, payload) {
            console.log('COMMIT_ARBEITSTAG', payload);
            payload.ende = new Date(payload.Ende);
            state.arbeitstag = payload;
            //state.date = new Date(payload.year, payload.month, payload.day);
            //state.year = payload.year;
            //state.month = payload.month;
            //state.day = payload.day;
        },
        updateArbeitstag(state, {
            field,
            value
        }) {
            state.arbeitstag[field] = value;
        },
    },
    actions: {
        loadArbeitJahr({
            commit,
            state
        }) {
            commit(SET_LOADING, true);
            axios('/arbeit/' + state.activeDate.year, {
                headers: {
                    Accept: 'application/json',
                },
            }).then(function (response) {
                commit(SET_ARBEITSJAHR, response.data.arbeitsjahr);
                commit(SET_LOADING, false);
            });
        },
        loadArbeitMonat({
            commit,
            state
        }) {
            commit(SET_LOADING, true);
            axios('/arbeit/' + state.activeDate.year + '/' + state.activeDate.month, {
                headers: {
                    Accept: 'application/json',
                },
            }).then(function () {
                //commit(SET_ARBEITSJAHR, response.data.arbeitstag);
                commit(SET_LOADING, false);
            });
        },
        loadArbeitstag({
            commit,
            state
        }) {
            commit(SET_LOADING, true);
            let year = state.activeDate.year,
                month = state.activeDate.month,
                day = state.activeDate.day;
            axios('/arbeit/' + year + '/' + month + '/' + day, {
                headers: {
                    Accept: 'application/json',
                },
            }).then(function (response) {
                // handle success

                commit(SET_ARBEITSTAG, response.data.arbeitstag);
                commit(SET_LOADING, false);
            });
        },
    },
});