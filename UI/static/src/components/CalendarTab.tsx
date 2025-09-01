import { useState, type FormEventHandler } from "react";
import { handleError, useAPI } from "../api/api";
import {
  Button,
  Container,
  FileUpload,
  Heading,
  Image,
  Stack,
} from "@chakra-ui/react";
import { toaster } from "./ui/toaster";

export const CalendarTab = () => {
  const api = useAPI();
  const [file, setFile] = useState<File | null>(null);
  const [url, setUrl] = useState("");

  const onChange: FormEventHandler<HTMLDivElement> = (e) => {
    const files = (e.target as HTMLInputElement).files;
    if (files && files.length > 0) {
      const file = files.item(0)!;
      setFile(file);
      setUrl(URL.createObjectURL(file));
    } else {
      setFile(null);
      setUrl("");
    }
  };

  const uploadCalendar = async () => {
    try {
      if (file) {
        await api.calendar.upload({ image: file });
        toaster.success({
          description: "Calendar image succesfully loaded!",
        });
      } else {
        toaster.error({
          description: "No calendar image provided",
        });
      }
    } catch (e) {
      handleError(e);
    }
  };

  return (
    <Container maxW="lg">
      <Stack>
        <Heading textAlign={"center"}>
          Load image for calendar of events
        </Heading>
        <Image
          src={url ? url : api.calendar.imageUrl()}
          aspectRatio={1 / 1}
          fit={"contain"}
        />
        <FileUpload.Root accept={"image/*"} onChange={onChange}>
          <FileUpload.HiddenInput />
          <FileUpload.Trigger asChild>
            <Button variant={"outline"} w={"full"}>
              Select image
            </Button>
          </FileUpload.Trigger>
        </FileUpload.Root>
        <Button disabled={!file} onClick={uploadCalendar}>
          Upload calendar
        </Button>
      </Stack>
    </Container>
  );
};
