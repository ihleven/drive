<template>
  <div class="application-wrapper">
    <arbeit-hero/>
    <router-view></router-view>
  </div>
</template>

<script>
import ArbeitHero from './ArbeitHero.vue';

export default {
    name: 'ArbeitApp',
    components: {
        ArbeitHero,
    },
    data() {
        return {
            menuOpen: false,
        };
    },
    methods: {},
    created() {
        //console.log('created ArbeitApp component');
    },
};
</script>
