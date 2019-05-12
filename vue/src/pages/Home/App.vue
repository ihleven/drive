<script>
export default {
    name: 'Clou11',
    data() {
        return {
            menuOpen: false,
        };
    },
    methods: {
        overlayNavigationToggler() {
            this.menuOpen = !this.menuOpen;
        },
    },
    created() {},
};
</script>

<template>
    <div class="application-wrapper">
        <nav class="navigation" :class="{ open: menuOpen }">
            <section class="hero is-dark is-fullheight">
                <div class="hero-body">
                    <div class="container">
                        <h3 class="title">Mallorca</h3>
                        <a class="subtitle" href="/alben/Mallorca">Mallorca</a>
                        <a class="subtitle">Hero subtitle</a>
                    </div>
                </div>
            </section>
        </nav>

        <button class="navigation-trigger" @click="overlayNavigationToggler">
            <svg class="feather-fat">
                <use xlink:href="/feather.svg#x"></use>
            </svg>
        </button>

        <div class="MainContainer">
            <section class="hero is-medium parallax">
                <div class="hero-body">
                    <div class="container">
                        <h1 class="title">Primary bold title</h1>
                        <h2 class="subtitle">Primary bold subtitle</h2>
                    </div>
                </div>
            </section>

            <div class="ContentContainer">
                <section class="section is-transparent" style="background:transparent;position:absolute;top:-200px">
                    <div class="container">
                        <h1 class="title">Section</h1>
                        <h2 class="subtitle">
                            A simple container to divide your page into
                            <strong>sections</strong>, like the one you're currently reading
                        </h2>
                    </div>
                </section>
                <section class="section has-background-primary">
                    <div class="container">
                        <div class="content">
                            <p>
                                Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam nunc tortor, porttitor
                                in nisl quis, commodo laoreet odio. Pellentesque eu pulvinar lacus. Curabitur venenatis
                                vel diam vitae molestie. Suspendisse ultricies mi a elit ultricies consequat. Ut
                                pharetra rhoncus facilisis. Integer scelerisque malesuada justo sit amet sollicitudin.
                                Aenean consectetur molestie dolor et gravida. Pellentesque nec felis in nunc aliquam
                                porta. Suspendisse laoreet justo vitae scelerisque faucibus.
                            </p>

                            <p>
                                Aliquam condimentum ex sem, eu fermentum nisi malesuada non. Vestibulum ullamcorper
                                neque non tellus molestie ultrices. Nullam vel sapien finibus, egestas velit vitae,
                                porttitor tellus. Vivamus ut eleifend turpis. Nulla eu enim at ligula luctus porta sed
                                id arcu. Nullam vitae convallis elit. In diam justo, aliquet ultrices elementum vitae,
                                aliquet et augue. Sed eu turpis nec nibh condimentum ornare ut at justo. Nulla hendrerit
                                blandit dictum. Nullam cursus sit amet est quis sollicitudin.
                            </p>

                            <p>
                                Mauris eu cursus sapien, quis egestas risus. In sit amet tempus velit. Proin in dui ut
                                ex consectetur laoreet ornare at leo. Nunc tincidunt, lorem sed fermentum pellentesque,
                                ante quam mollis massa, ut finibus ligula ipsum id augue. Etiam et feugiat sem, at
                                elementum erat. Quisque neque nibh, auctor dictum nunc sit amet, finibus condimentum
                                justo. Suspendisse volutpat sapien volutpat dui eleifend, ut ultricies lacus placerat.
                            </p>

                            <p>
                                Nulla sodales nec tellus a egestas. Donec vitae nulla viverra, molestie leo ut,
                                condimentum augue. In hac habitasse platea dictumst. Nunc a neque rhoncus, ultricies
                                nisl ut, tincidunt orci. Sed sollicitudin augue vitae erat commodo convallis. In vel
                                lectus feugiat, interdum elit aliquet, vehicula nunc. Vestibulum non nisl purus. Nunc
                                dapibus lobortis nibh non fermentum. Fusce sagittis laoreet justo id aliquam. Nulla sit
                                amet enim in neque ornare consequat. Vestibulum accumsan auctor elementum. Ut tincidunt
                                augue eget felis aliquet, ut malesuada dolor vulputate. Proin bibendum dui non ipsum
                                viverra auctor. Pellentesque habitant morbi tristique senectus et netus et malesuada
                                fames ac turpis egestas.
                            </p>

                            <p>
                                In id lacus nec magna imperdiet ullamcorper. Etiam in finibus turpis. Suspendisse
                                bibendum faucibus posuere. Phasellus vulputate mi sed neque accumsan, id tincidunt magna
                                scelerisque. Nunc sed felis nulla. Aliquam molestie risus odio, ut accumsan tortor
                                consectetur a. Duis at nisi sed nunc porta blandit. Phasellus varius ex ante, eu
                                accumsan ex mattis in. Aliquam erat volutpat. Nunc iaculis ex quis odio lacinia, id
                                condimentum enim ultrices. Proin tempus fringilla magna rhoncus hendrerit. Praesent
                                lorem ante, imperdiet in metus porta, finibus rutrum urna. Etiam tincidunt porta tellus,
                                facilisis consectetur diam maximus eu. Etiam elementum sed ipsum nec mattis. Nullam ut
                                sapien elit.
                            </p>
                            <ol type="1">
                                <li>Coffee</li>
                                <li>Tea</li>
                                <li>Milk</li>
                            </ol>
                            <ol type="A">
                                <li>Coffee</li>
                                <li>Tea</li>
                                <li>Milk</li>
                            </ol>
                            <ol type="a">
                                <li>Coffee</li>
                                <li>Tea</li>
                                <li>Milk</li>
                            </ol>
                            <ol type="I">
                                <li>Coffee</li>
                                <li>Tea</li>
                                <li>Milk</li>
                            </ol>
                            <ol type="i">
                                <li>Coffee</li>
                                <li>Tea</li>
                                <li>Milk</li>
                            </ol>
                        </div>
                    </div>
                </section>
            </div>
        </div>
    </div>
