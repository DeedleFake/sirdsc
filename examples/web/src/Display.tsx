import { useMemo } from "react";

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
  const query = useMemo(
    () =>
      Object.entries(params).reduce((query, [name, val]) => {
        query.set(name, `${val}`);
        return query;
      }, new URLSearchParams()),
    [params],
  );
  const src = useMemo(() => `/generate?${query}`, [query]);

  return src ? (
    <img className="flex-1 m-4" alt="Display" src={src} />
  ) : null;
}

export default Display;
