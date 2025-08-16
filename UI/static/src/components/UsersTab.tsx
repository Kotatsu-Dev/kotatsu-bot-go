import { handleError, useAPI } from "../api/api";
import {
  Badge,
  Button,
  Card,
  Container,
  DataList,
  Heading,
  Link,
  SimpleGrid,
  Stack,
  Tabs,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { toaster } from "./ui/toaster";
import type { User } from "@/api/users";
import { formatDistanceToNow } from "date-fns";

const UserComponent = (props: { value: User; reload: () => void }) => {
  const api = useAPI();
  const user = props.value;

  const acceptRequest = async () => {
    if (user.my_request) {
      await api.requests.accept({ id: user.my_request.id });
      toaster.success({
        description: "Succesfully accepted request",
      });
      props.reload();
    }
  };

  const rejectRequest = async () => {
    if (user.my_request) {
      await api.requests.reject({ id: user.my_request.id });
      toaster.success({
        description: "Succesfully accepted request",
      });
      props.reload();
    }
  };

  const kickUser = async () => {
    try {
      await api.users.setMemberStatus({
        user_tg_id: user.user_tg_id,
        is_club_member: false,
      });
      toaster.success({
        description: "Successfully kicked user",
      });
    } catch (e) {
      handleError(e);
    }
    props.reload();
  };

  const addUser = async () => {
    try {
      await api.users.setMemberStatus({
        user_tg_id: user.user_tg_id,
        is_club_member: true,
      });
      toaster.success({
        description: "Successfully added user",
      });
    } catch (e) {
      handleError(e);
    }
    props.reload();
  };

  return (
    <Card.Root key={user.id}>
      <Card.Body>
        <DataList.Root orientation={"horizontal"}>
          <DataList.Item>
            <DataList.ItemLabel>Status</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.is_club_member ? (
                <Badge colorPalette={"green"}>Member</Badge>
              ) : user.my_request ? (
                <Badge colorPalette={"yellow"}>
                  Submitted ({formatDistanceToNow(user.my_request.created_at)}{" "}
                  ago)
                </Badge>
              ) : (
                <Badge colorPalette={"red"}>Not member</Badge>
              )}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Full name</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.full_name ? (
                user.full_name
              ) : (
                <Badge colorPalette={"red"}>Unknown</Badge>
              )}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Telegram ID</DataList.ItemLabel>
            <DataList.ItemValue>{user.user_tg_id}</DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Username</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.user_name.length > 0 ? (
                <Link href={`https://t.me/${user.user_name}`}>
                  @{user.user_name}
                </Link>
              ) : (
                <Badge colorPalette={"red"}>None</Badge>
              )}
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
            <DataList.ItemValue>
              {user.isu ? user.isu : <Badge colorPalette={"red"}>None</Badge>}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Phone</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.phone_number.length ? (
                user.phone_number
              ) : (
                <Badge colorPalette={"red"}>None</Badge>
              )}
            </DataList.ItemValue>
          </DataList.Item>
        </DataList.Root>
      </Card.Body>
      <Card.Footer>
        {user.my_request ? (
          user.my_request.status == 0 ? (
            <SimpleGrid columns={2} w="full" gap={"2"}>
              <Button colorPalette={"red"} onClick={rejectRequest}>
                Reject
              </Button>
              <Button colorPalette={"green"} onClick={acceptRequest}>
                Accept
              </Button>
            </SimpleGrid>
          ) : null
        ) : user.is_club_member ? (
          <Button w="full" colorPalette={"red"} onClick={kickUser}>
            Kick user
          </Button>
        ) : (
          <Button w="full" colorPalette={"green"} onClick={addUser}>
            Add user to club
          </Button>
        )}
      </Card.Footer>
    </Card.Root>
  );
};

export const UsersTab = () => {
  const api = useAPI();
  const [users, setUsers] = useState<User[]>([]);

  const members = users.filter((user) => user.is_club_member);
  const requests = users.filter((user) => !!user.my_request);

  const loadUsers = async () => {
    setUsers(await api.users.getAll());
  };

  useEffect(() => {
    loadUsers();
  }, []);

  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>User management</Heading>
        <Tabs.Root fitted variant={"enclosed"} defaultValue={"all"}>
          <Tabs.List>
            <Tabs.Trigger value="all">Everybody</Tabs.Trigger>
            <Tabs.Trigger value="members">Members</Tabs.Trigger>
            <Tabs.Trigger value="requests">Requests</Tabs.Trigger>
          </Tabs.List>
          <Tabs.Content value="all">
            {users.map((user) => (
              <Stack key={user.id}>
                <UserComponent value={user} reload={loadUsers} />
              </Stack>
            ))}
          </Tabs.Content>
          <Tabs.Content value="members">
            {members.map((user) => (
              <Stack key={user.id}>
                <UserComponent value={user} reload={loadUsers} />
              </Stack>
            ))}
          </Tabs.Content>
          <Tabs.Content value="requests">
            {requests.map((user) => (
              <Stack key={user.id}>
                <UserComponent value={user} reload={loadUsers} />
              </Stack>
            ))}
          </Tabs.Content>
        </Tabs.Root>
      </Stack>
    </Container>
  );
};
