<script>
export default {
    props: {
        name: {},
        size: {
            default: 'normal',
        },
        modifier: {
            default: null,
        },
        fixAlign: {
            default: true,
        },
    },
};
</script>

<template>
    <div class="icon" :class="[`icon--${name}`, `icon--${size}`, { 'has-align-fix': fixAlign }]">
        <svg class="icon__svg">
            <use xmlns:xlink="http://www.w3.org/1999/xlink" :xlink:href="'/feather.svg#' + name"></use>
        </svg>
    </div>
</template>

<style lang="scss" scoped>
.icon {
    position: relative;
    display: inline-block;
    vertical-align: middle;
    width: 1rem;
    height: 1rem;
    margin: 0 0.3rem;
    top: -0.05rem;
    fill: currentColor;

    // &.has-align-fix {
    // 	top: -.1rem;
    // }

    &__svg {
        display: inline-block;
        vertical-align: top;
        width: 100%;
        height: 100%;
    }

    &:first-child {
        margin-left: 0;
    }

    &:last-child {
        margin-right: 0;
    }
}

// svg sprite
body > svg,
.icon use > svg,
symbol {
    path,
    rect,
    circle,
    g {
        fill: none;
        stroke: currentColor;
    }

    *[d='M0 0h24v24H0z'] {
        display: none;
    }
}
</style>
