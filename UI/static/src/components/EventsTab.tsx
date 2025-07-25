import { useAPI } from "../api/api";
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
import { useForm, type SubmitHandler } from "react-hook-form";
import { toaster } from "./ui/toaster";

type Inputs = {
  title: string;
  date_meeting: Date;
  description: string;
  location: string;
  send_images: FileList;
};

export const EventsTab = () => {
  const api = useAPI();
  const { register, handleSubmit, reset } = useForm<Inputs>();

  const createEvent: SubmitHandler<Inputs> = async (e) => {
    await api.activities.create(e);
    toaster.success({
      description: "Event successfully created!",
    });
    reset();
  };

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
    </Container>
  );
};
