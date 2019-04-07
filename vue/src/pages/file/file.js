import Vue from 'vue';
import CodeHighlighter from '@/components/CodeHighlighter.vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
import '@/directives/highlighter.js';

import '@/assets/bulma-customize.scss';

//import(/* webpackPreload: true */ 'typeface-clear-sans/index.css');
import 'typeface-clear-sans/index.css';

new Vue({
    el: '#app',
    components: {
        CodeHighlighter,
        MarkdownEditor,
    },
    data() {
        return {
            file: JSON.parse(document.getElementById('data').innerHTML),
            isActive: true,
            edit: false,
        };
    },
    mounted() {
        console.log('file =>', this.file);
        console.log('mounted', this.content);
    },
});
