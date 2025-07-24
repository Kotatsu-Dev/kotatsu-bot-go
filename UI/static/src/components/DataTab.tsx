import { useAPI } from "../api/api";
import {
  Button,
  Center,
  CloseButton,
  Container,
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

export const DataTab = () => {
  const api = useAPI();
  const [openUsers, setOpenUsers] = useState(false);
  const [users, setUsers] = useState<User[]>([]);
  const [usersType, setUsersType] = useState<"members" | "subscribers">(
    "members"
  );
  const [openEvents, setOpenEvents] = useState(false);

  const loadClubMembers = async () => {
    const users = await api.users.getAll();
    const clubMembers = users.filter((user) => user.is_club_member);
    if (clubMembers.length <= 0) {
      toaster.error({
        description: "No members in club",
      });
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
        description: "No members in club",
      });
    }
    setUsers(subscribers);
    setUsersType("subscribers");
    setOpenUsers(true);
  };

  const loadEvents = async () => {
    // const users = await api.users.getAll();
    // const subscribers = users.filter((user) => user.is_subscribe_newsletter);
    // if (subscribers.length <= 0) {
    //   toaster.error({
    //     description: "No members in club",
    //   });
    // }
    // setUsers(subscribers);
    setOpenEvents(true);
  };

  return (
    <Container>
      <Heading textAlign={"center"} pb={3}>
        Dashboard
      </Heading>
      <Center>
        <Stack w="lg">
          <Button onClick={loadClubMembers}>Show club members</Button>
          <Button onClick={loadNewsletterSubscribers}>
            Show newsletter subscribers
          </Button>
          <Button onClick={loadEvents}>Show events list</Button>
        </Stack>
      </Center>
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
              <Dialog.Body>EVENTS</Dialog.Body>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </Container>
  );
};
