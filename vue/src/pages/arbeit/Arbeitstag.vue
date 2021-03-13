<template>
    <section class="section">
        <div class="container">
            <div class="field is-grouped">
                <div class="control">
                    <button class="button is-link" @click=save>Submit</button>
                </div>
                <div class="control">
                    <button class="button is-text">Cancel</button>
                </div>
            </div>

            <div class="columns">
                <div class="column is-one-third">
                <div class="field">
                    <label class="label">Status: {{status}}</label>
                    <div class="control">
                    <div class="select">
                        <select :value="status" @input="updateStatus">
                        <option value="">W</option>
                        <option value="A">Arbeitstag</option>
                        <option value="H">Arbeitstag (halb)</option>
                        <option value="S">Sonntag</option>
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
                        <option value="Z">Zeitausgleich</option>
                        <option value="U">Urlaub</option>
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
                        <label class="label">Arbeitsbeginn: {{arbeitstag.beginn}}</label>
                        <div class="control has-icons-left">
                            <!--<input class="input" type="text" placeholder="Arbeitsbeginn" v-model.lazy="start">-->
                            <b-timepicker  editable v-model="start"></b-timepicker>
                            <feather-icon name="clock" size="small" class="is-small is-left"/>
                        </div>
                    </div>
                </div>
                <div class="column is-one-third">
                    <div class="field">
                        <label class="label">Feierabend: {{arbeitstag.ende}}</label>
                        <div class="control has-icons-left">
                            <b-timepicker  editable v-model="end"></b-timepicker>
                            <feather-icon name="clock" size="small" class="is-small is-left"/>
                        </div>
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
                <button @click="addPause"></button>
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

            <div class="columns" v-for="z in arbeitstag.Zeitspannen">
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
                                <input class="input" type="text" placeholder="Text input" v-model.lazy="z.Von">
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
                                <input class="input" type="text" v-model.lazy="z.Bis">
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
                <label class="label">Kommentar</label>
                <div class="control">
                <textarea class="textarea" placeholder="Kommentar"></textarea>
                </div>
            </div>

            <pre>
            {{ arbeitstag }}</pre>
            <pre>
            {{ date }}
            </pre>
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
            status: state => state.arbeitstag.status,
        }),
        typ: {
            get() {
                return this.$store.state.arbeitstag.kategorie;
            },
            set(value) {
                this.$store.commit('updateArbeitstag', { field: 'kategorie', value: value });
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
        startalt: {
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
        start: {
            get() {
                return this.$store.state.arbeitstag.beginn; // new Date(this.$store.state.arbeitstag.beginn);
            },
            set(value) {
                this.$store.commit('updateArbeitstag', { field: 'beginn', value: value });
            },
        },
        end: {
            get() {
                return this.$store.state.arbeitstag.ende; //new Date(this.$store.state.arbeitstag.ende);
            },
            set(value) {
                console.log("ende:", value);
                this.$store.commit('updateArbeitstag', { field: 'ende', value: value });
            },
        },
    },

    beforeRouteUpdate(to, from, next) {
        //console.log('Arbeitstag In-Component Guard => beforeRouteUpdate:', to.params);
        this.$store.commit('SET_ACTIVE', to.params);
        this.$store.dispatch('loadArbeitstag');
        next();
    },

    methods: {
        clear() {
            console.log('clear');
        },
        save() {
            if (typeof this.$store.state.arbeitstag.beginn === 'object' && this.$store.state.arbeitstag.beginn) {
                let b = this.$store.state.arbeitstag.beginn;
                    b.setFullYear(this.$store.state.arbeitstag.kalendertag.year);
                    b.setMonth(this.$store.state.arbeitstag.kalendertag.month);
                    b.setDate(this.$store.state.arbeitstag.kalendertag.day);
            }
            if (typeof this.$store.state.arbeitstag.ende === 'object' && this.$store.state.arbeitstag.ende) {
                let e = this.$store.state.arbeitstag.ende;
                    e.setFullYear(this.$store.state.arbeitstag.kalendertag.year);
                    e.setMonth(this.$store.state.arbeitstag.kalendertag.month);
                    e.setDate(this.$store.state.arbeitstag.kalendertag.day);
            }
            console.log("save:", this.$store.state.arbeitstag, typeof this.$store.state.arbeitstag.beginn);
            this.$store.dispatch('saveArbeitstag');
        },
        updateStatus(e) {
            //this.$store.commit('updateStatus', e.target.value);
            this.$store.commit('updateArbeitstag', { field: 'Status', value: e.target.value });
        },
        addPause() {
            //this.$store.commit('updateStatus', e.target.value);
            //let zeitspannen = [{Nr: 1, id: 20180506001}]
            this.$store.commit('updateArbeitstag', { 
                field: 'zeitspannen', 
                value: this.$store.state.arbeitstag.Zeitspannen.push({
                    Nr: this.$store.state.arbeitstag.Zeitspannen.length + 1
                })
            });
        },
    },

};
</script>

<style lang="stylus" scoped></style>
