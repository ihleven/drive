// vue.config.js
module.exports = {
    publicPath: '/',
    outputDir: '../_static',
    assetsDir: 'assets',
    pages: {
        drive: {
            entry: './src/pages/drive/drive.js',
            template: '../templates/drive.html',
            //filename: 'templates/drive.html',
            filename: 'index.html',
            minify: false,
        },
        arbeit: {
            entry: './src/pages/arbeit/arbeit.js',
            template: './public/templates/arbeit/arbeit.html',
            filename: 'templates/arbeit.html',
            minify: false,
        },
        //index: './src/pages/Home/main.js',
        
        
        // error: {
        //     entry: './src/pages/error.js',
        //     template: './public/templates/error.html',
        //     filename: 'templates/error.html',
        //     minify: false,
        //     chunks: ['chunk-vendors', 'chunk-common', 'error'],
        // },
        

       
    },
    runtimeCompiler: true,
    devServer: {
        proxy: 'http://localhost:3000'
    },
    css: {
        sourceMap: true
    },
    configureWebpack: (config) => {
        config.devtool = 'source-map'
    },
};