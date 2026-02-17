export type DisplayProps = {
  params: Params;
};

export type Params = {
  src: string;
  pat: string;
  seed: number;
  partsize: number;
  depth: number;
  sym: boolean;
  inverse: boolean;
  flat: boolean;
};

export function Display({ params }: DisplayProps) {
  return <div>Not implemented.</div>;
}

export default Display;
