import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from "./components/Home";

Vue.use(VueRouter);

export default new VueRouter({
    routes: [
        {
            path: '*',
            name: 'home',
            component: Home,
            meta: {
                keepAlive: false
            }
        },
    ]
})
