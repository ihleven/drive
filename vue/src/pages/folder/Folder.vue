/* eslint-disable prettier/prettier */
<template>
  <section class="section">
    <button type="button" class="btn" @click="showModal">Open Modal!</button>
    <upload-modal :visible.sync="isModalVisible" :url="folder.path"/>

    <div class="container" v-if="folder">
      <table class="table is-striped is-hoverable is-fullwidth">
        <thead>
          <tr>
            <th></th>
            <th>Name</th>

            <th colspan="2" class="has-text-centered">MIME</th>
            <th>Mode</th>
            <th class="has-text-right">Owner</th>
            <th>Group</th>
            <th>Perm</th>
            <th>Size</th>
            <th>Created</th>
            <th>Modified</th>
            <th>Accessed</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="f in folder.entries" :key="f.name">
            <td>
              <span class="icon is-small">
                <svg class="feather" v-if="f.mime.Type == 'dir'">
                  <use xlink:href="/feather.svg#folder"></use>
                </svg>
                <svg class="feather" v-else>
                  <use xlink:href="/feather.svg#file"></use>
                </svg>
              </span>
            </td>
            <td>
              <a :href="f.path" :title="f.path">{{ f.name }}</a>
            </td>

            <td class="has-text-right" style="padding-right:0">{{ f.mime.Type }}</td>
            <td class="nowrap">{{ f.mime.Subtype }}</td>
            <td class="nowrap">{{ f.Mode }}</td>
            <td
              class="has-text-right"
              :class="f.permissions.IsOwner ? 'is-user-group' : 'is-not-user-group'"
            >{{ f.owner.name }}</td>
            <td
              :class="f.permissions.IsOwner ? 'is-user-group' : 'is-not-user-group'"
            >{{ f.group.Name }}</td>
            <td>
              <span class="icon is-small">
                <svg class="feather" v-if="f.permissions.Read">
                  <use xlink:href="/feather.svg#eye"></use>
                </svg>
                <svg class="feather" v-else>
                  <use xlink:href="/feather.svg#eye-off"></use>
                </svg>
              </span>
              <span class="icon is-small">
                <svg class="feather" v-if="f.permissions.Write">
                  <use xlink:href="/feather.svg#edit"></use>
                </svg>
                <svg class="feather" v-else>
                  <use xlink:href="/feather.svg#lock"></use>
                </svg>
              </span>
            </td>

            <td class="nowrap" title="{{f.size}} Bytes">{{ f.size }}</td>
            <td class="overflow">
              <a class="button is-small" @click="deleteFile(f)">
                <span class="icon is-small">
                  <svg class="feather">
                    <use xlink:href="/feather.svg#trash"></use>
                  </svg>
                </span>
              </a>
            </td>
            <td class="overflow">{{ f.created }}</td>
            <td class="overflow">{{ f.modified }}</td>
            <td class="overflow" title="{{f.accessed}}">{{ f.accessed }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script>
import UploadModal from '@/components/modal.vue';
export default {
    name: 'Folder',
    components: {
        UploadModal,
    },
    props: {
        // kann als Root-Element keine Props mehr haben
        //folder: Object,
    },
    data() {
        return {
            folder: JSON.parse(document.getElementById('data').innerHTML),

            isModalVisible: false,
        };
    },
    methods: {
        deleteFile(file) {
            fetch(
                new Request(file.path, {
                    method: 'DELETE',
                })
            )
                //.then(response => response.json())
                //.then(body => console.log(body))
                .catch(function(foo) {
                    console.log('FAILURE!!', foo);
                });
        },

        showModal() {
            this.isModalVisible = true;
        },
    },
    mounted() {
        console.log('Folder.vue =>', this.folder);
    },
};
</script>

<style lang="scss">
#app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
}
#nav {
    padding: 30px;
    a {
        font-weight: bold;
        color: #2c3e50;
        &.router-link-exact-active {
            color: #42b983;
        }
    }
}
</style>
