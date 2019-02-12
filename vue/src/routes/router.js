import Vue from "vue";
import Router from "vue-router";
import Home from "./Home.vue";

Vue.use(Router);

export default new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [{
      path: "/",
      name: "home",
      component: Home
    },
    {
      path: "/dashboard",
      name: "dashboard",
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () =>
        import( /* webpackChunkName: "about" */ "./Dashboard.vue")
    },
    {
      path: "/fotos",
      name: "alben",
      component: () => import("./Alben.vue")
    },
    {
      // will match anything starting with `/user-`
      path: '/filebox*'
    }

  ],
  linkActiveClass: 'is-active' /* change to Bulma's active nav link */
});