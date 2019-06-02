<script>
import Icon from '@/components/Icon';
import { Editor, EditorContent, EditorMenuBubble, EditorMenuBar } from 'tiptap';
import { HardBreak, Heading, Bold, Code, Italic, Link, Strike, Underline, History } from 'tiptap-extensions';
import Image from './ImageNode';

import showdown from 'showdown';
import ImageSelectModal from './imageSelectModal.vue';

export default {
    name: 'TipTap',
    components: {
        Icon,
        EditorContent,
        EditorMenuBubble,
        EditorMenuBar,
        ImageSelectModal,
    },

    props: {
        source: String,
        readonly: { type: Boolean, default: true },
    },
    model: {
        prop: 'source',
        event: 'update:source',
    },
    data() {
        return {
            editor: null,
            metadata: null,
            isImageSelectVisible: false,
            imageCommand: null,
            //viewsource: false,
            converter: new showdown.Converter({ metadata: true }),
        };
    },
    watch: {
        source(md) {
            // so cursor doesn't jump to start on typing
            if (md !== this.markdown) {
                //this.editor.setContent(val);
                let html = this.converter.makeHtml(md);
                this.metadata = this.converter.getMetadata();
                //console.log('showdown', nv, ov);
                this.editor.setContent(html, false);
                //this.html = html;
            }
        },
    },

    mounted() {
        this.editor = new Editor({
            //content: this.html,
            //editable: true,
            extensions: [
                new HardBreak(),
                new Heading({ levels: [1, 2, 3] }),
                new Image(),
                new Bold(),
                new Code(),
                new Italic(),
                new Link(),
                new Strike(),
                new Underline(),
                new History(),
            ],
            onUpdate: ({ getHTML }) => {
                this.markdown = this.converter.makeMarkdown(getHTML(), this.metadata);
                this.$emit('update:source', this.markdown);
            },
        });
        this.editor.setContent(this.converter.makeHtml(this.source));
    },
    beforeDestroy() {
        this.editor.destroy();
    },
    methods: {
        reset() {
            // this.markdown = this.source;
            //this.editor.setContent(Marked(this.markdown), false);
        },
        showImagePrompt(command) {
            this.isImageSelectVisible = true;
            this.imageCommand = command;
        },

        imagesSelected(images) {
            console.log('images selected', images);
            images.forEach(image => {
                console.log('loop image', image);
                this.imageCommand({ src: image.source + '/' + image.name });
            });
        },
    },
};
</script>

<template>
  <div class="container">
    {{ metadata }}
    <br>------- Meta
    <br>
    <editor-menu-bar :editor="editor">
      <div slot-scope="{ commands, isActive }">
        <div class="buttons has-addons is-centered">
          <span
            class="button is-outlined"
            :class="{ 'is-selected': isActive.bold() }"
            @click="commands.bold"
          >Bold</span>
          <span
            class="button is-outlined"
            :class="{ 'is-selected': isActive.italic() }"
            @click="commands.italic"
          >I</span>
          <span
            class="button is-outlined"
            :class="{ 'is-active': isActive.heading({ level: 2 }) }"
            @click="commands.heading({ level: 2 })"
          >H2</span>
          <button class="button is-outlined" @click="showImagePrompt(commands.image)">Image</button>

          <button class="button is-outlined" @click="commands.undo">Undo</button>

          <button class="button is-outlined" @click="commands.redo">Redo</button>
        </div>
      </div>
    </editor-menu-bar>
    <editor-menu-bubble :editor="editor">
      <div
        slot-scope="{ commands, isActive, menu }"
        class="menububble"
        :class="{ 'is-active': menu.isActive }"
        :style="`left: ${menu.left}px; bottom: ${menu.bottom}px;`"
      >
        <button
          class="menububble__button"
          :class="{ 'is-active': isActive.bold() }"
          @click="commands.bold"
        >B</button>

        <button
          class="menububble__button"
          :class="{ 'is-active': isActive.italic() }"
          @click="commands.italic"
        >I</button>

        <button
          class="menububble__button"
          :class="{ 'is-active': isActive.code() }"
          @click="commands.code"
        >C</button>
      </div>
    </editor-menu-bubble>
    <editor-content class="content" :editor="editor"/>

    <image-select-modal :visible.sync="isImageSelectVisible" @selected="imagesSelected"/>
  </div>
</template>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
contenteditable {
    outline: none;
}

$color-black: #000;
$color-white: #fefefe;
.menububble {
    position: absolute;
    display: flex;
    z-index: 20;
    background: $color-black;
    border-radius: 5px;
    padding: 0.3rem;
    margin-bottom: 0.5rem;
    transform: translateX(-50%);
    visibility: hidden;
    opacity: 0;
    transition: opacity 0.2s, visibility 0.2s;

    &.is-active {
        opacity: 1;
        visibility: visible;
    }

    &__button {
        display: inline-flex;
        background: transparent;
        border: 0;
        color: $color-white;
        padding: 0.2rem 0.5rem;
        margin-right: 0.2rem;
        border-radius: 3px;
        cursor: pointer;

        &:last-child {
            margin-right: 0;
        }

        &:hover {
            background-color: rgba($color-white, 0.1);
        }

        &.is-active {
            background-color: rgba($color-white, 0.2);
        }
    }

    &__form {
        display: flex;
        align-items: center;
    }

    &__input {
        font: inherit;
        border: none;
        background: transparent;
        color: $color-white;
    }
}
</style>
