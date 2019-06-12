<template>
<div >
    <nav class="navbar" role="navigation" aria-label="main navigation">
        <div class="navbar-brand">
            <a class="navbar-item" href="https://bulma.io">
                <img src="https://bulma.io/images/bulma-logo.png" width="112" height="28">
            </a>
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
                <a class="navbar-item">Home</a>
            </div>

            <div class="navbar-end">
                <div class="navbar-item">
                    <div class="buttons">
                        <a class="button is-primary">
                            <strong>Sign up</strong>
                        </a>
                        <a class="button is-light">
                            Log in
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </nav>
    
            <prism-editor :code="content" :mime="'text/javascript'" @update="updateContent"></prism-editor>
        
        <section class="section">
            <div class="container" v-cloak style="border:1px solid black">
               
                    <!--  <markdown-editor :readonly="{{ not .File.Permissions.Write}}">{{.Content}}</markdown-editor>
                    
                    <code-highlighter language="{{.File.MIME.Subtype}}" :readonly="false">{{.Content}}
                    </code-highlighter>
                    
                  <pre class="text"><code class="language-go" v-highlight>{{.Content}}</code></pre>-->
                   
             


                <div class="content" v-text="content"></div>
              





            </div>
        </section>
       
       
   </div>
</template>

<script>
import { mapState, mapActions } from 'vuex'
import Cloud11Page from '@/components/Cloud11Page.vue';
import Parallax from '@/components/Parallax.vue';
import PrismEditor from './PrismEditor';



export default {
    name: "File",
    components: {
        Cloud11Page,
        Parallax,
        PrismEditor,
    },
    
    data() {
        return {
            menuOpen: false
        };
    },

    computed: {
        ...mapState(['account', 'file', 'content']),
        
    },
    
    created() {   
        console.log('File.vue =>', this.file);
     },
    methods: {
        updateContent(content) {
            console.log("onChange", content)
        }
    },
    
};
</script>
    
<style lang="scss">
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
