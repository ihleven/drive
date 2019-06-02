<template>
  <div>
    <section class="section has-background-light" v-if="!activeDate.day">
      <div class="container">
        <feather-icon name="sun"/>
        Monat: {{ activeDate.monat }}, {{ year }}, {{ month }}
      </div>
    </section>
    <router-view></router-view>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import FeatherIcon from '@/components/FeatherIcon.vue';

export default {
    name: 'Arbeitsmonat',
    components: {
        FeatherIcon,
    },
    props: ['year', 'month'],
    computed: {
        ...mapState(['activeDate']),
    },

    created() {
        //console.log('Arbeitsmonat => created:', this.year, this.month);
    },

    beforeRouteUpdate(to, from, next) {
        if (!to.params.day) {
            this.$store.commit('SET_ACTIVE', {
                year: to.params.year,
                month: to.params.month,
            });
        }
        //this.$store.dispatch('loadArbeitMonat', to.params);
        if (to.params.month != from.params.month) {
            //this.month = to.params.month;
            //console.log('Arbeitsmonat => beforeRouteUpdate:', to.params.month);
        }
        next();
    },
};
</script>

<style lang="stylus" scoped></style>
