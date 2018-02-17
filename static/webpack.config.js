/* eslint-env node */
const webpack = require('webpack');
const path = require('path');
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');

const prod = process.env.NODE_ENV === 'production';

module.exports = {
  entry: {
    app: './src/app.tsx',
    css: './scss/main.scss'
  },

  output: {
    path: path.resolve(__dirname, 'public'),
    filename: '[name].js'
  },

  devtool: 'inline-source-map',

  devServer: {
    publicPath: '/public/',
    port: 8081,
    historyApiFallback: true
  },

  resolve: {
    modules: [path.resolve(__dirname, 'src'), path.resolve(__dirname, 'scss'), 'node_modules'],
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
      },

      {
        test: /\.scss$/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: [
            {
              loader: 'css-loader',
              options: { minimize: prod }
            }, 
            'sass-loader'
          ]
        })
      }
    ]
  },

  plugins: [
    new webpack.DefinePlugin({
      API_URL: JSON.stringify(prod ? 'https://api-dot-tt-ladder.appspot.com' : 'http://localhost:8080'),
      'process.env': {
        'NODE_ENV': JSON.stringify(prod ? 'production' : '')
      },

      FIREBASE_CONFIG: JSON.stringify(prod ? require('./firebase_config.js').prod : require('./firebase_config.js').dev)
    }),

    prod ? new UglifyJsPlugin({
      parallel: true
    }) : () => null,

    new ExtractTextPlugin('main.css')
  ]
};
