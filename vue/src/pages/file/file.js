import Vue from "vue";
import CodeHighlighter from "@/components/CodeHighlighter.vue";

import '@/assets/bulma-customize.scss';


new Vue({
    el: "#app",
    components: {
        CodeHighlighter
    },
    data() {
        return {
            edit: false
        }
    },
    mounted() {
        console.log("mounted", this.content)
    }
});