import Vue from 'vue';
import CodeHighlighter from '@/components/CodeHighlighter.vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';

import '@/assets/bulma-customize.scss';

//Vue.config.productionTip = false;

//new Vue({
//  render: h => h(App)
//}).$mount("#app");

new Vue({
    el: '#app',
    components: {
        CodeHighlighter,
        MarkdownEditor,
    },
    data() {
        return {};
    },
    mounted() {
        //console.log("mounted", this.content)
    },
});
