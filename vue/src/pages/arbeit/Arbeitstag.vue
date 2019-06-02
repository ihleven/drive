<template>
  <section class="section">
    <div class="container">
      {{ arbeitstag }} {{ date }}
      <div class="columns">
        <div class="column is-one-third">
          <div class="field">
            <label class="label">Status:</label>
            <div class="control">
              <div class="select">
                <select :value="status" @input="updateStatus">
                  <option value>---</option>
                  <option value="A">Arbeitstag</option>
                  <option value="U">Urlaub</option>
                  <option value="K">Krank</option>
                  <option value="W">Wochenende</option>
                  <option value="F">Feiertag</option>
                </select>
              </div>
            </div>
          </div>
        </div>
        <div class="column is-one-third">
          <div class="field">
            <label class="label">Typ:</label>
            <div class="control">
              <div class="select">
                <select v-model="typ">
                  <option value>---</option>
                  <option value="B">Büro</option>
                  <option value="H">Homeoffice</option>
                  <option value="D">Dienstreise</option>
                  <option value="K">Krank</option>
                </select>
              </div>
            </div>
          </div>
        </div>
        <div class="column is-one-third">
          <div class="field">
            <label class="label">Soll:</label>
            <div class="field has-addons">
              <div class="control has-icons-left">
                <input
                  class="input"
                  type="number"
                  placeholder="Sollarbeitszeit"
                  v-model.number="soll"
                >
                <feather-icon name="clock" size="small" class="is-left"/>
              </div>
              <p class="control">
                <a class="button is-light">
                  <feather-icon name="x-circle" @click="clear"/>
                </a>
              </p>
            </div>
          </div>
        </div>
      </div>
      <div class="columns">
        <div class="column is-one-third">
          <div class="field">
            <label class="label">Arbeitsbeginn:</label>
            <div class="control has-icons-left">
              <input class="input" type="text" placeholder="Arbeitsbeginn" v-model.lazy="start">
              <feather-icon name="clock" size="small" class="is-small is-left"/>
            </div>
          </div>
        </div>
        <div class="column is-one-third">
          <div class="field">
            <label class="label">Feierabend: {{ende}}</label>
            <b-timepicker placeholder="Type or select a date..." editable v-model="ende"></b-timepicker>
            <feather-icon name="clock" size="small" class="is-small is-left"/>
          </div>
        </div>
        <div class="column is-one-third">
          <div class="field">
            <label class="label">Bruttoarbeitszeit:</label>
            <div class="control">{{ arbeitstag.Brutto }}</div>
          </div>
        </div>
      </div>
      <div class="columns">
        <div class="column is-one-third">
          <h3 class="title is-3">Pausen</h3>
        </div>
        <div class="column is-one-third">
          <div class="field">
            <label class="label">Summe Pausen:</label>
            <div class="control">{{ arbeitstag.Pausen }}</div>
          </div>
        </div>
        <div class="column is-one-third">
          <div class="field">
            <label class="label">Nettoarbeitszeit:</label>
            <div class="control">{{ arbeitstag.Netto }}</div>
          </div>
        </div>
      </div>

      <div class="columns">
        <div class="column is-one-fifth">
          <div class="field">
            <label class="label">Typ:</label>
            <div class="control">
              <div class="select">
                <select>
                  <option>Pause</option>
                  <option>Restpausenabzug</option>
                  <option>Meeting</option>
                  <option>Homeoffice</option>
                  <option>???</option>
                </select>
              </div>
            </div>
          </div>
        </div>
        <div class="column is-one-fifth">
          <div class="field">
            <label class="label">Anfang:</label>
            <div class="field has-addons">
              <div class="control has-icons-left">
                <input class="input" type="text" placeholder="Text input">
                <feather-icon name="clock" size="small" class="is-left"/>
              </div>
              <p class="control">
                <a class="button is-light">
                  <feather-icon name="x-circle" @click="clear"/>
                </a>
              </p>
            </div>
          </div>
        </div>

        <div class="column is-one-fifth">
          <div class="field">
            <label class="label">Ende:</label>
            <div class="field has-addons">
              <div class="control has-icons-left">
                <input class="input" type="text" placeholder="Text input">
                <feather-icon name="clock" size="small" class="is-left"/>
              </div>
              <p class="control">
                <a class="button is-light">
                  <feather-icon name="x-circle" @click="clear"/>
                </a>
              </p>
            </div>
          </div>
        </div>

        <div class="column is-one-fifth"></div>

        <div class="column is-one-fifth"></div>
      </div>

      <div class="field">
        <label class="label">Email</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input is-danger" type="email" placeholder="Email input" value="hello@">
          <span class="icon is-small is-left">
            <i class="fas fa-envelope"></i>
          </span>
          <span class="icon is-small is-right">
            <i class="fas fa-exclamation-triangle"></i>
          </span>
        </div>
        <p class="help is-danger">This email is invalid</p>
      </div>

      <div class="field">
        <label class="label">Message</label>
        <div class="control">
          <textarea class="textarea" placeholder="Textarea"></textarea>
        </div>
      </div>

      <div class="field">
        <div class="control">
          <label class="checkbox">
            <input type="checkbox">
            I agree to the
            <a href="#">terms and conditions</a>
          </label>
        </div>
      </div>

      <div class="field">
        <div class="control">
          <label class="radio">
            <input type="radio" name="question">
            Yes
          </label>
          <label class="radio">
            <input type="radio" name="question">
            No
          </label>
        </div>
      </div>

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link">Submit</button>
        </div>
        <div class="control">
          <button class="button is-text">Cancel</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import { mapState } from 'vuex';

