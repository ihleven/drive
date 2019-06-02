<script>
export default {
    name: 'Parallax',
    props: {
        image: String,
        ratio: Number,
        perspective: {
            type: Number,
            default() {
                return 1;
            },
        },
        translateZ: {
            type: Number,
            default() {
                return -2;
            },
        },
    },
    computed: {
        scale() {
            return 1 + (this.translateZ * -1) / this.perspective;
        },
        styles() {
            let styles = {
                transform: 'translateZ(' + this.translateZ + 'px) scale(' + this.scale + ')',
            };
            if (this.image) {
                styles['background-image'] = 'url(' + this.image + ')';
            }
            if (this.ratio) {
                styles.height = this.ratio + 'vw';
            }
            return styles;
        },
        mainStyles() {
            return {
                perspective: this.perspective + 'px',
            };
        },
    },
    created() {
        console.log('Parallax =>', this.perspective, this.translateZ, this.scale, this.ratio);
    },
};
</script>

<template>
  <div class="MainContainer" :style="mainStyles">
    <div class="parallax" :style="styles">
      <slot name="header"></slot>
    </div>

    <div class="ContentContainer">
      <slot name="attached"></slot>

      <slot></slot>
    </div>
  </div>
</template>


<style lang="css" scoped>
.MainContainer {
    perspective: 1px;
    perspective-origin: top left;
    transform-style: preserve-3d;
    height: 100vh;
    max-height: 100vh;
    overflow-x: hidden;
    overflow-y: scroll;
    -webkit-overflow-scrolling: touch;
}

.parallax {
    /*height: 66vh;*/
    max-height: 100vh;
    z-index: -1;
    position: relative;
    transform: translateZ(-1px) scale(2);
    transform-origin: 0 0;

    /*background: url(https://www.toptal.com/designers/subtlepatterns/patterns/sakura.png);*/
    background-color: rgb(250, 228, 216);
    background-position: center 0%;
    background-repeat: no-repeat;
    /*background-attachment: fixed;*/
    background-size: 100vw auto;
    background-color: #999;
}

.ContentContainer {
    display: block;
    position: relative;
    z-index: 1;
}

.attached {
    position: absolute;
    transform: translateY(-100%);
    width: 100%;
}
</style>
