// vue.config.js
module.exports = {
    publicPath: "/",
    outputDir: "../_static",
    assetsDir: "assets",
    pages: {
        index: "./src/pages/Home/main.js",
        folder: {
            entry: "./src/pages/folder/folder.js",
            template: './public/templates/directory.html',
            filename: 'templates/directory.html',
            minify: false,
            chunks: ['chunk-vendors', 'chunk-common', 'folder']
        },
        file: {
            entry: "./src/pages/file/file.js",
            template: './public/templates/file.html',
            filename: 'templates/file.html',
            minify: false,
            chunks: ['chunk-vendors', 'chunk-common', 'file']
        },
        textfile: {
            entry: "./src/pages/Textfile/textfile.js",
            template: './public/templates/textfile.html',
            filename: 'templates/textfile.html',
            minify: false,
            chunks: ['chunk-vendors', 'chunk-common', 'textfile']
        }
    },
    runtimeCompiler: true
}