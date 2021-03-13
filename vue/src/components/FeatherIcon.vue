<template functional>
    <span class="icon" :class="[data.class, data.staticClass, 'is-' + props.size]">
        <svg :class="props.sprite">
            <use :xlink:href="'/' + props.sprite + '.svg#' + props.name"></use>
        </svg>
    </span>
</template>



<script>
export default {
    name: 'FeatherIcon',
    props: {
      sprite: {
            required: false,
            type: String,
            default: 'feather',
        },
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