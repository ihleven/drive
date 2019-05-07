<script>
export default {
    name: 'Cloud11Figure',
    props: {
        src: {
            required: true,
            type: String,
        },
        type: {
            required: false,
            type: String,
        },
        tags: {
            required: false,
            type: Array,
        },
    },
};
</script>

<template functional>
    <div class="card">
        <div class="card-image">
            <figure class="image">
                <img :src="props.src" />
            </figure>
        </div>
        <div class="card-content">
            <button class="button" @click="listeners.select">Select</button>
            <span v-if="$slots.default">
                <slot></slot>
            </span>
            {{ props.type }}
            <div class="tags" v-if="props.tags && props.type != 'framed'">
                <span v-for="tag in props.tags" :key="tag">{{ tag }}</span>
            </div>
        </div>
    </div>
</template>

<style scoped>
.card {
    display: inline-block;
    margin: 1rem;
    max-width: 300px;
}
</style>
