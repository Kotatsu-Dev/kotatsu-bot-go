import { Tabs } from "@chakra-ui/react";
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

function App() {
  return (
    <Provider>
      <APIProvider>
        <Tabs.Root defaultValue={"data"} colorPalette={"orange"}>
          <Tabs.List maxW={"100%"} overflowX={"scroll"} scrollbarWidth={"none"}>
            <Tabs.Trigger value="data" flexShrink={0}>
              Data
            </Tabs.Trigger>
            <Tabs.Trigger value="events" flexShrink={0}>
              Events
            </Tabs.Trigger>
            <Tabs.Trigger value="calendar" flexShrink={0}>
              Calendar
            </Tabs.Trigger>
            <Tabs.Trigger value="users" flexShrink={0}>
              Users
            </Tabs.Trigger>
            <Tabs.Trigger value="roulettes" flexShrink={0}>
              Roulettes
            </Tabs.Trigger>
            <Tabs.Trigger value="messages" flexShrink={0}>
              Messages
            </Tabs.Trigger>
            <Tabs.Trigger value="requests" flexShrink={0}>
              Requests
            </Tabs.Trigger>
            <Tabs.Trigger value="export" flexShrink={0}>
              Export
            </Tabs.Trigger>
            <Tabs.Trigger value="deletion" flexShrink={0}>
              Deletion
            </Tabs.Trigger>
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
