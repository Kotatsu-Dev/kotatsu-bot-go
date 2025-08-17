import {
  Box,
  Container,
  Drawer,
  Flex,
  Icon,
  IconButton,
  Portal,
  Stack,
  Tabs,
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
import {
  FaBars,
  FaCalendarPlus,
  FaClock,
  FaFileExport,
  FaTrash,
  FaUsers,
} from "react-icons/fa";
import { FaCalendarDays, FaMessage, FaShuffle } from "react-icons/fa6";

const tabs = [
  {
    icon: <FaClock />,
    value: "data",
    title: "Data",
  },
  {
    icon: <FaCalendarPlus />,
    value: "events",
    title: "Events",
  },
  {
    icon: <FaCalendarDays />,
    value: "calendar",
    title: "Calendar",
  },
  {
    icon: <FaUsers />,
    value: "users",
    title: "Users",
  },
  {
    icon: <FaShuffle />,
    value: "roulettes",
    title: "Roulettes",
  },
  {
    icon: <FaMessage />,
    value: "messages",
    title: "Messages",
  },
  // {
  //   value: "requests",
  //   title: "Requests",
  // },
  {
    icon: <FaFileExport />,
    value: "export",
    title: "Export",
  },
  {
    icon: <FaTrash />,
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
                              <Flex gap={2} alignItems={"center"}>
                                <Icon size="md">{tab.icon}</Icon>
                                {tab.title}
                              </Flex>
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
                <Icon>{tab.icon}</Icon>
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
