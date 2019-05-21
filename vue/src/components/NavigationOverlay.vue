<script>
export default {
    name: 'NavigationOverlay',
    props: {
        open: Boolean,
    },
    methods: {
        close() {
            this.$emit('update:open', false);
        },
    },
};
</script>

<template>
    <nav class="navigation" :class="{ open }">
        <section class="hero is-dark is-bold is-fullheight">
            <div class="hero-head">
                <nav class="navbar is-transparent">
                    <a class="navbar-item hamburger" @click="close">
                        <svg>
                            <use xlink:href="/feather.svg#x"></use>
                        </svg>
                    </a>
                </nav>
            </div>
            <div class="hero-body">
                <div class="container">
                    <div class="columns">
                        <div class="column content">
                            <h3>ihle.cloud</h3>
                            <ul>
                                <li>
                                    <a href="/home">home</a>
                                </li>
                                <li>
                                    <a href="/public">public</a>
                                </li>
                            </ul>
                        </div>
                        <div class="column">
                            <slot></slot>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    </nav>
</template>

<style lang="css" scoped>
.navigation {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    background: #fff;
    width: 100vw;
    height: 100vh;
    text-align: left;
    opacity: 0;
    pointer-events: none;
    z-index: 1000;
    line-height: 1;
    transition: all 0.3s ease-in-out 0s;
    overflow-y: scroll;
}

.navigation.open {
    opacity: 1;
    pointer-events: auto;
}

.navigation .navbar {
    justify-content: flex-end;
}
.navigation .navbar .hamburger:hover {
    background-color: rgba(255, 255, 255, 0.5);
    padding: 0.25rem 0.5rem;
    transition: all 0.2s ease-in-out 0s;
}
.navigation .navbar .hamburger svg {
    stroke: white;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
    width: 2rem;
    height: 2rem;
    fill: none;
}
.navigation .navbar .hamburger:hover svg {
    stroke-width: 2.5;
    width: 2.5rem;
    height: 2.5rem;
    transition: all 0.2s ease-in-out 0s;
}
</style>
