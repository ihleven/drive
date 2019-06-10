import Vue from 'vue';
import store from './drive-store.js';
import { mapState } from 'vuex';


import PrismEditor from './PrismEditor';


import './drive-styles.scss';


import 'typeface-clear-sans/index.css';

store.dispatch('loadInitialData');
new Vue({
    store,
    el: '#app',
    components: {
//        CodeHighlighter,
//        MarkdownEditor,
        PrismEditor,
    },
    data() {
        return {
            menuOpen: false,
            edit: false,
        }
    },
    computed: {
        ...mapState(['file', 'content', 'breadcrumbs']),
    },
    mounted() {
        console.log('file =>', this.file);
    },
    methods: {
        onChange(content) {
            console.log(content)
            this.$store.commit('updateContent', content);
        }
    }
});
