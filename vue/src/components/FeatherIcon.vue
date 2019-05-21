<template functional>
  <span class="icon" :class="[data.class, data.staticClass, 'is-' + props.size]">
    <svg class="feather">
      <use :xlink:href="'/feather.svg#' + props.name"></use>
    </svg>
  </span>
</template>

<script>
export default {
    name: 'FeatherIcon',
    props: {
        name: {
            required: true,
            type: String,
        },
        size: {
            required: false,
            type: String,
            default: 'normal',
        },
    },
};
</script>

<style lang="css">
.feather {
    fill: none;
    stroke: currentcolor;
    stroke-width: 1;
    stroke-linecap: round;
    stroke-linejoin: round;
}
</style>