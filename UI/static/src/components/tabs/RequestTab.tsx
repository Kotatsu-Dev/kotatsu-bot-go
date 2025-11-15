import { handleError, useAPI } from "../../api/api";
import { type User } from "../../api/users";
import {
  Badge,
  Button,
  Card,
  CloseButton,
  Container,
  DataList,
  Dialog,
  Heading,
  Input,
  Link,
  NativeSelect,
  Portal,
  Separator,
  Stack,
} from "@chakra-ui/react";
import { useState } from "react";
import { toaster } from "../ui/toaster";

export const RequestTab = () => {
  const api = useAPI();
  const [open, setOpen] = useState(false);
  const [category, setCategory] = useState<"all" | "itmo" | "not_itmo">("all");
  const [users, setUsers] = useState<User[]>([]);
  const [requestId, setRequestId] = useState(0);

  const fetchUsers = async () => {
    try {
      const fetchedUsers = await api.users.getAll();
      const requestUsers = fetchedUsers.filter((user) => !!user.my_request);
      let filteredUsers;
      switch (category) {
        case "all":
          filteredUsers = requestUsers;
          break;
        case "itmo":
          filteredUsers = requestUsers.filter((user) => user.is_itmo);
          break;
        case "not_itmo":
          filteredUsers = requestUsers.filter((user) => !user.is_itmo);
          break;
      }
      setUsers(filteredUsers);
      if (filteredUsers.length > 0) {
        setOpen(true);
      } else {
        toaster.error({
          description: "No requests in this category",
        });
      }
    } catch (e) {
      handleError(e);
    }
  };

  const acceptRequest = async () => {
    try {
      await api.requests.accept({ id: requestId });
      toaster.success({
        description: "Succesfully accepted request",
      });
    } catch (e) {
      handleError(e);
    }
  };

  const rejectRequest = async () => {
    try {
      await api.requests.reject({ id: requestId });
      toaster.success({
        description: "Succesfully accepted request",
      });
    } catch (e) {
      handleError(e);
    }
  };

  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>View club requests</Heading>
        <NativeSelect.Root>
          <NativeSelect.Field
            value={category}
            onChange={(e) => setCategory(e.currentTarget.value as any)}
          >
            <option value={"all"}>All users</option>
            <option value={"itmo"}>Users from ITMO</option>
            <option value={"not_itmo"}>Users not from ITMO</option>
          </NativeSelect.Field>
          <NativeSelect.Indicator />
        </NativeSelect.Root>
        <Button onClick={fetchUsers}>Fetch requests</Button>
        <Separator />
        <Heading textAlign={"center"}>Request approval form</Heading>
        <Input
          placeholder="Enter request ID"
          type="number"
          value={requestId}
          onChange={(e) => setRequestId(e.currentTarget.valueAsNumber)}
        />
        <Stack direction={"row"}>
          <Button
            flexGrow={1}
            colorPalette={"red"}
            variant={"outline"}
            onClick={rejectRequest}
          >
            Reject
          </Button>
          <Button flexGrow={1} colorPalette={"green"} onClick={acceptRequest}>
            Accept
          </Button>
        </Stack>
      </Stack>
      <Dialog.Root open={open} onOpenChange={(e) => setOpen(e.open)}>
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.Header>
                <Dialog.Title>User requests</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <Stack>
                  {users.map((user) => (
                    <Card.Root key={user.id}>
                      <Card.Header>
                        <Heading>Request {user.my_request?.id}</Heading>
                      </Card.Header>
                      <Card.Body>
                        <DataList.Root orientation={"horizontal"}>
                          <DataList.Item>
                            <DataList.ItemLabel>Request ID</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {user.my_request?.id}
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>
                              Request date
                            </DataList.ItemLabel>
                            <DataList.ItemValue>
                              {user.my_request?.created_at}
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>Telegram ID</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {user.user_tg_id}
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>Username</DataList.ItemLabel>
                            <DataList.ItemValue>
                              <Link target="_blank" href={`https://t.me/${user.user_name}`}>
                                @{user.user_name}
                              </Link>
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>From ITMO</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {user.is_itmo ? (
                                <Badge colorPalette={"green"}>Yes</Badge>
                              ) : (
                                <Badge colorPalette={"red"}>No</Badge>
                              )}
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>ISU</DataList.ItemLabel>
                            <DataList.ItemValue>{user.isu}</DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>Full name</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {user.full_name}
                            </DataList.ItemValue>
                          </DataList.Item>
                          <DataList.Item>
                            <DataList.ItemLabel>Phone</DataList.ItemLabel>
                            <DataList.ItemValue>
                              {user.phone_number}
                            </DataList.ItemValue>
                          </DataList.Item>
                        </DataList.Root>
                      </Card.Body>
                    </Card.Root>
                  ))}
                </Stack>
              </Dialog.Body>
              <Dialog.CloseTrigger asChild>
                <CloseButton size="sm" />
              </Dialog.CloseTrigger>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </Container>
  );
};
