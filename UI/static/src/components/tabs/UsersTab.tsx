import { handleError, useAPI } from "../../api/api";
import {
  Button,
  Card,
  Container,
  DataList,
  Heading,
  Input,
  Link,
  SimpleGrid,
  Stack,
  Status,
  Tabs,
} from "@chakra-ui/react";
import { memo, useCallback, useEffect, useMemo, useState } from "react";
import { toaster } from "../ui/toaster";
import type { User } from "@/api/users";
import { formatDistanceToNow } from "date-fns";
import { PaginatedList } from "./PaginatedList";
import { useDebounceValue } from "usehooks-ts";
import Fuse from "fuse.js";

const UserCard = memo((props: { value: User; reload: () => void }) => {
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
                <Status.Root colorPalette={"green"}>
                  <Status.Indicator />
                  Member
                </Status.Root>
              ) : user.my_request ? (
                <Status.Root colorPalette={"yellow"}>
                  <Status.Indicator />
                  {`Submitted (${formatDistanceToNow(
                    user.my_request.created_at
                  )} ago)`}
                </Status.Root>
              ) : (
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  Not member
                </Status.Root>
              )}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Full name</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.full_name ? (
                user.full_name
              ) : (
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  Unknown
                </Status.Root>
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
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  None
                </Status.Root>
              )}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>From ITMO</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.is_itmo ? (
                <Status.Root colorPalette={"green"}>
                  <Status.Indicator />
                  Yes
                </Status.Root>
              ) : (
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  No
                </Status.Root>
              )}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>ISU</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.isu ? (
                user.isu
              ) : (
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  None
                </Status.Root>
              )}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Phone</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.phone_number.length ? (
                user.phone_number
              ) : (
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  None
                </Status.Root>
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
});

export const UsersTab = () => {
  const api = useAPI();
  const [users, setUsers] = useState<User[]>([]);
  const [search, setSearch] = useDebounceValue("", 500);

  const currentUsers = useMemo(
    () =>
      search.length == 0
        ? users
        : new Fuse(users, {
            keys: [
              "user_name",
              "full_tg_name",
              "isu",
              "full_name",
              "phone_number",
            ],
          })
            .search(search)
            .map((e) => e.item),
    [users, search]
  );

  const members = useMemo(
    () => currentUsers.filter((user) => user.is_club_member),
    [currentUsers]
  );
  const requests = useMemo(
    () => currentUsers.filter((user) => !!user.my_request),
    [currentUsers]
  );

  const loadUsers = useCallback(async () => {
    setUsers(await api.users.getAll());
  }, []);

  useEffect(() => {
    loadUsers();
  }, []);

  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>User management</Heading>
        <Input
          defaultValue={search}
          onChange={(e) => setSearch(e.currentTarget.value)}
        />
        <Tabs.Root fitted variant={"enclosed"} defaultValue={"all"}>
          <Tabs.List>
            <Tabs.Trigger value="all">Everybody</Tabs.Trigger>
            <Tabs.Trigger value="members">Members</Tabs.Trigger>
            <Tabs.Trigger value="requests">Requests</Tabs.Trigger>
          </Tabs.List>
          <Tabs.Content value="all">
            <PaginatedList
              items={currentUsers}
              pageSize={5}
              render={(user) => (
                <UserCard key={user.id} value={user} reload={loadUsers} />
              )}
            />
          </Tabs.Content>
          <Tabs.Content value="members">
            <PaginatedList
              items={members}
              pageSize={5}
              render={(user) => (
                <UserCard key={user.id} value={user} reload={loadUsers} />
              )}
            />
          </Tabs.Content>
          <Tabs.Content value="requests">
            <PaginatedList
              items={requests}
              pageSize={5}
              render={(user) => (
                <UserCard key={user.id} value={user} reload={loadUsers} />
              )}
            />
          </Tabs.Content>
        </Tabs.Root>
      </Stack>
    </Container>
  );
};
