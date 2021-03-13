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
            sortedBy: ['type', 'name'],
            clonedFile: null,
        };
    },
    computed: {
        ...mapState(['folder', 'account', 'breadcrumbs', "error"]),
        entries() {
            
            const sortBy = (keysarray) => {
                let ret = 0;
                return (a, b) => {
                    for (let key of keysarray) {
                        if (key.charAt(0) == '-') {
                            key = key.substring(1)
                            ret = (a[key] > b[key]) ? -1 : ((b[key] > a[key]) ? 1 : 0);
                        } else {
                            ret = (a[key] > b[key]) ? 1 : ((b[key] > a[key]) ? -1 : 0);
                        }
                        if (ret != 0) return ret;
                    }
                    return 0;
                }
            };
            const addType = entry => Object.assign(entry, { type: entry.mime.Type == 'dir' ? 'D' : 'F' })
            if (!this.folder.entries) return [];
            return this.folder.entries.concat().map(addType).sort(sortBy(this.sortedBy));;
        }
    },
    mounted() {
        console.log('Folder.vue =>', this.folder);
    },
    methods: {
        renameFile(file) {
            this.clonedFile = Object.assign({}, file);
            console.log("renameFile", this.clonedFile);
            this.clonedFile.rename = newName => {
                this.$store.dispatch("submitFileForm", {url: file.path, name: newName});

            }
        },
        
        deleteFile(file) {
            fetch(
                new Request(file.path, {
                    method: 'DELETE',
                })
            )
                .then(response => console.log(response))
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
        errorClear() {
            this.$store.commit("error", null);
        },
        newFolder(event) {
            console.log("newFolder:", event.target.value)
            this.$store.commit("newFile", {name: event.target.value, mime: {Type: "dir"}, permissions: {},
            group: {}, owner: {}, path: event.target.value})
        }
    },
    filters: {
        timestamp(v) {
            return v;
        }
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
                    <nav class="breadcrumb" aria-label="breadcrumbs" style="margin-bottom: .25rem">
                        <ul>
                            <li>
                                <router-link to="/" class="navbar-item">
                                    <feather-icon sprite="entypo" name="database"/>
                                </router-link>
                            </li>
                            <li v-for="item in breadcrumbs" :key="item.Path">
                            <router-link :to="item.Path" class="navbar-item" v-text="item.Name"></router-link>

                            </li>
                        </ul>
                    </nav>
                    <div class="columns">
                        <div class="column is-three-fifths">

                            <h1 class="title">{{ folder.name }}</h1>
                            <h4 class="subtitle">{{folder.mime.Type}}  {{entries.length}} items</h4>
                        </div>
                        <div class="column">
<span class="mime" v-text="folder.mime.Type"></span>
              <span class="mime-sub">{{ folder.mime.Subtype }}</span>
              <span class="mode">{{folder.permissions.Notation}}</span>
              <span class="has-text-right" 
              :class="folder.permissions.IsOwner ? 'is-user-group' : 'is-not-user-group'">
                  {{ folder.owner.name }}
              </span>
              <span :class="folder.permissions.IsOwner ? 'is-user-group' : 'is-not-user-group'">
                {{ folder.group.name }}
            </span>
              <span class="icons">
                  <span class="navbar-item">
                  <svg><use :xlink:href="'/entypo-icons.svg#' + (folder.permissions.Read ? 'icon-eye' : 'eye-with-line')"></use></svg>
                  <svg><use :xlink:href="'/feather.svg#' + (folder.permissions.Write ? 'edit-3' : 'minus')"></use></svg>
                  </span>
              </span>
, {{folder.permissions.Notation}}
<span class="ts">{{ folder.created | timeformat }}</span>
              <span class="ts">{{ folder.modified | timeformat }}</span>
              <span class="ts" :title="folder.accessed">{{ folder.accessed | timeformat }}</span>
                        </div>
                    </div>
              
                
                    <div class="notification is-danger" v-if="error">
                        <button class="delete" @click="errorClear"></button>
                        <strong>{{error.status}}: {{error.statusText}}</strong><br>
                        <small style="white-space: pre;">{{error.data}}</small>
                    </div>
                </div>
            </section>
        </template>

        <template #attached>
            <nav class="navbar attached">
                <a class="navbar-item" @click="showModal">
                    <feather-icon name="upload-cloud" />&nbsp;
                    <span>upload</span>
                </a>

                <a class="navbar-item">
                    <feather-icon name="activity" />
                </a>

                <router-link :to="folder.path ? folder.path.replace('home', 'alben') : ''" class="navbar-item">
                    <feather-icon name="camera"/>
                </router-link>

                <div class="navbar-item has-dropdown" :class="{'is-active':blah}" @click="blah=!blah">
                    <a class="navbar-link is-arrowless">
                        <feather-icon name="more-horizontal"/>
                    </a>
                    <div class="navbar-dropdown is-boxed">
                        <a class="navbar-item" @click="sortedBy=['type', 'name']">type, name</a>
                        <a class="navbar-item" @click="sortedBy=['modified']">mod</a>
                        <a class="navbar-item" @click="sortedBy=['type', 'size']">size</a>
                   </div>
                    
                </div>
                <div class="navbar-item">
                    <input @change="newFolder" />
                </div>
            </nav>
        </template>

      <section class="section has-background-white"></section>

      <div class="entries">
        <table class="table is-striped is-fullwidth ">
          <thead>
            <tr>
              <th></th>
              <th>Name</th>

              <th></th>
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
            <tr v-for="f in entries" :key="f.name" draggable="true">
              <td class="">
                <feather-icon :sprite="f.type == 'D' ? 'entypo-icons' : 'feather'" :name="f.type == 'D' ? 'icon-folder' : 'file'"/>
              </td>
              <td class="name">
                <input v-if="clonedFile && f.name==clonedFile.name"
                        :value="clonedFile.name" @change="clonedFile.rename($event.target.value)" />
                <router-link v-else :to="f.path" class="navbar-item" :title="f.path" v-text="f.name"></router-link>
              </td>
              <td class="icons">
                <span class="navbar-item has-dropdown is-hoverable">
                    <a class="navbar-link is-arrowless">
                    <svg><use xlink:href="/feather.svg#more-vertical"></use></svg>
                    </a>
                    <div class="navbar-dropdown is-boxed">
                        <a class="navbar-item" @click="renameFile(f)">
                            <svg><use xlink:href="/feather.svg#edit-2"></use></svg>
                            <span class="menu-entry">rename</span>
                        </a>
                        <a class="navbar-item">
                            <svg><use xlink:href="/feather.svg#sliders"></use></svg>
                            <span class="menu-entry">permissions</span>
                        </a>
                        <hr class="navbar-divider">
                        <a class="navbar-item">
                            <svg><use xlink:href="/feather.svg#git-branch"></use></svg>
                            <span class="menu-entry">move</span>
                        </a>
                        <a class="navbar-item" @click="deleteFile(f)">
                            <svg><use xlink:href="/feather.svg#trash"></use></svg>
                            <span class="menu-entry">delete</span>
                        </a>
                    </div>
                </span>
              </td>

              <td class="mime" v-text="f.mime.Type"></td>
              <td class="mime-sub">{{ f.mime.Subtype }}</td>
              <td class="mode">{{f.permissions.Notation}}</td>
              <td class="has-text-right" 
              :class="f.permissions.IsOwner ? 'is-user-group' : 'is-not-user-group'">
                  {{ f.owner.name }}
              </td>
              <td :class="f.permissions.IsOwner ? 'is-user-group' : 'is-not-user-group'">
                {{ f.group.name }}
            </td>
              <td class="icons">
                  <span class="navbar-item">
                  <svg><use :xlink:href="'/entypo-icons.svg#' + (f.permissions.Read ? 'icon-eye' : 'eye-with-line')"></use></svg>
                  <svg><use :xlink:href="'/feather.svg#' + (f.permissions.Write ? 'edit-3' : 'minus')"></use></svg>
                  </span>
              </td>
              <td class="nowrap" :title="f.size + 'Bytes'">
                  {{ f.size | bytes }}
              </td>
              <td class="ts">{{ f.created | timeformat }}</td>
              <td class="ts">{{ f.modified | timeformat }}</td>
              <td class="ts" :title="f.accessed">{{ f.accessed | timeformat }}</td>
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


.navbar.attached {
    min-height: 2.5rem;
    background: transparent;
   .navbar-item {
        background-color: transparent;
        color: white;
        &:hover, &:active, &:focus, &.is-active {
        background-color: #fafafa;
        color: black;
        }
        .navbar-link { color: inherit; background-color: inherit;}
    }
    .navbar-dropdown {
        .navbar-item {
            color: black;
            background-color: white;
        
            &:hover, &.is-active {
                background-color: grey;
                color: dark;
            }
        }
    }
}


.is-user-group {
    text-decoration: underline;
    font-weight: 400;
}

.is-not-user-group {
    text-decoration: line-through;
}
.entries {
    width: 100vw;
    overflow-x: auto;
    overflow-y: visible;
    /*margin-top: 1rem;*/
margin-bottom: 1rem;
    .table {overflow: hidden;}
    tr {
        cursor: pointer;
        &:hover {
            background-color: #f0f0f0!important;
        }
        td {
            vertical-align: middle;
        }
        .name {
            padding: 0;
        }
        .mime {
            
            text-align: right;
            padding-right: .375rem;
            
        }
        .mime-sub {
            position: relative;
            padding-left: .25rem;
            &::before {
                position: absolute;
                left: -.275rem;
                content: "/"
            }
            &:empty {
                &::before {
                content: ""
            }
            }
        }
        .icons {
            padding: 0;
            vertical-align: middle;
            svg {
                width: 1.25rem;
                height: 1.25rem;
                fill: none;
                stroke: currentcolor;
                stroke-width: 1;
                stroke-linecap: round;
                stroke-linejoin: round;
            }
            .menu-entry { 
                padding-left: .5rem;
                font-size: 1rem; 
            }
        }
        .ts {
            vertical-align: middle;
            white-space: nowrap;
            font-size: .875rem;
            color: #bdbdbd;
            &:hover {
                color: #aaaaaa;
            }
        }
        .mode {
            white-space: nowrap;
            font-family: monospace;
            font-size: 1rem;
            vertical-align: middle;
        }
    }
}
</style>
