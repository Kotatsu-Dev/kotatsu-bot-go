import { handleError, useAPI } from "../../api/api";
import {
  Button,
  Card,
  Checkbox,
  CloseButton,
  Container,
  DataList,
  Dialog,
  Field,
  Fieldset,
  Group,
  Heading,
  Input,
  Link,
  Portal,
  RadioGroup,
  Stack,
  Status,
  Tabs,
} from "@chakra-ui/react";
import { memo, useCallback, useEffect, useMemo, useState } from "react";
import { toaster } from "../ui/toaster";
import type { User } from "@/api/users";
import { formatDate, formatDistanceToNow } from "date-fns";
import { PaginatedList } from "./PaginatedList";
import { useDebounceValue } from "usehooks-ts";
import Fuse from "fuse.js";
import { Controller, useForm, type SubmitHandler } from "react-hook-form";

type Inputs = {
  full_name: string;
  phone_number: string;
  is_club_member: boolean;
  gender: string;
};

const genders = [
  { value: "male", label: "Male" },
  { value: "female", label: "Female" },
  { value: "", label: "Unknown" },
];

const UserEditDialog = (props: { value: User; reload: () => void }) => {
  const api = useAPI();
  const [open, setOpen] = useState(false);
  const { register, handleSubmit, control } = useForm<Inputs>({
    defaultValues: {
      full_name: props.value.full_name,
      phone_number: props.value.phone_number,
      is_club_member: props.value.is_club_member,
      gender: props.value.gender ?? "",
    },
  });

  const editUser: SubmitHandler<Inputs> = async (data, event) => {
    try {
      await api.users.update({ ...data, user_tg_id: props.value.user_tg_id });
      toaster.success({
        description: "Event successfully created!",
      });
      console.log(data);
      setOpen(false);
      props.reload();
    } catch (e) {
      handleError(e);
      event?.stopPropagation();
    }
  };
  return (
    <Dialog.Root open={open} onOpenChange={({ open }) => setOpen(open)}>
      <Dialog.Trigger asChild>
        <Button colorPalette={"green"}>Edit</Button>
      </Dialog.Trigger>
      <Portal>
        <Dialog.Backdrop />
        <Dialog.Positioner colorPalette={"orange"}>
          <Dialog.Content>
            <Dialog.Header>
              <Dialog.Title>Edit user</Dialog.Title>
            </Dialog.Header>
            <Dialog.Body as={"form"} onSubmit={handleSubmit(editUser)}>
              <Fieldset.Root>
                <Fieldset.Content>
                  <Field.Root>
                    <Field.Label>Full name</Field.Label>
                    <Input {...register("full_name")} />
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>Phone</Field.Label>
                    <Input {...register("phone_number")} />
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>ITMO status</Field.Label>
                    <RadioGroup.Root>
                      <Stack>
                        <RadioGroup.Item value={"guest"}>
                          <RadioGroup.ItemHiddenInput />
                          <RadioGroup.ItemIndicator />
                          <RadioGroup.ItemText>Guest</RadioGroup.ItemText>
                        </RadioGroup.Item>
                        <RadioGroup.Item value={"student"}>
                          <RadioGroup.ItemHiddenInput />
                          <RadioGroup.ItemIndicator />
                          <RadioGroup.ItemText>Student</RadioGroup.ItemText>
                        </RadioGroup.Item>
                        <RadioGroup.Item value={"employee"}>
                          <RadioGroup.ItemHiddenInput />
                          <RadioGroup.ItemIndicator />
                          <RadioGroup.ItemText>Employee</RadioGroup.ItemText>
                        </RadioGroup.Item>
                      </Stack>
                    </RadioGroup.Root>
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>Gender</Field.Label>
                    <RadioGroup.Root>
                      <Controller
                        name="gender"
                        control={control}
                        render={({ field }) => (
                          <RadioGroup.Root
                            name={field.name}
                            value={field.value}
                            onValueChange={({ value }) => {
                              field.onChange(value);
                            }}
                          >
                            <Stack>
                              {genders.map((gender) => (
                                <RadioGroup.Item
                                  key={gender.value}
                                  value={gender.value}
                                >
                                  <RadioGroup.ItemHiddenInput
                                    onBlur={field.onBlur}
                                  />
                                  <RadioGroup.ItemIndicator />
                                  <RadioGroup.ItemText>
                                    {gender.label}
                                  </RadioGroup.ItemText>
                                </RadioGroup.Item>
                              ))}
                            </Stack>
                          </RadioGroup.Root>
                        )}
                      />
                    </RadioGroup.Root>
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>Status</Field.Label>
                    <Controller
                      control={control}
                      name="is_club_member"
                      render={({ field }) => (
                        <Checkbox.Root
                          checked={field.value}
                          onCheckedChange={({ checked }) =>
                            field.onChange(checked)
                          }
                          onBlur={field.onBlur}
                        >
                          <Checkbox.HiddenInput
                            ref={field.ref}
                            name={field.name}
                          />
                          <Checkbox.Control />
                          <Checkbox.Label>Member</Checkbox.Label>
                        </Checkbox.Root>
                      )}
                    />
                  </Field.Root>
                </Fieldset.Content>
                <Group>
                  <Button
                    variant="outline"
                    type="button"
                    onClick={() => setOpen(false)}
                  >
                    Cancel
                  </Button>
                  <Button type="submit">Save</Button>
                </Group>
              </Fieldset.Root>
            </Dialog.Body>
            <Dialog.CloseTrigger asChild>
              <CloseButton size="sm" />
            </Dialog.CloseTrigger>
          </Dialog.Content>
        </Dialog.Positioner>
      </Portal>
    </Dialog.Root>
  );
};

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
                  Member{" "}
                  {user.club_member_since
                    ? `(since ${formatDate(user.club_member_since, "dd-MM-yyyy")})`
                    : null}
                </Status.Root>
              ) : user.my_request ? (
                <Status.Root colorPalette={"yellow"}>
                  <Status.Indicator />
                  {`Submitted (${formatDistanceToNow(
                    user.my_request.created_at,
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
            <DataList.ItemLabel>Gender</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.gender ? (
                genders.find((v) => v.value == user.gender)?.label
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
                <Link target="_blank" href={`https://t.me/${user.user_name}`}>
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
        <Group grow w={"full"}>
          {user.my_request ? (
            user.my_request.status == 0 ? (
              <>
                <Button colorPalette={"red"} onClick={rejectRequest}>
                  Reject
                </Button>
                <Button colorPalette={"green"} onClick={acceptRequest}>
                  Accept
                </Button>
              </>
            ) : null
          ) : user.is_club_member ? (
            <Button colorPalette={"red"} onClick={kickUser}>
              Kick user
            </Button>
          ) : (
            <Button colorPalette={"green"} onClick={addUser}>
              Add user to club
            </Button>
          )}
          {/* <Button colorPalette={"red"}>Delete</Button> */}
          <UserEditDialog {...props} />
        </Group>
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
    [users, search],
  );

  const members = useMemo(
    () => currentUsers.filter((user) => user.is_club_member),
    [currentUsers],
  );
  const requests = useMemo(
    () => currentUsers.filter((user) => !!user.my_request),
    [currentUsers],
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