import FeatherIcon from '@/components/FeatherIcon.vue';

export default {
    name: 'Arbeitstag',
    components: {
        FeatherIcon,
    },
    data() {
        return {};
    },
    computed: {
        ...mapState(['activeDate', 'arbeitstag']),
        date() {
            return this.$store.state.activeDate;
        },
        ...mapState({
            status: state => state.arbeitstag.Status,
        }),
        typ: {
            get() {
                return this.$store.state.arbeitstag.Typ;
            },
            set(value) {
                this.$store.commit('updateArbeitstag', { field: 'Typ', value: value });
            },
        },
        soll: {
            get() {
                return this.$store.state.arbeitstag.Soll;
            },
            set(value) {
                this.$store.commit('updateArbeitstag', { field: 'Soll', value: value });
            },
        },
        start: {
            get() {
                let start = this.$store.state.arbeitstag.Start;

                if (start) {
                    let s = new Date(start),
                        hours = ('0' + s.getHours().toString()).slice(-2),
                        minutes = ('0' + s.getMinutes().toString()).slice(-2),
                        seconds = s.getSeconds(),
                        h_m = hours + ':' + minutes;
                    return h_m + (seconds ? ':' + ('0' + seconds.toString()).slice(-2) : '');
                }
                return null;
            },
            set(value) {
                let re = /^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])(:([0-9]|[0-5][0-9]))?$/,
                    match = value.match(re);
                if (match) {
                    let d = new Date(this.$store.state.arbeitstag.Start);
                    d.setMinutes(match[2]);
                    d.setHours(match[1]);
                    d.setSeconds(match[4] ? match[4] : 0);
                    d.setMilliseconds(0);
                    console.log(match);
                    this.$store.commit('updateArbeitstag', { field: 'Start', value: d.toISOString() });
                }
            },
        },
        ende: {
            get() {
                return new Date(this.$store.state.arbeitstag.Ende);
            },
            set(value) {
                this.$store.commit('updateArbeitstag', { field: 'Ende', value: value });
            },
        },
    },

    methods: {
        clear() {
            console.log('clear');
        },

        updateStatus(e) {
            //this.$store.commit('updateStatus', e.target.value);
            this.$store.commit('updateArbeitstag', { field: 'Status', value: e.target.value });
        },
    },

    beforeRouteUpdate(to, from, next) {
        //console.log('Arbeitstag In-Component Guard => beforeRouteUpdate:', to.params);
        this.$store.commit('SET_ACTIVE', to.params);
        this.$store.dispatch('loadArbeitstag');
        next();
    },
};
</script>

<style lang="stylus" scoped></style>
