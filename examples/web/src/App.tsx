import type { ChangeEvent } from "react";
import { useState } from "react";

type InputData = { [name: string]: string | boolean | number };

type Input<T> = {
  type: string;
  value: T;
  onChange: (ev: ChangeEvent<HTMLInputElement>) => void;
  [name: string]: any;
};

interface Inputs {
  text(name: string, initialValue?: string): Input<string>;
  checkbox(name: string, initialValue?: boolean): Input<boolean>;
  number(name: string, initialValue?: number): Input<number>;
  range(
    name: string,
    min: number,
    max: number,
    initialValue?: number,
  ): Input<number>;
}

function useInputs(): [InputData, Inputs] {
  const [data, setData] = useState<InputData>({});
  return [
    data,
    {
      text(name: string, initialValue = ""): Input<string> {
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

      checkbox(name: string, initialValue = false): Input<boolean> {
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

      number(name: string, initialValue = 0): Input<number> {
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
        initialValue = 0,
      ): Input<number> {
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

function Input({
  label,
  type,
  value,
  ...props
}: {
  label: string;
  type: string;
  value: any;
  [name: string]: any;
}) {
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

export function App() {
  const [params, inputs] = useInputs();

  return (
    <div className="flex flex-col items-center">
      <div className="flex flex-col bg-base-300 m-4 p-4 pt-25 transform-[translate(0px,-100px)] rounded-lg shadow">
        <div className="flex flex-row justify-end mx-0 my-2">
          <Input label="Depth Map" {...inputs.text("src")} />
          <Input label="Pattern" {...inputs.text("pat")} />
          <Input label="Seed" {...inputs.number("pat")} />
          <Input label="Part Size" {...inputs.range("partsize", 0, 500, 100)} />
        </div>
      </div>
    </div>
  );
}

export default App;
