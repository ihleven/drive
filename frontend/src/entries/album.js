import Vue from 'vue'
import router from "../router";

import PhotoSwipePlugin from '../PhotoSwipePlugin'

import PswpGallery from "@/components/PswpGallery.vue";


import UIkit from 'uikit'
import Icons from 'uikit/dist/js/uikit-icons'

UIkit.use(Icons)

Vue.config.productionTip = false
Vue.use(PhotoSwipePlugin)

new Vue({
    router,
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