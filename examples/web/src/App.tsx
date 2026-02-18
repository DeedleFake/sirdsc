import Input, { useInputs } from "./Input.tsx";
import Display from "./Display.tsx";
import type { Params } from "./Display.tsx";

export function App() {
  const [params, inputs] = useInputs();

  return (
    <div className="flex flex-col items-center">
      <div className="flex flex-col bg-base-300 m-4 p-4 pt-25 transform-[translate(0px,-100px)] rounded-lg shadow-lg">
        <div className="flex flex-row justify-end mx-0 my-2">
          <Input label="Depth Map" {...inputs.text("src")} />
          <Input label="Pattern" {...inputs.text("pat")} />
          <Input label="Seed" {...inputs.number("seed")} />
          <Input label="Part Size" {...inputs.range("partsize", 0, 500, 100)} />
          <Input label="Max Depth" {...inputs.range("depth", 0, 50, 40)} />
          <Input
            label="Symmetric Random Generation"
            {...inputs.checkbox("sym")}
          />
          <Input label="Inverse" {...inputs.checkbox("inverse")} />
          <Input label="Flat" {...inputs.checkbox("flat")} />
        </div>
      </div>

      <Display params={params as Params} />
    </div>
  );
}

export default App;
