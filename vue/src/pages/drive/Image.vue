<template>
    
    <cloud11-page :account="account">

        <Parallax :image="'/serve' + file.path" :ratio="Math.floor(image.Height / image.Width * 100)">

            <template #attached>

                <section class="hero attached" style="z-index:1">
                    <div class="hero-body">
                        <div class="container">
                            <h1 class="title outline4">{{ image.Title }}</h1>
                            <h2 class="subtitle outline4">{{ image.Caption }}</h2>
                        </div>
                    </div>
                    <div class="hero-foot">
                        <nav class="navbar ">
                            
                             <router-link :to="href(siblings.parent)" class="navbar-item"><feather-icon name="arrow-left"/></router-link>
                            <a class="navbar-item">
                                <feather-icon name="activity"/>
                            </a>
                            <a class="navbar-item">
                                <feather-icon name="database"/>
                            </a>
                             <router-link :to="href(siblings.prev)" class="navbar-item"><feather-icon name="chevron-left"/></router-link>
                             <router-link :to="href(siblings.next)" class="navbar-item"><feather-icon name="chevron-right"/></router-link>
                        </nav>
                    </div>
                </section>
            </template>


            <section class="section has-background-light">

                <form method="POST">
                <div class="columns">
                    <div class="column">
                        
                            <div class="field">
                                <label class="label">Title</label>
                                <div class="control">
                                    <input class="input" type="text" placeholder="Title" name="title" :value="image.Title" />
                                </div>
                            </div>

                            <div class="field">
                                <label class="label">Caption</label>
                                <div class="control">
                                    <textarea name="caption" class="textarea" placeholder="Caption">{{image.Caption}}</textarea>
                                </div>
                                <p class="help">This is a help text</p>
                            </div>
                            <div class="field">
                                <label class="label">Cutline</label>
                                <div class="control">
                                    <textarea name="cutline" class="textarea" placeholder="Cutline">{{image.Cutline}}</textarea>
                                </div>
                            </div>
                        
                    </div>


                    <div class="column">
                        <div class="card">
                            <div class="card-image">
                                <figure class="image is-4by3" :style="{'padding-top': image.Ratio+'%'}">
                                    <img :src="imageSrc" :alt="image.Name">
                                </figure>
                            </div>
                            <div class="card-content">
                                <div class="media">
                                    <div class="media-left">
                                        <figure class="image is-48x48">
                                            <img :src="'/serve' + file.Path" :alt="image.Title">
                                        </figure>
                                    </div>
                                    <div class="media-content">
                                        <p class="title is-4">{{image.Title}}</p>
                                        <p class="subtitle is-6">{{image.Caption}}</p>
                                    </div>
                                </div>

                                <div class="content">
                                    {{image.Cutline}}
                                    <br>
                                    <time datetime="2016-1-1">11:09 PM - 1 Jan 2016</time>
                                </div>
                            </div>
                            <footer class="card-footer">
                                <button type="submit" href="#" class="card-footer-item">Save</button>
                                <a href="#" class="card-footer-item">Edit</a>
                                <a href="#" class="card-footer-item">Delete</a>
                            </footer>
                        </div>
                    </div>
                    <div class="column">
                        <nav class="panel">
                            <p class="panel-heading">
                                Exif-Daten:
                            </p>
                            <div class="panel-block">
                                {{image.Format}} {{image.Width}}x{{image.Height}}
                                {{image.ColorModel}}
                            </div>
                            <a class="panel-block">
                                <span class="panel-icon">
                                    <i class="fas fa-book" aria-hidden="true"></i>
                                </span>
                                {{image.Exif.Taken}}
                            </a>
                            <div class="panel-block">

                                {{image.Exif.Lat}},
                                {{image.Exif.Lng}}
                            </div>
                            <div class="panel-block is-flex">
                                <div class="has-text-left">
                                    Model:
                                </div>
                                <span class="has-text-right">{{image.Exif.Model}}</span>

                            </div>
                        </nav>
                    </div>

                </div>
                </form>
            </section>
       </Parallax>    
    </cloud11-page>
   
</template>

<script>
import Cloud11Figure from '@/components/Cloud11Figure.vue';
import { mapState, mapActions } from 'vuex'
import AlbumImage from '../album/AlbumImage.vue';
import Parallax from '@/components/Parallax.vue';
import NavigationOverlay from '@/components/NavigationOverlay.vue';
import NavbarTop from '@/components/NavbarTop.vue';
import Cloud11Page from '@/components/Cloud11Page.vue';
import FeatherIcon from '@/components/FeatherIcon.vue';

// So to sum it up. I love taking photographs and the pictures you see above are the result of this hobby of mine and are uploaded as images on Quora.

//  image is the root or universal term and is used for any kind of 2 dimensional visual presentation that can be on paper, canvas, display screen or simply what is in front of your eyes.

// picture implies that a person has made it, but doesn't specify how while 
// photo means a camera captured an image in the real world at a specific time and place.
//Photos are always images, but there are many images that are not photos.  An image is what you see; a photograph is one example of an image.
// https://www.quora.com/Whats-the-difference-between-a-picture-an-image-and-a-photo
// PICTURE IS A PAINTING OR DRAWING. 

export default {
    name: "Image",
    components: {
        FeatherIcon,
        Cloud11Figure,
        AlbumImage,
        Parallax, NavigationOverlay, NavbarTop, Cloud11Page
    },
    
    data() {
        return {
            account: {},
            //form: Object.assign({}, src);
        };
    },

    computed: {
        ...mapState(['file', 'image', 'siblings']),
        imageSrc() {
            console.log(this.file);
            return '/serve' + this.file.path
        },
        
        
    },
    
    created() {   
        console.log('Image.vue 2 =>', this.image.path);
     },
    methods: {
        href(path) {
            return "/home" + path
        }
    },
    
};
</script>
    
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
