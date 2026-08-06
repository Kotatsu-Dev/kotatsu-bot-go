import { handleError, useAPI } from "../../api/api";
import {
  Button,
  Card,
  Clipboard,
  CloseButton,
  Container,
  DataList,
  Dialog,
  DownloadTrigger,
  Field,
  Fieldset,
  FileUpload,
  Flex,
  Group,
  Heading,
  IconButton,
  Input,
  Portal,
  Stack,
  Status,
  Table,
  Tabs,
  Text,
  Textarea,
} from "@chakra-ui/react";
import { Controller, useForm, type SubmitHandler } from "react-hook-form";
import { toaster } from "../ui/toaster";
import { type Activity } from "../../api/activities";
import type { User } from "../../api/users";
import { useEffect, useMemo, useState } from "react";
import { isFuture, isPast } from "date-fns";
import { Workbook } from "exceljs";
import { FaDownload, FaEye, FaTimes } from "react-icons/fa";
import { Calendar } from "../Calendar";
import { PaginatedList } from "./PaginatedList";
import { useDebounceValue } from "usehooks-ts";
import { searchUsers } from "../../lib/userSearch";

const PAGE_SIZE = 10;

const sortByDateDesc = (list: Activity[]) =>
  [...list].sort(
    (a, b) =>
      new Date(b.date_meeting).getTime() - new Date(a.date_meeting).getTime(),
  );

const exportExcel = async (event: Activity) => {
  const wb = new Workbook();
  const sheet = wb.addWorksheet("СЗ");
  sheet.addRow([
    "Корпус:",
    null,
    "Дата, время:",
    null,
    "Название мероприятия:",
    null,
    "Ответственный подразделения:",
    "Контактное лицо:",
  ]);
  sheet.addRow([
    "№",
    "Фамилия",
    "Имя",
    "Отчество",
    "Серия и номер паспорта",
    "Номер телефона",
    null,
    null,
  ]);
  for (const [i, p] of event.participants.entries()) {
    const names = p.full_name.split(/\s+/);
    sheet.addRow([
      i + 1,
      names[0],
      names[1],
      names.slice(2).join(" "),
      null,
      p.phone_number,
    ]);
  }

  // Styles
  const grayCells = [
    "A1",
    "C1",
    "E1",
    "G1",
    "H1",
    "A2",
    "B2",
    "C2",
    "D2",
    "E2",
    "F2",
    "G2",
    "G3",
  ];

  for (const cell of grayCells) {
    sheet.getCell(cell).fill = {
      type: "pattern",
      pattern: "solid",
      fgColor: { argb: "E7E6E6" },
    };
  }

  for (let i = 0; i < Math.max(event.participants.length + 2, 3); i++) {
    const row = sheet.getRow(i + 1);
    for (let j = 0; j < (i < 3 ? 8 : 6); j++) {
      const cell = row.getCell(j + 1);
      cell.border = {
        top: { style: "thin" },
        left: { style: "thin" },
        bottom: { style: "thin" },
        right: { style: "thin" },
      };
      cell.font =
        i == 0
          ? {
              size: 11,
              color: { theme: 1 },
              name: "Calibri",
              family: 2,
              charset: 204,
              scheme: "minor",
            }
          : {
              size: 14,
              color: { theme: 1 },
              name: "Times New Roman",
              family: 1,
              charset: 204,
            };
    }
  }

  for (const cell of grayCells.slice(0, 3)) {
    sheet.getCell(cell).alignment = {
      horizontal: "right",
    };
  }

  sheet.getRow(2).eachCell((cell) => {
    cell.alignment = {
      ...cell.alignment,
      wrapText: true,
      vertical: "middle",
    };
  });

  for (let i = 0; i < event.participants.length; i++) {
    const row = sheet.getRow(i + 3);
    for (let j = 1; j < 6; j++) {
      const cell = row.getCell(j + 1);
      cell.alignment = {
        ...cell.alignment,
        horizontal: "center",
      };
    }
  }

  sheet.columns = [
    { width: 8 },
    { width: 19.6640625 },
    { width: 16.5 },
    { width: 21.83203125 },
    { width: 23.6640625 },
    { width: 27.83203125 },
    { width: 38.1640625 },
    { width: 29.6640625 },
  ];

  sheet.getRow(2).height = 76;

  const buff = await wb.xlsx.writeBuffer();
  const blob = new Blob([buff], {
    type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  });

  return blob;
};

