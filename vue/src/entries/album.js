import Vue from 'vue'
import Router from "vue-router";
import Home from "../routes/Home.vue";

import PhotoSwipePlugin from '../PhotoSwipePlugin'

import PswpGallery from "@/components/PswpGallery.vue";


import UIkit from 'uikit'
import Icons from 'uikit/dist/js/uikit-icons'

UIkit.use(Icons)

Vue.config.productionTip = false
Vue.use(PhotoSwipePlugin)


Vue.use(Router);

console.log("router", process.env.BASE_URL)

export default new Router({
    mode: "history",
    base: process.env.BASE_URL,
    routes: [{
            path: "/",
            name: "home",
            component: Home
        },
        {
            path: "/about",
            name: "about",
            // route level code-splitting
            // this generates a separate chunk (about.[hash].js) for this route
            // which is lazy-loaded when the route is visited.
            component: () =>
                import( /* webpackChunkName: "about" */ "../routes/About.vue")
        }
    ]
});


new Vue({


    el: '#vueapp',
    components: {
        PswpGallery
    },
    data: function () {
        return {
            directory: JSON.parse(document.getElementById('data-album').innerHTML)
        }

    },
    mounted() {
        console.log('mounted album:', this.directory);
    },
}).$mount('#vueapp')