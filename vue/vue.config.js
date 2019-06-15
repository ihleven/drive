// vue.config.js
module.exports = {
    publicPath: '/',
    outputDir: '../_static',
    assetsDir: 'assets',
    pages: {
        //index: './src/pages/Home/main.js',
        
        
        // error: {
        //     entry: './src/pages/error.js',
        //     template: './public/templates/error.html',
        //     filename: 'templates/error.html',
        //     minify: false,
        //     chunks: ['chunk-vendors', 'chunk-common', 'error'],
        // },
        drive: {
            entry: './src/pages/drive/drive.js',
            template: './public/templates/drive.html',
            //filename: 'templates/drive.html',
            filename: 'index.html',
            minify: false,
//            chunks: ['chunk-vendors', 'chunk-common', 'drive'],
        },

        
        // arbeit: {
        //     entry: './src/pages/arbeit/arbeit.js',
        //     template: './public/templates/arbeit/arbeit.html',
        //     filename: 'templates/arbeit.html',
        //     minify: false,
        //     //chunks: ['chunk-vendors', 'chunk-common', 'album'],
        // },
    },
    runtimeCompiler: true,
    devServer: {
        proxy: 'http://localhost:3000'
    }
};