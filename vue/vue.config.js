// vue.config.js
module.exports = {
  runtimeCompiler: true,
  baseUrl: '/dist/',
  //outputDir: './dist/',
  pages: {
    index: {
      entry: 'src/main.js',
      template: 'public/index.html',
      filename: 'index.html',
      title: 'Index Page',
      chunks: ['chunk-vendors', 'chunk-common', 'index']
    },
    file: {
      entry: 'src/entries/file.js',
      template: '../templates/file.html',
      filename: 'file.html',
      title: 'File Page',
      chunks: ['chunk-vendors', 'chunk-common', 'file']
    },
    directory: {
      entry: 'src/entries/directory.js',
      template: '../templates/directory.html',
    },
    album: {
      entry: 'src/entries/album.js',
      //template: '../../templates/album.html',
      template: '../templates/album.html',
      //filename: 'album.html',
    }
  }
}