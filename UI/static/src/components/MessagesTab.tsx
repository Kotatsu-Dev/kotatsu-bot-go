import { useAPI } from "../api/api";
import {
  Button,
  Container,
  Field,
  FileUpload,
  Heading,
  Stack,
  Textarea,
} from "@chakra-ui/react";
import { useForm, type SubmitHandler } from "react-hook-form";
import { toaster } from "./ui/toaster";

type Inputs = {
  message: string;
  files: FileList;
};

export const MessagesTab = () => {
  const api = useAPI();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<Inputs>();

  const sendBroadcast: SubmitHandler<Inputs> = async (e) => {
    await api.broadcast.send(e);
    toaster.success({
      description: "Message broadcasted successfully",
    });
    reset();
  };

  return (
    <Container maxW="lg">
      <Stack as={"form"} onSubmit={handleSubmit(sendBroadcast)}>
        <Heading textAlign={"center"}>Message broadcasting</Heading>
        <Field.Root invalid={!!errors.message}>
          <Field.Label>
            Message <Field.RequiredIndicator />
          </Field.Label>
          <Textarea
            {...register("message", { required: "Message is required" })}
          />
          <Field.ErrorText>{errors.message?.message}</Field.ErrorText>
        </Field.Root>
        <FileUpload.Root maxFiles={5} accept={"image/*"} {...register("files")}>
          <FileUpload.HiddenInput />
          <FileUpload.Trigger asChild>
            <Button variant={"outline"} w={"full"}>
              Select images
            </Button>
          </FileUpload.Trigger>
          <FileUpload.List />
        </FileUpload.Root>
        <Button type="submit">Send broadcast</Button>
      </Stack>
    </Container>
  );
};
