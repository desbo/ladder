/* eslint-env node */

const path = require('path');

module.exports = {
  entry: './src/app.js',

  output: {
    path: path.resolve(__dirname, 'public'),
    filename: 'app.js'
  },

  devtool: 'inline-source-map',

  devServer: {
    port: 8081
  },

  resolve: {
    modules: [path.resolve(__dirname, 'src'), 'node_modules'],
    extensions: ['.js', '.jsx']
  },

  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env', '@babel/preset-react']
          }
        }
      }
    ]
  },
};