</template>

<style lang="css">
html,
body,
.application-wrapper {
    max-height: 100vh;
    height: 100vh;
    overflow: hidden !important;
}

svg {
    width: 100%;
    height: 100%;
}
</style>

<style lang="scss">
.MainContainer {
    perspective: 1px;
    perspective-origin: top left;
    transform-style: preserve-3d;
    height: 100vh;
    max-height: 100vh;
    overflow-x: hidden;
    overflow-y: scroll;
    -webkit-overflow-scrolling: touch;
}

.parallax {
    height: 66vh;
    z-index: -1;
    position: relative;
    transform: translateZ(-1px) scale(2);
    transform-origin: 0 0;

    background: url(https://www.toptal.com/designers/subtlepatterns/patterns/sakura.png);
    background-color: rgb(250, 228, 216);
    background-position: center 0%;
    background-repeat: no-repeat;
    /*background-attachment: fixed;*/
    background-size: 100vw auto;
    background-color: #999;
    border: 1px solid green;
}

.ContentContainer {
    display: block;
    position: relative;
    z-index: 1;
    border: 1px solid red;
}

.navigation {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100vw;
    height: 100vh;
    text-align: right;
    opacity: 0;
    pointer-events: none;
    z-index: 1;
    line-height: 1;
    transition: all 0.3s ease-in-out 0s;
}

.navigation.open {
    opacity: 1;
    pointer-events: auto;
    overflow: hidden;
}
.navigation-trigger {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    padding: 1rem;
    z-index: 1000000;
    cursor: pointer;
    transition-property: opacity, filter;
    transition-duration: 0.15s;
    transition-timing-function: linear;
    background-color: transparent;
    border: 0;
    margin: 0;
    outline: 0;
    height: 3rem;
    width: 3rem;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    box-sizing: content-box;
}

.feather-fat {
    fill: none;
    stroke: grey;
    stroke-width: 2.5;
    stroke-linecap: round;
    stroke-linejoin: round;
}
</style>
