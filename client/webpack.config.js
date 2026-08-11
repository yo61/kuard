var path = require('path');

var BUILD_DIR = path.resolve(__dirname, '../sitedata/built');

// Exported as a function so the Babel config can see the build mode. preset-react lives here
// rather than in .babelrc because of its `development` flag: Babel picks its envName from
// NODE_ENV, which `webpack --mode=production` does NOT set for the build process -- it only
// defines process.env.NODE_ENV inside the bundle. Babel would therefore default to development
// and emit jsxDEV(), which React's production jsx-dev-runtime does not export, and the app would
// fail to mount with an empty page. Deriving it from argv.mode cannot drift.
module.exports = function (env, argv) {
  var isProduction = argv.mode === 'production';

  return {
    entry: './src/index.jsx',
    module: {
      rules: [
        {
          test: /\.(js|jsx)$/,
          exclude: /node_modules/,
          use: [
            {
              loader: 'babel-loader',
              options: {
                presets: [['@babel/preset-react', { development: !isProduction }]]
              }
            }
          ]
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
};
