import Vue from 'vue';

import bytes from './bytes.js';
import symbolic from './symbolic.js';
import timeformat from './timeformat.js';

Vue.filter('bytes', bytes);
Vue.filter('symbolic', symbolic);
Vue.filter('timeformat', timeformat);
