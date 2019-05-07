<script>
export default {
    name: 'AlbumImage',
    props: {
        src: {
            required: true,
            type: String,
        },
    },
    mounted() {
        console.log('AlbumImage', this.src);
    },
};
</script>

<template functional>
    <img class="thumb" :src="props.src" />
</template>

<style scoped>
.thumb {
    display: inline-block;
    margin: 1rem;
    max-width: 300px;
}
</style>
