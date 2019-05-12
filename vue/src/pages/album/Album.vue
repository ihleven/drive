<script>
import UploadModal from '@/components/modal.vue';
import PswpGallery from '@/components/PswpGallery.vue';
import Cloud11Figure from '@/components/Cloud11Figure.vue';
import TipTap from './TipTap.vue';
import { mapState, mapActions } from 'vuex'
import AlbumImage from './AlbumImage.vue';
import Parallax from './Parallax.vue';
import NavigationOverlay from './NavigationOverlay.vue';

export default {
    name: "Album",
    components: {
        UploadModal,
        PswpGallery,
        Cloud11Figure,
        TipTap,
        AlbumImage,
        Parallax, NavigationOverlay
    },
    provide: {
        prefix:"/serve/home",
        album: "Mallorca"
    },
    
    data() {
        return {
            isModalVisible: false,
            selectedSource: null,
            markdown: null,
            menuOpen:false,
        };
    },
    computed: {

        appStyles() {
            return {
            }
        },
        ...mapState({
            album: state => state.album,
            diaries: state => state.album.diaries,
            pages: state => state.album.pages,
            diaryNames: state => state.album.diaryNames,
        }),
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
      selectSource(source) {
          if (source.name === this.selectedSource) {
              this.selectedSource = null;
          } else {
              this.selectedSource=source.name
          }
      },
      onUpdateMarkdown(event) {
          console.log("onupdateMarkdown", event);
      },
      put(source) {
          fetch(this.album.file.path + '/' + source.name, {
            method: 'put',
            headers: {
                'Accept': 'application/json, text/plain, */*',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({CreateThumbnails: true})
            })
            .then(res=>res.json())
            .then(res => console.log(res));

      },
      overlayNavigationToggler() {
        this.menuOpen = !this.menuOpen;
      }
    },
    created() {
        this.markdown = this.pages[0]
        console.log('Album.vue =>', this.menu);
    },
};
</script>

<template>
    
    <div class="application-wrapper">

        <navigation-overlay :open.sync="menuOpen">
            <h3 class="title is-spaced"><a href="/alben">Fotoalben</a></h3>
            <h4 class="subtitle"><a href="/alben/Mallorca">Mallorca</a></h4>
            <a class="subtitle" href="/alben/hochzeitsreise">Hochzeitsreise</a>
            <h2 class="subtitle">Fullheight subtitle</h2>
        </navigation-overlay>
                
        <Parallax :image="'/serve/home/Mallorca/Roger/DSC03702.JPG'" :ratio="75" 
                  :perspective="1" :translate-z="-1">

            <template #header>
                <nav class="navbar is-transparent top" role="navigation" aria-label="main navigation">            
                        
                        <a class="navbar-item" href="/">
                            <!--<span class="login">matthias@</span>-->
                            <span class="domain">ihle.</span>
                            <svg class="cloud-icon">
                                <use xlink:href="/feather.svg#cloud"></use>
                            </svg>
                        </a>
                        
                    
                    <a class="navbar-item hamburger" @click="overlayNavigationToggler">
                        <svg class="">
                            <use xlink:href="/feather.svg#menu"></use>
                        </svg>
                    </a>

                </nav>
            </template>


            <section class="hero album-meta" >
                <div class="hero-body">
                    <div class="container">
                        <h1 class="title outline4">{{ album.title }}</h1>
                        <h2 class="subtitle outline">{{ album.subtitle }}</h2>
                    </div>
                </div>
                <div class="hero-foot">
                    <nav class="tabs is-right is-boxed is-dark">
                        <div class="container">
                            <ul>
                                <li v-for="source in album.sources" :key="source.name" 
                                    :class="{'is-active': selectedSource==source.name}" 
                                    @click="selectSource(source)" >
                                    <a>{{source.name}}</a>
                                </li>
                            </ul>
                        </div>
                    </nav>
                </div>
            </section>

            <section class="section has-background-white">
                <div class="container">
                    <div class="timeline is-centered">
                        <header class="timeline-header">
                            <a class="tag is-medium is-primary" href="/home/60/">Hotel Friday Attitude</a>
                        </header>
                        <div class="timeline-item is-primary" v-for="moment in album.moments" :key="moment.title">
                            <div class="timeline-marker is-primary is-image is-64x64" 
                            :style="{background: 'url( /serve/home' + moment.image + ')', 'background-size':'cover'}"
                            ></div>
                            <div class="timeline-content" :style="{'padding':'1.5rem 3rem 0 3rem'}">
                                <p class="heading">{{moment.title}}</p>
                                <p>Timeline content</p>
                            </div>
                        </div>
                        <header class="timeline-header">
                            <span class="tag is-primary">2017</span>
                        </header>
                        <div class="timeline-item is-danger">
                            <div class="timeline-marker is-danger is-icon">
                                <i class="fa fa-flag"></i>
                            </div>
                            <div class="timeline-content">
                                <p class="heading">March 2017</p>
                                <p>Timeline content - Can include any HTML element</p>
                            </div>
                        </div>
                        <header class="timeline-header">
                            <span class="tag is-medium is-primary">End</span>
                        </header>
                    </div>
                </div>
            </section>

            <section class="section has-background-dark" v-if="!selectedSource">
                <div class="container">
                    <h1 class="title">Section</h1>
                    <h2 class="subtitle">
                        A simple container to divide your page into
                        <strong>sections</strong>, like the one you're currently reading
                    </h2>
                    <pswp-gallery :images="album.images"></pswp-gallery>
                </div>
            </section>


            <section class="section has-background-light" v-for="moment in album.moments" :key="moment.title">
                <div class="container">
                    <h1 class="title">{{moment.title}}</h1>
                    <h2 class="subtitle">
                        {{moment.title}}
                    </h2>
                    <pswp-gallery :images="moment.images"></pswp-gallery>
                </div>
            </section>
            <section class="section has-background-primary" v-for="name in diaryNames" :key="name">
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
            </section>

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
                        <div v-for="source in album.sources" :key="source.name">
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
                        </div>

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
                    <upload-modal :visible.sync="isModalVisible" :url="album.file.path" />
                </div>
            </footer>

        </Parallax>

        

    </div>
</template>


<style lang="scss" scoped>


 
    .navbar.top {
        max-width: 100%;
        margin: 0 auto;
        background-color: transparent;
        justify-content: space-between;
    }
    .navbar.top .login {
        color: white;
        font-size: 1rem;
        font-weight: 200;
    }
    .navbar.top .domain {
        color: white;
        font-size: 1.25rem;
        font-weight: 400;
    }
    .navbar.top .cloud-icon {
        width:1.5rem;
        height:1.5rem;
        fill: none;
        stroke: white;
        stroke-width: 2;
        stroke-linecap: round;
        stroke-linejoin: round;
    }
    .navbar.top .hamburger:hover {
        background-color: rgba(255, 255, 255, 0.26)!important;;
        padding: 0.25rem 0.5rem;
        transition: all .2s ease-in-out 0s;
    }
    .navbar.top .hamburger svg {
        stroke: white;
        stroke-width: 2;
        stroke-linecap: round;
        stroke-linejoin: round;
        width:2rem;
        height:2rem;
        fill: none;
    }
    .navbar.top .hamburger:hover svg {
        stroke-width: 2.5;
        width:2.5rem;
        height:2.5rem;
        transition: all .2s ease-in-out 0s;
    }
    
</style>

<style lang="css">

.album-meta {
    position:absolute;
    transform:translateY(-100%);
    width: 100%;
}

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
