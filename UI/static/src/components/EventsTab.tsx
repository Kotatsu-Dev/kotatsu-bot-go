import { handleError, useAPI } from "../api/api";
import {
  Button,
  Card,
  Container,
  DataList,
  Field,
  Fieldset,
  FileUpload,
  Heading,
  Input,
  Stack,
  Status,
  Tabs,
  Textarea,
} from "@chakra-ui/react";
import { useForm, type SubmitHandler } from "react-hook-form";
import { toaster } from "./ui/toaster";
import { type Activity } from "../api/activities";
import { useEffect, useState } from "react";
import { isFuture, isPast } from "date-fns";
import { Workbook } from "exceljs";

const downloadExcel = async (event: Activity) => {
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

  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "СЗ.xlsx";
  document.body.appendChild(a);
  a.style = "display: none";
  a.click();
  a.remove();
};

const EventCard = (props: { value: Activity; reload: () => void }) => {
  const api = useAPI();
  const event = props.value;

  const deleteEvent = async (event: Activity) => {
    try {
      await api.activities.setStatus({ id: event.id, status: false });
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
                  Active
                </Status.Root>
              ) : (
                <Status.Root colorPalette={"red"}>
                  <Status.Indicator />
                  Inactive
                </Status.Root>
              )}
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
            <DataList.ItemLabel>Location</DataList.ItemLabel>
            <DataList.ItemValue>{event.location}</DataList.ItemValue>
          </DataList.Item>
        </DataList.Root>
      </Card.Body>
      <Card.Footer>
        <Button variant={"outline"} onClick={() => downloadExcel(event)}>
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
  );
};

type Inputs = {
  title: string;
  date_meeting: Date;
  description: string;
  location: string;
  send_images: FileList;
};

export const EventsTab = () => {
  const api = useAPI();
  const [events, setEvents] = useState<Activity[]>([]);
  const { register, handleSubmit, reset } = useForm<Inputs>();

  const createEvent: SubmitHandler<Inputs> = async (e) => {
    try {
      await api.activities.create(e);
      toaster.success({
        description: "Event successfully created!",
      });
      reset();
    } catch (e) {
      handleError(e);
    }
  };

  const upcoming = events.filter(
    (event) => isFuture(event.date_meeting) && event.status
  );
  const past = events.filter(
    (event) => isPast(event.date_meeting) && event.status
  );
  const inactive = events.filter((event) => !event.status);

  const loadEvents = async () => {
    try {
      setEvents(await api.activities.getAll());
    } catch (e) {
      handleError(e);
    }
  };

  useEffect(() => {
    loadEvents();
  }, []);

  return (
    <Container maxW={"lg"} as={"form"} onSubmit={handleSubmit(createEvent)}>
      <Fieldset.Root maxW="lg">
        <Fieldset.Legend>
          <Heading>Create new event</Heading>
        </Fieldset.Legend>
        <Fieldset.Content>
          <Field.Root>
            <Field.Label>Event title</Field.Label>
            <Input placeholder="Enter event title" {...register("title")} />
          </Field.Root>
          <Field.Root>
            <Field.Label>Date</Field.Label>
            <Input type={"datetime-local"} {...register("date_meeting")} />
            {/* TODO: Dayzed/react-datepicker + chakra or https://github.com/hiwllc/datepicker */}
          </Field.Root>
          <Field.Root>
            <Field.Label>Location</Field.Label>
            <Input placeholder="Enter location" {...register("location")} />
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
      <Tabs.Root fitted variant={"enclosed"} defaultValue={"all"}>
        <Tabs.List>
          <Tabs.Trigger value="upcoming">Upcoming</Tabs.Trigger>
          <Tabs.Trigger value="past">Past</Tabs.Trigger>
          <Tabs.Trigger value="inactive">Inactive</Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content value="upcoming">
          <Stack>
            {upcoming.map((event) => (
              <EventCard key={event.id} value={event} reload={loadEvents} />
            ))}
          </Stack>
        </Tabs.Content>
        <Tabs.Content value="past">
          <Stack>
            {past.map((event) => (
              <EventCard key={event.id} value={event} reload={loadEvents} />
            ))}
          </Stack>
        </Tabs.Content>
        <Tabs.Content value="inactive">
          <Stack>
            {inactive.map((event) => (
              <EventCard key={event.id} value={event} reload={loadEvents} />
            ))}
          </Stack>
        </Tabs.Content>
      </Tabs.Root>
    </Container>
  );
};
