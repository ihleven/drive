import Vue from 'vue';
import VueRouter from 'vue-router';
import Arbeitsjahr from './Arbeitsjahr';
import Arbeitsmonat from './Arbeitsmonat';
import Arbeitstag from './Arbeitstag';
Vue.use(VueRouter);
import store from './arbeit-store';

const KW = {
    template: '<div>kw</div>',
};

const routes = [{
        path: '/arbeit',
        component: {
            template: '<div>defualt</div>',
        },
        // redirect: () => {
        //     const today = new Date();
        //     return {
        //         name: 'arbeitstag',
        //         params: {
        //             year: today.getFullYear(),
        //             month: today.getMonth(),
        //             day: today.getDate(),
        //         },
        //     };
        // },
    },
    {
        path: '/:year',
        name: 'arbeitsjahr',
        component: Arbeitsjahr,
        beforeEnter: (to, from, next) => {
            if (!to.params.month) {
                store.commit('SET_ACTIVE', {
                    year: to.params.year,
                });
                store.dispatch('loadArbeitJahr', to.params);
            }
            next();
        },
        children: [{
                path: 'kw:kw',
                name: 'arbeitswoche',
                component: KW,
            },
            {
                path: ':month',
                name: 'arbeitsmonat',
                component: Arbeitsmonat,
                props: true,
                beforeEnter: (to, from, next) => {
                    if (!to.params.day) {
                        store.commit('SET_ACTIVE', {
                            year: to.params.year,
                            month: to.params.month,
                        });
                    }
                    //store.dispatch('loadArbeitMonat', to.params);
                    next();
                },
                children: [{
                    path: ':day',
                    name: 'arbeitstag',
                    component: Arbeitstag,
                    beforeEnter: (to, from, next) => {
                        //console.log('Arbeitstag Per-Route Guard', to.params);
                        let d = new Date(to.params.year, to.params.month - 1, to.params.day);
                        store.commit('SET_ACTIVE_DATE', d);
                        store.dispatch('loadArbeitstag', to.params);
                        next();
                    },
                }, ],
            },
        ],
    },
];

const router = new VueRouter({
    routes, // short for `routes: routes`
});

router.beforeEach((to, from, next) => {
    //console.log('router.beforeEach', to.params);
    next();
});

export default router;