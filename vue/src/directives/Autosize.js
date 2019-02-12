import autosize from 'autosize'

export default {
    inserted: function (el) {
        console.log('inserted:', el)
        var tagName = el.tagName
        if (tagName == 'TEXTAREA') {
            autosize(el)
        }
    },
    update: function (el) {
        console.log('update:', el)
        if (el.tagName == 'TEXTAREA') {
            autosize.update(el)
        }
    },
    componentUpdated: function (el) {
        console.log('componentUpdated:', el)
        var tagName = el.tagName
        if (tagName == 'TEXTAREA') {
            autosize.update(el)
        }
    },
    unbind: function (el) {
        autosize.destroy(el)
    }
}