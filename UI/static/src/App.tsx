import {
  Box,
  Container,
  Drawer,
  IconButton,
  Portal,
  Stack,
  Tabs,
  Text,
} from "@chakra-ui/react";
import "./App.css";
import { Provider } from "./components/ui/provider";
import { APIProvider } from "./api/api";
import { DataTab } from "./components/DataTab";
import { Toaster } from "./components/ui/toaster";
import { EventsTab } from "./components/EventsTab";
import { CalendarTab } from "./components/CalendarTab";
import { UsersTab } from "./components/UsersTab";
import { MessagesTab } from "./components/MessagesTab";
import { RequestTab } from "./components/RequestTab";
import { ExportTab } from "./components/ExportTab";
import { DeletionTab } from "./components/DeletionTab";
import { RouletteTab } from "./components/RouletteTab";
import { useState } from "react";
import { FaBars } from "react-icons/fa";

const tabs = [
  {
    value: "data",
    title: "Data",
  },
  {
    value: "events",
    title: "Events",
  },
  {
    value: "calendar",
    title: "Calendar",
  },
  {
    value: "users",
    title: "Users",
  },
  {
    value: "roulettes",
    title: "Roulettes",
  },
  {
    value: "messages",
    title: "Messages",
  },
  // {
  //   value: "requests",
  //   title: "Requests",
  // },
  {
    value: "export",
    title: "Export",
  },
  {
    value: "deletion",
    title: "Deletion",
  },
];

function App() {
  const [tab, setTab] = useState("data");

  return (
    <Provider>
      <APIProvider>
        <Box
          borderBottom={"colorPalette.500"}
          borderBottomWidth={1}
          display={{ base: "flex", md: "none" }}
        >
          <Container maxW={"lg"} pt={2} pb={2} colorPalette={"orange"}>
            <Drawer.Root placement={"start"}>
              <Drawer.Trigger asChild>
                <IconButton>
                  <FaBars />
                </IconButton>
              </Drawer.Trigger>
              <Portal>
                <Drawer.Backdrop />
                <Drawer.Positioner>
                  <Drawer.Content>
                    <Drawer.Header>
                      <Drawer.Title>Menu</Drawer.Title>
                    </Drawer.Header>
                    <Drawer.Body>
                      <Stack>
                        {tabs.map((tab) => (
                          <Drawer.ActionTrigger asChild key={tab.value}>
                            <Box
                              textStyle="lg"
                              cursor={"pointer"}
                              onClick={() => setTab(tab.value)}
                            >
                              {tab.title}
                            </Box>
                          </Drawer.ActionTrigger>
                        ))}
                      </Stack>
                    </Drawer.Body>
                  </Drawer.Content>
                </Drawer.Positioner>
              </Portal>
            </Drawer.Root>
          </Container>
        </Box>
        <Tabs.Root
          value={tab}
          onValueChange={(e) => setTab(e.value)}
          colorPalette={"orange"}
        >
          <Tabs.List
            maxW={"100%"}
            overflowX={"scroll"}
            scrollbarWidth={"none"}
            display={{ base: "none", md: "flex" }}
          >
            {tabs.map((tab) => (
              <Tabs.Trigger key={tab.value} value={tab.value} flexShrink={0}>
                {tab.title}
              </Tabs.Trigger>
            ))}
          </Tabs.List>
          <Tabs.Content value="data">
            <DataTab />
          </Tabs.Content>
          <Tabs.Content value="events">
            <EventsTab />
          </Tabs.Content>
          <Tabs.Content value="calendar">
            <CalendarTab />
          </Tabs.Content>
          <Tabs.Content value="users">
            <UsersTab />
          </Tabs.Content>
          <Tabs.Content value="roulettes">
            <RouletteTab />
          </Tabs.Content>
          <Tabs.Content value="messages">
            <MessagesTab />
          </Tabs.Content>
          <Tabs.Content value="requests">
            <RequestTab />
          </Tabs.Content>
          <Tabs.Content value="export">
            <ExportTab />
          </Tabs.Content>
          <Tabs.Content value="deletion">
            <DeletionTab />
          </Tabs.Content>
        </Tabs.Root>
        <Toaster />
      </APIProvider>
    </Provider>
  );
}

export default App;
