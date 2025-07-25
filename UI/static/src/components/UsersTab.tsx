import { Button, Container, Heading, Input, Stack } from "@chakra-ui/react";

export const UsersTab = () => {
  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>User management</Heading>
        <Input placeholder="Enter user's Telegram ID" />
        <Stack direction={"row"}>
          <Button variant={"outline"} colorPalette={"red"} flexGrow={1}>
            Kick user from club
          </Button>
          <Button colorPalette={"green"} flexGrow={1}>
            Add user to club
          </Button>
        </Stack>
      </Stack>
    </Container>
  );
};
