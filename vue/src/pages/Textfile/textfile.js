import Vue from 'vue';
import CodeHighlighter from '@/components/CodeHighlighter.vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';

import '@/assets/bulma-customize.scss';

const data = JSON.parse(document.getElementById('data').innerHTML);

new Vue({
    el: '#app',
    components: {
        CodeHighlighter,
        MarkdownEditor,
    },
    data() {
        return {
            file: data,
        };
    },
    mounted() {
        console.log('file =>', this.file);
    },
});