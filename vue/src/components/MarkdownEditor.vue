<template>
    <div class="hello">
        <div class="tabs">
            <ul>
                <li
                    :class="{ 'is-active': !viewsource }"
                    @click="viewsource = false"
                >
                    <a>Preview</a>
                </li>
                <li
                    :class="{ 'is-active': viewsource }"
                    @click="viewsource = true"
                >
                    <a>Source</a>
                </li>
            </ul>

            <a
                class="button is-info is-outlined"
                @click="reset"
                :disabled="!dirty"
                >Reset</a
            >
            <a class="button is-info is-outlined" :disabled="!dirty">Save</a>
        </div>
        <slot name="source"></slot>
        <div v-show="viewsource">
            <textarea
                class="textarea"
                v-model="markdown"
                @input="onMarkdownChange"
            ></textarea>
        </div>
        <div v-show="!viewsource">
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
                    >
                        <icon name="bold" />
                    </button>

                    <button
                        class="menububble__button"
                        :class="{ 'is-active': isActive.italic() }"
                        @click="commands.italic"
                    >
                        <icon name="italic" />
                    </button>

                    <button
                        class="menububble__button"
                        :class="{ 'is-active': isActive.code() }"
                        @click="commands.code"
                    >
                        <icon name="code" />
                    </button>
                </div>
            </editor-menu-bubble>
            <editor-content class="box content" :editor="editor" />
        </div>
        <!--<div v-html="html"></div>-->
    </div>
</template>

<script>
import { Editor, EditorContent, EditorMenuBubble } from 'tiptap';
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
    History,
} from 'tiptap-extensions';
import TurndownService from 'turndown';
import Marked from 'marked';

Marked.setOptions({
    renderer: new Marked.Renderer(),
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
    xhtml: false,
});

let turndownService = new TurndownService();

export default {
    name: 'MarkdownEditor',
    components: {
        EditorContent,
        EditorMenuBubble,
    },
    props: {
        msg: String,
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
                    new History(),
                ],
                content: this.html,
                onUpdate: event => {
                    this.markdown = turndownService.turndown(event.getHTML());
                    //this.html = event.getHTML();
                    console.log(' html => ', event.getHTML());
                },
            }),
            markdown: '',
            viewsource: false,
        };
    },
    computed: {
        dirty() {
            return this.source != this.markdown;
        },
    },
    mounted() {
        this.source = this.markdown = this.$slots.default[0].text;
        this.editor.setContent(Marked(this.markdown), false);
    },
    beforeDestroy() {
        this.editor.destroy();
    },
    methods: {
        onMarkdownChange(event) {
            console.log('onMarkdownChange', event.target.value);
            this.editor.setContent(Marked(event.target.value), false);
        },
        reset() {
            this.markdown = this.source;
            this.editor.setContent(Marked(this.markdown), false);
        },
    },
};
</script>

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
