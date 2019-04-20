<template>
    <transition name="modal-fade">
        <div class="modal" :class="{ 'is-active': visible }">
            <div class="modal-background"></div>
            <div class="modal-content">
                <div class="box dropbox">
                    <article class="media" v-for="file in files" :key="file.id">
                        <div class="media-left">
                            <figure class="image is-128x128">
                                <img src="https://bulma.io/images/placeholders/128x128.png" :ref="file.id" />
                            </figure>
                        </div>
                        <div class="media-content">
                            <div class="content">
                                <strong>{{ file.file.name }}</strong>
                                <br />
                                <small>{{ file.file.size }}</small>
                                <br />
                                <small>{{ file.file.lastModified }}</small>
                            </div>
                            <button @click="upload(file.file)">Upload</button>
                        </div>
                    </article>
                    <input
                        type="file"
                        multiple
                        :name="uploadFieldName"
                        :disabled="isSaving"
                        @change="onFileSelect"
                        accept="image/*"
                        class="input-file"
                    />
                    <p>
                        Drag your file(s) here to begin
                        <br />or click to browse
                    </p>
                    <p v-if="isSaving">Uploading {{ fileCount }} files...</p>

                    <button class="button">Cancel</button>
                </div>
            </div>
            <button class="modal-close is-large" aria-label="close" @click="close">close</button>
        </div>
    </transition>
</template>

<script>
export default {
    name: 'UploadModal',
    props: {
        visible: Boolean,
        url: String,
    },
    data() {
        return { files: [] };
    },
    methods: {
        close() {
            this.$emit('update:visible', false);
        },
        onFileSelect(event) {
            let files = event.target.files;
            for (let i = 0; i < files.length; i++) {
                let file = files[i],
                    id = this.id();

                console.group('File ' + i);
                console.log('name: ' + file.name, this.$refs.preview);
                console.groupEnd();
                this.files.push({
                    id: id,
                    file: file,
                });
                this.previewImage(file, id);
            }
        },
        previewImage(file, refid) {
            if (file) {
                var reader = new FileReader();
                reader.addEventListener(
                    'load',
                    () => {
                        this.$refs[refid][0].src = reader.result;
                    },
                    false
                );
                reader.readAsDataURL(file);
            }
        },
        id() {
            // Math.random should be unique because of its seeding algorithm.
            // Convert it to base 36 (numbers + letters), and grab the first 9 characters
            // after the decimal.
            return (
                '_' +
                Math.random()
                    .toString(36)
                    .substr(2, 9)
            );
        },
        upload(file) {
            let formData = new FormData();
            formData.append('file', file);
            const request = new Request(this.url, {
                method: 'POST',
                body: formData,
                headers: new Headers({ Accept: 'application/json' }),
            });

            fetch(request)
                .then(response => response.json())
                .then(body => console.log(body))
                .catch(function() {
                    console.log('FAILURE!!');
                });

            /* axios
                .post(this.folder.path, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                    },
                })
                .then(function() {
                    console.log('SUCCESS!!');
                }) */
        },
    },
    mounted() {
        console.log('upload modal mounted', this.url);
    },
};
</script>

<style>
.modal-fade-enter,
.modal-fade-leave-active {
    opacity: 0;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
    transition: opacity 0.5s ease;
}
</style>

<style lang="scss">
.dropbox {
    outline: 2px dashed grey; /* the dash box */
    outline-offset: -10px;
    background: lightcyan;
    color: dimgray;
    padding: 20px 20px;
    min-height: 200px; /* minimum height */
    position: relative;
    overflow-x: hidden;
}

.input-file {
    opacity: 0; /* invisible but it's there! */
    width: 100%;
    height: 200px;
    position: absolute;
    cursor: pointer;
}

.dropbox .input-file:hover {
    background: lightblue; /* when mouse over to the drop zone, change color */
}

.dropbox p {
    font-size: 1.2em;
    text-align: center;
    padding: 50px 0;
}
</style>
