import type { ChangeEvent } from "react";
import { useState } from "react";

export type InputData = { [name: string]: string | boolean | number };

export type InputPropData<T> = {
  type: string;
  value: T;
  onChange: (ev: ChangeEvent<HTMLInputElement>) => void;
  [name: string]: any;
};

export interface Inputs {
  text(name: string, initialValue?: string): InputPropData<string>;
  checkbox(name: string, initialValue?: boolean): InputPropData<boolean>;
  number(name: string, initialValue?: number): InputPropData<number>;
  range(
    name: string,
    min: number,
    max: number,
    initialValue?: number,
  ): InputPropData<number>;
}

export function useInputs(): [InputData, Inputs] {
  const [data, setData] = useState<InputData>({});

  return [
    data,
    {
      text(name: string, initialValue = ""): InputPropData<string> {
        return {
          type: "text",
          value: (data[name] as string) ?? initialValue,
          onChange: (ev) => {
            setData((data) => ({
              ...data,
              [name]: ev.target.value,
            }));
          },
        };
      },

      checkbox(name: string, initialValue = false): InputPropData<boolean> {
        return {
          type: "checkbox",
          value: (data[name] as boolean) ?? initialValue,
          onChange: (ev) => {
            setData((data) => ({
              ...data,
              [name]: ev.target.checked,
            }));
          },
        };
      },

      number(name: string, initialValue = 0): InputPropData<number> {
        return {
          type: "number",
          value: (data[name] as number) ?? initialValue,
          onChange: (ev) => {
            let val = parseFloat(ev.target.value);
            if (!isNaN(val)) {
              setData((data) => ({
                ...data,
                [name]: val,
              }));
            }
          },
        };
      },

      range(
        name: string,
        min: number,
        max: number,
        initialValue = (min + max) / 2,
      ): InputPropData<number> {
        return {
          type: "range",
          value: (data[name] as number) ?? initialValue,
          min,
          max,
          onChange: (ev) => {
            let val = parseFloat(ev.target.value);
            if (!isNaN(val)) {
              setData((data) => ({
                ...data,
                [name]: val,
              }));
            }
          },
        };
      },
    },
  ];
}

export type InputProps = {
  label: string;
  type: string;
  value: any;
  [name: string]: any;
};

export function Input({ label, type, value, ...props }: InputProps) {
  const className =
    type === "checkbox" ? "checkbox" : type === "range" ? "range" : "input";

  return (
    <div className="flex flex-col justify-center items-center m-2 gap-0.5">
      <div>{label}</div>
      <input className={className} type={type} value={value} {...props} />
      {type === "range" ? <div className="text-sm">{value}</div> : null}
    </div>
  );
}

export default Input;
