<template>
    <div class="pswp" tabindex="-1" role="dialog" aria-hidden="true">
        <div class="pswp__bg"></div>

        <div class="pswp__scroll-wrap">
            <div class="pswp__container">
                <div class="pswp__item"></div>
                <div class="pswp__item"></div>
                <div class="pswp__item"></div>
            </div>

            <div class="pswp__ui pswp__ui--hidden">
                <div class="pswp__top-bar">
                    <div class="pswp__counter"></div>

                    <button class="pswp__button pswp__button--close" title="Close (Esc)"></button>

                    <button class="pswp__button pswp__button--share" title="Share"></button>

                    <button class="pswp__button pswp__button--fs" title="Toggle fullscreen"></button>

                    <button class="pswp__button pswp__button--zoom" title="Zoom in/out"></button>

                    <div class="pswp__preloader">
                        <div class="pswp__preloader__icn">
                            <div class="pswp__preloader__cut">
                                <div class="pswp__preloader__donut"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="pswp__share-modal pswp__share-modal--hidden pswp__single-tap">
                    <div class="pswp__share-tooltip"></div>
                </div>

                <button class="pswp__button pswp__button--arrow--left" title="Previous (arrow left)"></button>

                <button class="pswp__button pswp__button--arrow--right" title="Next (arrow right)"></button>

                <div class="pswp__caption">
                    <div class="pswp__caption__center"></div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';
import PhotoSwipe from 'photoswipe/dist/photoswipe';
import PhotoSwipeDefaultUI from 'photoswipe/dist/photoswipe-ui-default';
//import PhotoSwipeDefaultUI from './ui'

export default {
    methods: {
        open(index, items, options) {
            const opts = Object.assign(
                {
                    index: index,
                    getThumbBoundsFn: function(index) {
                        // Good guide on how to get element coordinates:
                        // http://javascript.info/tutorial/coordinates
                        let thumbnail = document.querySelectorAll('.preview-img-item .item')[index],
                            pageYScroll = window.pageYOffset || document.documentElement.scrollTop, // || document.body.scrollTop || 0;
                            rect = thumbnail.getBoundingClientRect();

                        return { x: rect.left, y: rect.top + pageYScroll, w: rect.width };
                    },

                    // captionEl: false,
                    fullscreenEl: true,
                    history: false,
                    shareEl: false,
                    // Tap on sliding area should close gallery
                    tapToClose: false,

                    // Tap should toggle visibility of controls
                    tapToToggleControls: true,
                    preload: [2, 2],
                },
                options
            );

            this.photoswipe = new PhotoSwipe(this.$el, PhotoSwipeDefaultUI, items, opts);
            this.photoswipe.init();
        },

        close() {
            this.photoswipe.close();
        },
    },
};
</script>
