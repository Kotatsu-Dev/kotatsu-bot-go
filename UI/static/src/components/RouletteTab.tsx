import {
  Button,
  Container,
  Field,
  Fieldset,
  Heading,
  Input,
  Separator,
  Stack,
} from "@chakra-ui/react";

export const RouletteTab = () => {
  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>Roulettes management</Heading>
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
