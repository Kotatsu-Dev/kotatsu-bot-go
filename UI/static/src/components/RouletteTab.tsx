import type { Roulette } from "@/api/roulettes";
import { useAPI } from "../api/api";
import {
  Button,
  Card,
  Container,
  Field,
  Fieldset,
  Heading,
  Input,
  Separator,
  Stack,
  Tabs,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { format, isFuture, isPast, parse } from "date-fns";
import z from "zod";

const RouletteComponent = (props: { defaultValue: Roulette }) => {
  const api = useAPI();
  const [roulette, setRoulette] = useState(props.defaultValue);
  const toValue = (d: Date) => format(d, "yyyy-MM-dd HH:mm");

  const save = async () => {}

  return (
    <Card.Root key={roulette.id}>
      <Card.Body>
        <Fieldset.Root>
          <Fieldset.Content>
            <Field.Root>
              <Field.Label>Start date</Field.Label>
              <Input
                type="datetime-local"
                onChange={(e) =>
                  setRoulette((roulette) => ({
                    ...roulette,
                    start_date: parse(
                      e.target.value,
                      "yyyy-MM-dd'T'HH:mm",
                      new Date()
                    ),
                  }))
                }
                defaultValue={toValue(props.defaultValue.start_date)}
              />
            </Field.Root>
            <Field.Root>
              <Field.Label>Theme publication date</Field.Label>
              <Input
                type="datetime-local"
                onChange={(e) =>
                  setRoulette((roulette) => ({
                    ...roulette,
                    announce_date: parse(
                      e.target.value,
                      "yyyy-MM-dd'T'HH:mm",
                      new Date()
                    ),
                  }))
                }
                defaultValue={toValue(props.defaultValue.announce_date)}
              />
            </Field.Root>
            <Field.Root>
              <Field.Label>Distribution date</Field.Label>
              <Input
                type="datetime-local"
                onChange={(e) =>
                  setRoulette((roulette) => ({
                    ...roulette,
                    distribution_date: parse(
                      e.target.value,
                      "yyyy-MM-dd'T'HH:mm",
                      new Date()
                    ),
                  }))
                }
                defaultValue={toValue(props.defaultValue.distribution_date)}
              />
            </Field.Root>
            <Field.Root>
              <Field.Label>End date</Field.Label>
              <Input
                type="datetime-local"
                onChange={(e) =>
                  setRoulette((roulette) => ({
                    ...roulette,
                    end_date: parse(
                      e.target.value,
                      "yyyy-MM-dd'T'HH:mm",
                      new Date()
                    ),
                  }))
                }
                defaultValue={toValue(props.defaultValue.end_date)}
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
                defaultValue={props.defaultValue.theme}
              />
            </Field.Root>
          </Fieldset.Content>
        </Fieldset.Root>
      </Card.Body>
      <Card.Footer>
        <Button colorPalette={"red"} disabled>
          Delete
        </Button>
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

  const active = roulettes.filter(
    (r) => isPast(r.start_date) && isFuture(r.end_date)
  );

  useEffect(() => {
    api.roulettes
      .getAll()
      .then(setRoulettes)
      .catch((e) => console.log(z.prettifyError(e)));
  }, []);

  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>Roulettes management</Heading>
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
          <Tabs.Content value="past">Past</Tabs.Content>
          <Tabs.Content value="upcoming">Upcoming</Tabs.Content>
        </Tabs.Root>
        <Button colorPalette={"green"}>Create anime roulette</Button>
        <Button colorPalette={"red"}>Delete anime roulette</Button>
        <Button>Anime roulette stats</Button>
        <Fieldset.Root>
          <Fieldset.Content>
            <Field.Root>
              <Field.Label>Start date</Field.Label>
              <Input type="datetime-local" />
            </Field.Root>
            <Field.Root>
              <Field.Label>Theme publication date</Field.Label>
              <Input type="datetime-local" />
            </Field.Root>
            <Field.Root>
              <Field.Label>Distribution date</Field.Label>
              <Input type="datetime-local" />
            </Field.Root>
            <Field.Root>
              <Field.Label>End date</Field.Label>
              <Input type="datetime-local" />
            </Field.Root>
          </Fieldset.Content>
        </Fieldset.Root>
        <Separator />
        <Heading textAlign={"center"}>Set roulette theme</Heading>
        <Input placeholder="Write theme" />
      </Stack>
    </Container>
  );
};
