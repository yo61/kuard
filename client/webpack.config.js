var path = require('path');

var BUILD_DIR = path.resolve(__dirname, '../sitedata/built');

module.exports = {
  entry: './src/index.jsx',
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: ['babel-loader']
      }
    ]
  },
  resolve: {
    extensions: ['.js', '.jsx']
  },
  output: {
    path: BUILD_DIR,
    publicPath: '/built/',
    filename: 'bundle.js'
  },
  performance: { hints: false },
  devServer: {
    hot: true,
    port: 8081,
    // Catch-all proxy to the kuard Go server. The dev server's own middleware serves /built/
    // before the proxy runs, so the bundle stays local and everything else reaches the backend.
    proxy: [
      {
        context: function () { return true; },
        target: 'http://localhost:8080'
      }
    ]
  }
};
