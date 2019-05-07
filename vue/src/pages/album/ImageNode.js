import { Node, Plugin } from 'tiptap';
import AlbumImage from './AlbumImage.vue';

export default class Image extends Node {
    get name() {
        return 'image';
    }

    get schema() {
        return {
            inline: true,
            attrs: {
                src: {
                    default: null,
                },
                alt: {
                    default: null,
                },
                title: {
                    default: null,
                },
            },
            group: 'inline',
            draggable: true,
            parseDOM: [
                {
                    tag: 'img[src]',
                    getAttrs: dom => ({
                        src: dom.getAttribute('src'),
                        title: dom.getAttribute('title'),
                        alt: dom.getAttribute('alt'),
                    }),
                },
            ],
            toDOM: node => ['img', node.attrs],
        };
    }

    // return a vue component
    // this can be an object or an imported component
    get view() {
        return {
            // there are some props available
            // `node` is a Prosemirror Node Object
            // `updateAttrs` is a function to update attributes defined in `schema`
            // `editable` is the global editor prop whether the content can be edited
            // `options` is an array of your extension options
            // `selected`
            name: 'TipTapImage',
            props: ['node', 'updateAttrs', 'editable'],
            inject: ['prefix', 'album'],
            computed: {
                src: {
                    get() {
                        //return this.node.attrs.src;
                        return this.prefix + '/' + this.album + '/' + this.node.attrs.src;
                    },
                    set(src) {
                        // we cannot update `src` itself because `this.node.attrs` is immutable
                        this.updateAttrs({
                            src,
                        });
                    },
                },
                url() {},
            },
            created() {
                console.log(this.prefix, this.src); // => "bar"
            },
            mounted() {
                console.log('TipTapImage', this.$parent);
            },
            methods: {
                open_pswp() {
                    console.log('node image', this.getComponent('TipTap'));
                    this.$pswp_open(0, [
                        {
                            src: this.src,
                            w: 4000,
                            h: 3000,
                        },
                    ]);
                },
                getComponent(componentName) {
                    let component = null;
                    let parent = this.$parent;
                    while (parent && !component) {
                        if (parent.$options.name === componentName) {
                            component = parent;
                        }
                        parent = parent.$parent;
                    }
                    return component;
                },
            },
            template: `
            <img class="thumb" :src="src" @click="open_pswp" />
      `,
        };
    }

    commands({ type }) {
        return attrs => (state, dispatch) => {
            const { selection } = state;
            const position = selection.$cursor ? selection.$cursor.pos : selection.$to.pos;
            const node = type.create(attrs);
            const transaction = state.tr.insert(position, node);
            dispatch(transaction);
        };
    }

    get plugins() {
        return [
            new Plugin({
                props: {
                    handleDOMEvents: {
                        drop(view, event) {
                            const hasFiles =
                                event.dataTransfer && event.dataTransfer.files && event.dataTransfer.files.length;

                            if (!hasFiles) {
                                return;
                            }

                            const images = Array.from(event.dataTransfer.files).filter(file =>
                                /image/i.test(file.type)
                            );

                            if (images.length === 0) {
                                return;
                            }

                            event.preventDefault();

                            const { schema } = view.state;
                            const coordinates = view.posAtCoords({
                                left: event.clientX,
                                top: event.clientY,
                            });

                            images.forEach(image => {
                                const reader = new FileReader();

                                reader.onload = readerEvent => {
                                    const node = schema.nodes.image.create({
                                        src: readerEvent.target.result,
                                    });
                                    const transaction = view.state.tr.insert(coordinates.pos, node);
                                    view.dispatch(transaction);
                                };
                                reader.readAsDataURL(image);
                            });
                        },
                    },
                },
            }),
        ];
    }
}
