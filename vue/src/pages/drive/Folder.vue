<script>
import UploadModal from '@/components/modal.vue';
import Parallax from '@/components/Parallax.vue';
import NavigationOverlay from '@/components/NavigationOverlay.vue';
import NavbarTop from '@/components/NavbarTop.vue';
import FeatherIcon from '@/components/FeatherIcon.vue';
import { mapState } from 'vuex';

export default {
    name: 'Folder',
    components: {
        UploadModal,
        Parallax,
        NavigationOverlay,
        NavbarTop,
        FeatherIcon,
    },
    data() {
        return {
            blah: false,
            isModalVisible: false,
            menuOpen: false,
        };
    },
    computed: {
        ...mapState(['folder', 'account', 'breadcrumbs']),
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
        overlayNavigationToggler() {
            this.menuOpen = !this.menuOpen;
        },
    },
    mounted() {
        console.log('Folder.vue =>', this.folder);
    },
};
</script>

<template>
  <div class="application-wrapper">
    <upload-modal :visible.sync="isModalVisible" :url="folder ? folder.path : ''"/>

    <navigation-overlay :open.sync="menuOpen">
      <h3 class="title is-spaced">
        <a href="/alben">Fotoalben</a>
      </h3>
      <h4 class="subtitle">
        <a href="/alben/Mallorca">Mallorca</a>
      </h4>
      <a class="subtitle" href="/alben/hochzeitsreise">Hochzeitsreise</a>
      <h2 class="subtitle">Fullheight subtitle</h2>
    </navigation-overlay>

    <navbar-top :account="account" @click:hamburger="overlayNavigationToggler()"/>

    <Parallax>
      <template #header>
        <section class="hero is-medium is-dark is-bold">
          <div class="hero-body">
            <div class="container">
              <h1 class="title">{{ folder.name }}</h1>
              <h3 class="subtitle">
                <nav class="breadcrumb" aria-label="breadcrumbs">
                  <ul>
                    <li v-for="item in breadcrumbs" :key="item.Path">
                      <router-link :to="item.Path" class="navbar-item" v-text="item.Name"></router-link>

                    </li>
                  </ul>
                </nav>
              </h3>
            </div>
          </div>
        </section>
      </template>

      <template #attached>
        <nav class="navbar attached">
          <a class="navbar-item">
            <feather-icon name="upload-cloud"/>&nbsp;
            <span>upload</span>
          </a>

          <a class="navbar-item">
            <feather-icon name="activity"/>
          </a>
          <router-link :to="folder.path ? folder.path.replace('home', 'alben') : ''" class="navbar-item">
            <feather-icon name="camera"/>
          </router-link>
          <div class="navbar-item has-dropdown" :class="{'is-active':blah}" @click="blah=!blah">
            <a class="navbar-link is-arrowless">
              <feather-icon name="more-horizontal"/>
            </a>

            <div class="navbar-dropdown">
              <a class="navbar-item">
                <div class="dropdown-content">
                  <div class="dropdown-item">
                    <p>
                      You can insert
                      <strong>any type of content</strong> within the dropdown menu.
                    </p>
                  </div>
                </div>
              </a>
              <a class="navbar-item">Elements</a>
              <a class="navbar-item">Components</a>
              <hr class="navbar-divider">
              <div class="navbar-item">Version 0.7.4</div>
            </div>
          </div>
        </nav>
      </template>

      <section class="section has-background-white"></section>

      <div class="scrollwrapper">
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
                <feather-icon :name="f.mime.Type == 'dir' ? 'folder' : 'file'" :size="'small'"/>
              </td>
              <td>
                <!--<a :href="folder.name + '/' + f.name" :title="f.path">{{ f.name }}</a>-->
                                      <router-link :to="f.path" class="navbar-item" :title="f.path" v-text="f.name"></router-link>

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

              <td class="nowrap" :title="f.size + 'Bytes'">{{ f.size }}</td>
              <td class="overflow">
                <span class="icon is-small" @click="deleteFile(f)">
                  <svg class="feather">
                    <use xlink:href="/feather.svg#trash"></use>
                  </svg>
                </span>
              </td>
              <td class="overflow">{{ f.created }}</td>
              <td class="overflow">{{ f.modified }}</td>
              <td class="overflow" :title="f.accessed">{{ f.accessed }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <section>
        <pre>
                    {{ folder }}
                </pre>
      </section>
    </Parallax>
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
</style>
<style lang="scss" scoped>
.scrollwrapper {
    width: 100vw;
    overflow: auto;
}
.navbar.attached {
    min-height: 2.5rem;
    background: transparent;
}
.navbar.attached .navbar-item,
.navbar.attached .navbar-link {
    background: transparent;
    color: rgb(187, 176, 176);
}
.navbar.attached .navbar-item:hover,
.navbar.attached .navbar-item.is-active {
    background-color: #ffffff32;
    color: rgb(248, 248, 248);
}
.navbar.attached .navbar-link:hover,
.navbar.attached .navbar-link.is-active {
    background-color: transparent;
    color: rgb(248, 248, 248);
}
</style>
