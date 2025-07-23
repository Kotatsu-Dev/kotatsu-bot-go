import { createContext, useContext, type ReactNode } from "react";
import { createUsersApi } from "./users";
import axios from "axios";

const createApi = (_ctx: null) => {
  const $ = axios.create({
    baseURL: `https://test.bot.kotatsu.spb.ru/api/`,
  });
  
  return {
    users: createUsersApi($),
  };
};

const APIContext = createContext(null);

export const APIProvider = (props: { children: ReactNode[] | ReactNode }) => {
  return (
    <APIContext.Provider value={null}>{props.children}</APIContext.Provider>
  );
};

export const useAPI = () => {
  return createApi(useContext(APIContext));
};
