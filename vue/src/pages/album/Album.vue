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
  <div class>
    <section class="hero is-medium is-transparent">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">{{ album.title }}</h1>
          <h2 class="subtitle">{{ album.subtitle }}</h2>
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
    <section class="section">
      <div class="container">
        <h1 class="title">Section</h1>
        <h2 class="subtitle">
          A simple container to divide your page into
          <strong>sections</strong>, like the one you're currently reading
        </h2>
        <pswp-gallery :images="album.images"></pswp-gallery>
        <button type="button" class="button" @click="isModalVisible = true">Open Modal!</button>
        <upload-modal :visible.sync="isModalVisible" :url="album.file.path"/>
      </div>
    </section>
  </div>
</template>

<style lang="css">
.has-bg-img {
    background: url('/serve/home/14/DSC02359.jpg') center center;
    background-size: cover;
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
