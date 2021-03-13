<template>
    <div>
        <section class="section has-background-light" v-if="!activeDate.day">
            <div class="container">
                <feather-icon name="sun"/>
                Monat: {{ activeDate.monat }}, {{ year }}, {{ month }}
            </div>
        </section>
        <router-view></router-view>

        <section class="section">
            <div class="container">
                <table class="table is-striped is-narrow is-hoverable is-fullwidth">
                
                    <router-link :to="{ name: 'arbeitstag', params: a.kalendertag }" tag="tr" 
                        v-for="a in arbeitstage" :key="a.id"
                        class="day" :class="a.status + ' ' + a.kategorie" >
                        <td class="date">{{ new Date(a.kalendertag.Datum).toLocaleDateString('de-DE', { month: 'short', day: 'numeric' }) }}</td>
                        <td class="date">{{ new Date(a.kalendertag.Datum).toLocaleDateString('de-DE', { weekday: 'short' }) }}</td>
                        <td>{{a.status == 'A' ? 1 : a.status=='H' ? "&#189;" : a.status}}</td>
                        <td>{{a.kategorie}}</td>
                        <td v-text="a.Urlaubstage || ''"></td>
                        <td v-text="a.Soll || ''"></td>
                        <td><span v-if="a.beginn">{{a.beginn.toLocaleTimeString('de')}}</span></td>
                        <td><span v-if="a.ende">{{a.ende.toLocaleTimeString('de')}}</span></td>
                        <td v-text="a.Brutto"></td>
                        <td v-text="a.Pausen"></td>
                        <td v-text="a.Extra"></td>
                        <td v-text="a.Netto"></td>
                        <td v-text="a.Differenz"></td>
                        <td>{{a.Saldo}}</td>
                    </router-link>
                
                </table>
      
            </div>
        </section>
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
        ...mapState(['activeDate', 'arbeitstage']),
    },

    created() {
        //console.log('Arbeitsmonat => created:', this.year, this.month);
    },
    beforeRouteUpdate(to, from, next) {
        //console.log('Arbeitstag In-Component Guard => beforeRouteUpdate:', to.params);

        if (!to.params.day) {
            this.$store.commit('SET_ACTIVE', {
                year: to.params.year,
                month: to.params.month,
            });
        }
        this.$store.dispatch('loadArbeitMonat', to.params);
        if (to.params.month != from.params.month) {
            //this.month = to.params.month;
            //console.log('Arbeitsmonat => beforeRouteUpdate:', to.params.month);
        }
        next();
    },
};
</script>

<style lang="stylus" scoped>
    .day {
        &:hover {
            color: hsl(0, 0%, 20%);
            background-color: hsl(0, 0%, 85%);
            border-bottom: 1px solid hsl(0, 0%, 10%);
            cursor: pointer;
        }
    }
    .date {
        white-space: nowrap;
    }
    .U {
        background-color: hsl(141, 71%, 48%);
    }
    .Z {
        background-color: #00d1b2;
    }
    .K {
        background-color: hsl(48, 100%, 67%);
    }
    .W {
        background-color: hsl(348, 100%, 85%);
    }
    .F, .S {
        background-color: hsl(348, 100%, 91%);
    }
</style>
