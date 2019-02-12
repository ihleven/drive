<template>
  <div id="app" class="has-navbar-fixed-top">
    <nav class="navbar is-fixed-top is-primary" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a class="navbar-item" href="https://bulma.io">
          <img
            src="./assets/logo.png"
            alt="Bulma: Free, open source, & modern CSS framework based on Flexbox"
            width="112"
            height="28"
          >
          <a class="navbar-item">Home</a>
          <router-link class="navbar-item" to="/about">Home</router-link>
        </a>
        
        <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false">
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>
    </nav>

    <section class="hero is-primary is-with-navbar">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">Title</h1>
          <h2 class="subtitle">Subtitle</h2>
        </div>
      </div>

      <div class="hero-foot">
        <nav class="tabs is-boxed">
          <div class="container">
            <ul>
              <li>
                <router-link class="navbar-item" to="/">Home</router-link>
              </li>
              <li>
                <router-link class="navbar-item" to="/about">About</router-link>
              </li>
              <li>
                <router-link class="navbar-item" to="/fotos">Fotos</router-link>
              </li>
            </ul>
          </div>
        </nav>
      </div>
    </section>

    <section class="section">
      <transition>
        <router-view/>
      </transition>
    </section>

    <transition>
      <bug-report :showModal="openBugReport" v-on:close="openBugReport = false"></bug-report>
    </transition>
  </div>
</template>

