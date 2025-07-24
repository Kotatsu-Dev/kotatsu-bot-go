import {
  Button,
  Center,
  Container,
  Field,
  Fieldset,
  FileUpload,
  Heading,
  Input,
  Textarea,
} from "@chakra-ui/react";

export const EventsTab = () => {
  return (
    <Container>
      <Center>
        <Fieldset.Root maxW="lg">
          <Fieldset.Legend>
            <Heading>Create new event</Heading>
          </Fieldset.Legend>
          <Fieldset.Content>
            <Field.Root>
              <Field.Label>Event title</Field.Label>
              <Input placeholder="Enter event title" />
            </Field.Root>
            <Field.Root>
              <Field.Label>Date</Field.Label>
              <Input type={"datetime-local"} />
              {/* TODO: Dayzed/react-datepicker + chakra or https://github.com/hiwllc/datepicker */}
            </Field.Root>
            <Field.Root>
              <Field.Label>Location</Field.Label>
              <Input placeholder="Enter location" />
            </Field.Root>
            <Field.Root>
              <Field.Label>Description</Field.Label>
              <Textarea placeholder="Enter event description" />
            </Field.Root>
            <FileUpload.Root maxFiles={5}>
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
      </Center>
    </Container>
  );
};
