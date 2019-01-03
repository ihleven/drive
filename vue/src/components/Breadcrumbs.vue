<template>
  <ul class="uk-breadcrumb">
    <li v-for="e in elems">
      <a href="/">{{e}}</a>
    </li>
  </ul>
</template>

<script>
export default {
  name: 'Breadcrumbs',
  props: {
    path: String
  },
  data: function() {
    return {
      elems: null
    }
  },
  mounted() {
    this.elems = this.path.split('/').filter(e => e != '')
    console.log('Hello World mounted', this.elems)
  }
}
</script>


<style scoped>
</style>
