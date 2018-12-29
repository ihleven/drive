import Vue from 'vue'


Vue.config.productionTip = false

import UIkit from 'uikit'
import Icons from 'uikit/dist/js/uikit-icons'
import 'uikit/dist/css/uikit.css'
// loads the Icon plugin
UIkit.use(Icons)


new Vue({
    data: function () {
        return JSON.parse(document.getElementById('data').innerHTML)
    },
    mounted() {
        console.log(this.org)
    }
}).$mount('#vue')