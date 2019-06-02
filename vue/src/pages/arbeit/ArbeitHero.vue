<template>
  <section class="hero is-primary is-bold">
    <div class="hero-body">
      <div class="left">
        <div>
          <feather-icon name="clock"/>
        </div>
        <div class="wrapper">
          <span class="day" v-show="activeDate.day">{{ activeDate.day + '.' }}</span>
          <span class="monat">
            <router-link
              :to="{ name: 'arbeitsmonat', params: { year: activeDate.year, month:activeDate.month } }"
            >
              {{
              activeDate.monat
              }}
            </router-link>
            <router-link class="kw" :to="'/arbeit/' + week">{{ week }}</router-link>
            <router-link
              class="jahr"
              :to="{ name: 'arbeitsjahr', params: { year: activeDate.year } }"
            >{{ activeDate.year }}</router-link>
          </span>
        </div>
        <div class="wochentag">{{ activeDate.wochentag }}</div>
      </div>
      <div class="right">
        <div ref="cal"></div>
      </div>
    </div>
  </section>
</template>

<script>
import { mapState } from 'vuex';

import bulmaCalendar from 'bulma-calendar'; ///src/js/index.js';

export default {
    name: 'ArbeitHero',
    data() {
        return {
            week: 'KW23',
        };
    },
    computed: {
        ...mapState(['activeDate']),
    },
    created() {
        if (!this.activeDate.date) {
            this.$store.commit('SET_ACTIVE_DATE', new Date());
        }
    },
    mounted() {
        this.setupCalendar();
    },
    methods: {
        setupCalendar() {
            const calendar = new bulmaCalendar(this.$refs.cal, {
                startDate: this.activeDate.date,
                displayMode: 'inline',
                showHeader: false,
                showFooter: false,
                weekStart: 1,
                type: 'date',
            });
            calendar.on('select', this.calendarSelectEvent);
        },
        calendarSelectEvent(e) {
            let d = e.data.date.start,
                params = { year: d.getFullYear(), month: d.getMonth() + 1, day: d.getDate() };
            this.$router.push({ name: 'arbeitstag', ...{ params } });
        },
    },
};
</script>

<style lang="stylus" scoped>
.hero-body {
  padding: 2rem;
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  justify-content: space-around;
}

.left {
  padding-left: 3rem;
}

.day {
  background-color: rgba(100, 30, 50, 0);
  font: ultra-condensed normal 900 6rem / 1 'Raleway';
}

.wrapper {
  display: inline-block;
  background-color: rgba(100, 30, 50, 0);
  position: relative;
}

.monat {
  position: relative;
  top: -2rem;
  left: -1.4rem;
  background-color: rgba(0, 150, 20, 0);
  font: 600 2rem / 1 'Raleway';
}

.kw {
  position: relative;
  top: -1rem;
  left: 1rem;
  background-color: rgba(0, 150, 20, 0);
  font: 300 1.5rem / 1 'Raleway';
}

.jahr {
  position: absolute;
  top: 1.8rem;
  left: 2rem;
  background-color: rgba(0, 150, 20, 0);
  font: 700 2.5rem / 1 'Raleway';
}

.wochentag {
  font: 600 2rem / 1 'Raleway';
}
</style>
