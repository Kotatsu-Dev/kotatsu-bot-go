import { useAPI } from "../api/api";
import { Button, Container, Heading, Stack } from "@chakra-ui/react";
import { toaster } from "./ui/toaster";

export const DeletionTab = () => {
  const api = useAPI();

  const wipeDb = async () => {
    await api.db.wipe();
    toaster.success({
      description: "DB succefully deleted",
    });
  };

  const wipeUsers = async () => {
    await api.users.wipe();
    toaster.success({
      description: "All users succefully deleted",
    });
  };

  const wipeEvents = async () => {
    await api.activities.wipe();
    toaster.success({
      description: "All events succefully deleted",
    });
  };

  const wipeRoulettes = async () => {
    await api.roulettes.wipe();
    toaster.success({
      description: "All roulettes succefully deleted",
    });
  };

  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>Delete data</Heading>
        <Button colorPalette={"red"} onClick={wipeDb}>
          WIPE DB
        </Button>
        <Button colorPalette={"red"} onClick={wipeUsers}>
          Delete all users
        </Button>
        <Button colorPalette={"red"} onClick={wipeEvents}>
          Delete all events
        </Button>
        <Button colorPalette={"red"} onClick={wipeRoulettes}>
          Delete all roulettes
        </Button>
      </Stack>
    </Container>
  );
};
