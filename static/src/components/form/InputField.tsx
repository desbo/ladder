import * as React from 'react';
import { Ref } from 'react';

export default class InputField extends React.Component {
  props: {
    label: string,
    type: string,
    inputRef?: Ref<HTMLInputElement>,
    [propName: string]: any
  }

  render() {
    const {
      label,
      type,
      inputRef,
      ...inputProps
    } = this.props;

    return (
      <div className="field">
        <label className="label">{label}</label>
        <div className="control">
          <input className="input" type={type} ref={inputRef} {...inputProps} />
        </div>
      </div>
    );
  }
};