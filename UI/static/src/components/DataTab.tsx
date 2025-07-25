import { useAPI } from "../api/api";
import {
  Badge,
  Button,
  Card,
  Center,
  CloseButton,
  Container,
  DataList,
  Dialog,
  Heading,
  Link,
  Portal,
  Stack,
  Table,
} from "@chakra-ui/react";
import { toaster } from "./ui/toaster";
import { useState } from "react";
import type { User } from "@/api/users";
import type { Activity } from "@/api/activities";
import z, { ZodError } from "zod";

export const DataTab = () => {
  const api = useAPI();
  const [openUsers, setOpenUsers] = useState(false);
  const [users, setUsers] = useState<User[]>([]);
  const [usersType, setUsersType] = useState<"members" | "subscribers">(
    "members"
  );
  const [openEvents, setOpenEvents] = useState(false);
  const [events, setEvents] = useState<Activity[]>([]);

  const loadClubMembers = async () => {
    const users = await api.users.getAll();
    const clubMembers = users.filter((user) => user.is_club_member);
    if (clubMembers.length <= 0) {
      toaster.error({
        description: "No members in club",
      });
      return;
    }
    setUsers(clubMembers);
    setUsersType("members");
    setOpenUsers(true);
  };

  const loadNewsletterSubscribers = async () => {
    const users = await api.users.getAll();
    const subscribers = users.filter((user) => user.is_subscribe_newsletter);
    if (subscribers.length <= 0) {
      toaster.error({
        description: "No subscribers on newsletter",
      });
      return;
    }
    setUsers(subscribers);
    setUsersType("subscribers");
    setOpenUsers(true);
  };

  const loadEvents = async () => {
    try {
      const events = await api.activities.getAll();
      if (events.length <= 0) {
        toaster.error({
          description: "No events",
        });
        return;
      }
      setEvents(events);
      setOpenEvents(true);
    } catch (e) {
      if (e instanceof ZodError) {
        console.log(z.prettifyError(e));
      }
    }
  };

  const deleteEvent = async (event: Activity) => {
    await api.activities.setStatus({ id: event.id, status: false });
    await loadEvents();
  };

  return (
    <Container maxW={"lg"}>
      <Heading textAlign={"center"} pb={3}>
        Dashboard
      </Heading>
      <Stack>
        <Button onClick={loadClubMembers}>Show club members</Button>
        <Button onClick={loadNewsletterSubscribers}>
          Show newsletter subscribers
        </Button>
        <Button onClick={loadEvents}>Show events list</Button>
      </Stack>
      <Dialog.Root open={openUsers} onOpenChange={(e) => setOpenUsers(e.open)}>
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.CloseTrigger asChild>
                <CloseButton size="sm" />
              </Dialog.CloseTrigger>
              <Dialog.Header>
                <Dialog.Title>
                  {usersType === "members"
                    ? "Club members"
                    : "Newsletter subscribers"}
                </Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <Table.Root interactive>
                  <Table.Header>
                    <Table.Row>
                      <Table.ColumnHeader>Username</Table.ColumnHeader>
                      <Table.ColumnHeader>Telegram ID</Table.ColumnHeader>
                    </Table.Row>
                  </Table.Header>
                  <Table.Body>
                    {users.map((user) => (
                      <Table.Row key={user.id}>
                        <Table.Cell>
                          <Link href={`https://t.me/${user.user_name}`}>
                            @{user.user_name}
                          </Link>
                        </Table.Cell>
                        <Table.Cell>{user.user_tg_id}</Table.Cell>
                      </Table.Row>
                    ))}
                  </Table.Body>
                </Table.Root>
              </Dialog.Body>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
      <Dialog.Root
        open={openEvents}
        onOpenChange={(e) => setOpenEvents(e.open)}
      >
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.CloseTrigger asChild>
                <CloseButton size="sm" />
              </Dialog.CloseTrigger>
              <Dialog.Header>
                <Dialog.Title>Events list</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <Stack>
                  {events.map((event) => (
                    <Card.Root key={event.id}>
                      <Card.Header>
                        <Heading>{event.title}</Heading>
                      </Card.Header>
                      <Card.Body>
                        <DataList.Root orientation={"horizontal"}>
                          <DataList.Item>
                            <DataList.ItemLabel>Status</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {event.status ? (
                                <Badge colorPalette={"green"}>Active</Badge>
                              ) : (
                                <Badge colorPalette={"red"}>Inactive</Badge>
                              )}
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>Description</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {event.description}
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>Date</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {event.date_meeting}
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>Location</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {event.location}
                            </DataList.ItemValue>
                          </DataList.Item>
                        </DataList.Root>
                      </Card.Body>
                      <Card.Footer>
                        <Button disabled variant={"outline"}>
                          Download signed up
                        </Button>
                        <Button
                          disabled={!event.status}
                          colorPalette={"red"}
                          onClick={() => deleteEvent(event)}
                        >
                          Delete
                        </Button>
                      </Card.Footer>
                    </Card.Root>
                  ))}
                </Stack>
              </Dialog.Body>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </Container>
  );
};
