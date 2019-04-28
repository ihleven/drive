<script>
import UploadModal from '@/components/modal.vue';
import PswpGallery from '@/components/PswpGallery.vue';

export default {
    components: {
        UploadModal,
        PswpGallery,
    },
    props: {
        album: Object,
    },
    data() {
        return {
            isModalVisible: false,
        };
    },
    methods: {},
    mounted() {
        console.log('Album.vue =>', this.album);
    },
};
</script>
<template>
    <div class="bg" style="background-image:url(/serve/home/14/DSC02257.jpg)">
        <section class="hero is-medium is-transparent">
            <div class="hero-head">
                <nav class="tabs">
                    <div class="container">
                        <ul>
                            <li class="is-active">
                                <a>Overview</a>
                            </li>
                            <li>
                                <a>Modifiers</a>
                            </li>
                            <li>
                                <a>Grid</a>
                            </li>
                            <li>
                                <a class="button is-dark" href="/home/14">folder</a>
                            </li>
                            <li>
                                <a>Components</a>
                            </li>
                            <li>
                                <a>Layout</a>
                            </li>
                        </ul>
                    </div>
                </nav>
            </div>
            <div class="hero-body">
                <div class="container">
                    <h1 class="title outline4">{{ album.title }}</h1>
                    <h2 class="subtitle outline">{{ album.subtitle }}</h2>
                </div>
            </div>
            <div class="hero-foot">
                <nav class="tabs">
                    <div class="container">
                        <ul>
                            <li class="is-active">
                                <a>Overview</a>
                            </li>
                            <li>
                                <a>Modifiers</a>
                            </li>
                            <li>
                                <a>Grid</a>
                            </li>
                            <li>
                                <a class="button is-dark" href="/home/14">folder</a>
                            </li>
                            <li>
                                <a>Components</a>
                            </li>
                            <li>
                                <a>Layout</a>
                            </li>
                        </ul>
                    </div>
                </nav>
            </div>
        </section>

        <section class="section has-background-dark">
            <div class="container">
                <h1 class="title">Section</h1>
                <h2 class="subtitle">
                    A simple container to divide your page into
                    <strong>sections</strong>, like the one you're currently reading
                </h2>
                <pswp-gallery :images="album.images"></pswp-gallery>
            </div>
        </section>

        <section class="section">
            <div class="container">
                <div class="timeline is-centered">
                    <header class="timeline-header">
                        <a class="tag is-medium is-primary" href="/home/60/">Hotel Friday Attitude</a>
                    </header>
                    <div class="timeline-item is-primary">
                        <div class="timeline-marker is-primary"></div>
                        <div class="timeline-content">
                            <p class="heading">January 2016</p>
                            <p>Timeline content - Can include any HTML element</p>
                        </div>
                    </div>
                    <div class="timeline-item is-warning">
                        <div class="timeline-marker is-warning is-image is-64x64">
                            <img src="/serve/home/14/DSC02257.jpg" />
                        </div>
                        <div class="timeline-content">
                            <p class="heading">February 2016</p>
                            <p>Timeline content - Can include any HTML element</p>
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
                <div class="column">Auto</div>
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
    </div>
</template>

<style lang="css">
.has-bg-img {
    background: url('/serve/home/14/DSC02359.jpg') center center;
    background-size: cover;
}
.bg {
    background-size: contain;
    background-position: 0 0;
    background-repeat: no-repeat;
}
.outline {
    color: black !important;
    text-shadow: 1px 1px 0px white, 1px -1px 0px white, -1px 1px 0px white, -1px -1px 0px white;
}
.outline4 {
    color: white !important;
    text-shadow: 1px 1px 0px black, 1px -1px 0px black, -1px 1px 0px black, -1px -1px 0px black;
}
.html {
    overflow-y: auto !important;
}
.j {
    perspective: 1px;
    transform-style: preserve-3d;
    height: 100vh;
    overflow-x: hidden;
    overflow-y: scroll;
}
.ParallaxContainer {
    position: relative;
    /*height: 100vh;*/
    transform: translateZ(-1px) scale(2);
    transform-origin: 50% 100%;
    z-index: -1;
}

.ContentContainer {
    position: relative;
    background-color: white;
    z-index: 1;
}
</style>
