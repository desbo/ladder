import * as React from 'react';
import { Ref } from 'react';

export default class InlineInputField extends React.Component {
  props: {
    type: string,
    inputRef?: Ref<HTMLInputElement>,
    [propName: string]: any
  }

  render() {
    const {
      type,
      inputRef,
      ...inputProps
    } = this.props;

    return (
      <div className="field">
        <div className="control">
          <input
            className="input is-small" 
            type={type} 
            ref={inputRef} 
            {...inputProps} />
        </div>
      </div>
    );
  }
};