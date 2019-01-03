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

  props: ['images'],

  data() {
    return {
      items: []
    }
  },

  created() {
    console.log(this.images)
    this.items = this.images.map(item => {
      if (typeof item == 'string') {
        return this.imageMetaFromFilename(item)
      }
      return {
        msrc: '/serve/' + item.path.replace(item.name, 'thumbs/' + item.name),
        src: '/serve/' + item.path,
        w: item.Width,
        h: item.Height,
        title: item.Title
      }
    })
  },

  methods: {
    imageMetaFromFilename(filename) {
      let regex = /^hochzeit_(\d+)_(\d+)x(\d+).jpg$/,
        match = filename.match(regex),
        thumb = match
          ? '/media/hochzeit-1800/thumbs/hochzeit_' + match[1] + '_x100.jpg'
          : filename

      return {
        msrc: thumb,
        src: '/media/hochzeit-1800/' + filename,
        w: match ? parseInt(match[2]) : 0,
        h: match ? parseInt(match[3]) : 0,
        title: `Bild ${filename}`
      }
    },

    imageErrorHandler(event) {
      //window.location.href = '/hello'
      console.log('imageErrorHandler(', event, ')')
    }
  }
}
</script>
