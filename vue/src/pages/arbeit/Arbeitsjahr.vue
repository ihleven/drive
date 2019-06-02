<template>
  <div>
    <section class="section has-background-grey-lighter" v-if="!activeDate.month">
      <div class="container">
        <feather-icon name="star"/>
        Jahr: {{ activeDate.year }}
        <br>
        {{ arbeitsjahr }}
      </div>
    </section>
    <router-view></router-view>
  </div>
</template>

<script>
import { mapState } from 'vuex';

import FeatherIcon from '@/components/FeatherIcon.vue';

export default {
    name: 'Arbeitsjahr',
    components: {
        FeatherIcon,
    },
    computed: {
        ...mapState(['activeDate', 'arbeitsjahr']),
    },

    created() {
        //console.log('Arbeitsjahr => created:', this.activeDate.year);
    },
    beforeRouteUpdate(to, from, next) {
        if (!to.params.month) {
            this.$store.commit('SET_ACTIVE', {
                year: to.params.year,
            });
            this.$store.dispatch('loadArbeitJahr', to.params);
        }
        next();
    },
};
</script>

<style lang="stylus" scoped></style>
