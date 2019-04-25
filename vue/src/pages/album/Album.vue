<script>
import UploadModal from '@/components/modal.vue';

export default {
    components: {
        UploadModal,
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
  <div>
    <section class="hero is-medium is-primary is-bold">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">{{ album.title }}</h1>
          <h2 class="subtitle">{{ album.subtitle }}</h2>
        </div>
      </div>
    </section>
    <section class="section">
      <button type="button" class="button" @click="isModalVisible = true">Open Modal!</button>
      <upload-modal :visible.sync="isModalVisible" :url="album.file.path"/>
      {{ album.file }}
    </section>

    <pswp-gallery :images="album.images"></pswp-gallery>
  </div>
</template>
<style lang="scss"></style>
