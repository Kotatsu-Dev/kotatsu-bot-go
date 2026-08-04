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
  Listbox,
  Portal,
  RadioGroup,
  Stack,
  Status,
  Text,
} from "@chakra-ui/react";
import { memo, useCallback, useEffect, useMemo, useState } from "react";
import { toaster } from "../ui/toaster";
import type { User } from "@/api/users";
import { formatDate, formatDistanceToNow } from "date-fns";
import { PaginatedList } from "./PaginatedList";
import { useDebounceValue } from "usehooks-ts";
import { Controller, useForm, type SubmitHandler } from "react-hook-form";
import {
  genders,
  itmoStatuses as statuses,
  itmoTraitCollection,
  traitsOf,
  type ItmoTrait,
} from "../../constants/users";
import { searchUsers } from "../../lib/userSearch";

type Inputs = {
  full_name: string;
  phone_number: string;
  is_club_member: boolean;
  gender: string;
  itmo_status: string;
};

const PAGE_SIZE = 10;

type ClubFilter = "all" | "member" | "not_member";
type GenderFilter = "all" | "male" | "female" | "unknown";

const UserEditDialog = (props: { value: User; reload: () => void }) => {
  const api = useAPI();
  const [open, setOpen] = useState(false);
  const { register, handleSubmit, control } = useForm<Inputs>({
    defaultValues: {
      full_name: props.value.full_name,
      phone_number: props.value.phone_number,
      is_club_member: props.value.is_club_member,
      gender: props.value.gender ?? "",
      itmo_status: props.value.itmo_status ?? "",
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
                    <Controller
                      name="itmo_status"
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
                            {statuses.map((status) => (
                              <RadioGroup.Item
                                key={status.value}
                                value={status.value}
                              >
                                <RadioGroup.ItemHiddenInput
                                  onBlur={field.onBlur}
                                />
                                <RadioGroup.ItemIndicator />
                                <RadioGroup.ItemText>
                                  {status.label}
                                </RadioGroup.ItemText>
                              </RadioGroup.Item>
                            ))}
                          </Stack>
                        </RadioGroup.Root>
                      )}
                    />
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>Gender</Field.Label>
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
            <DataList.ItemLabel>ITMO status</DataList.ItemLabel>
            <DataList.ItemValue>
              {user.itmo_status ? (
                statuses.find((v) => v.value == user.itmo_status)?.label
              ) : (
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  Unknown
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
  const [searchInput, setSearchInput] = useState("");
  const [search, setSearch] = useDebounceValue("", 500);
  const [clubFilter, setClubFilter] = useState<ClubFilter>("all");
  const [genderFilter, setGenderFilter] = useState<GenderFilter>("all");
  const [itmoTraitFilter, setItmoTraitFilter] = useState<string[]>([]);
  const [onlyRequests, setOnlyRequests] = useState(false);

  const currentUsers = useMemo(() => {
    let result = searchUsers(users, search);

    if (clubFilter !== "all") {
      const wanted = clubFilter === "member";
      result = result.filter((user) => user.is_club_member === wanted);
    }

    if (genderFilter !== "all") {
      const wanted = genderFilter === "unknown" ? "" : genderFilter;
      result = result.filter((user) => (user.gender ?? "") === wanted);
    }

    if (itmoTraitFilter.length > 0) {
      result = result.filter((user) => {
        const traits = traitsOf(user.itmo_status);
        return itmoTraitFilter.every((trait) =>
          traits.includes(trait as ItmoTrait),
        );
      });
    }

    if (onlyRequests) {
      result = result.filter((user) => !!user.my_request);
    }

    return result;
  }, [
    users,
    search,
    clubFilter,
    genderFilter,
    itmoTraitFilter,
    onlyRequests,
  ]);

  const resetFilters = () => {
    setSearchInput("");
    setSearch("");
    setClubFilter("all");
    setGenderFilter("all");
    setItmoTraitFilter([]);
    setOnlyRequests(false);
  };

  const loadUsers = useCallback(async () => {
    setUsers(await api.users.getAll());
  }, []);

  useEffect(() => {
    loadUsers();
  }, []);

  return (
    <Container maxW={"lg"} mb={5}>
      <Stack>
        <Heading textAlign={"center"}>User management</Heading>
        <Input
          placeholder="Search by name, ISU, phone, @username or Telegram ID"
          value={searchInput}
          onChange={(e) => {
            setSearchInput(e.currentTarget.value);
            setSearch(e.currentTarget.value);
          }}
        />
        <Card.Root>
          <Card.Body>
            <Stack gap={4}>
              <Field.Root>
                <Field.Label>Club membership</Field.Label>
                <RadioGroup.Root
                  value={clubFilter}
                  onValueChange={({ value }) =>
                    setClubFilter((value as ClubFilter) ?? "all")
                  }
                >
                  <Group wrap={"wrap"}>
                    <RadioGroup.Item value="all">
                      <RadioGroup.ItemHiddenInput />
                      <RadioGroup.ItemIndicator />
                      <RadioGroup.ItemText>All</RadioGroup.ItemText>
                    </RadioGroup.Item>
                    <RadioGroup.Item value="member">
                      <RadioGroup.ItemHiddenInput />
                      <RadioGroup.ItemIndicator />
                      <RadioGroup.ItemText>Members</RadioGroup.ItemText>
                    </RadioGroup.Item>
                    <RadioGroup.Item value="not_member">
                      <RadioGroup.ItemHiddenInput />
                      <RadioGroup.ItemIndicator />
                      <RadioGroup.ItemText>Non-members</RadioGroup.ItemText>
                    </RadioGroup.Item>
                  </Group>
                </RadioGroup.Root>
              </Field.Root>

              <Field.Root>
                <Field.Label>Gender</Field.Label>
                <RadioGroup.Root
                  value={genderFilter}
                  onValueChange={({ value }) =>
                    setGenderFilter((value as GenderFilter) ?? "all")
                  }
                >
                  <Group wrap={"wrap"}>
                    <RadioGroup.Item value="all">
                      <RadioGroup.ItemHiddenInput />
                      <RadioGroup.ItemIndicator />
                      <RadioGroup.ItemText>All</RadioGroup.ItemText>
                    </RadioGroup.Item>
                    <RadioGroup.Item value="male">
                      <RadioGroup.ItemHiddenInput />
                      <RadioGroup.ItemIndicator />
                      <RadioGroup.ItemText>Male</RadioGroup.ItemText>
                    </RadioGroup.Item>
                    <RadioGroup.Item value="female">
                      <RadioGroup.ItemHiddenInput />
                      <RadioGroup.ItemIndicator />
                      <RadioGroup.ItemText>Female</RadioGroup.ItemText>
                    </RadioGroup.Item>
                    <RadioGroup.Item value="unknown">
                      <RadioGroup.ItemHiddenInput />
                      <RadioGroup.ItemIndicator />
                      <RadioGroup.ItemText>Unknown</RadioGroup.ItemText>
                    </RadioGroup.Item>
                  </Group>
                </RadioGroup.Root>
              </Field.Root>

              <Field.Root>
                <Field.Label>ITMO status</Field.Label>
                <Field.HelperText>
                  All selected traits must match. Nothing selected means no
                  filtering.
                </Field.HelperText>
                <Listbox.Root
                  collection={itmoTraitCollection}
                  selectionMode="multiple"
                  value={itmoTraitFilter}
                  onValueChange={({ value }) => setItmoTraitFilter(value)}
                >
                  <Listbox.Content>
                    {itmoTraitCollection.items.map((trait) => (
                      <Listbox.Item item={trait} key={trait.value}>
                        <Listbox.ItemText>{trait.label}</Listbox.ItemText>
                        <Listbox.ItemIndicator />
                      </Listbox.Item>
                    ))}
                  </Listbox.Content>
                </Listbox.Root>
              </Field.Root>

              <Checkbox.Root
                checked={onlyRequests}
                onCheckedChange={({ checked }) => setOnlyRequests(!!checked)}
              >
                <Checkbox.HiddenInput />
                <Checkbox.Control />
                <Checkbox.Label>Pending request only</Checkbox.Label>
              </Checkbox.Root>

              <Button variant={"outline"} onClick={resetFilters}>
                Reset filters
              </Button>
            </Stack>
          </Card.Body>
        </Card.Root>

        <Text textAlign={"center"} color={"fg.muted"}>
          Showing {currentUsers.length} of {users.length} users
        </Text>

        {currentUsers.length === 0 ? (
          <Text textAlign={"center"} py={6}>
            No users match these filters
          </Text>
        ) : (
          <PaginatedList
            items={currentUsers}
            pageSize={PAGE_SIZE}
            render={(user) => (
              <UserCard key={user.id} value={user} reload={loadUsers} />
            )}
          />
        )}
      </Stack>
    </Container>
  );
};
