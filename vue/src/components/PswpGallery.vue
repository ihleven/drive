<style>
.preview-img-item {
    cursor: pointer;
}

.item {
    display: inline-block;
    background-color: #fff;
    border-radius: 5px;
    margin: 7px;
}

.shadow {
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
    position: relative;
    border-radius: 2px;
    -webkit-transform: translateY(0);
    -webkit-transition: all 0.6s cubic-bezier(0.165, 0.84, 0.44, 1);
    transition: all 0.6s cubic-bezier(0.165, 0.84, 0.44, 1);
}
.shadow:after {
    content: '';
    border-radius: 2px;
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.6);
    opacity: 0;
    -webkit-transition: all 0.6s cubic-bezier(0.165, 0.84, 0.44, 1);
    transition: all 0.6s cubic-bezier(0.165, 0.84, 0.44, 1);
}
.shadow:hover {
    -webkit-transform: scale(1.1, 1.1);
    transform: scale(1.1, 1.1);
}
.shadow:hover:after {
    opacity: 1;
}
</style>

<template>
  <div class="preview-img-list">
    <a
      class="preview-img-item"
      v-for="(item, index) in items"
      :key="index"
      @click="$pswp_open(index, items)"
    >
      <img :src="item.msrc" class="item shadow" @error="imageErrorHandler">
    </a>
  </div>
</template>

<script>
export default {
    name: 'pswp-gallery',

    props: ['images', 'src'],

    data() {
        return {};
    },
    computed: {
        items() {
            if (!this.images) return [];
            return this.images.map(item => {
                return {
                    msrc: this.src + '/' + item.source + '/thumbs/x100/' + item.name,
                    src: item.src,
                    w: item.w,
                    h: item.h,
                    title: item.name,
                };
            });
        },
    },
    created() {
        console.log('images:', this.src);
    },

    methods: {
        imageErrorHandler(event) {
            //window.location.href = '/hello'
            console.log('imageErrorHandler(', event, ')');
        },
    },
};
</script>
