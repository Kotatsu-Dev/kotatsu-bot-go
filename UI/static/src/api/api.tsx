import { createContext, useContext, type ReactNode } from "react";

const APIContext = createContext(null);

export const APIProvider = (props: { children: ReactNode[] }) => {
  return (
    <APIContext.Provider value={null}>{props.children}</APIContext.Provider>
  );
};

export const useAPI = () => {
  return useContext(APIContext);
};