const EventEditDialog = (props: { value: Activity; reload: () => void }) => {
  const api = useAPI();
  const [open, setOpen] = useState(false);
  const { register, handleSubmit, control } = useForm<Inputs>({
    defaultValues: {
      title: props.value.title,
      date_meeting: new Date(props.value.date_meeting),
      guest_registration_until: props.value.guest_registration_until
        ? new Date(props.value.guest_registration_until)
        : undefined,
      description: props.value.description,
      location: props.value.location,
      send_images: undefined,
      status: props.value.status,
    },
  });

  const editEvent: SubmitHandler<Inputs> = async (data, event) => {
    try {
      console.log(data);
      await api.activities.update({ ...data, id: props.value.id });
      toaster.success({
        description: "Event successfully edited!",
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
              <Dialog.Title>Edit event</Dialog.Title>
            </Dialog.Header>
            <Dialog.Body as={"form"} onSubmit={handleSubmit(editEvent)}>
              <Fieldset.Root>
                <Fieldset.Content>
                  <Field.Root>
                    <Field.Label>Title</Field.Label>
                    <Input
                      placeholder="Enter event title"
                      {...register("title")}
                    />
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>Date</Field.Label>
                    <Controller
                      control={control}
                      name="date_meeting"
                      render={({ field }) => <Calendar {...field} />}
                    />
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>Guest registration</Field.Label>
                    <Controller
                      control={control}
                      name="guest_registration_until"
                      render={({ field }) => <Calendar {...field} />}
                    />
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>Location</Field.Label>
                    <Input
                      placeholder="Enter location"
                      {...register("location")}
                    />
                  </Field.Root>
                  <Field.Root>
                    <Field.Label>Description</Field.Label>
                    <Textarea
                      placeholder="Enter event description"
                      {...register("description")}
                    />
                  </Field.Root>
                  <FileUpload.Root
                    maxFiles={5}
                    accept={"image/*"}
                    {...register("send_images")}
                  >
                    <FileUpload.HiddenInput />
                    <FileUpload.Trigger asChild>
                      <Button variant="outline" w="full">
                        Upload images for event
                      </Button>
                    </FileUpload.Trigger>
                    <FileUpload.List showSize clearable />
                  </FileUpload.Root>
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

const EventCard = (props: {
  value: Activity;
  reload: () => void;
  allUsers: User[];
}) => {
  const api = useAPI();
  const event = props.value;

  // event.participants comes back with unreliable id/created_at (backend TODO in
  // db/Activities.go ToRead()) — resolve each row via user_tg_id against the
  // canonical user list instead of trusting participant.id directly.
  const byTgId = useMemo(
    () => new Map(props.allUsers.map((u) => [u.user_tg_id, u])),
    [props.allUsers],
  );

  const participants = useMemo(() => {
    const result: User[] = [];
    for (const raw of event.participants) {
      const user = byTgId.get(raw.user_tg_id);
      if (user) result.push(user);
    }
    return result;
  }, [event.participants, byTgId]);

  const participantTgIds = useMemo(
    () => new Set(participants.map((u) => u.user_tg_id)),
    [participants],
  );

  const [participantSearchInput, setParticipantSearchInput] = useState("");
  const [participantSearch, setParticipantSearch] = useDebounceValue("", 500);

  const participantSearchResults = useMemo(() => {
    if (participantSearch.trim().length === 0) return [];
    return searchUsers(props.allUsers, participantSearch)
      .filter((u) => !participantTgIds.has(u.user_tg_id))
      .slice(0, 8);
  }, [props.allUsers, participantSearch, participantTgIds]);

  const addParticipant = async (user: User) => {
    try {
      await api.activities.addParticipant({
        activityId: event.id,
        userId: user.id,
      });
      toaster.success({ description: "Participant added" });
      setParticipantSearchInput("");
      setParticipantSearch("");
      props.reload();
    } catch (e) {
      handleError(e);
    }
  };

  const removeParticipant = async (user: User) => {
    try {
      await api.activities.removeParticipant({
        activityId: event.id,
        userId: user.id,
      });
      toaster.success({ description: "Participant removed" });
      props.reload();
    } catch (e) {
      handleError(e);
    }
  };

  const deactivateEvent = async (event: Activity) => {
    try {
      await api.activities.setStatus({ id: event.id, status: false });
      props.reload();
    } catch (e) {
      handleError(e);
    }
  };

  const reactivateEvent = async (event: Activity) => {
    try {
      await api.activities.setStatus({ id: event.id, status: true });
      props.reload();
    } catch (e) {
      handleError(e);
    }
  };

  return (
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
                <Status.Root colorPalette={"green"}>
                  <Status.Indicator />
                  Shown
                </Status.Root>
              ) : (
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  Hidden
                </Status.Root>
              )}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Link</DataList.ItemLabel>
            <DataList.ItemValue>
              <Clipboard.Root
                value={`https://t.me/${import.meta.env.VITE_BOT_USERNAME}?start=${props.value.id}`}
              >
                <Clipboard.Trigger asChild>
                  <IconButton variant="surface" size="xs">
                    <Clipboard.Indicator />
                  </IconButton>
                </Clipboard.Trigger>
              </Clipboard.Root>
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Description</DataList.ItemLabel>
            <DataList.ItemValue>{event.description}</DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Date</DataList.ItemLabel>
            <DataList.ItemValue>{event.date_meeting}</DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Guest registration</DataList.ItemLabel>
            <DataList.ItemValue>
              {event.guest_registration_until}
            </DataList.ItemValue>
          </DataList.Item>
          <DataList.Item>
            <DataList.ItemLabel>Location</DataList.ItemLabel>
            <DataList.ItemValue>{event.location}</DataList.ItemValue>
          </DataList.Item>
        </DataList.Root>
      </Card.Body>
      <Card.Footer>
        {event.status ? (
          <Button
            colorPalette={"red"}
            flexGrow={1}
            onClick={() => deactivateEvent(event)}
          >
            Hide
          </Button>
        ) : (
          <Button
            colorPalette={"red"}
            flexGrow={1}
            onClick={() => reactivateEvent(event)}
          >
            Show
          </Button>
        )}
        <EventEditDialog {...props} />
        <DownloadTrigger
          data={() => exportExcel(event)}
          fileName="СЗ.xlsx"
          mimeType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
          asChild
        >
          <IconButton aria-label="Download signed up" variant={"outline"}>
            <FaDownload />
          </IconButton>
        </DownloadTrigger>
        <Dialog.Root>
          <Dialog.Trigger asChild>
            <IconButton aria-label="Show signed up" variant={"outline"}>
              <FaEye />
            </IconButton>
          </Dialog.Trigger>
          <Portal>
            <Dialog.Backdrop />
            <Dialog.Positioner colorPalette={"orange"}>
              <Dialog.Content>
                <Dialog.Header>
                  <Dialog.Title>Signed up for "{event.title}"</Dialog.Title>
                </Dialog.Header>
                <Dialog.Body>
                  <Table.Root>
                    <Table.Header>
                      <Table.Row>
                        <Table.ColumnHeader>Name</Table.ColumnHeader>
                        <Table.ColumnHeader>Telegram</Table.ColumnHeader>
                        <Table.ColumnHeader>From ITMO</Table.ColumnHeader>
                        <Table.ColumnHeader>Phone Number</Table.ColumnHeader>
                        <Table.ColumnHeader></Table.ColumnHeader>
                      </Table.Row>
                    </Table.Header>
                    <Table.Body>
                      {participants.map((user) => (
                        <Table.Row key={user.id}>
                          <Table.Cell>{user.full_name}</Table.Cell>
                          <Table.Cell>{user.user_name}</Table.Cell>
                          <Table.Cell>
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
                          </Table.Cell>
                          <Table.Cell>{user.phone_number}</Table.Cell>
                          <Table.Cell>
                            <IconButton
                              aria-label="Remove participant"
                              size="xs"
                              variant="ghost"
                              onClick={() => removeParticipant(user)}
                            >
                              <FaTimes />
                            </IconButton>
                          </Table.Cell>
                        </Table.Row>
                      ))}
                    </Table.Body>
                  </Table.Root>

                  <Stack gap={2} mt={4}>
                    <Text fontWeight={"medium"}>Add participant</Text>
                    <Input
                      placeholder="Search by name, ISU, phone, @username or Telegram ID"
                      value={participantSearchInput}
                      onChange={(e) => {
                        setParticipantSearchInput(e.currentTarget.value);
                        setParticipantSearch(e.currentTarget.value);
                      }}
                    />
                    {participantSearchResults.length > 0 ? (
                      <Stack gap={1}>
                        {participantSearchResults.map((user) => (
                          <Flex
                            key={user.id}
                            justify={"space-between"}
                            align={"center"}
                            gap={2}
                            px={2}
                            py={1}
                            borderRadius={"md"}
                            cursor={"pointer"}
                            _hover={{ bg: "bg.muted" }}
                            onClick={() => addParticipant(user)}
                          >
                            <Stack gap={0}>
                              <Text>
                                {user.full_name ||
                                  (user.user_name
                                    ? `@${user.user_name}`
                                    : `Telegram ID ${user.user_tg_id}`)}
                              </Text>
                              <Text fontSize={"sm"} color={"fg.muted"}>
                                {[
                                  user.isu,
                                  user.phone_number,
                                  user.user_name ? `@${user.user_name}` : "",
                                ]
                                  .filter(Boolean)
                                  .join(" · ")}
                              </Text>
                            </Stack>
                          </Flex>
                        ))}
                      </Stack>
                    ) : participantSearchInput.trim().length > 0 ? (
                      <Text color={"fg.muted"}>No users found</Text>
                    ) : null}
                  </Stack>
                </Dialog.Body>
                <Dialog.CloseTrigger asChild>
                  <CloseButton size="sm" />
                </Dialog.CloseTrigger>
              </Dialog.Content>
            </Dialog.Positioner>
          </Portal>
        </Dialog.Root>
      </Card.Footer>
    </Card.Root>
  );
};

type Inputs = {
  title: string;
  date_meeting: Date;
  guest_registration_until?: Date;
  description: string;
  location: string;
  send_images: FileList;
  status: boolean;
};

export const EventsTab = () => {
  const api = useAPI();
  const [events, setEvents] = useState<Activity[]>([]);
  const [allUsers, setAllUsers] = useState<User[]>([]);
  const [open, setOpen] = useState(false);
  const [searchInput, setSearchInput] = useState("");
  const [search, setSearch] = useDebounceValue("", 500);
  const { register, handleSubmit, reset, control } = useForm<Inputs>();

  const createEvent: SubmitHandler<Inputs> = async (data, event) => {
    try {
      await api.activities.create(data);
      toaster.success({
        description: "Event successfully created!",
      });
      reset();
      setOpen(false);
      loadEvents();
    } catch (e) {
      handleError(e);
      event?.stopPropagation();
    }
  };

  const filteredEvents = useMemo(() => {
    const needle = search.trim().toLowerCase();
    if (needle.length === 0) return events;
    return events.filter((event) =>
      event.title.toLowerCase().includes(needle),
    );
  }, [events, search]);

  const upcoming = useMemo(
    () =>
      sortByDateDesc(
        filteredEvents.filter(
          (event) => isFuture(event.date_meeting) && event.status,
        ),
      ),
    [filteredEvents],
  );
  const past = useMemo(
    () =>
      sortByDateDesc(
        filteredEvents.filter(
          (event) => isPast(event.date_meeting) && event.status,
        ),
      ),
    [filteredEvents],
  );
  const inactive = useMemo(
    () => sortByDateDesc(filteredEvents.filter((event) => !event.status)),
    [filteredEvents],
  );

  const loadEvents = async () => {
    try {
      setEvents(await api.activities.getAll());
    } catch (e) {
      handleError(e);
    }
  };

  const loadUsers = async () => {
    try {
      setAllUsers(await api.users.getAll());
    } catch (e) {
      handleError(e);
    }
  };

  useEffect(() => {
    loadEvents();
    loadUsers();
  }, []);

  return (
    <Container maxW={"lg"}>
      <Stack>
        <Dialog.Root open={open} onOpenChange={(e) => setOpen(e.open)}>
          <Dialog.Trigger asChild>
            <Button>New event</Button>
          </Dialog.Trigger>
          <Portal>
            <Dialog.Backdrop />
            <Dialog.Positioner colorPalette={"orange"}>
              <Dialog.Content>
                <Dialog.Header>
                  <Dialog.Title>Create new event</Dialog.Title>
                </Dialog.Header>
                <Dialog.Body as={"form"} onSubmit={handleSubmit(createEvent)}>
                  <Fieldset.Root>
                    <Fieldset.Content>
                      <Field.Root>
                        <Field.Label>Event title</Field.Label>
                        <Input
                          placeholder="Enter event title"
                          {...register("title")}
                        />
                      </Field.Root>
                      <Field.Root>
                        <Field.Label>Date</Field.Label>
                        <Controller
                          control={control}
                          name="date_meeting"
                          render={({ field }) => <Calendar {...field} />}
                        />
                      </Field.Root>
                      <Field.Root>
                        <Field.Label>Guest registration</Field.Label>
                        <Controller
                          control={control}
                          name="guest_registration_until"
                          render={({ field }) => <Calendar {...field} />}
                        />
                      </Field.Root>
                      <Field.Root>
                        <Field.Label>Location</Field.Label>
                        <Input
                          placeholder="Enter location"
                          {...register("location")}
                        />
                      </Field.Root>
                      <Field.Root>
                        <Field.Label>Description</Field.Label>
                        <Textarea
                          placeholder="Enter event description"
                          {...register("description")}
                        />
                      </Field.Root>
                      <FileUpload.Root
                        maxFiles={5}
                        accept={"image/*"}
                        {...register("send_images")}
                      >
                        <FileUpload.HiddenInput />
                        <FileUpload.Trigger asChild>
                          <Button variant="outline" w="full">
                            Upload images for event
                          </Button>
                        </FileUpload.Trigger>
                        <FileUpload.List showSize clearable />
                      </FileUpload.Root>
                    </Fieldset.Content>
                    <Button type={"submit"}>Create event</Button>
                  </Fieldset.Root>
                </Dialog.Body>
                <Dialog.CloseTrigger asChild>
                  <CloseButton size="sm" />
                </Dialog.CloseTrigger>
              </Dialog.Content>
            </Dialog.Positioner>
          </Portal>
        </Dialog.Root>
        <Input
          placeholder="Search by event title"
          value={searchInput}
          onChange={(e) => {
            setSearchInput(e.currentTarget.value);
            setSearch(e.currentTarget.value);
          }}
        />
        <Tabs.Root fitted variant={"enclosed"} defaultValue={"upcoming"}>
          <Tabs.List>
            <Tabs.Trigger value="upcoming">Upcoming</Tabs.Trigger>
            <Tabs.Trigger value="past">Past</Tabs.Trigger>
            <Tabs.Trigger value="inactive">Inactive</Tabs.Trigger>
          </Tabs.List>
          <Tabs.Content value="upcoming">
            <PaginatedList
              items={upcoming}
              pageSize={PAGE_SIZE}
              render={(event) => (
                <EventCard
                  key={event.id}
                  value={event}
                  reload={loadEvents}
                  allUsers={allUsers}
                />
              )}
            />
          </Tabs.Content>
          <Tabs.Content value="past">
            <PaginatedList
              items={past}
              pageSize={PAGE_SIZE}
              render={(event) => (
                <EventCard
                  key={event.id}
                  value={event}
                  reload={loadEvents}
                  allUsers={allUsers}
                />
              )}
            />
          </Tabs.Content>
          <Tabs.Content value="inactive">
            <PaginatedList
              items={inactive}
              pageSize={PAGE_SIZE}
              render={(event) => (
                <EventCard
                  key={event.id}
                  value={event}
                  reload={loadEvents}
                  allUsers={allUsers}
                />
              )}
            />
          </Tabs.Content>
        </Tabs.Root>
      </Stack>
    </Container>
  );
};
