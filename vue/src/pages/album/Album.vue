<script>
import UploadModal from '@/components/modal.vue';
import PswpGallery from '@/components/PswpGallery.vue';
import Cloud11Figure from '@/components/Cloud11Figure.vue';
import TipTap from './TipTap.vue';
import { mapState, mapActions } from 'vuex'
import AlbumImage from './AlbumImage.vue';
import Parallax from '@/components/Parallax.vue';
import NavigationOverlay from '@/components/NavigationOverlay.vue';
import NavbarTop from '@/components/NavbarTop.vue';
import Cloud11Page from '@/components/Cloud11Page.vue';
import FeatherIcon from '@/components/FeatherIcon.vue';


export default {
    name: "Album",
    components: {
        FeatherIcon,
        UploadModal,
        PswpGallery,
        Cloud11Figure,
        TipTap,
        AlbumImage,
        Parallax, NavigationOverlay, NavbarTop, Cloud11Page
    },
    provide: {
        prefix:"/serve/home",
        album: "Mallorca"
    },
    
    data() {
        return {
            account: {},
            isModalVisible: false,
            selectedSource: "",
            markdown: null,
            isSourceDropdownOpen: false
        };
    },

    computed: {
        ...mapState(["image", "images", "serveURL", "baseURL", "meta", "sources"]),
        source() {
            return this.selectedSource ? this.sources.filter(s => s.name == this.selectedSource)[0] : {photographer:null, camera:null}
        },
        filteredImages() {
            return this.images.filter(image => {
                return image.source == this.selectedSource;      
            })
        },

        appStyles() {
            return {}
        },

    },
    created() {
        
        console.log('Album.vue =>', this.serveURL);
    },
    methods: {
      handleSelect(i) {

         this.$set(i, "selected", !i.selected)  

          this.$store.commit('diaryImage', {
            diaryName: 'undefined',
            image: i,
            mode: i.selected ? 'add' : 'remove',
          })
      },
      selectSource(name) {
          this.isSourceDropdownOpen = !this.isSourceDropdownOpen
          console.log("selectSource", name)
          this.selectedSource = name;
      },
      onUpdateMarkdown(event) {
          console.log("onupdateMarkdown", event);
      },
      put(source) {
          const url = this.baseURL + (source.name ? '/' + source.name : "");
          console.log("put", url);
          fetch(url, {
            method: 'put',
            headers: {
                'Accept': 'application/json, text/plain, */*',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({CreateThumbnails: true})
            })
            .then(res=>res.json())
            .then(res => console.log(res));

      }
    },
    
};
</script>

<template>
    
    <cloud11-page :account="account">
    
        <Parallax :image="image.src" :ratio="Math.floor(image.h / image.w * 100)">

            <template #attached>

                <section class="hero attached" style="z-index:1">
                    <div class="hero-body">
                        <div class="container">
                            <h1 class="title outline4">{{ meta.title }}</h1>
                            <h2 class="subtitle outline4">{{ meta.subtitle }}</h2>
                        </div>
                    </div>
                    <div class="hero-foot">
                        <nav class="navbar ">
                            
                            <a :href="baseURL" class="navbar-item" >
                                <feather-icon name="folder"/> directory
                            </a>
                            <a class="navbar-item">
                                <feather-icon name="activity"/>
                            </a>
                            <a class="navbar-item">
                                <feather-icon name="database"/>
                            </a>
                            <div class="navbar-item has-dropdown" :class="{'is-active': isSourceDropdownOpen}">
                                <a class="navbar-link is-arrowless" @click="isSourceDropdownOpen=!isSourceDropdownOpen" >
                                    <feather-icon name="database"/> {{selectedSource||"source"}}
                                </a>

                                <div class="navbar-dropdown">
                                <a class="navbar-item" 
                                    v-for="s in sources" :key="s.name" 
                                    :class="{'is-active': selectedSource==s.name}" 
                                    @click="selectedSource=s.name" >
                                    {{s.name}}
                                </a>
                                <hr class="navbar-divider">
                                <div class="navbar-item">
                                    Version 0.7.5
                                </div>
                                </div>
                            </div>
                        </nav>
                    </div>
                </section>
            </template>



            <section class="section has-background-light">
                <div class="container">
                    <div class="tabs">
                        <ul>
                            <li :class="{'is-active':selectedSource==s.name}" @click="selectedSource=s.name"
                                v-for="s in sources" :key="s.name">
                                <a>{{s.name}}</a>
                            </li>
                        </ul>
                    </div>
                    <h1 class="title">{{selectedSource}}</h1>
                    <h2 class="subtitle">
                       
                    </h2>
                    <button @click="put(source)">put</button>
                    <pswp-gallery :images="filteredImages" :src="serveURL" ></pswp-gallery>
                </div>
            </section>


          
            
            <!-- <section class="section has-background-primary" v-for="name in diaryNames" :key="name">
                <div class="container">
                    <h1 class="title">{{diaries[name].title}}</h1>
                    <h2 class="subtitle">
                        {{diaries[name].title}}
                    </h2>
                    {{ diaries[name] }}
                    <img :src="'/serve/home' + i.URL" v-for="i in diaries[name].images" :key="i.name" style="width: 100px" />
                    <pswp-gallery :images="diaries[name].images"></pswp-gallery>
                </div>
            </section> 

            <section class="section has-background-light" >
                <div class="container">
                    <h1 class="title"></h1>
                    <textarea v-model="markdown"></textarea>
                    <tip-tap  v-model="markdown"></tip-tap>
                </div>
            </section>-->

            <footer class="footer is-dark has-text-white">
                <div class="columns">
                    <ul class="column is-one-quarter">
                        <li class="subtitle is-4">Title 5</h5>
                        <li class="is-5">Subtitle 5</h5>
                        <li class="is-5">Subtitle 5</li>
                        <li class=" is-5">Subtitle 5</li>
                        <li class=" is-6">Subtitle 6</li>
                        <li class="is-6">Subtitle 6</li>
                        <li class=" is-6">Subtitle 6</li>
                        <li class="subtitle is-6">Subtitle 6</li>
                    </ul>
                    <div class="column">Auto</div>
                    <div class="column">
                        <!-- <div v-for="source in album.sources" :key="source.name">
                            <button @click="put(source)">thumbs</button>
                            <section class="section has-background-dark" v-if="selectedSource==source.name">
                                <div class="container" >
                                    <h1 class="title">{{source.name}}</h1>
                                    <h2 class="subtitle">{{source.camera}} / {{source.photographer}}</h2>

                                    <div v-for="(i, index) in source.images" :key="i.URL">
                                        <cloud11-figure 
                                        :src="'/serve/home/' + i.URL" 
                                        :type="i.selected ? 'sel' : 'not'" 
                                        :tags="['leaf','plant','forest','green']"
                                        @select="handleSelect(i)" 
                                        /> 
                                        <label class="checkbox">
                                            <input type="checkbox" v-model="i.selected">
                                            <span style="color:white">selected: {{i.selected}}</span>
                                        </label>
                                    </div>
                                </div>
                            </section>
                        </div> -->

                    </div>
                </div>

                <div class="content has-text-centered">
                    <p>
                        <strong>Bulma</strong> by <a href="https://jgthms.com">Jeremy Thomas</a>. The source code is
                        licensed <a href="http://opensource.org/licenses/mit-license.php">MIT</a>. The website content is
                        licensed
                        <a href="http://creativecommons.org/licenses/by-nc-sa/4.0/">CC BY NC SA 4.0</a>
                        .
                    </p>
                    <button type="button" class="button" @click="isModalVisible = true">Open Modal!</button>
                    <!-- <upload-modal :visible.sync="isModalVisible" :url="album.file.path" /> -->
                </div>
            </footer>

            <section class="section has-background-light" >
                <div class="container">
                    <pre>album</pre>
                </div>
            </section>

       </Parallax>

        
</cloud11-page>
   
</template>


    
<style lang="scss">
    .hero.attached { 
        .title {
            //display: inline-block;
            float:left;
            position: relative;
            background-color: violet;
            padding: .25rem .75rem;
            z-index: 0;
        }
        .subtitle {
            //display: inline-block;
            float:left;
            clear:left;
            position: relative;
            top:-.5rem;
            background-color: yellow;
            padding: .125rem .5rem;
            margin: .5rem;
        }
        .navbar {
        
            .navbar-item, .navbar-link {
                color: white;
                //text-shadow: 1px 1px 0px black, 1px -1px 0px black, -1px 1px 0px black, -1px -1px 0px black;
                &:hover, &.is-active, &.has-dropdown { background: rgba(206, 198, 198, 0.5); }
                .icon {
                    margin-left: 0.25rem;
                    margin-right: 0.25rem;
                }
            }
        }
    }
</style>

<style lang="css">


.outline {
    color: black !important;
    text-shadow: 1px 1px 0px white, 1px -1px 0px white, -1px 1px 0px white, -1px -1px 0px white;
}
.outline4 {
    color: white !important;
    text-shadow: 1px 1px 0px black, 1px -1px 0px black, -1px 1px 0px black, -1px -1px 0px black;
}




.thumb {
    display: inline-block;
    margin: 1rem;
    max-height: 100px;
}

.selected {
  border: 5px solid white;
}


</style>
