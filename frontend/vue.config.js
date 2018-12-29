// vue.config.js
module.exports = {
  runtimeCompiler: true,
  baseUrl: '/static/',
  outputDir: '../static/',
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
      template: 'src/templates/file.html',
      filename: 'file.html',
      title: 'File Page',
      chunks: ['chunk-vendors', 'chunk-common', 'file']
    },
    directory: {
      entry: 'src/entries/directory.js',
      template: 'src/templates/directory.html',
      filename: 'directory.html',
      title: 'directory Page',
      chunks: ['chunk-vendors', 'chunk-common', 'directory']
    }
  }
}
