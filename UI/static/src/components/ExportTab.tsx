import { handleError, useAPI } from "../api/api";
import { exportJSON } from "../misc";
import { Button, Container, Heading, Stack } from "@chakra-ui/react";

export const ExportTab = () => {
  const api = useAPI();

  const allUsers = async () => {
    try {
      const users = await api.users.getAll();
      exportJSON(users);
    } catch (e) {
      handleError(e);
    }
  };

  const allEvents = async () => {
    try {
      const events = await api.activities.getAll();
      exportJSON(events);
    } catch (e) {
      handleError(e);
    }
  };

  const allRequests = async () => {
    try {
      const requests = await api.requests.getAll();
      exportJSON(requests);
    } catch (e) {
      handleError(e);
    }
  };

  const allRoulettes = async () => {
    try {
      const roulettes = await api.roulettes.getAll();
      exportJSON(roulettes);
    } catch (e) {
      handleError(e);
    }
  };

  const allClubMembers = async () => {
    try {
      const users = await api.users.getAll();
      const clubMembers = users.filter((user) => user.is_club_member);
      exportJSON(clubMembers);
    } catch (e) {
      handleError(e);
    }
  };

  const allNonClubMembers = async () => {
    try {
      const users = await api.users.getAll();
      const clubMembers = users.filter((user) => !user.is_club_member);
      exportJSON(clubMembers);
    } catch (e) {
      handleError(e);
    }
  };

  return (
    <Container maxW="lg">
      <Stack>
        <Heading textAlign={"center"}>Export data in JSON</Heading>
        <Button onClick={allUsers}>Get all users JSON</Button>
        <Button onClick={allEvents}>Get all events JSON</Button>
        <Button onClick={allRoulettes}>Get all roulettes JSON</Button>
        <Button onClick={allRequests}>Get all requests JSON</Button>
        <Button onClick={allClubMembers}>Get all club members JSON</Button>
        <Button onClick={allNonClubMembers}>
          Get all non club members JSON
        </Button>
      </Stack>
    </Container>
  );
};
