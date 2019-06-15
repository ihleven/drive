<template>
<div>
    <nav class="navbar" role="navigation" aria-label="main navigation">
        <div class="navbar-brand">
            
            <a class="navbar-item " href="https://ihle.cloud">
                    <strong class="brand">ihle.</strong>
                    <svg class="feather feather-cloud sc-dnqmqq jxshSx" xmlns="http://www.w3.org/2000/svg" width="24"
                        height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                        stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" data-reactid="351">
                        <path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z"></path>
                    </svg>
            </a>
            
            <a role="button" class="navbar-burger burger" aria-label="menu" aria-expanded="false" data-target="navbarBasicExample"
                @click="menuOpen=!menuOpen" :class="{'is-active': menuOpen}">
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
            </a>
        </div>

        <div id="navbarBasicExample" class="navbar-menu" :class="{'is-active': menuOpen}">
            <div class="navbar-start">
                <div class="navbar-item">
<nav class="breadcrumb" aria-label="breadcrumbs">
                  <ul>
                    <li v-for="item in breadcrumbs" :key="item.Path">
                      <router-link :to="item.Path" class="navbar-item" v-text="item.Name"></router-link>

                    </li>
                  </ul>
                </nav>
            </div>
            </div>

            <div class="navbar-end">
                <div class="navbar-item">
                    <div class="buttons">
                        <a class="button is-primary" @click="save" :disabled="pristine">
                            <strong>save</strong>
                        </a>
                        <a class="button is-light"  @click="reset" :disabled="pristine">
                            reset
                        </a>
                        <a class="button is-light"  :href="'/serve' + file.path" target="_blank">
                            open in new tab
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </nav>
    <div class="wrapper">
        
            <div class="container is-fullhd">
                           
                <prism-editor v-model="fileContent" :mime="'text/javascript'" @update="updateContent"></prism-editor>

                    <!--  <markdown-editor :readonly="{{ not .File.Permissions.Write}}">{{.Content}}</markdown-editor>
                    
                    <code-highlighter language="{{.File.MIME.Subtype}}" :readonly="false">{{.Content}}
                    </code-highlighter>
                    
                  <pre class="text"><code class="language-go" v-highlight>{{.Content}}</code></pre>-->
                   
            
            </div>
            <section class="section has-background-light" v-if="mimeType=='markdown'">
                <div class="container">
                    <h1 class="title"></h1>
                    <textarea v-model="fileContent"></textarea>
                    <tip-tap  v-model="fileContent"></tip-tap>
                </div>
            </section>
     </div>  
   </div>
</template>

<script>
import axios from 'axios';

import { mapState, mapActions } from 'vuex'
import Cloud11Page from '@/components/Cloud11Page.vue';
import Parallax from '@/components/Parallax.vue';
import PrismEditor from './PrismEditor';
import TipTap from './TipTap';



export default {
    name: "File",
    components: {
        Cloud11Page,
        Parallax,
        PrismEditor,
        TipTap
    },
    
    data() {
        return {
            menuOpen: false,
            fileContent: null
        };
    },

    computed: {
        ...mapState(['account', 'file', 'content', 'breadcrumbs']), 
        pristine() {
            return this.fileContent == this.content;
        },
        mimeType() {
            return this.file && this.file.mime && this.file.mime.Value == "text/markdown" ? "markdown" : "";
        }
    },
    
    created() {   
        this.fileContent = this.content;
        console.log('File.vue =>', this.file);
     },
    methods: {
        save() {
            this.$store.dispatch("submitFileForm", this.fileContent);
        },
        reset() {
            this.fileContent = this.content;
        }
    },
    
};
</script>
    
<style lang="scss">
    .wrapper {
        overflow-y: auto;
        height: calc(100vh - 52px);
    }
   .brand {
       font-size: 1.3rem;
       font-weight: 500;
   }
.navbar-burger.burger {
    span {
        background-color: white;
        height: 4px;
        border: 1px solid black;
        border-radius: 2px;
        left: calc(50% - 12px);
        width: 24px;
    }
    span:nth-child(1) {
        top: calc(50% - 9px);
        
    } 
    span:nth-child(2) {
        top: calc(50% - 1px);
    }
    span:nth-child(3) {
        top: calc(50% + 7px);
    }
    &:active, &:focus, &:hover {
        background: transparent;
    }
    &.is-active {
        
        span:nth-child(1) {
            transform: translateY(9px) rotate(45deg);
        } 
        span:nth-child(3) {
            transform: translateY(-7px) rotate(-45deg);
        }
    }
}
</style>
