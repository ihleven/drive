import Vue from 'vue'

import 'semantic-ui-css/semantic.min.css';



import PhotoSwipePlugin from '../PhotoSwipePlugin'

import PswpGallery from "@/components/PswpGallery.vue";

Vue.config.productionTip = false
Vue.use(PhotoSwipePlugin)

new Vue({
    el: '#vueapp',
    components: {
        PswpGallery
    },
    data: function () {
        return {
            diary: JSON.parse(document.getElementById('data-diary').innerHTML)
        }

    },
    mounted() {
        console.log('mounted diary:', this.diary);
    },
})