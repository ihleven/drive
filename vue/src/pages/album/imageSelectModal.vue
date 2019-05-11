<template>
  <transition name="modal-fade">
    <div class="modal" :class="{ 'is-active': visible }">
      <div class="modal-background"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <div class="tabs is-fullwidth">
            <ul>
              <li
                :class="{ 'is-active': s.name == source.name }"
                v-for="s in sources"
                :key="s.name"
                @click="selectSource(s)"
              >
                <a>{{ s.name }}</a>
              </li>
            </ul>
          </div>
        </header>
        <section class="modal-card-body">
          {{ images }}
          <div class>
            <img
              v-for="image in source.images"
              :key="image.name"
              class="image"
              :src="'/serve' + source.path + '/thumbs/x100/' + image.name"
              style="max-height:5rem;"
              @click="toggleImage(image)"
              :class="{ selected: image.selected }"
            >
          </div>
        </section>

        <footer class="modal-card-foot">
          <button class="button is-success" @click="close">Save changes</button>
          <button class="button">Cancel</button>
        </footer>
      </div>

      <button class="modal-close is-large" aria-label="close" @click="close">close</button>
    </div>
  </transition>
</template>

<script>
import { mapState } from 'vuex';

export default {
    name: 'ImageSelectModal',
    inject: ['prefix', 'album'],
    props: {
        visible: Boolean,
        url: String,
    },
    data() {
        return {
            source: { images: [] },
            images: [],
        };
    },
    computed: mapState({
        sources: state => state.album.sources,
    }),
    methods: {
        selectSource(source) {
            this.source = source;
        },
        close() {
            this.$emit('update:visible', false);
            this.$emit('selected', this.images);
        },
        toggleImage(image) {
            console.log('toggle:', image.name);
            if (image.selected) {
                let index = this.images.findIndex(i => i.name == image.name);
                this.images.splice(index, 1);
            } else {
                this.images.push({ source: this.source.name, name: image.name });
            }
            image.selected = !image.selected;
        },
    },
    mounted() {
        console.log('upload modal mounted', this.sources);
    },
};
</script>

<style>
.modal-fade-enter,
.modal-fade-leave-active {
    opacity: 0;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
    transition: opacity 0.5s ease;
}
.selected {
    border: 4px solid black !important;
}
.image {
    display: inline-block !important;
    margin: 0 10px 10px 0;
}
</style>

<style lang="scss"></style>
