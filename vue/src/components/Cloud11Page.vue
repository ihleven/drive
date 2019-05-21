<template>
  <div class="application-wrapper">
    <navigation-overlay :open.sync="menuOpen"></navigation-overlay>

    <navbar-top :account="account" @click:hamburger="overlayNavigationToggler()"/>

    <slot></slot>
  </div>
</template>

<script>
import NavigationOverlay from '@/components/NavigationOverlay.vue';
import NavbarTop from '@/components/NavbarTop.vue';

export default {
    name: 'Cloud11Page',
    components: {
        NavigationOverlay,
        NavbarTop,
    },
    props: {
        account: Object,
    },
    data() {
        return {
            menuOpen: false,
        };
    },
    methods: {
        overlayNavigationToggler() {
            this.menuOpen = !this.menuOpen;
        },
    },
};
</script>

<style lang="css">
html,
body,
.application-wrapper {
    max-height: 100vh;
    height: 100vh;
    overflow: hidden !important;
}
</style>
            menuOpen:false,
