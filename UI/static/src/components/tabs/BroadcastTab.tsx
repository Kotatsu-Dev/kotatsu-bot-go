import { handleError, useAPI } from "../../api/api";
import {
  Button,
  Combobox,
  Container,
  createListCollection,
  Field,
  Fieldset,
  Heading,
  Listbox,
  Portal,
  Textarea,
  useFilter,
  useListCollection,
} from "@chakra-ui/react";
import { Controller, useForm } from "react-hook-form";
import { useEffect } from "react";
import { toaster } from "../ui/toaster";

const clubMemberStatus = createListCollection({
  items: [
    { label: "Club member", value: "club_member" },
    { label: "Not club member", value: "not_club_member" },
  ],
});

const itmoStatus = createListCollection({
  items: [
    { label: "Guest", value: "guest" },
    { label: "Student", value: "student" },
    { label: "Graduate", value: "graduate" },
    { label: "Employee", value: "employee" },
    { label: "Student and employee", value: "student_employee" },
    { label: "Graduate and employee", value: "graduate_employee" },
    { label: "Unknown", value: "" },
  ],
});

export const BroadcastTab = () => {
  const api = useAPI();
  const { contains } = useFilter({ sensitivity: "base" });

  type BroadcastFormValues = {
    events: number[];
    users: number[];
    roulettes: number[];
    club_member_status: string[];
    itmo_status: string[];
    message: string;
  };

  const { control, handleSubmit, register, reset } =
    useForm<BroadcastFormValues>({
      defaultValues: {
        events: [],
        users: [],
        roulettes: [],
        club_member_status: [],
        itmo_status: [],
        message: "",
      },
    });

  const onSubmit = async (data: BroadcastFormValues) => {
    try {
      await api.broadcast.sendBroadcast({
        events: data.events,
        users: data.users,
        roulettes: data.roulettes,
        club_member_status: transformClubMemberStatus(data.club_member_status),
        itmo_status: data.itmo_status,
        message: data.message,
      });

      toaster.success({
        title: "Broadcast sent successfully",
      });

      reset();
    } catch (error) {
      handleError(error);
    }
  };

  // Transform club member status from string array to boolean or null
  const transformClubMemberStatus = (statusArray: string[]): boolean | null => {
    const isClubMember = statusArray.includes("club_member");
    const isNotClubMember = statusArray.includes("not_club_member");

    if (isClubMember && !isNotClubMember) {
      return true;
    }

    if (!isClubMember && isNotClubMember) {
      return false;
    }

    return null;
  };

  const {
    collection: events,
    filter: filterEvents,
    set: setEvents,
  } = useListCollection<{ label: string; value: number }>({
    initialItems: [],
    filter: contains,
  });

  const {
    collection: users,
    filter: filterUsers,
    set: setUsers,
  } = useListCollection<{ label: string; value: number }>({
    initialItems: [],
    filter: contains,
  });

  const {
    collection: roulettes,
    filter: filterRoulettes,
    set: setRoulettes,
  } = useListCollection<{ label: string; value: number }>({
    initialItems: [],
    filter: contains,
  });

  const loadEvents = async () => {
    const allEvents = await api.activities.getAll();
    const sortedEvents = allEvents.sort((a, b) => b.id - a.id);
    setEvents(
      sortedEvents.map((event) => ({
        label: event.title,
        value: event.id,
      })),
    );
  };

  const loadUsers = async () => {
    const allUsers = await api.users.getAll();
    setUsers(
      allUsers.map((user) => {
        const parts = [];

        if (user.full_name && user.full_name.trim()) {
          parts.push(user.full_name);
        }

        if (user.user_name && user.user_name.trim()) {
          parts.push(`@${user.user_name}`);
        }

        parts.push(`(${user.user_tg_id})`);

        return {
          label: parts.join(" "),
          value: user.id,
        };
      }),
    );
  };

  const loadRoulettes = async () => {
    const allRoulettes = await api.roulettes.getAll();
    setRoulettes(
      allRoulettes.map((roulette) => ({
        label: roulette.theme,
        value: roulette.id,
      })),
    );
  };

  useEffect(() => {
    loadEvents();
    loadUsers();
    loadRoulettes();
  }, []);

  return (
    <Container maxW={"lg"} mb={5} as="form" onSubmit={handleSubmit(onSubmit)}>
      <Heading textAlign={"center"} pb={3}>
        Broadcast to users
      </Heading>
      <Fieldset.Root>
        <Fieldset.Content>
          <Field.Root>
            <Field.Label>Events</Field.Label>
            <Controller
              name="events"
              control={control}
              render={({ field }) => (
                <Combobox.Root
                  collection={events}
                  onInputValueChange={(e) => filterEvents(e.inputValue)}
                  onValueChange={({ value }) => field.onChange(value)}
                  value={field.value as any as string[]}
                >
                  <Combobox.Control>
                    <Combobox.Input placeholder="Type to search" />
                    <Combobox.IndicatorGroup>
                      <Combobox.ClearTrigger />
                      <Combobox.Trigger />
                    </Combobox.IndicatorGroup>
                  </Combobox.Control>
                  <Portal>
                    <Combobox.Positioner>
                      <Combobox.Content>
                        <Combobox.Empty>No items found</Combobox.Empty>
                        {events.items.map((item) => (
                          <Combobox.Item item={item} key={item.value}>
                            {item.label}
                            <Combobox.ItemIndicator />
                          </Combobox.Item>
                        ))}
                      </Combobox.Content>
                    </Combobox.Positioner>
                  </Portal>
                </Combobox.Root>
              )}
            />
          </Field.Root>
          <Field.Root>
            <Field.Label>Users</Field.Label>
            <Controller
              name="users"
              control={control}
              render={({ field }) => (
                <Combobox.Root
                  collection={users}
                  onInputValueChange={(e) => filterUsers(e.inputValue)}
                  onValueChange={({ value }) => field.onChange(value)}
                  value={field.value as any as string[]}
                >
                  <Combobox.Control>
                    <Combobox.Input placeholder="Type to search" />
                    <Combobox.IndicatorGroup>
                      <Combobox.ClearTrigger />
                      <Combobox.Trigger />
                    </Combobox.IndicatorGroup>
                  </Combobox.Control>
                  <Portal>
                    <Combobox.Positioner>
                      <Combobox.Content>
                        <Combobox.Empty>No items found</Combobox.Empty>
                        {users.items.map((item) => (
                          <Combobox.Item item={item} key={item.value}>
                            {item.label}
                            <Combobox.ItemIndicator />
                          </Combobox.Item>
                        ))}
                      </Combobox.Content>
                    </Combobox.Positioner>
                  </Portal>
                </Combobox.Root>
              )}
            />
          </Field.Root>
          <Field.Root>
            <Field.Label>Roulettes</Field.Label>
            <Controller
              name="roulettes"
              control={control}
              render={({ field }) => (
                <Combobox.Root
                  collection={roulettes}
                  onInputValueChange={(e) => filterRoulettes(e.inputValue)}
                  onValueChange={({ value }) => field.onChange(value)}
                  value={field.value as any as string[]}
                >
                  <Combobox.Control>
                    <Combobox.Input placeholder="Type to search" />
                    <Combobox.IndicatorGroup>
                      <Combobox.ClearTrigger />
                      <Combobox.Trigger />
                    </Combobox.IndicatorGroup>
                  </Combobox.Control>
                  <Portal>
                    <Combobox.Positioner>
                      <Combobox.Content>
                        <Combobox.Empty>No items found</Combobox.Empty>
                        {roulettes.items.map((item) => (
                          <Combobox.Item item={item} key={item.value}>
                            {item.label}
                            <Combobox.ItemIndicator />
                          </Combobox.Item>
                        ))}
                      </Combobox.Content>
                    </Combobox.Positioner>
                  </Portal>
                </Combobox.Root>
              )}
            />
          </Field.Root>
          <Field.Root>
            <Field.Label>Club membership</Field.Label>
            <Controller
              name="club_member_status"
              control={control}
              render={({ field }) => (
                <Listbox.Root
                  collection={clubMemberStatus}
                  selectionMode="multiple"
                  value={field.value as any as string[]}
                  onValueChange={({ value }) => field.onChange(value)}
                >
                  <Listbox.Content>
                    {clubMemberStatus.items.map((membership) => (
                      <Listbox.Item item={membership} key={membership.value}>
                        <Listbox.ItemText>{membership.label}</Listbox.ItemText>
                        <Listbox.ItemIndicator />
                      </Listbox.Item>
                    ))}
                  </Listbox.Content>
                </Listbox.Root>
              )}
            />
          </Field.Root>
          <Field.Root>
            <Field.Label>ITMO status</Field.Label>
            <Controller
              name="itmo_status"
              control={control}
              render={({ field }) => (
                <Listbox.Root
                  collection={itmoStatus}
                  selectionMode="multiple"
                  value={field.value as any as string[]}
                  onValueChange={({ value }) => field.onChange(value)}
                >
                  <Listbox.Content>
                    {itmoStatus.items.map((itmoStatus) => (
                      <Listbox.Item item={itmoStatus} key={itmoStatus.value}>
                        <Listbox.ItemText>{itmoStatus.label}</Listbox.ItemText>
                        <Listbox.ItemIndicator />
                      </Listbox.Item>
                    ))}
                  </Listbox.Content>
                </Listbox.Root>
              )}
            />
          </Field.Root>
          <Field.Root>
            <Field.Label>Message</Field.Label>
            <Textarea {...register("message")} />
          </Field.Root>
          <Button type="submit">Submit</Button>
        </Fieldset.Content>
      </Fieldset.Root>
    </Container>
  );
};
