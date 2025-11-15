import { handleError, useAPI } from "../../api/api";
import {
  Button,
  Card,
  CloseButton,
  Container,
  DataList,
  Dialog,
  Heading,
  Link,
  Portal,
  Stack,
  Status,
  Table,
} from "@chakra-ui/react";
import { toaster } from "../ui/toaster";
import { useState } from "react";
import type { User } from "@/api/users";
import type { Activity } from "@/api/activities";

import { Workbook } from "exceljs";

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
    try {
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
    } catch (e) {
      handleError(e);
    }
  };

  const loadNewsletterSubscribers = async () => {
    try {
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
    } catch (e) {
      handleError(e);
    }
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
      handleError(e);
    }
  };

  const deleteEvent = async (event: Activity) => {
    try {
      await api.activities.setStatus({ id: event.id, status: false });
      await loadEvents();
    } catch (e) {
      handleError(e);
    }
  };

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
                          <Link target="_blank" href={`https://t.me/${user.user_name}`}>
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
                        <Button
                          variant={"outline"}
                          onClick={() => downloadExcel(event)}
                        >
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
