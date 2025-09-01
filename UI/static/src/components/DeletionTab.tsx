import { handleError, useAPI } from "../api/api";
import { Button, Container, Heading, Stack } from "@chakra-ui/react";
import { toaster } from "./ui/toaster";

export const DeletionTab = () => {
  const api = useAPI();

  const wipeDb = async () => {
    try {
      await api.db.wipe();
      toaster.success({
        description: "DB succefully deleted",
      });
    } catch (e) {
      handleError(e);
    }
  };

  const wipeUsers = async () => {
    try {
      await api.users.wipe();
      toaster.success({
        description: "All users succefully deleted",
      });
    } catch (e) {
      handleError(e);
    }
  };

  const wipeEvents = async () => {
    try {
      await api.activities.wipe();
      toaster.success({
        description: "All events succefully deleted",
      });
    } catch (e) {
      handleError(e);
    }
  };

  const wipeRoulettes = async () => {
    try {
      await api.roulettes.wipe();
      toaster.success({
        description: "All roulettes succefully deleted",
      });
    } catch (e) {
      handleError(e);
    }
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
