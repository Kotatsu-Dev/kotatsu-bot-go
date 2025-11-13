import type { Roulette } from "@/api/roulettes";
import { handleError, useAPI } from "../../api/api";
import {
  Button,
  Card,
  CloseButton,
  Container,
  Dialog,
  DownloadTrigger,
  Field,
  Fieldset,
  Flex,
  Heading,
  IconButton,
  Input,
  Portal,
  Stack,
  Tabs,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { isFuture, isPast } from "date-fns";
import { toaster } from "../ui/toaster";
import { Calendar } from "../Calendar";
import { FaDownload } from "react-icons/fa";
import { Workbook } from "exceljs";

const exportExcel = async (roulette: Roulette) => {
  const wb = new Workbook();
  const sheet = wb.addWorksheet("СЗ");
  sheet.addRow(["username", "title", "list"]);

  for (const participant of roulette.participants ?? []) {
    sheet.addRow([
      `@${participant.user_name}`,
      participant.enigmatic_title,
      participant.link_my_anime_list,
    ]);
  }

  const buff = await wb.xlsx.writeBuffer();
  const blob = new Blob([buff], {
    type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  });

  return blob;
};

const RouletteComponent = (props: {
  defaultValue?: Roulette;
  onSave?: () => void;
}) => {
  const api = useAPI();
  const [roulette, setRoulette] = useState<
    Omit<Roulette, "id"> & { id?: number }
  >(
    props.defaultValue ?? {
      created_at: new Date(),
      start_date: new Date(),
      announce_date: new Date(),
      distribution_date: new Date(),
      end_date: new Date(),
      theme: "",
      participants: [],
      distribution: null,
    }
  );

  const save = async () => {
    try {
      if (roulette.id) {
        const id = roulette.id;
        await api.roulettes.update({
          ...roulette,
          id,
        });
        toaster.success({
          description: "Roulette succesfully updated",
        });
      } else {
        await api.roulettes.create(roulette);
        toaster.success({
          description: "Roulette succesfully created",
        });
      }
      props.onSave?.();
    } catch (e) {
      handleError(e);
    }
  };

  return (
    <Card.Root key={roulette.id}>
      <Card.Body>
        <Fieldset.Root>
          <Fieldset.Content>
            <Field.Root>
              <Field.Label>Start date</Field.Label>
              <Calendar
                onChange={(date) => {
                  if (date) {
                    setRoulette((roulette) => ({
                      ...roulette,
                      start_date: date,
                    }));
                  }
                }}
                defaultValue={roulette.start_date}
              />
            </Field.Root>
            <Field.Root>
              <Field.Label>Theme publication date</Field.Label>
              <Calendar
                onChange={(date) => {
                  if (date) {
                    setRoulette((roulette) => ({
                      ...roulette,
                      announce_date: date,
                    }));
                  }
                }}
                defaultValue={roulette.announce_date}
              />
            </Field.Root>
            <Field.Root>
              <Field.Label>Distribution date</Field.Label>
              <Calendar
                onChange={(date) => {
                  if (date) {
                    setRoulette((roulette) => ({
                      ...roulette,
                      distribution_date: date,
                    }));
                  }
                }}
                defaultValue={roulette.distribution_date}
              />
            </Field.Root>
            <Field.Root>
              <Field.Label>End date</Field.Label>
              <Calendar
                onChange={(date) => {
                  if (date) {
                    setRoulette((roulette) => ({
                      ...roulette,
                      end_date: date,
                    }));
                  }
                }}
                defaultValue={roulette.end_date}
              />
            </Field.Root>
            <Field.Root>
              <Field.Label>Theme</Field.Label>
              <Input
                onChange={(e) =>
                  setRoulette((roulette) => ({
                    ...roulette,
                    theme: e.target.value,
                  }))
                }
                defaultValue={roulette.theme}
              />
            </Field.Root>
          </Fieldset.Content>
        </Fieldset.Root>
      </Card.Body>
      <Card.Footer>
        {roulette.id && (
          <>
            <DownloadTrigger
              data={() => exportExcel(roulette as Roulette)}
              fileName={`Roulette-${roulette.id!}.xlsx`}
              mimeType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
              asChild
            >
              <IconButton aria-label="Download signed up" variant={"outline"}>
                <FaDownload />
              </IconButton>
            </DownloadTrigger>
            <Button colorPalette={"red"} disabled>
              Delete
            </Button>
          </>
        )}
        <Button colorPalette={"green"} onClick={save}>
          Save
        </Button>
      </Card.Footer>
    </Card.Root>
  );
};

export const RouletteTab = () => {
  const api = useAPI();
  const [roulettes, setRoulettes] = useState<Roulette[]>([]);
  const [open, setOpen] = useState(false);

  const active = roulettes.filter(
    (r) => isPast(r.start_date) && isFuture(r.end_date)
  );
  const past = roulettes.filter((r) => isPast(r.end_date));
  const upcoming = roulettes.filter((r) => isFuture(r.start_date));

  const loadRoulettes = async () => {
    try {
      setRoulettes(await api.roulettes.getAll());
    } catch (e) {
      handleError(e);
    }
  };

  useEffect(() => {
    loadRoulettes();
  }, []);

  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>Roulettes management</Heading>
        <Flex gap={"2"}>
          <Button
            colorPalette={"green"}
            flexGrow={1}
            onClick={() => setOpen(true)}
          >
            Create anime roulette
          </Button>
        </Flex>
        <Tabs.Root variant={"enclosed"} defaultValue={"active"} fitted>
          <Tabs.List>
            <Tabs.Trigger value="active">Active</Tabs.Trigger>
            <Tabs.Trigger value="past">Past</Tabs.Trigger>
            <Tabs.Trigger value="upcoming">Upcoming</Tabs.Trigger>
          </Tabs.List>
          <Tabs.Content value="active">
            <Stack>
              {active.map((roulette) => (
                <RouletteComponent defaultValue={roulette} key={roulette.id} />
              ))}
            </Stack>
          </Tabs.Content>
          <Tabs.Content value="past">
            <Stack>
              {past.map((roulette) => (
                <RouletteComponent defaultValue={roulette} key={roulette.id} />
              ))}
            </Stack>
          </Tabs.Content>
          <Tabs.Content value="upcoming">
            <Stack>
              {upcoming.map((roulette) => (
                <RouletteComponent defaultValue={roulette} key={roulette.id} />
              ))}
            </Stack>
          </Tabs.Content>
        </Tabs.Root>
      </Stack>
      <Dialog.Root open={open} onOpenChange={(e) => setOpen(e.open)}>
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content colorPalette={"orange"}>
              <Dialog.Header>
                <Dialog.Title>Create new roulette</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <RouletteComponent
                  onSave={() => {
                    setOpen(false);
                    loadRoulettes();
                  }}
                />
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
