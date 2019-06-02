<script>
export default {
    name: 'NavbarTop',
    props: {
        account: { type: Object, default: () => new Object() },
    },
    methods: {
        close() {
            this.$emit('update:open', false);
        },
        hamburgerClick() {
            this.$emit('click:hamburger');
        },
    },
};
</script>

<template>
  <nav class="navbar is-fixed-top top" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <a class="navbar-item" href="/" style="padding-right:0;">
        <small class="login">matthias</small>
        <svg class="cloud-icon" style="width: 1em;height: 1em;margin-right: .125em;">
          <use xlink:href="/feather.svg#at-sign"></use>
        </svg>
      </a>
      <a class="navbar-item" href="/" style="padding-left:0;">
        <span class="domain">ihle.</span>
        <svg class="cloud-icon">
          <use xlink:href="/feather.svg#cloud"></use>
        </svg>
      </a>
    </div>
    <a class="navbar-item hamburger" @click="hamburgerClick">
      <svg class>
        <use xlink:href="/feather.svg#menu"></use>
      </svg>
    </a>
  </nav>
</template>

<style lang="scss" scoped>
.navbar.top {
    max-width: 100%;
    margin: 0 auto;
    background-color: transparent;
    justify-content: space-between;
    margin-right: 15px;
    border-right: 0px solid black;
}

.navbar.top .login {
    color: white;
    font-size: 1rem;
    font-weight: 200;
}
.navbar.top .domain {
    color: white;
    font-size: 1.25rem;
    font-weight: 400;
}
.navbar.top .cloud-icon {
    width: 1.5rem;
    height: 1.5rem;
    fill: none;
    stroke: white;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.navbar.top .hamburger:hover {
    background-color: rgba(255, 255, 255, 0.26) !important;
    padding: 0.25rem 0.5rem;
    transition: all 0.2s ease-in-out 0s;
}
.navbar.top .hamburger svg {
    stroke: white;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
    width: 2rem;
    height: 2rem;
    fill: none;
}
.navbar.top .hamburger:hover svg {
    stroke-width: 2.5;
    width: 2.5rem;
    height: 2.5rem;
    transition: all 0.2s ease-in-out 0s;
}
</style>
