import { createContext, useContext, useEffect, type ReactNode } from "react";
import axios, { AxiosError } from "axios";
import { createUsersApi } from "./users";
import { createActivitiesApi } from "./activities";
import { createCalendarApi } from "./calendar";
import { createBroadcastApi } from "./broadcast";
import { createRequestsApi } from "./requests";
import { createDbApi } from "./db";
import { createRoulettesApi } from "./roulettes";
import { toaster } from "../components/ui/toaster";
import z, { ZodError } from "zod";
import { useLocalStorage } from "usehooks-ts";
import { Button, Center } from "@chakra-ui/react";

const ErrorData = z.object({
  status: z.object({
    code: z.int(),
    message: z.string(),
  }),
});

const BASE_URL =
  import.meta.env.VITE_BASE_URL ??
  new URL("/", location.toString()).toString().slice(0, -1);
const API_URL =
  import.meta.env.VITE_API_URL ??
  new URL("/api", location.toString()).toString();

const createApi = (token: string) => {
  const $ = axios.create({
    headers: {
      Authorization: `Bearer ${token}`,
    },
    baseURL: API_URL,
  });

  return {
    users: createUsersApi($),
    activities: createActivitiesApi($),
    calendar: createCalendarApi($, BASE_URL),
    broadcast: createBroadcastApi($, BASE_URL),
    requests: createRequestsApi($),
    roulettes: createRoulettesApi($),
    db: createDbApi($, BASE_URL),
  };
};

const parseJwt = (token: string) => {
  var base64Url = token.split(".")[1];
  var base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  var jsonPayload = decodeURIComponent(
    window
      .atob(base64)
      .split("")
      .map(function (c) {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join(""),
  );

  return JSON.parse(jsonPayload);
};

const APIContext = createContext<ReturnType<typeof createApi> | null>(null);

export const APIProvider = (props: { children: ReactNode[] | ReactNode }) => {
  const [token, setToken, _] = useLocalStorage("token", "");

  useEffect(() => {
    (window as any).Telegram.Login.init({
      client_id: import.meta.env.VITE_BOT_ID,
      request_access: ["write"],
    });
  }, []);

  const login = () => {
    (window as any).Telegram.Login.open((result: any) => {
      if ("error" in result) {
        return;
      }
      const base = import.meta.env.PROD
        ? new URL("/", location.toString()).toString().slice(0, -1)
        : `http://localhost:8006`;

      // TODO: Model check
      axios
        .get(`${base}/api/login`, {
          headers: {
            Authorization: `Bearer ${result.id_token}`,
          },
        })
        .then((data) => setToken(data.data.token));
    });
  };

  if (token != "") {
    const parsed = parseJwt(token);
    if (parsed.exp * 1000 >= Date.now()) {
      // Not expired
      return (
        <APIContext.Provider value={createApi(token)}>
          {props.children}
        </APIContext.Provider>
      );
    }
  }

  return (
    <Center h={"100vh"}>
      <Button colorPalette={"cyan"} onClick={login}>
        Login with Telegram
      </Button>
    </Center>
  );
};

export const useAPI = () => {
  return useContext(APIContext)!;
};

export const handleError = (e: unknown) => {
  if (e instanceof ZodError) {
    console.log(`Error parsing data:\n${z.prettifyError(e)}`);
    console.log(e.issues);
    toaster.error({
      description: `Error parsing data:\n${z.prettifyError(e)}`,
    });
    return;
  }

  if (e instanceof AxiosError) {
    if (e.response && e.response.data) {
      const data = ErrorData.safeParse(e.response.data);
      if (data.success) {
        toaster.error({
          description: `Error: ${data.data.status.message}`,
        });
        return;
      }
    }
    toaster.error({
      description: `HTTP Error ${e.code ?? ""}`,
    });
    return;
  }

  toaster.error({
    description: `Unknown error`,
  });
};
