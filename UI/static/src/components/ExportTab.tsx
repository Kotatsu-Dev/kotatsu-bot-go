import { useAPI } from "../api/api";
import { exportJSON } from "../misc";
import { Button, Container, Heading, Stack } from "@chakra-ui/react";

export const ExportTab = () => {
  const api = useAPI();

  const allUsers = async () => {
    const users = await api.users.getAll();
    exportJSON(users);
  };

  const allEvents = async () => {
    const events = await api.activities.getAll();
    exportJSON(events);
  };

  const allRequests = async () => {
    const requests = await api.requests.getAll();
    exportJSON(requests);
  };

  const allClubMembers = async () => {
    const users = await api.users.getAll();
    const clubMembers = users.filter((user) => user.is_club_member);
    exportJSON(clubMembers);
  };

  const allNonClubMembers = async () => {
    const users = await api.users.getAll();
    const clubMembers = users.filter((user) => !user.is_club_member);
    exportJSON(clubMembers);
  };

  return (
    <Container maxW="lg">
      <Stack>
        <Heading textAlign={"center"}>Export data in JSON</Heading>
        <Button onClick={allUsers}>Get all users JSON</Button>
        <Button onClick={allEvents}>Get all events JSON</Button>
        {/* TODO */}
        <Button disabled>Get all roulettes JSON</Button>
        <Button onClick={allRequests}>Get all requests JSON</Button>
        <Button onClick={allClubMembers}>Get all club members JSON</Button>
        <Button onClick={allNonClubMembers}>
          Get all non club members JSON
        </Button>
      </Stack>
    </Container>
  );
};
