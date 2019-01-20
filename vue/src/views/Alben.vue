<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <HelloWorld msg="Welcome to Your Vue.js App"/>
    {{ foo }} |
    <a @click="load()">qwer</a> |
    <a href="/tmp">tmp</a>
  </div>
</template>

<script>
import axios from 'axios'
import HelloWorld from '@/components/HelloWorld.vue'

export default {
    name: 'alben',
    components: { HelloWorld },
    data: function() {
        return {
            foo: 'bar',
        }
    },
    mounted() {
        console.log('mounted Alben.vue')
        axios({
            method: 'get',
            url: '/tmp/14/DSC02007.jpg',
            headers: {
                Accept: 'application/json',
                foo: 'bar',
            },
            data: {
                firstName: 'Fred',
                lastName: 'Flintstone',
            },
        })
    },
    methods: {
        load: function() {
            console.log('load')
            const instance = axios.create({
                baseURL: '/tmp/14/DSC02007.jpg',
                timeout: 1000,
                headers: { Accept: 'application/json', foo: 'bar' },
            })

            instance
                .get('/tmp/14/DSC02007.jpg')
                .then(function(response) {
                    console.log(response)
                })
                .catch(function(error) {
                    console.log(error)
                })
                .then(function() {
                    // always executed
                })
        },
    },
}
</script>
