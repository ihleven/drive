<template>
  <div class="about">
    <h1>This is an about page</h1>

    <pre ref="md">
# Why I built tiptap
I was looking for a text 
* editor 
* fofound 

### some solutions that didn't really satisfy me. The editor shoossation
    </pre>
    <textarea
      class="textarea"
      placeholder="e.g. Hello world"
      v-model="md"
      @input="onMarkdownChange"
    ></textarea>

    <editor-content class="content" :editor="editor"/>
    <!--<div v-html="html"></div>-->
  </div>
</template>


<script>
import { Editor, EditorContent, EditorMenuBubble } from "tiptap";
import {
  Blockquote,
  BulletList,
  CodeBlock,
  HardBreak,
  Heading,
  ListItem,
  OrderedList,
  TodoItem,
  TodoList,
  Bold,
  Code,
  Italic,
  Link,
  Strike,
  Underline,
  History
} from "tiptap-extensions";

import TurndownService from "turndown";
import mymarked from "marked";
mymarked.setOptions({
  renderer: new mymarked.Renderer(),
  highlight: function(code) {
    return code;
  },
  pedantic: false,
  gfm: true,
  tables: true,
  breaks: false,
  sanitize: false,
  smartLists: true,
  smartypants: false,
  xhtml: false
});

let turndownService = new TurndownService();

export default {
  components: {
    EditorContent,
    EditorMenuBubble
  },
  data() {
    return {
      editor: new Editor({
        extensions: [
          new Blockquote(),
          new BulletList(),
          new CodeBlock(),
          new HardBreak(),
          new Heading({ levels: [1, 2, 3] }),
          new ListItem(),
          new OrderedList(),
          new TodoItem(),
          new TodoList(),
          new Bold(),
          new Code(),
          new Italic(),
          new Link(),
          new Strike(),
          new Underline(),
          new History()
        ],
        content: this.html,
        onUpdate: event => {
          this.md = turndownService.turndown(event.getHTML());
          //this.html = event.getHTML();
          console.log(" html => ", event.getHTML());
        }
      }),

      md: "lökjlökjlöj"
    };
  },

  mounted() {
    this.md = this.$refs.md.innerHTML;
    this.editor.setContent(mymarked(this.md), false);
    //this.editor.setContent(this.$refs.child.innerHTML);
    //console.log("md", this.md);
  },
  beforeDestroy() {
    this.editor.destroy();
  },
  methods: {
    clearContent() {
      this.editor.clearContent(true);
      this.editor.focus();
    },
    toMarkdown() {
      return this.turndownService.turndown(this.html);
    },
    onMarkdownChange(event, a, b, c) {
      console.log("onMarkdownChange", event.target.value);
      this.editor.setContent(mymarked(event.target.value), false);
    }
  }
};
</script>

<style lang="scss" scoped>
.about {
  text-align: left;
}
$color-black: #000000;
$color-white: #ffffff;
$color-grey: #dddddd;

.actions {
  max-width: 30rem;
  margin: 0 auto 2rem auto;
}

.export {
  max-width: 30rem;
  margin: 0 auto 2rem auto;

  pre {
    padding: 1rem;
    border-radius: 5px;
    font-size: 0.8rem;
    font-weight: bold;
    background: rgba($color-black, 0.05);
    color: rgba($color-black, 0.8);
  }

  code {
    display: block;
    white-space: pre-wrap;
  }
}
</style>