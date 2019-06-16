import Vue from 'vue';

import bytes from './bytes.js';
import timeformat from './timeformat.js';

Vue.filter('bytes', bytes);
Vue.filter('timeformat', timeformat);
