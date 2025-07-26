import { Button, Container, Heading, Stack } from "@chakra-ui/react";

export const ExportTab = () => {
  return (
    <Container maxW="lg">
      <Stack>
        <Heading textAlign={"center"}>Export data in JSON</Heading>
        <Button>Get all users JSON</Button>
        <Button>Get all events JSON</Button>
        <Button>Get all roulettes JSON</Button>
        <Button>Get all requests JSON</Button>
        <Button>Get all club members JSON</Button>
        <Button>Get all non club members JSON</Button>
      </Stack>
    </Container>
  );
};
