import { useState, type FormEventHandler } from "react";
import { useAPI } from "../api/api";
import { Button, Container, FileUpload, Image, Stack } from "@chakra-ui/react";

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

  return (
    <Container maxW="lg">
      <Stack>
        <Image
          src={url ? url : `${api.getRoot()}/get-calendar-file`}
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
        <Button>Upload calendar</Button>
      </Stack>
    </Container>
  );
};
