/* eslint-disable prettier/prettier */
<template>
  <div class="home">
    <img alt="Vue logo" src="../../../assets/logo.png">

    <div class="section">
      <div class="container">
        <markdown-editor>
          # Why I built tiptap I was looking for a text * editor * fofound ### some solutions that didn't
          really satisfy me. The editor shoossation
        </markdown-editor>
      </div>
    </div>
    <div class="footer">
      <div class="container">
        <HelloWorld msg="Welcome to Your Vue.js App"/>
      </div>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
export default {
    name: 'home',
    components: {
        HelloWorld,
        MarkdownEditor,
    },
};
</script>
