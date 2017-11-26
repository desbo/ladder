/* eslint-env node */

const path = require('path');

module.exports = {
  entry: './src/app.tsx',

  output: {
    path: path.resolve(__dirname, 'public'),
    filename: 'app.js'
  },

  devtool: 'inline-source-map',

  devServer: {
    publicPath: '/public/',
    port: 8081
  },

  resolve: {
    modules: [path.resolve(__dirname, 'src'), 'node_modules'],
    extensions: ['.ts', '.tsx', '.js', '.jsx']
  },

  module: {
    rules: [
      { 
        test: /\.tsx?$/, 
        loader: 'awesome-typescript-loader' 
      },

      { 
        enforce: 'pre', 
        test: /\.js$/, 
        loader: 'source-map-loader' 
      }
    ]
  }
};