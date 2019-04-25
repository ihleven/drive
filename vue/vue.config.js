// vue.config.js
module.exports = {
    publicPath: '/',
    outputDir: '../_static',
    assetsDir: 'assets',
    pages: {
        index: './src/pages/Home/main.js',
        error: {
            entry: './src/pages/error.js',
            template: './public/templates/error.html',
            filename: 'templates/error.html',
            minify: false,
            chunks: ['chunk-vendors', 'chunk-common', 'error'],
        },
        file: {
            entry: './src/pages/file/file.js',
            template: './public/templates/file.html',
            filename: 'templates/file.html',
            minify: false,
            chunks: ['chunk-vendors', 'chunk-common', 'file'],
        },
        folder: {
            entry: './src/pages/folder/folder.js',
            template: './public/templates/folder.html',
            filename: 'templates/folder.html',
            minify: false,
            chunks: ['chunk-vendors', 'chunk-common', 'folder'],
        },
        image: {
            entry: './src/pages/image/image.js',
            template: './public/templates/image.html',
            filename: 'templates/image.html',
            minify: false,
            chunks: ['chunk-vendors', 'chunk-common', 'image'],
        },
        album: {
            entry: './src/pages/album/album.js',
            template: './public/templates/album.html',
            filename: 'templates/album.html',
            minify: false,
            chunks: ['chunk-vendors', 'chunk-common', 'album'],
        },
    },
    runtimeCompiler: true,
};