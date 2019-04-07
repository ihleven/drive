import Vue from 'vue';

import Prism from 'prismjs';
//import 'prismjs/themes/prism.css';

const highlighter = {
    // When the bound element is inserted into the DOM...
    inserted: function(element, binding, vnode) {
        Prism.highlightElement(element, false, console.log);

        //Prism.highlight(text, grammar, language)
        console.log('highlighter');
    },
};

Vue.directive('highlight', highlighter);

export default highlighter;
