<template>
    <div class="content">
        <prism-editor v-model="code" :language="lang" :lineNumbers="linenumbers" :readonly="readonly">
            <slot></slot>
        </prism-editor>
    </div>
</template>

<script>
import 'prismjs';
import 'prismjs/themes/prism.css';
import PrismEditor from 'vue-prism-editor';
import 'vue-prism-editor/dist/VuePrismEditor.css'; // import the styles

const langDict = {
    javascript: 'js',
    golang: 'go',
};

export default {
    components: {
        PrismEditor,
    },
    props: {
        language: String,
        linenumbers: {
            type: Boolean,
            default: true,
        },
        readonly: {
            type: Boolean,
            default: true,
        },
    },
    data() {
        return {
            code: this.$slots.default[0].text,
        };
    },
    computed: {
        lang() {
            return langDict[this.language];
        },
    },
    mounted() {
        //this.code = this.$slots.default[0].text;
        console.log('language:', this.lang);
        console.log('linenumbers:', this.linenumbers);
        console.log('readonly:', this.readonly);
    },
    methods: {
        save() {
            console.log(this.code);
        },
    },
};
</script>

<style>
/* https://github.com/jgthms/bulma/issues/1708 */
.content .tag,
.content .number {
    display: inline;
    padding: inherit;
    font-size: inherit;
    line-height: inherit;
    text-align: inherit;
    vertical-align: inherit;
    border-radius: inherit;
    font-weight: inherit;
    white-space: inherit;
    background: inherit;
    margin: inherit;
}
</style>
