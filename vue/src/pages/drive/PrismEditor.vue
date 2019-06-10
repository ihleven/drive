<template>
    <div class="prism-editor">
        <vue-prism-editor 
            :code="code"
            :language="lang" 
            :readonly="false"
            :line-numbers="true"
            :auto-style-line-numbers="true"
            :emit-events="false"
            @change="onChange"
            @keyup="keyup"
            @keydown="keydown"
            @editorClick="editorClick"></vue-prism-editor>
    </div>
</template>

<script>
// zunaechst prism.js
//nicht import "prismjs" da sonst css/prism.css mit geladen wird
import "prismjs/prism.js";
// theme: Tomorrow Night
import "prismjs/themes/prism.css";

// dann der Editor:
import VuePrismEditor from 'vue-prism-editor'
// css damit z.B. die linenumbers richtig aussehen
import "vue-prism-editor/dist/VuePrismEditor.css";


export default {
    components: {
        VuePrismEditor
    },
    props: ['code', 'lang'],
    mounted() {},
    methods: {
        onChange(e) {
            this.$emit("update", e)
        },
        keyup(e) {
            console.log("event keyup => ", e);
        },
        keydown(e) {
            console.log("event keydown => ", e);
        },
        editorClick(e) {
            console.log("event editorClick => ", e);
        }
    }
}
</script>

<style>

/* 
    https://github.com/jgthms/bulma/issues/1708 
    Bulma conflicts with Prism.js syntax highlighting plugin
*/

.token.number {
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