/* eslint-env node */
const webpack = require('webpack');
const path = require('path');

const prod = process.env.NODE_ENV === 'production';

module.exports = {
  entry: './src/app.tsx',

  output: {
    path: path.resolve(__dirname, 'public'),
    filename: 'app.js'
  },

  devtool: 'inline-source-map',

  devServer: {
    publicPath: '/public/',
    port: 8081,
    historyApiFallback: true
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
  },

  plugins: [
    new webpack.DefinePlugin({
      API_URL: JSON.stringify(prod ? 'https://something.com' : 'http://localhost:8080')
    })
  ]
};