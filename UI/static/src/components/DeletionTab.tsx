import { Button, Container, Heading, Stack } from "@chakra-ui/react";

export const DeletionTab = () => {
  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>Delete data</Heading>
        <Button colorPalette={"red"}>WIPE DB</Button>
        <Button colorPalette={"red"}>Delete all users</Button>
        <Button colorPalette={"red"}>Delete all events</Button>
        <Button colorPalette={"red"}>Delete all roulettes</Button>
      </Stack>
    </Container>
  );
};
