import React from 'react';
import ReactDOM from 'react-dom';

export function init() {
  return ReactDOM.render(
    <h1>Hello, world!</h1>,
    document.getElementById('app')
  );
}
